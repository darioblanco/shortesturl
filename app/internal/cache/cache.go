package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/darioblanco/shortesturl/app/internal/config"
	"github.com/darioblanco/shortesturl/app/internal/logging"
	"github.com/go-redis/redis/v8"
)

// A Cache exposes functions from an in-memory store
type Cache interface {
	// Get gets the value of the given key.
	// If the key does not exist, the returned string will be empty ("").
	Get(ctx context.Context, key string) (string, error)
	// SetIfNotExists set key to hold string value if key does not exist (returning true).
	// If key exists, no operation is performed and false is returned as a result.
	// An expiration for the key has to be set. Zero expiration means the key is there forever.
	SetIfNotExists(
		ctx context.Context, key string, value string, expiration time.Duration,
	) (bool, error)
}

type cache struct {
	client    *redis.Client
	logger    logging.Logger
	keyPrefix string
}

// New creates a new cache instance
func New(
	ctx context.Context, conf *config.Values, logger logging.Logger,
) (Cache, error) {
	var opts *redis.Options
	if conf.IsDevelopment {
		mr, _ := miniredis.Run()
		opts = &redis.Options{
			Addr: fmt.Sprintf("%s:%s", mr.Host(), mr.Port()),
		}
		logger.Debug("Loaded miniredis configuration")
	} else {
		opts = &redis.Options{
			Addr: fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort),
		}
		logger.Debug("Loaded production redis configuration")
	}
	client := redis.NewClient(opts)
	// If redis is not available, it will hang here until it times out
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	logger.Info("Connected to cache", "address", client.Options().Addr)
	return &cache{
		client: client,
		logger: logger,
	}, nil
}

func (c cache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(
		ctx,
		key,
	).Result()
	if err == redis.Nil {
		// The given key does not exist, in our architecture, each key MUST always have a url
		return "", nil
	}
	return value, err
}

func (c cache) SetIfNotExists(
	ctx context.Context, key string, value string, expiration time.Duration,
) (bool, error) {
	collisionError := errors.New("url collision")
	txf := func(tx *redis.Tx) error {
		// Check if the key candidate is already stored
		previousValue, err := c.Get(ctx, key)
		if err != nil {
			return err
		}
		if previousValue != "" {
			// The key is already being used
			if value == previousValue {
				// The key was already stored in the service, no need to set it again
				return nil
			}
			// The key is used for a different value: collision
			return collisionError
		}

		// Runs only if the key to be stored is not set by a different process
		// This prevents race conditions, where other processes might store a different
		// value in the same key, which would render the system in a corrupted state
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, value, expiration)
			return err
		})
		return err
	}

	// In case of race condition, the operation will be retried four more times
	// As collisions are very rare between different urls, and encoding the same
	// url multiple times might cause a single race condition but it is unlikely
	// it will trigger more, an optimistic lock approach like this should be
	// sufficient to prevent a corrupted state even in scenarios with high parallelism
	for retries := 0; retries < 5; retries++ {
		err := c.client.Watch(ctx, txf, key)
		if err != redis.TxFailedErr {
			if errors.Is(err, collisionError) {
				return false, nil
			}
			return true, err
		}
		// optimistic lock lost
	}
	return false, errors.New("max retries reached (4)")
}

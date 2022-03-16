package main

import (
	"context"

	"github.com/darioblanco/shortesturl/app"
)

func main() {
	// This part is very simple now, but for a production application it can grow
	// rather quickly. Therefore, creating an application abstraction will allow
	// the creation of multiple configuration parameters (e.g. to execute migrations,
	// to disable different storage backends, etc...) easily
	a := app.New(context.Background())
	a.Serve()
}

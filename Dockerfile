# ------------------------------------------------------
#                       Dockerfile
# ------------------------------------------------------
# image:    shortesturl
# tag:      <COMMIT HASH>
# name:     darioblanco/shortesturl
# requires: golang:1.16
# authors:  dblancoit@gmail.com
# ------------------------------------------------------

# BUILDER - Artifacts build for production
FROM golang:1.17 as builder

WORKDIR /go/src/github.com/darioblanco/shortesturl/
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN make -r build PRODUCTION=1

# RUNNER - Production image
FROM gcr.io/distroless/static-debian11

WORKDIR /bin/
COPY --from=builder /go/src/github.com/darioblanco/shortesturl/tmp/shortesturl .
COPY --from=builder /go/src/github.com/darioblanco/shortesturl/cmd/config.default.yaml ./cmd/
ENTRYPOINT [ "/bin/shortesturl" ]

# Tests

To run the tests (with its respective coverage):

```sh
make test
```

This project aims for a 98% test coverage of all files, with the exception of
`cmd/server/main.go` and `app/app.go`. Those files are considered entry points
of the application, thus they are tested with `make run`.

In addition, `app/internal/cache/cache.go` is
not properly shown in the coverage due to go limitations, as technically the `app/internal/http/api_test.go`
file is using the optimistic lock and eventually could trigger its retry mechanism. Therefore,
a threshold of 98% test coverage has been set instead of 100%. It would be possible to test the
retry and error branch in the `SetIfNotExists` function, but that would require creating a mock
over the cache inherent client instead of using miniredis. I decided not to spend extra time on
that mock wrapper as the benefit for this task is little and I consider that the code is production
ready. However, if the project would grow, I would strongly consider even reaching 100% test coverage
with unit tests there.

All files ending in `_test.go` are performing unit tests, except `api_test.go`
that follows an integration test approach for the `/encode` and `/decode` endpoints.

Files ending in `_mock.go` are thought to expose mocks and stubs to different packages
that might need them.
Therefore, these mocks will be defined within the package of the real implementation,
even if it does not use them. [Mockery](https://github.com/vektra/mockery) is an interesting tool to
automatically generate such mocks, though this project does not implement it as it is not complex enough.
The mock files are automatically ignored from the coverage report.

## Manual tests

You can manually check the endpoints, using the `shortesturl.http` file and
[VisualStudioCode/humao.rest-client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client).

I consider that strong automated testing is the way to go in order to have a
codebase that is maintainable (especially when developing collaboratively), but
manual tests is helpful also during development and for demo purposes.

## Benchmarks

In addition, I have create a few basic benchmarks that would help to analyze the speed of the
encode and decode functions. These are set in `app/internal/http/benchmarks_test.go`.

These are not tests! Therefore, no asserts are done in this file. However, they will be very
helpful to analyze the impact of possible improvements in the algorithms.

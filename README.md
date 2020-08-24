# bazel-recipes

This repo is a collection of recipes that I have found useful in using bazel as the main component of a continuous 
integration build.

Taking on the ideas from the [Continuous Delivery book](https://continuousdelivery.com/), a “commit stage” pipeline
runs before multiple “acceptance test” pipelines.  In the commit stage, all _all_ artifacts are built (to containers
where possible) - this includes services, and all acceptance tests.  All acceptance tests should be baked into
containers with all their required dependencies so that running them in the "acceptance test” pipeline is pure
scheduling.

Bazel is a great fit for the commit stage pipeline in a cloud-native deployment:
* it builds containers natively (“rules_docker”)
* it is polyglot/can build anything (by using/defining the appropriate rule: java, c++, python, golang, docker, dhall,
 clojure)
* the pipelne complexity is in the bazel setup, not in custom bash scripts or CI tool config. This also has the
 desirable consequence that all developers can test the build using just bazel.
* it is fast

The recipes are:

| Recipe                         | Description |
|--------------------------------|-------------|
| [gcp_cloud_function](/gcp_cloud_function/README.md) | Building and testing a GCP cloud function |
| [go_openapi_service](/go_openapi_service/README.md) | Building and testing a OpenAPI 3 service using golang |
| [integration_test_with_services](/integration_test_with_services/README.md) | Building a container that includes dependencies (in this case PostgreSQL) |
| [turps](/turps/README.md)  | Building a golang grpc service, with e2e acceptance tests in a container |

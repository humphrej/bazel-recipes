Go OpenAPI Service
------------------
### Introduction
This recipe demonstrates how to build a simple REST service inside a container using **both** Bazel and the "go" tool
from the same sources.

The code from this example is taken directly from https://goswagger.io/tutorial/todo-list.html

### Context:
* Bazel rules_go https://github.com/bazelbuild/rules_go
* Bazel rules_docker https://github.com/bazelbuild/rules_docker
* Go Swagger https://goswagger.io

### Solution: 
Use Bazel to run the cloud function tests as part of a regular build.

To regenerate the server:

    swagger generate server -A todo-list -f ./swagger.yml

To update pinned dependencies in bazel:

    bazel run :gazelle -- update-repos -from_file=$PWD/go.mod -to_macro=go_dependencies.bzl%go_repositories 

To run the container using bazel (MacOS): 

    bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64  //go_openapi_service/cmd/todo-list-server:todo-list-server

To run the container using bazel (Linux): 

    bazel run //go_openapi_service/cmd/todo-list-server:todo-list-server

#### Solution Limitations:
* When dependencies are changed in go.mod, then the bazel dependencies must be regenerated using the gazelle command above.
* Dependencies that are no longer referenced from go.mod can be removed by appending -prune to the gazelle update-repos command
* go_dependencies.bzl is shared between all go modules in the build

#### Consequences:
* Positive
    * Go based REST service is only 29MB, compared to 200+MB for a Java equivalent
* Negative
    * None

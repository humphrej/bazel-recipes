Turps
------------------------------
### Introduction
This recipe demonstrates
* building a simple grpc service in golang to store test results.
* building acceptance test suite that is built into a container with dependencies

### Context:
* Using the ideas from the Continuous Delivery book, a deployment pipeline should be split into parts:
    * a _commit stage_ that builds software artifacts, **including** acceptance tests
    * an _acceptance test stage_ that runs the test artifacts built in the _commit stage_
* There is a need for a durable store of test results from the _acceptance test stage_ so that ops can make go/no-go
 release decisions
    
### Solution:
* Based heavily on [integraton test with services](../integration_test_with_services).
* The turpsd service consists of a simple grpc controller plus a PostgreSQL storage backend, written in golang.  
* A command line tool, turps, is used to create change lists, attach test results to change lists, and query change lists.

### Development
#### How to build
    bazel build //...
#### How to run the tests
1. Load the acceptance tests into the local docker:

    (MacOS) bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //turps/test/image:test_image
    (Linux) bazel run //turps/test/image:test_image

1. Run the acceptance tests:
    
    docker run -ti bazel/turps/test/image:test_image
    
#### How to update the BUILD files from the golang source
gazelle (see [rules_go](https://github.com/bazelbuild/rules_go)) is used to update the BUILD files with dependencies
based on their usage in the go code. This needs to be run every time any new dependencies are added to a go source file.
Any dependencies must also be defined in go_dependencies.bzl (see
[How to update project dependencies](#How-to-update-project-dependencies))

From project root:

    bazel run :gazelle update -- turps
    
#### How to update project dependencies
gazelle (see [rules_go](https://github.com/bazelbuild/rules_go)) is used pin external golang dependency version based
on the contents of the go.mod.

From project root:

    bazel run :gazelle -- update-repos -from_file=$PWD/turps/go.mod -to_macro=go_dependencies.bzl%go_repositories

#### Solution Limitations:

#### Consequences:
* Positive

* Negative
    * Slow to build
    * gazelle isn't great in a monorepo project
    * Goland sucks a bit

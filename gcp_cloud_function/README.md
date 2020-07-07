GCP Cloud Function
------------------
### Introduction
This recipe demonstrates how to build/test a GCP Cloud function using **both** Bazel and the go tool
from the same sources.

Cloud functions should be tested and built just like any other deployment artifact.

In this example, a Cloud Function listens for notifications from Cloud Build and sends a notification
email containing the build status.  

### Context:
* Cloud Functions https://cloud.google.com/functions
* Cloud Build https://cloud.google.com/cloud-build
* Bazel rules_go https://github.com/bazelbuild/rules_go
* Sendgrid email API https://github.com/sendgrid/sendgrid-go

### Solution: 
Use Bazel to run the cloud function tests as part of a regular build.

To test using go tool:
  
    cd gcp_cloud_function
    go test

To update pinned dependencies in bazel:

    bazel run :gazelle -- update-repos -from_file=$PWD/go.mod -to_macro=go_dependencies.bzl%go_repositories -prune

To run using Bazel:

    bazel test //...

To deploy to gcp:

    gcloud functions deploy NotifyBuild --trigger-topic cloud-builds --runtime go113 --set-env-vars=SENDGRID_API_KEY=<insert>,EMAIL_FROM=bob@foo.com,EMAIL_TO=alice@foo.com

#### Solution Limitations:
When dependencies are changed in go.mod, then the bazel dependencies must be regenerated using the gazelle command above.

#### Consequences:
* Positive
    * consistent testng of cloud functions with other sources
* Negative
    * None

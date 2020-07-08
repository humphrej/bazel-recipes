load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

# gazelle:prefix foo.com/build-notifier

go_library(
    name = "go_default_library",
    srcs = [
        "email.go",
        "notify.go",
    ],
    importpath = "foo.com/build-notifier",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_sendgrid_sendgrid_go//:go_default_library",
        "@com_github_sendgrid_sendgrid_go//helpers/mail:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["notify_test.go"],
    embed = [":go_default_library"],
)

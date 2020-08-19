load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "turps/cmd/client",
    visibility = ["//visibility:private"],
    deps = [
        "//turps/api:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_binary(
    name = "client",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

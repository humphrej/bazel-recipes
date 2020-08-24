load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "turps/cmd/client",
    visibility = ["//visibility:private"],
    deps = [
        "//turps/api:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
    ],
)

go_binary(
    name = "turps",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

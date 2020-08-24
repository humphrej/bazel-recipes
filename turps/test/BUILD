load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "driver_api.go",
        "driver_cli.go",
        "spec.go",
    ],
    importpath = "turps/test",
    visibility = ["//visibility:public"],
    deps = [
        "//turps/api:go_default_library",
        "//turps/internal:go_default_library",
        "//turps/internal/grpc:go_default_library",
        "//turps/internal/postgres:go_default_library",
        "//turps/internal/postgres/testing:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
        "@org_golang_google_protobuf//proto:go_default_library",
    ],
)

go_test(
    name = "go_acceptance_test",
    srcs = ["e2e_test.go"],
    data = [
        "//turps/cmd/client:turps",
    ],
    embed = [":go_default_library"],
    tags = [
        "external",
        "manual",
    ],
    visibility = ["//turps:__subpackages__"],
)

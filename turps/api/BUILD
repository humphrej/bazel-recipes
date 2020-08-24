load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "turps_proto",
    srcs = [
        "turps.proto",
    ],
    deps = ["@com_google_protobuf//:timestamp_proto"],
)

go_proto_library(
    name = "api_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "turps/api",
    proto = ":turps_proto",
)

go_library(
    name = "go_default_library",
    embed = [":api_go_proto"],
    importpath = "turps/api",
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["storage.go"],
    importpath = "turps/internal/postgres",
    visibility = ["//turps:__subpackages__"],
    deps = [
        "//turps/internal:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_jackc_pgx_v4//pgxpool:go_default_library",
    ],
)

go_test(
    name = "go_postgres_storage_test",
    srcs = ["storage_test.go"],
    embed = [":go_default_library"],
    tags = ["manual"],
    visibility = ["//turps:__subpackages__"],
    deps = [
        "//turps/internal:go_default_library",
        "//turps/internal/postgres/testing:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["storage_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//turps/internal:go_default_library",
        "//turps/internal/postgres/testing:go_default_library",
    ],
)

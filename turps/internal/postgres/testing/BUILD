load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["testing.go"],
    importpath = "turps/internal/postgres/testing",
    visibility = ["//turps:__subpackages__"],
    deps = ["@com_github_jackc_pgx_v4//pgxpool:go_default_library"],
)

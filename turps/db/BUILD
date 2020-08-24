package(default_visibility = ["//visibility:public"])

load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")

pkg_tar(
    name = "db_migrations_tar",
    srcs = glob(["**/*.sql"]),
    package_dir = "/db/migrations",
)

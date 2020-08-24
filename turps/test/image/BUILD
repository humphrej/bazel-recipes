load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_docker//container:container.bzl", "container_image")

pkg_tar(
    name = "golang_migrate_tar",
    srcs = ["@golang-migrate//:files"],
)

pkg_tar(
    name = "prestart_scripts_tar",
    srcs = glob(["prestart.d/**"]),
    package_dir = "/etc/prestart.d",
)

container_image(
    name = "test_image",
    testonly = True,
    base = "//integration_test_with_services:postgres_bare_image",
    directory = "/",
    entrypoint = ["/entrypoint.sh"],
    files = [
        "entrypoint.sh",
        "env.sh",
        "workload.sh",
        "//turps/cmd/client:turps",
        "//turps/cmd/turpsd",
        "//turps/internal/postgres:go_postgres_storage_test",
        "//turps/test:go_acceptance_test",
    ],
    tars = [
        ":golang_migrate_tar",
        ":prestart_scripts_tar",
        "//turps/db:db_migrations_tar",
    ],
)

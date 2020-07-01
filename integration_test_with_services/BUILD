package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_docker//docker/package_managers:install_pkgs.bzl", "install_pkgs")
load("@io_bazel_rules_docker//docker/package_managers:download_pkgs.bzl", "download_pkgs")
load("@io_bazel_rules_docker//container:container.bzl", "container_image")
load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")

download_pkgs(
  name = "postgres_pkgs",
  image_tar = "@ubuntu2004//image",
  packages = [
      "postgresql",
      #"openjdk-11-jdk"
  ]
)

install_pkgs(
    name = "postgres_bare_image",
    image_tar = "@ubuntu2004//image",
    installables_tar = ":postgres_pkgs.tar",
    output_image_name = "postgres_bare_image",
)

# define prestart directory as a tar so that bazel will trigger re-build on change
pkg_tar(
    name = "prestart_scripts_tar",
    srcs = glob(["prestart.d/**"]),
    package_dir = "/etc/prestart.d",
)

container_image(
    name = "postgresql_test_image",
    base = ":postgres_bare_image",
    directory = "/",
    tars = [
        ":prestart_scripts_tar",
    ],
    files = [
        "env.sh",
        "workload.sh",
        "entrypoint.sh"
    ],
    entrypoint = [ "/entrypoint.sh" ]
)


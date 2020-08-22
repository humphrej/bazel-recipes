load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "app.go",
        "server.go",
        "transform.go",
    ],
    importpath = "turps/internal/grpc",
    visibility = ["//turps:__subpackages__"],
    deps = [
        "//turps/api:go_default_library",
        "//turps/internal:go_default_library",
        "//turps/internal/postgres:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        "@com_github_jackc_pgx_v4//pgxpool:go_default_library",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "server_test.go",
        "transform_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["small"],
    deps = [
        "//turps/api:go_default_library",
        "//turps/internal:go_default_library",
        "//turps/internal/postgres:go_default_library",
        "//turps/internal/postgres/testing:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        "@com_github_leanovate_gopter//:go_default_library",
        "@com_github_leanovate_gopter//gen:go_default_library",
        "@com_github_leanovate_gopter//prop:go_default_library",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
    ],
)

go_test(
    name = "go_acceptance_test",
    srcs = ["server_test.go"],
    #    tags = ["manual"],
    embed = [":go_default_library"],
    visibility = ["//turps:__subpackages__"],
    deps = [
        "//turps/api:go_default_library",
        "//turps/internal:go_default_library",
        "//turps/internal/postgres:go_default_library",
        "//turps/internal/postgres/testing:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
    ],
)

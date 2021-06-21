load("@rules_java//java:defs.bzl", "java_library")

java_library(
    name = "bazeljunit5",
    srcs = glob(["**/*.java"]),
    visibility = [
        "//visibility:public",
    ],
    deps = [
        "@maven//:com_google_guava_guava",
        "@maven//:org_junit_platform_junit_platform_commons",
        "@maven//:org_junit_platform_junit_platform_console",
        "@maven//:org_junit_platform_junit_platform_reporting",
        #        "@maven//:org_junit_platform_junit_platform_runner",
        #        "@maven//:org_junit_platform_junit_platform_suite_api",
    ],
)

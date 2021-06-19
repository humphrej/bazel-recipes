java_library(
    name = "junit4container_library",
    srcs = glob(["*.java"]),
    deps = [
        "//junit-in-container/src/main/java/junitcontainer:junit_library",
        "//junit-in-container/src/test/java/junitcontainer:junitcontainer_test_library",
        "@maven//:ant_ant_junit",
        "@maven//:com_google_truth_truth",
        "@maven//:junit_junit",
        "@maven//:org_apache_ant_ant",
    ],
)

java_test(
    name = "junit4container",
    test_class = "junitcontainer.junit4.TestSuite",
    runtime_deps = [
        ":junit4container_library",
        "@maven//:junit_junit",
    ],
)

java_binary(
    name = "junit4container_main",
    main_class = "junit4container.RunTests",
    runtime_deps = [
        ":junit4container_library",
    ],
)

load("@io_bazel_rules_docker//java:image.bzl", "java_image")

# Creates a container that runs the tests using the junit launcher
java_image(
    name = "junit4container_launcher_image",
    srcs = None,
    args = [
        "--select-class",
        "junitcontainer.junit4.TestSuite",
    ],
    main_class = "com.flexport.bazeljunit5.BazelJUnit5ConsoleLauncher",
    runtime_deps = [
        ":junit4container_library",
        "//junit-in-container/tools/bazeljunit5/src/main/java/com/flexport/bazeljunit5",
        "@maven//:org_apiguardian_apiguardian_api",
        "@maven//:org_junit_jupiter_junit_jupiter_api",
        "@maven//:org_junit_jupiter_junit_jupiter_params",
        "@maven//:org_junit_platform_junit_platform_launcher",
        "@maven//:org_junit_platform_junit_platform_suite_api",
        "@maven//:org_junit_vintage_junit_vintage_engine",
        "@maven//:org_opentest4j_opentest4j",
    ],
)

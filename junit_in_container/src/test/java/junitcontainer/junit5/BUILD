java_library(
    name = "junit5container_library",
    srcs = glob(["*.java"]),
    deps = [
        "//junit_in_container/src/main/java/junitcontainer:junit_library",
        "//junit_in_container/src/test/java/junitcontainer:junitcontainer_test_library",
        "@maven//:com_google_truth_truth",
        "@maven//:org_junit_jupiter_junit_jupiter_api",
        "@maven//:org_junit_jupiter_junit_jupiter_engine",
        "@maven//:org_junit_jupiter_junit_jupiter_params",
        "@maven//:org_junit_platform_junit_platform_console",
        "@maven//:org_junit_platform_junit_platform_engine",
        "@maven//:org_junit_platform_junit_platform_launcher",
        "@maven//:org_junit_platform_junit_platform_reporting",
    ],
)

load("//:junit_in_container/junit5.bzl", "java_junit5_test")

# Uses java_junit_test to create a test target
java_junit5_test(
    name = "junit5container_using_macro",
    srcs = None,
    main_class = "com.flexport.bazeljunit5.BazelJUnit5ConsoleLauncher",
    test_package = "junit5container",
    runtime_deps = [
        ":junit5container_library",
        "//junit_in_container/tools/bazeljunit5/src/main/java/com/flexport/bazeljunit5",
    ],
)

# Uses java_test to create a test target
java_test(
    name = "junit5container",
    srcs = None,
    args = [
        "--select-package",
        "junit5container",
    ],
    main_class = "com.flexport.bazeljunit5.BazelJUnit5ConsoleLauncher",
    use_testrunner = False,
    runtime_deps = [
        ":junit5container_library",
        "//junit_in_container/src/main/java/junitcontainer:junit_library",
        "//junit_in_container/tools/bazeljunit5/src/main/java/com/flexport/bazeljunit5",
        "@maven//:com_google_truth_truth",
        "@maven//:org_apiguardian_apiguardian_api",
        "@maven//:org_junit_jupiter_junit_jupiter_api",
        "@maven//:org_junit_jupiter_junit_jupiter_engine",
        "@maven//:org_junit_jupiter_junit_jupiter_params",
        "@maven//:org_junit_platform_junit_platform_launcher",
        "@maven//:org_junit_platform_junit_platform_suite_api",
        "@maven//:org_opentest4j_opentest4j",
    ],
)

# Creates a binary target to run the tests programatically
java_binary(
    name = "junit5container_main",
    main_class = "junitcontainer.junit5.RunTests",
    runtime_deps = [
        ":junit5container_library",
    ],
)

load("@io_bazel_rules_docker//java:image.bzl", "java_image")

# Creates a container that runs the tests
java_image(
    name = "junit5container_image",
    main_class = "junit5container.RunTests",
    runtime_deps = [
        ":junit5container_library",
    ],
)

# Creates a container that runs the tests using the junit launcher
java_image(
    name = "junit5container_launcher_image",
    srcs = None,
    args = [
        "--select-package",
        "junitcontainer.junit5",
    ],
    main_class = "com.flexport.bazeljunit5.BazelJUnit5ConsoleLauncher",
    runtime_deps = [
        ":junit5container_library",
        "//junit_in_container/src/main/java/junitcontainer:junit_library",
        "//junit_in_container/tools/bazeljunit5/src/main/java/com/flexport/bazeljunit5",
        "@maven//:com_google_truth_truth",
        "@maven//:org_apiguardian_apiguardian_api",
        "@maven//:org_junit_jupiter_junit_jupiter_api",
        "@maven//:org_junit_jupiter_junit_jupiter_engine",
        "@maven//:org_junit_jupiter_junit_jupiter_params",
        "@maven//:org_junit_platform_junit_platform_launcher",
        "@maven//:org_junit_platform_junit_platform_suite_api",
        "@maven//:org_opentest4j_opentest4j",
    ],
)

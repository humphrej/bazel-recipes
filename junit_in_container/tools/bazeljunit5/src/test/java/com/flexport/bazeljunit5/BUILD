load("//:junit_in_container/junit5.bzl", "java_junit5_test")

filegroup(
    name = "java_test_files",
    srcs = glob(["**/*.java"]),
)

java_junit5_test(
    name = "test",
    srcs = [
        ":java_test_files",
    ],
    jvm_flags = [],
    main_class = "com.flexport.bazeljunit5.BazelJUnit5ConsoleLauncher",
    test_package = "com.flexport.bazeljunit5",
    deps = [
        # Internal dependencies
        "//junit_in_container/tools/bazeljunit5/src/main/java/com/flexport/bazeljunit5",

        # External dependencies
        "@maven//:com_google_truth_truth",
        "@maven//:org_junit_jupiter_junit_jupiter_api",
        "@maven//:org_junit_jupiter_junit_jupiter_params",
    ],
)

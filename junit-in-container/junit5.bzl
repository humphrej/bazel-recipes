"""External dependencies & java_junit5_test rule"""

JUNIT_JUPITER_GROUP_ID = "org.junit.jupiter"
JUNIT_JUPITER_ARTIFACT_ID_LIST = [
    #    "junit-jupiter-api",
    "junit-jupiter-engine",
    "junit-jupiter-params",
]

JUNIT_PLATFORM_GROUP_ID = "org.junit.platform"
JUNIT_PLATFORM_ARTIFACT_ID_LIST = [
    "junit-platform-commons",
    "junit-platform-console",
    "junit-platform-engine",
    "junit-platform-launcher",
    "junit-platform-runner",
    "junit-platform-suite-api",
]

JUNIT_VINTAGE_GROUP_ID = "org.junit.vintage"
JUNIT_VINTAGE_ARTIFACT_ID_LIST = [
    "junit-vintage-engine",
]

JUNIT_EXTRA_DEPENDENCIES = [
    ("org.apiguardian", "apiguardian-api", "1.0.0"),
    ("org.opentest4j", "opentest4j", "1.1.1"),
]

def _format_maven_coordinates_dep_name(group_id, artifact_id, version):
    return "%s:%s:%s" % (group_id, artifact_id, version)

JUNIT_JUPITER_VERSION = "5.7.2"
JUNIT_PLATFORM_VERSION = "1.6.2"

def make_deps():
    JUNIT_JUPITER_DEPS = [
        _format_maven_coordinates_dep_name(JUNIT_JUPITER_GROUP_ID, artifact_id, JUNIT_JUPITER_VERSION)
        for artifact_id in JUNIT_JUPITER_ARTIFACT_ID_LIST
    ]
    JUNIT_PLATFORM_DEPS = [
        _format_maven_coordinates_dep_name(JUNIT_PLATFORM_GROUP_ID, artifact_id, JUNIT_PLATFORM_VERSION)
        for artifact_id in JUNIT_PLATFORM_ARTIFACT_ID_LIST
    ]
    JUNIT_VINTAGE_DEPS = [
        _format_maven_coordinates_dep_name(JUNIT_VINTAGE_GROUP_ID, artifact_id, JUNIT_JUPITER_VERSION)
        for artifact_id in JUNIT_VINTAGE_ARTIFACT_ID_LIST
    ]
    JUNIT_EXTRA_DEPS = [
        _format_maven_coordinates_dep_name(t[0], t[1], t[2])
        for t in JUNIT_EXTRA_DEPENDENCIES
    ]

    return {
        "jupiter": JUNIT_JUPITER_DEPS,
        "platform": JUNIT_PLATFORM_DEPS,
        "vintage": JUNIT_VINTAGE_DEPS,
        "extras": JUNIT_EXTRA_DEPS,
    }

def java_junit5_test(name, srcs, test_package, main_class = None, deps = None, runtime_deps = [], **kwargs):
    FILTER_KWARGS = [
        "main_class",
        "use_testrunner",
        "args",
    ]

    for arg in FILTER_KWARGS:
        if arg in kwargs.keys():
            kwargs.pop(arg)

    junit_console_args = []
    if test_package:
        junit_console_args += ["--select-package", test_package]
    else:
        fail("must specify 'test_package'")

    # add the junit 5 dependencies to runtime_deps
    runtime_deps = runtime_deps + [
        _format_maven_jar_dep_name(JUNIT_JUPITER_GROUP_ID, artifact_id)
        for artifact_id in JUNIT_JUPITER_ARTIFACT_ID_LIST
    ] + [
        _format_maven_jar_dep_name(JUNIT_PLATFORM_GROUP_ID, "junit-platform-suite-api"),
    ] + [
        _format_maven_jar_dep_name(t[0], t[1])
        for t in JUNIT_EXTRA_DEPENDENCIES
    ] + [
        _format_maven_jar_dep_name(JUNIT_PLATFORM_GROUP_ID, artifact_id)
        for artifact_id in JUNIT_PLATFORM_ARTIFACT_ID_LIST
    ]

    # remove any duplicates from runtime_deps
    extra_runtime_deps = {}
    for rd in runtime_deps:
        extra_runtime_deps[rd] = True

    #print(extra_runtime_deps.keys())

    native.java_test(
        name = name,
        srcs = srcs,
        use_testrunner = False,
        main_class = main_class or "org.junit.platform.console.ConsoleLauncher",
        args = junit_console_args,
        deps = deps,
        runtime_deps = extra_runtime_deps.keys(),
        **kwargs
    )

def _format_maven_jar_name(group_id, artifact_id):
    return ("%s_%s" % (group_id, artifact_id)).replace(".", "_").replace("-", "_")

def _format_maven_jar_dep_name(group_id, artifact_id):
    return "@maven//:%s" % _format_maven_jar_name(group_id, artifact_id)

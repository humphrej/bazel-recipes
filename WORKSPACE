workspace(name = "bazel-recipes")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

# rules_docker section ------------------------------------
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "4521794f0fba2e20f3bf15846ab5e01d5332e587e9ce81629c7f96c793bb7036",
    strip_prefix = "rules_docker-0.14.4",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.14.4/rules_docker-v0.14.4.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

# This is NOT needed when going through the language lang_image
# "repositories" function(s).
load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

load("@io_bazel_rules_docker//repositories:pip_repositories.bzl", "pip_deps")

pip_deps()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
)
load(
    "@io_bazel_rules_docker//java:image.bzl",
    _java_image_repos = "repositories",
)

_java_image_repos()

container_pull(
    name = "ubuntu2004",
    digest = "sha256:93fd0705706e5bdda6cc450b384d8d5afb18fecc19e054fe3d7a2c8c2aeb2c83",
    registry = "index.docker.io",
    repository = "library/ubuntu",
)
# rules_docker section ------------------------------------

# rules_go section ------------------------------------
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "a8d6b1b354d371a646d2f7927319974e0f9e52f73a2452d2b3877118169eb6bb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.23.3/rules_go-v0.23.3.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.23.3/rules_go-v0.23.3.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

http_archive(
    name = "bazel_gazelle",
    sha256 = "cdb02a887a7187ea4d5a27452311a75ed8637379a1287d8eeb952138ea485f7d",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.1/bazel-gazelle-v0.21.1.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.1/bazel-gazelle-v0.21.1.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

load("//:go_dependencies.bzl", "go_repositories")

# gazelle:repository_macro go_dependencies.bzl%go_repositories
go_repositories()

# rules_go section ------------------------------------

# rules_proto_grpc section ------------------------------------

http_archive(
    name = "rules_proto_grpc",
    sha256 = "5f0f2fc0199810c65a2de148a52ba0aff14d631d4e8202f41aff6a9d590a471b",
    strip_prefix = "rules_proto_grpc-1.0.2",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/1.0.2.tar.gz"],
)

load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_repos", "rules_proto_grpc_toolchains")

rules_proto_grpc_toolchains()

rules_proto_grpc_repos()

load("@rules_proto_grpc//java:repositories.bzl", rules_proto_grpc_java_repos = "java_repos")

rules_proto_grpc_java_repos()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories(
    omit_bazel_skylib = True,
    omit_com_google_guava = True,
    omit_com_google_protobuf = True,
    omit_com_google_protobuf_javalite = False,
    omit_net_zlib = True,
)

load("@bazel_tools//tools/build_defs/repo:jvm.bzl", "jvm_maven_import_external")

# This is to upgrade the version of guava used by grpc_java (see omit above)
GUAVA_VERSION = "28.1-jre"

jvm_maven_import_external(
    name = "com_google_guava_guava",
    artifact = "com.google.guava:guava:%s" % GUAVA_VERSION,
    artifact_sha256 = "30beb8b8527bd07c6e747e77f1a92122c2f29d57ce347461a4a55eb26e382da4",
    licenses = ["notice"],  # Apache 2.0
    server_urls = ["https://repo1.maven.org/maven2"],
)

load("@rules_proto_grpc//:repositories.bzl", "bazel_gazelle", "io_bazel_rules_go")

io_bazel_rules_go()

bazel_gazelle()

load("@rules_proto_grpc//go:repositories.bzl", rules_proto_grpc_go_repos = "go_repos")

rules_proto_grpc_go_repos()

# rules_proto_grpc section ------------------------------------

# golang-migrate section ------------------------------------
http_archive(
    name = "golang-migrate",
    build_file_content = 'filegroup( name="files", srcs= [ "migrate.linux-amd64" ], visibility=["//visibility:public"]  )',
    sha256 = "e4d66ada9e98cd502b5ba9c467d641ecbc7a7ba5a8d1e92051aacdcbc10528f0",
    urls = ["https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz"],
)
# golang-migrate section ------------------------------------

# rules_clojure section ------------------------------------

#RULES_CLOJURE_SHA = "0a2b6b06802263e5ce7f0e903e667ae4c103c6fc"
RULES_CLOJURE_SHA = "bad7ead30e3426425d4ae44d974a2bfa868d61e8"
http_archive(name = "rules_clojure",
             strip_prefix = "rules_clojure-%s" % RULES_CLOJURE_SHA,
             sha256 = "a8245b81226cd70a54eae1048b83986b3168a9127930a18f7ac00868385b1bb3",
             url = "https://github.com/griffinbank/rules_clojure/archive/%s.zip" % RULES_CLOJURE_SHA)

load("@rules_clojure//:repositories.bzl", "rules_clojure_deps")
rules_clojure_deps()

load("@rules_clojure//:setup.bzl", "rules_clojure_setup")
rules_clojure_setup()

# rules_clojure section ------------------------------------

load("//:junit_in_container/junit5.bzl", "make_deps")

JUNIT_DEPS_ = make_deps()

# rules_jvm_external section ------------------------------------
RULES_JVM_EXTERNAL_TAG = "4.0"

RULES_JVM_EXTERNAL_SHA = "31701ad93dbfe544d597dbe62c9a1fdd76d81d8a9150c2bf1ecf928ecdf97169"

http_archive(
    name = "rules_jvm_external",
    sha256 = RULES_JVM_EXTERNAL_SHA,
    strip_prefix = "rules_jvm_external-%s" % RULES_JVM_EXTERNAL_TAG,
    url = "https://github.com/bazelbuild/rules_jvm_external/archive/%s.zip" % RULES_JVM_EXTERNAL_TAG,
)

load("@rules_jvm_external//:defs.bzl", "maven_install")
load("@rules_jvm_external//:specs.bzl", "maven")

JUNIT_DEPS = JUNIT_DEPS_["jupiter"] + JUNIT_DEPS_["platform"] + JUNIT_DEPS_["vintage"] + JUNIT_DEPS_["extras"]

maven_install(
    artifacts =
        JUNIT_DEPS + [
          maven.artifact(
              group = "org.clojure",
              artifact = "clojure",
              version = "1.11.1",
              exclusions = [
                  "org.clojure:spec.alpha",
                  "org.clojure:core.specs.alpha"
              ]
          ),
          maven.artifact(
              group = "org.clojure",
              artifact = "spec.alpha",
              version = "0.3.218",
              exclusions = ["org.clojure:clojure"]
            ),
          maven.artifact(
              group = "org.clojure",
              artifact = "core.specs.alpha",
              version = "0.2.62",
              exclusions = [
                  "org.clojure:clojure",
                  "org.clojure:spec.alpha"
            ]),
          "ant:ant-junit:1.6.5",
          "com.google.guava:guava:jar:30.1.1-jre",
          "com.google.truth:truth:1.1.3",
          "junit:junit:4.12",
          "org.apache.ant:ant:1.10.10",
        ],
    fail_on_missing_checksum = False,
    fetch_sources = True,
    maven_install_json = "@bazel-recipes//third_party:maven_install.json",
    repositories = [
        "https://maven.google.com",
        "https://repo1.maven.org/maven2",
    ],
    version_conflict_policy = "pinned",
)

load("@maven//:defs.bzl", "pinned_maven_install")

pinned_maven_install()

# rules_jvm_external section ------------------------------------

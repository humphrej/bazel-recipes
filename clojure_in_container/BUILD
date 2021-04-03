load("@rules_clojure//:rules.bzl", "clojure_binary", "clojure_java_library", "clojure_library")
load("@io_bazel_rules_docker//java:image.bzl", "java_image")

clojure_library(
    name = "clojure_in_container_clojure_library",
    srcs = ["clojure_in_container.clj"],
)

clojure_java_library(
    name = "clojure_in_container_clojure_java_library",
    namespaces = ["clojure_in_container"],
    deps = [
        ":clojure_in_container_clojure_library",
    ],
)

java_image(
    name = "clojure_in_container_image",
    main_class = "clojure_in_container",
    runtime_deps = [
        ":clojure_in_container_clojure_java_library",
    ],
)

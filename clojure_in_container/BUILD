load("@rules_clojure//:rules.bzl", "clojure_library")
load("@io_bazel_rules_docker//java:image.bzl", "java_image")

clojure_library(
    name = "clojure_in_container_clojure_library",
    srcs = ["clojure_in_container.clj"],
    resource_strip_prefix = "clojure_in_container",
    aot = ["clojure_in_container"],
    deps = [ 
            "@maven//:org_clojure_clojure",
            "@maven//:org_clojure_spec_alpha",
            "@maven//:org_clojure_core_specs_alpha",
    ]
)

java_image(
    name = "clojure_in_container_image",
    main_class = "clojure_in_container",
    runtime_deps = [
        ":clojure_in_container_clojure_library",
    ],
)

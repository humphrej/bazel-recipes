load("@io_bazel_rules_docker//java:image.bzl", "java_image")

java_binary(
    name = "java_in_container",
    srcs = glob(["*.java"]),
    main_class = "incontainer.JavaInContainer",
)

java_image(
    name = "java_in_container_image",
    srcs = glob(["*.java"]),
    main_class = "incontainer.JavaInContainer",
)

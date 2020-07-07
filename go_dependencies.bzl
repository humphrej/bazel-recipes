load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_repositories():
    go_repository(
        name = "com_github_sendgrid_rest",
        importpath = "github.com/sendgrid/rest",
        sum = "h1:a2tyRVS0S5kcY6fVq5ihxOTJiGTQROrqf7SkKbmpYzs=",
        version = "v2.6.0+incompatible",
    )
    go_repository(
        name = "com_github_sendgrid_sendgrid_go",
        importpath = "github.com/sendgrid/sendgrid-go",
        sum = "h1:ZSvuHv7JuLnC9iaDubK7iNCfmR4vt8tcbuzsAYQUumo=",
        version = "v3.6.0+incompatible",
    )

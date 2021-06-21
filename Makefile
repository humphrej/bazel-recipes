
current_dir = $(shell pwd)

all: ;

format:
	docker run -ti -u nonroot -v $(current_dir):/workspace bazel/image:xformat_image -v 1 --ignore_directories=data,.git,.ijwb /workspace

repin:
	bazel run @unpinned_maven//:pin


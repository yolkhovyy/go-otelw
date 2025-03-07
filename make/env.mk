GIT_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "none")

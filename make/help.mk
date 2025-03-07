.PHONY: display-help
display-help:
	@echo 'Usage: $(MAKE) <target>'
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
			if (/^## /) { \
				print "  " substr($$1, 4) \
			} else if (/^[[:alnum:]_-]+:.*?##/) { \
				printf "    %-22s%s\n", $$1, $$2 \
			} \
		}' $(MAKEFILE_LIST)

export COMPOSE_DOCKER_CLI_BUILD = 1
export DOCKER_BUILDKIT = 1
export COMPOSE_BAKE = true

DOCO = docker compose -f docker-compose.yml
export DEPENDENCIES = otel-collector

ifdef NR		#----- newrelic
	export OTEL_COLLECTOR_CONFIG = config-newrelic.yml
else ifdef UPT	#----- uptrace
	DOCO := $(DOCO) -f docker-compose.uptrace.yml
	DEPENDENCIES := $(DEPENDENCIES) postgres clickhouse uptrace
	export OTEL_COLLECTOR_CONFIG = config-uptrace.yml
else ifdef GCL	#----- grafana cloud
	export OTEL_COLLECTOR_CONFIG = config-gcloud.yml
else ifdef ALY	#----- grafana cloud via alloy
	DOCO := $(DOCO) -f docker-compose.alloy.yml
	DEPENDENCIES := $(DEPENDENCIES) alloy
	export OTEL_COLLECTOR_CONFIG = config-alloy.yml
else ifdef TPL	#----- tempo, prometheus, loki, promtail, grafana
	DOCO := $(DOCO) -f docker-compose.tplpg.yml
	DEPENDENCIES := $(DEPENDENCIES) tempo prometheus loki promtail grafana
	export OTEL_COLLECTOR_CONFIG = config-tplpg.yml
else			#----- jaeger, prometheus, loki, promtail, grafana
	DOCO := $(DOCO) -f docker-compose.jplpg.yml
	DEPENDENCIES := $(DEPENDENCIES) jaeger prometheus loki promtail grafana
	export OTEL_COLLECTOR_CONFIG = config-jplpg.yml
endif

export SERVICES = $(filter-out ${DEPENDENCIES}, $(shell ${DOCO} config --services))

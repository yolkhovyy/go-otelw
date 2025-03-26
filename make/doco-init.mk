export COMPOSE_DOCKER_CLI_BUILD = 1
export DOCKER_BUILDKIT = 1
export COMPOSE_BAKE = true

DOCO = docker compose -f docker-compose.yml
export DEPENDENCIES = otel-collector

ifdef NR		#----- newrelic
	export OTEL_COLLECTOR_CONFIG = config-newrelic.yml
else ifdef DD	#----- datadog
	export OTEL_COLLECTOR_CONFIG = config-datadog.yml
else ifdef UPT	#----- uptrace
	DOCO := $(DOCO) -f docker-compose.uptrace.yml
	DEPENDENCIES := $(DEPENDENCIES) postgres clickhouse uptrace
	export OTEL_COLLECTOR_CONFIG = config-uptrace.yml
else ifdef EK	#----- elastic, kibana
	DOCO := $(DOCO) -f docker-compose.elastic-kibana.yml
	DEPENDENCIES := $(DEPENDENCIES) elasticsearch kibana
	export OTEL_COLLECTOR_CONFIG = config-elastic-kibana.yml
else ifdef GCL	#----- grafana cloud
	export OTEL_COLLECTOR_CONFIG = config-grafana-cloud.yml
else ifdef ALY	#----- grafana cloud alloy
	DOCO := $(DOCO) -f docker-compose.grafana-cloud-alloy.yml
	DEPENDENCIES := $(DEPENDENCIES) alloy
	export OTEL_COLLECTOR_CONFIG = config-grafana-cloud-alloy.yml
else ifdef GLT	#----- grafana loki, tempo, prometheus
	DOCO := $(DOCO) -f docker-compose.grafana-loki-tempo.yml
	DEPENDENCIES := $(DEPENDENCIES) grafana loki tempo prometheus promtail
	export OTEL_COLLECTOR_CONFIG = config-grafana-loki-tempo.yml
else			#----- grafana loki, jaeger, prometheus
	DOCO := $(DOCO) -f docker-compose.grafana-loki-jaeger.yml
	DEPENDENCIES := $(DEPENDENCIES) grafana loki jaeger prometheus promtail
	export OTEL_COLLECTOR_CONFIG = config-grafana-loki-jaeger.yml
endif

export SERVICES = $(filter-out ${DEPENDENCIES}, $(shell ${DOCO} config --services))

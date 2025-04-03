export DOCKER_BUILDKIT = 1
export COMPOSE_DOCKER_CLI_BUILD = 1
export COMPOSE_BAKE = true

DOCO = docker compose -f docker-compose.yml
export DEPENDENCIES = otel-collector

ifdef NR		#----- newrelic
	export OTEL_COLLECTOR_CONFIG = newrelic.yml
else ifdef HC	#----- honeycomb
	export OTEL_COLLECTOR_CONFIG = honeycomb.yml
else ifdef DD	#----- datadog
	export OTEL_COLLECTOR_CONFIG = datadog.yml
else ifdef DT	#----- dynatrace
	export OTEL_COLLECTOR_CONFIG = dynatrace.yml
else ifdef OO	#----- openobserve
	export OTEL_COLLECTOR_CONFIG = openobserve.yml
else ifdef UPT	#----- uptrace
	DOCO := $(DOCO) -f docker-compose.uptrace.yml
	DEPENDENCIES := $(DEPENDENCIES) postgres clickhouse uptrace
	export OTEL_COLLECTOR_CONFIG = uptrace.yml
else ifdef EK	#----- elastic, kibana
	DOCO := $(DOCO) -f docker-compose.elastic-kibana.yml
	DEPENDENCIES := $(DEPENDENCIES) elasticsearch kibana
	export OTEL_COLLECTOR_CONFIG = elastic-kibana.yml
else ifdef GCL	#----- grafana cloud
	export OTEL_COLLECTOR_CONFIG = grafana-cloud.yml
else ifdef ALY	#----- grafana cloud alloy
	DOCO := $(DOCO) -f docker-compose.grafana-cloud-alloy.yml
	DEPENDENCIES := $(DEPENDENCIES) alloy
	export OTEL_COLLECTOR_CONFIG = grafana-cloud-alloy.yml
else ifdef GLT	#----- grafana loki, tempo, prometheus
	DOCO := $(DOCO) -f docker-compose.grafana-loki-tempo.yml
	DEPENDENCIES := $(DEPENDENCIES) grafana loki tempo prometheus promtail
	export OTEL_COLLECTOR_CONFIG = grafana-loki-tempo.yml
else			#----- grafana loki, jaeger, prometheus
	DOCO := $(DOCO) -f docker-compose.grafana-loki-jaeger.yml
	DEPENDENCIES := $(DEPENDENCIES) grafana loki jaeger prometheus promtail
	export OTEL_COLLECTOR_CONFIG = grafana-loki-jaeger.yml
endif

export SERVICES = $(filter-out ${DEPENDENCIES}, $(shell ${DOCO} config --services))

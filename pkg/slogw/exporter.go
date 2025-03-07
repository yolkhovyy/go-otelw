package slogw

import (
	"context"
	"fmt"
	"strconv"

	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/log"
)

type WithSeverityText struct {
	log.Exporter
}

func (e *WithSeverityText) Export(ctx context.Context, records []log.Record) error {
	for i := range records {
		mapSeverityText(&records[i])
	}

	if err := e.Exporter.Export(ctx, records); err != nil {
		return fmt.Errorf("error exporting: %w", err)
	}

	return nil
}

//nolint:cyclop
func mapSeverityText(record *log.Record) {
	severity := record.Severity()
	switch severity {
	case otellog.SeverityUndefined:
		record.SetSeverityText("UNDEFINED")

	case otellog.SeverityTrace1:
		record.SetSeverityText("TRACE")
	case otellog.SeverityTrace2, otellog.SeverityTrace3, otellog.SeverityTrace4:
		record.SetSeverityText("TRACE" + strconv.Itoa(int(severity)))

	case otellog.SeverityDebug1:
		record.SetSeverityText("DEBUG")
	case otellog.SeverityDebug2, otellog.SeverityDebug3, otellog.SeverityDebug4:
		record.SetSeverityText("DEBUG" + strconv.Itoa(int(severity)))

	case otellog.SeverityInfo1:
		record.SetSeverityText("INFO")
	case otellog.SeverityInfo2, otellog.SeverityInfo3, otellog.SeverityInfo4:
		record.SetSeverityText("INFO" + strconv.Itoa(int(severity)))

	case otellog.SeverityWarn1:
		record.SetSeverityText("WARN")
	case otellog.SeverityWarn2, otellog.SeverityWarn3, otellog.SeverityWarn4:
		record.SetSeverityText("WARN" + strconv.Itoa(int(severity)))

	case otellog.SeverityError1:
		record.SetSeverityText("ERROR")
	case otellog.SeverityError2, otellog.SeverityError3, otellog.SeverityError4:
		record.SetSeverityText("ERROR" + strconv.Itoa(int(severity)))

	case otellog.SeverityFatal1:
		record.SetSeverityText("FATAL")
	case otellog.SeverityFatal2, otellog.SeverityFatal3, otellog.SeverityFatal4:
		record.SetSeverityText("FATAL" + strconv.Itoa(int(severity)))

	default:
		record.SetSeverityText("UNKNOWN")
	}
}

package slogw

import (
	"context"
	"fmt"
	"strconv"

	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type WithSeverityText struct {
	sdklog.Exporter
}

func (e *WithSeverityText) Export(ctx context.Context, records []sdklog.Record) error {
	for i := range records {
		mapSeverityText(&records[i])
	}

	if err := e.Exporter.Export(ctx, records); err != nil {
		return fmt.Errorf("error exporting: %w", err)
	}

	return nil
}

//nolint:cyclop
func mapSeverityText(record *sdklog.Record) {
	severity := record.Severity()
	switch severity {
	case log.SeverityUndefined:
		record.SetSeverityText("UNDEFINED")

	case log.SeverityTrace1:
		record.SetSeverityText("TRACE")
	case log.SeverityTrace2, log.SeverityTrace3, log.SeverityTrace4:
		record.SetSeverityText("TRACE" + strconv.Itoa(int(severity)))

	case log.SeverityDebug1:
		record.SetSeverityText("DEBUG")
	case log.SeverityDebug2, log.SeverityDebug3, log.SeverityDebug4:
		record.SetSeverityText("DEBUG" + strconv.Itoa(int(severity)))

	case log.SeverityInfo1:
		record.SetSeverityText("INFO")
	case log.SeverityInfo2, log.SeverityInfo3, log.SeverityInfo4:
		record.SetSeverityText("INFO" + strconv.Itoa(int(severity)))

	case log.SeverityWarn1:
		record.SetSeverityText("WARN")
	case log.SeverityWarn2, log.SeverityWarn3, log.SeverityWarn4:
		record.SetSeverityText("WARN" + strconv.Itoa(int(severity)))

	case log.SeverityError1:
		record.SetSeverityText("ERROR")
	case log.SeverityError2, log.SeverityError3, log.SeverityError4:
		record.SetSeverityText("ERROR" + strconv.Itoa(int(severity)))

	case log.SeverityFatal1:
		record.SetSeverityText("FATAL")
	case log.SeverityFatal2, log.SeverityFatal3, log.SeverityFatal4:
		record.SetSeverityText("FATAL" + strconv.Itoa(int(severity)))

	default:
		record.SetSeverityText("UNKNOWN")
	}
}

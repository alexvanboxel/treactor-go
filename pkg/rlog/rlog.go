package rlog

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"go.opencensus.io/trace"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
	"log"
	"net/http"
)

type RLogger struct {
	logger            *logging.Logger
	monitoredResource *mrpb.MonitoredResource
	projectId         string
}

func NewRLogger(projectId string, logger *logging.Logger, monitoredResource *mrpb.MonitoredResource) *RLogger {
	return &RLogger{
		logger:            logger,
		monitoredResource: monitoredResource,
		projectId:         projectId,
	}
}

func (l *RLogger) addSpan(ctx context.Context, entry *logging.Entry) *logging.Entry {
	span := trace.FromContext(ctx)
	if span != nil {
		entry.Trace = fmt.Sprintf("projects/%s/traces/%s", l.projectId, span.SpanContext().TraceID.String())
		entry.SpanID = span.SpanContext().SpanID.String()
	}
	return entry
}

func (l *RLogger) InfoF(ctx context.Context, format string, a ...interface{}) {
	l.log(ctx, logging.Info, fmt.Sprintf(format, a...))
}

func (l *RLogger) Info(ctx context.Context, message string) {
	l.log(ctx, logging.Info, message)
}

func (l *RLogger) WarningF(ctx context.Context, format string, a ...interface{}) {
	l.log(ctx, logging.Warning, fmt.Sprintf(format, a...))
}

func (l *RLogger) Warning(ctx context.Context, message string) {
	l.log(ctx, logging.Warning, message)
}

func (l *RLogger) log(ctx context.Context, severity logging.Severity, message string) {
	entry := l.addSpan(ctx, &logging.Entry{
		Severity: severity,
		Payload:  message,
		Resource: l.monitoredResource,
	})
	l.logger.Log(*entry)
}

func (l *RLogger) Error(ctx context.Context, r *http.Request, message string) {
	entry := newREntry(l, logging.Error)
	entry.addPayLoad(r, message)
	entry.addSpan(ctx)
	entry.addErrorLocation()
	entry.addStackTrace()
	entry.log()
}

func (l *RLogger) ErrorErr(ctx context.Context, r *http.Request, message string, err error) {
	l.Error(ctx, r, fmt.Sprintf("%s: %s", message, err.Error()))
}

func (l *RLogger) ErrorF(ctx context.Context, r *http.Request, format string, a ...interface{}) {
	l.Error(ctx, r, fmt.Sprintf(format, a...))
}

func (l *RLogger) Flush() {
	err := l.logger.Flush()
	if err != nil {
		log.Println("Error flushing log")
		log.Println(err)
	}
}

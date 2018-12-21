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

func NewRLogger(projectId string, logger *logging.Logger, monitoredResource *mrpb.MonitoredResource) (*RLogger) {
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

func (l *RLogger) Info(ctx context.Context, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	entry := l.addSpan(ctx, &logging.Entry{
		Severity: logging.Info,
		Payload:  message,
		Resource: l.monitoredResource,
	})
	l.logger.Log(*entry)
}

func (l *RLogger) Error(ctx context.Context, r *http.Request, format string, a ...interface{}) {
	entry := newREntry(l, logging.Error)
	entry.addPayLoad(r, format, a)
	entry.addSpan(ctx)
	entry.addErrorLocation()
	entry.addStackTrace()
	entry.log()
}

func (l *RLogger) Warning(ctx context.Context, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	entry := l.addSpan(ctx, &logging.Entry{
		Severity: logging.Warning,
		Payload:  message,
		Resource: l.monitoredResource,
	})
	l.logger.Log(*entry)
}

func (l *RLogger) Flush() {
	err := l.logger.Flush()
	if err != nil {
		log.Println("Error flushing log")
		log.Println(err)
	}
}

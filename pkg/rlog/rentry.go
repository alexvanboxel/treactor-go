package rlog

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"go.opencensus.io/trace"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
	"log"
	"net/http"
	"runtime"
)

type rPayLoad struct {
	EventTime      string          `json:"eventTime,omitempty"`
	ServiceContext rServiceContext `json:"serviceContext,omitempty"`
	Message        string          `json:"message,omitempty"`
	Context        rContext        `json:"context,omitempty"`
}

type rServiceContext struct {
	Service string `json:"service,omitempty"`
	Version string `json:"version,omitempty"`
}

type rContext struct {
	HttpRequest    rHttpRequest    `json:"httpRequest,omitempty"`
	User           string          `json:"user,omitempty"`
	ReportLocation rReportLocation `json:"reportLocation,omitempty"`
}

type rHttpRequest struct {
	Method             string `json:"method,omitempty"`
	Url                string `json:"url,omitempty"`
	UserAgent          string `json:"userAgent,omitempty"`
	Referrer           string `json:"referrer,omitempty"`
	ResponseStatusCode int    `json:"responseStatusCode,omitempty"`
	RemoteIp           string `json:"remoteIp,omitempty"`
}

type rReportLocation struct {
	FilePath     string `json:"filePath,omitempty"`
	LineNumber   int    `json:"lineNumber,omitempty"`
	FunctionName string `json:"functionName,omitempty"`
}

type REntry struct {
	projectId string
	Entry     logging.Entry
	PayLoad   rPayLoad
}

func (e *REntry) addSpan(ctx context.Context) {
	span := trace.FromContext(ctx)
	if span != nil {
		e.Entry.Trace = fmt.Sprintf("projects/%s/traces/%s", e.projectId, span.SpanContext().TraceID.String())
		e.Entry.SpanId = span.SpanContext().SpanID.String()
	}
}

func (e *REntry) addPayLoad(request *http.Request, s string, a []interface{}) {
	_, fn, line, _ := runtime.Caller(2)

	e.Entry.Payload = rPayLoad{
		Message: s,
		ServiceContext: rServiceContext{
			Service: "reactor",
			Version: "1",
		},
		Context: rContext{
			HttpRequest: rHttpRequest{
				Url: request.URL.String(),
			},
			ReportLocation: rReportLocation{
				FilePath:     fn,
				LineNumber:   line,
				FunctionName: "test",
			},
		},
	}
}

func (e *REntry) addErrorLocation() {
}

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
		entry.SpanId = span.SpanContext().SpanID.String()
	}
	return entry
}

/*
{
  "eventTime": string,
  "rServiceContext": {
    "service": string,     // Required.
    "version": string
  },
  "message": string,       // Required. Should contain the full exception
                           // message, including the stack trace.
  "context": {
    "rHttpRequest": {
      "method": string,
      "url": string,
      "userAgent": string,
      "referrer": string,
      "responseStatusCode": number,
      "remoteIp": string
    },
    "user": string,
    "rReportLocation": {    // Required if no stack trace in 'message'.
      "filePath": string,
      "lineNumber": number,
      "functionName": string
    }
  }
}

*/

func (l *RLogger) addRequest(entry *logging.Entry) *logging.Entry {
	return entry
}

func (l *RLogger) addPayload(entry *logging.Entry) *logging.Entry {
	return entry
}

func (l *RLogger) addResource(entry *logging.Entry) *logging.Entry {
	return entry
}

//func  (l *RLogger) log(ctx context.Context, r *http.Request, server format string, a ...interface{})

func (l *RLogger) Info(ctx context.Context, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a)
	entry := l.addSpan(ctx, &logging.Entry{
		Severity: logging.Info,
		Payload:  message,
		Resource: l.monitoredResource,
	})
	l.logger.Log(*entry)
}

func (l *RLogger) Error(ctx context.Context, r *http.Request, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a)
	entry := &REntry{
		Entry: logging.Entry{
			Severity: logging.Error,
			Payload:  message,
			Resource: l.monitoredResource,
		},
	}
	entry.addSpan(ctx)
	entry.addPayLoad(r, format, a)
	entry.addErrorLocation()

	l.logger.Log(entry.Entry)
}

func (l *RLogger) Warning(ctx context.Context, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a)
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

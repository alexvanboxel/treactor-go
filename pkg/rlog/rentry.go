package rlog

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"go.opencensus.io/trace"
	"net/http"
	"runtime"
	"runtime/debug"
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

type rEntry struct {
	entry   logging.Entry
	payLoad rPayLoad
	rLogger *RLogger
}

func newREntry(rLogger *RLogger, severity logging.Severity) *rEntry {
	return &rEntry{
		rLogger: rLogger,
		entry: logging.Entry{
			Severity: logging.Error,
			Resource: rLogger.monitoredResource,
		},
	}
}

func (e *rEntry) addSpan(ctx context.Context) {
	span := trace.FromContext(ctx)
	if span != nil {
		e.entry.Trace = fmt.Sprintf("projects/%s/traces/%s", e.rLogger.projectId, span.SpanContext().TraceID.String())
		e.entry.SpanID = span.SpanContext().SpanID.String()
	}
}

func (e *rEntry) addPayLoad(request *http.Request, s string, a []interface{}) {
	e.payLoad = rPayLoad{
		Message: s,
		ServiceContext: rServiceContext{
			Service: "reactor",
			Version: "1",
		},
		Context: rContext{
			HttpRequest: rHttpRequest{
				Method: request.Method,
				Url: request.URL.String(),

			},
		},
	}
}

func (e *rEntry) addErrorLocation() {
	_, fn, line, _ := runtime.Caller(2)
	e.payLoad.Context.ReportLocation = rReportLocation{
		FilePath:     fn,
		LineNumber:   line,
		FunctionName: "<undefined>",
	}
}

func (e *rEntry) log() {
	entry := e.entry
	entry.Payload = e.payLoad
	e.rLogger.logger.Log(entry)
}

func (e *rEntry) addStackTrace() {
	e.payLoad.Message = fmt.Sprintf("%s\n%s", e.payLoad.Message, string(debug.Stack()))
}

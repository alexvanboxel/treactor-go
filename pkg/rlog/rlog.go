package rlog

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
	"log"
)

type RLogger struct {
	logger            *logging.Logger
	monitoredResource *mrpb.MonitoredResource
}

func NewQLogger(logger *logging.Logger, monitoredResource *mrpb.MonitoredResource) (*RLogger) {
	return &RLogger{
		logger:            logger,
		monitoredResource: monitoredResource,
	}
}

func (l *RLogger) Info(ctx context.Context, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a)
	l.logger.Log(logging.Entry{
		Severity: logging.Info,
		Payload:  message,
		Resource: l.monitoredResource,
	})
	log.Print(message)
}

func (l *RLogger) Warning(ctx context.Context, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a)
	l.logger.Log(logging.Entry{
		Severity: logging.Warning,
		Payload:  message,
		Resource: l.monitoredResource,
	})
	log.Print(message)
}

func (l *RLogger) Flush() {
	err := l.logger.Flush()
	if err != nil {
		log.Println("Errror flushig log")
		log.Println(err)
	}
}

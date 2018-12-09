package client

import (
	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
	"log"
	"net/http"
	"sync"
)

var (
	LoggingClient *logging.Client
)

type PubSub struct {
	client       *pubsub.Client
	PostCaptured *pubsub.Topic
}

func NewPubSub() (*PubSub, error) {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := ""

	// Creates a client.
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// "consumer-ns.consumer.producer-ns.producer.topic-xxx"

	postCaptured := client.TopicInProject("reactor.reactor.request-captured", projectID)

	ensure := &PubSub{
		client:       client,
		PostCaptured: postCaptured,
	}

	//ensure.SyncCache()

	return ensure, nil
}

func initCencus(wg *sync.WaitGroup, projectId string) {
	defer wg.Done()
	sd, err := stackdriver.NewExporter(stackdriver.Options{
		//ProjectID: projectId,
	})
	if err != nil {
		log.Fatal(err)
	}

	trace.RegisterExporter(sd)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	exporter, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		log.Fatal(err)
	}
	view.RegisterExporter(exporter)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", exporter)
		zpages.Handle(mux, "/debug")
		log.Fatal(http.ListenAndServe("127.0.0.1:4999", mux))
	}()

}

func GoogleCloudInit() {
	projectID := "quantum-research" //os.Getenv("GOOGLE_PROJECT_ID")
	//initMonitoredResource(projectID)
	wg := sync.WaitGroup{}
	wg.Add(1)
	//initProfiler(&wg)
	//initLoggingClient(&wg, projectID)
	//initPubSubCient(&wg, projectID)
	initCencus(&wg, projectID)
	wg.Wait()
}

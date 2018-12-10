package client

import (
	"cloud.google.com/go/logging"
	"cloud.google.com/go/profiler"
	"cloud.google.com/go/pubsub"
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/alexvanboxel/reactor/pkg/rlog"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	PubsubClient      *pubsub.Client
	LoggingClient     *logging.Client
	MonitoredResource *mrpb.MonitoredResource
	Logger            *rlog.RLogger
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

func initLoggingClient(wg *sync.WaitGroup, projectId string) () {
	defer wg.Done()
	var err error
	ctx := context.Background()
	LoggingClient, err = logging.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("Failed to create pubsub client: %v", err)
	}
	Logger = rlog.NewRLogger(projectId, LoggingClient.Logger("reactor"), MonitoredResource)
}

func initMonitoredResource(projectId string) {
	MonitoredResource = &mrpb.MonitoredResource{
		Type:   "global",
		Labels: make(map[string]string),
	}
	MonitoredResource.Labels["project_id"] = projectId

	//MonitoredResource = &mrpb.MonitoredResource{
	//	Type:   "k8s_container",
	//	Labels: make(map[string]string),
	//}
	//MonitoredResource.Labels["project_id"] = projectId
	//MonitoredResource.Labels["namespace_name"] = os.Getenv("POD_NAMESPACE")
	//MonitoredResource.Labels["pod_name"] = os.Getenv("POD_NAME")
	//clusterName, err := metadata.InstanceAttributeValue("cluster-name")
	//if err != nil {
	//	log.Fatalf("Failed to get cluster_name from meta data server: %v", err)
	//}
	//MonitoredResource.Labels["cluster_name"] = clusterName
	//clusterLocation, err := metadata.InstanceAttributeValue("cluster-location")
	//if err != nil {
	//	log.Fatalf("Failed to get cluster_location from meta data server: %v", err)
	//}
	//MonitoredResource.Labels["location"] = clusterLocation
	//MonitoredResource.Labels["container_name"] = "proton-gdpr-api"
}

// Get a single pubsub client and attach it to the background context
func initPubSubCient(wg *sync.WaitGroup, projectId string) () {
	defer wg.Done()
	var err error
	ctx := context.Background()
	PubsubClient, err = pubsub.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("Failed to create pubsub client: %v", err)
	}
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

func initProfiler(wg *sync.WaitGroup) {
	defer wg.Done()
	if err := profiler.Start(profiler.Config{Service: "gdpr-api", ServiceVersion: "0"}); err != nil {
		log.Fatal(err)
	}
}

func GoogleCloudInit() {
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	initMonitoredResource(projectID)
	wg := sync.WaitGroup{}
	wg.Add(3)
	//initProfiler(&wg)
	initLoggingClient(&wg, projectID)
	initPubSubCient(&wg, projectID)
	initCencus(&wg, projectID)
	wg.Wait()
}

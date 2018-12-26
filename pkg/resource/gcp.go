package resource

import (
	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/profiler"
	"cloud.google.com/go/pubsub"
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"fmt"
	"github.com/alexvanboxel/reactor/pkg/config"
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
	fmt.Printf("* Google Cloud:  Logging initialized.\n")
}

func initMonitoredResource(projectId string) {
	if config.IsKubernetesMode() {
		initMonitoredResourceAsK8SContainer(projectId)
	}

	if MonitoredResource == nil {
		initMonitoredResourceAsGlobal(projectId)
	}
}

func initMonitoredResourceAsK8SContainer(projectId string) {
	if metadata.OnGCE() {
		// see: https://cloud.google.com/monitoring/api/resources#tag_k8s_container
		mr := &mrpb.MonitoredResource{
			Type:   "k8s_container",
			Labels: make(map[string]string),
		}
		mr.Labels["namespace_name"] = os.Getenv("POD_NAMESPACE")
		mr.Labels["pod_name"] = os.Getenv("POD_NAME")
		mr.Labels["container_name"] = config.AppName
		mr.Labels["project_id"] = projectId
		clusterName, err := metadata.InstanceAttributeValue("cluster-name")
		if err != nil {
			fmt.Printf("! Google Cloud: Failed to get cluster_name from meta data server: %v\n", err)
			return
		}
		mr.Labels["cluster_name"] = clusterName
		clusterLocation, err := metadata.InstanceAttributeValue("cluster-location")
		if err != nil {
			fmt.Printf("! Google Cloud:Failed to get cluster_location from meta data server: %v\n", err)
			return
		}
		mr.Labels["location"] = clusterLocation
		MonitoredResource = mr
		fmt.Printf("* Google Cloud: MonitoredResource set to k8s_container.\n")
	} else {
		fmt.Printf("! Google Cloud: Meta Server not accessible, couldn't configure MonitoredResource.\n")
	}

}

func initMonitoredResourceAsGlobal(projectId string) {
	// see: https://cloud.google.com/monitoring/api/resources#tag_global
	MonitoredResource = &mrpb.MonitoredResource{
		Type:   "global",
		Labels: make(map[string]string),
	}
	MonitoredResource.Labels["project_id"] = projectId
	fmt.Printf("* Google Cloud: MonitoredResource set to global.\n")

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
	fmt.Printf("* Google Cloud: PubSub initialized.\n")
}

func initCencus(wg *sync.WaitGroup, projectId string) {
	defer wg.Done()
	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectId,
	})
	if err != nil {
		log.Fatal(err)
	}

	trace.RegisterExporter(sd)
	if config.IsDebug() {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		fmt.Printf("* OpenCensus: Trace running in debug, all requests will be traced.\n")
	} else {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1e-2)})
		fmt.Printf("* OpenCensus: Trace probability sampling is active.\n")
	}

	exporter, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		log.Fatal(err)
	}
	view.RegisterExporter(exporter)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", exporter)
		zpages.Handle(mux, "/debug")
		fmt.Printf("* OpenCensus zPages: Running on 127.0.0.1:4999.\n")
		log.Fatal(http.ListenAndServe("127.0.0.1:4999", mux))
	}()

}

func initProfiler(wg *sync.WaitGroup) {
	defer wg.Done()
	if config.IsProfiling() {
		fmt.Printf("* StackDriver: Profiler active.\n")
		if err := profiler.Start(profiler.Config{Service: config.AppName, ServiceVersion: config.AppVersion}); err != nil {
			log.Fatal(err)
		}
	}
}

func GoogleCloudInit() {
	fmt.Printf("Start initializing reactor.\n")
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	initMonitoredResource(projectID)
	wg := sync.WaitGroup{}
	wg.Add(3)
	go initProfiler(&wg)
	go initLoggingClient(&wg, projectID)
	//go initPubSubCient(&wg, projectID)
	go initCencus(&wg, projectID)
	wg.Wait()
	fmt.Printf("Finished initializing reactor.\n")
}

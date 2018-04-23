package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/tsloughter/grafana-operator/pkg/controller"
	"github.com/tsloughter/grafana-operator/pkg/grafana"
	"github.com/tsloughter/grafana-operator/pkg/kubernetes"
)

var (
	grafanaUrl        = flag.String("grafana-url", "", "The url to issue requests to update dashboards to.")
	runOutsideCluster = flag.Bool("run-outside-cluster", false, "Set this flag when running outside of the cluster.")
)

func main() {
	flag.Parse()

	// Set logging output to standard console out
	log.SetOutput(os.Stdout)

	if *grafanaUrl == "" {
		log.Println("Missing grafana-url")
		flag.Usage()
		os.Exit(1)
	}

	gUrl, err := url.Parse(*grafanaUrl)
	if err != nil {
		log.Fatalf("Grafana URL could not be parsed: %s", *grafanaUrl)
	}

	if os.Getenv("GRAFANA_USER") != "" && os.Getenv("GRAFANA_PASSWORD") == "" {
		gUrl.User = url.User(os.Getenv("GRAFANA_USER"))
	}

	if os.Getenv("GRAFANA_USER") != "" && os.Getenv("GRAFANA_PASSWORD") != "" {
		gUrl.User = url.UserPassword(os.Getenv("GRAFANA_USER"), os.Getenv("GRAFANA_PASSWORD"))
	}

	g := grafana.New(gUrl)

	sigs := make(chan os.Signal, 1) // Create channel to receive OS signals
	stop := make(chan struct{})     // Create channel to receive stop signal

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT) // Register the sigs channel to receieve SIGTERM

	wg := &sync.WaitGroup{} // Goroutines can add themselves to this to be waited on so that they finish

	// Create clientset for interacting with the kubernetes cluster
	clientset, err := kubernetes.NewClientSet(*runOutsideCluster)

	if err != nil {
		panic(err.Error())
	}

	controller.NewConfigMapController(clientset, g).Run(stop, wg)

	<-sigs // Wait for signals (this hangs until a signal arrives)
	log.Printf("Shutting down...")

	close(stop) // Tell goroutines to stop themselves
	wg.Wait()   // Wait for all to be stopped
}

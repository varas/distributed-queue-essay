package main

import (
	"context"

	"log"

	"flag"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"runtime/pprof"

	"github.com/varas/distributed-queue-essay/pkg/cluster"
	"github.com/varas/distributed-queue-essay/pkg/node"
)

var (
	port        = flag.Int("port", node.DefaultPort, fmt.Sprintf("Queue endpoint -port %d", node.DefaultPort))
	clusterPort = flag.Int("cluster-port", cluster.DefaultPort, fmt.Sprintf("Cluster synchronization port -cluster-port %d", cluster.DefaultPort))
	cpuProfile  = flag.String("cpuprofile", "", "write cpu profile to file, -cpuprofile profile.cpu")
)

func init() {
	flag.Parse()
}

func main() {

	// profiling
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// cluster initialization

	// hardcoded some initial known members, could also be arguments
	queueCluster := cluster.NewCluster(
		*clusterPort,
		*cluster.NewMember("11.11.11.11", cluster.DefaultPort),
		*cluster.NewMember("11.11.11.12", cluster.DefaultPort),
	)

	if err := queueCluster.Join(); err != nil {
		panic(err)
	}
	defer func() {
		_ = queueCluster.Leave()
	}()

	// node initialization

	node := node.New(*port, cluster.NewClusterPublisher(*clusterPort, queueCluster))

	// launch daemon
	go node.Run(context.Background())

	// wait for runtime start
	go func() {
		<-node.Ready
		log.Printf("queue node listening on tcp/%d", *port)

	}()

	// wait for ctrl+c to stop the daemon
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-c
	close(node.Stop)
}

package cluster

// QueueEvent event occurred on the queue
type QueueEvent struct {
	// TODO
}

// Publisher publishes queue and dequeue events on the cluster
type Publisher interface {
	Publish(QueueEvent) error
}

// NewClusterPublisher cluster queue-command publisher that would manage cluster consistency in background mode (ie. goroutine on instantiation).
// For this illustrative example it just shows a sample constructor to explicit that it will manage a cluster and
// receive cluster operations listening on a given port.
func NewClusterPublisher(port int, cl Cluster) Publisher {

	// go ClusterManager.Run()

	return &noopClusterPublisher{}
}

// noopClusterPublisher does nothing, for testing
type noopClusterPublisher struct{}

func (q *noopClusterPublisher) Publish(QueueEvent) error { return nil }

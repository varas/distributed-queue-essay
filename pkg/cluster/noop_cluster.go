package cluster

// NoopCluster (no-operation cluster) does nothing, for testing
type NoopCluster struct{}

// NewCluster instantiates a new cluster manager
// For this illustrative example just shows a sample constructor and does nothing.
func NewCluster(port int, knownMembers ...Member) Cluster {
	return &NoopCluster{}
}

// Join ...
func (c *NoopCluster) Join() error {
	return nil
}

// Leave ...
func (c *NoopCluster) Leave() error {
	return nil
}

// EventReceivers ...
func (c *NoopCluster) EventReceivers() (MemberJoined, MemberLeft) {
	return make(<-chan Member), make(<-chan Member)
}

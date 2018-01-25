// Package cluster decouples from: io, transport, consistency model and/or 3rd party membership/store libraries.
package cluster

// DefaultPort default cluster management port
const DefaultPort = 5555

// Member models a cluster peer endpoint
type Member struct {
	ip   string
	port uint
}

// NewMember instantiates a member struct
func NewMember(ip string, port int) *Member {
	return &Member{
		ip:   ip,
		port: uint(port),
	}
}

// Cluster manages cluster actions: for now Membership, but may include cluster Runtime (start/stop/status)
// It also provide event receivers to hook on.
// This is an illustrative approach, again as a distributed-system, this component may vary its API based on
// its member consistency model
type Cluster interface {
	Join() error
	Leave() error
	EventReceivers() (MemberJoined, MemberLeft)
}

// MemberJoined returns a member that joined the queue cluster
type MemberJoined <-chan Member

// MemberLeft returns a member that left the cluster
type MemberLeft <-chan Member

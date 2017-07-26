package warewulf

import (
	"fmt"
	"time"

	"github.com/altairsix/eventsource"
)

//Node represents a physical or virtual system that is to be managed, provision, etc
type Node struct {
	ID        string
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	State     string
	Arch      string
	Bootstrap *Bootstrap
	VNFS      *VNFS
	Netdevs   map[string]*Netdev // Key of map[string]*Netdev is CIDR subnet, eg. 196.168.1.0/16
}

//Netdev reprents a physical or virtual network adapter in a node
type Netdev struct {
	HWAddr  string
	Name    string
	IP      string
	Netmask string
	Gateway string
	Domain  string
	// VLAN    int
	// SubNetDev []*Netdev
}

// NodeCreated type represents the event of a node creation
type NodeCreated struct {
	eventsource.Model
	State string
}

// NodeDisabled type represents the event disabling a node
type NodeDisabled struct {
	eventsource.Model
	State string
}

// NodeDelete type represents the event disabling a node
type NodeDelete struct {
	eventsource.CommandModel
	State string
}

//NodeArchSet type represents the event of the architecture of a node being set
type NodeArchSet struct {
	eventsource.Model
	Arch string
}

//NodeBootstrapSet type represents the event of a bootstrap of a node being set
type NodeBootstrapSet struct {
	eventsource.Model
	Bootstrap *Bootstrap
}

//NodeVNFSSet type represents the event of a VNFS of a node being set
type NodeVNFSSet struct {
	eventsource.Model
	VNFS *VNFS
}

//NodeNetdevsSet type represents the event of a Netdev of a node being set
type NodeNetdevsSet struct {
	eventsource.Model
	Netdevs map[string]*Netdev
}

//On parses an event and applies the event's changes to the Node object
func (n *Node) On(event eventsource.Event) error {
	switch e := event.(type) {
	case *NodeCreated:
		n.Version = e.Model.Version
		n.ID = e.Model.ID

	case *NodeArchSet:
		n.Version = e.Model.Version
		n.Arch = e.Arch

	case *NodeBootstrapSet:
		n.Version = e.Model.Version
		n.Bootstrap = e.Bootstrap

	case *NodeVNFSSet:
		n.Version = e.Model.Version
		n.VNFS = e.VNFS

	case *NodeNetdevsSet:
		n.Version = e.Model.Version
		n.Netdevs = e.Netdevs

	default:
		return fmt.Errorf("unhandled event, %v", e)
	}

	return nil
}

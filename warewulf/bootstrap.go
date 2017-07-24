package warewulf

import (
	"fmt"
	"time"

	"github.com/altairsix/eventsource"
)

// Bootstrap represents a kernel and initramfs
type Bootstrap struct {
	ID           string
	Version      int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	State        string
	Arch         string
	Path         string
	Checksum     string
	Size         int64
	CompressAlgo string
}

//BootstrapCreated represents the event of the bootstrap being created
type BootstrapCreated struct {
	eventsource.Model
}

//BootstrapChanged represents the event of the bootstrap files being changed
type BootstrapChanged struct {
	eventsource.Model
	Arch         string
	Path         string
	Checksum     string
	Size         int64
	CompressAlgo string
}

//BootstrapDeleted represents the event of the bootstrap files being deleted
type BootstrapDeleted struct {
	eventsource.Model
	State string
}

//On parses an event and applies the event's changes to the Bootstrap object
func (b *Bootstrap) On(event eventsource.Event) error {
	switch e := event.(type) {
	case *BootstrapCreated:
		b.Version = e.Model.Version
		b.ID = e.Model.ID

	case *BootstrapChanged:
		b.Version = e.Model.Version

		if e.Arch != "" {
			b.Arch = e.Arch
		}

		if e.Path != "" {
			b.Path = e.Path
		}

		if e.Checksum != "" {
			b.Checksum = e.Checksum
		}

		if e.Size != 0 {
			b.Size = e.Size
		}

		if e.CompressAlgo != "" {
			b.CompressAlgo = e.CompressAlgo
		}
	case *BootstrapDeleted:
		b.Version = e.Model.Version
		b.UpdatedAt = e.At
		b.State = "Deleted"

	default:
		return fmt.Errorf("unhandled event, %v", e)
	}

	return nil
}

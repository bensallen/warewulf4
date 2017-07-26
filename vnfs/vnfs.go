package warewulf

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/altairsix/eventsource"
)

// VNFS represents a userland OS image compressed CPIO format
type VNFS struct {
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

// Create saves a new VNFS by building a CreateVNFS command and appling it against the repository.
func (v *VNFS) Create(ctx context.Context, repo *eventsource.Repository) error {
	if v.ID == "" {
		return fmt.Errorf("ID of VNFS must be specified")
	}

	createVNFS := CreateVNFS{
		CommandModel: eventsource.CommandModel{ID: v.ID},
		Arch:         v.Arch,
		Path:         v.Path,
		Checksum:     v.Checksum,
		Size:         v.Size,
		CompressAlgo: v.CompressAlgo,
	}

	_, err := repo.Apply(ctx, createVNFS)
	return err
}

// Read attemps to fetch the VNFS aggregrate from the event repository. v.ID must be specified
// as it is used the aggregate ID.
func (v *VNFS) Read(ctx context.Context, repo *eventsource.Repository) error {
	if v.ID == "" {
		return fmt.Errorf("ID of VNFS must be specified")
	}

	aggregate, err := repo.Load(ctx, v.ID)
	if aggregate == nil {
		return fmt.Errorf("VNFS not found")
	}

	vnfs, ok := aggregate.(*VNFS)

	if !ok {
		return fmt.Errorf("ID returned an aggregate that is not a VNFS")
	}

	// Copy values of casted aggregigate to *v
	*v = *vnfs

	return err
}

func (v *VNFS) Update(repo *eventsource.Repository, ctx context.Context) error {
	if v.ID == "" {
		return fmt.Errorf("ID of VNFS must be specified")
	}
	updateVNFS := UpdateVNFS{
		CommandModel: eventsource.CommandModel{ID: v.ID},
		Arch:         v.Arch,
		Path:         v.Path,
		Checksum:     v.Checksum,
		Size:         v.Size,
		CompressAlgo: v.CompressAlgo,
	}

	_, err := repo.Apply(ctx, updateVNFS)
	return err

}

func (v *VNFS) Delete(repo *eventsource.Repository, ctx context.Context) error {
	if v.ID == "" {
		return fmt.Errorf("ID of VNFS must be specified")
	}
	deleteVNFS := DeleteVNFS{
		CommandModel: eventsource.CommandModel{ID: v.ID},
	}
	_, err := repo.Apply(ctx, deleteVNFS)
	return err
}

//VNFSCreated represents the event of the bootstrap being created
type VNFSCreated struct {
	eventsource.Model
	State        string
	Arch         string
	Path         string
	Checksum     string
	Size         int64
	CompressAlgo string
}

//VNFSUpdated represents the event of the VNFS files being updated
type VNFSUpdated struct {
	eventsource.Model
	Arch         string
	Path         string
	Checksum     string
	Size         int64
	CompressAlgo string
}

//VNFSDeleted represents the event of the VNFS files being deleted
type VNFSDeleted struct {
	eventsource.Model
	State string
}

//On parses event types and applies the event's changes to the VNFS object
func (v *VNFS) On(event eventsource.Event) error {
	switch e := event.(type) {
	case *VNFSCreated:
		v.Version = e.Model.Version
		v.ID = e.Model.ID
		v.State = "Created"
		v.CreatedAt = e.At
		v.UpdatedAt = e.At

		if e.Arch != "" {
			v.Arch = e.Arch
		}

		if e.Path != "" {
			v.Path = e.Path
		}

		if e.Checksum != "" {
			v.Checksum = e.Checksum
		}

		if e.Size != 0 {
			v.Size = e.Size
		}

		if e.CompressAlgo != "" {
			v.CompressAlgo = e.CompressAlgo
		}

	case *VNFSUpdated:
		v.Version = e.Model.Version
		v.UpdatedAt = e.At

		if e.Arch != "" {
			v.Arch = e.Arch
		}

		if e.Path != "" {
			v.Path = e.Path
		}

		if e.Checksum != "" {
			v.Checksum = e.Checksum
		}

		if e.Size != 0 {
			v.Size = e.Size
		}

		if e.CompressAlgo != "" {
			v.CompressAlgo = e.CompressAlgo
		}

	case *VNFSDeleted:
		v.Version = e.Model.Version
		v.UpdatedAt = e.At
		v.State = "Deleted"

	default:
		return fmt.Errorf("unhandled event, %v, type: %s", e, reflect.TypeOf(e))
	}

	return nil
}

//CreateVNFS represents the command to create a VNFS
type CreateVNFS struct {
	eventsource.CommandModel
	Arch         string
	Path         string
	Checksum     string
	Size         int64
	CompressAlgo string
}

//UpdateVNFS represents the command to create a VNFS
type UpdateVNFS struct {
	eventsource.CommandModel
	Arch         string
	Path         string
	Checksum     string
	Size         int64
	CompressAlgo string
}

//DeleteVNFS represents the command to delete VNFS files
type DeleteVNFS struct {
	eventsource.CommandModel
}

//Apply implements the CommandHandler interface for VNFS
func (v *VNFS) Apply(ctx context.Context, command eventsource.Command) ([]eventsource.Event, error) {
	switch c := command.(type) {
	case *CreateVNFS:
		if v.State != "" {
			return nil, fmt.Errorf("VNFS, %v, already exists, use an UpdateVNFS type instead", command.AggregateID())
		}
		vnfsCreated := &VNFSCreated{
			Model:        eventsource.Model{ID: command.AggregateID(), Version: v.Version + 1, At: time.Now()},
			Arch:         c.Arch,
			Path:         c.Path,
			Checksum:     c.Checksum,
			Size:         c.Size,
			CompressAlgo: c.CompressAlgo,
		}
		return []eventsource.Event{vnfsCreated}, nil

	case *UpdateVNFS:
		vnfsUpdated := &VNFSUpdated{
			Model:        eventsource.Model{ID: command.AggregateID(), Version: v.Version + 1, At: time.Now()},
			Arch:         c.Arch,
			Path:         c.Path,
			Checksum:     c.Checksum,
			Size:         c.Size,
			CompressAlgo: c.CompressAlgo,
		}
		return []eventsource.Event{vnfsUpdated}, nil

	case *DeleteVNFS:
		if v.State == "Deleted" {
			return nil, fmt.Errorf("VNFS, %v, is already deleted", command.AggregateID())
		}
		vnfsDeleted := &VNFSDeleted{
			Model: eventsource.Model{ID: command.AggregateID(), Version: v.Version + 1, At: time.Now()},
		}
		return []eventsource.Event{vnfsDeleted}, nil

	default:
		return nil, fmt.Errorf("unhandled command, %v", c)
	}
}

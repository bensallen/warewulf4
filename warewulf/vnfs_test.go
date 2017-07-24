package warewulf

import (
	"context"
	"testing"
	"time"

	"reflect"

	"github.com/altairsix/eventsource"
)

func TestVNFSOn(t *testing.T) {

	t.Run("VNFSCreated", func(t *testing.T) {
		v1 := VNFS{}
		timeNow := time.Now()
		v2 := VNFS{
			ID:        "test",
			Version:   1,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			State:     "Created",
		}
		vnfsCreated := VNFSCreated{
			Model: eventsource.Model{ID: v2.ID, Version: v2.Version, At: v2.CreatedAt},
		}
		err := v1.On(&vnfsCreated)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
		if v1.State != "Created" {
			t.Fatalf("State not Created, %s instead", v1.State)
		}
		if !reflect.DeepEqual(v1, v2) {
			t.Fatalf("Mismatch: %v", v1)
		}

	})

	t.Run("VNFSUpdated", func(t *testing.T) {
		v1 := VNFS{
			ID:        "test",
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			State:     "Created",
		}
		v2 := VNFS{
			ID:        "test",
			Version:   2,
			CreatedAt: v1.CreatedAt,
			UpdatedAt: time.Now(),
			State:     "Created",
			Arch:      "x86_64",
			Checksum:  "0f82992b7693500212e8dd9d7542a5e960c8f787a25f7ce2497a51202c16e558cf4b61a48eade93135bf96c79aafb12dc6dbabdb10907f84571c886157714a74",
			Size:      123456,
			Path:      "/test/path/to/nothing",
		}

		vnfsUpdated := VNFSUpdated{
			Model:    eventsource.Model{ID: v2.ID, Version: v2.Version, At: v2.UpdatedAt},
			Arch:     v2.Arch,
			Checksum: v2.Checksum,
			Size:     v2.Size,
			Path:     v2.Path,
		}
		err := v1.On(&vnfsUpdated)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
		if !reflect.DeepEqual(v1, v2) {
			t.Fatalf("Mismatch: %v", v1)
		}

	})

	t.Run("VNFSDeleted", func(t *testing.T) {
		v1 := VNFS{
			ID: "test",
		}
		v2 := VNFS{
			ID:        "test",
			Version:   10,
			UpdatedAt: time.Now(),
			State:     "Deleted",
		}
		vnfsDeleted := VNFSDeleted{
			Model: eventsource.Model{ID: v2.ID, Version: v2.Version, At: v2.UpdatedAt},
			State: "Created",
		}
		err := v1.On(&vnfsDeleted)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
		if !reflect.DeepEqual(v1, v2) {
			t.Fatalf("Mismatch: %v", v1)
		}

	})

}

func TestVNFSApply(t *testing.T) {
	vnfsID := "test"

	serializer := eventsource.NewJSONSerializer(
		VNFSCreated{},
		VNFSUpdated{},
		VNFSDeleted{},
	)
	repo := eventsource.New(&VNFS{},
		eventsource.WithSerializer(serializer),
	)
	ctx := context.Background()

	t.Run("CreateVNFS", func(t *testing.T) {
		v2 := VNFS{
			Arch:     "x86_64",
			Checksum: "0f82992b7693500212e8dd9d7542a5e960c8f787a25f7ce2497a51202c16e558cf4b61a48eade93135bf96c79aafb12dc6dbabdb10907f84571c886157714a74",
			Size:     123456,
			Path:     "/test/path/to/nothing",
		}
		vers, err := repo.Apply(ctx, &CreateVNFS{
			CommandModel: eventsource.CommandModel{ID: vnfsID},
			Arch:         v2.Arch,
			Checksum:     v2.Checksum,
			Size:         v2.Size,
			Path:         v2.Path,
		})
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if vers != 1 {
			t.Fatalf("Version not 1, %d instead", vers)
		}
		aggregate, err := repo.Load(ctx, vnfsID)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		vnfs := aggregate.(*VNFS)
		if vnfs.State != "Created" {
			t.Fatalf("State not set to created, set to %s instead", vnfs.State)
		}
		if vnfs.Arch != v2.Arch {
			t.Fatalf("Arch mismatch, set to %s instead", vnfs.Arch)
		}
		if vnfs.Checksum != v2.Checksum {
			t.Fatalf("Checksum mismatch, set to %s instead", vnfs.Checksum)
		}
		if vnfs.Size != v2.Size {
			t.Fatalf("Size mismatch, set to %d instead", vnfs.Size)
		}
		if vnfs.Path != v2.Path {
			t.Fatalf("Path mismatch, set to %s instead", vnfs.Path)
		}

		_, err = repo.Apply(ctx, &CreateVNFS{
			CommandModel: eventsource.CommandModel{ID: vnfsID},
			Arch:         v2.Arch,
			Checksum:     v2.Checksum,
			Size:         v2.Size,
			Path:         v2.Path,
		})

		if err == nil {
			t.Fatal("Shoud have failed with already exists error instead")
		}
	})

	t.Run("UpdateVNFS", func(t *testing.T) {
		v2 := VNFS{
			Arch:     "aarch64",
			Checksum: "202c16e558cf4b61a48eade93135bf96c79aafb12dc6dbabdb10907f84571c886157714a740f82992b7693500212e8dd9d7542a5e960c8f787a25f7ce2497a51",
			Size:     654321,
			Path:     "/test/path/to/everything",
		}
		vers, err := repo.Apply(ctx, &UpdateVNFS{
			CommandModel: eventsource.CommandModel{ID: vnfsID},
			Arch:         v2.Arch,
			Checksum:     v2.Checksum,
			Size:         v2.Size,
			Path:         v2.Path,
		})
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if vers != 2 {
			t.Fatalf("Version not 2, %d instead", vers)
		}
		aggregate, err := repo.Load(ctx, vnfsID)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		vnfs := aggregate.(*VNFS)
		if vnfs.State != "Created" {
			t.Fatalf("State not set to created, set to %s instead", vnfs.State)
		}
		if vnfs.Arch != v2.Arch {
			t.Fatalf("Arch mismatch, set to %s instead", vnfs.Arch)
		}
		if vnfs.Checksum != v2.Checksum {
			t.Fatalf("Checksum mismatch, set to %s instead", vnfs.Checksum)
		}
		if vnfs.Size != v2.Size {
			t.Fatalf("Size mismatch, set to %d instead", vnfs.Size)
		}
		if vnfs.Path != v2.Path {
			t.Fatalf("Path mismatch, set to %s instead", vnfs.Path)
		}
	})
	t.Run("DeleteVNFS", func(t *testing.T) {

		vers, err := repo.Apply(ctx, &DeleteVNFS{
			CommandModel: eventsource.CommandModel{ID: vnfsID},
		})
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if vers != 3 {
			t.Fatalf("Version, %d, not incremented on apply", vers)
		}
		aggregate, err := repo.Load(ctx, vnfsID)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		vnfs := aggregate.(*VNFS)
		if vnfs.State != "Deleted" {
			t.Fatalf("State not updated to deleted, set to %s instead", vnfs.State)
		}
		_, err = repo.Apply(ctx, &DeleteVNFS{
			CommandModel: eventsource.CommandModel{ID: vnfsID},
		})
		if err == nil {
			t.Fatal("Should have failed with already deleted error")
		}
	})
}

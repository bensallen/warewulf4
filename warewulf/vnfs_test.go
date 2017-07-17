package warewulf

import (
	"testing"
	"time"

	"github.com/altairsix/eventsource"
)

func testCompareModel(t *testing.T, version int, at time.Time, id string, e eventsource.Model) {
	t.Run("Version equal", func(t *testing.T) {
		if version != e.Version {
			t.Fatalf("Item does not matct input, %v instead", version)
		}
	})
	t.Run("At equal", func(t *testing.T) {
		if at != e.At {
			t.Fatalf("Item does not matct input, %v instead", at)
		}
	})
	t.Run("ID equal", func(t *testing.T) {
		if id != e.ID {
			t.Fatalf("Item does not matct input, %v instead", id)
		}
	})
}
func TestVNFSOn(t *testing.T) {

	t.Run("VNFSCreated", func(t *testing.T) {
		v := VNFS{}

		model := eventsource.Model{ID: "test", Version: 1, At: time.Now()}
		vnfsCreated := VNFSCreated{
			Model: model,
		}
		err := v.On(&vnfsCreated)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
		if v.State != "Created" {
			t.Fatalf("State not Created, %s instead", v.State)
		}
		testCompareModel(t, v.Version, v.CreatedAt, v.ID, model)

	})

	t.Run("VNFSDeleted", func(t *testing.T) {
		v := VNFS{}

		model := eventsource.Model{Version: 10, At: time.Now()}
		vnfsDeleted := VNFSDeleted{
			Model: model,
			State: "Created",
		}
		err := v.On(&vnfsDeleted)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
		testCompareModel(t, v.Version, v.UpdatedAt, "", model)

		if v.State != "Deleted" {
			t.Fatalf("State not Deleted, %s instead", v.State)
		}

	})

}

package redisDb

import (
	"fmt"
	"testing"

	objects "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Objects"
)

func TestQmint(t *testing.T) {
	db, err := NewDBinstance()
	if err != nil {
		t.Error(err)
	}

	var tests = []struct {
		address   string
		contentID string
	}{
		{"0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f", "thing1"},
		{"0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f", "thing1"},
		{"0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f", "thing1"},
		{"0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f", "thing1"},
		{"0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f", "thing1"},
		{"0xEd5E90a45476706A70B9e87Da147988Fdd0e9F6f", "thing1"},
	}
	for _, test := range tests {
		err := db.Qmint(test.address, test.contentID)
		if err != nil {
			t.Error(err)
		}

	}

}

func TestGetState(t *testing.T) {
	db, err := NewDBinstance()
	if err != nil {
		t.Error(err)
	}
	State, err := db.GetState()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(State.HighestFinalizedBlock)
	fmt.Println(State.HighestProcessedBlock)
}
func TestUpdateState(t *testing.T) {
	db, err := NewDBinstance()
	if err != nil {
		t.Error(err)
	}
	State := &objects.State{
		HighestFinalizedBlock: 1,
		HighestProcessedBlock: 2,
	}
	err = db.UpdateState(State)
	if err != nil {
		t.Error(err)
	}
}

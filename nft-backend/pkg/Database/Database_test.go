package redisDb

import (
	"fmt"
	"testing"
)

func TestGetState(t *testing.T) {
	db, err := NewDBinstance()
	if err != nil {
		t.Error(err)
	}
	State, err := db.GetState()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(State.HighestProcessedBlock)
}
func TestUpdateState(t *testing.T) {
	db, err := NewDBinstance()
	if err != nil {
		t.Error(err)
	}

	err = db.UpdateProcessedState(123)
	if err != nil {
		t.Error(err)
	}
	//get the state
	state, err := db.GetState()
	if err != nil {
		t.Error(err)
	}
	_ = state
	//assert.Equal(t, state, State)
}

func TestAddItem(t *testing.T) {
	db, err := NewDBinstance()
	if err != nil {
		t.Error(err)
	}

	var tests = []struct {
		contentID string
	}{
		{"thing1"},
		{"thing2"},
		{"thing3"},
		{"thing4"},
		{"thing5"},
		{"thing6"},
	}
	for _, test := range tests {
		err := db.AddItem(test.contentID)
		if err != nil {
			t.Error(err)
		}
	}
}

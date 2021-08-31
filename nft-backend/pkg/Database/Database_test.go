package redisDb

import (
	"testing"
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

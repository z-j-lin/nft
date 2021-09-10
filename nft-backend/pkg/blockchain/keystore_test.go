package blockchain

import (
	"testing"
)

func TestStoreKs(t *testing.T) {
	var tests = []struct {
		key      string
		passcode string
	}{
		{"d555a6762ae4351f8dd01df268fd8680503dd0ace8ed87df4aedebc0e87a42e8", "pineapple"},
		{"c28403ccef9f382510cce6132c95d3f6f2d072f5cab8c5efe25abdd82bd4b340", "pineapple"},
		{"0e50bcd23b8af3d5a80d822c93cf3bbdcf70cf38d3e6dae79d8280b68b7d5035", "pineapple"},
		{"3d70d4c1d865ae886c9b998e9d0b682f7c2cef4c1b3c0d3cfb15bd222e3d909c", "pineapple"},
		{"9daa330bc31fb43aeccb6e7db8d250ce1a36c5078cf32358fe1f1853113485b7", "pineapple"},
		{"a40cb0f2c62b439366dda0312f113a34788e5d41239e267612e08bc068f02209", "pineapple"},
	}
	for _, test := range tests {
		err := StoreKs(test.passcode, test.key)
		if err != nil {
			t.Error(err)
		}
	}

}

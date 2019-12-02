// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2015-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base58

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// hexToBytes is a wrapper around hex.DecodeString that panics if there is an
// error.  It MUST only be used with hard coded values in the tests.
func hexToBytes(origHex string) []byte {
	buf, err := hex.DecodeString(origHex)
	if err != nil {
		panic(err)
	}
	return buf
}

// TestBase58Coding ensures Decode and Encode produces the expected results for
// both strings and raw byte data converted from hex.
func TestBase58Coding(t *testing.T) {
	tests := []struct {
		decoded []byte
		encoded string
	}{
		// String inputs.
		{[]byte(""), ""},
		{[]byte(" "), "Z"},
		{[]byte("-"), "n"},
		{[]byte("0"), "q"},
		{[]byte("1"), "r"},
		{[]byte("-1"), "4SU"},
		{[]byte("11"), "4k8"},
		{[]byte("abc"), "ZiCa"},
		{[]byte("1234598760"), "3mJr7AoUXx2Wqd"},
		{[]byte("abcdefghijklmnopqrstuvwxyz"), "3yxU3u1igY8WkgtjK92fbJQCd4BZiiT1v25f"},
		{[]byte("00000000000000000000000000000000000000000000000000000000000000"), "3sN2THZeE9Eh9eYrwkvZqNstbHGvrxSAM7gXUXvyFQP8XvQLUqNCS27icwUeDT7ckHm4FUHM2mTVh1vbLmk7y"},

		// Hex inputs.
		{hexToBytes("61"), "2g"},
		{hexToBytes("626262"), "a3gV"},
		{hexToBytes("636363"), "aPEr"},
		{hexToBytes("73696d706c792061206c6f6e6720737472696e67"), "2cFupjhnEsSn59qHXstmK2ffpLv2"},
		{hexToBytes("00eb15231dfceb60925886b67d065299925915aeb172c06647"), "1NS17iag9jJgTHD1VXjvLCEnZuQ3rJDE9L"},
		{hexToBytes("516b6fcd0f"), "ABnLTmg"},
		{hexToBytes("bf4f89001e670274dd"), "3SEo3LWLoPntC"},
		{hexToBytes("572e4794"), "3EFU7m"},
		{hexToBytes("ecac89cad93923c02321"), "EJDM8drfXA6uyA"},
		{hexToBytes("10c8511e"), "Rt5zm"},
		{hexToBytes("00000000000000000000"), "1111111111"},
	}

	for i, test := range tests {
		if res := Encode(test.decoded); res != test.encoded {
			t.Errorf("Encode test #%d failed: got: %q, want: %q", i, res,
				test.encoded)
			continue
		}

		if res := Decode(test.encoded); !bytes.Equal(res, test.decoded) {
			t.Errorf("Decode test #%d failed: got: %q, want: %q", i, res,
				test.decoded)
			continue
		}
	}
}

// TestBase58DecodeInvalid ensures Decode produces an empty result when provided
// with invalid base58 encodings.
func TestBase58DecodeInvalid(t *testing.T) {
	tests := []string{
		"0", "O", "I", "l", "3mJr0", "O3yxU", "3sNI", "4kl8", "0OIl",
		"!@#$%^&*()-_=+~`",
	}

	for i, test := range tests {
		if res := Decode(test); !bytes.Equal(res, []byte("")) {
			t.Errorf("Decode test #%d failed: got: %q, want: %q", i, res, test)
			continue
		}
	}
}

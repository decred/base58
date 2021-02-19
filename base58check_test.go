// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2015-2021 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base58

import (
	"errors"
	"strings"
	"testing"
)

// TestBase58Check ensures CheckDecode and CheckEncode produces the expected
// results for inputs with a given version as well as decoding errors.
func TestBase58Check(t *testing.T) {
	tests := []struct {
		version [2]byte
		decoded string
		encoded string
	}{
		{[2]byte{20, 0}, "", "Axk2WA6L"},
		{[2]byte{20, 0}, " ", "kxg5DGCa1"},
		{[2]byte{20, 0}, "-", "kxhWqwoTY"},
		{[2]byte{20, 0}, "0", "kxhrrcZDw"},
		{[2]byte{20, 0}, "1", "kxhzgbzwe"},
		{[2]byte{20, 0}, "-1", "4M2qnQVfVwu"},
		{[2]byte{20, 0}, "11", "4M2smzp65NR"},
		{[2]byte{20, 0}, "abc", "FmT72s9HXyp6"},
		{[2]byte{20, 0}, "1234598760", "3UFLKR4oYrL1hSX1Eu2W3F"},
		{[2]byte{20, 0}, "abcdefghijklmnopqrstuvwxyz", "2M5VSfthNqvveeGWTcKRgY4Rm258o4ZDKBZGkAQ799jp"},
		{[2]byte{20, 0}, "00000000000000000000000000000000000000000000000000000000000000", "3cmTs9hNQGCVmurJUgS7UokKFYZCCJWvWfYRBCaox5hXDn3Giiy1u9AEKn7vLS8K87BcDr6Ckr4JYRnnaSMRDsB49i3eU"},
	}

	for i, test := range tests {
		// Test encoding.
		gotEncoded := CheckEncode([]byte(test.decoded), test.version)
		if gotEncoded != test.encoded {
			t.Errorf("CheckEncode test #%d failed: got %q, want: %q", i,
				gotEncoded, test.encoded)
			continue
		}

		// Test decoding.
		gotDecoded, version, err := CheckDecode(test.encoded)
		switch {
		case err != nil:
			t.Errorf("CheckDecode test #%d failed with err: %v", i, err)
		case version != test.version:
			t.Errorf("CheckDecode test #%d failed: got version: %x, want: %x",
				i, version, test.version)
		case string(gotDecoded) != test.decoded:
			t.Errorf("CheckDecode test #%d failed: got: %q, want: %q", i,
				gotDecoded, test.decoded)
		}
	}

	// Test the two decoding failure cases:
	// case 1: Checksum error.
	_, _, err := CheckDecode("Axk2WA6M")
	if !errors.Is(err, ErrChecksum) {
		t.Error("Checkdecode test failed, expected ErrChecksum")
	}
	// case 2: invalid formats (string lengths below 6 mean the version byte
	// and/or the checksum bytes are missing).
	for size := 0; size < 6; size++ {
		testString := strings.Repeat("1", size)
		_, _, err = CheckDecode(testString)
		if !errors.Is(err, ErrInvalidFormat) {
			t.Error("Checkdecode test failed, expected ErrInvalidFormat")
		}
	}
}

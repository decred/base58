// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2015-2021 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base58

import (
	"testing"

	"github.com/decred/dcrd/crypto/blake256"
)

// BenchmarkBase58Encode benchmarks how long it takes to perform a base58 encode
// on a typical input.
func BenchmarkBase58Encode(b *testing.B) {
	var input [20]byte
	hash := blake256.Sum256(input[:])
	data := hash[:]

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(data)
	}
}

// BenchmarkBase58Decode benchmarks how long it takes to perform a base58 decode
// on a typical input.
func BenchmarkBase58Decode(b *testing.B) {
	var input [20]byte
	hash := blake256.Sum256(input[:])
	encoded := Encode(hash[:])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(encoded)
	}
}

// BenchmarkCheckDecode benchmarks how long it takes to perform a base58 check
// decode on a typical input.
func BenchmarkCheckDecode(b *testing.B) {
	var input [20]byte
	encoded := CheckEncode(input[:], [2]byte{0x07, 0x3f})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		_, _, err = CheckDecode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCheckEncode benchmarks how long it takes to perform a base58 check
// encode on a typical input.
func BenchmarkCheckEncode(b *testing.B) {
	var input [20]byte
	version := [2]byte{0x07, 0x3f}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckEncode(input[:], version)
	}
}

// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base58

import (
	"bytes"
	"testing"
)

// base58BenchTest describes tests that are used for the various base58
// benchmarks.  It is defined separately so the same tests can easily be used in
// comparison benchmarks.
type base58BenchTest struct {
	name    string // benchmark description
	encoded string // encoded value
	data    []byte // decoded bytes
}

// makeBase58Benches returns a slice of tests that consist of an encoded base58
// string and its associated decoded data for use in the base58 benchmarks.
func makeBase58Benches() []base58BenchTest {
	addrHash := hexToBytes("2789d58cfa0957d206f025c2af056fc8a77cebb0")
	wif := "PmQdMn8xafwaQouk8ngs1CccRCB1ZmsqQxBaxNR4vhQi5a5QB5716"
	extKey := "dpubZ9169KDAEUnyoBhjjmT2VaEodr6pUTDoqCEAeqgbfr2JfkB88BbK77jbTY" +
		"bcYXb2FVz7DKBdW4P618yd51MwF8DjKVopSbS7Lkgi6bowX5w"
	fiftyZeros := bytes.Repeat([]byte{0x00}, 50)
	large := hexToBytes("0102030405060708090a0b0c0d0e0f101112131415161718191a" +
		"1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c" +
		"3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e" +
		"5f606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f80" +
		"8182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2" +
		"a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4" +
		"c5c6c7c8")

	return []base58BenchTest{
		{name: "20_bytes_addrhash", encoded: Encode(addrHash), data: addrHash},
		{name: "53_chars_wif", encoded: wif, data: Decode(wif)},
		{name: "111_chars_extkey", encoded: extKey, data: Decode(extKey)},
		{name: "50_zeros", encoded: Encode(fiftyZeros), data: fiftyZeros},
		{name: "200_bytes_large", encoded: Encode(large), data: large},
	}
}

// BenchmarkBase58Encode benchmarks how long it takes to perform a base58 encode
// on typical inputs.
func BenchmarkBase58Encode(b *testing.B) {
	benches := makeBase58Benches()
	for _, bench := range benches {
		b.Run(bench.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				Encode(bench.data)
			}
		})
	}
}

// BenchmarkBase58Decode benchmarks how long it takes to perform a base58 decode
// on typical inputs.
func BenchmarkBase58Decode(b *testing.B) {
	benches := makeBase58Benches()
	for _, bench := range benches {
		b.Run(bench.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				Decode(bench.encoded)
			}
		})
	}
}

var (
	noElideResult  []byte
	noElideVersion [2]byte
	noElideEncoded string
)

// BenchmarkCheckDecode benchmarks how long it takes to perform a base58 check
// decode on a typical input.
func BenchmarkCheckDecode(b *testing.B) {
	var input [20]byte
	encoded := CheckEncode(input[:], [2]byte{0x07, 0x3f})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		noElideResult, noElideVersion, err = CheckDecode(encoded)
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
		noElideEncoded = CheckEncode(input[:], version)
	}
}

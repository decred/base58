// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2015-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base58

import (
	"bytes"
	"testing"
)

func BenchmarkBase58Encode(b *testing.B) {
	data := bytes.Repeat([]byte{0xff}, 5000)
	b.SetBytes(int64(len(data)))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(data)
	}
}

func BenchmarkBase58Decode(b *testing.B) {
	data := bytes.Repeat([]byte{0xff}, 5000)
	encoded := Encode(data)
	b.SetBytes(int64(len(encoded)))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(encoded)
	}
}

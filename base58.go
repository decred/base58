// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base58

import (
	"math/big"
)

//go:generate go run genalphabet.go

var bigRadix = big.NewInt(58)
var bigZero = big.NewInt(0)

// Decode decodes a modified base58 string to a byte slice.
func Decode(input string) []byte {
	if len(input) == 0 {
		return []byte("")
	}

	// The max possible output size is when a base58 encoding consists of
	// nothing but the alphabet character at index 0 which would result in the
	// same number of bytes as the number of input chars.
	output := make([]byte, len(input))

	// Encode to base256 in reverse order to avoid extra calculations to
	// determine the final output size in favor of just keeping track while
	// iterating.
	var index int
	for _, r := range input {
		// Invalid base58 character.
		val := uint32(b58[r])
		if val == 255 {
			return []byte("")
		}

		// Multiply each byte in the output by 58 and encode to base256 while
		// propagating the carry.
		for i, b := range output[:index] {
			val += uint32(b) * 58
			output[i] = byte(val)
			val >>= 8
		}
		for ; val > 0; val >>= 8 {
			output[index] = byte(val)
			index++
		}
	}

	// Account for the leading zeros in the input.  They are appended since the
	// encoding is happening in reverse order.
	for _, r := range input {
		if r != alphabetIdx0 {
			break
		}

		output[index] = 0
		index++
	}

	// Truncate the output buffer to the actual number of decoded bytes and
	// reverse it since it was calculated in reverse order.
	output = output[:index:index]
	for i := 0; i < index/2; i++ {
		output[i], output[index-1-i] = output[index-1-i], output[i]
	}

	return output
}

// Encode encodes a byte slice to a modified base58 string.
func Encode(b []byte) string {
	x := new(big.Int)
	x.SetBytes(b)

	answer := make([]byte, 0, len(b)*136/100)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		answer = append(answer, alphabet[mod.Int64()])
	}

	// leading zero bytes
	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, alphabetIdx0)
	}

	// reverse
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	return string(answer)
}

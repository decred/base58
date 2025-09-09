// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base58

import (
	"math/bits"
)

//go:generate go run genalphabet.go

// Decode decodes a modified base58 string to a byte slice.
func Decode(input string) []byte {
	if len(input) == 0 {
		return []byte("")
	}

	// Determine the maximum possible output size.
	//
	// Since the conversion is from base58 to base256, the max possible number
	// of bytes of output per input byte, excluding the leading zeros, is
	// log_256(58).  Therefore, the max total output size is the number of
	// leading zero bytes plus ceil(inputSizeMinusLeadingZeros * log_256(58)).
	//
	// Note that log_256(58) ~= 0.7322 < 47/64 which is within 0.3% of the true
	// value and efficient to compute as it only involves division by a power of
	// 2 and thus serves as a good approximation.  So, the calculation below is
	// the integer division equivalent of nlz + ceil(len(input[nlz:]) * 47/64).
	//
	// Finally, in order to avoid additional conditional branches in the
	// conversion from uint32s to bytes, the max output size is rounded up to
	// the next multiple of 4.
	var nlz int
	for i := 0; i < len(input) && input[i] == alphabetIdx0; i++ {
		nlz++
	}
	maxOutputSizeNoLZ := (len(input[nlz:])*47 + 63) / 64
	maxOutputSize := nlz + maxOutputSizeNoLZ
	maxOutputSize = ((maxOutputSize + 3) / 4) * 4
	output := make([]byte, maxOutputSize)

	// The algorithm below performs the calculations with uint32s for better
	// performance and the total number of uint32s is ceil(maxOutputSizeNoLZ /
	// 4).  Note that the leading zeros are skipped here, so the calculation is
	// based on the max output size excluding them.
	//
	// In order to avoid an additional heap allocation for the vast majority of
	// typical cases, use an array on the stack for inputs of up to 120 chars
	// (plus any leading zeros) and fall back to a heap alloc for larger inputs.
	// Note that 120 input chars, excluding leading zeros, equates to a max
	// output size of 92 when applying the same calculations as above.
	//
	// This value was chosen because it provides a good balance between alloc
	// size, speed, and the max chars in the vast majority of inputs decoded in
	// the most common use cases.
	const maxOut32StackAlloc = 92 / 4
	maxOut32Size := (maxOutputSizeNoLZ + 3) / 4
	var out32 []uint32
	if maxOut32Size <= maxOut32StackAlloc {
		var out32Arr [maxOut32StackAlloc]uint32
		out32 = out32Arr[:maxOut32Size]
	} else {
		out32 = make([]uint32, maxOut32Size)
	}

	// Decode to base256 in reverse order to reduce the total number of overall
	// calculations.
	var out32Idx int
	for _, r := range []byte(input[nlz:]) {
		// Invalid base58 character.
		val := uint64(b58[r])
		if val == 255 {
			return []byte("")
		}

		for i, ui32 := range out32[:out32Idx] {
			val += uint64(ui32) * 58
			out32[i] = uint32(val) // nolint:gosec
			val >>= 32
		}
		if val > 0 {
			out32[out32Idx] = uint32(val) // nolint:gosec
			out32Idx++
		}
	}

	// Convert uint32 words to bytes.
	var index int
	for _, ui32 := range out32[:out32Idx] {
		output[index] = byte(ui32)
		output[index+1] = byte(ui32 >> 8)
		output[index+2] = byte(ui32 >> 16)
		output[index+3] = byte(ui32 >> 24)
		index += 4
	}

	// Adjust the output index to the position of the most significant byte and
	// to account for the leading zeros in the input.  They come last since the
	// decoding is happening in reverse order.
	if out32Idx > 0 {
		index -= bits.LeadingZeros32(out32[out32Idx-1]) / 8
	}
	index += nlz

	// Truncate the output buffer to the actual number of decoded bytes and
	// reverse it since it was calculated in reverse order.
	output = output[:index:index]
	for i := 0; i < index/2; i++ {
		output[i], output[index-1-i] = output[index-1-i], output[i]
	}

	return output
}

// Encode encodes a byte slice to a modified base58 string.
func Encode(input []byte) string {
	// Since the conversion is from base256 to base58, the max possible number
	// of bytes of output per input byte is log_58(256) ~= 1.37.  Thus, the max
	// total output size is ceil(len(input) * 137/100).  Rather than worrying
	// about the ceiling, just add one even if it isn't needed since the final
	// output is truncated to the right size at the end.
	output := make([]byte, (len(input)*137/100)+1)

	// Encode to base58 in reverse order to avoid extra calculations to
	// determine the final output size in favor of just keeping track while
	// iterating.
	var index int
	for _, r := range input {
		// Multiply each byte in the output by 256 and encode to base58 while
		// propagating the carry.
		val := uint32(r)
		for i, b := range output[:index] {
			val += uint32(b) << 8
			output[i] = byte(val % 58)
			val /= 58
		}
		for ; val > 0; val /= 58 {
			output[index] = byte(val % 58)
			index++
		}
	}

	// Replace the calculated remainders with their corresponding base58 digit.
	for i, b := range output[:index] {
		output[i] = alphabet[b]
	}

	// Account for the leading zeros in the input.  They are appended since the
	// encoding is happening in reverse order.
	for _, r := range input {
		if r != 0 {
			break
		}

		output[index] = alphabetIdx0
		index++
	}

	// Truncate the output buffer to the actual number of encoded bytes and
	// reverse it since it was calculated in reverse order.
	output = output[:index:index]
	for i := 0; i < index/2; i++ {
		output[i], output[index-1-i] = output[index-1-i], output[i]
	}

	return string(output)
}

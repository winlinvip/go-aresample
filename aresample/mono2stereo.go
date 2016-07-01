// The MIT License (MIT)
//
// Copyright (c) 2016 winlin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// The PCM resample.
package aresample

import "fmt"

// Transform the mono pcm to stereo npcm, where len(npcm)===2*len(pcm).
// @remark the pcm must be s16le(16bits PCM in little-endian).
func PcmS16leMono2Stereo(pcm, npcm []byte) (err error) {
	if len(pcm) == 0 {
		return fmt.Errorf("PCM empty")
	}
	if (len(pcm) % 2) != 0 {
		return fmt.Errorf("PCM size=%v not s16le", len(pcm))
	}
	if len(npcm) != 2*len(pcm) {
		return fmt.Errorf("NPCM size=%v invalid", len(npcm))
	}

	// The value of pcm is v(16bits int little-endian),
	// then the energy e=v*v, when we transform mono to stereo,
	// we must make sure the e is not changed, that is:
	//		v0*v0+v0*v0=e=v*v
	//		2*v0*v0=v*v
	// that is:
	//		v0 = square_root((v*v)/2)
	//		v0 = v*square_root(1/2)
	//		v0 = v*square_root(2)/2
	// we can use fast int transform:
	//		v0 = v * 1.4142135623731 / 2
	//		v0 = v * 0.7071067811865499
	for i:=0; i<len(pcm); i+=2 {
		// 16bits le sample
		v := (int16(pcm[i])) | (int16(pcm[i+1]) << 8)
		// use float32, no need float64, slower than int64
		//		PcmS16leMono2Stereo_int64, loop=8000000, diff=2.188711601s
		//  	PcmS16leMono2Stereo_float64, loop=8000000, diff=2.017123758s
		//		PcmS16leMono2Stereo_float32, loop=8000000, diff=1.749901193s
		v = int16(float32(v) * 0.7071)

		// L
		npcm[i*2] = byte(v)
		npcm[i*2 + 1] = byte(v >> 8)

		// R
		npcm[i*2 + 2] = byte(v)
		npcm[i*2 + 3] = byte(v >> 8)
	}

	return
}

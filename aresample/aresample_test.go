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

package aresample

import (
	"testing"
)

func TestPcmS16leMono2Stereo(t *testing.T) {
	if err := PcmS16leMono2Stereo(make([]byte, 1), make([]byte, 2)); err == nil {
		t.Error("invalid pcm")
		return
	}

	if err := PcmS16leMono2Stereo(make([]byte, 1), make([]byte, 3)); err == nil {
		t.Error("invalid pcm")
		return
	}

	if err := PcmS16leMono2Stereo(make([]byte, 0), nil); err == nil {
		t.Error("invalid pcm")
		return
	}

	b := []byte{0x01, 0x02}
	b0 := make([]byte, len(b) * 2)
	if err := PcmS16leMono2Stereo(b, b0); err != nil {
		t.Error("resample failed, err is", err)
		return
	} else if len(b0) != 2*len(b) {
		t.Error("invalid resample", len(b0))
		return
	} else if b[0] == b0[0] || b[1] == b0[1] || b[0] == b0[2] || b[1] == b0[3] {
		t.Error("invalid resample", b0)
		return
	}
}

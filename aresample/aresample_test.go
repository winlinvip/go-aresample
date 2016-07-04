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
	"bytes"
)

func TestPcmS16leMono2Stereo(t *testing.T) {
	if err := PcmS16leMono2Stereo(make([]byte, 1), make([]byte, 2)); err == nil {
		t.Error("invalid pcm, err is", err)
	}

	if err := PcmS16leMono2Stereo(make([]byte, 1), make([]byte, 3)); err == nil {
		t.Error("invalid pcm, err is", err)
	}

	if err := PcmS16leMono2Stereo(make([]byte, 0), nil); err == nil {
		t.Error("invalid pcm, err is", err)
	}

	b := []byte{0x01, 0x02}
	b0 := make([]byte, len(b) * 2)
	if err := PcmS16leMono2Stereo(b, b0); err != nil {
		t.Error("resample failed, err is", err)
	} else if len(b0) != 2*len(b) {
		t.Error("invalid resample", len(b0))
	} else if bytes.Compare(b0[:2], b0[2:]) != 0 {
		t.Error("invalid resample", b0)
	} else if bytes.Compare(b0[:2], b) == 0 || bytes.Compare(b0[2:], b) == 0 {
		t.Error("invalid resample", b0)
	}
}

func TestPcmS16leResample_Basic(t *testing.T) {
	if _,err := NewPcmS16leResampler(0, 0, 0); err == nil {
		t.Error("invalid resampler")
	}
	if _,err := NewPcmS16leResampler(3, 0, 0); err == nil {
		t.Error("invalid resampler")
	}
	if _,err := NewPcmS16leResampler(1, 0, 0); err == nil {
		t.Error("invalid resampler")
	}
	if _,err := NewPcmS16leResampler(1, 44100, 0); err == nil {
		t.Error("invalid resampler")
	}
	if _,err := NewPcmS16leResampler(1, 0, 44100); err == nil {
		t.Error("invalid resampler")
	}

	pfn0 := func(pcm []byte, f func(pcm, npcm []byte, err error)) {
		if v,err := NewPcmS16leResampler(1, 44100, 44100); err != nil {
			t.Error("invalid resampler")
		} else {
			npcm,err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn0(nil, func(pcm,npcm []byte, err error){
		if err == nil {
			t.Error("invalid pcm, err is", err)
		}
	})
	pfn0([]byte{}, func(pcm,npcm []byte, err error){
		if err == nil {
			t.Error("invalid pcm, err is", err)
		}
	})
	pfn0([]byte{0x00}, func(pcm,npcm []byte, err error){
		if err == nil {
			t.Error("invalid pcm, err is", err)
		}
	})

	pfn1 := func(pcm []byte, f func(pcm, npcm []byte, err error)) {
		if v,err := NewPcmS16leResampler(2, 44100, 44100); err != nil {
			t.Error("invalid resampler")
		} else {
			npcm,err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn1(nil, func(pcm,npcm []byte, err error){
		if err == nil {
			t.Error("invalid pcm, err is", err)
		}
	})
	pfn1([]byte{}, func(pcm,npcm []byte, err error){
		if err == nil {
			t.Error("invalid pcm, err is", err)
		}
	})
	pfn1([]byte{0x00}, func(pcm,npcm []byte, err error){
		if err == nil {
			t.Error("invalid pcm, err is", err)
		}
	})
	pfn1([]byte{0x00,0x01}, func(pcm,npcm []byte, err error){
		if err == nil {
			t.Error("invalid pcm, err is", err)
		}
	})
	pfn1([]byte{0x00,0x01,0x02}, func(pcm,npcm []byte, err error){
		if err == nil {
			t.Error("invalid pcm, err is", err)
		}
	})
}

func TestPcmS16leResample_Mono(t *testing.T) {
	pfn0 := func(pcm []byte, f func(pcm, npcm []byte, err error)) {
		if v,err := NewPcmS16leResampler(1, 44100, 44100); err != nil {
			t.Error("invalid resampler")
		} else {
			npcm,err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn0([]byte{0x00,0x00}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn0([]byte{0x00,0x01, 0x02,0x03}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 4 {
			t.Error("invalid pcm", npcm)
		}
	})

	pfn1 := func(pcm []byte, f func(pcm,npcm []byte, err error)){
		if v,err := NewPcmS16leResampler(1, 44100, 22050); err != nil {
			t.Error("invalid pcm, err is", err)
		} else {
			npcm,err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn1([]byte{0x00,0x01, 0x02,0x03}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn1([]byte{0x00,0x01, 0x02,0x03, 0x04,0x05, 0x06,0x07}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 4 {
			t.Error("invalid pcm", npcm)
		}
	})

	pfn2 := func(pcm []byte, f func(pcm,npcm []byte, err error)){
		if v,err := NewPcmS16leResampler(1, 22050, 44100); err != nil {
			t.Error("invalid pcm, err is", err)
		} else {
			npcm,err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn2([]byte{0x00,0x01}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 4 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn2([]byte{0x00,0x01, 0x02,0x03,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 8 {
			t.Error("invalid pcm", npcm)
		}
	})

	pfn3 := func(pcm []byte, f func(pcm,npcm []byte, err error)){
		if v,err := NewPcmS16leResampler(1, 22050, 32000); err != nil {
			t.Error("invalid pcm, err is", err)
		} else {
			npcm,err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn3([]byte{0x00,0x01}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn3([]byte{0x00,0x01, 0x02,0x03}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 4 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn3([]byte{0x00,0x01, 0x02,0x03, 0x04,0x05}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 8 {
			t.Error("invalid pcm", npcm)
		}
	})

	pfn4 := func(pcm []byte, f func(pcm,npcm []byte, err error)){
		if v,err := NewPcmS16leResampler(1, 32000, 22050); err != nil {
			t.Error("invalid pcm, err is", err)
		} else {
			npcm,err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn4([]byte{0x00,0x01}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 0 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn4([]byte{0x00,0x01, 0x02,0x03}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn4([]byte{0x00,0x01, 0x02,0x03, 0x04,0x05}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 4 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn4([]byte{0x00,0x01, 0x02,0x03, 0x04,0x05, 0x06,0x07}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 4 {
			t.Error("invalid pcm", npcm)
		}
	})
}

func TestPcmS16leResample_MonoSamples(t *testing.T) {
	var err error
	var v ResampleSampleRate
	if v,err = NewPcmS16leResampler(1, 44100, 22050); err != nil {
		t.Error("invalid resampler")
	}

	var npcm []byte
	if npcm,err = v.Resample([]byte{0x00,0x01}); err != nil {
		t.Error("invalid pcm, err is", err)
	}
	if len(npcm) != 0 {
		t.Error("invalid pcm", npcm)
	}

	if npcm,err = v.Resample([]byte{0x00,0x02}); err != nil {
		t.Error("invalid pcm, err is", err)
	}
	if len(npcm) != 2 {
		t.Error("invalid pcm", npcm)
	}

	if npcm,err = v.Resample([]byte{0x00,0x02, 0x03,0x04}); err != nil {
		t.Error("invalid pcm, err is", err)
	}
	if len(npcm) != 2 {
		t.Error("invalid pcm", npcm)
	}

	if npcm,err = v.Resample([]byte{0x00,0x02}); err != nil {
		t.Error("invalid pcm, err is", err)
	}
	if len(npcm) != 0 {
		t.Error("invalid pcm", npcm)
	}

	//////////////////////////////////////////////////////////////////////
	if v,err = NewPcmS16leResampler(1, 22050, 32000); err != nil {
		t.Error("invalid resampler")
	}
	if npcm,err = v.Resample([]byte{0x00,0x01}); err != nil {
		t.Error("invalid pcm, err is", err)
	}
	if len(npcm) != 2 {
		t.Error("invalid pcm", npcm)
	}

	if npcm,err = v.Resample([]byte{0x00,0x02}); err != nil {
		t.Error("invalid pcm, err is", err)
	}
	if len(npcm) != 2 {
		t.Error("invalid pcm", npcm)
	}

	if npcm,err = v.Resample([]byte{0x00,0x03}); err != nil {
		t.Error("invalid pcm, err is", err)
	}
	if len(npcm) != 4 {
		t.Error("invalid pcm", npcm)
	}

	if npcm,err = v.Resample([]byte{0x00,0x04}); err != nil {
		t.Error("invalid pcm, err is", err)
	}
	if len(npcm) != 2 {
		t.Error("invalid pcm", npcm)
	}
}

func TestPcmS16leResample_Stereo(t *testing.T) {
	pfn0 := func(pcm []byte, f func(pcm, npcm []byte, err error)) {
		if v,err := NewPcmS16leResampler(2, 44100, 44100); err != nil {
			t.Error("invalid resampler")
		} else {
			npcm,err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn0([]byte{0x00,0x00,0x00,0x00}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn0([]byte{0x00,0x01,0x00,0x01, 0x02,0x03,0x02,0x03}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 4*2 {
			t.Error("invalid pcm", npcm)
		}
	})

	pfn1 := func(pcm []byte, f func(pcm,npcm []byte, err error)){
		if v,err := NewPcmS16leResampler(2, 44100, 22050); err != nil {
			t.Error("invalid pcm, err is", err)
		} else {
			npcm,err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn1([]byte{0x00,0x01,0x00,0x01, 0x02,0x03,0x02,0x03,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn1([]byte{0x00,0x01,0x00,0x01, 0x02,0x03,0x02,0x03, 0x04,0x05,0x04,0x05, 0x06,0x07,0x06,0x07,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 4*2 {
			t.Error("invalid pcm", npcm)
		}
	})

	pfn2 := func(pcm []byte, f func(pcm,npcm []byte, err error)){
		if v,err := NewPcmS16leResampler(2, 22050, 44100); err != nil {
			t.Error("invalid pcm, err is", err)
		} else {
			npcm,err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn2([]byte{0x00,0x01,0x00,0x01,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 4*2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn2([]byte{0x00,0x01,0x00,0x01, 0x02,0x03,0x02,0x03,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 8*2 {
			t.Error("invalid pcm", npcm)
		}
	})

	pfn3 := func(pcm []byte, f func(pcm,npcm []byte, err error)){
		if v,err := NewPcmS16leResampler(2, 22050, 32000); err != nil {
			t.Error("invalid pcm, err is", err)
		} else {
			npcm,err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn3([]byte{0x00,0x01,0x00,0x01,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn3([]byte{0x00,0x01,0x00,0x01, 0x02,0x03,0x02,0x03,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 4*2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn3([]byte{0x00,0x01,0x00,0x01, 0x02,0x03,0x02,0x03, 0x04,0x05,0x04,0x05,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 8*2 {
			t.Error("invalid pcm", npcm)
		}
	})

	pfn4 := func(pcm []byte, f func(pcm,npcm []byte, err error)){
		if v,err := NewPcmS16leResampler(2, 32000, 22050); err != nil {
			t.Error("invalid pcm, err is", err)
		} else {
			npcm,err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn4([]byte{0x00,0x01,0x00,0x01,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 0*2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn4([]byte{0x00,0x01,0x00,0x01, 0x02,0x03,0x02,0x03,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn4([]byte{0x00,0x01,0x00,0x01, 0x02,0x03,0x02,0x03, 0x04,0x05,0x04,0x05,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 4*2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn4([]byte{0x00,0x01,0x00,0x01, 0x02,0x03,0x02,0x03, 0x04,0x05,0x04,0x05, 0x06,0x07,0x06,0x07,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 4*2 {
			t.Error("invalid pcm", npcm)
		}
	})
}

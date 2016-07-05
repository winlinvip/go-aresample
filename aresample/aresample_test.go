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
	"fmt"
)

func TestSpline(t *testing.T) {
	if spline(nil,nil,nil,nil) == nil {
		t.Error("invalid data")
	}

	if spline([]float64{},nil,nil,nil) == nil {
		t.Error("invalid data")
	}
	if spline([]float64{1},nil,nil,nil) == nil {
		t.Error("invalid data")
	}
	if spline([]float64{1,2},nil,nil,nil) == nil {
		t.Error("invalid data")
	}
	if spline([]float64{1,2,3},nil,nil,nil) == nil {
		t.Error("invalid data")
	}
	if spline([]float64{1,2,3,4,5},nil,nil,nil) == nil {
		t.Error("invalid data")
	}

	if spline([]float64{1,2,3,4},[]float64{},nil,nil) == nil {
		t.Error("invalid data")
	}
	if spline([]float64{1,2,3,4},[]float64{1},nil,nil) == nil {
		t.Error("invalid data")
	}
	if spline([]float64{1,2,3,4},[]float64{1,2},nil,nil) == nil {
		t.Error("invalid data")
	}
	if spline([]float64{1,2,3,4},[]float64{1,2,3},nil,nil) == nil {
		t.Error("invalid data")
	}
	if spline([]float64{1,2,3,4},[]float64{1,2,3,5},nil,nil) == nil {
		t.Error("invalid data")
	}

	if spline([]float64{1,2,3,4},[]float64{1,2,3,4},[]float64{},nil) == nil {
		t.Error("invalid data")
	}
	if spline([]float64{1,2,3,4},[]float64{1,2,3,4},[]float64{},[]float64{}) == nil {
		t.Error("invalid data")
	}

	if spline([]float64{1,2,3,4},[]float64{1,2,3,4},[]float64{1.5},[]float64{0}) != nil {
		t.Error("invalid data")
	}

	xi := []float64{1,2,4,5}
	x0,x1,x2,x3 := xi[0],xi[1],xi[2],xi[3]
	h0,h1,h2,l1,u1,l2,u2 := spline_lu(xi)
	if h0 != 1 || h1 != 2 || h2 != 1 {
		t.Error("invalid h", []float64{h0,h1,h2})
	}
	if l1 != float64(1/3.0) || u1 != float64(2/3.0) {
		t.Error("invalid l1/u1", []float64{l1,u1})
	}
	if l2 != float64(2/3.0) || u2 != float64(1/3.0) {
		t.Error("invalid l2/u2", []float64{l2,u2})
	}

	yi := []float64{1,3,4,2}
	y0,y1,y2,y3 := yi[0],yi[1],yi[2],yi[3]
	c1,c2 := spline_c1(yi,h0,h1), spline_c2(yi,h1,h2)
	if c1 != -3 || c2 != -5 {
		t.Error("invalid c", []float64{c1,c2})
	}

	m1,m2 := spline_m1(c1,c2,u1,l2),spline_m2(c1,c2,u1,l2)
	if m1 != -3/4.0 || m2 != -9/4.0 {
		t.Error("invalid m", []float64{m1,m2})
	}

	x := 1.0
	v := spline_z0(m1,h0,x0,x1,y0,y1,x)
	ev := -1.0*(x*x*x)/8 + 3.0*(x*x)/8 + 7.0*x/4 - 1 // 1
	if v != ev {
		t.Error("z0(1.0) ev is", ev, "and v is", v)
	}

	x = 2.0
	v = spline_z0(m1,h0,x0,x1,y0,y1,x)
	ev = -1.0*(x*x*x)/8 + 3.0*(x*x)/8 + 7.0*x/4 - 1 // 3
	if v != ev {
		t.Error("z0(2.0) ev is", ev, "and v is", v)
	}

	x = 2.0
	v = spline_z1(m1,m2,h1,x1,x2,y1,y2,x)
	ev = -1.0*(x*x*x)/8 + 3*(x*x)/8 + 7.0*x/4 - 1 // 3
	if v != ev {
		t.Error("z1(2.0) ev is", ev, "and v is", v)
	}

	x = 4.0
	v = spline_z1(m1,m2,h1,x1,x2,y1,y2,x)
	ev = -1.0*(x*x*x)/8 + 3*(x*x)/8 + 7.0*x/4 - 1 // 4
	if v != ev {
		t.Error("z1(4.0) ev is", ev, "and v is", v)
	}

	x = 4.0
	v = spline_z2(m2,h2,x2,x3,y2,y3,x)
	ev = 3.0*(x*x*x)/8 - 45.0*(x*x)/8 + 103.0*x/4 - 33.0 // 4
	if v != ev {
		t.Error("z2(4.0) ev is", ev, "and v is", v)
	}

	x = 5.0
	v = spline_z2(m2,h2,x2,x3,y2,y3,x)
	ev = 3.0*(x*x*x)/8 - 45.0*(x*x)/8 + 103.0*x/4 - 33.0 // 2
	if v != ev {
		t.Error("z2(5.0) ev is", ev, "and v is", v)
	}

	x = 1.5
	v = spline_z0(m1,h0,x0,x1,y0,y1,x)
	ev = -1.0*(x*x*x)/8 + 3.0*(x*x)/8 + 7.0*x/4 - 1 // 2.046875
	if v != ev {
		t.Error("z0(1.5) ev is", ev, "and v is", v)
	}

	x = 2.5
	v = spline_z1(m1,m2,h1,x1,x2,y1,y2,x)
	ev = -1.0*(x*x*x)/8 + 3*(x*x)/8 + 7.0*x/4 - 1 // 3.765625
	if v != ev {
		t.Error("z1(2.5) ev is", ev, "and v is", v)
	}

	x = 4.5
	v = spline_z2(m2,h2,x2,x3,y2,y3,x)
	ev = 3.0*(x*x*x)/8 - 45.0*(x*x)/8 + 103.0*x/4 - 33.0 // 3.140625
	if v != ev {
		t.Error("z2(4.5) ev is", ev, "and v is", v)
	}
}

func TestSpline_Resample(t *testing.T) {
	xi := []float64{1,2,4,5}
	yi := []float64{1,3,4,2}
	xo := []float64{1, 1.5, 2, 2.5, 4, 4.5, 5}
	yo := make([]float64, len(xo))
	if err := spline(xi,yi,xo,yo); err != nil {
		t.Error("spline failed, err is", err)
	} else if len(yo) != len(xo) {
		t.Error("invalid yo", yo)
	} else if yo[0] != 1 || yo[1] != 2.046875 || yo[2] != 3 || yo[3] != 3.765625 || yo[4] != 4 || yo[5] != 3.140625 {
		t.Error("invalid yo", yo)
	}

	xi = []float64{0,1,2,3}
	yi = []float64{17,9,33,5}
	xo = []float64{0, 0.5, 1, 1.5, 2, 2.5, 3}
	yo = make([]float64, len(xo))
	if err := spline(xi,yi,xo,yo); err != nil {
		t.Error("spline failed, err is", err)
	} else if len(yo) != len(xo) {
		t.Error("invalid yo", yo)
	} else if yo[0] != 17 || yo[2] != 9 || yo[4] != 33 || yo[6] != 5 {
		t.Error("invalid yo", yo)
	} else if yo[1] != 8.5 || yo[3] != 22.5 || yo[5] != 25 {
		t.Error("invalid yo", yo)
	}

	xi = []float64{1,2,3,4}
	yi = []float64{17,9,33,5}
	xo = []float64{1, 1.5, 2, 2.5, 3, 3.5, 4}
	yo = make([]float64, len(xo))
	if err := spline(xi,yi,xo,yo); err != nil {
		t.Error("spline failed, err is", err)
	} else if len(yo) != len(xo) {
		t.Error("invalid yo", yo)
	} else if yo[0] != 17 || yo[2] != 9 || yo[4] != 33 || yo[6] != 5 {
		t.Error("invalid yo", yo)
	} else if yo[1] != 8.5 || yo[3] != 22.5 || yo[5] != 25 {
		t.Error("invalid yo", yo)
	}
}

func TestPcmS16leResample_Kernel(t *testing.T) {
	if npcm := resampler_init_channel([]byte{0x00,0x01, 0x02,0x03}, 1, 1); len(npcm) != 0 {
		t.Error("invalid channel", len(npcm))
	}

	if npcm := resampler_init_channel([]byte{0x00,0x01, 0x02,0x03}, 1, 0); len(npcm) != 2 {
		t.Error("invalid channel", len(npcm))
	} else if npcm[0] != 0x0100 || npcm[1] != 0x0302 {
		t.Error("invalid channel", npcm)
	}

	if npcm := resampler_init_channel([]byte{0x00,0x01, 0x02,0x03}, 2, 0); len(npcm) != 1 {
		t.Error("invalid channel", len(npcm))
	} else if npcm[0] != 0x0100 {
		t.Error("invalid channel", npcm)
	}

	if npcm := resampler_init_channel([]byte{0x00,0x01, 0x02,0x03}, 2, 1); len(npcm) != 1 {
		t.Error("invalid channel", len(npcm))
	} else if npcm[0] != 0x0302 {
		t.Error("invalid channel", npcm)
	}

	ipcm,ds,isr,osr := []int16{17,9,33,5},0,16000,32000
	nbM,nbP,nds,x := resample_calc_x(ipcm,ds,isr,osr)
	if nbM != 4 || nbP != 7 || nds != 1 {
		t.Error("invalid info", []int{nbM,nbP,nds})
	}
	if len(x) != 7 {
		t.Error("invalid postion", len(x), x)
	} else if x[0] != 0 || x[1] != 0.5 || x[2] != 1 || x[3] != 1.5 || x[4] != 2 || x[5] != 2.5 || x[6] != 3 {
		t.Error("invalid position", x)
	}

	if yo,err := resample_channel(ipcm, x); len(yo) != 7 || err != nil {
		t.Error("invalid yo", len(yo), yo)
	} else if yo[0] != 17 || yo[2] != 9 || yo[4] != 33 || yo[6] != 5 {
		t.Error("invalid yo", yo)
	} else if yo[1] != 8 || yo[3] != 22 || yo[5] != 25 {
		t.Error("invalid yo", yo)
	}

	if npcm := resample_merge([]int16{0x01},nil); len(npcm) != 2 {
		t.Error("invalid merged data", len(npcm))
	}
	if npcm := resample_merge([]int16{0x01},[]int16{0x02}); len(npcm) != 4 {
		t.Error("invalid merged data", len(npcm))
	}
}

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
		if v,err := NewPcmS16leResampler(1, 44100, 22010); err != nil {
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
	pfn0([]byte{0x00,0x00,0x00,0x00,0x00,0x00,0x00}, func(pcm,npcm []byte, err error){
		if err == nil {
			t.Error("invalid pcm, err is", err)
		}
	})
	pfn0([]byte{0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
	})

	pfn1 := func(pcm []byte, f func(pcm, npcm []byte, err error)) {
		if v,err := NewPcmS16leResampler(2, 44100, 22010); err != nil {
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
	pfn1([]byte{0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,}, func(pcm,npcm []byte, err error){
		if err == nil {
			t.Error("invalid pcm, err is", err)
		}
	})
	pfn1([]byte{0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
	})
}

func TestPcmS16leResample_Mono(t *testing.T) {
	pfn0 := func(pcm []byte, f func(pcm, npcm []byte, err error)) {
		if v, err := NewPcmS16leResampler(1, 44100, 44100); err != nil {
			t.Error("invalid resampler")
		} else {
			npcm, err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn0([]byte{0x00, 0x00}, func(pcm, npcm []byte, err error) {
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn0([]byte{0x00, 0x01, 0x02, 0x03}, func(pcm, npcm []byte, err error) {
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 4 {
			t.Error("invalid pcm", npcm)
		}
	})

	pfn1 := func(pcm []byte, f func(pcm, npcm []byte, err error)) {
		if v, err := NewPcmS16leResampler(1, 44100, 22050); err != nil {
			t.Error("invalid pcm, err is", err)
		} else {
			npcm, err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn1([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, func(pcm, npcm []byte, err error) {
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2 * 2 {
			t.Error("invalid pcm", npcm)
		}
	})

	pfn2 := func(pcm []byte, f func(pcm, npcm []byte, err error)) {
		if v, err := NewPcmS16leResampler(1, 22050, 44100); err != nil {
			t.Error("invalid pcm, err is", err)
		} else {
			npcm, err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn2([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, func(pcm, npcm []byte, err error) {
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2 * (8 - 1) {
			t.Error("invalid pcm", len(npcm), npcm)
		}
	})

	pfn3 := func(pcm []byte, f func(pcm, npcm []byte, err error)) {
		if v, err := NewPcmS16leResampler(1, 22050, 32000); err != nil {
			t.Error("invalid pcm, err is", err)
		} else {
			npcm, err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn3([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, func(pcm, npcm []byte, err error) {
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2 * (5 - 1) {
			t.Error("invalid pcm", len(npcm), npcm)
		}
	})

	pfn4 := func(pcm []byte, f func(pcm, npcm []byte, err error)) {
		if v, err := NewPcmS16leResampler(1, 32000, 22050); err != nil {
			t.Error("invalid pcm, err is", err)
		} else {
			npcm, err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn4([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, func(pcm, npcm []byte, err error) {
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2 * 2 {
			t.Error("invalid pcm", npcm)
		}
	})
}

func TestPcmS16leResample_MonoFFMPEG(t *testing.T) {
	// FFMPEG 16KHZ to 32KHZ
	pfn5 := func(pcm []byte, f func(pcm,npcm []byte, err error)){
		if v,err := NewPcmS16leResampler(1, 16000, 32000); err != nil {
			t.Error("invalid pcm, err is", err)
		} else {
			npcm,err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}

	pfn5([]byte{
		0xce,0x0c, 0x6e,0x0d, 0xef,0x0e, 0x93,0x13,
	}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*(8) {
			t.Error("invalid pcm", len(npcm), npcm)
		} else {
			if v := int16(npcm[2])|(int16(npcm[3])<<8); v != (3391-36) {
				t.Error("invalid v", v, npcm)
			}
			if v := int16(npcm[6])|(int16(npcm[7])<<8); v != (3539+14) {
				t.Error("invalid v", v, npcm)
			}
			if v := int16(npcm[10])|(int16(npcm[11])<<8); v != (4345-3) {
				t.Error("invalid v", v, npcm)
			}
			if v := int16(npcm[14])|(int16(npcm[15])<<8); v != (5642-43) {
				t.Error("invalid v", v, npcm)
			}
		}
	})
	return

	pfn5([]byte{
		0xce,0x0c, 0x6e,0x0d, 0xef,0x0e, 0x93,0x13, 0xe7,0x17, 0x1f,0x1b, 0x58,0x1f, 0xa3,0x1d,
		0xbb,0x10, 0x41,0x02, 0x5d,0xfe, 0xc3,0x00, 0x56,0x01,
	}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*(26-1) {
			t.Error("invalid pcm", len(npcm), npcm)
		} else {
			if npcm[2] != (0x3f-36) || npcm[3] != 0x0d {
				t.Error("invalid p[2][3]", npcm[2], npcm[3], fmt.Sprintf("%x", npcm[2:]))
			}
			if npcm[6] != (0xd3+14) || npcm[7] != 0x0d {
				t.Error("invalid p[6][7]", npcm[6], npcm[7], fmt.Sprintf("%x", npcm[6:]))
			}
			if npcm[10] != (0xf9-3) || npcm[11] != 0x10 {
				t.Error("invalid p[10][11]", npcm[10], npcm[11], fmt.Sprintf("%x", npcm[10:]))
			}
		}
	})
}

func TestPcmS16leResample_MonoFrames(t *testing.T) {
	pcm := []byte{
		0x3d, 0xdc, 0x20, 0x13, 0xf3, 0x00, 0x00, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x80, 0x61, 0xbf, 0xff, 0xf7, 0x6e, 0x3a, 0xa2, 0xff, 0xff, 0xf6, 0x01, 0x37, 0xd7, 0x49, 0x9d, 0xf7, 0xdf, 0xf2,
		0x6f, 0x63, 0xda, 0xcd, 0xa4, 0x18, 0x47, 0xe6, 0x19, 0x47, 0x96, 0xf4, 0x32, 0xe6, 0x21, 0x26, 0x8d, 0x12, 0xee,
		0x6d, 0x7c, 0x5b, 0x3f, 0x3c, 0x5f, 0xd7, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0x6a, 0xba,
		0xb8, 0x4a, 0x74, 0x9a, 0xb4, 0x2d, 0xd8, 0xd8, 0xe1, 0xc3, 0x47, 0x25, 0xe8, 0x05, 0xa3, 0xbb, 0xd7, 0x66, 0x3a,
		0x1b, 0xb7, 0xa4, 0x7d, 0xa2, 0xab, 0xfe, 0xd9, 0x08, 0x2c, 0x47,
	}

	var err error
	var r ResampleSampleRate
	if r,err = NewPcmS16leResampler(1, 16000, 32000); err != nil {
		t.Error("aresample failed, err is", err)
		return
	}

	var npcm []byte
	if npcm,err = r.Resample(pcm); err != nil {
		t.Error("aresample failed, err is", err)
		return
	}
	if len(npcm) != 2*(106-1) {
		t.Error("invalid pcm", len(npcm), npcm)
	}

	if npcm,err = r.Resample(pcm); err != nil {
		t.Error("aresample failed, err is", err)
		return
	}
	if len(npcm) != 2*(106) {
		t.Error("invalid pcm", len(npcm), npcm)
	}

	if npcm,err = r.Resample(pcm); err != nil {
		t.Error("aresample failed, err is", err)
		return
	}
	if len(npcm) != 2*(106 - 1) {
		t.Error("invalid pcm", len(npcm), npcm)
	}
}

func TestPcmS16leResample_Stereo(t *testing.T) {
	pfn0 := func(pcm []byte, f func(pcm, npcm []byte, err error)) {
		if v, err := NewPcmS16leResampler(2, 44100, 44100); err != nil {
			t.Error("invalid resampler")
		} else {
			npcm, err := v.Resample(pcm)
			f(pcm, npcm, err)
		}
	}
	pfn0([]byte{0x00,0x00,0x00,0x00,}, func(pcm, npcm []byte, err error) {
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn0([]byte{0x00,0x01,0x02,0x03, 0x00,0x01,0x02,0x03,}, func(pcm, npcm []byte, err error) {
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*4 {
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
	pfn1([]byte{0x00,0x01,0x02,0x03, 0x04,0x05,0x06,0x07, 0x08,0x09,0x0a,0x0b, 0x0c,0x0d,0x0e,0x0f,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*2*2 {
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
	pfn2([]byte{0x00,0x01,0x02,0x03, 0x04,0x05,0x06,0x07, 0x08,0x09,0x0a,0x0b, 0x0c,0x0d,0x0e,0x0f,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*2*(8-1) {
			t.Error("invalid pcm", len(npcm), npcm)
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
	pfn3([]byte{0x00,0x01,0x02,0x03, 0x04,0x05,0x06,0x07, 0x08,0x09,0x0a,0x0b, 0x0c,0x0d,0x0e,0x0f,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*2*(5-1) {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn3([]byte{0x00,0x01,0x00,0x01, 0x02,0x03,0x02,0x03, 0x04,0x05,0x04,0x05, 0x06,0x07,0x06,0x07,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*2*(5-1) {
			t.Error("invalid pcm", npcm)
		}
		for i:=0;i<len(npcm);i+=4 {
			if npcm[i] != npcm[i+2] || npcm[i+1] != npcm[i+3] {
				t.Error("invalid pcm at", i, npcm[i:])
			}
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
	pfn4([]byte{0x00,0x01,0x02,0x03, 0x04,0x05,0x06,0x07, 0x08,0x09,0x0a,0x0b, 0x0c,0x0d,0x0e,0x0f,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*2*2 {
			t.Error("invalid pcm", npcm)
		}
	})
	pfn4([]byte{0x00,0x01,0x00,0x01, 0x02,0x03,0x02,0x03, 0x04,0x05,0x04,0x05, 0x06,0x07,0x06,0x07,}, func(pcm,npcm []byte, err error){
		if err != nil {
			t.Error("invalid pcm, err is", err)
		}
		if len(npcm) != 2*2*2 {
			t.Error("invalid pcm", npcm)
		}
		for i:=0;i<len(npcm);i+=4 {
			if npcm[i] != npcm[i+2] || npcm[i+1] != npcm[i+3] {
				t.Error("invalid pcm at", i, npcm[i:])
			}
		}
	})
}

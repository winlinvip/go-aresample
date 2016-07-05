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

import (
	"fmt"
)

type ResampleSampleRate interface {
	// Resample the pcm to npcm, which contains len(pcm)/2 samples.
	// @remark each sample is 16bits in short int.
	// @reamrk pcm must align to 2, atleast 4 samples.
	Resample(pcm []byte) (npcm []byte, err error)
}

// sample rate resampler.
type srResampler struct {
	channels        int    // Channels, L or LR
	sampleRate      int    // Transform from this sample rate.
	nSampleRate     int    // Transform to this sample rate.

	// Assume:
	//		M is input samples,
	//		isr is input sample rate,
	//		N is output samples,
	//		P is the actual output samples,
	//		osr is output sample rate,
	//		DS is the delta samples
	// Then:
	//		N = (M+DS)*osr/isr
	//		P = int(N) - int(N)%2
	//		DS = M - int(P*isr/osr)
	// @remark Initialize the DS to 0
	deltaSamples 	int // The delta samples of previous resample.
}

// Create resampler to transform pcm
// from sampleRate to nSampleRate, where pcm contains number of channels
// @remark each sample is 16bits in short int.
func NewPcmS16leResampler(channels, sampleRate int, nSampleRate int) (ResampleSampleRate, error) {
	if channels < 1 || channels > 2 {
		return nil,fmt.Errorf("invalid channels=%v", channels)
	}
	if sampleRate <= 0 {
		return nil,fmt.Errorf("invalid sampleRate=%v", sampleRate)
	}
	if nSampleRate <= 0 {
		return nil,fmt.Errorf("invalid nSampleRate=%v", nSampleRate)
	}

	v := &srResampler{
		channels: channels,
		sampleRate: sampleRate,
		nSampleRate: nSampleRate,
	}

	return v,nil
}

func (v *srResampler) Resample(pcm []byte) (npcm []byte, err error) {
	if len(pcm) == 0 {
		return nil,fmt.Errorf("empty pcm")
	}
	if (len(pcm)%(2*v.channels)) != 0 {
		return nil,fmt.Errorf("invalid pcm, should mod(%v)", 2*v.channels)
	}

	if v.sampleRate == v.nSampleRate {
		return pcm[:],nil
	}

	// Atleast 4samples when not init.
	if nbSamles := len(pcm) / 2 / v.channels; nbSamles < 4 {
		return nil,fmt.Errorf("invalid pcm, atleast 4samples, actual %vsamples", nbSamles)
	}

	// Convert pcm to int16 values
	ipcmLeft := resampler_init_channel(pcm, v.channels, 0)
	ipcmRight := resampler_init_channel(pcm, v.channels, 1)
	if ipcmRight != nil && len(ipcmLeft) != len(ipcmRight) {
		return nil,fmt.Errorf("invalid pcm, L%v!=%v", len(ipcmLeft), len(ipcmRight))
	}

	// Resample all channels
	ds := v.deltaSamples
	_,_,nds,x := resample_calc_x(ipcmLeft,ds,v.sampleRate,v.nSampleRate)
	var opcmLeft []int16
	if opcmLeft,err = resample_channel(ipcmLeft, x); err != nil {
		return nil,err
	}
	v.deltaSamples = nds

	var opcmRight []int16
	if ipcmRight != nil {
		_,_,nds,x = resample_calc_x(ipcmRight,ds,v.sampleRate,v.nSampleRate)
		if opcmRight,err = resample_channel(ipcmLeft, x); err != nil {
			return nil,err
		}
	}

	// Convert int16 samples to bytes.
	npcm = resample_merge(opcmLeft, opcmRight)

	return
}

// merge left and right(can be nil).
func resample_merge(left,right []int16) (npcm []byte) {
	npcm = []byte{}
	for i:=0; i<len(left); i++ {
		v := left[i]
		npcm = append(npcm, byte(v))
		npcm = append(npcm, byte(v >> 8))

		if right != nil {
			v = right[i]
			npcm = append(npcm, byte(v))
			npcm = append(npcm, byte(v >> 8))
		}
	}
	return
}

// x is the position of output pcm
func resample_channel(ipcm []int16, x []float32) (opcm []int16, err error) {
	xi := make([]float64, 4)
	yi := make([]float64, 4)

	var p int
	for i:=0; i<len(ipcm); i+= 3 {
		// Complete when left only 1byte
		if i + 1 == len(ipcm) {
			break
		}

		// Skip back to always use 4samples as input.
		if i + 3 >= len(ipcm) {
			i = len(ipcm) - 4
		}

		for j:=0; j<4; j++ {
			xi[j] = float64(i+j)
			yi[j] = float64(ipcm[i+j])
		}

		var xo []float64
		x0,x3 := float32(xi[0]),float32(xi[3])
		for j:=p; j<len(x); j++ {
			if x[j] >= x0 && x[j] <= x3+1 {
				xo = append(xo, float64(x[j]))
			}
		}
		p += len(xo)

		// Completed for this block
		if len(xo) == 0 {
			continue
		}

		yo := make([]float64, len(xo))
		if err = spline(xi,yi,xo,yo); err != nil {
			return nil,err
		}

		for _,v := range yo {
			opcm = append(opcm, int16(v))
		}
	}

	if len(x) != len(opcm) {
		return nil,fmt.Errorf("invalid yo=%v, x=%v", opcm, len(x)*2)
	}

	return
}

// nbM is M, the input samples.
// nbP is P, the actual output samples.
// nds is DS, the updated DS.
// x is the x positions of new samples.
func resample_calc_x(ipcm []int16, ds,isr,osr int) (nbM,nbP,nds int, x []float32) {
	nbM = len(ipcm)
	nbN := (nbM + ds) * osr / isr
	nbP = int(nbN) - int(nbN)%2

	step := float32(nbM) / float32(nbP)
	for i:=float32(0.0); i<float32(nbM); i+=step {
		x = append(x, i)
	}

	nbP = len(x)
	nds = nbM - nbP*isr/osr
fmt.Println(fmt.Sprintf("ds=%v, m=%v, p=%v, nds=%v", ds, nbM, nbP, nds))
	return
}

// resampler_init_channel([]byte{...}, 1, 0)
// resampler_init_channel([]byte{...}, 2, 0)
// resampler_init_channel([]byte{...}, 2, 1)
func resampler_init_channel(pcm []byte, channels, channel int) (ipcm []int16) {
	if channel >= channels {
		return
	}

	ipcm = []int16{}
	for i:=2*channel; i<len(pcm); i+=2*channels {
		// 16bits le sample
		v := (int16(pcm[i])) | (int16(pcm[i + 1]) << 8)
		ipcm = append(ipcm, v)
	}

	return
}

// xi must be [x0, x1, x2, x3] which is [1, 2, 3, 4]
// yi must be [y0, y1, y2, y3] which corresponding to xi
// xo the output insert position of x, must in [x0, x3]
// yo is the inserted value corresponding to xo
// For example:
//		spline([1,2,3,4], [7,9,2,5], [1.5,2.5,3.5], [?,?,?])
// which will fill the yo with values.
func spline(xi,yi,xo,yo []float64) (err error) {
	if len(xi) != 4 {
		return fmt.Errorf("invalid xi")
	}
	if len(yi) != 4 {
		return fmt.Errorf("invalid yi")
	}
	if len(xo) == 0 {
		return fmt.Errorf("invalid xo")
	}
	if len(yo) != len(xo) {
		return fmt.Errorf("invalid yo")
	}

	x0,x1,x2,x3 := xi[0],xi[1],xi[2],xi[3]
	y0,y1,y2,y3 := yi[0],yi[1],yi[2],yi[3]
	h0,h1,h2,_,u1,l2,_ := spline_lu(xi)
	c1,c2 := spline_c1(yi,h0,h1), spline_c2(yi,h1, h2)
	m1,m2 := spline_m1(c1,c2,u1,l2), spline_m2(c1,c2,u1,l2) // m0=m3=0

	for k,v := range xo {
		if v <= x1 {
			yo[k] = spline_z0(m1,h0,x0,x1,y0,y1,v)
		} else if v <= x2 {
			yo[k] = spline_z1(m1,m2,h1,x1,x2,y1,y2,v)
		} else {
			yo[k] = spline_z2(m2,h2,x2,x3,y2,y3,v)
		}
	}

	return
}

func spline_z0(m1,h0,x0,x1,y0,y1,x float64) float64 {
	v0 := 0.0
	v1 := (x-x0)*(x-x0)*(x-x0)*m1/(6*h0)
	v2 := -1.0*y0*(x-x1)/h0
	v3 := (y1 - h0*h0*m1/6)*(x-x0)/h0
	return v0+v1+v2+v3
}

func spline_z1(m1,m2,h1,x1,x2,y1,y2,x float64) float64 {
	v0 := -1.0*(x-x2)*(x-x2)*(x-x2)*m1/(6*h1)
	v1 := (x-x1)*(x-x1)*(x-x1)*m2/(6*h1)
	v2 := -1.0*(y1-h1*h1*m1/6)*(x-x2)/h1
	v3 := (y2-h1*h1*m2/6)*(x-x1)/h1
	return v0+v1+v2+v3
}

func spline_z2(m2,h2,x2,x3,y2,y3,x float64) float64 {
	v0 := -1.0*(x-x3)*(x-x3)*(x-x3)*m2/(6*h2)
	v1 := 0.0
	v2 := -1.0*(y2-h2*h2*m2/6)*(x-x3)/h2
	v3 := y3*(x-x2)/h2
	return v0+v1+v2+v3
}

func spline_m1(c1,c2,u1,l2 float64) float64 {
	return (c1/u1 - c2/2) / (2/u1 - l2/2)
}

func spline_m2(c1,c2,u1,l2 float64) float64 {
	return (c1/2 - c2/l2) / (u1/2 - 2/l2)
}

func spline_c1(yi []float64, h0,h1 float64) float64 {
	y0,y1,y2,_ := yi[0],yi[1],yi[2],yi[3]
	return 6.0 / (h0 + h1) * ((y2 - y1)/h1 - (y1 - y0)/h0)
}

func spline_c2(yi []float64, h1,h2 float64) float64 {
	_,y1,y2,y3 := yi[0],yi[1],yi[2],yi[3]
	return 6.0 / (h1 + h2) * ((y3-y2)/h2 - (y2-y1)/h1)
}

func spline_lu(xi []float64) (h0,h1,h2,l1,u1,l2,u2 float64) {
	x0,x1,x2,x3 := xi[0],xi[1],xi[2],xi[3]

	h0,h1,h2 = x1-x0,x2-x1,x3-x2

	l1 = h0 / (h1 + h0) // lambada1
	u1 = h1 / (h1 + h0)

	l2 = h1 / (h2 + h1) // lambada2
	u2 = h2 / (h2 + h1)

	return

}

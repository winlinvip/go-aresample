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
	Resample(pcm []byte) (npcm []byte, err error)
}

// sample rate resampler.
type srResampler struct {
	channels        int    // Channels, L or LR
	sampleRate      int    // Transform from this sample rate.
	nSampleRate     int    // Transform to this sample rate.

	nbSamples       uint64 // Total samples we got from input.
	nbOutputSamples uint64 // Total samples we put to output.
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

	// Convert pcm to int16 values
	ipcm := []int16{}

	for i:=0; i<len(pcm); i+=2 {
		// 16bits le sample
		v := (int16(pcm[i])) | (int16(pcm[i+1]) << 8)

		ipcm = append(ipcm, v)
	}

	// Resample the ipcm
	opcm := []int16{}
	for i:=0; i<len(ipcm); i+=v.channels {
		// How many samples should we output
		v.nbSamples++
		outputSamples := uint64(float64(v.nbSamples) * float64(v.nSampleRate) / float64(v.sampleRate))

		// Drop samples
		if outputSamples <= v.nbOutputSamples {
			continue
		}

		// Normally insert current sample.
		if outputSamples == v.nbOutputSamples + 1 {
			for j:=0; j<v.channels; j++ {
				opcm = append(opcm, ipcm[i + j])
			}
			v.nbOutputSamples++
			continue
		}

		// Duplicate current sample N times.
		// TODO: FIXME: Insert reasonable value for the duplicated value introduce lots of nosie.
		for k:=v.nbOutputSamples; k < outputSamples; k++ {
			for j:=0; j<v.channels; j++ {
				opcm = append(opcm, ipcm[i+j])
			}
			v.nbOutputSamples++
		}
	}

	// Convert int16 samples to bytes.
	npcm = []byte{}
	for _,v := range opcm {
		npcm = append(npcm, byte(v))
		npcm = append(npcm, byte(v >> 8))
	}

	return
}

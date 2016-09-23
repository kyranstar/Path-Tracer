package lib

import (
	"image"
	"math"
)

type Channel int

const (
	ColorChannel = iota
	VarianceChannel
	StandardDeviationChannel
	SamplesChannel
)

type Pixel struct {
	Samples int
	M, V    RGB
}

func (p *Pixel) AddSample(sample RGB) {
	p.Samples++
	if p.Samples == 1 {
		p.M = sample
		return
	}
	m := p.M
	p.M = p.M.Add(sample.Sub(p.M).DivScalar(float64(p.Samples)))
	p.V = p.V.Add(sample.Sub(m).Multiply(sample.Sub(p.M)))
}

func (p *Pixel) Color() RGB {
	return p.M
}

func (p *Pixel) Variance() RGB {
	if p.Samples < 2 {
		return RGB{}
	}
	return p.V.DivScalar(float64(p.Samples - 1))
}

func (p *Pixel) StandardDeviation() RGB {
	return p.Variance().Sqrt()
}

type Buffer struct {
	W, H   int
	Pixels []Pixel
}

func NewBuffer(w, h int) *Buffer {
	pixels := make([]Pixel, w*h)
	return &Buffer{w, h, pixels}
}

func (b *Buffer) Copy() *Buffer {
	pixels := make([]Pixel, b.W*b.H)
	copy(pixels, b.Pixels)
	return &Buffer{b.W, b.H, pixels}
}

func (b *Buffer) AddSample(x, y int, sample RGB) {
	b.Pixels[y*b.W+x].AddSample(sample)
}

func (b *Buffer) Samples(x, y int) int {
	return b.Pixels[y*b.W+x].Samples
}

func (b *Buffer) Color(x, y int) RGB {
	return b.Pixels[y*b.W+x].Color()
}

func (b *Buffer) Variance(x, y int) RGB {
	return b.Pixels[y*b.W+x].Variance()
}

func (b *Buffer) StandardDeviation(x, y int) RGB {
	return b.Pixels[y*b.W+x].StandardDeviation()
}

func (b *Buffer) Image(channel Channel) image.Image {
	result := image.NewRGBA64(image.Rect(0, 0, b.W, b.H))
	var maxSamples float64
	if channel == SamplesChannel {
		for _, pixel := range b.Pixels {
			maxSamples = math.Max(maxSamples, float64(pixel.Samples))
		}
	}
	for y := 0; y < b.H; y++ {
		for x := 0; x < b.W; x++ {
			var c RGB
			switch channel {
			case ColorChannel:
				// gamma correction
				c = b.Pixels[y*b.W+x].Color().Pow(1 / 2.2)
			case VarianceChannel:
				c = b.Pixels[y*b.W+x].Variance()
			case StandardDeviationChannel:
				c = b.Pixels[y*b.W+x].StandardDeviation()
			case SamplesChannel:
				p := float64(b.Pixels[y*b.W+x].Samples) / maxSamples
				c = RGB{p, p, p}
			}
			result.Set(b.W-1-x, b.H-1-y, c.RGBA())
		}
	}
	return result
}

package lib

import (
	"image/color"
	"math"
)

type RGB struct {
	R, G, B float64
}

func (c RGB) Multiply(c1 RGB) RGB {
	return RGB{c.R * c1.R, c.G * c1.G, c.B * c1.B}
}
func (c RGB) DivScalar(f float64) RGB {
	return RGB{c.R / f, c.G / f, c.B / f}
}

func (c RGB) Add(c1 RGB) RGB {
	return RGB{c.R + c1.R, c.G + c1.G, c.B + c1.B}
}
func (c RGB) Sub(c1 RGB) RGB {
	return RGB{c.R - c1.R, c.G - c1.G, c.B - c1.B}
}
func (c RGB) MultiplyScalar(f float64) RGB {
	return RGB{c.R * f, c.G * f, c.B * f}
}
func (c RGB) Pow(f float64) RGB {
	return RGB{math.Pow(c.R, f), math.Pow(c.G, f), math.Pow(c.B, f)}
}
func (c RGB) Sqrt() RGB {
	return RGB{math.Sqrt(c.R), math.Sqrt(c.G), math.Sqrt(c.B)}
}
func (c RGB) MaxComponent() float64 {
	return math.Max(c.R, math.Max(c.G, c.B))
}
func (c RGB) RGBA() color.RGBA {
	return color.RGBA{uint8(c.R * 255.0), uint8(c.G * 255.0), uint8(c.B * 255.0), uint8(255)}
}

//func (c RGB) clamp() RGB {
//	r := c.R > 1 ? 1 : c.R < 0 ? 0 : c.R
//	g := c.G > 1 ? 1 : c.G < 0 ? 0 : c.G
//	b := c.B > 1 ? 1 : c.B < 0 ? 0 : c.B

//	return RGB{r, g, b}
//}

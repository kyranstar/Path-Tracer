package lib

import (
	"math"
	"math/rand"
)

type Material struct {
	Col          RGB
	Index        float64 // refractive index
	Reflectivity float64
	Transparency float64 // the amount of light to let through
	Gloss        float64 // reflection cone angle in radians
	Emittance    float64
	Tint         float64
}
type BounceType uint8

const (
	BounceTypeAny BounceType = iota
	BounceTypeSpecular
	BounceTypeDiffuse
)

func (m *Material) Color() RGB {
	return m.Col
}

func (m *Material) Bounce(r Ray, fu, fv float64, bounceType BounceType, hit Hit, rnd *rand.Rand) (reflectedRay Ray, reflected bool, weight float64) {

	var direction Vector

	p := m.Reflectivity
	switch bounceType {
	case BounceTypeAny:
		reflected = rnd.Float64() < p
	case BounceTypeSpecular:
		reflected = true
	case BounceTypeDiffuse:
		reflected = false
	}

	if reflected {
		// we should reflect
		reflectDirection := r.Direction.Reflect(hit.Normal)
		direction = Cone(reflectDirection, m.Gloss, fu, fv, rnd)
		weight = p
	} else if m.Transparency > 0 {
		direction = hit.Normal.Refract(r.Direction, m.Index)
		direction = Cone(direction, m.Gloss, fu, fv, rnd)
		hit.Point = hit.Point.Add(direction.MultiplyScalar(1e-4))
		weight = 1 - p
	} else {
		direction = hit.Normal.Add(VectorInUnitSphere(rnd))
		weight = 1 - p
	}

	reflectedRay = Ray{hit.Point, direction.Normalize()}
	return
}

func Lambertian(c RGB) *Material {
	return &Material{Col: c}
}
func Metal(c RGB, gloss, reflectivity, tint float64) *Material {
	return &Material{Col: c, Index: 1, Reflectivity: reflectivity, Gloss: gloss, Tint: tint}
}
func Transparent(c RGB, index, gloss, reflectivity, transparency float64) *Material {
	return &Material{Col: c, Index: index, Gloss: gloss, Reflectivity: reflectivity, Transparency: transparency}
}
func Light(c RGB, emittance float64) *Material {
	return &Material{Col: c, Emittance: emittance, Reflectivity: -1}
}

func Cone(direction Vector, theta, u, v float64, rnd *rand.Rand) Vector {
	if theta < EPS {
		return direction
	}
	theta = theta * (1 - (2 * math.Acos(u) / math.Pi))
	m1 := math.Sin(theta)
	m2 := math.Cos(theta)
	a := v * 2 * math.Pi
	q := VectorInUnitSphere(rnd)
	s := direction.Cross(q)
	t := direction.Cross(s)
	d := Vector{}
	d = d.Add(s.MultiplyScalar(m1 * math.Cos(a)))
	d = d.Add(t.MultiplyScalar(m1 * math.Sin(a)))
	d = d.Add(direction.MultiplyScalar(m2))
	d = d.Normalize()
	return d
}

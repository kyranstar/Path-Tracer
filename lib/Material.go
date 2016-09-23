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
}

func (m *Material) Color() RGB {
	return m.Col
}

func (m *Material) Bounce(r Ray, fu, fv float64, hit Hit, rnd *rand.Rand) (bool, Ray) {

	var direction Vector

	if m.Reflectivity > 0 && rnd.Float64() < m.Reflectivity {
		// we should reflect
		reflectDirection := r.Direction.Reflect(hit.Normal)
		direction = Cone(reflectDirection, m.Gloss, fu, fv, rnd)

	} else if m.Transparency > 0 && rnd.Float64() < m.Transparency {
		direction = hit.Normal.Refract(r.Direction, m.Index)
		direction = Cone(direction, m.Gloss, fu, fv, rnd)
		hit.Point = hit.Point.Add(direction.MultiplyScalar(1e-4))
	} else {
		direction = hit.Normal.Add(VectorInUnitSphere(rnd))
	}

	return true, Ray{hit.Point, direction.Normalize()}
}

func Lambertian(c RGB) *Material {
	return &Material{Col: c}
}
func Metal(c RGB, gloss, reflectivity float64) *Material {
	return &Material{Col: c, Index: 1, Reflectivity: reflectivity, Gloss: gloss}
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

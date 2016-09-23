package lib

import (
	"math"
	"math/rand"
)

type Sphere struct {
	Center Vector
	Radius float64
	Mat    *Material
}

func (s *Sphere) Material() *Material {
	return s.Mat
}
func (s *Sphere) RandomPoint(rnd *rand.Rand, point Vector) Vector {
	dir := s.Center.Min(point).Normalize()
	hem := dir.Add(VectorInUnitSphere(rnd)).Normalize()
	return hem.MultiplyScalar(s.Radius).Add(s.Center)
}

func (s *Sphere) Hit(r Ray, tMin float64, tMax float64) (bool, Hit) {
	centerToRay := r.Origin.Subtract(s.Center)
	a := r.Direction.SquaredLength()
	b := centerToRay.Dot(r.Direction)
	c := centerToRay.SquaredLength() - s.Radius*s.Radius

	discriminant := b*b - a*c

	if discriminant > 0 {
		hit := Hit{Material: s.Material(), Ray: r}
		sqrtDiscrim := math.Sqrt(discriminant)

		temp := (-b - sqrtDiscrim) / a
		if temp < tMax && temp > tMin {
			hit.T = temp
			hit.Point = r.Step(temp)
			hit.Normal = hit.Point.Subtract(s.Center).DivideScalar(s.Radius)
			return true, hit
		}
		temp = (-b + sqrtDiscrim) / a
		if temp < tMax && temp > tMin {
			hit.T = temp
			hit.Point = r.Step(temp)
			hit.Normal = hit.Point.Subtract(s.Center).DivideScalar(s.Radius)
			return true, hit
		}
	}
	return false, Hit{}
}
func (s *Sphere) BoundingBox() Box {
	rad := Vector{s.Radius, s.Radius, s.Radius}
	return Box{s.Center.Subtract(rad), s.Center.Add(rad)}
}

func (s *Sphere) MidPoint() Vector {
	return s.Center
}

package lib

import (
	"math"
	"math/rand"
)

type Vector struct {
	X, Y, Z float64
}

var UnitVector = Vector{1, 1, 1}

func VectorInUnitSphere(rnd *rand.Rand) Vector {
	for {
		r := Vector{rnd.Float64(), rnd.Float64(), rnd.Float64()}
		p := r.MultiplyScalar(2.0).Subtract(UnitVector)
		if p.SquaredLength() >= 1.0 {
			return p
		}
	}
}
func (v Vector) Get(a Axis) float64 {
	switch a {
	case AxisX:
		return v.X
	case AxisY:
		return v.Y
	case AxisZ:
		return v.Z
	}
	return -1
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.SquaredLength())
}

func (v Vector) SquaredLength() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vector) Normalize() Vector {
	return v.DivideScalar(v.Length())
}

func (v Vector) Dot(ov Vector) float64 {
	return v.X*ov.X + v.Y*ov.Y + v.Z*ov.Z
}
func (v Vector) Cross(ov Vector) Vector {
	return Vector{
		v.Y*ov.Z - v.Z*ov.Y,
		v.Z*ov.X - v.X*ov.Z,
		v.X*ov.Y - v.Y*ov.X,
	}
}

func (v Vector) Add(v2 Vector) Vector {
	return Vector{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

func (v Vector) Subtract(ov Vector) Vector {
	return Vector{v.X - ov.X, v.Y - ov.Y, v.Z - ov.Z}
}

func (v Vector) Multiply(ov Vector) Vector {
	return Vector{v.X * ov.X, v.Y * ov.Y, v.Z * ov.Z}
}

func (v Vector) Divide(ov Vector) Vector {
	return Vector{v.X / ov.X, v.Y / ov.Y, v.Z / ov.Z}
}

func (v Vector) AddScalar(t float64) Vector {
	return Vector{v.X + t, v.Y + t, v.Z + t}
}

func (v Vector) SubtractScalar(t float64) Vector {
	return Vector{v.X - t, v.Y - t, v.Z - t}
}

func (v Vector) MultiplyScalar(t float64) Vector {
	return Vector{v.X * t, v.Y * t, v.Z * t}
}

func (v Vector) DivideScalar(t float64) Vector {
	return Vector{v.X / t, v.Y / t, v.Z / t}
}

func (v Vector) Min(v2 Vector) Vector {
	return Vector{math.Min(v.X, v2.X), math.Min(v.Y, v2.Y), math.Min(v.Z, v2.Z)}
}
func (v Vector) Max(v2 Vector) Vector {
	return Vector{math.Max(v.X, v2.X), math.Max(v.Y, v2.Y), math.Max(v.Z, v2.Z)}
}

func (v Vector) Reflect(ov Vector) Vector {
	b := 2 * v.Dot(ov)
	return v.Subtract(ov.MultiplyScalar(b))
}

func (n Vector) Refract(i Vector, ior float64) Vector {
	cosI := n.Dot(i)
	var n1, n2 float64
	if cosI > 0 {
		n1 = ior
		n2 = 1.0
		n = n.MultiplyScalar(-1)
	} else {
		n1 = 1.0
		n2 = ior
		cosI = -cosI
	}

	nr := n1 / n2
	sinT2 := nr * nr * (1 - cosI*cosI)
	//	if sinT2 > 1 {
	//		return Vector{}
	//	}
	cosT := math.Sqrt(1 - sinT2)
	if 1-sinT2 < 0 {
		// total internal reflection
		return i.Reflect(n)
	}
	return i.MultiplyScalar(nr).Add(n.MultiplyScalar(nr*cosI - cosT))
}

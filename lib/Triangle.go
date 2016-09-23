package lib

import (
	"math"
	"math/rand"
)

type Triangle struct {
	Mat        *Material
	V1, V2, V3 Vector
	T1, T2, T3 Vector
	N1, N2, N3 Vector
	Area       float64
}

func NewTriangle(v1, v2, v3, n1, n2, n3 Vector, mat *Material) *Triangle {
	t := Triangle{}
	t.V1 = v1
	t.V2 = v2
	t.V3 = v3
	t.N1 = n1
	t.N2 = n2
	t.N3 = n3
	t.Mat = mat
	t.FixNormals()
	t.Area = .5 * (v3.Subtract(v1)).Cross(v3.Subtract(v2)).Length()
	return &t
}

func (tri *Triangle) RandomPoint(rnd *rand.Rand, point Vector) Vector {
	sum := math.Sqrt(rnd.Float64()) // takes varying line length s+t = sum into account
	t := rnd.Float64() * sum
	s := sum - t
	r := 1.0 - s - t

	return tri.V1.MultiplyScalar(r).Add(tri.V2.MultiplyScalar(s)).Add(tri.V3.MultiplyScalar(t))
}

func (t *Triangle) ScaleAndTranslate(scale float64, translate Vector) {
	t.V1 = t.V1.MultiplyScalar(scale).Add(translate)
	t.V2 = t.V2.MultiplyScalar(scale).Add(translate)
	t.V3 = t.V3.MultiplyScalar(scale).Add(translate)
}

func (t *Triangle) Material() *Material {
	return t.Mat
}

func (t *Triangle) Hit(r Ray, tMin, tMax float64) (bool, Hit) {
	e1x := t.V2.X - t.V1.X
	e1y := t.V2.Y - t.V1.Y
	e1z := t.V2.Z - t.V1.Z
	e2x := t.V3.X - t.V1.X
	e2y := t.V3.Y - t.V1.Y
	e2z := t.V3.Z - t.V1.Z
	px := r.Direction.Y*e2z - r.Direction.Z*e2y
	py := r.Direction.Z*e2x - r.Direction.X*e2z
	pz := r.Direction.X*e2y - r.Direction.Y*e2x
	det := e1x*px + e1y*py + e1z*pz
	if det > -EPS && det < EPS {
		return false, Hit{}
	}
	inv := 1 / det
	tx := r.Origin.X - t.V1.X
	ty := r.Origin.Y - t.V1.Y
	tz := r.Origin.Z - t.V1.Z
	u := (tx*px + ty*py + tz*pz) * inv
	if u < 0 || u > 1 {
		return false, Hit{}
	}
	qx := ty*e1z - tz*e1y
	qy := tz*e1x - tx*e1z
	qz := tx*e1y - ty*e1x
	v := (r.Direction.X*qx + r.Direction.Y*qy + r.Direction.Z*qz) * inv
	if v < 0 || u+v > 1 {
		return false, Hit{}
	}
	d := (e2x*qx + e2y*qy + e2z*qz) * inv
	if d < tMin {
		return false, Hit{}
	}
	return true, Hit{T: d, Normal: t.Normal(), Material: t.Material(), Point: r.Step(d), Ray: r}
}
func (t *Triangle) Normal() Vector {
	return (t.N1.Add(t.N2).Add(t.N3)).DivideScalar(3)

}
func (t *Triangle) BoundingBox() Box {
	min := t.V1.Min(t.V2).Min(t.V3)
	max := t.V1.Max(t.V2).Max(t.V3)
	return Box{min, max}
}
func (t *Triangle) MidPoint() Vector {
	return Vector{(t.V1.X + t.V2.X + t.V3.X) / 3.0, (t.V1.Y + t.V2.Y + t.V3.Y) / 3.0, (t.V1.Z + t.V2.Z + t.V3.Z) / 3.0}
}

func (t *Triangle) FixNormals() {
	e1 := t.V2.Subtract(t.V1)
	e2 := t.V3.Subtract(t.V1)
	n := e1.Cross(e2).Normalize()
	zero := Vector{}
	if t.N1 == zero {
		t.N1 = n
	}
	if t.N2 == zero {
		t.N2 = n
	}
	if t.N3 == zero {
		t.N3 = n
	}
}

package lib

import (
	"math"
	"math/rand"
)

type Hit struct {
	T             float64
	Point, Normal Vector
	Ray           Ray
	*Material
}

type Hittable interface {
	Hit(r Ray, tMin float64, tMax float64) (bool, Hit)
	BoundingBox() Box
	MidPoint() Vector
	Material() *Material
	RandomPoint(rnd *rand.Rand, point Vector) Vector
}

type Box struct {
	Min, Max Vector
}

/*func (b *Box) Hit(r Ray, tMin float64, tMax float64) (bool, Hit) {

}*/
func (b *Box) MidPoint() Vector {
	return Vector{(b.Min.X + b.Max.X) / 2.0, (b.Min.Y + b.Max.Y) / 2.0, (b.Min.Z + b.Max.Z) / 2.0}
}

func (b *Box) BoundingBox() Box {
	return *b
}
func (a *Box) Extend(b Box) {
	a.Min = a.Min.Min(b.Min)
	a.Max = a.Max.Max(b.Max)
}
func (a *Box) SurfaceArea() float64 {
	length := a.Max.X - a.Min.X
	width := a.Max.Z - a.Min.Z
	height := a.Max.Y - a.Min.Y

	return 2*length*width + 2*length*height + 2*width*height
}
func (b *Box) Intersect(r Ray) (float64, float64) {
	x1 := (b.Min.X - r.Origin.X) / r.Direction.X
	y1 := (b.Min.Y - r.Origin.Y) / r.Direction.Y
	z1 := (b.Min.Z - r.Origin.Z) / r.Direction.Z
	x2 := (b.Max.X - r.Origin.X) / r.Direction.X
	y2 := (b.Max.Y - r.Origin.Y) / r.Direction.Y
	z2 := (b.Max.Z - r.Origin.Z) / r.Direction.Z
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	if z1 > z2 {
		z1, z2 = z2, z1
	}
	t1 := math.Max(math.Max(x1, y1), z1)
	t2 := math.Min(math.Min(x2, y2), z2)
	return t1, t2
}
func (a *Box) Intersects(r Ray) bool {
	dir := r.Direction

	tmin := (a.Min.X - r.Origin.X) / dir.X
	tmax := (a.Max.X - r.Origin.X) / dir.X

	if tmin > tmax {
		tmin, tmax = tmax, tmin
	}

	tymin := (a.Min.Y - r.Origin.Y) / dir.Y
	tymax := (a.Max.Y - r.Origin.Y) / dir.Y

	if tymin > tymax {
		tymin, tymax = tymax, tymin
	}

	if (tmin > tymax) || (tymin > tmax) {
		return false
	}

	if tymin > tmin {
		tmin = tymin
	}

	if tymax < tmax {
		tmax = tymax
	}

	tzmin := (a.Min.Z - r.Origin.Z) / dir.Z
	tzmax := (a.Max.Z - r.Origin.Z) / dir.Z

	if tzmin > tzmax {
		tzmin, tzmax = tzmax, tzmin
	}

	if (tmin > tzmax) || (tzmin > tmax) {
		return false
	}

	return true
}

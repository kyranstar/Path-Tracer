package lib

import (
	"math/rand"
)

type Func func(x, y float64) float64

type Volume struct {
	Function *Func
	Box      *Box
	Mat      *Material
}

func NewVolume(function *Func, box *Box, mat *Material) *Volume {
	return &Volume{function, box, mat}
}
func (v *Volume) Hit(r Ray, tMin float64, tMax float64) (bool, Hit) {
	return true, Hit{}
}
func (v *Volume) RandomPoint(rnd *rand.Rand) Vector {
	x := v.Box.Min.X + rnd.Float64()*(v.Box.Max.X-v.Box.Min.X)
	y := v.Box.Min.Y + rnd.Float64()*(v.Box.Max.Y-v.Box.Min.Y)
	z := v.Box.Min.Z + rnd.Float64()*(v.Box.Max.Z-v.Box.Min.Z)
	return Vector{x, y, z}
}
func (v *Volume) BoundingBox() Box {
	return *v.Box
}
func (v *Volume) MidPoint() Vector {
	return v.Box.MidPoint()
}
func (v *Volume) Material() *Material {
	return v.Mat
}

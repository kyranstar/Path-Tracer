package lib

type Ray struct {
	Origin, Direction Vector
}

func (r Ray) Step(t float64) Vector {
	return r.Origin.Add(r.Direction.MultiplyScalar(t))
}

package lib

import (
	"math/rand"
)

type Scene struct {
	objects []Hittable
	Lights  []Hittable
	KDTree  *KDNode
}

func (s *Scene) Add(h Hittable) {
	s.objects = append(s.objects, h)
	if h.Material().Emittance > 0 {
		s.Lights = append(s.Lights, h)
	}
	s.KDTree = build(s.objects, 0)

}
func (s *Scene) AddAll(hittables []Hittable) {
	for _, h := range hittables {
		s.objects = append(s.objects, h)
		if h.Material().Emittance > 0 {
			s.Lights = append(s.Lights, h)
		}
	}
	s.KDTree = build(s.objects, 0)
}
func (s *Scene) RayToRandomLight(p Vector, rnd *rand.Rand) Vector {
	light := s.Lights[rnd.Intn(len(s.Lights))]

	return light.RandomPoint(rnd, p).Subtract(p).Normalize()
}

func (w *Scene) Count() int {
	return len(w.objects)
}

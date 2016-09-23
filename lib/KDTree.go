package lib

import (
	"math"
)

const SAHRes = 32 // the number of possible planes to check for SAH

type KDNode struct {
	BoundingBox Box
	Axis        Axis
	Left, Right *KDNode
	objects     []Hittable
}

func MakeKDTree(objects []Hittable) *KDNode {
	return build(objects, 0)
}

func build(objects []Hittable, depth int) *KDNode {
	if len(objects) == 0 {
		return &KDNode{}
	}

	parent := KDNode{BoundingBox: objects[0].BoundingBox(), objects: objects}

	if len(objects) == 1 {
		parent.Left = &KDNode{objects: make([]Hittable, 0, 0)}
		parent.Right = &KDNode{objects: make([]Hittable, 0, 0)}
		return &parent
	}

	// calculate boundingbox
	for _, v := range objects {
		parent.BoundingBox.Extend(v.BoundingBox())
	}
	// sort objects into left and right by longest axis
	leftObjects := make([]Hittable, 0, len(objects)/2)
	rightObjects := make([]Hittable, 0, len(objects)/2)

	p, axis := parent.FindOptimalSplit()

	for _, object := range objects {
		if object.MidPoint().Get(axis) < p {
			leftObjects = append(leftObjects, object)
		} else {
			rightObjects = append(rightObjects, object)
		}
	}

	// if 50% of objects match, don't subdivide
	var matches int
	for _, a := range leftObjects {
		for _, b := range rightObjects {
			if a == b {
				matches++
			}
		}
	}
	if float64(matches)/float64(len(leftObjects)) < 0.5 && float64(matches)/float64(len(rightObjects)) < 0.5 {
		parent.Left = build(leftObjects, depth+1)
		parent.Right = build(rightObjects, depth+1)
	} else {
		parent.Left = &KDNode{objects: make([]Hittable, 0, 0)}
		parent.Right = &KDNode{objects: make([]Hittable, 0, 0)}
	}
	parent.Left.Axis = axis
	parent.Right.Axis = axis

	return &parent
}
func (n *KDNode) FindOptimalSplit() (float64, Axis) {
	var bestAxis Axis
	var bestP float64
	bestSAH := math.MaxFloat64

	for i := 0; i < SAHRes; i++ {
		p := n.BoundingBox.Min.X + (n.BoundingBox.Max.X-n.BoundingBox.Min.X)*float64(i)/float64(SAHRes)
		val := SAHValue(n.objects, p, AxisX)
		if val < bestSAH {
			bestAxis = AxisX
			bestSAH = val
			bestP = p
		}
	}
	for i := 0; i < SAHRes; i++ {
		p := n.BoundingBox.Min.Y + (n.BoundingBox.Max.Y-n.BoundingBox.Min.Y)*float64(i)/float64(SAHRes)
		val := SAHValue(n.objects, p, AxisY)
		if val < bestSAH {
			bestAxis = AxisY
			bestSAH = val
			bestP = p
		}
	}
	for i := 0; i < SAHRes; i++ {
		p := n.BoundingBox.Min.Z + (n.BoundingBox.Max.Z-n.BoundingBox.Min.Z)*float64(i)/float64(SAHRes)
		val := SAHValue(n.objects, p, AxisZ)
		if val < bestSAH {
			bestAxis = AxisZ
			bestSAH = val
			bestP = p
		}
	}

	return bestP, bestAxis
}
func SAHValue(objects []Hittable, p float64, a Axis) float64 {
	var nL, nR float64
	var bL, bR Box
	for _, v := range objects {
		if v.MidPoint().Get(a) < p {
			bL.Extend(v.BoundingBox())
			nL++
		} else {
			bR.Extend(v.BoundingBox())
			nR++
		}
	}
	return nL*bL.SurfaceArea() + nR*bR.SurfaceArea()
}
func (node *KDNode) Hit(r Ray, tMin, tMax float64, intersections *int) (bool, Hit) {
	return node.FindHit(r, tMin, tMax, intersections, true)
}
func (node *KDNode) Intersects(r Ray, tMin, tMax float64, intersections *int) bool {
	b, _ := node.FindHit(r, tMin, tMax, intersections, false)
	return b
}

func (node *KDNode) FindHit(r Ray, tMin, tMax float64, intersections *int, lookForClosest bool) (bool, Hit) {
	if !node.BoundingBox.Intersects(r) {
		return false, Hit{}
	}
	if len(node.Left.objects) > 0 || len(node.Right.objects) > 0 {
		bL, hL := node.Left.Hit(r, tMin, tMax, intersections)
		bR, hR := node.Right.Hit(r, tMin, tMax, intersections)
		if bL && bR {
			if hL.T < hR.T {
				return true, hL
			} else {
				return true, hR
			}
		}
		if bL {
			return true, hL
		}
		if bR {
			return true, hR
		}
		return false, Hit{}
	} else {
		// we have reached a leaf
		return node.IntersectShapes(r, tMin, tMax, intersections, lookForClosest)
	}
}

func (node *KDNode) IntersectShapes(r Ray, tMin, tMax float64, intersections *int, lookForClosest bool) (bool, Hit) {
	hit := Hit{}
	intersected := false
	for _, shape := range node.objects {
		b, h := shape.Hit(r, tMin, tMax)

		(*intersections)++
		if b && (!intersected || h.T < hit.T) {
			if !lookForClosest {
				return true, h
			}
			hit = h
			intersected = true
		}
	}
	return intersected, hit
}

package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	. "./lib"
)
import "runtime/pprof"

var CamPosition = Vector{2.5, 2, -4}
var CamDirection = CamPosition.Add(Vector{0, -.4, 1})

const Width = 640
const Height = 580
const Fov = 50.0
const ApertureDiameter = 0.000001
const SPP = 20 // samples per pixel
const OutputFile = "img.png"
const tMin = .001
const tMax = math.MaxFloat64
const MaxDepth = 5
const ShadowRays = 25

const AdaptiveSamples = 0
const AdaptiveThreshold = 1
const AdaptiveExponent = 3

var NumCPU = runtime.NumCPU()

var random = rand.New(rand.NewSource(0))

var flagCpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *flagCpuprofile != "" {
		f, err := os.Create(*flagCpuprofile)
		if err != nil {
			fmt.Println(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	buf := NewBuffer(Width, Height)
	scene := &Scene{}
	//mat := Metal(RGB{0.9, 1.0, 0.9}, math.Pi/8)
	//mat2 := Lambertian(RGB{0.1, 1.0, 1.0})
	//mat3 := Lambertian(RGB{0.1, 0.1, 1.0})
	//scene.Add(&Sphere{Radius: 1, Center: Vector{7, 0, 0}, Material: mat})
	//scene.Add(&Sphere{Radius: 2, Center: Vector{7, 3, 1}, Material: mat2})
	//scene.Add(&Sphere{Radius: 2, Center: Vector{9, 2, -1}, Material: mat3})

	var objects []Hittable
	//	for i := 0; i < 5; i++ {
	//		for j := 0; j < 5; j++ {
	//			var mat *Material
	//			switch rand.Intn(5) {
	//			case 0, 1:
	//				mat = Transparent(RGB{random.Float64(), random.Float64(), random.Float64()}, random.Float64()*.5+1, 0, .3, .7)
	//			case 2:
	//				mat = Metal(RGB{random.Float64(), random.Float64(), random.Float64()}, random.Float64(), 0.85)
	//			case 3, 4:
	//				mat = Light(RGB{random.Float64(), 1, random.Float64()}, 1)
	//			}
	//			objects = append(objects, &Sphere{Radius: random.Float64() * .5, Center: Vector{float64(i), 0, float64(j)}, Mat: mat})
	//		}
	//	}
	bl := Vector{0, 0, 0}
	br := Vector{5, 0, 0}
	tl := Vector{0, 0, 5}
	tr := Vector{5, 0, 5}
	objects = append(objects, NewTriangle(bl, br, tl, Vector{0, 1, 0}, Vector{0, 1, 0}, Vector{0, 1, 0}, Lambertian(RGB{.5, .5, .5})))
	objects = append(objects, NewTriangle(tl, tr, br, Vector{0, 1, 0}, Vector{0, 1, 0}, Vector{0, 1, 0}, Lambertian(RGB{.5, .5, .5})))
	objects = append(objects, &Sphere{Center: Vector{1, 2, 4}, Radius: .1, Mat: Light(RGB{1, 1, 1}, 1)})
	objects = append(objects, &Sphere{Center: Vector{2, 2, 4}, Radius: .1, Mat: Light(RGB{1, 1, 1}, 1)})
	objects = append(objects, &Sphere{Center: Vector{4, 2, 4}, Radius: .1, Mat: Light(RGB{1, 1, 1}, 1)})
	objects = append(objects, &Sphere{Center: Vector{4, 2, 2}, Radius: .1, Mat: Light(RGB{1, 1, 1}, 1)})
	objects = append(objects, &Sphere{Center: Vector{1.25, .5, 3}, Radius: .5, Mat: Lambertian(RGB{.8, .1, .1})})
	//barrel, _ := LoadOBJ("barrel.obj", Vector{1.5, 1, 1.5}, .5, *Light(RGB{.8, .6, .2}, .75))
	teapot, _ := LoadOBJ("teapot.obj", Vector{2.4, .8, -1}, .25, *Transparent(RGB{.9, 1, .9}, 1.5, 0, .3, .7))

	//objects = append(objects, barrel)
	objects = append(objects, teapot)

	scene.AddAll(objects)

	cam := NewCamera(CamPosition, CamDirection, Fov, Width/Height, ApertureDiameter)
	fmt.Println("Rendering", Width, "x", Height, "image with", SPP, "samples per pixel.")
	start := time.Now()

	render(scene, cam, buf)

	fmt.Println("Render time: ", (time.Since(start)))
	WritePng(OutputFile, buf.Image(ColorChannel))

}

func render(scene *Scene, cam *Camera, buf *Buffer) {
	intersections := 0
	runtime.GOMAXPROCS(NumCPU)
	ch := make(chan int, Height)
	for i := 0; i < NumCPU; i++ {
		go func(i int) {
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			for y := i; y < Height; y += NumCPU {
				for x := 0; x < Width; x++ {
					for i := 0; i < SPP; i++ {
						u := (float64(x) + rnd.Float64()) / float64(Width)
						v := (float64(y) + rnd.Float64()) / float64(Height)

						r := cam.RayAt(u, v, rnd)

						c := getColor(r, scene, 0, rnd, &intersections)
						buf.AddSample(x, y, c)
					}
					if AdaptiveSamples > 0 {
						v := buf.StandardDeviation(x, y).MaxComponent()
						v = v / AdaptiveThreshold
						if v > 1 {
							v = 1
						} else if v < 0 {
							v = 0
						}
						v = math.Pow(v, AdaptiveExponent)
						samples := int(v * float64(AdaptiveSamples))
						for i := 0; i < samples; i++ {
							u := (float64(x) + rnd.Float64()) / float64(Width)
							v := (float64(y) + rnd.Float64()) / float64(Height)
							r := cam.RayAt(u, v, rnd)

							c := getColor(r, scene, 0, rnd, &intersections)
							buf.AddSample(x, y, c)
						}
					}
				}
				ch <- y
			}
		}(i)
	}
	for j := 0; j < Height; j++ {
		row := <-ch
		fmt.Println("Finished row", row, "out of", Height, ",", (float64(j) / Height * 100), "% done")
	}
	fmt.Println("Intersections: ", intersections)

}

func getColor(r Ray, scene *Scene, depth int, rnd *rand.Rand, intersections *int) RGB {
	if depth > MaxDepth {
		return background(r)
	}
	b, hit := scene.KDTree.Hit(r, tMin, tMax, intersections)

	if b {
		if hit.Material.Emittance > 0.0 {
			return hit.Material.Color()
		}
		mode := BounceTypeAny
		bouncedRay, reflected, p := hit.Bounce(r, rnd.Float64(), rnd.Float64(), mode, hit, rnd)
		if mode == BounceTypeAny {
			p = 1
		}
		if p > 0 && reflected {
			// specular
			indirectLight := getColor(bouncedRay, scene, depth+1, rnd, intersections)
			tinted := indirectLight.Mix(hit.Material.Color().Multiply(indirectLight), hit.Material.Tint)
			return tinted.MultiplyScalar(p)
		} else if p > 0 && !reflected {
			//diffuse
			indirectLight := getColor(bouncedRay, scene, depth+1, rnd, intersections)
			directLight := getLighting(scene, hit, bouncedRay, rnd)
			return hit.Material.Color().Multiply(directLight.Add(indirectLight)).MultiplyScalar(p)
		}
		return RGB{}
	}
	return background(r)
}
func getLighting(scene *Scene, hit Hit, bounce Ray, rnd *rand.Rand) RGB {
	var intersections int
	var contrib RGB
	for _, light := range scene.Lights {
		for i := 0; i < ShadowRays; i++ {
			L_i := light.Material().Color().MultiplyScalar(light.Material().Emittance)
			pToL := light.RandomPoint(rnd, hit.Point).Subtract(hit.Point)
			occluded := scene.KDTree.Intersects(Ray{Origin: hit.Point, Direction: pToL}, tMin, tMax, &intersections)
			if !occluded {
				dot := math.Max(0.0, hit.Normal.Dot(pToL.MultiplyScalar(-1)))
				contrib = contrib.Add(L_i.MultiplyScalar(dot))
			}
		}
	}
	return contrib
}

func background(r Ray) RGB {
	//return RGB{0, 0, 0}
	return RGB{0, .3, .5}.MultiplyScalar(math.Max(0.0, r.Direction.Dot(Vector{0, 1, 1})))
}

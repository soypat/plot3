package plot3

import (
	"image/color"
	"math"

	"github.com/fogleman/ln/ln"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Curve defines a 3D Curve. Points are ordered
type Curve struct {
	XYZs
	color.Color
	minXYZ struct{ X, Y, Z float64 }
	maxXYZ struct{ X, Y, Z float64 }
}

// NewCurve Creates new curve ready to graph from XYZer
func NewCurve(xyz XYZer) Curve {
	c := Curve{}
	c.minXYZ.X, c.minXYZ.Y, c.minXYZ.Z = xyz.XYZ(0)
	c.maxXYZ.X, c.maxXYZ.Y, c.maxXYZ.Z = xyz.XYZ(0)
	c.XYZs = CopyXYZs(xyz)
	for i := range c.XYZs {
		x, y, z := c.XYZ(i)
		updateBounds(&c, x, y, z)
	}
	return c
}

// Plot implements gonum's Plotter interface
func (c Curve) Plot(canvas draw.Canvas, plt *plot.Plot) {
	const FOVAngle, nearestPoint = 90., 0.1
	plt.X.Min, plt.Y.Min, _ = c.Min()
	plt.X.Max, plt.Y.Max, _ = c.Max()

	trX, trY := plt.Transforms(&canvas)

	scene := ln.Scene{}
	scene.Add(c)
	// calculate viewing limits
	xl, yl, zl := c.Min()
	xg, yg, zg := c.Max()
	maxAbs := math.Max(xg, math.Max(yg, zg))
	distance := math.Sqrt(math.Pow(xg-xl, 2) + math.Pow(yg-yl, 2) + math.Pow(zg-zl, 2))

	// ISO view
	eye := ln.Vector{X: maxAbs, Y: maxAbs, Z: maxAbs}
	center := ln.Vector{X: 0, Y: 0, Z: 0}
	up := ln.Vector{X: 0, Y: 0, Z: 1}
	size := canvas.Size()

	trfmPaths := scene.Render(eye, center, up,
		size.X.Points(), size.Y.Points(), FOVAngle, nearestPoint, distance, 0.01)

	canvas.SetColor(c.Color)
	// plt lims hack

	for _, path := range trfmPaths {
		var p vg.Path
		for i := range path[:len(path)-2] {
			x1, x2 := trX(path[i].X), trX(path[i+1].X)
			y1, y2 := trY(path[i].Y), trY(path[i+1].Y)
			p.Move(vg.Point{X: x1, Y: y1})
			p.Line(vg.Point{X: x2, Y: y2})
		}
		p.Close()
		canvas.Fill(p)
	}
}

// Plot does a 3D lineplot using fogleman/ln library
func Plot(filename string, xyz XYZer) {
	scene := ln.Scene{}
	c := NewCurve(xyz)
	// calculate viewing limits
	xl, yl, zl := c.Min()
	xg, yg, zg := c.Max()
	maxAbs := math.Max(xg, math.Max(yg, zg))
	distance := math.Sqrt(math.Pow(xg-xl, 2) + math.Pow(yg-yl, 2) + math.Pow(zg-zl, 2))

	scene.Add(c)
	width := 750.0
	height := 750.0
	// ISO view
	eye := ln.Vector{X: maxAbs, Y: maxAbs, Z: maxAbs}
	center := ln.Vector{X: 0, Y: 0, Z: 0}
	up := ln.Vector{X: 0, Y: 0, Z: 1}
	// add axis
	paths := scene.Render(eye, center, up, width, height, 90, 0.1, distance, 0.01)
	paths.WriteToPNG(filename, width, height)
}

// Min returns minimum x, y and z values of a Curve.
//
// The numbers returned do NOT define a point of the Curve.
func (c Curve) Min() (float64, float64, float64) {
	return c.minXYZ.X, c.minXYZ.Y, c.minXYZ.Z
}

// Max Like Min but maximum values
func (c Curve) Max() (float64, float64, float64) {
	return c.maxXYZ.X, c.maxXYZ.Y, c.maxXYZ.Z
}

// BoundingBox ln.Shape interface implementation
func (c Curve) BoundingBox() ln.Box {
	return ln.Box{Min: ln.Vector(c.maxXYZ), Max: ln.Vector(c.maxXYZ)}
}

// Compile ln.Shape interface implementation
func (c Curve) Compile() {}

// Contains ln.Shape interface implementation
func (c Curve) Contains(ln.Vector, float64) bool { return false }

// Intersect ln.Shape interface implementation
func (c Curve) Intersect(ln.Ray) ln.Hit { return ln.NoHit }

// Paths ln.Shape interface implementation
func (c Curve) Paths() ln.Paths { return ln.Paths{ln.Path(c.XYZs)} }

// updateBounds min/max values of curve to then save compute time on
// visualization calculations
func updateBounds(c *Curve, x, y, z float64) {
	c.minXYZ.X, c.maxXYZ.X = math.Min(c.minXYZ.X, x), math.Max(c.maxXYZ.X, x)
	c.minXYZ.Y, c.maxXYZ.Y = math.Min(c.minXYZ.Y, y), math.Max(c.maxXYZ.Y, y)
	c.minXYZ.Z, c.maxXYZ.Z = math.Min(c.minXYZ.Z, z), math.Max(c.maxXYZ.Z, z)
}

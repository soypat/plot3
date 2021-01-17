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
	// hide the load of interfaces normal users don't need to use
	points path
	color.Color
}

// path hides away 3D graphing utility interfaces to reduce
// Method bloat in normal use of Curve type
type path struct {
	XYZs
	minXYZ struct{ X, Y, Z float64 }
	maxXYZ struct{ X, Y, Z float64 }
	path   vg.Path
}

// NewCurve Creates new curve ready to graph from XYZer
func NewCurve(xyz XYZer) Curve {
	c := Curve{}
	c.points.minXYZ.X, c.points.minXYZ.Y, c.points.minXYZ.Z = xyz.XYZ(0)
	c.points.maxXYZ.X, c.points.maxXYZ.Y, c.points.maxXYZ.Z = xyz.XYZ(0)
	c.points.XYZs = CopyXYZs(xyz)
	for i := range c.points.XYZs {
		x, y, z := c.points.XYZ(i)
		updateBounds(&c.points, x, y, z)
	}
	return c
}

// Plot implements gonum's Plotter interface
func (c Curve) Plot(canvas draw.Canvas, plt *plot.Plot) {
	const FOVAngle, nearestPoint = 90., 0.1
	// plt.X.Min, plt.Y.Min, _ = c.Min()
	// plt.X.Max, plt.Y.Max, _ = c.Max()

	trX, trY := plt.Transforms(&canvas)
	size := canvas.Size()
	trfmPaths := c.planarPaths(size.X.Points(), size.Y.Points())
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
		canvas.Stroke(p)
	}
}

func (c Curve) planarPaths(width, height float64) ln.Paths {
	const FOVAngle, nearestPoint, stepDetail = 90., 0.1, 0.01
	scene := ln.Scene{}
	scene.Add(c.points)
	// calculate viewing limits
	xl, yl, zl := c.points.Min()
	xg, yg, zg := c.points.Max()
	maxAbs := math.Max(xg, math.Max(yg, zg))
	distance := math.Sqrt(math.Pow(xg-xl, 2) + math.Pow(yg-yl, 2) + math.Pow(zg-zl, 2))

	// ISO view
	eye, center, up := isoview(maxAbs)
	return scene.Render(eye, center, up,
		width, height, FOVAngle, nearestPoint, distance, stepDetail)
}

// DataRange implements DataRanger interface
func (c Curve) DataRange() (xmin, xmax, ymin, ymax float64) {
	paths := c.planarPaths(c.points.CharLength(), c.points.CharLength())
	box := paths.BoundingBox()
	xmax, ymax = box.Max.X, box.Max.Y
	xmin, ymin = box.Min.X, box.Min.Y
	return xmin, xmax, ymin, ymax
}

// Plot does a 3D lineplot using fogleman/ln library
func Plot(filename string, xyz XYZer) {
	scene := ln.Scene{}
	c := NewCurve(xyz)
	// calculate viewing limits
	xl, yl, zl := c.points.Min()
	xg, yg, zg := c.points.Max()
	maxAbs := math.Max(xg, math.Max(yg, zg))
	distance := math.Sqrt(math.Pow(xg-xl, 2) + math.Pow(yg-yl, 2) + math.Pow(zg-zl, 2))

	scene.Add(c.points)
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
func (p path) Min() (float64, float64, float64) {
	return p.minXYZ.X, p.minXYZ.Y, p.minXYZ.Z
}

// Max Like Min but maximum values
func (p path) Max() (float64, float64, float64) {
	return p.maxXYZ.X, p.maxXYZ.Y, p.maxXYZ.Z
}

// BoundingBox ln.Shape interface implementation
func (p path) BoundingBox() ln.Box {
	return ln.Box{Min: ln.Vector(p.maxXYZ), Max: ln.Vector(p.maxXYZ)}
}

// Compile ln.Shape interface implementation
func (p path) Compile() {}

// Contains ln.Shape interface implementation
func (p path) Contains(ln.Vector, float64) bool { return false }

// Intersect ln.Shape interface implementation
func (p path) Intersect(ln.Ray) ln.Hit { return ln.NoHit }

// Paths ln.Shape interface implementation
func (p path) Paths() ln.Paths { return ln.Paths{ln.Path(p.XYZs)} }

// updateBounds min/max values of curve to then save compute time on
// visualization calculations
func updateBounds(p *path, x, y, z float64) {
	p.minXYZ.X, p.maxXYZ.X = math.Min(p.minXYZ.X, x), math.Max(p.maxXYZ.X, x)
	p.minXYZ.Y, p.maxXYZ.Y = math.Min(p.minXYZ.Y, y), math.Max(p.maxXYZ.Y, y)
	p.minXYZ.Z, p.maxXYZ.Z = math.Min(p.minXYZ.Z, z), math.Max(p.maxXYZ.Z, z)
}

func (p path) CompileVGPath() vg.Path {
	const FOVAngle, nearestPoint = 90., 0.1

	scene := ln.Scene{}
	scene.Add(p)
	// calculate viewing limits
	xl, yl, zl := p.Min()
	xg, yg, zg := p.Max()
	maxAbs := math.Max(xg, math.Max(yg, zg))

	distance := math.Sqrt(math.Pow(xg-xl, 2) + math.Pow(yg-yl, 2) + math.Pow(zg-zl, 2))

	// ISO view
	eye, center, up := isoview(maxAbs)

	trfmPaths := scene.Render(eye, center, up,
		1000, 1000, FOVAngle, nearestPoint, distance, 0.01)

	var vgp vg.Path
	// There should be only one path
	path := trfmPaths[0]
	for i := range trfmPaths[0][:len(path)-2] {
		x1, x2 := vg.Length(path[i].X), vg.Length(path[i+1].X)
		y1, y2 := vg.Length(path[i].Y), vg.Length(path[i+1].Y)
		vgp.Move(vg.Point{X: x1, Y: y1})
		vgp.Line(vg.Point{X: x2, Y: y2})
	}
	return vgp
}

// isoview takes a characteristic distance of the graph
// and returns a normalized view configuration (ISO view)
func isoview(λ float64) (eye ln.Vector, center ln.Vector, up ln.Vector) {
	eye = ln.Vector{X: λ, Y: λ, Z: λ}
	center = ln.Vector{X: 0, Y: 0, Z: 0}
	up = ln.Vector{X: 0, Y: 0, Z: 1}
	return eye, center, up
}

//
func (p path) CharLength() float64 {
	return math.Sqrt(math.Pow(p.maxXYZ.X-p.minXYZ.X, 2) +
		math.Pow(p.maxXYZ.Y-p.minXYZ.Y, 2) +
		math.Pow(p.maxXYZ.Z-p.minXYZ.Z, 2))
}

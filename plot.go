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

	trX, trY := plt.Transforms(&canvas)
	size := canvas.Size()
	trfmPaths := c.planarPaths(size.X.Points(), size.Y.Points())
	canvas.SetColor(c.Color)

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

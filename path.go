package plot3

import (
	"math"

	"github.com/fogleman/ln/ln"
)

// path hides away 3D graphing utility interfaces to reduce
// Method bloat in normal use of Curve type
type path struct {
	XYZs
	minXYZ struct{ X, Y, Z float64 }
	maxXYZ struct{ X, Y, Z float64 }
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

// isoview takes a characteristic distance of the graph
// and returns a normalized view configuration (ISO view)
func isoview(λ float64) (eye ln.Vector, center ln.Vector, up ln.Vector) {
	eye = ln.Vector{X: λ, Y: λ, Z: λ}
	center = ln.Vector{X: 0, Y: 0, Z: 0}
	up = ln.Vector{X: 0, Y: 0, Z: 1}
	return eye, center, up
}

// CharLength returns the bounding box's characteristic length
//
// It can be thought of as the longest straight line that could fit in the box.
func (p path) CharLength() float64 {
	return math.Sqrt(math.Pow(p.maxXYZ.X-p.minXYZ.X, 2) +
		math.Pow(p.maxXYZ.Y-p.minXYZ.Y, 2) +
		math.Pow(p.maxXYZ.Z-p.minXYZ.Z, 2))
}

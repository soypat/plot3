package plot3

import "github.com/fogleman/ln/ln"

// XYZer wraps the Len and XYZ methods.
type XYZer interface {
	// Len returns the number of x, y, z triples.
	Len() int

	// XYZ returns an x, y, z triple.
	XYZ(int) (float64, float64, float64)
}

// XYZs implements the XYZer interface using a slice.
type XYZs ln.Path

// Len implements the Len method of the XYZer interface.
func (xyz XYZs) Len() int {
	return len(xyz)
}

// XYZ returns a point in space (3 element tuple)
func (xyz XYZs) XYZ(i int) (float64, float64, float64) {
	return xyz[i].X, xyz[i].Y, xyz[i].Z
}

// CopyXYZs copies an XYZer.
func CopyXYZs(data XYZer) XYZs {
	cpy := make(XYZs, data.Len())
	for i := range cpy {
		cpy[i].X, cpy[i].Y, cpy[i].Z = data.XYZ(i)
	}
	return cpy
}

// XYZerFromSlices generates XYZer interface from a set of
// points (x[i] corresponds to y[i] and z[i])
func XYZerFromSlices(x, y, z []float64) XYZer {

	if len(x) != len(y) || len(y) != len(z) {
		panic("length of slices unequal")
	}
	xyz := make(XYZs, len(x))
	for i := range x {
		xyz[i].X, xyz[i].Y, xyz[i].Z = x[i], y[i], z[i]
	}
	return XYZer(xyz)
}

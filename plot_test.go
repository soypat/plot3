package plot3_test

import (
	"math"
	"testing"

	"github.com/soypat/plot3"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

// XYZFunc defines a parametrized 3D curve
type XYZFunc struct {
	X, Y, Z    func(in float64) float64
	start, end float64
	steps      int
}

// Len Returns amount of sampled points.
// part of XYZer interface implementation
func (f XYZFunc) Len() int { return f.steps }

// XYZ returns a point in space (3 element tuple)
func (f XYZFunc) XYZ(i int) (float64, float64, float64) {
	if i >= f.steps {
		panic("out of bounds")
	}
	t := float64(i)*(f.end-f.start)/float64(f.steps) + f.start
	return f.X(t), f.Y(t), f.Z(t)
}

var cos, sin, exp = math.Cos, math.Sin, math.Exp

// for log spiral
const a_l, b_l = 1., 0.25

var logSpiral = XYZFunc{
	start: 0, end: 20., steps: 100,
	X: func(in float64) float64 { return a_l * exp(b_l*in) * cos(in) },
	Y: func(in float64) float64 { return a_l * exp(b_l*in) * sin(in) },
	Z: func(in float64) float64 { return in },
}

func TestPlotter(t *testing.T) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "plot3 using Gonum Plotter Interface (WIP)"
	p.X.Label.Text = "X "
	p.Y.Label.Text = "Y"
	c := plot3.NewCurve(logSpiral)
	// c.Color = color.RGBA{R: 196, B: 128, A: 255}
	p.Add(c)
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "logspiral.png"); err != nil {
		panic(err)
	}
}

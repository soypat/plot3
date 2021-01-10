package plot3_test

import (
	"image/color"
	"math"
	"plot3"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

var cos, sin, exp = math.Cos, math.Sin, math.Exp

func TestPlot(t *testing.T) {
	a, b := 1., 0.25 // log spiral constants
	time := linspace(0.0, 20.0, 100)
	x := newFromFunc(time, func(in float64) float64 { return a * exp(b*in) * cos(in) })
	y := newFromFunc(time, func(in float64) float64 { return a * exp(b*in) * sin(in) })
	z := newFromFunc(time, func(in float64) float64 { return in })
	plot3.Plot("plotspiral.png", x, y, z)
}

func TestPlotter(t *testing.T) {
	a, b := 1., 0.25 // log spiral constants
	time := linspace(0.0, 20.0, 100)
	x := newFromFunc(time, func(in float64) float64 { return a * exp(b*in) * cos(in) })
	y := newFromFunc(time, func(in float64) float64 { return a * exp(b*in) * sin(in) })
	z := newFromFunc(time, func(in float64) float64 { return in })

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "LogSpiral using gonum's Plotter interface"
	p.X.Label.Text = "X (does not work yet)"
	p.Y.Label.Text = "Y"
	c := plot3.NewCurve(plot3.XYZerFromSlices(x, y, z))
	c.Color = color.RGBA{R: 196, B: 128, A: 255}
	p.Add(c)
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "logspiral.png"); err != nil {
		panic(err)
	}
}

func newFromFunc(x []float64, f func(float64) float64) []float64 {
	y := make([]float64, len(x))
	for i := range y {
		y[i] = f(x[i])
	}
	return y
}

func linspace(start, end float64, samples int) []float64 {
	if samples <= 0 {
		panic("linspace: sampling 0 or negative")
	}
	dx := (end - start) / float64(samples)
	sli := make([]float64, samples)
	for i := range sli {
		sli[i] = dx*float64(i) + start
	}
	return sli
}

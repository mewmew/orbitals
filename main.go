// References.
//
//    [1]: https://chemistrygod.com/atomic-orbital

package main

import (
	"fmt"
	"log"
	"math"
	"math/cmplx"

	"github.com/pkg/errors"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// Convert meter to picometer.
const pm = 0.000000000001 // 1.0 * 10^{-12} m

// Bohr radius with unit m.
const a0 = 52.9177210903 * pm // 52.9 pm

func main() {
	// 1s-orbital.
	var lines []Line
	{
		const (
			n = 1 // principal quantum number
			l = 0 // azimuthal quantum number
			m = 0 // magnetic quantum number
		)
		line := getLine(n, l, m)
		lines = append(lines, line)
	}
	// 2s-orbital.
	{
		const (
			n = 2 // principal quantum number
			l = 0 // azimuthal quantum number
			m = 0 // magnetic quantum number
		)
		line := getLine(n, l, m)
		lines = append(lines, line)
	}
	// 3s-orbital.
	{
		const (
			n = 3 // principal quantum number
			l = 0 // azimuthal quantum number
			m = 0 // magnetic quantum number
		)
		line := getLine(n, l, m)
		lines = append(lines, line)
	}
	// 2p-orbitals.
	{
		const (
			n = 2 // principal quantum number
			l = 1 // azimuthal quantum number
			//m = 0 // magnetic quantum number
		)
		for m := -1; m <= l; m++ {
			line := getLine(n, l, m)
			lines = append(lines, line)
		}
	}
	if err := genPlot("radial_probability.png", lines...); err != nil {
		log.Fatalf("%+v", err)
	}
}

// getLine returns a plotter line of the specified (n, l, m)-orbital.
func getLine(n, l, m int) Line {
	vals := getValues(n, l, m)
	legend := getLegend(n, l, m)
	line := Line{
		XYs:    vals,
		Legend: legend,
	}
	return line
}

// getLegend returns a legend for the plotter line of the specified (n, l,
// m)-orbital.
func getLegend(n, l, m int) string {
	return fmt.Sprintf("(n=%d, l=%d, m=%d)-orbital", n, l, m)
}

// getValues returns the radial probability values of the (n, l, m)-orbital
// based on the specified principal quantum number, n, azimuthal quantum number,
// l, and magnetic quantum number, m.
func getValues(n, l, m int) plotter.XYs {
	Psi := Orbitals(n, l, m)
	var xys plotter.XYs
	const (
		//rho = a0   // radial distance (radius)
		theta = 15 // inclination (angular)
		phi   = 30 // azimuth (angular)
	)
	for r := 0.0 * pm; r < 1300*pm; r += 1.0 * pm {
		rho := r
		psi := Psi(rho, theta, phi)
		//psi2 := math.Pow(psi, 2)
		radialProb := RadialProb(rho, psi)
		//fmt.Printf("r=%.2g pm\n", rho/pm)
		//fmt.Printf("   psi^2:       %v\n", psi2)
		//fmt.Printf("   radial prob: %v\n", radialProb)
		//fmt.Println()
		xy := plotter.XY{
			X: rho,
			Y: radialProb,
		}
		xys = append(xys, xy)
	}
	// Normalize Y-values.
	total := 0.0
	for i := range xys {
		total += xys[i].Y
	}
	if total != 0 {
		for i := range xys {
			xys[i].Y /= total
		}
	}
	// Use pm for X-axis.
	for i := range xys {
		xys[i].X /= pm
	}
	return xys
}

// Line is a plotter line.
type Line struct {
	// X- and Y-values.
	XYs plotter.XYs
	// Plotter line legend.
	Legend string
}

// getPlot generates a plot containing the given plotter lines.
func genPlot(dstPath string, elems ...Line) error {
	p, err := plot.New()
	if err != nil {
		return errors.WithStack(err)
	}
	p.Title.Text = "Radial probability"
	p.X.Label.Text = "Radius (pm)"
	p.Y.Label.Text = "radial probability"

	for i, elem := range elems {
		line, err := plotter.NewLine(elem.XYs)
		if err != nil {
			return errors.WithStack(err)
		}
		line.Color = plotutil.Color(i)
		line.LineStyle.Width = vg.Points(2)
		// Add values.
		p.Add(line)
		p.Legend.Add(elem.Legend, line)
		p.Legend.ThumbnailWidth = 5 * vg.Centimeter
		p.Legend.Top = true
	}
	// Store plot as PNG image.
	fmt.Printf("creating %q\n", dstPath)
	if err := p.Save(36*vg.Centimeter, 36*vg.Centimeter, dstPath); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// RadialProb returns the radial probability based on the given radius, r, and
// psi for the s-orbital.
//
// NOTE: the returned probability is not yet normalized.
func RadialProb(r, psi float64) float64 {
	// area of sphere.
	area := 4.0 * math.Pi * math.Pow(r, 2)
	// radial probability = area * psi^2.
	return area * math.Pow(psi, 2)
}

// psi1SOrbital returns the psi function of the 1s-orbital (n=1, l=0, m=0).
//
// ref: https://chemistrygod.com/atomic-orbital#s-orbital
func psi1SOrbital(rho, theta, phi float64) float64 {
	// a_0 is the Bohr radius, and rho is the radius.
	return (1.0 / math.SqrtPi) * math.Pow(1.0/a0, 3.0/2.0) * math.Exp(-rho/a0)
}

// psi2SOrbital returns the psi function of the 2s-orbital (n=2, l=0, m=0).
//
// ref: https://chemistrygod.com/atomic-orbital#s-orbital
func psi2SOrbital(rho, theta, phi float64) float64 {
	// a_0 is the Bohr radius, and rho is the radius.
	return (1.0 / (math.Sqrt(32) * math.SqrtPi)) * math.Pow(1.0/a0, 3.0/2.0) * (2.0 - rho/a0) * math.Exp(-rho/(2*a0))
}

// psi3SOrbital returns the psi function of the 2s-orbital (n=3, l=0, m=0).
//
// ref: https://chemistrygod.com/atomic-orbital#s-orbital
func psi3SOrbital(rho, theta, phi float64) float64 {
	// a_0 is the Bohr radius, and rho is the radius.
	return (1.0 / (81 * math.Sqrt(3) * math.SqrtPi)) * math.Pow(1.0/a0, 3.0/2.0) * (27.0 - (18.0*rho)/a0 + (2*math.Pow(rho, 2))/math.Pow(a0, 2)) * math.Exp(-rho/(3*a0))
}

// psi2POrbital returns the psi function of the 2p-orbital (n=2, l=1,
// m={-1,0,1})
//
// ref: https://chemistrygod.com/atomic-orbital#p-orbital
func psi2POrbital(m int) func(rho, theta, phi float64) float64 {
	// a_0 is the Bohr radius, and rho is the radius.
	switch m {
	case 0:
		// 2p-orbital (n=2, l=1, m=0)
		return func(rho, theta, phi float64) float64 {
			return (1.0 / (math.Sqrt(32) * math.SqrtPi)) * math.Pow(1.0/a0, 3.0/2.0) * (rho / a0) * math.Exp(-rho/(2*a0)) * math.Cos(theta)
		}
	case -1, +1:
		// 2p-orbitals (n=2, l=1, m=+-1)
		sign := float64(m)
		return func(rho, theta, phi float64) float64 {
			return (1.0 / (math.Sqrt(64) * math.SqrtPi)) * math.Pow(1.0/a0, 3.0/2.0) * (rho / a0) * math.Exp(-rho/(2*a0)) * math.Sin(theta) * real(cmplx.Exp(complex(sign, 0)*1i*complex(phi, 0)))
		}
	}
	panic("unreachable")
}

// Orbitals returns the psi function of the orbital with the specified principal
// quantum number, n, azimuthal quantum number, l, and magnetic quantum number,
// m.
//
//    rho (ρ):   radial distance (radius)
//    theta (θ): inclination (angular)
//    phi (φ):   azimuth (angular)
func Orbitals(n, l, m int) func(rho, theta, phi float64) float64 {
	if !(n >= 1) {
		panic(fmt.Errorf("invalid n; expected n >= 1, got %d", n))
	}
	if !(0 <= l && l < n) {
		panic(fmt.Errorf("invalid l; expected 0 <= l < n, got %d", l))
	}
	if !(-l <= m && m <= l) {
		panic(fmt.Errorf("invalid m; expected -l <= m <= +l, got %d", m))
	}
	switch n {
	case 1:
		// 1s-orbital (n=1, l=0, m=0)
		return psi1SOrbital
	case 2:
		switch l {
		case 0:
			// 2s-orbital (n=2, l=0, m=0)
			return psi2SOrbital
		case 1:
			// 2p-orbitals (n=2, l=1, m={-1,0,1})
			return psi2POrbital(m)
		}
	case 3:
		switch l {
		case 0:
			// 3s-orbital (n=3, l=0, m=0)
			return psi3SOrbital
		case 1:
			panic(fmt.Errorf("support for (n=%d, l=%d, m=%d)-orbital not yet implemented", n, l, m))
		case 2:
			panic(fmt.Errorf("support for (n=%d, l=%d, m=%d)-orbital not yet implemented", n, l, m))
		}
	}
	panic(fmt.Errorf("support for (n=%d, l=%d, m=%d)-orbital not yet implemented", n, l, m))
}

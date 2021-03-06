// References.
//
//    [1]: https://chemistrygod.com/atomic-orbital

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/cmplx"
	"os"

	"github.com/mewmew/orbitals/orb"
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
	// Generate plot of radial probability for the 1s-, 2s-, 3s-, 2p-, 3p- and
	// 3d-orbitals.
	//lines := getLines()
	//if err := genPlot("radial_probability.png", lines...); err != nil {
	//	log.Fatalf("%+v", err)
	//}

	// Generate 3D-models visualizing the probability distribution of the 1s-,
	// 2s-, 3s-, 2p-, 3p- and 3d-orbitals.
	if err := genModels(); err != nil {
		log.Fatalf("%+v", err)
	}
}

// genModels generates 3D-models visualizing the probability distribution of the
// 1s-, 2s-, 3s-, 2p-, 3p- and 3d-orbitals.
func genModels() error {
	//genModel := genSphericModel
	genModel := genCartesianModel
	// 1s-orbital.
	{
		const (
			n = 1 // principal quantum number
			l = 0 // azimuthal quantum number
			m = 0 // magnetic quantum number
		)
		if err := genModel(n, l, m); err != nil {
			return errors.WithStack(err)
		}
	}
	// 2s-orbital.
	{
		const (
			n = 2 // principal quantum number
			l = 0 // azimuthal quantum number
			m = 0 // magnetic quantum number
		)
		if err := genModel(n, l, m); err != nil {
			return errors.WithStack(err)
		}
	}
	// 3s-orbital.
	{
		const (
			n = 3 // principal quantum number
			l = 0 // azimuthal quantum number
			m = 0 // magnetic quantum number
		)
		if err := genModel(n, l, m); err != nil {
			return errors.WithStack(err)
		}
	}
	// 2p-orbitals.
	{
		const (
			n = 2 // principal quantum number
			l = 1 // azimuthal quantum number
			//m = 0 // magnetic quantum number
		)
		for m := -l; m <= l; m++ {
			if err := genModel(n, l, m); err != nil {
				return errors.WithStack(err)
			}
		}
	}
	// 3p-orbitals.
	{
		const (
			n = 3 // principal quantum number
			l = 1 // azimuthal quantum number
			//m = 0 // magnetic quantum number
		)
		for m := -l; m <= l; m++ {
			if err := genModel(n, l, m); err != nil {
				return errors.WithStack(err)
			}
		}
	}
	// 3d-orbitals.
	{
		const (
			n = 3 // principal quantum number
			l = 2 // azimuthal quantum number
			//m = 0 // magnetic quantum number
		)
		for m := -l; m <= l; m++ {
			if err := genModel(n, l, m); err != nil {
				return errors.WithStack(err)
			}
		}
	}
	// Hybrid orbitals.
	//
	// sp
	for i, Psi := range psiSPHybridOrbitals {
		dstPath := fmt.Sprintf("hybrid_orbital_sp_%d.obj", i)
		if err := genCartesianHybridModel(Psi, dstPath); err != nil {
			return errors.WithStack(err)
		}
	}
	// sp^2
	for i, Psi := range psiSP2HybridOrbitals {
		dstPath := fmt.Sprintf("hybrid_orbital_sp^2_%d.obj", i)
		if err := genCartesianHybridModel(Psi, dstPath); err != nil {
			return errors.WithStack(err)
		}
	}
	// sp^3
	for i, Psi := range psiSP3HybridOrbitals {
		dstPath := fmt.Sprintf("hybrid_orbital_sp^3_%d.obj", i)
		if err := genCartesianHybridModel(Psi, dstPath); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Probability threshold.
//const threshold = 1.0e-6
const threshold = 1.0e-11

// genSphericModel generates a 3D-model visualizing the probability distribution
// of the specified (n, l, m)-orbital.
func genSphericModel(n, l, m int) error {
	pts := getSphericModel(n, l, m)
	ps := pruneSphericModel(pts, threshold)
	dstPath := getObjModelName(n, l, m)
	fmt.Printf("creating %q\n", dstPath)
	if err := writeObjFile(dstPath, ps); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// genCartesianModel generates a 3D-model visualizing the probability
// distribution of the specified (n, l, m)-orbital.
func genCartesianModel(n, l, m int) error {
	pts := getCartesianModel(n, l, m)
	ps := pruneCartesianModel(pts, threshold)
	dstPath := getObjModelName(n, l, m)
	fmt.Printf("creating %q\n", dstPath)
	if err := writeObjFile(dstPath, ps); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// genCartesianHybridModel generates a 3D-model visualizing the probability
// distribution of the specified hybrid wave function psi.
func genCartesianHybridModel(Psi func(rho, theta, phi float64) float64, dstPath string) error {
	pts := getCartesianModelWithPsi(Psi)
	ps := pruneCartesianModel(pts, threshold)
	fmt.Printf("creating %q\n", dstPath)
	if err := writeObjFile(dstPath, ps); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// getObjModelName returns a OBJ output file name for the specified (n, l,
// m)-orbital.
func getObjModelName(n, l, m int) string {
	return fmt.Sprintf("orbital_n_%d_l_%d_m_%d.obj", n, l, m)
}

// getJsonModelName returns a JSON output file name for the specified (n, l,
// m)-orbital.
func getJsonModelName(n, l, m int) string {
	return fmt.Sprintf("orbital_n_%d_l_%d_m_%d.json", n, l, m)
}

// getLines returns plotter lines for the 1s-, 2s-, 3s-, 2p-, 3p- and
// 3d-orbitals.
func getLines() []Line {
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
		for m := -l; m <= l; m++ {
			line := getLine(n, l, m)
			lines = append(lines, line)
		}
	}
	// 3p-orbitals.
	{
		const (
			n = 3 // principal quantum number
			l = 1 // azimuthal quantum number
			//m = 0 // magnetic quantum number
		)
		for m := -l; m <= l; m++ {
			line := getLine(n, l, m)
			lines = append(lines, line)
		}
	}
	// 3d-orbitals.
	{
		const (
			n = 3 // principal quantum number
			l = 2 // azimuthal quantum number
			//m = 0 // magnetic quantum number
		)
		for m := -l; m <= l; m++ {
			line := getLine(n, l, m)
			lines = append(lines, line)
		}
	}
	return lines
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
	p.X.Label.Text = "radius (pm)"
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
			// 3p-orbitals (n=3, l=1, m={-1,0,1})
			return psi3POrbital(m)
		case 2:
			// 3d-orbitals (n=3, l=2, m={-2,-1,0,1,2})
			return psi3DOrbital(m)
		}
	}
	panic(fmt.Errorf("support for (n=%d, l=%d, m=%d)-orbital not yet implemented", n, l, m))
}

// === [ s-orbitals ] ==========================================================

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

// === [ p-orbitals ] ==========================================================

// psi2POrbital returns the psi function of the 2p-orbitals (n=2, l=1,
// m={-1,0,1})
//
//    2p_z: k=0
//    2p_x: k=+1
//    2p_y: k=-1
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
		return func(rho, theta, phi float64) float64 {
			return (1.0 / (math.Sqrt(64) * math.SqrtPi)) * math.Pow(1.0/a0, 3.0/2.0) * (rho / a0) * math.Exp(-rho/(2*a0)) * math.Sin(theta) * real(cmplx.Exp(complex(float64(m), 0)*1i*complex(phi, 0)))
		}
	}
	panic("unreachable")
}

// psi3POrbital returns the psi function of the 3p-orbitals (n=3, l=1,
// m={-1,0,1})
//
// ref: https://chemistrygod.com/atomic-orbital#p-orbital
func psi3POrbital(m int) func(rho, theta, phi float64) float64 {
	// a_0 is the Bohr radius, and rho is the radius.
	switch m {
	case 0:
		// 3p-orbital (n=3, l=1, m=0)
		return func(rho, theta, phi float64) float64 {
			return (1.0 / 81.0) * (math.Sqrt(2) / math.SqrtPi) * math.Pow(1.0/a0, 3.0/2.0) * (6*rho/a0 - math.Pow(rho, 2)/math.Pow(a0, 2)) * math.Exp(-rho/(3*a0)) * math.Cos(theta)
		}
	case -1, +1:
		// 3p-orbitals (n=3, l=1, m=+-1)
		return func(rho, theta, phi float64) float64 {
			// TODO: verify if `e^{-i phi}` should be `e^{+-i phi}`
			return (1.0 / (81.0 * math.SqrtPi)) * math.Pow(1.0/a0, 3.0/2.0) * (6*rho/a0 - math.Pow(rho, 2)/math.Pow(a0, 2)) * math.Exp(-rho/(3*a0)) * math.Sin(theta) * real(cmplx.Exp(complex(float64(m), 0)*1i*complex(phi, 0)))
		}
	}
	panic("unreachable")
}

// === [ d-orbitals ] ==========================================================

// psi3DOrbital returns the psi function of the 3d-orbitals (n=3, l=2,
// m={-2,-1,0,1,2})
//
// ref: https://chemistrygod.com/atomic-orbital#d-orbital
func psi3DOrbital(m int) func(rho, theta, phi float64) float64 {
	// a_0 is the Bohr radius, and rho is the radius.
	switch m {
	case 0:
		// 3d-orbital (n=3, l=2, m=0)
		return func(rho, theta, phi float64) float64 {
			return (1.0 / (81.0 * math.Sqrt(6) * math.SqrtPi)) * math.Pow(1.0/a0, 3.0/2.0) * math.Pow(rho/a0, 2) * math.Exp(-rho/(3*a0)) * (3*math.Pow(math.Cos(theta), 2) - 1)
		}
	case -1, +1:
		// 3d-orbitals (n=3, l=2, m=+-1)
		return func(rho, theta, phi float64) float64 {
			return (1.0 / (81.0 * math.SqrtPi)) * math.Pow(1.0/a0, 3.0/2.0) * math.Pow(rho/a0, 2) * math.Exp(-rho/(3*a0)) * math.Sin(theta) * math.Cos(theta) * real(cmplx.Exp(complex(float64(m), 0)*1i*complex(phi, 0)))
		}
	case -2, +2:
		// 3d-orbitals (n=3, l=2, m=+-2)
		return func(rho, theta, phi float64) float64 {
			return (1.0 / (162.0 * math.SqrtPi)) * math.Pow(1.0/a0, 3.0/2.0) * math.Pow(rho/a0, 2) * math.Exp(-rho/(3*a0)) * math.Pow(math.Sin(theta), 2) * real(cmplx.Exp(complex(float64(m), 0)*1i*complex(phi, 0)))
		}
	}
	panic("unreachable")
}

// === [ Hybrid wave functions ] ===============================================

// --- [ sp hybrid wave function ] ---------------------------------------------

// psiSPHybridOrbitals holds the psi functions of the sp hybrid orbitals.
//
// ref: https://winter.group.shef.ac.uk/orbitron/AO-hybrids/sp/equations.html
var psiSPHybridOrbitals = []func(rho, theta, phi float64) float64{
	// sp_1
	func(rho, theta, phi float64) float64 {
		const m_z = 0 // z: m=0
		return (1.0 / math.Sqrt2) * (psi2SOrbital(rho, theta, phi) + psi2POrbital(m_z)(rho, theta, phi))
	},
	// sp_2
	func(rho, theta, phi float64) float64 {
		const m_z = 0 // z: m=0
		return (1.0 / math.Sqrt2) * (psi2SOrbital(rho, theta, phi) - psi2POrbital(m_z)(rho, theta, phi))
	},
}

// --- [ sp^2 hybrid wave function ] -------------------------------------------

// psiSP2HybridOrbitals holds the psi functions of the sp^2 hybrid orbitals.
//
// ref: https://winter.group.shef.ac.uk/orbitron/AO-hybrids/sp2/equations.html
var psiSP2HybridOrbitals = []func(rho, theta, phi float64) float64{
	// sp^2_1
	func(rho, theta, phi float64) float64 {
		const m_z = 0 // z: m=0
		return (1.0 / math.Sqrt(3)) * (psi2SOrbital(rho, theta, phi) + math.Sqrt2*psi2POrbital(m_z)(rho, theta, phi))
	},
	// sp^2_2
	func(rho, theta, phi float64) float64 {
		const (
			m_x = +1 // x: m=+1
			m_y = -1 // y: m=-1
		)
		return (1.0 / math.Sqrt(3)) * (psi2SOrbital(rho, theta, phi) - (1.0/math.Sqrt2)*psi2POrbital(m_x)(rho, theta, phi) + math.Sqrt(3.0/2.0)*psi2POrbital(m_y)(rho, theta, phi))
	},
	// sp^2_3
	func(rho, theta, phi float64) float64 {
		const (
			m_x = +1 // x: m=+1
			m_y = -1 // y: m=-1
		)
		return (1.0 / math.Sqrt(3)) * (psi2SOrbital(rho, theta, phi) - (1.0/math.Sqrt2)*psi2POrbital(m_x)(rho, theta, phi) - math.Sqrt(3.0/2.0)*psi2POrbital(m_y)(rho, theta, phi))
	},
}

// --- [ sp^3 hybrid wave function ] -------------------------------------------

// psiSP3HybridOrbitals holds the psi functions of the sp^3 hybrid orbitals.
//
// ref: https://winter.group.shef.ac.uk/orbitron/AO-hybrids/sp3/equations.html
var psiSP3HybridOrbitals = []func(rho, theta, phi float64) float64{
	// sp^3_1
	func(rho, theta, phi float64) float64 {
		const (
			m_z = 0  // z: m=0
			m_x = +1 // x: m=+1
			m_y = -1 // y: m=-1
		)
		return (1.0 / 2.0) * (psi2SOrbital(rho, theta, phi) + psi2POrbital(m_x)(rho, theta, phi) + psi2POrbital(m_y)(rho, theta, phi) + psi2POrbital(m_z)(rho, theta, phi))
	},
	// sp^3_2
	func(rho, theta, phi float64) float64 {
		const (
			m_z = 0  // z: m=0
			m_x = +1 // x: m=+1
			m_y = -1 // y: m=-1
		)
		return (1.0 / 2.0) * (psi2SOrbital(rho, theta, phi) + psi2POrbital(m_x)(rho, theta, phi) - psi2POrbital(m_y)(rho, theta, phi) - psi2POrbital(m_z)(rho, theta, phi))
	},
	// sp^3_3
	func(rho, theta, phi float64) float64 {
		const (
			m_z = 0  // z: m=0
			m_x = +1 // x: m=+1
			m_y = -1 // y: m=-1
		)
		return (1.0 / 2.0) * (psi2SOrbital(rho, theta, phi) - psi2POrbital(m_x)(rho, theta, phi) + psi2POrbital(m_y)(rho, theta, phi) - psi2POrbital(m_z)(rho, theta, phi))
	},
	// sp^3_4
	func(rho, theta, phi float64) float64 {
		const (
			m_z = 0  // z: m=0
			m_x = +1 // x: m=+1
			m_y = -1 // y: m=-1
		)
		return (1.0 / 2.0) * (psi2SOrbital(rho, theta, phi) - psi2POrbital(m_x)(rho, theta, phi) - psi2POrbital(m_y)(rho, theta, phi) + psi2POrbital(m_z)(rho, theta, phi))
	},
}

// ### [ Helper functions ] ####################################################

// writeJsonFile marshals ps into JSON format, writing to dstPath.
func writeJsonFile(dstPath string, ps []orb.CartesianPoint) error {
	f, err := os.Create(dstPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, p := range ps {
		if err := enc.Encode(p); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// writeObjFile stroes the points in OBJ format.
//
// Example file:
//
//    v 2.00000 0.00000 0.00000
//    v 2.00000 1.00000 0.00000
//    v 1.99037 0.00000 0.19603
func writeObjFile(dstPath string, ps []orb.CartesianPoint) error {
	f, err := os.Create(dstPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	bw := bufio.NewWriter(f)
	defer bw.Flush()
	for _, p := range ps {
		// TODO: Also include probablility? Perhaps as colour or transparency?
		if _, err := fmt.Fprintf(bw, "v %.1f %.1f %.1f\n", float64(p.X), float64(p.Y), float64(p.Z)); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

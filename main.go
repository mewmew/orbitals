// References.
//
//    [1]: https://chemistrygod.com/atomic-orbital
//    [2]: https://chemistrygod.com/atomic-orbital

package main

import (
	"fmt"
	"log"
	"math"
)

// Bohr radius with unit m.
const a0 = 0.0000000000529177210903 // 52.9 pm

func main() {
	const (
		n = 1 // principal quantum number
		l = 0 // azimuthal quantum number
		m = 0 // magnetic quantum number
	)
	Psi, err := Orbitals(n, l, m)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	const (
		//rho = a0   // radial distance (radius)
		theta = 0 // inclination (angular)
		phi   = 0 // azimuth (angular)
	)
	for r := 0.0; r < 395*a0; r += 0.0000000000001 {
		rho := r
		psi := Psi(rho, theta, phi)
		psi2 := math.Pow(psi, 2)
		radialProb := RadialProb(rho, psi)
		fmt.Printf("r=%.2g\n", rho)
		fmt.Printf("   psi^2:    %v\n", psi2)
		fmt.Printf("radial prob: %v\n", radialProb)
		fmt.Println()
	}
}

// RadialProb returns the radial probability based on the given radius, r, and
// psi for the s-orbital.
func RadialProb(r, psi float64) float64 {
	// area of sphere.
	area := 4.0 * math.Pi * math.Pow(r, 2)
	// radial probability = area * psi^2.
	return area * math.Pow(psi, 2)
}

// psiSOrbital returns the psi function of the s-orbital (n=1, l=0, m=0).
//
// ref: https://chemistrygod.com/atomic-orbital#s-orbital
func psiSOrbital(rho, theta, phi float64) float64 {
	// n = 1, l = 0, m = 0.

	// a_0 is the Bohr radius, and rho is the radius.
	return (1.0 / math.SqrtPi) * math.Pow(1.0/a0, 3.0/2.0) * math.Exp(-rho/a0)
}

// Orbitals returns the psi function of the orbital with the specified principal
// quantum number, n, azimuthal quantum number, l, and magnetic quantum number,
// m.
//
//    rho (ρ):   radial distance (radius)
//    theta (θ): inclination (angular)
//    phi (φ):   azimuth (angular)
func Orbitals(n, l, m int) (func(rho, theta, phi float64) float64, error) {
	if !(n >= 1) {
		return nil, fmt.Errorf("invalid n; expected n >= 1, got %d", n)
	}
	if !(0 <= l && l < n) {
		return nil, fmt.Errorf("invalid l; expected 0 <= l < n, got %d", l)
	}
	if !(-l <= m && m <= l) {
		return nil, fmt.Errorf("invalid m; expected -l <= m <= +l, got %d", m)
	}
	switch n {
	case 1:
		// n = 1, l = 0, m = 0.
		return psiSOrbital, nil
	}
	return nil, fmt.Errorf("support for (n=%d, l=%d, m=%d)-orbital not yet implemented", n, l, m)
}

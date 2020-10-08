package main

import (
	//"fmt"
	"math"

	"github.com/mewmew/orbitals/orb"
)

// convert radial to degree.
const degToRad = 2 * math.Pi / 360.0

// cartesianCoordFromSphericalCoord returns the Cartesian (x, y, z)-coordinate
// corresponding to the given spherical (rho, theta, phi)-coordinate.
func cartesianCoordFromSphericalCoord(rho, theta, phi float64) (x, y, z float64) {
	x = rho * math.Sin(phi) * math.Cos(theta)
	y = rho * math.Sin(phi) * math.Sin(theta)
	z = rho * math.Cos(phi)
	return x, y, z
}

// getModel returns a 3D-model visualizing the probability distribution of the
// electron orbital with the specified principal quantum number, n, azimuthal
// quantum number, l, and magnetic quantum number, m.
func getModel(n, l, m int) []orb.SphericalPoint {
	Psi := Orbitals(n, l, m)
	var pts []orb.SphericalPoint
	for theta := 0.0; theta <= 2*math.Pi; theta += 4.0 * degToRad {
		//fmt.Println("theta:", theta/degToRad)
		for phi := 0.0; phi <= 2*math.Pi; phi += 4.0 * degToRad {
			for rho := 0.0 * pm; rho < 1300*pm; rho += 1.0 * pm {
				psi := Psi(rho, theta, phi)
				//psi2 := math.Pow(psi, 2)
				radialProb := RadialProb(rho, psi)
				//fmt.Printf("rho=%.2g pm\n", rho/pm)
				//fmt.Printf("   psi^2:       %v\n", psi2)
				//fmt.Printf("   radial prob: %v\n", radialProb)
				//fmt.Println()
				pt := orb.SphericalPoint{
					//SphericalCoord: SphericalCoord{
					Rho:   rho,
					Theta: theta,
					Phi:   phi,
					//},
					Prob: radialProb,
				}
				pts = append(pts, pt)
			}
		}
	}
	// Normalize probability.
	totalProb := 0.0
	for i := range pts {
		totalProb += pts[i].Prob
	}
	if totalProb != 0 {
		for i := range pts {
			pts[i].Prob /= totalProb
		}
	}
	return pts
}

// pruneModel prunes points below the given threshold probability and converts
// the points from spherical coordinates to Cartesian coordinates.
func pruneModel(pts []orb.SphericalPoint, threshold float64) []orb.CartesianPoint {
	var ps []orb.CartesianPoint
	for _, pt := range pts {
		if pt.Prob < threshold {
			continue
		}
		x, y, z := cartesianCoordFromSphericalCoord(pt.Rho, pt.Theta, pt.Phi)
		p := orb.CartesianPoint{
			X:    int(math.Round(x / pm)),
			Y:    int(math.Round(y / pm)),
			Z:    int(math.Round(z / pm)),
			Prob: pt.Prob,
		}
		ps = append(ps, p)
	}
	return ps
}

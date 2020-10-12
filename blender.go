package main

import (
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

// sphericalCoordCoordFromCartesian returns the spherical (rho, theta,
// phi)-coordinate corresponding to the given Cartesian (x, y, z)-coordinate.
//
// ref: https://en.wikipedia.org/wiki/Spherical_coordinate_system#Cartesian_coordinates
func sphericalCoordCoordFromCartesian(x, y, z float64) (rho, theta, phi float64) {
	rho = math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2) + math.Pow(z, 2))
	theta = math.Atan(y / x)
	phi = math.Atan(math.Sqrt(math.Pow(x, 2)+math.Pow(y, 2)) / z)
	return rho, theta, phi
}

// getSphericModel returns a 3D-model visualizing the probability distribution
// of the electron orbital with the specified principal quantum number, n,
// azimuthal quantum number, l, and magnetic quantum number, m.
func getSphericModel(n, l, m int) []orb.SphericalPoint {
	Psi := Orbitals(n, l, m)
	var pts []orb.SphericalPoint
	for theta := 0.0; theta <= math.Pi; theta += 4.0 * degToRad {
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

// pruneSphericModel prunes points below the given threshold probability and
// converts the points from spherical coordinates to Cartesian coordinates.
func pruneSphericModel(pts []orb.SphericalPoint, threshold float64) []orb.CartesianPoint {
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

// getCartesianModel returns a 3D-model visualizing the probability distribution
// of the electron orbital with the specified principal quantum number, n,
// azimuthal quantum number, l, and magnetic quantum number, m.
func getCartesianModel(n, l, m int) []orb.CartesianPoint {
	Psi := Orbitals(n, l, m)
	var pts []orb.CartesianPoint
	const (
		step = 15 * pm
		max  = 3000 * pm
	)
	for x := -max; x <= max; x += step {
		for y := -max; y <= max; y += step {
			for z := -max; z <= max; z += step {
				rho, theta, phi := sphericalCoordCoordFromCartesian(x, y, z)
				psi := Psi(rho, theta, phi)
				//psi2 := math.Pow(psi, 2)
				radialProb := RadialProb(rho, psi)
				//fmt.Printf("rho=%.2g pm\n", rho/pm)
				//fmt.Printf("   psi^2:       %v\n", psi2)
				//fmt.Printf("   radial prob: %v\n", radialProb)
				//fmt.Println()
				pt := orb.CartesianPoint{
					X:    int(math.Round(x / step)),
					Y:    int(math.Round(y / step)),
					Z:    int(math.Round(z / step)),
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

// pruneCartesianModel prunes points below the given threshold probability.
func pruneCartesianModel(pts []orb.CartesianPoint, threshold float64) []orb.CartesianPoint {
	var ps []orb.CartesianPoint
	for _, pt := range pts {
		if pt.Prob < threshold {
			continue
		}
		ps = append(ps, pt)
	}
	return ps
}

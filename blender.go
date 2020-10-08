package main

import (
	//"fmt"
	"math"
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

// SphericalCoord is a spherical (rho, theta, phi)-coordinate.
type SphericalCoord struct {
	// Radial distance (radius)
	Rho float64
	// Inclination (angular)
	Theta float64
	// Azimuth (angular)
	Phi float64
}

// SphericalPoint is a spherical coordinate with a probability.
type SphericalPoint struct {
	// Spherical coordinate of electron.
	//SphericalCoord

	// Radial distance (radius)
	Rho float64
	// Inclination (angular)
	Theta float64
	// Azimuth (angular)
	Phi float64

	// Probability of electron occurence at the spherical coordinate.
	Prob float64
}

// CartesianPoint is a Cartesian coordinate with a probability.
type CartesianPoint struct {
	// X-, Y-, Z-coordinate in picometer.
	X, Y, Z int
	// Probability of electron occurence at the Cartesian coordinate.
	Prob float64
}

// genModel generates a 3D-model visualizing the probability distribution of the
// electron orbital with the specified principal quantum number, n, azimuthal
// quantum number, l, and magnetic quantum number, m.
func genModel(n, l, m int) []SphericalPoint {
	Psi := Orbitals(n, l, m)
	var pts []SphericalPoint
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
				pt := SphericalPoint{
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
func pruneModel(pts []SphericalPoint, threshold float64) []CartesianPoint {
	var ps []CartesianPoint
	for _, pt := range pts {
		if pt.Prob < threshold {
			continue
		}
		x, y, z := cartesianCoordFromSphericalCoord(pt.Rho, pt.Theta, pt.Phi)
		p := CartesianPoint{
			X:    int(math.Round(x / pm)),
			Y:    int(math.Round(y / pm)),
			Z:    int(math.Round(z / pm)),
			Prob: pt.Prob,
		}
		ps = append(ps, p)
	}
	return ps
}

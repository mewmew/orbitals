package main

import "math"

// cartesianCoordFromSphericalCoord returns the Cartesian (x, y, z)-coordinate
// corresponding to the given spherical (rho, theta, phi)-coordinate.
func cartesianCoordFromSphericalCoord(rho, theta, phi float64) (x, y, z float64) {
	x = rho * math.Sin(phi) * math.Cos(theta)
	y = rho * math.Sin(phi) * math.Sin(theta)
	z = rho * math.Cos(phi)
	return x, y, z
}

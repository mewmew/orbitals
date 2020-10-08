package orb

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

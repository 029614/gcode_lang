// lib/arcs.go
package gcode

import (
	"fmt"
	"math"
)

// CalculateArcCenter calculates the absolute I and J values (center of the arc) for a G03 arc
func CalculateArcCenter(startX, startY, endX, endY, radius float64) (centerX, centerY float64, err error) {
	// Step 1: Calculate the chord length (d)
	d := math.Sqrt(math.Pow(endX-startX, 2) + math.Pow(endY-startY, 2))

	// If the distance between the start and end point is greater than the diameter, return an error
	if d > 2*radius {
		err = fmt.Errorf("impossible arc: distance between points is greater than the diameter")
		return 0, 0, err
	}

	// Step 2: Find the midpoint of the chord
	midX := (startX + endX) / 2
	midY := (startY + endY) / 2

	// Step 3: Calculate the distance from the midpoint to the center of the arc (h)
	h := math.Sqrt(radius*radius - (d/2)*(d/2))

	// Step 4: Calculate the vector perpendicular to the chord
	perpX := -(endY - startY)
	perpY := endX - startX

	// Normalize the perpendicular vector
	perpLength := math.Sqrt(perpX*perpX + perpY*perpY)
	perpX /= perpLength
	perpY /= perpLength

	// Step 5: Calculate the center of the arc (X_center, Y_center)
	// We move along the perpendicular vector by distance h from the midpoint
	centerX = midX + perpX*h
	centerY = midY + perpY*h

	return centerX, centerY, nil
}

// GetRadiusFromIJ calculates the radius of an arc based on the start point (startX, startY)
// and the center point (I, J).
func GetRadiusFromIJ(startX, startY, centerX, centerY float64) float64 {
	// Use the Pythagorean theorem to calculate the radius
	radius := math.Sqrt(math.Pow(centerX-startX, 2) + math.Pow(centerY-startY, 2))
	return radius
}

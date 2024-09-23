// lib/arcs.go
package lib

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

// Line segment structure
type LineSegment struct {
	StartX, StartY float64
	EndX, EndY     float64
}

// GetIntersection finds the intersection point of two line segments
func GetIntersection(seg1, seg2 LineSegment) (float64, float64, bool) {
	A1 := seg1.EndY - seg1.StartY
	B1 := seg1.StartX - seg1.EndX
	C1 := A1*seg1.StartX + B1*seg1.StartY

	A2 := seg2.EndY - seg2.StartY
	B2 := seg2.StartX - seg2.EndX
	C2 := A2*seg2.StartX + B2*seg2.StartY

	det := A1*B2 - A2*B1
	if det == 0 {
		return 0, 0, false // Lines are parallel
	}

	x := (B2*C1 - B1*C2) / det
	y := (A1*C2 - A2*C1) / det

	return x, y, true
}

// Calculate the rolling arc parameters
func CalculateRollingArc(seg1, seg2 LineSegment, cutterDiameter float64) (startX, startY, endX, endY, i, j float64) {
	// Find intersection point
	intersectionX, intersectionY, exists := GetIntersection(seg1, seg2)
	if !exists {
		fmt.Println("No intersection found.")
		return 0, 0, 0, 0, 0, 0
	}
	fmt.Printf("intersectionX: (%.4f), intersectionY: (%.4f)\n", intersectionX, intersectionY)

	// Calculate direction vectors for both segments
	direction1X := seg1.EndX - seg1.StartX
	direction1Y := seg1.EndY - seg1.StartY
	direction2X := seg2.EndX - seg2.StartX
	direction2Y := seg2.EndY - seg2.StartY

	// Calculate lengths of both segments
	length1 := math.Sqrt(direction1X*direction1X + direction1Y*direction1Y)
	length2 := math.Sqrt(direction2X*direction2X + direction2Y*direction2Y)
	fmt.Printf("length1: (%.2f), length2: (%.2f)\n", length1, length2)

	// Normalize the direction vectors
	direction1X /= length1
	direction1Y /= length1
	direction2X /= length2
	direction2Y /= length2

	// Calculate offset from intersection point
	offset := cutterDiameter * 0.5

	// Calculate the start and end points for the arc
	startX = intersectionX + direction1X*-offset
	startY = intersectionY + direction1Y*-offset
	endX = intersectionX + direction2X*offset
	endY = intersectionY + direction2Y*offset

	normal1X, normal1Y, normal2X, normal2Y := getSegmentsNormals(seg1, seg2)

	// Calculate the center coordinates for I and J
	i = (startX - endX) * 0.5
	j = (startY + endY) * 0.5

	return startX, startY, endX, endY, i, j
}

// PointOnSegment calculates the point on a segment that is a given distance from the end of the segment.
// The segment is defined by (startX, startY) to (endX, endY).
func PointOnSegment(seg LineSegment, distance float64) (float64, float64) {
	// Calculate the direction vector from start to end
	dirX := seg.StartX - seg.EndX
	dirY := seg.StartY - seg.EndY

	// Calculate the length of the segment
	length := math.Sqrt(dirX*dirX + dirY*dirY)

	// Normalize the direction vector
	unitDirX := dirX / length
	unitDirY := dirY / length

	// Calculate the point that is 'distance' away from the end
	pointX := seg.EndX + unitDirX*distance
	pointY := seg.EndY + unitDirY*distance

	return pointX, pointY
}

// getSegmentsNormals calculates the normal vectors for two segments
func getSegmentsNormals(seg1, seg2 LineSegment) (normal1X, normal1Y, normal2X, normal2Y float64) {
	// Calculate the direction vector for segment 1
	dir1X := seg1.EndX - seg1.StartX
	dir1Y := seg1.EndY - seg1.StartY

	// Normalize the direction vector for segment 1
	length1 := math.Sqrt(dir1X*dir1X + dir1Y*dir1Y)
	unitDir1X := dir1X / length1
	unitDir1Y := dir1Y / length1

	// The counterclockwise normal vector for segment 1
	normal1X = -unitDir1Y
	normal1Y = unitDir1X

	// Calculate the direction vector for segment 2
	dir2X := seg2.EndX - seg2.StartX
	dir2Y := seg2.EndY - seg2.StartY

	// Normalize the direction vector for segment 2
	length2 := math.Sqrt(dir2X*dir2X + dir2Y*dir2Y)
	unitDir2X := dir2X / length2
	unitDir2Y := dir2Y / length2

	// The counterclockwise normal vector for segment 2
	normal2X = -unitDir2Y
	normal2Y = unitDir2X

	return normal1X, normal1Y, normal2X, normal2Y
}

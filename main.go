// main.go
package main

import (
	"fmt" // Update with your actual module path
	"gcode_lang/lib"
	"log"
)

func main() {
	// Test case
	// startX := 37.75
	// startY := 60.4375
	// endX := 37.5625
	// endY := 60.25

	startX := 96.8825
	startY := 60.1174
	endX := 96.9374
	endY := 60.25

	radius := 0.1875

	centerX, centerY, err := lib.CalculateArcCenter(startX, startY, endX, endY, radius)
	if err != nil {
		log.Fatal("Error calculating arc center:", err)
	} else {
		fmt.Printf("Center X (I): %.4f, Center Y (J): %.4f\n", centerX, centerY)
	}

	// Test values: startX, startY, I (centerX), J(centerY)
	// Calculate radius from I and J
	derived_radius := lib.GetRadiusFromIJ(startX, startY, centerX, centerY)
	fmt.Printf("Calculated radius: %.4f\n", derived_radius)

	seg1 := lib.LineSegment{StartX: 67.3825, StartY: 30.6174, EndX: 96.8825, EndY: 60.1174}
	seg2 := lib.LineSegment{StartX: 96.7499, StartY: 60.4375, EndX: 37.5625, EndY: 60.4375}
	cutterDiameter := 0.375

	startX, startY, endX, endY, i, j := lib.CalculateRollingArc(seg1, seg2, cutterDiameter)
	fmt.Printf("Start: (%.4f, %.4f), End: (%.4f, %.4f), I: %.4f, J: %.4f\n", startX, startY, endX, endY, i, j)
}

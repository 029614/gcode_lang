package main

import (
	"math/rand"
	"time"

	"github.com/029614/gcode_lang/pkg/patterns"
)

// rtrlib := data.GetRouterLibrary()
// tllib := data.GetToolLibrary()
// rtr, err := rtrlib.GetRouterByName("Multicam")
// if err != nil {
// 	panic(err)
// }
// proc := processor.NewMulticamProcessor(rtr, tllib)
// // Example demonstrating declarative use
// // println(scode.DrillingExample().GetScript())

// // Example demonstrating dynamic use with loops
// pts := make([][2]float32, 0)
// pts = append(pts, [2]float32{0, 0})
// pts = append(pts, [2]float32{10, 0})
// pts = append(pts, [2]float32{10, 10})
// pts = append(pts, [2]float32{0, 10})
// pts = append(pts, [2]float32{0, 0})
// ex := scode.CuttingExample(pts...)

// //println(ex.GetScript())

// // proc := processor.MulticamProcessor{}
// proc.PostProcess(ex)

// println(ex.GetScript())

func main() {

	pattern := [][2]float64{
		{0, 4}, // Drill at origin
		{0, 3}, // Drill at origin
		{0, 2}, // Drill at origin
		{0, 1}, // Drill at origin
		{0, 0}, // Drill at origin
		{1, 0}, // Drill one unit to the right
		{2, 0}, // Drill two units to the right
		{3, 0}, // Drill two units to the right
		{4, 0}, // Drill two units to the right
	}

	universe := randomSample(pattern, 500)

	solution := patterns.SetCover(universe, pattern)
	println("solution size: ", len(solution), ", original size: ", len(universe))
	for _, dp := range solution {
		println("solution:")
		for _, pt := range dp {
			println("(", pt[0], ", ", pt[1], ")")
		}
	}
}

func randomSample(pattern [][2]float64, groupCount int) [][2]float64 {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Map to track unique points and prevent duplicates
	pointMap := make(map[[2]float64]struct{})

	// The size of the universe (can be larger or smaller depending on your test)
	numGroups := rand.Intn(groupCount) + 5 // Number of groups
	maxOffsetX := rand.Intn(10) + 5        // Max random X offset
	maxOffsetY := rand.Intn(10) + 5        // Max random Y offset

	universe := make([][2]float64, 0)

	for i := 0; i < numGroups; i++ {
		// Create a random offset for the group
		offsetX := float64(rand.Intn(maxOffsetX))
		offsetY := float64(rand.Intn(maxOffsetY))

		// Randomly choose 1 to all elements from the pattern
		numElements := rand.Intn(len(pattern)) + 1

		// Randomly shuffle the pattern
		shuffledPattern := make([][2]float64, len(pattern))
		copy(shuffledPattern, pattern)
		rand.Shuffle(len(shuffledPattern), func(i, j int) {
			shuffledPattern[i], shuffledPattern[j] = shuffledPattern[j], shuffledPattern[i]
		})

		// Add randomly selected elements of the shuffled pattern to the universe
		for j := 0; j < numElements; j++ {
			pt := shuffledPattern[j]
			newPoint := [2]float64{
				pt[0] + offsetX,
				pt[1] + offsetY,
			}

			// Ensure the point is unique before adding it to the universe
			if _, exists := pointMap[newPoint]; !exists {
				universe = append(universe, newPoint)
				pointMap[newPoint] = struct{}{}
			}
		}
	}

	// Shuffle the universe to mix the groupings
	rand.Shuffle(len(universe), func(i, j int) {
		universe[i], universe[j] = universe[j], universe[i]
	})

	return universe
}

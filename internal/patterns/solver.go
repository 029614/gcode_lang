package patterns

import (
	"math"
	"sort"
)

// SetCover calculates the smallest set of drill patterns that cover all points in the universe.
func SetCover(universe [][2]float64, pattern [][2]float64) [][][2]float64 {
	// Sort the universe by X and then by Y
	universe = sortPoints(universe)

	// Sort the pattern by X and then by Y
	pattern = sortPoints(pattern)

	// Create a map of the pattern permutations where each permutation is a slice of indices
	permutations := getPermutations(universe, pattern)

	// Use the set cover algorithm to reduce the permutations to the smallest set that covers the universe
	solution := GreedyCoverage(permutations)
	var result [][][2]float64

	// Return the values of the smallest set
	for _, subsets := range solution {
		s := make([][2]float64, 0)
		for _, value := range subsets {
			s = append(s, universe[value])
		}
		result = append(result, s)
	}

	return result
}

// sortPoints sorts the given slice of points by X coordinate, then by Y coordinate.
func sortPoints(slice [][2]float64) [][2]float64 {
	// Sort the slice using a custom sorting function
	sort.Slice(slice, func(i, j int) bool {
		if slice[i][0] == slice[j][0] {
			return slice[i][1] < slice[j][1] // Sort by Y if X coordinates are equal
		}
		return slice[i][0] < slice[j][0] // Sort by X
	})
	return slice
}

// getPermutations creates a slice of indices for the permutations of the pattern that can cover the universe.
func getPermutations(universe [][2]float64, pattern [][2]float64) [][]int {
	var permutations [][]int

	//var uOffset [2]float64
	//var uOrigin [2]float64

	for _, iPat := range pattern {
		for i, iPoint := range universe {
			// permutation increment
			//uOffset = add(iPoint, iPat)
			//uOrigin = getOrigin(iPoint, uOffset)
			perm := []int{i}

			for _, jPat := range pattern {
				if isEqualApprox(jPat, iPat) {
					continue
				}
				for j, jPoint := range universe {
					if isEqualApprox(jPoint, iPoint) {
						continue
					}
					if canCover(jPoint, add(jPat, iPoint)) {
						perm = append(perm, j)
					}
				}
			}

			permutations = append(permutations, perm)
		}
	}
	return permutations
}

// canCover checks if a specific point in the universe can be covered by a pattern point.
func canCover(point [2]float64, patternPoint [2]float64) bool {
	// Define the coverage logic based on your requirements.
	// Here, we can simply consider that the pattern covers the point if they are at the same coordinates.
	return isEqualApproxF(point[0], patternPoint[0]) && isEqualApproxF(point[1], patternPoint[1])
}

func isEqualApproxF(a, b float64) bool {
	return math.Abs(a-b) < 0.0001
}

func isEqualApprox(a, b [2]float64) bool {
	return isEqualApproxF(a[0], b[0]) && isEqualApproxF(a[1], b[1])
}

func sub(p1, p2 [2]float64) [2]float64 {
	return [2]float64{p1[0] - p2[0], p1[1] - p2[1]}
}

func add(p1, p2 [2]float64) [2]float64 {
	return [2]float64{p1[0] + p2[0], p1[1] + p2[1]}
}

func getOrigin(global, offset [2]float64) [2]float64 {
	return sub(global, offset)
}

func getGlobal(origin, offset [2]float64) [2]float64 {
	return add(origin, offset)
}

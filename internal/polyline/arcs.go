package polyline

import (
	"math"

	"github.com/029614/gcode_lang/internal/arc"
	gd "github.com/Anaxarchus/zero-gdscript"
	"github.com/Anaxarchus/zero-gdscript/pkg/vector2"
)

func NewPathFromBulge(points [][3]float64, closed bool, maxInterval float64, minArcPoints int) *Path {

	// Check if the first and the last points are the same
	s1 := points[0][0] + points[0][1] + points[0][2]
	s2 := points[len(points)-1][0] + points[len(points)-1][1] + points[len(points)-1][2]
	if math.Abs(s1-s2) < 0.0001 { // IsEqualApprox
		// Remove the last point by slicing up to the second-to-last element
		points = points[:len(points)-1]
	}

	//pts := LoopA(points, closed, maxInterval, minSteps)
	pts := []vector2.Vector2{}
	for i := 0; i < len(points); i++ {
		nxt := gd.Wrapi(i+1, 0, len(points)-1)
		if points[i][2] > 0 {
			arc := GetArcTo(points[i], points[nxt])
			pts = append(pts, arc.Discretize(maxInterval, minArcPoints)...)
		} else {
			pts = append(pts, vector2.New(points[i][0], points[i][1]))
		}
	}

	return &Path{
		Points: pts,
		Closed: closed,
	}
}

func GetArcTo(from, to [3]float64) *arc.Arc {
	return arc.New(bulgeToArc(vector2.New(to[0], to[1]), vector2.New(from[0], from[1]), from[2]))
}

func bulgeToArc(p1, p2 vector2.Vector2, bulge float64) (vector2.Vector2, float64, float64, float64) {
	// Calculate the distance between the points
	distance := p1.DistanceTo(p2)

	// Calculate the radius of the arc
	b := bulge
	r := (distance * (1 + b*b)) / (4 * b)

	// Calculate the angle from point 1 to point 2
	angleP1P2 := math.Atan2(p2.Y-p1.Y, p2.X-p1.X)
	//angleP1P2 := p2.AngleTo(p1)

	// Calculate the center of the arc
	theta := angleP1P2 - (math.Pi/2 - 2*math.Atan(b))
	cX := p1.X + (r * math.Cos(theta))
	cY := p1.Y + (r * math.Sin(theta))
	center := vector2.New(cX, cY)

	// Calculate the start and end angles
	var startAngle, endAngle float64
	if b > 0 {
		startAngle = math.Atan2(p2.Y-cY, p2.X-cX)
		endAngle = math.Atan2(p1.Y-cY, p1.X-cX)
	} else {
		startAngle = math.Atan2(p1.Y-cY, p1.X-cX)
		endAngle = math.Atan2(p2.Y-cY, p2.X-cX)
	}

	return center, math.Abs(r), startAngle, endAngle
}

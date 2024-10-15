package path

import (
	"errors"
	"math"

	zerogdscript "github.com/Anaxarchus/zero-gdscript"
	"github.com/Anaxarchus/zero-gdscript/pkg/geometry2d"
	"github.com/Anaxarchus/zero-gdscript/pkg/vector2"
)

const MaxArcSegmentLength = 0.1
const MaxArcSegmentAngle = 0.05

type Waypoint struct {
	Position vector2.Vector2
	Index    int
	Distance float64
}

type Arc struct {
	Position             vector2.Vector2
	Radius               float64
	Points               []int
	AngleStart, AngleEnd float64
	path                 *Path
}

type Path struct {
	Points []vector2.Vector2
	Arcs   []*Arc
	Closed bool
}

func NewPath(points []vector2.Vector2, closed bool) *Path {
	return &Path{
		Points: points,
		Closed: closed,
	}
}

func NewPathFromBulgePoints(closed bool, bpoints ...[3]float64) *Path {
	var points []vector2.Vector2

	for i := 0; i < len(bpoints); i += 2 {
		j := zerogdscript.Wrapi(i+1, 0, len(bpoints))
		this := vector2.New(bpoints[i][0], bpoints[i][1])
		next := vector2.New(bpoints[j][0], bpoints[j][1])
		if bpoints[i][2] != 0.0 {
			arc := BulgeToArc(this, next, bpoints[i][2])
			arcPoints := ArcToPoints(arc.Position, arc.Radius, arc.AngleStart, arc.AngleEnd)
			points = append(points, arcPoints...)
		} else {
			points = append(points, this, next)
		}
	}

	return NewPath(points, closed)
}

func (p *Path) Walk(from vector2.Vector2, distance float64) []Waypoint {
	var result []Waypoint
	idx := -1

	// Get nearest edge and position on the path
	edge := p.GetNearestEdge(from)
	nidx := edge[1]
	pos := p.GetNearestPosition(from)

	// Determine the starting index of the nearest position
	if pos == p.Points[edge[0]] {
		idx = edge[0]
	} else if pos == p.Points[edge[1]] {
		idx = edge[1]
	}

	npos := p.Points[nidx]
	d2 := distance * distance // Work with squared distances for efficiency
	var d float64

	// Start walking along the path until the distance is covered
	for d2 > 0.0 {
		// Append the current waypoint
		result = append(result, *p.NewWaypoint(pos, idx, p.LengthToIndex(idx)))

		// Calculate the squared distance to the next point
		d = pos.DistanceSquaredTo(npos)
		if d < d2 {
			// Move to the next point if the remaining distance is larger
			pos = npos
			idx = nidx

			// Handle closed paths (looping back to the beginning)
			if p.Closed {
				nidx = zerogdscript.Wrapi(nidx+1, 0, len(p.Points))
			} else {
				// Ensure we don't go out of bounds for open paths
				nidx = min(nidx+1, len(p.Points)-1)
			}
			npos = p.Points[nidx]

			// Break if we've reached the last point
			if npos == pos {
				break
			}
		} else {
			// If the remaining distance fits within the current segment
			direction := pos.DirectionTo(npos)
			pos = pos.Add(direction.Mulf(math.Sqrt(d2))) // Use d2 directly
			result = append(result, *p.NewWaypoint(pos, -1, p.LengthToIndex(idx)+math.Sqrt(d2)))
			break
		}

		// Update remaining distance
		d2 -= d
	}

	return result
}

func (p *Path) NewWaypoint(position vector2.Vector2, index int, distance float64) *Waypoint {
	return &Waypoint{
		Position: position,
		Index:    index,
		Distance: distance,
	}
}

func (p *Path) NewArc(points []int) (*Arc, error) {
	arc := &Arc{
		path:   p,
		Points: points,
	}

	err := arc.Fit()
	if err != nil {
		return nil, err
	}

	p.Arcs = append(p.Arcs, arc)
	return arc, nil
}

func (p *Path) SetPoints(points []vector2.Vector2) {
	p.Points = points
	p.Arcs = []*Arc{}
	arcPoints := p.FindArcs()
	for _, pts := range arcPoints {
		_, err := p.NewArc(pts)
		if err != nil {
			println(err.Error())
		}
	}
}

func (p Path) Offset(delta float64, rollingPath bool) *Path {
	var offset []vector2.Vector2

	if p.Closed {

		joinType := geometry2d.JoinTypeMiter
		if rollingPath {
			joinType = geometry2d.JoinTypeRound
		}

		offset = geometry2d.OffsetPolygon(p.Points, delta, joinType)[0]
		if len(offset) > 0 {
			offset = append(offset, offset[0])
		}

		p.SetPoints(offset)

	}
	return &p
}

func (p *Path) LengthToIndex(index int) float64 {
	var d float64
	for i := range index - 1 {
		d += p.Points[i].DistanceSquaredTo(p.Points[i+1])
	}
	return math.Sqrt(d)
}

func (p *Path) Length() float64 {
	var l float64
	for i := 0; i < len(p.Points); i++ {
		l += p.Points[i].DistanceSquaredTo(p.Points[zerogdscript.Wrapi(i+1, 0, len(p.Points))])
	}
	return math.Sqrt(l)
}

func (p *Path) GetNearestEdge(to vector2.Vector2) [2]int {
	var result [2]int
	dist := math.Inf(1)
	d := 0.0
	var i1, i2 vector2.Vector2
	for i := range len(p.Points) {
		i1 = p.Points[i]
		j := zerogdscript.Wrapi(i+1, 0, len(p.Points)-1)
		i2 = p.Points[j]
		d = to.DistanceSquaredTo(geometry2d.GetClosestPointToSegment(to, [2]vector2.Vector2{i1, i2}))
		if d < dist {
			dist = d
			result = [2]int{i, j}
		}
	}
	return result
}

func (p *Path) GetNearestPosition(to vector2.Vector2) vector2.Vector2 {
	result := vector2.Zero()
	dist := math.Inf(1)
	var i1, i2, n vector2.Vector2
	d := 0.0
	for i := 0; i < len(p.Points); i++ {
		i1 = p.Points[i]
		i2 = p.Points[zerogdscript.Wrapi(i+1, 0, len(p.Points))]
		n = geometry2d.GetClosestPointToSegment(to, [2]vector2.Vector2{i1, i2})
		d = to.DistanceSquaredTo(n)
		if d < dist {
			dist = d
			result = n
		}
	}
	return result
}

// # Function to detect arcs from a mix of arcs and straight lines
// # Returns a list of arcs (each arc is an array of point indices)
func (p *Path) FindArcs() [][]int {
	arcs := [][]int{}
	curArc := []int{}
	lastN := vector2.Vector2{}
	d2 := MaxArcSegmentLength * MaxArcSegmentLength

	for i, pt := range p.Points {
		pn := p.Points[(i+1)%len(p.Points)]
		curN := pt.DirectionTo(pn)
		curArc = append(curArc, FindPoint(pt, p.Points))
		deltaAngle := curN.AngleTo(lastN)
		deltaLength := pt.DistanceSquaredTo(pn)
		if deltaLength > d2 || math.Abs(deltaAngle) > 0.1 {
			if len(curArc) >= 3 {
				arcs = append(arcs, curArc)
			}
			curArc = []int{}
		}
		lastN = curN
	}

	return arcs
}

func (arc *Arc) Fit() error {
	points := arc.path.Points
	e1 := points[0]             //# Start point of the arc
	e2 := points[len(points)-1] //# End point of the arc
	var np vector2.Vector2

	var arcDist float64
	for i := 0; i < len(points)-1; i++ {
		arcDist += points[i].DistanceSquaredTo(points[i+1])
	}

	var midDist = arcDist * 0.5
	var cur int
	var d float64

	for cur = 0; cur < len(points)-1; cur++ {

		//# break if end of list
		if midDist <= 0.0 {
			println("error while fitting arc")
			break
		}

		//# continue if mid dist is greater than distance to next point
		d = points[cur].DistanceSquaredTo(points[cur+1])

		if d < midDist {
			midDist -= d
			continue
		} else {
			dir := points[cur].DirectionTo(points[cur+1])
			np = points[cur].Add(dir.Mulf(midDist))
			break
		}
	}

	//# Calculate distances of the sides of the triangle formed by e1, e2, and np
	a := e1.DistanceTo(e2) //# Distance between end points
	b := e1.DistanceTo(np) //# Distance from start point to nearest point
	c := e2.DistanceTo(np) //# Distance from end point to nearest point

	//# Check for degenerate triangle (invalid product for radius calculation)
	product := (a + b + c) * (a + b - c) * (a - b + c) * (b + c - a)
	if product <= 0.0 {
		arc.Position = np
		arc.Radius = 0.0
		return errors.New("degenerate arc, aborting fit") //# Return degenerate case
	}

	//# Calculate the radius using Heron's formula and triangle area
	abc := a * b * c
	radius := abc / math.Sqrt(product)

	//# Calculate the center of the circumscribed circle
	A := e1.Sub(np)
	B := e2.Sub(np)

	//# Midpoints between e1 and np, and e2 and np
	D := A.Dot(e1.Add(np).Mulf(0.5))
	E := B.Dot(e2.Add(np).Mulf(0.5))

	//# Calculate the center using determinants (solving linear system)
	denominator := A.X*B.Y - B.X*A.Y
	if math.Abs(denominator) < 1e-7 {
		arc.Position = np
		arc.Radius = 0.0
		return errors.New("degenerate arc, aborting fit") //# Return degenerate case
	}

	center_x := (D*B.Y - E*A.Y) / denominator
	center_y := (A.X*E - B.X*D) / denominator
	center := vector2.New(center_x, center_y)

	//# Return the center and radius of the circle
	arc.Position = center
	arc.Radius = radius
	return nil
}

func ArcToPoints(position vector2.Vector2, radius, a1, a2 float64) []vector2.Vector2 {
	var points []vector2.Vector2

	//# Compute the angle between consecutive points such that the angle between
	//# point1, point2, and point3 does not exceed max_angle.
	//# Using the relation between arc length and angle, we calculate the step angle.

	//# The distance between two consecutive points (chord length)
	var chLen = 2 * radius * math.Sin(MaxArcSegmentAngle/2)

	//# Step angle corresponding to that chord length
	var stAngle = 2 * math.Asin(chLen/(2*radius))

	//# Total number of points needed to cover the arc (adjust based on the angle range)
	var aSpan = math.Abs(a2 - a1)
	var nPoints = int(math.Ceil(aSpan/stAngle)) + 1

	//# Generate points evenly spaced along the arc
	for i := range nPoints {
		var angle = a1 + float64(i)*stAngle
		if angle > a2 {
			break
		}
		var x = position.X + radius*math.Cos(angle)
		var y = position.Y + radius*math.Sin(angle)
		points = append(points, vector2.New(x, y))
	}

	return points
}

// ## @tutorial(source): http://www.lee-mac.com/bulgeconversion.html#bulgearc
func BulgeToArc(bulge_start, bulge_end vector2.Vector2, bulge float64) *Arc {
	if bulge < 1e-7 {
		return nil
	}

	//# Calculate the radius of the arc
	distance := bulge_start.DistanceTo(bulge_end)
	radius := (distance * (1.0 + bulge*bulge)) / (4.0 * bulge)

	//# Emulating the 'polar' function from AutoLISP
	angle := bulge_start.AngleToPoint(bulge_end) + math.Pi/2.0 - 2.0*math.Atan(bulge)
	center := bulge_start.Add(vector2.New(math.Cos(angle), math.Sin(angle)).Mulf(radius))

	startAngle := center.AngleToPoint(bulge_start) * 180.0 / math.Pi
	endAngle := center.AngleToPoint(bulge_end) * 180.0 / math.Pi

	if bulge < 0.0 {
		return &Arc{Position: center, AngleStart: endAngle, AngleEnd: startAngle, Radius: math.Abs(radius)}
	} else {
		return &Arc{Position: center, AngleStart: startAngle, AngleEnd: endAngle, Radius: math.Abs(radius)}
	}
}

func FindPoint(point vector2.Vector2, points []vector2.Vector2) int {
	for i, pt := range points {
		if pt.IsEqualApprox(point) {
			return i
		}
	}
	return -1
}

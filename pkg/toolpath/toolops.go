package toolpath

import (
	"math"

	"github.com/029614/gcode_lang/internal/data"
	nestparser "github.com/029614/gcode_lang/internal/parser/nest"
	"github.com/029614/gcode_lang/internal/path"
	"github.com/Anaxarchus/zero-gdscript/pkg/vector2"
)

type ArcPoint struct {
	Position vector2.Vector2
	Radius   float64
	Start    float64
	End      float64
}

type ToolpathPoint [4]float64

type ToolpathOperation struct {
	Operation *data.Operation
	Instance  []*nestparser.Operation
	Toolpath  []ToolpathPoint
}

type Polygon []vector2.Vector2

func toolpathClosedChain(toolop *ToolpathOperation) error {
	// calculate offset
	return nil
}

func toolpathChain(toolop *ToolpathOperation, closed bool) error {
	offset := getCompensation(toolop)
	paths := getPaths(toolop)
	rampDist := getRampIn(toolop)
	for _, p := range paths {

		// get ramp
		path := p.Offset(offset, true)
		if rampDist > 0 {
			if closed {
				// ramp in
			} else {
				// ramp in
			}
		}

		// get rest of path

	}
	return nil
}

func toolpathPocket(toolop *ToolpathOperation) error {
	return nil
}

func toolpathArc(toolop *ToolpathOperation) error {
	return nil
}

func sortPaths(paths []*path.Path) []*path.Path {
	// sort paths by nearest neighbor
	return nil
}

func getCompensation(toolop *ToolpathOperation) float64 {
	// calculate offset

	// get the tool instance
	t, err := data.GetToolLibrary().GetToolByName(toolop.Operation.Tool)
	if err != nil {
		return 0.0
	}

	// get tool radius
	tcomp := t.CutDiameter * 0.5

	// calculate offset geometries
	if toolop.Operation.Offset == "right" {
		return tcomp
	} else if toolop.Operation.Offset == "left" {
		return -tcomp
	}
	return 0.0
}

func getPaths(toolop *ToolpathOperation) []*path.Path {
	pths := make([]*path.Path, 0)
	for _, ins := range toolop.Instance {
		var pts [][3]float64
		for _, pt := range ins.Geometry.(nestparser.ChainGeometry).Points {
			pts = append(pts, [3]float64{pt.X, pt.Y, pt.Bulge})
		}
		path := path.NewPathFromBulgePoints(ins.Geometry.(nestparser.ChainGeometry).Closed == 1, pts...)
		pths = append(pths, path)
	}
	return pths
}

func getRampIn(toolop *ToolpathOperation) float64 {

}

// # Helper function to calculate the ramp length given a height difference and angle
func getRampLength(to, from, rampRadians float64) float64 {
	return (from - to) / math.Tan(rampRadians)
}

package toolpath

import (
	"github.com/029614/gcode_lang/internal/data"
	nestparser "github.com/029614/gcode_lang/internal/parser/nest"
	"github.com/029614/gcode_lang/internal/polyline"
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
	for _, path := range paths {
		err := path.Offset(offset, false)
		if err != nil {
			println("error offsetting path in toolpathChain")
			return err
		}
	}
	return nil
}

func toolpathPocket(toolop *ToolpathOperation) error {
	return nil
}

func toolpathArc(toolop *ToolpathOperation) error {
	return nil
}

func sortPaths(paths []*polyline.Path) []*polyline.Path {
	// sort paths by nearest neighbor
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

func getPaths(toolop *ToolpathOperation) []*polyline.Path {
	pths := make([]*polyline.Path, 0)
	for _, ins := range toolop.Instance {
		var pts [][3]float64
		for _, pt := range ins.Geometry.(nestparser.ChainGeometry).Points {
			pts = append(pts, [3]float64{pt.X, pt.Y, pt.Bulge})
		}
		path := polyline.NewPathFromBulge(pts, ins.Geometry.(nestparser.ChainGeometry).Closed == 1, 0.1, 64)
		pths = append(pths, path)
	}
	return pths
}

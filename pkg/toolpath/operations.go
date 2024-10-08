package toolpath

import (
	nest "github.com/029614/gcode_lang/internal/parser/nest"
	nestparser "github.com/029614/gcode_lang/internal/parser/nest"
	"github.com/029614/gcode_lang/pkg/data"
	"github.com/Anaxarchus/zero-gdscript/pkg/vector2"
)

type ToolpathOperation struct {
	Operation *data.Operation
	Instance  []*nestparser.Operation
	Toolpath  []ToolpathPoint
}

type Polygon []vector2.Vector2

func toolpathPartCut(to *ToolpathOperation) error {
	// Toolpath function

	// get bounding box
	bb := nest.GetOperationSliceBoundingBox(to.Instance)

	// calculate category
	// cat := getCategory(bb)

	// check if should onion skin
	os := shouldOnion(bb)

	// check if should down cut
	dc := shouldDownCut(bb)

	return nil
}

func toolpathBlockDrillSystem(to *ToolpathOperation) error {
	// Toolpath function
	return nil
}

func toolpathBlockDrillPilot(to *ToolpathOperation) error {
	// Toolpath function
	return nil
}

func toolpathRabbet(to *ToolpathOperation) error {
	// Toolpath function
	return nil
}

func toolpathGroove(to *ToolpathOperation) error {
	// Toolpath function
	return nil
}

func toolpathDadoBack(to *ToolpathOperation) error {
	// Toolpath function
	return nil
}

func toolpathDrawBolts(to *ToolpathOperation) error {
	// Toolpath function
	return nil
}

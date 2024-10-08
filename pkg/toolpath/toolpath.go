package toolpath

import (
	"errors"
	"fmt"

	nestparser "github.com/029614/gcode_lang/internal/parser/nest"
	"github.com/029614/gcode_lang/pkg/data"
	"github.com/Anaxarchus/zero-gdscript/pkg/rect2"
)

const OnionMinLength = 7.25
const OnionMinArea = 700.0

type Rect2 = rect2.Rect2

type ToolpathPoint [4]float64

type ToolpathSolution []*ToolpathSheet

type ToolpathSheet []*ToolpathOperation

type PartCategory int

const (
	PartCategoryHuge PartCategory = iota
	PartCategoryLarge
	PartCategoryMedium
	PartCategorySmall
	PartCategoryTiny
)

func Toolpath(nest *nestparser.Nest, data *data.Data) (*ToolpathSolution, error) {
	sol := ToolpathSolution{}
	// Toolpath function
	for _, sheet := range nest.Sheets {
		tsh, err := toolpathSheet(sheet, data)
		sol = append(sol, tsh)
		if err != nil {
			return nil, err
		}
	}
	return &sol, nil
}

func toolpathSheet(sheet *nestparser.Sheet, data *data.Data) (*ToolpathSheet, error) {
	tsh := ToolpathSheet{}

	// Toolpath function
	var opMap = make(map[string][]*nestparser.Operation)
	for _, opName := range data.OperationLibrary.ListOperationsByName() {
		opMap[opName] = make([]*nestparser.Operation, 0)
	}

	// Compile chains, arcs, and points
	for _, part := range sheet.Parts {
		for _, chain := range part.Geometry.Chains {
			opMap[chain.Operation] = append(opMap[chain.Operation], &chain)
		}
		for _, arc := range part.Geometry.Arcs {
			opMap[arc.Operation] = append(opMap[arc.Operation], &arc)
		}
		for _, point := range part.Geometry.Points {
			opMap[point.Operation] = append(opMap[point.Operation], &point)
		}
	}

	for opName, operations := range opMap {
		dop, _ := data.OperationLibrary.GetOperationByName(opName)
		top := ToolpathOperation{
			Operation: dop,
			Instance:  operations,
		}

		err := top.toolpath()
		if err != nil {
			fmt.Println(err)
		}

		tsh = append(tsh, &top)
	}

	return &tsh, nil
}

func (to *ToolpathOperation) toolpath() error {
	switch to.Operation.Name {
	case "PartCut":
		return toolpathPartCut(to)
	case "BLOCKDRILLSYSTEM":
		return toolpathBlockDrillSystem(to)
	case "BLOCKDRILLPILOT":
		return toolpathBlockDrillPilot(to)
	case "RABBET2525":
		return toolpathRabbet(to)
	case "GROOVE25":
		return toolpathGroove(to)
	case "DadoBack":
		return toolpathDadoBack(to)
	case "DRAWBOLTS":
		return toolpathDrawBolts(to)
	default:
		return errors.New("invalid operation")
	}
}

func getCategory(rect Rect2) PartCategory {
	a := rect.GetArea()
	if a < 150.0 {
		return PartCategoryTiny
	} else if a < 300.0 || rect.Size.X < 6.0 || rect.Size.Y < 6.0 {
		return PartCategorySmall
	} else if a < 700.0 || rect.Size.X < 7.25 || rect.Size.Y < 7.25 {
		return PartCategoryMedium
	} else {
		return PartCategoryLarge
	}
}

func getRampDistance(angle, height float64) (float64, error) {
	// Toolpath function
	return 0.0, nil
}

func shouldOnion(rect Rect2) bool {
	switch getCategory(rect) {
	case PartCategoryHuge, PartCategoryLarge:
		return false
	default:
		return true
	}
}

func shouldDownCut(rect Rect2) bool {
	switch getCategory(rect) {
	case PartCategorySmall, PartCategoryTiny:
		return true
	default:
		return false
	}
}

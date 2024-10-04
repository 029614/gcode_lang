package toolpath

import (
	nestparser "github.com/029614/gcode_lang/internal/parser/nest"
	"github.com/029614/gcode_lang/pkg/data"
	"github.com/Anaxarchus/zero-gdscript/pkg/vector3"
)

type OpPath struct {
	points []vector3.Vector3
	tool   *data.Tool
}

type Path struct {
	points []vector3.Vector3
}

func Toolpath(nest *nestparser.Nest) error {
	// Toolpath function
	for _, sheet := range nest.Sheets {
		err := toolpathSheet(sheet)
		if err != nil {
			return err
		}
	}
	return nil
}

func toolpathSheet(sheet *nestparser.Sheet) error {
	// Toolpath function
	chains := make([]*nestparser.Operation, 0)
	arcs := make([]*nestparser.Operation, 0)
	points := make([]*nestparser.Operation, 0)

	// compile chains, arcs, and points
	for _, part := range sheet.Parts {
		for _, chain := range part.Geometry.Chains {
			chains = append(chains, &chain)
		}
		for _, arc := range part.Geometry.Arcs {
			arcs = append(arcs, &arc)
		}
		for _, point := range part.Geometry.Points {
			points = append(points, &point)
		}
	}

	// process chains, arcs, and points
	err := toolpathChains(chains)
	if err != nil {
		return err
	}
	err = toolpathArcs(arcs)
	if err != nil {
		return err
	}
	err = toolpathPoints(points)
	if err != nil {
		return err
	}
	return nil
}

func toolpathChains(chains []*nestparser.Operation) error {
	// Toolpath function
	ops := map[string]int{}
	for _, chain := range chains {
		// do something with the chain
		ops[chain.Operation]++
	}
	println("chains:")
	// print the keys
	for op := range ops {
		println("\t", op)
	}
	return nil
}

func toolpathArcs(arcs []*nestparser.Operation) error {
	// Toolpath function
	println("arcs:")
	ops := map[string]int{}
	for _, arc := range arcs {
		// do something with the arc
		ops[arc.Operation]++
	}
	// print the keys
	for op := range ops {
		println("\t", op)
	}
	return nil
}

func toolpathPoints(points []*nestparser.Operation) error {
	// Toolpath function
	ops := map[string]int{}
	for _, point := range points {
		// do something with the point
		ops[point.Operation]++
	}
	// print the keys
	for op := range ops {
		println("\t", op)
	}
	return nil
}

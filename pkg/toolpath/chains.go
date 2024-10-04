package toolpath

import (
	nestparser "github.com/029614/gcode_lang/internal/parser/nest"
	"github.com/Anaxarchus/zero-gdscript/pkg/vector2"
)

func toolpathOpenChain(chain *nestparser.Operation, from vector2.Vector2) error {
	path := make([][3]float64, 0)
	// Toolpath function
	return nil
}

func toolpathClosedChain(chain *nestparser.Operation, from vector2.Vector2) error {
	path := make([][3]float64, 0)
	// Toolpath function
	return nil
}

func calculateRampDistance(angle, height float64) (float64, error) {
	// Toolpath function
	return 0.0, nil
}

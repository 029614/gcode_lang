package main

import (
	"errors"

	nestparser "github.com/029614/gcode_lang/internal/parser/nest"
	"github.com/029614/gcode_lang/pkg/data"
	"github.com/029614/gcode_lang/pkg/toolpath"
)

func main() {
	data := data.NewData()
	parts, err := nestparser.LoadPartList("/Users/samuelmcclure/dev/GitHub/gcode_lang/tests/sawboxtestingCasework/output/PartOutput_casework_parts.json")
	if err != nil {
		panic(err)
	}
	println("parts loaded")

	err = parts.LoadNest("/Users/samuelmcclure/dev/GitHub/gcode_lang/tests/sawboxtestingCasework/output/PartOutput_casework_parts_output.json")
	if err != nil {
		panic(err)
	}
	println("nest loaded")

	tp, err := toolpath.Toolpath(parts.Nest, data)
	if err != nil {
		panic(err)
	}

	println(tp)
}

func ValidateParts(partlist *nestparser.PartList) error {
	for _, sheet := range partlist.Nest.Sheets {
		for _, part := range sheet.Parts {
			if !Contains(partlist.Parts, part) {
				return errors.New("part not found in PartList")
			}
		}
	}
	return nil
}

func Contains(pslice []*nestparser.Part, part *nestparser.Part) bool {
	for _, p := range pslice {
		if p == part {
			return true
		}
	}
	return false
}

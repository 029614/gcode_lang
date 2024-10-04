package main

import (
	"errors"

	nestparser "github.com/029614/gcode_lang/internal/parser/nest"
	"github.com/029614/gcode_lang/pkg/toolpath"
)

func main() {
	partlist, err := nestparser.LoadPartList("/Users/samuelmcclure/dev/GitHub/gcode_lang/tests/sawboxtestingCasework/output/PartOutput_casework_parts.json")
	if err != nil {
		panic(err)
	}

	err = partlist.LoadNest("/Users/samuelmcclure/dev/GitHub/gcode_lang/tests/sawboxtestingCasework/output/PartOutput_casework_parts_output.json")
	if err != nil {
		panic(err)
	}

	err = ValidateParts(partlist)
	if err != nil {
		panic(err)
	}

	err = toolpath.Toolpath(partlist.Nest)
	if err != nil {
		panic(err)
	}

	println("done")
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

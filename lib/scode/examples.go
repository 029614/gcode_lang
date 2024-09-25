package scode

import "fmt"

func DrillingExample() string {
	ot := NewOperationTree(

		// Starting the job
		NewOperation(
			NewCommand(
				NewInstruction(NewToken(ID_COMMENT, "Drilling Example")),
				NewInstruction(NewToken(ID_JOB_START, "")),
			),
		),

		// Configuring the gang drill and drilling holes
		NewOperation(
			NewCommand(
				NewInstruction(
					NewToken(ID_DRILL, ""),
					NewToken(ID_PARAMETER_TOOL, "1"),
					NewToken(ID_PARAMETER_FEED, "100"),
					NewToken(ID_PARAMETER_SPEED, "100"),
				),
				NewInstruction(
					NewToken(ID_DRILL, ""),
					NewToken(ID_PARAMETER_X, "10"),
					NewToken(ID_PARAMETER_Y, "10"),
					NewToken(ID_PARAMETER_Z, "0"),
				),
				NewInstruction(
					NewToken(ID_DRILL, ""),
					NewToken(ID_PARAMETER_X, "20"),
					NewToken(ID_PARAMETER_Y, "20"),
					NewToken(ID_PARAMETER_Z, "0"),
				),
				NewInstruction(
					NewToken(ID_DRILL, ""),
					NewToken(ID_PARAMETER_X, "30"),
					NewToken(ID_PARAMETER_Y, "30"),
					NewToken(ID_PARAMETER_Z, "0"),
				),
			),
		),

		// Ending the job
		NewOperation(
			NewCommand(
				NewInstruction(
					NewToken(ID_JOB_END, ""),
					NewToken(ID_COMMENT, "End of Drilling Example"),
				),
			),
		),
	)

	return ot.GetScript()
}

func CuttingExample(pts ...[2]float32) string {

	ot := NewOperationTree()

	// Starting the job
	op := ot.NewOperation()
	com := op.NewCommand()
	com.NewInstruction(NewToken(ID_COMMENT, "// Starting the job"))
	com.NewInstruction(
		NewToken(ID_JOB_START, ""),
		NewToken(ID_COMMENT, "Cutting Example"),
	)

	// Configuring the spindle
	op = ot.NewOperation()
	com = op.NewCommand()
	com.NewInstruction(NewToken(ID_COMMENT, "// Configuring the spindle"))
	com.NewInstruction(
		NewToken(ID_SPINDLE, ""),
		NewToken(ID_PARAMETER_TOOL, "1"),
		NewToken(ID_PARAMETER_FEED, "100"),
		NewToken(ID_PARAMETER_SPEED, "100"),
	)

	// Slewing into position
	com.NewInstruction(NewToken(ID_LINEBREAK, ""), NewToken(ID_COMMENT, "// Slewing into position"))
	com.NewInstruction(
		NewToken(ID_MOVE, ""),
		NewToken(ID_PARAMETER_X, fmt.Sprintf("%f", pts[0][0])),
		NewToken(ID_PARAMETER_Y, fmt.Sprintf("%f", pts[0][1])),
		NewToken(ID_PARAMETER_Z, "1"),
	)

	// Iterating over waypoints and cutting through them
	com.NewInstruction(NewToken(ID_LINEBREAK, ""), NewToken(ID_COMMENT, "// Iterating over waypoints and cutting through them"))
	for _, pt := range pts {
		com.NewInstruction(
			NewToken(ID_CUT, ""),
			NewToken(ID_PARAMETER_X, fmt.Sprintf("%f", pt[0])),
			NewToken(ID_PARAMETER_Y, fmt.Sprintf("%f", pt[1])),
			NewToken(ID_PARAMETER_Z, "0"),
		)
	}
	// Retracting the spindle
	com.NewInstruction(NewToken(ID_LINEBREAK, ""), NewToken(ID_COMMENT, "// Retracting the spindle"))
	com.NewInstruction(
		NewToken(ID_MOVE, ""),
		NewToken(ID_PARAMETER_X, fmt.Sprintf("%f", pts[len(pts)-1][0])),
		NewToken(ID_PARAMETER_Y, fmt.Sprintf("%f", pts[len(pts)-1][1])),
		NewToken(ID_PARAMETER_Z, "1"),
	)

	// Ending the job
	op = ot.NewOperation()
	com = op.NewCommand()
	com.NewInstruction(NewToken(ID_COMMENT, "// Ending the job"))
	com.NewInstruction(NewToken(ID_JOB_END, ""), NewToken(ID_COMMENT, "End of Cutting Example"))

	return ot.GetScript()
}

package scode

import "fmt"

func DrillingExample() *OperationTree {
	ot := NewOperationTree(

		// Starting the job
		NewOperation(OT_START,
			NewCommand(CT_START,
				NewInstruction(NewToken(ID_COMMENT, "Drilling Example")),
				NewInstruction(NewToken(ID_JOB_START, "")),
			),
		),

		// Configuring the gang drill and drilling holes
		NewOperation(OT_DRILL,
			NewCommand(CT_DRILLSET,
				NewInstruction(
					NewToken(ID_DRILL, ""),
					NewToken(ID_PARAMETER_TOOL, "1"),
					NewToken(ID_PARAMETER_FEED, "100"),
					NewToken(ID_PARAMETER_SPEED, "100"),
				),
			),
			NewCommand(CT_DRILLMOTION,
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
		NewOperation(OT_END,
			NewCommand(CT_STOP,
				NewInstruction(
					NewToken(ID_JOB_END, ""),
					NewToken(ID_COMMENT, "End of Drilling Example"),
				),
			),
		),
	)

	return ot
}

func CuttingExample(pts ...[2]float32) *OperationTree {

	ot := NewOperationTree()

	// Starting the job
	op := ot.NewOperation(OT_START)
	com := op.NewCommand(CT_START)
	com.NewInstruction(NewToken(ID_COMMENT, "// Starting the job"))
	com.NewInstruction(
		NewToken(ID_JOB_START, ""),
		NewToken(ID_COMMENT, "Cutting Example"),
	)

	// Configuring the spindle
	op = ot.NewOperation(OT_SPINDLE)
	com = op.NewCommand(CT_SPINDLESET)
	com.NewInstruction(NewToken(ID_COMMENT, "// Configuring the spindle"))
	com.NewInstruction(
		NewToken(ID_SPINDLE, ""),
		NewToken(ID_PARAMETER_TOOL, "1"),
		NewToken(ID_PARAMETER_SPEED, "100"),
	)

	// Slewing into position
	com = op.NewCommand(CT_SPINDLEMOTION)
	com.NewInstruction(NewToken(ID_LINEBREAK, ""), NewToken(ID_COMMENT, "// Slewing into position"))
	com.NewInstruction(
		NewToken(ID_MOVE, ""),
		NewToken(ID_PARAMETER_X, fmt.Sprintf("%f", pts[0][0])),
		NewToken(ID_PARAMETER_Y, fmt.Sprintf("%f", pts[0][1])),
		NewToken(ID_PARAMETER_Z, "1"),
		NewToken(ID_PARAMETER_FEED, "100"),
	)

	// Iterating over waypoints and cutting through them
	com.NewInstruction(NewToken(ID_LINEBREAK, ""), NewToken(ID_COMMENT, "// Iterating over waypoints and cutting through them"))
	for _, pt := range pts {
		com.NewInstruction(
			NewToken(ID_CUT, ""),
			NewToken(ID_PARAMETER_X, fmt.Sprintf("%f", pt[0])),
			NewToken(ID_PARAMETER_Y, fmt.Sprintf("%f", pt[1])),
			NewToken(ID_PARAMETER_Z, "0"),
			NewToken(ID_PARAMETER_FEED, "100"),
		)
	}
	// Retracting the spindle
	com.NewInstruction(NewToken(ID_LINEBREAK, ""), NewToken(ID_COMMENT, "// Retracting the spindle"))
	com.NewInstruction(
		NewToken(ID_MOVE, ""),
		NewToken(ID_PARAMETER_X, fmt.Sprintf("%f", pts[len(pts)-1][0])),
		NewToken(ID_PARAMETER_Y, fmt.Sprintf("%f", pts[len(pts)-1][1])),
		NewToken(ID_PARAMETER_Z, "1"),
		NewToken(ID_PARAMETER_FEED, "100"),
	)

	// Ending the job
	op = ot.NewOperation(OT_END)
	com = op.NewCommand(CT_STOP)
	com.NewInstruction(NewToken(ID_COMMENT, "// Ending the job"))
	com.NewInstruction(NewToken(ID_JOB_END, ""), NewToken(ID_COMMENT, "End of Cutting Example"))

	return ot
}

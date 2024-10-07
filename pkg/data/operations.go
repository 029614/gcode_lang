package data

import (
	"errors"
)

type OperationLibrary []*Operation

type Operation struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Tool       string  `json:"tool"`
	Ramp       float64 `json:"ramp"`
	FeedRate   int     `json:"feed_rate"`
	PlungeRate int     `json:"plunge_rate"`
	SpindleRPM int     `json:"spindle_rpm"`
	Offset     string  `json:"offset"`
	CutDepth   float64 `json:"cut_depth"`
	CutHeight  float64 `json:"cut_height"`
	FeedHeight float64 `json:"feed_height"`
}

func (ol *OperationLibrary) GetOperationByName(name string) (*Operation, error) {
	for _, operation := range *ol {
		if operation.Name == name {
			return operation, nil
		}
	}
	return nil, errors.New("Operation not found")
}

func (ol *OperationLibrary) ListOperationsByName() []string {
	var operationList []string
	for _, operation := range *ol {
		operationList = append(operationList, operation.Name)
	}
	return operationList
}

func GetOperationsLibrary() *OperationLibrary {
	filePath := "./tests/resources/operations.json"
	var operationLibrary OperationLibrary
	unmarshalJson(filePath, &operationLibrary)
	return &operationLibrary
}

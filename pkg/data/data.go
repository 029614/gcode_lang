package data

import (
	"encoding/json"
	"log"
	"os"
)

type Data struct {
	ToolLibrary      *ToolLibrary
	OperationLibrary *OperationLibrary
	RouterLibrary    *RouterLibrary
}

func NewData() *Data {
	return &Data{
		ToolLibrary:      GetToolLibrary(),
		OperationLibrary: GetOperationsLibrary(),
		RouterLibrary:    GetRouterLibrary(),
	}
}

func unmarshalJson(filepath string, v interface{}) {
	byteValue, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	if err := json.Unmarshal(byteValue, v); err != nil {
		log.Fatalf("failed to unmarshal file: %v", err)
	}
}

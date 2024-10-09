package data

import (
	"errors"
)

type ToolLibrary []*Tool

type Tool struct {
	ID            string                 `json:"id"`
	CutDiameter   float64                `json:"cut_diameter"`
	ShankDiameter float64                `json:"shank_diameter"`
	CutLength     float64                `json:"cut_length"`
	Flutes        int                    `json:"flutes"`
	FluteType     string                 `json:"flute_type"`
	Shape         string                 `json:"shape"`
	Material      string                 `json:"material"`
	MaxRPM        int                    `json:"max_rpm"`
	Name          string                 `json:"name"`
	Meta          map[string]interface{} `json:"meta"`
	Supplier      string                 `json:"supplier"`
	Model         string                 `json:"model"`
}

func (tl *ToolLibrary) GetToolByID(id string) (*Tool, error) {
	for _, tool := range *tl {
		if tool.ID == id {
			return tool, nil
		}
	}
	return nil, errors.New("Tool not found")
}

func (tl *ToolLibrary) GetToolByName(name string) (*Tool, error) {
	for _, tool := range *tl {
		if tool.Name == name {
			return tool, nil
		}
	}
	return nil, errors.New("Tool not found")
}

func (tl *ToolLibrary) ListToolsByName() []string {
	var toolList []string
	for _, tool := range *tl {
		toolList = append(toolList, tool.Name)
	}
	return toolList
}

// GetToolLibrary loads the ToolLibrary from the specified mock file.
func GetToolLibrary() *ToolLibrary {
	filePath := "./tests/resources/toollib.json"
	var toolLibrary ToolLibrary
	unmarshalJson(filePath, &toolLibrary)
	return &toolLibrary
}

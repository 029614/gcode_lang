package data

import (
	"encoding/json"
	"errors"
	"log"
	"os"
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

type RouterLibrary []*Router

type Router struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Spindle RouterSpindleData `json:"spindle"`
	Gang    RouterGangData    `json:"gangdrill"`
}

type SlotData string
type GangSlotData struct {
	ToolID  string  `json:"tool"`
	OffsetX float64 `json:"x"`
	OffsetY float64 `json:"y"`
}

type RouterSpindleData struct {
	One    SlotData `json:"1"`
	Two    SlotData `json:"2"`
	Three  SlotData `json:"3"`
	Four   SlotData `json:"4"`
	Five   SlotData `json:"5"`
	Six    SlotData `json:"6"`
	Seven  SlotData `json:"7"`
	Eight  SlotData `json:"8"`
	Nine   SlotData `json:"9"`
	Ten    SlotData `json:"10"`
	Eleven SlotData `json:"11"`
	Twelve SlotData `json:"12"`
}

type RouterGangData struct {
	One   GangSlotData `json:"1"`
	Two   GangSlotData `json:"2"`
	Three GangSlotData `json:"3"`
	Four  GangSlotData `json:"4"`
	Five  GangSlotData `json:"5"`
	Six   GangSlotData `json:"6"`
	Seven GangSlotData `json:"7"`
	Eight GangSlotData `json:"8"`
	Nine  GangSlotData `json:"9"`
}

func (tl *ToolLibrary) GetTool(id string) (*Tool, error) {
	for _, tool := range *tl {
		if tool.ID == id {
			return tool, nil
		}
	}
	return nil, errors.New("Tool not found")
}

func (rl *RouterLibrary) GetRouterByID(id string) (*Router, error) {
	for _, router := range *rl {
		if router.ID == id {
			return router, nil
		}
	}
	return nil, errors.New("Router not found")
}

func (rl *RouterLibrary) GetRouterByName(name string) (*Router, error) {
	for _, router := range *rl {
		if router.Name == name {
			return router, nil
		}
	}
	return nil, errors.New("Router not found")
}

func (r *Router) GetSpindleSlot(idx int) SlotData {
	switch idx {
	case 1:
		return r.Spindle.One
	case 2:
		return r.Spindle.Two
	case 3:
		return r.Spindle.Three
	case 4:
		return r.Spindle.Four
	case 5:
		return r.Spindle.Five
	case 6:
		return r.Spindle.Six
	case 7:
		return r.Spindle.Seven
	case 8:
		return r.Spindle.Eight
	case 9:
		return r.Spindle.Nine
	case 10:
		return r.Spindle.Ten
	case 11:
		return r.Spindle.Eleven
	case 12:
		return r.Spindle.Twelve
	default:
		return ""
	}
}

func (r *Router) GetGangSlot(idx int) GangSlotData {
	switch idx {
	case 1:
		return r.Gang.One
	case 2:
		return r.Gang.Two
	case 3:
		return r.Gang.Three
	case 4:
		return r.Gang.Four
	case 5:
		return r.Gang.Five
	case 6:
		return r.Gang.Six
	case 7:
		return r.Gang.Seven
	case 8:
		return r.Gang.Eight
	case 9:
		return r.Gang.Nine
	default:
		return GangSlotData{}
	}
}

// GetToolLibrary loads the ToolLibrary from the specified mock file.
func GetToolLibrary() *ToolLibrary {
	filePath := "./tests/toollib.json"
	//file, err := os.Open(filePath)
	//if err != nil {
	//	log.Fatalf("failed to open tool library file: %v", err)
	//}
	//defer file.Close()

	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read tool library file: %v", err)
	}

	var toolLibrary ToolLibrary
	if err := json.Unmarshal(byteValue, &toolLibrary); err != nil {
		log.Fatalf("failed to unmarshal tool library: %v", err)
	}

	return &toolLibrary
}

// GetRouterLibrary loads the RouterLibrary from the specified mock file.
func GetRouterLibrary() *RouterLibrary {
	filePath := "./tests/routerlib.json"
	//file, err := os.Open(filePath)
	//if err != nil {
	//    log.Fatalf("failed to open router library file: %v", err)
	//}
	//defer file.Close()

	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read router library file: %v", err)
	}

	var routerLibrary RouterLibrary
	if err := json.Unmarshal(byteValue, &routerLibrary); err != nil {
		log.Fatalf("failed to unmarshal router library: %v", err)
	}

	return &routerLibrary
}

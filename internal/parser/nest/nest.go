package nestparser

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Anaxarchus/zero-gdscript/pkg/vector2"
)

type Point struct {
	vector2.Vector2
	Bulge float64 `json:"bulge"`
}

type BaseGeometry struct{}

type ChainGeometry struct {
	BaseGeometry
	Points []Point `json:"points"`
	Closed int     `json:"closed"`
}

type ArcGeometry struct {
	BaseGeometry
	Radius     float64 `json:"radius"`
	StartAngle float64 `json:"start_angle"`
	Sweep      float64 `json:"sweep"`
	Position   Point   `json:"position"`
}

type Operation struct {
	Geometry  BaseGeometry `json:"geometry"`
	Operation string       `json:"operation"`
	Depth     float64      `json:"depth"`
}

type PartGeometry struct {
	Points []Operation `json:"Points"`
	Chains []Operation `json:"Chains"`
	Arcs   []Operation `json:"Arcs"`
}

type Part struct {
	Name        string          `json:"name"`
	Size        vector2.Vector2 `json:"size"`
	Origin      vector2.Vector2 `json:"origin"`
	Limit       vector2.Vector2 `json:"limit"`
	IsRotated   int             `json:"isRotated"`
	Area        float64         `json:"area"`
	SheetNumber int             `json:"sheet"`

	ID         string       `json:"partId,omitempty"`
	CanRotate  int          `json:"canRotate,omitempty"`
	UnitNumber string       `json:"unitNumber,omitempty"`
	Thickness  float64      `json:"thickness,omitempty"`
	Quantity   int          `json:"quantity,omitempty"`
	UnitLetter string       `json:"unitLetter,omitempty"`
	Geometry   PartGeometry `json:"geometry,omitempty"`
}

type Sheet struct {
	SheetNumber int     `json:"sheet_number"`
	Parts       []*Part `json:"parts"`
}

type Nest struct {
	Partgap   float64         `json:"partgap"`
	Sheetsize vector2.Vector2 `json:"sheetsize"`
	Sheets    []*Sheet        `json:"sheets"`
	Nofits    []*Part         `json:"nofits"`
	Material  string          `json:"material"`
	Jobname   string          `json:"jobname"`
}

type PartList struct {
	Nest     *Nest
	FilePath string
	Partgap  float64 `json:"partgap"`
	Jobname  string  `json:"jobname"`
	Parts    []*Part `json:"parts"`
}

func (pl *PartList) FindPart(id string) *Part {
	for _, part := range pl.Parts {
		if part.ID == id {
			return part
		}
	}
	return nil
}

func (pl *PartList) LoadNest(filepath string) error {
	nest, err := LoadNest(filepath)
	if err != nil {
		return err
	}

	// Create a mapping of original parts by their IDs
	partMap := make(map[string]*Part)
	for _, part := range pl.Parts {
		partMap[part.ID] = part
	}

	// Iterate through each sheet and its parts in the loaded nest
	for _, sheet := range nest.Sheets {
		for _, fragment := range sheet.Parts {
			// Find the original part based on the ID of the fragment
			originalPart := partMap[fragment.Name]
			if originalPart != nil {
				// Update the original part with properties from the fragment
				originalPart.SheetNumber = sheet.SheetNumber
				originalPart.Origin = fragment.Origin
				originalPart.IsRotated = fragment.IsRotated
				// Update any other properties as needed
			} else {
				// Log if the original part is not found
				fmt.Printf("Original part with ID %s not found in PartList\n", fragment.Name)
			}
		}
	}

	// After updating original parts, update the nest to point to them
	for x := range nest.Sheets {
		for y, fragment := range nest.Sheets[x].Parts {
			nest.Sheets[x].Parts[y] = partMap[fragment.Name] // Point to the original part
		}
	}

	pl.Nest = nest // Finally, set the updated nest in PartList
	return nil
}

func LoadNest(filepath string) (*Nest, error) {
	// Read the file content into a byte slice
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a Nest struct
	var nest Nest
	err = json.Unmarshal(bytes, &nest)
	if err != nil {
		return nil, err
	}

	// Return the Nest struct and no error
	return &nest, nil
}

func LoadPartList(filepath string) (*PartList, error) {
	// Read the file content into a byte slice
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a Nest struct
	var plist PartList
	err = json.Unmarshal(bytes, &plist)
	if err != nil {
		return nil, err
	}

	plist.FilePath = filepath

	// Return the Nest struct and no error
	return &plist, nil
}

// Save JSON data to a file
func saveJSONToFile(data interface{}, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filepath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: pretty-printing the JSON
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to write JSON to file %s: %v", filepath, err)
	}

	return nil
}

func (pl *PartList) GenerateNest(pathToInterface string, sheetLength, sheetHeight float64) error {
	// Step 1: Create a temporary directory at the program’s location
	tempDir, err := os.MkdirTemp("", "nestparser-")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Ensure the directory is deleted after processing

	// Define file paths within the temporary directory
	inputFilePath := filepath.Join(tempDir, "inputfile.json")
	outputFilePath := filepath.Join(tempDir, "outputfile.json")

	// Step 2: Marshal the PartList into JSON and save it to the temp input file
	if err := saveJSONToFile(pl, inputFilePath); err != nil {
		return fmt.Errorf("failed to save PartList to input file: %v", err)
	}

	// Step 3: Prepare arguments for the CLI command
	args := []string{
		"--gap", fmt.Sprintf("%f", pl.Partgap),
		"--length", fmt.Sprintf("%f", sheetLength),
		"--width", fmt.Sprintf("%f", sheetHeight),
		"--output", outputFilePath,
		inputFilePath, // This is the file containing part data
	}

	// Step 4: Execute the command
	cmd := exec.Command(pathToInterface, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute external command: %v", err)
	}

	// Step 5: Load the generated output file (outputfile.json) into the PartList
	if err := pl.LoadNest(outputFilePath); err != nil {
		return fmt.Errorf("failed to load generated nest from output file: %v", err)
	}

	// Temp directory and files will be deleted automatically by `defer os.RemoveAll(tempDir)`
	return nil
}

func (pg *PartGeometry) GetBoundingBox() [2]float64 {
	var minX, minY, maxX, maxY float64

	// Initialize min and max to suitable values
	minX, minY = math.MaxFloat64, math.MaxFloat64
	maxX, maxY = -math.MaxFloat64, -math.MaxFloat64

	for _, chain := range pg.Chains {
		geo := chain.Geometry

		// Check if chain.Geometry is of type ChainGeometry
		if chainGeometry, ok := geo.(ChainGeometry); ok {
			// Use chainGeometry for bounding box calculations
			// Example: Assuming ChainGeometry has fields for coordinates
			for _, point := range chainGeometry.Points {
				if point.X < minX {
					minX = point.X
				}
				if point.Y < minY {
					minY = point.Y
				}
				if point.X > maxX {
					maxX = point.X
				}
				if point.Y > maxY {
					maxY = point.Y
				}
			}
		}
	}

	return [2]float64{minX, minY, maxX, maxY}
}

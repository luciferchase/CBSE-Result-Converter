package backend

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

func Parse(inputFilePath, class string) (parsedData [][]string, schoolResult Result,
	missingSubjectCodes []string, err error) {
	rawData, err := read(inputFilePath, class)
	if err != nil {
		return nil, Result{}, nil, err
	}
	parsedData, schoolResult, missingSubjectCodes, err = parser(rawData, class)
	if err != nil {
		return nil, Result{}, nil, err
	}

	return parsedData, schoolResult, missingSubjectCodes, nil
}

func Write(parsedData [][]string, class string, path string) error {
	var (
		outputFileName string
		outputFile     *os.File
	)

	if class == "X" && len(path) == 0 {
		outputFileName = filepath.Join("..", "output", "class_10th_result.csv")
	} else if class == "XII" && len(path) == 0 {
		outputFileName = filepath.Join("..", "output", "class_12th_result.csv")
	}

	if len(path) == 0 {
		// Create output folder if not exists
		outputFolder := filepath.Join("..", "output")
		_ = os.Mkdir(outputFolder, os.ModePerm)

		outputFile, _ = os.Create(outputFileName)
		defer outputFile.Close()
	} else {
		outputFile, _ = os.Create(path)
	}

	writer := csv.NewWriter(outputFile)

	err := writer.WriteAll(parsedData)
	if err != nil {
		log.Println(err)
		return errors.New("Error occurred while writing to .csv")
	}

	return nil
}

// To update the subject codes files
func Update(codes map[string]string) error {
	// Add to basic subject codes
	for key, value := range codes {
		basicSubjectCodes[key] = value
	}

	json, _ := json.MarshalIndent(basicSubjectCodes, "  ", "  ")
	err := os.WriteFile(filepath.Join("backend", "subject_codes.json"), json, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

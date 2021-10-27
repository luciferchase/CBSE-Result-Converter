package backend

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

// To read and Parse from the raw Input Text File and convert it to Desired .csv file
func Parse(inputFilePath, class string) (parsedData [][]string, schoolResult Result,
	missingSubjectCodes []string, err error) {
	// Convert the input text file into []string
	rawData, err := read(inputFilePath, class)
	if err != nil {
		return nil, Result{}, nil, err
	}

	// Parse the result into desired format and analyse it
	parsedData, schoolResult, missingSubjectCodes, err = parser(rawData, class)
	if err != nil {
		return nil, Result{}, nil, err
	}

	return parsedData, schoolResult, missingSubjectCodes, nil
}

// Write the file to User-Specified Location
func Write(parsedData [][]string, class string, path string) error {
	var (
		outputFileName string
		outputFile     *os.File
	)

	// Save the file in the ./output/<class-wise-result>.csv
	if class == "X" && len(path) == 0 {
		// The ".." is there to access the "output" folder which should be outside of the
		// src folder, or where the app is stored
		outputFileName = filepath.Join("..", "output", "class_10th_result.csv")
	} else if class == "XII" && len(path) == 0 {
		outputFileName = filepath.Join("..", "output", "class_12th_result.csv")
	}

	// If no output file is specified
	if len(path) == 0 {
		// Create output folder if not exists
		outputFolder := filepath.Join("..", "output")
		_ = os.Mkdir(outputFolder, os.ModePerm)

		outputFile, _ = os.Create(outputFileName)
		defer outputFile.Close()
	} else {
		// Save the file in User-Specified Location
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
	
	// Using os.ModePerm so I don't have to resort to use Magic Numbers.
	err := os.WriteFile(filepath.Join("backend", "subject_codes.json"), json, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

package backend

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var basicSubjectCodes = map[string]string{
	// Class 10
	"002": "HINDI",
	"086": "SCIENCE",
	"087": "SOCIAL SC.",
	"122": "SANSKRIT",
	"184": "ENGLISH",
	"241": "MATHEMATICS - BASIC",
	"401": "RETAILING",
	"402": "IT",

	// Common
	"041": "MATHEMATICS",

	// Class 12
	"027": "HISTORY",
	"028": "POLITICAL SCIENCE",
	"029": "GEOGRAPHY",
	"030": "ECONOMICS",
	"042": "PHYSICS",
	"043": "CHEMISTRY",
	"044": "BIOLOGY",
	"048": "PHYSICAL EDUCATION",
	"049": "PAINTING",
	"054": "BUSINESS STUDIES",
	"055": "ACCOUNTANCY",
	"083": "COMPUTER SCIENCE",
	"301": "ENGLISH CORE",
	"302": "HINDI CORE",
	"806": "TOURISM",
	"809": "FOOD PRODUCTION",
} 

func Parse(inputFilePath, class string) (parsedData [][]string, 
	missingSubjectCodes []string, err error) {
	rawData, err := read(inputFilePath, class)
	if err != nil {
		return nil, nil, err
	}
	parsedData, missingSubjectCodes, err = parser(rawData, class)
	if err != nil {
		return nil, nil, err
	}

	return parsedData, missingSubjectCodes, nil
}

func Write(parsedData [][]string, class string) error {
	var outputFileName string
	if class == "X" {
		outputFileName = filepath.Join("..", "output", "class_10th_result.csv")
	} else if class == "XII" {
		outputFileName = filepath.Join("..", "output", "class_12th_result.csv")
	}

	// Create output folder if not exists
	outputFolder := filepath.Join("..", "output")
	_ = os.Mkdir(outputFolder, os.ModePerm)

	outputFile, _ := os.Create(outputFileName)
	defer outputFile.Close()

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

func read(path, class string) ([]string, error) {
	inputFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Incorrect path!")
	}
	defer inputFile.Close()

	// Accept only .txt file
	if strings.ToLower(filepath.Ext(path)) != ".txt" {
		return nil, errors.New("Invalid Input File! (Accepts Only .txt file)")
	}

	scanner := bufio.NewScanner(inputFile)
	var rawData []string

	for scanner.Scan() {
		records := strings.Fields(scanner.Text())
		rawData = append(rawData, records...)
	}

	// Check if the given file is actually result file or not
	var validFile int
	switch {
	case class == "X":
		validFile = strings.Compare(strings.Join(rawData[4:7], " "), 
		"SECONDARY SCHOOL EXAMINATION")
	case class == "XII":
		validFile = strings.Compare(strings.Join(rawData[4:8], " "), 
		"SENIOR SCHOOL CERTIFICATE EXAMINATION")
	}

	if validFile != 0 {
		return nil, errors.New("Invalid Input File!")
	}

	return rawData, nil
}

func getSubjectCode() (map[string]string, error) {
	subjectCodesfile, err := os.ReadFile(filepath.Join("backend", "subject_codes.json"))
	if err != nil {
		// If file doesn't exist, use the basic subject codes
		return basicSubjectCodes, nil
	}

	var subjectCodes map[string]string
	err = json.Unmarshal(subjectCodesfile, &subjectCodes)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Corrupt Subject Code File!")
	}
	return subjectCodes, nil
}

func notInSlice (slice []string, ele string) bool {
	for _, i := range slice {
		if i == ele { return false }
	}
	return true
}

func parser(rawData []string, class string) (parsedData [][]string, 
	missingSubjectCodes []string, err error) {
	// For Class X every field is 3 index before their counterpart in Class XII
	var modifyIndex int
	if class == "X" {
		modifyIndex = -3
	}

	rollRegex, _ := regexp.Compile("^\\d{8}$")
	nameRegex, _ := regexp.Compile("(([A-Z])\\w+)|([A-Z]){1}")

	subjectCodes, err := getSubjectCode()
	if err != nil {
		return nil, nil, err
	}

	parsedData = [][]string{
		{
			"ROLL", "GENDER", "NAME",
			"SUBJECT CODE 1", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
			"SUBJECT CODE 2", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
			"SUBJECT CODE 3", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
			"SUBJECT CODE 4", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
			"SUBJECT CODE 5", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
			"SUBJECT CODE 6", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
			"TOTAL", "PERCENTAGE", "RESULT",
		},
	}

	for i, _ := range rawData {
		if rollRegex.MatchString(rawData[i]) {
			// Parse individual student's data

			roll := rawData[i]
			gender := rawData[i+1]

			var name string
			switch {
			// Extremely rare but if a student has five word name
			case nameRegex.MatchString(rawData[i+6]):
				name = strings.Join(rawData[i+2:i+7], " ")
				i += 3

			// Rare but if a student has four word name
			case nameRegex.MatchString(rawData[i+5]):
				name = strings.Join(rawData[i+2:i+6], " ")
				i += 2

			// If a student has three word name
			case nameRegex.MatchString(rawData[i+4]):
				name = strings.Join(rawData[i+2:i+5], " ")
				i++

			// If a student has only one word name
			case !nameRegex.MatchString(rawData[i+3]):
				name = rawData[i+2]
				i--

			// Default two word names
			default:
				name = strings.Join(rawData[i+2:i+4], " ")
			}

			// Skip R.L or other results which doesn't have any data to it
			var modifySubject int
			result := rawData[i+13+modifyIndex]
			if result != "PASS" && result != "COMP" {
				// For those students who has only 5 subjects
				if rawData[i+12+modifyIndex] == "PASS" || rawData[i+12+modifyIndex] == "COMP" {
					modifySubject = -2
					i++
				} else {
					continue
				}
			}

			var (
				studentSubjects [][]string
				skipRecord      bool
			)
			for _, j := range rawData[i+4 : i+10+modifySubject] {
				// Check if the subject code is present in the database
				if name, ok := subjectCodes[j]; ok {
					studentSubjects = append(studentSubjects, []string{j, name})
				} else {
					if notInSlice(missingSubjectCodes, j) {
						missingSubjectCodes = append(missingSubjectCodes, j)
					}
					skipRecord = true
				}
			}
			if skipRecord { continue }

			var marks [][]string
			for j := 14; j < 25+modifySubject; j += 2 {
				marks = append(marks, 
					[]string{rawData[i+j+modifyIndex], rawData[i+j+modifyIndex+1]})
			}

			var total int
			for _, j := range marks {
				num, _ := strconv.Atoi(j[0])
				total += num
			}

			var percentage float64
			if len(studentSubjects) == 6 {
				additionalSubjectMarks, _ := strconv.Atoi(marks[5][0])
				percentage = (float64(total) - float64(additionalSubjectMarks)) / 5.00
			} else {
				percentage = float64(total) / 5.00
			}
			// Convert everything back to string
			strTotal := strconv.Itoa(total)
			strPercentage := strconv.FormatFloat(percentage, 'f', 2, 64)

			studentData := []string{roll, gender, name}
			for j := 0; j < 6+modifySubject; j++ {
				tempArray := []string{studentSubjects[j][0], studentSubjects[j][1], 
				marks[j][0], marks[j][1]}
				studentData = append(studentData, tempArray...)
			}
			tempArray := []string{strTotal, strPercentage, result}
			studentData = append(studentData, tempArray...)

			parsedData = append(parsedData, studentData)
		}
	}
	return parsedData, missingSubjectCodes, nil
}

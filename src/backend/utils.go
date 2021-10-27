package backend

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func read(path, class string) ([]string, error) {
	inputFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Incorrect path!")
	}
	defer inputFile.Close()

	// Accept only .txt file
	// Legacy Code
	if strings.ToLower(filepath.Ext(path)) != ".txt" {
		return nil, errors.New("Invalid Input File! (Accepts Only .txt file)")
	}

	// Convert the text file into []string
	scanner := bufio.NewScanner(inputFile)
	var rawData []string

	for scanner.Scan() {
		records := strings.Fields(scanner.Text())
		rawData = append(rawData, records...)
	}

	// Check if the given file is actually a result file or not
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
		return nil, errors.New("Mismatch between Input File and Class selected!")
	}

	return rawData, nil
}

func getSubjectCode() (map[string]string, error) {
	// If there is a User-Specific subject code file, then use that
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

// Helper function to check if a element is in a slice or not
func notInSlice(slice []string, ele string) bool {
	for _, i := range slice {
		if i == ele {
			return false
		}
	}
	return true
}

func parser(rawData []string, class string) (parsedData [][]string, schoolResult Result,
	missingSubjectCodes []string, err error) {
	// Set lowest percentage to be 100 initially
	schoolResult.LowestMarks.Percentage = 100.00

	// For Class X every field is 3 index before their counterpart in Class XII
	var modifyIndex int
	if class == "X" {
		modifyIndex = -3
	}

	// Checks if there are 8 digits or not
	rollRegex, _ := regexp.Compile("^\\d{8}$")
	// First group checks if there is a word comprising of A-Z
	// Second group checks if there is a standalone letter in case of initials
	nameRegex, _ := regexp.Compile("(([A-Z])\\w+)|([A-Z]){1}")

	subjectCodes, err := getSubjectCode()
	if err != nil {
		return nil, Result{}, nil, err
	}

	parsedData = [][]string{outputTemplate}

	for i, _ := range rawData {
		// Skip those elements which are not Roll Numbers
		if rollRegex.MatchString(rawData[i]) {
			// Parse individual student's data

			roll := rawData[i]
			gender := rawData[i+1]

			// Maybe there is a better way to do this
			// Needs Refactoring
			var name string
			switch {
			// Extremely rare but if a student has five word name
			case nameRegex.MatchString(rawData[i+6]):
				name = strings.Join(rawData[i+2:i+7], " ")
				// Adjust i because I have considered everything according to
				// default two word name.
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
					result = rawData[i+12+modifyIndex]
					// Adjust subject because I have made code for the majority 6 subject takers,
					// this will ensure in case of 5 subjects, program doesn't panic.
					modifySubject = -1
				} else {
					// In case of Fail or Other things, skip record.
					continue
				}
			}

			var (
				studentSubjects [][]string
				skipRecord      bool
			)
			for _, j := range rawData[i+4:i+10+modifySubject] {
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
			if skipRecord {
				continue
			}

			var marks [][]string
			// Step here is 2 because both marks and its corresponding grade is added to marks array.
			// First element = Marks
			// Second element = Grade
			for j := 14; j < 25+modifySubject; j += 2 {
				marks = append(marks,
					[]string{rawData[i+j+modifyIndex+modifySubject], rawData[i+j+modifyIndex+1]})
			}

			var total int
			for _, j := range marks {
				num, _ := strconv.Atoi(j[0])
				total += num
			}

			// Percentage is calculated by taking the main 5 subjects only
			// Additional subject is not calculated
			var percentage float64
			if len(studentSubjects) == 6 {
				additionalSubjectMarks, _ := strconv.Atoi(marks[5][0])
				percentage = (float64(total) - float64(additionalSubjectMarks)) / 5.00
			} else {
				percentage = float64(total) / 5.00
			}

			// Update School Result
			if percentage > schoolResult.Topper.Percentage {
				schoolResult.Topper.Roll = roll
				schoolResult.Topper.Name = name
				schoolResult.Topper.Percentage = percentage
			} else if percentage < schoolResult.LowestMarks.Percentage {
				schoolResult.LowestMarks.Roll = roll
				schoolResult.LowestMarks.Name = name
				schoolResult.LowestMarks.Percentage = percentage
			}

			// Convert everything back to string
			strTotal := strconv.Itoa(total)
			strPercentage := strconv.FormatFloat(percentage, 'f', 2, 64)

			// Compile everything
			studentData := []string{roll, gender, name}
			for j := 0; j < 6+modifySubject; j++ {
				tempArray := []string{studentSubjects[j][0], studentSubjects[j][1],
					marks[j][0], marks[j][1]}
				studentData = append(studentData, tempArray...)
			}
			tempArray := []string{strTotal, strPercentage, result}
			studentData = append(studentData, tempArray...)

			// Push it to main slice
			parsedData = append(parsedData, studentData)

		}
	}

	// All the Aggregate School Result is written at the last of the file
	tempArray := rawData[len(rawData) - 25: len(rawData) - 1]
	schoolResult.Absent = tempArray[20]
	schoolResult.Comptt = tempArray[11]
	schoolResult.Other = tempArray[23]
	schoolResult.Passed = tempArray[7]
	schoolResult.Repeat = tempArray[16]
	schoolResult.Total = tempArray[3]

	return parsedData, schoolResult, missingSubjectCodes, nil
}

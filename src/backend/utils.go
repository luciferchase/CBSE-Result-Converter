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
	// For Class X every field is 3 index before their counterpart in Class XII
	var modifyIndex int
	if class == "X" {
		modifyIndex = -3
	}

	rollRegex, _ := regexp.Compile("^\\d{8}$")
	nameRegex, _ := regexp.Compile("(([A-Z])\\w+)|([A-Z]){1}")

	subjectCodes, err := getSubjectCode()
	if err != nil {
		return nil, Result{}, nil, err
	}

	parsedData = [][]string{outputTemplate}

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
			if skipRecord {
				continue
			}

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

			if strings.Compare(strPercentage, schoolResult.Topper.Percentage) == 1 {
				schoolResult.Topper.Name = name
				schoolResult.Topper.Percentage = strPercentage
			} else if strings.Compare(strPercentage, schoolResult.LowestMarks.Percentage) == -1 {
				schoolResult.Topper.Name = name
				schoolResult.Topper.Percentage = strPercentage
			}

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

	tempArray := rawData[len(rawData) - 25: len(rawData) - 1]
	schoolResult.Absent = tempArray[20]
	schoolResult.Comptt = tempArray[11]
	schoolResult.Other = tempArray[23]
	schoolResult.Passed = tempArray[7]
	schoolResult.Repeat = tempArray[16]
	schoolResult.Total = tempArray[3]

	return parsedData, schoolResult, missingSubjectCodes, nil
}

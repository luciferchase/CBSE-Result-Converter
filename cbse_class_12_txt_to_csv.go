package main

import (
    "bufio"
    "encoding/csv"
    "errors"
    "fmt"
    "regexp"
    "strconv"
    "strings"
    "os"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Println("Enter the (absolute) path to the CBSE Class 12th result file",
         "(should be of the format {SCHOOL_CODE}.TXT).")
        fmt.Println("If this executable program is stored in the same directory, then enter its name only.")
        
        fmt.Print("\nPath: ")
        scanner.Scan()
        path := strings.ToLower(scanner.Text())

        err := writeToCSV(path)
        if err != nil {
            fmt.Println(err)
            fmt.Println("An error occurred. Please try again!\n")
            continue
        } else {
            fmt.Println("Records successfully written to 'class_12th_result.csv'!")
            break
        }
    }
}

func writeToCSV(path string) error {
    if path[len(path) - 4: len(path)] != ".txt" {
        return errors.New("It doesn't looks like a valid result file. Make sure the path is correct.")
    }

    inputFile, err := os.Open(path)
    if err != nil {
        fmt.Println(err)
        return errors.New("Incorrect path!")
    }
    defer inputFile.Close()

    scanner := bufio.NewScanner(inputFile)
    var rawData []string

    for scanner.Scan() {
        records := strings.Fields(scanner.Text())
        rawData = append(rawData, records...) 
    }

    // Check if the given file is actually 12th result file or not
    if strings.Join(rawData[4:8], " ") != "SENIOR SCHOOL CERTIFICATE EXAMINATION" {
        fmt.Println("It doesn't looks like a valid result file. Make sure the path is correct.")
        return errors.New("Not a valid 12th result file.")
    }

    parsedData := parser(rawData)

    outputFile, _ := os.Create("class_12th_result.csv")
    defer outputFile.Close()

    writer := csv.NewWriter(outputFile)
    
    err = writer.WriteAll(parsedData)
    if err != nil {
        fmt.Println(err)
        return errors.New("Error occurred while writing to .csv")
    } else {
        return nil
    }
}

func parser(rawData []string) [][]string {
    rollRegex, _ := regexp.Compile("^\\d{8}$")
    nameRegex, _ := regexp.Compile("([A-Z])\\w+")

    subjectCodes := map[string]string{
        "027": "HISTORY",
        "028": "POLITICAL SCIENCE",
        "029": "GEOGRAPHY",
        "030": "ECONOMICS",
        "041": "MATHEMATICS",
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

    parsedData := [][]string{
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
            gender := rawData[i + 1]

            var name string
            switch {
            // Extremely rare but if a student has five word name
            case nameRegex.MatchString(rawData[i + 6]):
                name = strings.Join(rawData[i + 2: i + 7], " ")
                i += 3

            // Rare but if a student has four word name
            case nameRegex.MatchString(rawData[i + 5]):
                name = strings.Join(rawData[i + 2: i + 6], " ")
                i += 2

            // If a student has three word name
            case nameRegex.MatchString(rawData[i + 4]):
                name = strings.Join(rawData[i + 2: i + 5], " ")
                i++

            // If a student has only one word name
            case !nameRegex.MatchString(rawData[i + 3]):
                name = rawData[i + 2]
                i--
            
            // Default two word names
            default:
                name = strings.Join(rawData[i + 2: i + 4], " ")
            }

            // Skip R.L or other results which doesn't have any data to it
            result := rawData[i + 13]
            if result != "PASS" && result != "COMP" {
                continue
            }

            var studentSubjects [][]string
            for _, j := range rawData[i + 4: i + 10] {
                studentSubjects = append(studentSubjects, []string{j, subjectCodes[j]})
            }

            var marks [][]string
            for j := 14; j < 25; j += 2 {
                marks = append(marks, []string{rawData[i + j], rawData[i + j + 1]})
            }

            var total int
            for _, j := range marks {
                num, _ := strconv.Atoi(j[0])
                total += num
            }
            additionalSubjectMarks, _ := strconv.Atoi(marks[5][0])
            percentage := (total - additionalSubjectMarks) / 5

            // Convert everything back to string
            strTotal := strconv.Itoa(total)
            strPercentage := strconv.Itoa(percentage)

            studentData := []string{roll, gender, name}
            for j := 0; j < 6; j++ {
                tempArray := []string{studentSubjects[j][0], studentSubjects[j][1], marks[j][0], marks[j][1]}
                studentData = append(studentData, tempArray...)
            }
            tempArray := []string{strTotal, strPercentage, result}
            studentData = append(studentData, tempArray...)

            parsedData = append(parsedData, studentData)
        }
    }
    return parsedData
}
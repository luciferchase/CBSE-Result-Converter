package main

import (
    "bufio"
    "encoding/csv"
    "errors"
    "fmt"
    "os"
    "regexp"
    "strconv"
    "strings"

    "github.com/fatih/color"

    // Useful packages but buggy hence not in use yet
    // "github.com/common-nighthawk/go-figure"
    // "github.com/manifoldco/promptui"
)

func main() {
    // Set colours
    asciiArt := color.New(color.FgCyan, color.Bold).PrintlnFunc()
    prompt := color.New(color.FgYellow, color.Bold).PrintFunc()
    warn := color.New(color.FgRed, color.Bold).PrintlnFunc()
    success := color.New(color.FgGreen, color.Bold).PrintfFunc()

    title := 
    `
        ========================================================================================
          ________  ________  ___               ____    _____                      __         
         / ___/ _ )/ __/ __/ / _ \___ ___ __ __/ / /_  / ___/__  ___ _  _____ ____/ /____ ____
        / /__/ _  |\ \/ _/  / , _/ -_|_-</ // / / __/ / /__/ _ \/ _ \ |/ / -_) __/ __/ -_) __/
        \___/____/___/___/ /_/|_|\__/___/\_,_/_/\__/  \___/\___/_//_/___/\__/_/  \__/\__/_/  

        ========================================================================================                                                                                         
    `
    asciiArt(title)

    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Println("Whose class' result do you want to convert?")
        fmt.Println("\t1. Class X")
        fmt.Println("\t2. Class XII")
        
        prompt("\nEnter your choice: ")
        scanner.Scan()
        choice := scanner.Text()
        
        var class string
        var outputFileName string
        switch {
        case choice == "1":
            class = "X"
            outputFileName = "class_10th_result.csv"
        case choice == "2":
            class = "XII"
            outputFileName = "class_12th_result.csv"
        default:
            warn("Invalid choice. Please try again!")
            continue
        }

        fmt.Println(
            "\nEnter the (absolute) path to the CBSE Class", class, 
            "result file (should be of the format {SCHOOL_CODE}.TXT).",
        )
        fmt.Println("If this executable program is stored in the same directory, then enter its name only.")

        prompt("\nPath: ")
        scanner.Scan()

        path := strings.ToLower(scanner.Text())

        err := writeToCSV(class, path)
        if err != nil {
            fmt.Println(err)
            warn("An error occurred. Please try again!\n")
            continue
        } else {
            success("Records successfully written to '%s'!\n", outputFileName)

            prompt("Press [Y] to convert another file or press [N] to exit the program: ")
            scanner.Scan()

            if strings.ToLower(scanner.Text()) == "y" {
                fmt.Println()
                continue
            } else {
                break
            }
        }
    }
}

func writeToCSV(class, path string) error {
    if len(path) == 0 || path[len(path)-4:len(path)] != ".txt" {
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

    // Check if the given file is actually result file or not
    var validFile int
    var outputFileName string
    switch {
    case class == "X":
        validFile = strings.Compare(strings.Join(rawData[4:7], " "), "SECONDARY SCHOOL EXAMINATION")
        outputFileName = "class_10th_result.csv"
    case class == "XII":
        validFile = strings.Compare(strings.Join(rawData[4:8], " "), "SENIOR SCHOOL CERTIFICATE EXAMINATION")
        outputFileName = "class_12th_result.csv"
    }

    if validFile != 0 {
        return errors.New("Invalid file!")
    }

    parsedData := parser(class, rawData)

    outputFile, _ := os.Create(outputFileName)
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

func parser(class string, rawData []string) [][]string {
    // For Class X every field is 3 index before their counterpart in Class XII
    var modifyIndex int
    if class == "X" {
        modifyIndex = -3
    }

    rollRegex, _ := regexp.Compile("^\\d{8}$")
    nameRegex, _ := regexp.Compile("(([A-Z])\\w+)|([A-Z]){1}")

    subjectCodes := map[string]string{
        // Class X
        "002": "HINDI",
        "086": "SCIENCE",
        "087": "SOCIAL SC.",
        "122": "SANSKRIT",
        "184": "ENGLISH",
        "241": "MATHEMATICS - BASIC",
        "401": "RETAILING",
        "402": "IT",

        // COMMON
        "041": "MATHEMATICS",

        // Class XII
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
                    modifySubject = -1
                    i++
                } else {
                    continue
                }
            }

            var studentSubjects [][]string
            for _, j := range rawData[i+4 : i+10+modifySubject] {
                studentSubjects = append(studentSubjects, []string{j, subjectCodes[j]})
            }

            var marks [][]string
            for j := 14; j < 25+modifySubject; j += 2 {
                marks = append(marks, []string{rawData[i+j+modifyIndex], rawData[i+j+modifyIndex+1]})
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
            strPercentage := fmt.Sprint(percentage)

            studentData := []string{roll, gender, name}
            for j := 0; j < 6+modifySubject; j++ {
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

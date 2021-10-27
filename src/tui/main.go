package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/sqweek/dialog"

	backend "github.com/luciferchase/CBSE-Result-Converter/src/backend"
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
		switch {
		case choice == "1":
			class = "X"
		case choice == "2":
			class = "XII"
		default:
			warn("Invalid choice. Please try again!")
			continue
		}

		prompt("Press any key to select file: ")
		scanner.Scan()

		path, err := dialog.File().Filter("Text File (*.txt, *.text)", "txt").Load()
		if err != nil {
			fmt.Println(err)
			warn("An error occurred. Please try again!\n")
			continue
		}

		parsedData, schoolResult, missingSubjectCodes, err := backend.Parse(path, class)
		if err != nil {
			fmt.Println(err)
			warn("An error occurred. Please try again!\n")
			continue
		}

		// If there some subject codes are missing, prompt the user to manually enter their names
		// and store it as a .json file
		if len(missingSubjectCodes) != 0 {
			fmt.Println("These Subject Codes are missing in the database:")
			fmt.Println(missingSubjectCodes)
			fmt.Println("Do you want to:",
				"\n1. Enter their names so next time this won't happen?",
				"\n2. Skip these records?")

			prompt("\nEnter your choice: ")
			scanner.Scan()
			choice = scanner.Text()

			if choice == "1" {
				subjectCodes := make(map[string]string)
				fmt.Println("Enter the names of the following codes:")

				for _, i := range missingSubjectCodes {
					prompt(i + ": ")
					scanner.Scan()
					subjectCodes[i] = strings.ToUpper(scanner.Text())
				}

				err = backend.Update(subjectCodes)
				if err != nil {
					fmt.Println(err)
					warn("An error occurred. Please try again!\n")
				}

				// Try again
				parsedData, schoolResult, missingSubjectCodes, err = backend.Parse(path, class)
				if err != nil {
					fmt.Println(err)
					warn("An error occurred. Please try again!\n")
					continue
				}
			} else {
				warn("Skipping the missing subject codes records")
			}
		}

		success("\nFile Converted Successfully!")

		fmt.Println("\nTOTAL CANDIDATES:\t", schoolResult.Total)
		fmt.Println("TOTAL PASSED:\t\t", schoolResult.Passed)
		fmt.Println("TOTAL ABSENT:\t\t", schoolResult.Absent)
		fmt.Println("TOTAL COMPTT:\t\t", schoolResult.Comptt)
		fmt.Println("TOTAL ESSENTIAL REPEAT:\t", schoolResult.Repeat)
		fmt.Println("OTHER:\t\t\t", schoolResult.Other)

		prompt("\nPress any key to save file: ")
		scanner.Scan()

		path, err = dialog.File().Filter("XLS Worksheet (*.csv)", "csv").Title("Save File").Save()
		if err != nil {
			fmt.Println(err)
			warn("An error occurred. Please try again!\n")
			continue
		} else if len(path) == 0 {
			// If no file is specifed, store it to ./output/ folder
			warn("Saving file to ./output/ folder")
		}

		err = backend.Write(parsedData, class, path)
		if err != nil {
			fmt.Println(err)
			warn("An error occurred. Please try again!\n")
			continue
		} else {
			success("Records successfully written to '%s'!\n", filepath.Base(path))
		}

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

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	backend "github.com/luciferchase/CBSE-Result-Converter/src/backend"

	"github.com/fatih/color"
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

		parsedData, missingSubjectCodes, err := backend.Parse(path, class)
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
				parsedData, missingSubjectCodes, err = backend.Parse(path, class)
				if err != nil {
					fmt.Println(err)
					warn("An error occurred. Please try again!\n")
					continue
				}
			} else {
				warn("Skipping the missing subject codes records")
			}
		}

		err = backend.Write(parsedData, class)
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

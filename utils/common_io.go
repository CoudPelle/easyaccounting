package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/exp/slices"
)

// Read a string input and parse it according to OS environment
func StrInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	// Remove CRLF and LF
	if runtime.GOOS == "windows" {
		text = strings.TrimRight(text, "\r\n")
	} else {
		text = strings.TrimRight(text, "\n")
	}
	return text
}

// Check if a user input is amoung a list of attended string, ask again for a prompt if not
// Returns choice
func promptChoice(choiceList []string) string {
	var choice string
	choice = StrInput()
	if !slices.Contains(choiceList, choice) {
		fmt.Println("Please make a choice between: ", choiceList)
		choice = StrInput()
	}
	return choice
}

// Ask for the cv filename, clean the string and returns it
func promptCsvName() string {

	fmt.Println("Enter the name of your SG account exported csv")
	fmt.Println("---------------------")

	return StrInput()
}

// Print all transaction types
func PromptTransactionTypes(types []string) {
	typesString := "- Afficher les descriptions de types(?)"
	for index := range types {
		if index%4 == 0 {
			typesString += "\n- "
		}
		typesString += fmt.Sprintf("%s(%d) ", types[index], index)
	}
	typesString = typesString[:len(typesString)-1]
	fmt.Println(typesString)
}

// Print descriptions for all transaction types
func PromptTransationTypesDescription(typesDesc []string) {
	typesDescString := "> Descriptions des types\n"
	for index := range typesDesc {
		typesDescString += fmt.Sprintf("%d - %s\n", index, typesDesc[index])
	}
	fmt.Println(typesDescString)
}

// check if a tmp file exists, returns true if yes, false if not
func checkForTmpSave(filePath string) bool {
	var IsTmpSave bool
	IsTmpSave = false
	// If filename Contains _tmp.csv, we check if original file can be found
	if strings.Contains(filePath, "_tmp.csv") {
		_, err := os.Stat(strings.Replace(filePath, "_tmp.csv", ".csv", 1))
		// test if error is PathError type
		if _, isFileMissing := err.(*os.PathError); !isFileMissing {
			IsTmpSave = true
		}
		// If a tmp save exists, alert user
		if !IsTmpSave {
			fmt.Println("Found a temporary save file without source data file.\nFile: '", filePath, "'")
		}
	}
	return IsTmpSave
}

// check list of files in input folder. Their must be only one returned
func checkInputFiles(files []string, localDir string) (string, string) {
	if len(files) > 0 {
		// fmt.Println("Found :")
		// for _, file := range files {
		// 	// print local path of input file
		// 	fmt.Println(strings.Replace(file, localDir, ".", 1))
		// }

		if len(files) != 1 {
			log.Fatal("Too many inputs, please provide only one csv file.")
		}
	} else {
		log.Fatal("Couldn't find any input")
	}
	return files[0], strings.Replace(files[0], localDir, ".", 1)
}

// Fetch a .csv in ./input folder
// Returns an error if it can't find any file, or too much files
// Else, return filepath and filename
func findBankCsv() (string, string) {
	var files []string
	var inputFolder string
	// Get local dir
	localDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// Seek files in input
	if runtime.GOOS == "windows" {
		inputFolder = "\\input"
	} else {
		inputFolder = "/input"
	}
	localDir = localDir + inputFolder
	err = filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		// Todo: detect when its a problem with folder not existing before the output below
		if err != nil {
			localFolderError := `Couldn't find a local input folder\n
				Please create an input and ouput folder in the same folder of easyaccounting binary.`
			fmt.Println(localFolderError)
			log.Fatal(err)
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			if checkForTmpSave(path) {
				fmt.Println("> Found file:", path)
				fmt.Println("> This is a temporary save file, would you like to open it and continue your work ?(y/n)")
				choice := promptChoice([]string{"y", "n"})
				if choice == "y" {
					// add only path in the files slice and end os.walk
					files = []string{path}
					return nil
				}
			}
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return checkInputFiles(files, localDir)
}

// Open csv and returns it as a 2d array of string
func readCSV(csvPath string) [][]string {
	file, err := os.Open(csvPath)
	if err != nil {
		panic(err)
	}
	row1, err := bufio.NewReader(file).ReadSlice('\n')
	if err != nil {
		panic(err)
	}
	_, err = file.Seek(int64(len(row1)), io.SeekStart)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true
	rows, err := reader.ReadAll()
	defer file.Close()
	if err != nil {
		panic(err)
	}
	return rows
}

// Custom NewWriter to set ";" as separator
func customWriter(w io.Writer) (writer *csv.Writer) {
	writer = csv.NewWriter(w)
	writer.Comma = ';'
	return
}
func WriteCSV(rows [][]string, csvName string) {
	fmt.Println("> Sauvegarde de votre fichier dans le dossier \"output\".")
	f, err := os.Create(csvName)
	if err != nil {
		log.Fatal(err)
	}
	err = customWriter(f).WriteAll(rows)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// Find csv, open it as a double string array
// return rows and csv path
func GetCSV() ([][]string, string) {
	csvPath, _ := findBankCsv()

	fmt.Printf("> Opening \"%s\"\n\n", csvPath)
	rows := readCSV(csvPath)
	return rows, csvPath
}

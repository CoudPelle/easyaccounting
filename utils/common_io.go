package utils

import (
	"encoding/csv"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
// Ask for the cv filename, clean the string and returns it
func promptCsvName() string {

	fmt.Println("Enter the name of your SG account exported csv")
	fmt.Println("---------------------")

	return StrInput()
}
// Print all transaction types 
func PromptTransactionTypes(types []string){
	typesString := "- Afficher les descriptions de types(?)"  
	for index, _ := range types {
		if index%4 == 0 {
			typesString += "\n- "
		}
		typesString += fmt.Sprintf("%s(%d) ", types[index], index)
	}
	typesString = typesString[:len(typesString)-1]
	fmt.Println(typesString)
}
// Print descriptions for all transaction types 
func PromptTransationTypesDescription(typesDesc []string){
	typesDescString := "> Descriptions des types\n"  
	for index, _ := range typesDesc {
		typesDescString += fmt.Sprintf("%d - %s\n", index, typesDesc[index])
	}
	typesDescString = typesDescString
	fmt.Println(typesDescString)
}
// Fetch a .csv in ./input folder
// Returns an error if it can't find any file, or too much files
// Else, return filepath and filename
func findBankCsv() (string, string){
	var files []string
	var inputFolder string
	// Get local dir
	local, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Seek files in input
	if runtime.GOOS == "windows" {
		inputFolder = "\\input"
	} else {
		inputFolder = "/input"
	}
	local = local + inputFolder
	err = filepath.Walk(local, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(files) > 0 {
		// fmt.Println("Found :")
		// for _, file := range files {
		// 	// print local path of input file
		// 	fmt.Println(strings.Replace(file, local, ".", 1))
		// }
		if len(files) != 1 {
			log.Fatal("Too many inputs, please provide only one csv file.")
		}
	} else {
		log.Fatal("Couldn't find any input")
	}
	return files[0], strings.Replace(files[0], local, ".", 1)
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
func GetCSV() ([][]string, string){
	csvPath, _ := findBankCsv()

	fmt.Printf("> Opening \"%s\"\n\n",csvPath)
	rows := readCSV(csvPath)
	return rows, csvPath
}


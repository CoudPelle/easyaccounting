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


// Ask for the cv filename, clean the string and returns it
func promptCsvName() string {

	fmt.Println("Enter the name of your SG account exported csv")
	fmt.Println("---------------------")

	return StrInput()
}

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

// Fetch a .csv in ./input folder
// Returns an error if it can't find any file, or too much files
// Else, return filepath
func findBankCsv() string {
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
	fmt.Printf(local)
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
		fmt.Println("Found :")
		for _, file := range files {
			// print local path of input file
			fmt.Println(strings.Replace(file, local, ".", 1))
		}
		if len(files) != 1 {
			log.Fatal("Too many inputs, please provide only one csv file.")
		}
	} else {
		log.Fatal("Couldn't find any input")
	}
	return files[0]
}

// Open csv and returns it as a 2d array of string
func readCSV(csv_name string) [][]string {
	file, err := os.Open(csv_name)
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

func GetCSV()[][]string {
	csv_path := findBankCsv()
	
	fmt.Println("opening " + csv_path)
	rows := readCSV(csv_path)
	// Remove column names
	values := rows[1:]
	return values
}
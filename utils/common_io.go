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

var (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// Desc: Read a string input with bufio lib and parse it according to OS environment
// Return: user input without new line character
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

// Desc: Check if a user input is amoung a list of expected string
// Request an input again if it is not the desired input
// Parameters: slice of expected input
// Return: choice
func getChoice(choiceList []string) string {
	var choice string
	choice = StrInput()
	if !slices.Contains(choiceList, choice) {
		fmt.Println("Choisissez parmis les options: ", choiceList)
		choice = StrInput()
	}
	return choice
}

// Desc: Ask for the cv filename, clean the string and returns it
// Return: csv filename given by the user
// TODO: check if csv filename exists
func promptCsvName() string {

	fmt.Println("Entrez le nom du csv de votre relevé de compte SG")
	fmt.Println("---------------------")

	return StrInput()
}

// Desc: Print all transaction types from given slice
// Parameters: transaction types as slice
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

// Desc: Print descriptions for all transaction types
// Parameters: types description as a slice of string
func PromptTransationTypesDescription(typesDesc []string) {
	typesDescString := "> Descriptions des types\n"
	for index := range typesDesc {
		typesDescString += fmt.Sprintf("%d - %s\n", index, typesDesc[index])
	}
	fmt.Println(typesDescString)
}

// Desc: Check if a tmp file exists, returns true if yes, false if not
// Parameters: filepath of file to test
// Return: a bool indication whether tmp file and original exists
func checkForTmpSave(filePath string) bool {
	var IsTmpSave bool
	IsTmpSave = false
	// If filename contains checkpoint file .tmp , we check if original file can be found
	if strings.Contains(filePath, ".tmp") {
		source_filename := strings.Replace(filePath, ".tmp", ".csv", 1)
		_, err := os.Stat(source_filename)
		// test if error is PathError type
		if _, isFileMissing := err.(*os.PathError); !isFileMissing {
			IsTmpSave = true
		}
		// If a tmp save exists, without source csv, alert user
		if !IsTmpSave {
			fmt.Println("> " + ColorRed + "Le fichier de la sauvegarde d'un travail en cours a été trouvé sans son fichier source." + ColorReset + "\nFichier manquant: '" + source_filename + "'")
		}
	}
	return IsTmpSave
}

// Desc: Check list of files in input folder. Their must be only one returned
// Parameters: a slice of filenames, the local directory path
// Return: input file path
func checkInputFiles(files []string, localDir string) string {
	if len(files) > 0 {
		if len(files) != 1 {
			log.Fatal("Trop de fichiers trouvés dans le dossier Input. Merci de fournir uniquement 1 fichier (excepté les .tmp)")
		}
	} else {
		log.Fatal("Aucun fichier trouvé dans le dossier Input")
	}
	return files[0]
}

// Desc: Fetch a .csv in ./input folder
// Returns an error if it can't find any file, or too much files
// Else, return filepath and filename
// Return: path of input file, load found tmp boolean
func findBankCsv() (string, bool) {
	var files []string
	var inputFolder string
	var loadTmp = false
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
			localFolderError := `Dossier local "Input" introuvable\n
				Merci de créer 2 dossiers "Input" et "Output" dans le dossier de l'application easyaccounting.`
			fmt.Println(localFolderError)
			log.Fatal(err)
			return nil
		}

		if !info.IsDir() && (filepath.Ext(path) == ".csv" || filepath.Ext(path) == ".tmp") {
			if checkForTmpSave(path) {
				fmt.Println("> Fichier temporaire trouvé:", path)
				fmt.Println("> " + ColorYellow + "Voulez-vous reprendre votre travail ?(y/n)" + ColorReset)
				choice := getChoice([]string{"y", "n"})
				if choice == "y" {
					// add only path in the files slice and end os.walk
					files = []string{strings.Replace(path, ".tmp", ".csv", 1)}
					loadTmp = true
					return nil
				}
			} else {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	filePath := checkInputFiles(files, localDir)
	return filePath, loadTmp
}

// Desc: Open csv and returns it as a 2d array of string
// Parameters: Path of the input CSV file
// Return: csv file as 2d array
func ReadCSV(csvPath string) [][]string {
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatal("Unable to read input file "+csvPath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+csvPath, err)
	}
	return rows
}

// Desc: Custom NewWriter to set ";" as separator
// Parameters: file output
// Return: csv writer
func customWriter(file io.Writer) (writer *csv.Writer) {
	writer = csv.NewWriter(file)
	writer.Comma = ';'
	return
}

// Desc: write a 2d array in a new csv file in a local output directory
// Parameters: csv file as 2d array, csvFilePath for the output file path
func WriteCSV(rows [][]string, csvFilePath string) {
	f, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	err = customWriter(f).WriteAll(rows)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// Desc: cycle map keys of given csv map and write a csv per key
// in a new folder with original filename
// Parameters: csv files as map of 2d array, inputFilePath
func saveCSVMap(csvMap map[string][][]string, inputFilePath string) {

	outputPath := strings.Replace(inputFilePath, "input", "output", 1)
	fileName := outputPath[strings.LastIndex(outputPath, "/")+1 : strings.Index(outputPath, ".csv")]
	outputPath = outputPath[:strings.LastIndex(outputPath, "/")+1] +
		fileName + "/"
	os.MkdirAll(outputPath, os.ModePerm)
	for transaction_type := range csvMap {
		outputFile := transaction_type + ".csv"
		WriteCSV(csvMap[transaction_type], outputPath+outputFile)
	}
}

// Desc: Global saving function, I may add other stuff to save later
// Parameters: csv files as map of 2d array, inputFilePath
func SaveResults(csvMap map[string][][]string, inputFilePath string) {
	saveCSVMap(csvMap, inputFilePath)
}

// Desc: delete a file at a given filePath
// Parameters: filePath of file to delete
func DeleteFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		log.Fatal(err)
	}
}

// Desc: Find csv, open it as a double string array
// Return: csv file as 2d array, csv path, loadTmp as boolean
func GetCSV() ([][]string, string, bool) {
	csvPath, loadTmp := findBankCsv()

	fmt.Printf("> Ouverture de \"%s\"\n", csvPath)
	rows := ReadCSV(csvPath)
	return rows, csvPath, loadTmp
}

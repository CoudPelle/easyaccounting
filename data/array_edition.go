package data

import (
	"easyaccounting/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Desired columns for final csv
var (
	COL_NAMES = []string{
		"Date transaction", "Date prelevement", "Label", "Montant", "Type"}
	TRANSACTION_TYPES = []string{
		"Cadeau", "Don", "Divers", "Vacances", "Administratif", "Logement",
		"Vetement", "Media", "Abonnement", "Transport", "Sortie", "Nourriture", "Loisir", "Sante"}
	TRANSACTION_TYPES_DESC = []string{
		"Cadeaux à des proches",
		"Dons à des associations, projets...",
		"Non catégorisable, dépense ponctuelle",
		"Logements de vacances, activités",
		"Impôts, fabrication papiers...",
		"Loyer, factures électricité/gaz/box internet, lavomatic, matériel d'entretien, éléctroménager, meubles",
		"Vêtements, chaussures, ce qui se porte",
		"Abonnements à des medias, numérique comme papier",
		"Autres abonnements : netflix, bitwarden, forfait mobile, spotify",
		"Abonnements transports, billet de trains, carburant",
		"Activités entre amis/famille, restaurant/fast-food accompagné ou non",
		"Marché, épicerie, boulangerie, fromager",
		"Toutes les activités ludiques non liés au travail : achats de jeux, livres, concert, cinéma, sport, musique, DIY",
		"Consultation médicales, médicaments, passage aux urgences"}
)

// Desc: Remove end of card number (present as a substring) in column colIndex, if it exists
// Parameters: row to be process, index of column where card number is
// Return: a new row with column cleaned
func removeCardNum(row []string, colIndex int) []string {
	var new_row []string
	new_row = row
	if strings.HasPrefix(row[colIndex], "CARTE X") {
		new_row[colIndex] = row[colIndex][12:]
	}
	return new_row
}

// Desc: Extract the transaction date contained in column labelColIndex for given row
// Add it as a new column at beginning of the row, Add NULL if no date is found
// Parameters: row to process
// Return: a new row with string column transaction date
func addTransactionDateCol(row []string, labelColIndex int) []string {
	var transactionDate string
	var new_row []string

	new_row = row

	_, err := time.Parse("02/01", new_row[labelColIndex][:5])
	if err != nil {
		new_row = append([]string{"NULL"}, new_row...)
	} else {
		transactionDate = new_row[labelColIndex][:5]
		//remove transaction date from label col
		new_row[labelColIndex] = new_row[labelColIndex][6:]
		new_row = append([]string{transactionDate}, new_row...)
	}
	return new_row
}

// Desc: Add a transaction type column
// convert user type choice int as string
// Parameters: row to process
// Return: row with new string column type
func addTypeColumn(row []string) []string {

	var choice int
	choice = getTypeColumn(row)
	row = append(row, TRANSACTION_TYPES[choice])
	return row
}

// Desc: Display the row and prompt the user to choose a type for the transaction
// Parameters: row to display
// Return: the column type as integer
func getTypeColumn(row []string) int {
	var choice int

	fmt.Println("<---------->")
	fmt.Printf(" Colonnes    %+q\n", COL_NAMES[:len(COL_NAMES)-1])
	fmt.Printf(" Transaction %+q\n", row)
	fmt.Print("<---------->\n\n")

	utils.PromptTransactionTypes(TRANSACTION_TYPES)
	input := utils.StrInput()

	if input == "?" {
		utils.PromptTransationTypesDescription(TRANSACTION_TYPES_DESC)
		choice = getTypeColumn(row)
	} else {
		var err error
		choice, err = strconv.Atoi(input)
		if err != nil || choice < 0 || choice > len(TRANSACTION_TYPES) {
			fmt.Print(utils.ColorRed + "\nErreur: Merci d'entrer une des valeurs proposées." + utils.ColorReset + "\n\n")
			choice = getTypeColumn(row)
		}
	}
	return choice
}

// Desc: Remove unwanted columns for a given row and column index
// Parameters: row to process, colindex to remove
// Return: a new row without this column
func cleanColumns(row []string, colIndex int) []string {
	var new_row []string
	new_row = row
	//create new row without unwanted column
	copy(new_row[colIndex:], new_row[colIndex+1:])

	return new_row[:len(new_row)-1]
}

// Desc: Save a checkpoint of work in progress, so we can work on a file, stop the program and continue later
// Parameters: csv file as 2d array, csvPath is the full path of csv file
func saveCheckpoint(values [][]string, csvPath string) {
	var tmp_filename string

	tmp_filename = csvPath
	tmp_filename = strings.Replace(csvPath, ".csv", ".tmp", 1)

	utils.WriteCSV(values, tmp_filename)
}

// Desc: Delete the checkpoint of work in progress
// Parameters: csvPath is the full path of csv file
func deleteCheckpoint(csvPath string) {
	var tmp_filePath string

	tmp_filePath = csvPath
	tmp_filePath = strings.Replace(csvPath, ".csv", ".tmp", 1)

	utils.DeleteFile(tmp_filePath)
}

// Desc: Cycle through array and process row 1-by-1
// May change the shape of the input array as it add/remove columns
// Parameters: csv file as 2d array, csvPath is the full path of csv file, load found tmp as boolean
// Return: a new array processed
func editColumns(values [][]string, csvPath string, loadTmp bool) [][]string {
	var values_cleaned [][]string
	if loadTmp == true {
		values_cleaned = utils.ReadCSV(strings.Replace(csvPath, ".csv", ".tmp", 1))
	}
	fmt.Print("> Choisissez à quel catégorie appartient chaque transaction:\n\n")
	for index, row := range values {
		// skip processed rows
		if index < len(values_cleaned) {
			continue
		}

		new_row := row
		//remove 2 last columns : currency, short label and empty column
		new_row = cleanColumns(new_row, 4)
		new_row = cleanColumns(new_row, 4)
		new_row = cleanColumns(new_row, 1)
		//remove card num if it exists in label col
		new_row = removeCardNum(new_row, 1)
		// add a transaction date column as the first column
		new_row = addTransactionDateCol(new_row, 1)
		// prompt to add type column
		new_row = addTypeColumn(new_row)
		values_cleaned = append(values_cleaned, new_row)
		saveCheckpoint(values_cleaned, csvPath)
	}
	return values_cleaned
}

// Desc: Main function for the editing of the accounting exported csv
// It removes original array column names, process the array and return the new array with new column names
// Parameters: csv file as 2d array, csvPath is the full path of csv file, load found tmp as boolean
// Output: csv processed as a new array
func FormatAccountingCSV(values [][]string, csvPath string, loadTmp bool) [][]string {
	// Remove column names
	values = values[1:]
	values = editColumns(values, csvPath, loadTmp)
	// Delete checkpoint file of work in progress
	deleteCheckpoint(csvPath)
	// add column names
	values = append([][]string{COL_NAMES}, values...)
	return values
}

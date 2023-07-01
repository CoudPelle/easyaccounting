package data

import (
	"easyaccounting/utils"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Desired columns for final csv
var (
	COL_NAMES = []string{
		"Date transaction", "Label", "Date prelevement", "Montant", "Type", "Categorie"}
	TRANSACTION_TYPES = []string{
		"Depense", "Versement"}
	TRANSACTION_CATEGORIES = []string{
		"Media",
		"Nourriture",
		"Sante",
		"Sortie",
		"Logement",
		"Loisir",
		"Transport",
		"Administratif",
		"Abonnement",
		"Divers",
		"Travail",
		"Vetement",
		"Don",
		"Cadeau",
		"Vacances",
	}
	TRANSACTION_CATEGORIES_DESC = []string{
		"Abonnements à des medias, numerique comme papier",
		"Marche, epicerie, boulangerie, fromager",
		"Consultation medicales, medicaments, passage aux urgences",
		"Activites entre amis/famille, restaurant/fast-food accompagne ou non",
		"Loyer, factures electricite/gaz/box internet, lavomatic, materiel d'entretien, electromenager, meubles",
		"Toutes les activites ludiques non lies au travail : achats de jeux, livres, concert, cinema, sport, musique, DIY",
		"Abonnements transports, billet de trains, carburant",
		"Impôts, fabrication papiers...",
		"Autres abonnements : netflix, bitwarden, forfait mobile, spotify",
		"Non categorisable, depense ponctuelle",
		"Depenses/Versements lies au travail (materiel, deplacements, salaires...)",
		"Vêtements, chaussures, ce qui se porte",
		"Dons à des associations, projets...",
		"Cadeaux à des proches",
		"Logements de vacances, activites",
	}
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

// Desc: Extract the transaction date contained in column labelIndex for given row
// Add it as a new column at beginning of the row, Add NULL if no date is found
// Parameters: row to process
// Return: a new row with string column transaction date
// TODO: replace hard coded substring to a date regex
func addTransactionDateCol(row []string, labelIndex int) []string {
	var transactionDate string
	var new_row []string

	new_row = row

	_, err := time.Parse("02/01", new_row[labelIndex][:5])
	if err != nil {
		new_row = append([]string{"NULL"}, new_row...)
	} else {
		transactionDate = new_row[labelIndex][:5]
		//remove transaction date from label col
		new_row[labelIndex] = new_row[labelIndex][6:]
		new_row = append([]string{transactionDate}, new_row...)
	}
	return new_row
}

// Desc: Add a transaction type column (deposit/versement or spent/depense)
// Parameters: row to process, amount column index
// Returns: new row with string column transaction type
func addTypeColumn(row []string, amountIndex int) []string {
	var transactionType string
	// Replace comma by dot is required to be parsed as float
	amountStr := strings.Replace(row[amountIndex], ",", ".", 1)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		log.Fatal("ERREUR FATALE: Impossible de recuperer le montant d'une de vos transactions\nTransaction: ", row, "\nErreur: ", err)
	}

	if amount < 0 {
		transactionType = TRANSACTION_TYPES[0]
	} else {
		transactionType = TRANSACTION_TYPES[1]
	}
	row = append(row, transactionType)
	return row
}

// Desc: Add a transaction category column
// convert user category choice int as string
// Parameters: row to process
// Return: row with new string column category
func addCategoryColumn(row []string) []string {

	var choice int
	choice = getCategoryColumn(row)
	row = append(row, TRANSACTION_CATEGORIES[choice])
	return row
}

// Desc: Display the row and prompt the user to choose a category for the transaction
// Parameters: row to display
// Return: the column category as integer
func getCategoryColumn(row []string) int {
	var choice int

	fmt.Println("<---------->")
	fmt.Printf(" Colonnes    %+q\n", COL_NAMES[:len(COL_NAMES)-1])
	fmt.Printf(" Transaction %+q\n", row)
	fmt.Printf("<---------->\n\n")

	utils.PromptTransactionTypes(TRANSACTION_CATEGORIES)
	input := utils.StrInput()

	if input == "?" {
		utils.PromptTransationTypesDescription(TRANSACTION_CATEGORIES_DESC)
		choice = getCategoryColumn(row)
	} else {
		var err error
		choice, err = strconv.Atoi(input)
		if err != nil || choice < 0 || choice > len(TRANSACTION_CATEGORIES) {
			fmt.Print(utils.ColorRed + "\nErreur: Merci d'entrer une des valeurs proposees." + utils.ColorReset + "\n\n")
			choice = getCategoryColumn(row)
		}
	}
	return choice
}

// Desc: Remove unwanted column for a given row and column index
// Parameters: row to process, colindex to remove
// Return: a new row without this column
func removeColumn(row []string, colIndex int) []string {
	new_row := make([]string, len(row))
	copy(new_row, row)
	//create new row without unwanted column
	copy(new_row[colIndex:], new_row[colIndex+1:])

	return new_row[:len(new_row)-1]
}

// Desc: Insert a column for a given row and column index
// Parameters: row to process, value to append, column index destination
// Return: a new row with new column
func insertColumn(row []string, value string, colIndex int) []string {
	return append(row[:colIndex], append([]string{value}, row[colIndex:]...)...)
}

// Desc: Move a column attribute for a given row from a source col to a destination col
// The current element at destination index will move to the left or destIndex-1
// Parameters: row to process, colindex to remove
// Return: a new row without this column
func moveColumn(row []string, srcIndex int, dstIndex int) []string {
	value := row[srcIndex]
	return insertColumn(removeColumn(row, srcIndex), value, dstIndex)
}

// Desc: Save a checkpoint of work in progress, so we can work on a file, stop the program and continue later
// Parameters: csv file as 2d array, csvPath is the full path of csv file
func saveCheckpoint(values [][]string, csvPath string) {
	var tmp_filename string

	tmp_filename = csvPath
	tmp_filename = strings.Replace(csvPath, ".csv", ".tmp", 1)

	utils.WriteCSV(values, tmp_filename)
}

// Desc: Create two 2d arrays by transaction type
// Parameters: csv file as 2d array, type column index
// Return: multiple csv files as 2d arrays
func discriminateByType(values [][]string, typeIndex int) map[string][][]string {
	values_discriminated := make(map[string][][]string)
	// initialise map
	for _, row := range values {
		// remove transaction type col
		new_row := removeColumn(row, 4)
		values_discriminated[row[typeIndex]] = append(values_discriminated[row[typeIndex]], new_row)
	}
	return values_discriminated
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
	if len(values_cleaned) != len(values) {
		fmt.Print("> Choisissez à quel categorie appartient chaque transaction:\n\n")

		for index, row := range values {
			// skip processed rows
			if index < len(values_cleaned) {
				continue
			}

			new_row := row
			//remove 2 last columns : currency, short label
			new_row = removeColumn(new_row, 4)
			new_row = removeColumn(new_row, 1)
			//remove card num if it exists in label col
			new_row = removeCardNum(new_row, 1)
			new_row = moveColumn(new_row, 0, 1)
			// add a transaction date column as the first column
			new_row = addTransactionDateCol(new_row, 0)
			// add a transaction type column
			new_row = addTypeColumn(new_row, 3)
			// prompt to add category column
			new_row = addCategoryColumn(new_row)
			values_cleaned = append(values_cleaned, new_row)
			saveCheckpoint(values_cleaned, csvPath)
		}
	}
	return values_cleaned
}

// Desc: Main function for the editing of the accounting exported csv
// It removes original array column names, process the array and return the new array with new column names
// Parameters: csv file as 2d array, csvPath is the full path of csv file, load found tmp as boolean
// Output: map of csv per transaction types
func FormatAccountingCSV(values [][]string, csvPath string, loadTmp bool) map[string][][]string {

	var formated_values_by_type map[string][][]string
	// Remove column names
	values = values[1:]
	values = editColumns(values, csvPath, loadTmp)
	formated_values_by_type = discriminateByType(values, 4)
	// Delete checkpoint file of work in progress
	//deleteCheckpoint(csvPath)
	// add column names
	values = append([][]string{COL_NAMES}, values...)
	return formated_values_by_type
}

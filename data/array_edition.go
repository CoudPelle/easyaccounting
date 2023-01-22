package data

import (
	"fmt"
	"easyaccounting/utils"
)

// Desired columns for final csv
var (
	COL_NAMES = []string{
		"Date prelevement", "Label", "Montant", "Date transaction", "Type"}
	TRANSACTION_TYPES = []string{
		"Cadeau", "Don", "Divers", "Vacances", "Administratif", "Logement",
		"Vetement" , "Media", "Abonnement", "Transport", "Sortie", "Bouffe" , "Loisir", "Sante" }
)

// rm end of card number from column 1
func removeCardNum(values [][]string) {
	for _, row := range values {
		row[1] = row[1][12:]
	}
}

// Extract the transaction date contained in column 1 for each row add it as a new column
func addTransactionDateCol(row []string) []string {
	var transDate string

	transDate = row[1][:5]
	row[1] = row[1][6:]
	row = append(row, transDate)

	return row
}

// Add a transaction type column.
// Prompt the user to choose in which type they want to classify this transaction
func addTypeColumn(row []string) []string {

	fmt.Println("Enter the name of your SG account exported csv")
	utils.StrInput()
	return row
}

// Add 2 columns which changes shape of the array
func addColumns(values [][]string) [][]string {
	var values_cleaned [][]string

	for _, row := range values {
		var new_row []string
		new_row = addTransactionDateCol(row)
		values_cleaned = append(values_cleaned, new_row)
	}
	return values_cleaned
}

// main function for the editing of the accounting exported csv
// it edits multiple fields and add columns
func FormatAccountingCSV (values [][]string) [][]string {
	removeCardNum(values)
	values = addColumns(values)
	return values
}
/*
This program aims to help to do your accounting.
It takes as input a csv export of your account transactions.

You have to provide it an "input" folder, located in the same folder of this binary

The program selects only desired columns and clean data
For now this program works only for Societe Generale exported csv

Usage: ./easyaccounting / ./easyaccounting.exe
*/
package main

import (
	"easyaccounting/data"
	"easyaccounting/utils"
	"fmt"
)

// main function of the module
func main() {
	var formated_values_by_type map[string][][]string
	values, csvPath, loadTmp := utils.GetCSV()
	formated_values_by_type = data.FormatAccountingCSV(values, csvPath, loadTmp)
	fmt.Println("> Sauvegarde des resultats dans le dossier \"output\".")
	utils.SaveResults(formated_values_by_type, csvPath)
}

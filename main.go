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
	"strings"
)

// main function of the module
func main() {
	values, csvPath, loadTmp := utils.GetCSV()
	values = data.FormatAccountingCSV(values, csvPath, loadTmp)
	fmt.Println("> Sauvegarde de votre fichier dans le dossier \"output\".")
	utils.WriteCSV(values, strings.Replace(csvPath, "input", "output", 1))
}

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
	"strings"
)
// main function of the module
func main() {
	values, csvPath := utils.GetCSV()
	values = data.FormatAccountingCSV(values, csvPath)
	utils.WriteCSV(values, strings.Replace(csvPath, "input", "output", 1))
}

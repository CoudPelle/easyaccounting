// This program aims to help to do your accounting.
// It takes as input a csv export of your account transactions.
// It selects only desired columns and clean data
// For now this program works only for Societe General exported csv

package main

import (
	"easyaccounting/data"
	"easyaccounting/utils"
	"strings"
)

func main() {
	values, csvPath := utils.GetCSV()
	values = data.FormatAccountingCSV(values, csvPath)
	utils.WriteCSV(values, strings.Replace(csvPath, "input", "output", 1))
}

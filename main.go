// This program aims to help to do your accounting.
// It takes as input a csv export of your account transactions.
// It selects only desired columns and clean data
// For now this program works only for Societe General exported csv

package main

import (
	"strconv"
	"easyaccounting/utils"
	"easyaccounting/data"
	"strings"
)

func appendSum(rows [][]string) {
	rows[0] = append(rows[0], "SUM")
	for i := 1; i < len(rows); i++ {
		rows[i] = append(rows[i], sum(rows[i]))
	}
}

func sum(row []string) string {
	sum := 0
	for _, s := range row {
		x, err := strconv.Atoi(s)
		if err != nil {
			return "NA"
		}
		sum += x
	}
	return strconv.Itoa(sum)
}

func main() {
	values, csvPath := utils.GetCSV()
	values = data.FormatAccountingCSV(values)
	
	utils.WriteCSV(values, strings.Replace(csvPath, "input", "output", 1))
}

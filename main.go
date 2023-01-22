// This program aims to help to do your accounting.
// It takes as input a csv export of your account transactions.
// It selects only desired columns and clean data
// For now this program works only for Societe General exported csv

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"easyaccounting/utils"
	"easyaccounting/data"
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

func writeChanges(rows [][]string) {
	f, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}
	err = csv.NewWriter(f).WriteAll(rows)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//columns := rows[0]
	values := utils.GetCSV()
	//newColNames(rows)
	values = data.FormatAccountingCSV(values)
	fmt.Println(values[0])
	// appendSum(rows)
	// writeChanges(rows)
}

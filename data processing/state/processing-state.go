package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func ProcessState(fileName string, states map[string]int) map[string]int {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("error opening csv file: %v\n", err)
		return states
	}
	defer file.Close()

	reader := csv.NewReader(file)
	//var book Book
	for {
		row, err2 := reader.Read() // row is an array of strings for each col in row
		if err2 != nil {
			return states
		}

		states[row[6]]++

		//fmt.Println(row[6])

	}

	//return states
}

func main() {
	stateMap := map[string]int{
		"Alaska":         0,
		"Alabama":        0,
		"Arizona":        0,
		"Arkansas":       0,
		"California":     0,
		"Colorado":       0,
		"Connecticut":    0,
		"Delaware":       0,
		"Florida":        0,
		"Georgia":        0,
		"Hawaii":         0,
		"Idaho":          0,
		"Illinois":       0,
		"Indiana":        0,
		"Iowa":           0,
		"Kansas":         0,
		"Kentucky":       0,
		"Louisiana":      0,
		"Maine":          0,
		"Maryland":       0,
		"Massachusetts":  0,
		"Michigan":       0,
		"Minnesota":      0,
		"Mississippi":    0,
		"Missouri":       0,
		"Montana":        0,
		"Nebraska":       0,
		"Nevada":         0,
		"New Hampshire":  0,
		"New Jersey":     0,
		"New Mexico":     0,
		"New York":       0,
		"North Carolina": 0,
		"North Dakota":   0,
		"Ohio":           0,
		"Oklahoma":       0,
		"Oregon":         0,
		"Pennsylvania":   0,
		"Rhode Island":   0,
		"South Carolina": 0,
		"South Dakota":   0,
		"Tennssee":       0,
		"Texas":          0,
		"Utah":           0,
		"Vermont":        0,
		"Virginia":       0,
		"Washington":     0,
		"West Virginia":  0,
		"Wisconsin":      0,
		"Wyoming":        0,
	}
	fileName := "2023-2024-bans.csv"
	stateMap = ProcessState(fileName, stateMap)

	file, err := os.Create("stateBookBans.csv")
	if err != nil {
		fmt.Printf("error creating csv file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"State", "Bans"})
	if err != nil {
		fmt.Printf("error writing categories csv file: %v\n", err)
		return
	}

	for key, num := range stateMap {
		csvRow := []string{key, strconv.Itoa(num)}
		err = writer.Write(csvRow)
		if err != nil {
			fmt.Printf("error writing body csv file: %v\n", err)
			break
		}
		fmt.Printf("%s : %d\n", key, num)
	}
}

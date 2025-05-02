package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Book struct {
	title        string
	author       string
	secondAuthor string
	timesBanned  int
	numStates    int
	statesBanned map[string]int // string is the state and int is the number of times banned in that state
	byAdmin      int
	byFormal     int
	byInformal   int
}

func NewBook(bookTitle string) Book {
	return Book{
		title:        bookTitle,
		author:       "",
		secondAuthor: "",
		timesBanned:  0,
		numStates:    0,
		statesBanned: make(map[string]int),
		byAdmin:      0,
		byFormal:     0,
		byInformal:   0,
	}
}

func (book *Book) incBans() {
	book.timesBanned++
}

func (book *Book) incState(state string) {
	states := book.statesBanned
	states[state]++
	if states[state] == 1 {
		book.incNumStates()
	}
}

func (book *Book) incNumStates() {
	book.numStates++
}

func (book *Book) incAdmin() {
	book.byAdmin++
}

func (book *Book) incFormal() {
	book.byFormal++
}

func (book *Book) incInformal() {
	book.byInformal++
}

func (book *Book) addAuthor(author string) {
	book.author = author
}

func (book *Book) hasAuthor() bool {
	// if book.author == "" {
	// 	return false
	// }

	// return true
	return book.author != ""
}

func (book *Book) addSecAuthor(author string) {
	book.secondAuthor = author
}

func (book *Book) hasSecAuthor() bool {
	// if book.secondAuthor == "" {
	// 	return false
	// }

	// return true
	return book.secondAuthor != ""
}

func getObj(title string, objs []Book) (Book, []Book) {
	for _, obj := range objs {
		if obj.title == title {
			return obj, objs
		}
	}
	new := NewBook(title)
	return new, append(objs, new)
}

func (book Book) UpdateList(objs []Book) []Book {
	for ind, obj := range objs {
		if obj.title == book.title {
			objs[ind] = book
			return objs
		}
	}
	return objs
}

func (book Book) MapToString() string {
	toReturn := ""

	for key, value := range book.statesBanned {
		toReturn += key + ":"
		toReturn += strconv.Itoa(value) + "/"
	}

	return toReturn
}

func Process(fileName string, objects []Book) []Book {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("error opening csv file: %v\n", err)
		return objects
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var book Book
	for {
		row, err2 := reader.Read() // row is an array of strings for each col in row
		if err2 != nil {
			return objects
		}
		book, objects = getObj(row[0], objects)
		if !book.hasAuthor() {
			book.addAuthor(row[1])
		}
		if !book.hasSecAuthor() {
			book.addSecAuthor(row[2])
		}

		book.incState(row[6])
		book.incBans()

		if strings.Contains(row[10], "Administration") {
			book.incAdmin()
		}
		if strings.Contains(row[10], "Formal Challenge") {
			book.incFormal()
		}
		if strings.Contains(row[10], "Informal Challenge") {
			book.incInformal()
		}

		objects = book.UpdateList(objects)

	}
}

func main() {
	objects := []Book{}
	fileName := "banned-books-2023-2024.csv"
	objects = Process(fileName, objects)

	file, err := os.Create("bannedBooks.csv")
	if err != nil {
		fmt.Printf("error creating csv file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Title", "Author", "SecondAuthor", "Times", "NumStates", "States", "Admin", "Formal", "Informal"})
	if err != nil {
		fmt.Printf("error writing categories csv file: %v\n", err)
		return
	}

	for _, obj := range objects {
		csvRow := []string{obj.title, obj.author, obj.secondAuthor, strconv.Itoa(obj.timesBanned), strconv.Itoa(obj.numStates), obj.MapToString(), strconv.Itoa(obj.byAdmin), strconv.Itoa(obj.byFormal), strconv.Itoa(obj.byInformal)}
		err = writer.Write(csvRow)
		if err != nil {
			fmt.Printf("error writing body csv file: %v\n", err)
			break
		}
		fmt.Println(obj)
	}
}

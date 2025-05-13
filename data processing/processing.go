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

// increases the number of times a book was banned
func (book *Book) incBans() {
	book.timesBanned++
}

/*
increases the ban count of the current book in the given state
increases the number of states the book has been banned in if that state is appearing for the first time
*/
func (book *Book) incState(state string) {
	states := book.statesBanned
	states[state]++
	if states[state] == 1 {
		book.incNumStates()
	}
}

// increases the number of states the book was challenged/banned in
func (book *Book) incNumStates() {
	book.numStates++
}

// increases the number of times a book was challenged by admininstarion
func (book *Book) incAdmin() {
	book.byAdmin++
}

// increases the number of times a book was challenged formally
func (book *Book) incFormal() {
	book.byFormal++
}

// increases the number of times a book was challenged informally
func (book *Book) incInformal() {
	book.byInformal++
}

// sets the author in the book object to the given name
func (book *Book) addAuthor(author string) {
	book.author = author
}

// returns a boolean value determining if the book's object has an author listed
func (book *Book) hasAuthor() bool {
	return book.author != ""
}

// sets the second author in the book object to the given name
func (book *Book) addSecAuthor(author string) {
	book.secondAuthor = author
}

// returns a boolean value determining if the book's object has a second author listed
func (book *Book) hasSecAuthor() bool {
	return book.secondAuthor != ""
}

/*
if an object with the given title is in the list, it returns it
otherwise, it makes a new object with that title and adds it to the list
*/
func getObj(title string, objs []Book) (Book, []Book) {
	for _, obj := range objs {
		if obj.title == title {
			return obj, objs
		}
	}
	new := NewBook(title)
	return new, append(objs, new)
}

// updates the list of books so it doesn't lose the current & previous data on the next iteration
func (book Book) UpdateList(objs []Book) []Book {
	for ind, obj := range objs {
		if obj.title == book.title {
			objs[ind] = book
			return objs
		}
	}
	return objs
}

// turning the map of states into a string so it can be written to the csv file
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

	// using methods to add relevant information from csv file to the objects

	// reading data from csv file line by line
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

/*
program does not filter out the heading data from PEN America so you will need to modify
the code or manually delete it
*/
func main() {
	// initializing list of struct objects
	objects := []Book{}

	// establishing the csv file to read data from
	fileName := "banned-books-2023-2024.csv" // you may need to move the data file into a different folder if it can't find it
	objects = Process(fileName, objects)

	// creating the new csv file
	file, err := os.Create("bannedBooks.csv")
	if err != nil {
		fmt.Printf("error creating csv file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// writing in the column names
	err = writer.Write([]string{"Title", "Author", "SecondAuthor", "Times", "NumStates", "States", "Admin", "Formal", "Informal"})
	if err != nil {
		fmt.Printf("error writing categories csv file: %v\n", err)
		return
	}

	// writing each object to the CSV file
	for _, obj := range objects {
		csvRow := []string{obj.title, obj.author, obj.secondAuthor, strconv.Itoa(obj.timesBanned), strconv.Itoa(obj.numStates), obj.MapToString(), strconv.Itoa(obj.byAdmin), strconv.Itoa(obj.byFormal), strconv.Itoa(obj.byInformal)}
		err = writer.Write(csvRow)
		if err != nil {
			fmt.Printf("error writing body csv file: %v\n", err)
			break
		}
	}
}

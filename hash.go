package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// var allbooks []inta
var scanned []bool
var libraries []library
var allBooks []book
var numDays int
var resString string

var libWins map[int]int

type library struct {
	id                  int
	books               []book
	mostRecentlyScanned int
	signedUp            bool
	signupTime          int
	maxPerDay           int
	libarayVal          float64
	allBooksScanned     bool
	booksScanned        string
	hasScanned          bool
}

type book struct {
	id      int
	val     int
	scanned bool
}

func selectRandomSqrt(libs []library) []library {
	size := int(math.Sqrt(float64(len(libs))))
	res := make([]library, size)

	idx := rand.Perm(size)
	for i := 0; i < size; i++ {
		res[i] = libs[idx[i]]
	}

	return res
}

func calculateUniquenessTable() {
	testLib := selectRandomSqrt(libraries)
	libWins = make(map[int]int)

	for _, libFor := range libraries {
		log.Println("Starting uniqueness for " + strconv.Itoa(libFor.id))

		for _, libAgainst := range testLib {
			if libFor.id != libAgainst.id {
				libForVals := libFor.books
				libAgainstVals := libAgainst.books

				libForUScore := 0
				for _, x := range libFor.books {
					if !contains(libAgainstVals, x.id) {
						libForUScore += x.val
					}
				}

				libAgainstUScore := 0
				for _, y := range libAgainst.books {
					if !contains(libForVals, y.id) {
						libAgainstUScore += y.val
					}
				}

				if libForUScore >= libAgainstUScore {
					libWins[libFor.id] = libWins[libFor.id] + 1
				}
			}
		}
		log.Println("Ending uniqueness for " + strconv.Itoa(libFor.id))
	}
}

func libVal(b []book, libID, maxPerDay, timeForSignup int) float64 {
	tot := 0
	for _, val := range b {
		tot += val.val
	}

	tot += libWins[libID]

	if timeForSignup >= numDays {
		return 0
	}

	return float64(tot) / float64(timeForSignup)
}

func recalcLibVals() {
	for i, l := range libraries {
		libraries[i].libarayVal = libVal(l.books, l.id, l.maxPerDay, l.signupTime)
	}
}

func newLibrary(id int, books []book, signUpTime int, maxPerDay int) library {
	sortByBookValue(books)
	return library{
		id:                  id,
		books:               books,
		signedUp:            false,
		signupTime:          signUpTime,
		maxPerDay:           maxPerDay,
		libarayVal:          0,
		allBooksScanned:     false,
		mostRecentlyScanned: len(books) - 1,
		booksScanned:        "",
		hasScanned:          false,
	}
}

func newBook(id, val int) book {
	return book{
		id:      id,
		val:     val,
		scanned: false,
	}
}

func scanFromCurrentLibs() int {
	sum := 0
	for index, lib := range libraries {
		// dont scan if not signed up / no books left
		if !lib.signedUp || lib.allBooksScanned {
			continue
		}
		i := lib.mostRecentlyScanned
		nbBooksScanned := 0
		bookIds := make([]int, 0)
		for ; i >= 0 && nbBooksScanned <= lib.maxPerDay; i-- {
			indexAllBooks := lib.books[i].id
			cBook := allBooks[indexAllBooks]
			if !cBook.scanned {
				sum += cBook.val
				allBooks[indexAllBooks].scanned = true
				nbBooksScanned++
				bookIds = append(bookIds, cBook.id)
				libraries[index].hasScanned = true
			}
		}
		// library is done
		if i == -1 {
			libraries[index].mostRecentlyScanned = -1
			libraries[index].allBooksScanned = true
		}

		if len(bookIds) != 0 {
			s := ""
			//s += strconv.Itoa(lib.id) + " " + strconv.Itoa(len(bookIds)) + "\n"
			for _, item := range bookIds {
				s += strconv.Itoa(item) + " "
			}
			libraries[index].booksScanned += s
		}
	}
	return sum
}

func simulate() {
	log.Println("Starting Uniquness calculation. . .")
	calculateUniquenessTable()
	recalcLibVals()
	log.Println("Uniquness calculation done. . .")

	// prepare libraries
	log.Println("Starting sorting Libraries. . .")
	sortLibrariesByLibraryVal()
	log.Println("Sorting Libraries done. . .")

	log.Println("Starting day simulation. . .")
	totalScore := 0
	currentSignupIndex := 0
	for numDays >= 0 {
		totalScore += scanFromCurrentLibs()
		if currentSignupIndex < len(libraries) {
			//Tick for signup progress
			libraries[currentSignupIndex].signupTime--

			//Check if signup complete
			if libraries[currentSignupIndex].signupTime == 0 {
				libraries[currentSignupIndex].signedUp = true
				currentSignupIndex++
			}
		}
		numDays--
		//Scan books from signed up libs
	}
	fmt.Printf("The total score is : %d \n", totalScore)
	countLibs := 0
	for _, item := range libraries {
		if item.hasScanned {
			countLibs++
		}
	}
	for _, l := range libraries {
		if len(l.booksScanned) == 0 {
			continue
		}
		resString += strconv.Itoa(l.id) + " " + strconv.Itoa(len(strings.Split(l.booksScanned, " "))-1) + "\n" + l.booksScanned + "\n"
		//fmt.Printf(resString)
	}
	writeToFile(strconv.Itoa(countLibs) + "\n" + resString)
}

func sortByBookValue(b []book) {
	// sort.Slice(libraries, func(i, j int) bool { return libraries[i] < libraries[j].val })
	sort.Slice(b, func(i, j int) bool { return b[i].val < b[j].val })
}

func sortLibrariesByLibraryVal() {
	sort.Slice(libraries, func(i, j int) bool { return libraries[i].libarayVal > libraries[j].libarayVal })
}

func main() {
	libraries = make([]library, 0)
	allBooks = make([]book, 0)

	file, err := os.Open("f_libraries_of_the_world.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	i := 0

	currentLib := 0
	currentLibSignupTime := 0
	currentLibMaxPerDay := 0

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")

		if i == 0 {
			numDays, _ = strconv.Atoi(line[2])
		} else if i == 1 {
			for index, elem := range line {
				temp, _ := strconv.Atoi(elem)
				allBooks = append(allBooks, newBook(index, temp))
			}
		} else {
			if i%2 == 0 {
				//lib info
				currentLibSignupTime, _ = strconv.Atoi(line[1])
				currentLibMaxPerDay, _ = strconv.Atoi(line[2])
			} else {
				//book ids for lib
				booksLocal := make([]book, 0)
				for _, elem := range line {
					temp, _ := strconv.Atoi(elem)
					booksLocal = append(booksLocal, allBooks[temp])
				}

				libraries = append(libraries, newLibrary(currentLib, booksLocal, currentLibSignupTime, currentLibMaxPerDay))
				currentLib++
			}
		}
		i++
	}

	simulate()
}

func writeToFile(output string) {
	ctime := time.Now().String()
	var outputFilename string = "output-" + ctime + ".txt"
	f, err := os.Create(outputFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(output)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

/*
func writeResults(pTypes []int, resList []int) {
	finStr := ""
	for _, final := range resList {
		finStr = finStr + strconv.Itoa(final) + " "
	}

	finStr = finStr[:len(finStr)-1]
	writeToFile(strconv.Itoa(len(resList)) + "\n" + finStr + "\n")
	diff := same(pTypes, resList)
	s, _ := json.Marshal(diff)
	writeToFile(string(s))
}
*/

func contains(s []book, id int) bool {
	for _, a := range s {
		if a.id == id {
			return true
		}
	}
	return false
}

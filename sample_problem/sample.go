package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Set Difference: A - B
func difference(a, b []int) (diff []int) {
	m := make(map[int]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
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

func writeResults(pTypes []int, resList []int) {
	finStr := ""
	for _, final := range resList {
		finStr = finStr + strconv.Itoa(final) + " "
	}

	finStr = finStr[:len(finStr)-1]
	writeToFile(strconv.Itoa(len(resList)) + "\n" + finStr + "\n")
	diff := difference(pTypes, resList)
	s, _ := json.Marshal(diff)
	writeToFile(string(s))
}

func calcPizza(goal int, pTypes []int) {
	current := 0

	resList := make([]int, 0)
	lastAdded := 0
	for _, e := range pTypes {
		if current+e <= goal {
			current += e
			lastAdded = e
			resList = append(resList, e)
		} else if (current-lastAdded)+e <= goal {
			current -= lastAdded
			resList = resList[:len(resList)-1]
			current += e
			lastAdded = e
			resList = append(resList, e)
		} else if current+e > goal {
			break
		}
	}

	fmt.Printf("%d closest you can get to %d\n", current, goal)

	writeResults(pTypes, resList)
}

func main() {
	file, err := os.Open("e_also_big.in")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	i := 0

	goal := 0
	numTypes := 0

	pTypes := make([]int, 0)

	for scanner.Scan() {
		switch i {
		case 0:
			line := strings.Split(scanner.Text(), " ")

			goal, _ = strconv.Atoi(line[0])
			numTypes, _ = strconv.Atoi(line[1])
		case 1:
			pTypes = make([]int, numTypes)

			line := strings.Split(scanner.Text(), " ")

			for index, elem := range line {
				temp, _ := strconv.Atoi(elem)

				pTypes[index] = temp
			}
		}
		i++
	}

	calcPizza(goal, pTypes)
}

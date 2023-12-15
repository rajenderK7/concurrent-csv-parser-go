package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func countChars(records [][]string) {
	count := 0
	for _, record := range records {
		for _, val := range record {
			count += len(val)
		}
	}
	fmt.Println("Total characters:", count)
}

func countCharsConcurrently(records [][]string, start int, countChan chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	end := start + 10
	count := 0
	for i := start; i < end; i++ {
		for _, val := range records[i] {
			count += len(val)
		}
	}
	countChan <- count
}

func main() {
	csvFile, err := os.Open("mock.csv")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records := make([][]string, 0)
	start := time.Now()
	// This is can be done using csv.ReadAll() too.
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err.Error())
		}
		records = append(records, record)
	}
	// find the total number of characters in the file

	// Sequential
	// countChars(records)
	// fmt.Println("Time taken:", time.Since(start))

	// Concurrently processing the records
	nGoroutines := len(records) / 10
	countChan := make(chan int)
	totalChars := 0
	wg := sync.WaitGroup{}
	go func() {
		for i := 0; i < nGoroutines; i++ {
			wg.Add(1)
			go countCharsConcurrently(records, i*10, countChan, &wg)
		}
		wg.Wait()
		close(countChan)
	}()
	for val := range countChan {
		totalChars += val
	}
	fmt.Println("Time taken:", time.Since(start))
	fmt.Println(totalChars)
}

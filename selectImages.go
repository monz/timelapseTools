package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"flag"
	"os"
	"log"
	"strconv"
	"sort"
	"time"
)

// read command line options
var file string
var picsPerDay int
func init() {
	flag.StringVar(&file, "f", "./timestamps.txt", "file containing timestamps")
	flag.IntVar(&picsPerDay, "p", 3, "pictures to select per day")
}

// create helper type to sort int64
type Int64Slice []int64
func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func selectImages(timestamps []int64, picsPerDay int) []int64 {
	if len(timestamps) <= 0 {
		return nil
	}

	sort.Sort(Int64Slice(timestamps))
	lowerBoundTimestamp := timestamps[0] // timestamp of first image
	day := getDay(lowerBoundTimestamp)
	lastDay := getDay(timestamps[len(timestamps)-1])

	var pics []int64
	for day < lastDay {
		for j := 0; j < picsPerDay; j++ {
			lowerBoundTimestamp = getNextImageTimestamp(timestamps, lowerBoundTimestamp)
			if lowerBoundTimestamp == -1 {
				break
			}
			pics = append(pics, lowerBoundTimestamp)
			lowerBoundTimestamp += 60 // plus one minute
		}
		if lowerBoundTimestamp == -1 {
			lowerBoundTimestamp = pics[len(pics)-1]
		}
		lowerBoundTimestamp += 60*60*24 // plus one days
		day = getDay(lowerBoundTimestamp)
	}

	return pics
}

func getNextImageTimestamp(timestamps []int64, lowerBoundTimestamp int64) int64 {
	upperBoundTimestamp := lowerBoundTimestamp + 5*60 // plus five minutes

	for i := 0; i < len(timestamps); i++ {
		timestamp := timestamps[i]
		if timestamp > lowerBoundTimestamp && timestamp <= upperBoundTimestamp {
			return timestamp
		} else if timestamp > upperBoundTimestamp {
			upperBoundTimestamp += 60*60*24 // plus one day
		}
	}

	return -1
}

func getDay(timestamp int64) int {
	t := time.Unix(timestamp, 0)
	return t.YearDay()
}

func readTimestamps() []int64 {
	// exit if file not exists
	_, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Could not read timestamps")
	}
	timestampStrings := strings.Split(string(content), "\n")
	var timestamps []int64
	for _, t := range timestampStrings {
		timestamp, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			log.Printf("Could not interpret timestamp as integer:'%v'", t)
			continue
		}
		timestamps = append(timestamps, timestamp)
	}
	return timestamps
}

func main() {
	flag.Parse()

	timestamps := readTimestamps()
	pics := selectImages(timestamps, picsPerDay)

	for _, p := range pics {
		fmt.Println(p)
	}
}

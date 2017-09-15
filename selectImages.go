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
var days int
func init() {
	flag.StringVar(&file, "f", "./timestamps.txt", "file containing timestamps")
	flag.IntVar(&picsPerDay, "p", 3, "pictures to select per day")
	flag.IntVar(&days, "d", 210, "number of days to consider")
}

// create helper type to sort int64
type Int64Slice []int64
func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func selectImages(timestamps []int64, days int, picsPerDay int) []int64 {
	if len(timestamps) <= 0 {
		return nil
	}

	sort.Sort(Int64Slice(timestamps))
	lowerBoundTimestamp := timestamps[0] // timestamp of first image
//fmt.Println("First timestamp:", lowerBoundTimestamp)
	var pics []int64

	for i := 0; i < days; i++ {
		for j := 0; j < picsPerDay; j++ {
//fmt.Println("Input lowerboundTimestamp", lowerBoundTimestamp)
			lowerBoundTimestamp = getNextImageTimestamp(timestamps, lowerBoundTimestamp)
			if lowerBoundTimestamp == -1 {
				//lowerBoundTimestamp = pics[len(pics)-1] + 60
//fmt.Println("CONTINUED")
				break
			}
//fmt.Println("Added pic:", lowerBoundTimestamp)
			pics = append(pics, lowerBoundTimestamp)
			lowerBoundTimestamp += 60 // plus one minute
		}
//fmt.Println("Lowerbound after pic loop:", lowerBoundTimestamp)
		if lowerBoundTimestamp == -1 {
			lowerBoundTimestamp = pics[len(pics)-1]
		}
		lowerBoundTimestamp += 60*60*24 // plus one day
	}

	return pics
}

func getNextImageTimestamp(timestamps []int64, lowerBoundTimestamp int64) int64 {
	day := getDay(lowerBoundTimestamp)
	lastDay := getDay(timestamps[len(timestamps)-1])
	upperBoundTimestamp := lowerBoundTimestamp + 5*60 // plus five minutes
//fmt.Println(day, lastDay, upperBoundTimestamp)
	var pic int64
	for day <= lastDay {
		for i := 0; i < len(timestamps); i++ {
			timestamp := timestamps[i]
			if timestamp > lowerBoundTimestamp && timestamp <= upperBoundTimestamp {
				pic = timestamp
				return pic
			} else if timestamp > upperBoundTimestamp {
				upperBoundTimestamp += 60*60*24 // plus one day
			}
		}
		day++
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

	//timestamps := []int64{
	//	1504802572,
	//	1504802652,
	//	1504802752,
	//	1504802832,
	//	1504802932,
	//	1504803012,
	//	1504803112,
	//	1504803187,
	//	1504803252,
	//	1504803372,
	//	1504803472,
	//	1504803552,
	//	1504803652,
	//	1504803732,
	//	1504803832,
	//	1504803912,
	//	1504804012,
	//	1504804087,
	//	1504804152,
	//	1504804272,
	//	1505399248,
	//	1505399352,
	//	1505399452,
	//	1505399527,
	//	1505399592,
	//	1505399686,
	//	1505399772,
	//	1505399889,
	//	1505399952,
	//	1505400072,
	//	1505400172,
	//	1505400252,
	//	1505400352,
	//}

	pics := selectImages(timestamps, days, picsPerDay)

	for _, p := range pics {
		fmt.Println(p)
	}
}

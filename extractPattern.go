package main

import (
	"log"
	"regexp"
	"os"
	"flag"
	"path/filepath"
	"fmt"
)

// read command line options
var basePath string
var recursive bool
var regex string
func init() {
	flag.StringVar(&basePath, "p", "./", "directory which contains names to extract")
	flag.BoolVar(&recursive, "r", false, "read files under each directory, recursive.")
	flag.StringVar(&regex, "e", "\\d{10}", "regex string to match")
}

func extractPattern(reg *regexp.Regexp, timestamps *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		if info.IsDir() && (path != basePath) && !recursive {
			return filepath.SkipDir
		} else if info.IsDir() {
			return nil
		}

		timestamp := reg.FindString(info.Name())
		if timestamp == "" {
			log.Println("Error in:", path)
			timestamp = "NA"
		}
		*timestamps = append(*timestamps, timestamp)
		return nil
	}
}

func main() {
	flag.Parse()

	log.Println("Extract Pattern")

	log.Println("Basepath:", basePath)
	log.Println("Recursive:", recursive)
	log.Println("Regex:", regex)

	reg, err := regexp.Compile(regex)
	if err != nil {
		log.Println("Could not compile regex")
		return
	}

	var timestamps []string
	err = filepath.Walk(basePath, extractPattern(reg, &timestamps))
	if err != nil {
		log.Println("Could not walk path:", err)
	}

	for _, v := range timestamps {
		//log.Println("Timestamp:", v)
		fmt.Println(v)
	}
}
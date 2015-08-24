package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

func parseYaml(file string) (entries Source, err error) {
	type SupFile struct {
		Supplements []Entry
	}

	contents, err := ioutil.ReadFile(file)
	var sup SupFile
	err = yaml.Unmarshal([]byte(contents), &sup)

	entries = sup.Supplements
	return
}

func parseMarkdown(file string) (entries Source, err error) {
	// e.g. 22.08.2015
	dayMatch, _ := regexp.Compile("^# ?(\\d{2}).(\\d{2}).(\\d{4})$")

	// e.g. * 12:25 250mg caffeine
	entryMatch, _ := regexp.Compile("^[\t ]*\\* (\\d{2}):(\\d{2}) ([\\.\\d]+)(g|mg) (.+)")

	h, err := os.Open(file)

	if err != nil {
		return
	}

	scanner := bufio.NewScanner(h)

	var day, month, year int

	for scanner.Scan() {
		if dayMatch.MatchString(scanner.Text()) {
			groups := dayMatch.FindAllStringSubmatch(scanner.Text(), -1)

			day, err = strconv.Atoi(groups[0][1])
			month, err = strconv.Atoi(groups[0][2])
			year, err = strconv.Atoi(groups[0][3])
		}

		if entryMatch.MatchString(scanner.Text()) {
			groups := entryMatch.FindAllStringSubmatch(scanner.Text(), -1)
			var hour, minute int
			var dose float64

			hour, err = strconv.Atoi(groups[0][1])
			minute, err = strconv.Atoi(groups[0][2])
			dose, err = strconv.ParseFloat(groups[0][3], 64)
			unit := groups[0][4]

			description := strings.Fields(groups[0][5])
			name := description[0]
			misc := description[1:]

			date := time.Date(year, time.Month(month), day, hour, minute, 0, 0,
				time.UTC)

			entries = append(entries, Entry{date, Unit(dose) *
				ParseUnit(unit), name, misc})
		}
		if err != nil {
			return
		}
	}
	return
}

func main() {
	inputFile := flag.String("merge", "", "Files")
	flag.Parse()

	// List of filenames to be merged
	inputFiles := []string{*inputFile}

	// Obtain a complete list of (existing) files to be merged
	for _, file := range flag.Args() {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			break
		}
		inputFiles = append(inputFiles, file)
	}

	// Sources to be merged
	sources := make([]Source, len(inputFiles))

	// Obtain an array of parsed sources to be merged
	for i, file := range inputFiles {
		var (
			err     error
			entries Source
		)
		switch filepath.Ext(file) {
		case ".md":
			entries, err = parseMarkdown(file)
		case ".yaml":
			entries, err = parseYaml(file)
		}
		if err == nil {
			sources[i] = entries
		}
	}

	yaml := merge(sources)
	fmt.Println(yaml)
}

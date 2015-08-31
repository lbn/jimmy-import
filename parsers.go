package main

import (
	"bufio"
	//"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Parser interface {
	Parse() (Source, error)
}

// Markdown parser
type MarkdownParser struct {
	File   string
	Schema Schema
}

func removeUnderscore(s string) string {
	return strings.Replace(s, "_", " ", -1)
}

func (parser MarkdownParser) Parse() (entries Source, err error) {

	// e.g. 22.08.2015
	dayMatch, _ := regexp.Compile("^# ?(\\d{2}).(\\d{2}).(\\d{4})$")

	// e.g. * 12:25 250mg caffeine
	entryMatch, _ := regexp.Compile("^[\t ]*\\* (\\d{2}):(\\d{2}) ([\\.\\d]+)(g|mg) (.+)")

	h, err := os.Open(parser.File)

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
			name := removeUnderscore(description[0])
			misc := make(map[string]string)

			// Schema
			if definition, ok := parser.Schema.Definitions[name]; ok {
				// number of space separated words
				if len(description[1:]) != len(definition.Order) {
					log.WithFields(log.Fields{
						"name":     name,
						"expected": definition.Order,
						"values":   description[1:],
					}).Fatal("Entry does not have the required parameters")
				}
				for i, prop := range definition.Order {
					misc[prop] = removeUnderscore(description[1+i])
				}
			}

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

type YamlParser struct {
	file string
}

func (parser YamlParser) Parse() (entries Source, err error) {
	type SupFile struct {
		Supplements []Entry
	}

	contents, err := ioutil.ReadFile(parser.file)
	var sup SupFile
	err = yaml.Unmarshal([]byte(contents), &sup)

	entries = sup.Supplements
	return
}

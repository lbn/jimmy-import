package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

func main() {

	inputFile := flag.String("merge", "", "Files")
	schemaFile := flag.String("schema", "", "YAML schema file")
	flag.Parse()

	if *inputFile == "" {
		log.Fatal("Provide at least one file to merge")
	}

	if *schemaFile == "" {
		log.Fatal("No schema file provided")
	}

	schema, err := NewSchema(*schemaFile)
	if err != nil {
		log.Fatal("Could not parse schema")
	}

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
		var parser Parser

		switch filepath.Ext(file) {
		case ".md":
			parser = MarkdownParser{file, schema}
		case ".yaml":
			parser = YamlParser{file}
		}
		if parser == nil {
			log.WithFields(log.Fields{
				"file": file,
			}).Fatal("Invalid file format given")
		}
		entries, err = parser.Parse()
		if err == nil {
			sources[i] = entries
		} else {
			log.WithFields(log.Fields{
				"file":    file,
				"message": err,
			}).Fatal("Could not parse file")
		}
	}

	yaml := merge(sources, &schema)
	fmt.Println(yaml)
}

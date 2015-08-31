package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Property struct {
	Type string
}

type Definition struct {
	Properties map[string]Property
	Order      []string
}

type Schema struct {
	Definitions map[string]Definition
}

func NewSchema(file string) (schema Schema, err error) {
	contents, err := ioutil.ReadFile(file)
	yaml.Unmarshal([]byte(contents), &schema)
	return
}

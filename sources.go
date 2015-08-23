package main

import (
	"gopkg.in/yaml.v2"
	"sort"
	"time"
)

type Entry struct {
	Date time.Time
	Dose Unit
	Name string
	Misc []string
}

// Source
type Source []*Entry

func (source *Source) ToYaml() string {
	type YamlNode map[string]interface{}

	nodes := make([]YamlNode, len(*source))
	for i, entry := range *source {
		m := make(YamlNode)
		m["date"] = entry.Date
		m["name"] = entry.Name
		m["dose"] = entry.Dose
		nodes[i] = m
	}

	root := make(YamlNode)
	root["supplements"] = nodes
	d, _ := yaml.Marshal(root)
	return string(d)
}

// Sort
type ByDate Source

func (entries ByDate) Len() int {
	return len(entries)
}
func (entries ByDate) Swap(i, j int) {
	entries[i], entries[j] = entries[j], entries[i]
}
func (entries ByDate) Less(i, j int) bool {
	return entries[i].Date.Before(entries[j].Date)
}

func merge(sources []Source) string {
	totalLen := 0
	for _, src := range sources {
		totalLen += len(src)
	}

	// Concatenate sources
	var combined Source
	for _, src := range sources {
		combined = append(combined, src...)
	}

	sort.Sort(ByDate(combined))

	return combined.ToYaml()
}

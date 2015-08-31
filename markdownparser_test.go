package main_test

import (
	. "."
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var parser MarkdownParser

var _ = Describe("MarkdownParser", func() {
	schema, _ := NewSchema("data/schema.yaml")
	parser := MarkdownParser{"data/example.md", schema}
	entries, _ := parser.Parse()

	Describe("Parsing a single file", func() {
		It("should parse at least some entries", func() {
			Expect(len(entries)).To(BeNumerically(">", 0))
		})
		It("should parse coffee entries", func() {
			for _, entry := range entries {
				if entry.Name == "coffee" {
					Expect(entry.Misc).To(HaveKeyWithValue("brand", "Taylors"))
					Expect(entry.Misc).To(HaveKeyWithValue("strength", "4"))
				}
			}
		})
	})
})

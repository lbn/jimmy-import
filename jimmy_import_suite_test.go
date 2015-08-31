package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestJimmyImport(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JimmyImport Suite")
}

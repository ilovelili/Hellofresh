package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHellofresh(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hellofresh Suite")
}

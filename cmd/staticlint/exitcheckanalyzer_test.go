package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestExitCheckAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ExitCheckAnalyzer, "./...")
}

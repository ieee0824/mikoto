package checker

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	old := maxLines
	maxLines = 0
	defer func() { maxLines = old }()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "example")
}

func TestFunctionLength(t *testing.T) {
	old := maxLines
	maxLines = 3
	defer func() { maxLines = old }()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "length")
}

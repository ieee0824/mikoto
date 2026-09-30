package checker

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	old := maxLines
	maxLines = 0
	defer func() { maxLines = old }()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "example", "genericrange", "pointeraccess", "calls")
}

func TestFunctionLengthDirectives(t *testing.T) {
	old := maxLines
	maxLines = 8
	defer func() { maxLines = old }()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "linedirectives")
}

func TestFunctionLength(t *testing.T) {
	old := maxLines
	maxLines = 3
	defer func() { maxLines = old }()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "length")
}

package main

import (
	"github.com/ieee0824/mikoto/internal/checker"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(checker.Analyzer) }

package main

import (
	"flag"
	"log"
	"video2music/internal"

	_ "go.uber.org/zap"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	inputDir := flag.String("in", "./content/input", "input dir")
	outputDir := flag.String("out", "./content/output", "output dir")

	flag.Parse()

	log.Println(*inputDir)
	log.Println(*outputDir)

	internal.NewConverter(*inputDir, *outputDir).Process()
}

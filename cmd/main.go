package main

import (
	"flag"
	"log"
	"os"

	"github.com/claudemuller/bin-patcher/internal/pkg"
)

var inFile = flag.String("in", "", "the binary file to patch")
var outFile = flag.String("out", "", "the destination of the patched binary file")
var sig = flag.String("sig", "", "the signature to patch")
var patch = flag.String("patch", "", "the patch apply to the binary file with")

func main() {
	log.SetPrefix("bin-patcher > ")
	flag.Parse()

	if *inFile == "" || *sig == "" || *patch == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *outFile == "" {
		*outFile = *inFile + ".patched"
	}

	if err := pkg.Patch(*inFile, *outFile, *sig, *patch); err != nil {
		log.Fatalf("patching %s failed: %v", *inFile, err)
	}
}

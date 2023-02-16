package main

import (
	"flag"
	l "log"
	"os"

	"github.com/claudemuller/bin-patcher/internal/pkg/gui"
	log "github.com/claudemuller/bin-patcher/internal/pkg/logger"
	"github.com/claudemuller/bin-patcher/internal/pkg/patcher"
)

func main() {
	inFile := flag.String("in", "", "the binary file to patch")
	outFile := flag.String("out", "", "the destination of the patched binary file")
	sig := flag.String("sig", "", "the signature to patch")
	patch := flag.String("patch", "", "the patch apply to the binary file with")
	flag.Parse()

	if *outFile == "" {
		*outFile = *inFile + ".patched"
	}

	app := gui.NewApp()

	if len(os.Args) > 1 {
		startCli(*inFile, *outFile, *sig, *patch, app.Logger)
		return
	}

	app.Run()
}

func startCli(inFile, outFile, sig, patch string, logger *log.Log) {
	if inFile == "" || sig == "" || patch == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := patcher.Patch(inFile, outFile, sig, patch, logger); err != nil {
		l.Fatalf("patching %s failed: %v", inFile, err)
	}
}

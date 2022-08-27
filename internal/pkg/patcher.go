package pkg

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
)

var ErrSignatureSizeMismatch = errors.New("signatures do not match in length")

var ErrSignatureNotFound = errors.New("signature was not found in the binary")

func Patch(inFile, outFile, sigStr, patchStr string) error {
	// Read the file.
	data, err := os.ReadFile(inFile)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	sig, patch, err := decodeSigs(sigStr, patchStr)
	if err != nil {
		return fmt.Errorf("error decoding signatures: %w", err)
	}

	// Check new sig matches old sig.
	if len(sig) != len(patch) {
		return ErrSignatureSizeMismatch
	}

	// Make a backup of the file.
	err = os.WriteFile(inFile + ".bak", data, 0770)
	if err != nil {
		return fmt.Errorf("error writing backup binary: %w", err)
	}

	// Patch.
	patchedData := patchData(data, sig, patch)
	if patchedData == nil {
		return ErrSignatureNotFound
	}

	// Write the patched file.
	log.Printf("writing patched binary to %s", outFile)

	err = os.WriteFile(outFile, patchedData, 0770)
	if err != nil {
		return fmt.Errorf("error writing patched binary: %w", err)
	}

	return nil
}

func checkSig(data []byte, sig []byte) bool {
	for i := range data {
		if data[i] != sig[i] {
			return false
		}
	}

	return true
}

func decodeSigs(sigStr, patchStr string) ([]byte, []byte, error) {
	sig, err := hex.DecodeString(sigStr)
	if err != nil {
		return nil, nil, err
	}

	patch, err := hex.DecodeString(patchStr)
	if err != nil {
		return nil, nil, err
	}

	return sig, patch, nil
}

func patchData(data, sig, patch []byte) []byte {
	for i := 0; i < len(data)-len(sig); i++ {
		if data[i] == sig[0] {
			if data[i+1] == sig[1] {
				fmt.Println(data[i], sig[0])
			}
		}

		if bytes.Compare(data[i:i+len(sig)], sig) == 0 {
			log.Printf("signature found at %#x, patching...", i)

			return doPatch(data, patch, i)
		}
	}

	return nil
}

func doPatch(data, patch []byte, loc int) []byte {
	var output = make([]byte, len(data))

	n := copy(output, data)
	if n < 0 {
		return []byte{}
	}

	c := 0

	for i := loc; i < loc+len(patch); i++ {
		output[i] = patch[c]
		c++
	}

	if bytes.Compare(data, output) != 0 {
		log.Printf("patched %v bytes", c)

		return output
	}

	return []byte{}
}
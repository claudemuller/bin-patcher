package pkg

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

var ErrSignatureSizeMismatch = errors.New("signatures do not match in length")

var ErrSignatureNotFound = errors.New("signature was not found in the binary")

func Patch(inFile, outFile, sigStr, patchStr string, logger *Log) error {
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
	if err = os.WriteFile(inFile+".bak", data, 0770); err != nil {
		return fmt.Errorf("error writing backup binary: %w", err)
	}

	// Patch.
	patchedData := patchData(data, sig, patch, logger)
	if patchedData == nil {
		return ErrSignatureNotFound
	}

	// Write the patched file.
	logger.log(fmt.Sprintf("writing patched binary to %s", outFile))

	if err = os.WriteFile(outFile, patchedData, 0770); err != nil {
		return fmt.Errorf("error writing patched binary: %w", err)
	}

	return nil
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

func patchData(data, sig, patch []byte, logger *Log) []byte {
	for i := 0; i < len(data)-len(sig); i++ {
		if bytes.Equal(data[i:i+len(sig)], sig) {
			logger.log(fmt.Sprintf("signature found at %#x, patching...", i))

			return doPatch(data, patch, i, logger)
		}
	}

	return nil
}

func doPatch(data, patch []byte, loc int, logger *Log) []byte {
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

	if !bytes.Equal(data, output) {
		logger.log(fmt.Sprintf("patched %v bytes", c))

		return output
	}

	return []byte{}
}

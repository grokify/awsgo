package inspector2util

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
)

// ReadFileListFindingsOutputStream streams findings from a ListFindingsOutput
// JSON file, invoking fn for each finding without loading the full findings
// array into memory.
func ReadFileListFindingsOutputStream(filename string, fn func(f types.Finding) error) error {
	if fn == nil {
		return errors.New("function cannot be nil")
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)

	// Read open curly brace
	token, err := decoder.Token()
	if err != nil {
		return err
	} else if token != json.Delim('{') {
		return errors.New("expected start of object")
	}

	// Loop through the top-level keys
	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return err
		}
		if key != "findings" {
			var discard json.RawMessage
			if err := decoder.Decode(&discard); err != nil {
				return err
			}
			continue
		}

		// Read opening bracket of the array
		token, err := decoder.Token()
		if err != nil {
			return err
		} else if token != json.Delim('[') {
			return errors.New("expected start of array")
		}

		for decoder.More() {
			var finding types.Finding
			if err := decoder.Decode(&finding); err != nil {
				return err
			} else if err := fn(finding); err != nil {
				return err
			}
		}

		// Read closing bracket of the array
		if _, err := decoder.Token(); err != nil {
			return err
		}
	}
	return nil
}

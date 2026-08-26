package solve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func ParseInput(data []byte) (*Input, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var in Input
	if err := dec.Decode(&in); err != nil {
		return nil, &ParseError{Reason: fmt.Sprintf("invalid JSON: %v", err)}
	}
	if err := ensureEOF(dec); err != nil {
		return nil, err
	}
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return &in, nil
}

func ParseInputFile(path string) (*Input, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, &ParseError{Reason: err.Error()}
	}
	return ParseInput(data)
}

func ensureEOF(dec *json.Decoder) error {
	var extra any
	err := dec.Decode(&extra)
	if err == io.EOF {
		return nil
	}
	if err == nil {
		return &ParseError{Reason: "unexpected extra JSON content after the first object"}
	}
	return &ParseError{Reason: fmt.Sprintf("invalid trailing content: %v", err)}
}

var readFile = osReadFile

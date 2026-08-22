package utils

import (
	"encoding/json"
	"io"
	"os"
)

func RetrieveJsonBytes(path string) ([]byte, error) {
	return []byte{}, nil
}

func WriteToFile[T any](path string, data T) error {
	jsonFile, err := os.Create(path)

	if err != nil {
		return err
	}

	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)

	err = encoder.Encode(data)

	return err
}

func DecodeReader[t any](r io.Reader) (t, error) {
	decoder := json.NewDecoder(r)

	var outputStruct t

	err := decoder.Decode(&outputStruct)

	if err != nil {
		return outputStruct, err
	}

	return outputStruct, nil
}

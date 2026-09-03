package utils

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
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

func GetBytesFromHTML(url string, userAgent string) ([]byte, error) {
	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func ParseBytesToXML[xmlSchema any](data []byte) (xmlSchema, error) {
	var parsedXml xmlSchema

	err := xml.Unmarshal(data, &parsedXml)

	return parsedXml, err
}

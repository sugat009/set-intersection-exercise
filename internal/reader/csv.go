package reader

import (
	"encoding/csv"
	"io"

	"github.com/pkg/errors"
)

// ReadKeysFromCsvIntoChannel reads csv content to find key for each row and push into the passed in channel
// returns when end of file is reached or when error
func ReadKeysFromCsvIntoChannel(key string, reader io.Reader, keysOuput chan<- string) error {
	if reader == nil {
		return errors.New("csv source is nil")
	}

	csvReader := csv.NewReader(reader)

	headerKeyIndex := -1

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.Wrap(err, "while reading from reader")
		}

		switch headerKeyIndex {
		case -1:
			headerKeyIndex, err = getIndex(row, key)
			if err != nil {
				return errors.Errorf("header: %s does not exist", key)
			}
		default:
			keysOuput <- row[headerKeyIndex]
		}
	}

	return nil
}

func getIndex(headers []string, key string) (int, error) {
	for idx, header := range headers {
		if header == key {
			return idx, nil
		}
	}
	return 0, errors.Errorf("key (%s) does not exist in header", key)
}

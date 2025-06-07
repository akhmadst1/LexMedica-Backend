package pkg

import (
	"encoding/json"
	"io/ioutil"
)

type Regulation struct {
	Document string `json:"document"`
	Article  string `json:"article"`
	Content  string `json:"content"`
}

type TestCase struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Regulations []Regulation `json:"regulations"`
	Reason      string       `json:"reason"`
	Disharmony  string       `json:"disharmony"`
}

// LoadTestCases loads the test cases from a JSON file
func LoadTestCases(filePath string) ([]TestCase, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var testCases []TestCase
	if err := json.Unmarshal(file, &testCases); err != nil {
		return nil, err
	}

	return testCases, nil
}

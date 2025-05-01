package pkg

import (
	"encoding/json"
	"io/ioutil"
)

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

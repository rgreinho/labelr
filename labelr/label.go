package labelr

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Label represents a label entry.
type Label struct {
	Name        string
	Description string
	Color       string
}

// Labels represents the content of a labelr file.
type Labels struct {
	Labels []Label `yaml:",flow"`
}

func (l *Label) String() string {
	return fmt.Sprintf("%v, %v, %v", l.Color, l.Name, l.Description)
}

func (l *Labels) String() string {
	var buffer bytes.Buffer

	for _, label := range l.Labels {
		buffer.WriteString(fmt.Sprintf("%s\n", label.String()))
	}

	return buffer.String()
}

func readFile(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ParseFile parses a file describing labels.
func ParseFile(file string) (*Labels, error) {
	data, err := readFile(file)
	if err != nil {
		return nil, err
	}
	return ParseDocument(data)
}

// ParseDocument parses a document describing labels.
func ParseDocument(document []byte) (*Labels, error) {
	labelDocument := Labels{}

	err := yaml.Unmarshal([]byte(document), &labelDocument)
	if err != nil {
		return nil, err
	}
	return &labelDocument, nil
}

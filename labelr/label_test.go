package labelr

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v1"
)

func TestParseDocument00(t *testing.T) {
	label := Label{
		Name:        "test",
		Description: "description",
		Color:       "#000000",
	}
	labels := &Labels{
		Labels: []Label{
			label,
		},
	}

	// Marshal the Labels.
	data, err := yaml.Marshal(&labels)
	if err != nil {
		t.Errorf("Error while marshalling the document: %s", err)
	}

	// Parse the document.
	got, err := ParseDocument(data)
	if err != nil {
		t.Errorf("Error while parsing the document: %s", err)
	}
	want := labels

	// Check.
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q want %q given, %q", got, want, data)
	}
}

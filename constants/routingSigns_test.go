package constants

import (
	"reflect"
	"testing"
)

func TestConstants(t *testing.T) {
	expected := map[string]string{
		"WildcardSign": WildcardSign,
		"ParamSign":    ParamSign,
		"OptionalSign": OptionalSign,
	}

	actual := map[string]string{
		"WildcardSign": "*",
		"ParamSign":    ":",
		"OptionalSign": "?",
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

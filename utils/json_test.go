package utils

import (
	"fmt"
	"testing"
)

func TestToJSON(t *testing.T) {
	t.Parallel()
	res, _ := ToJSON("MY/NAME/IS/:PARAM/*")

	fmt.Println(res)
}

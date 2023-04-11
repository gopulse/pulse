package utils

import (
	"fmt"
	"testing"
)

func TestToLower(t *testing.T) {
	t.Parallel()
	res := ToLower("MY/NAME/IS/:PARAM/*")
	res = ToLower("1MY/NAME/IS/:PARAM/*")
	res = ToLower("/MY2/NAME/IS/:PARAM/*")
	res = ToLower("/MY3/NAME/IS/:PARAM/*")
	res = ToLower("/MY4/NAME/IS/:PARAM/*")

	fmt.Println(res)
}

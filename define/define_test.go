package define

import (
	"fmt"
	"testing"
)

func TestTypeToName(t *testing.T) {
	res := TypeToName(2)
	fmt.Println(res)
}

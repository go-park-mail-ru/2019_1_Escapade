package cors

import (
	"testing"
)

func TestCors(t *testing.T) {
	IsAllowed("there", []string{"string", "is", "not", "there"})
}

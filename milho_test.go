package milho

import (
	"strings"
	"testing"
)

func Test_milho(t *testing.T) {
	response := Run("(+ 1 2 (- 3) 2)")

	if strings.TrimSpace(response) != "2" {
		t.Errorf("Expected 2, got '%s'", response)
	}
}

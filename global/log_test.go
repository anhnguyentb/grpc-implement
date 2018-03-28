package global

import "testing"

func TestLoadLogger(t *testing.T) {
	err := LoadLogger(true)
	if err != nil {
		t.Fatalf("Fatal error load logger: %s \n", err)
	}
}

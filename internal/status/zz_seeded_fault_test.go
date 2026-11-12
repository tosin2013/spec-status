package status

import "testing"

// Seeded fault for feature 003 T024; never merged.
func TestSeededFault(t *testing.T) { t.Fatal("seeded fault: CI must fail") }

package generator

import "testing"

func TestHashToUint64_Deterministic(t *testing.T) {
	a := HashToUint64("hello")
	b := HashToUint64("hello")

	if a != b {
		t.Errorf("same input produced different hashes: %d vs %d", a, b)
	}
}

func TestHashToUint64_DifferentInputs(t *testing.T) {
	a := HashToUint64("hello")
	b := HashToUint64("world")

	if a == b {
		t.Errorf("different inputs produced same hash: %d", a)
	}
}

func TestHashToUint64_KnownValue(t *testing.T) {
	got := HashToUint64("test")
	want := uint64(11495104353665842533)

	if got != want {
		t.Errorf("HashToUint64(%q) = %d, want %d", "test", got, want)
	}
}

func TestHashToUint64_EmptyString(t *testing.T) {
	// Should not panic and should return a consistent value
	a := HashToUint64("")
	b := HashToUint64("")

	if a != b {
		t.Errorf("empty string produced different hashes: %d vs %d", a, b)
	}
}

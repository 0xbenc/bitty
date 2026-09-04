package state

import "testing"

func TestIntroRoundTrip(t *testing.T) {
	dir := t.TempDir()
	if err := Write(dir, "v1.2.3"); err != nil {
		t.Fatal(err)
	}
	got, err := Read(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.SchemaVersion != 1 || got.LastSeen != "v1.2.3" {
		t.Fatalf("Read() = %+v", got)
	}
}

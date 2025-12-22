package util_test

import (
	"os"
	"path/filepath"
	"testing"

	"backtester/util"
)

func TestLoadCSVColumn_ReversesOrder(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "prices.csv")
	content := "Date,Price,CapMVRVCur\n2024-01-02,200,1.5\n2024-01-01,100,1.0\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	values, err := util.LoadCSVColumn(path, "Price", true)
	if err != nil {
		t.Fatalf("LoadCSVColumn returned error: %v", err)
	}
	want := []float64{100, 200}
	for i := range want {
		if values[i] != want[i] {
			t.Fatalf("index %d got %.2f want %.2f", i, values[i], want[i])
		}
	}
}

func TestLoadCSVColumnWithDefault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.csv")
	content := "Col\n1.0\n\"\"\nabc\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	values, err := util.LoadCSVColumnWithDefault(path, "Col", false, 5.0)
	if err != nil {
		t.Fatalf("LoadCSVColumnWithDefault returned error: %v", err)
	}
	want := []float64{1.0, 5.0, 5.0}
	if len(values) != len(want) {
		t.Fatalf("length mismatch got %d want %d", len(values), len(want))
	}
	for i := range want {
		if values[i] != want[i] {
			t.Fatalf("index %d got %.2f want %.2f", i, values[i], want[i])
		}
	}
}

func TestLoadCSVVolumeColumn(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vol.csv")
	content := "Vol.\n1.5K\n2M\n-\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	values, err := util.LoadCSVVolumeColumn(path, "Vol.", false)
	if err != nil {
		t.Fatalf("LoadCSVVolumeColumn returned error: %v", err)
	}
	want := []float64{1500, 2_000_000, 0}
	for i := range want {
		if values[i] != want[i] {
			t.Fatalf("index %d got %.0f want %.0f", i, values[i], want[i])
		}
	}
}

func TestParseDateFlexible(t *testing.T) {
	cases := []struct {
		input string
		ok    bool
	}{
		{"2024-01-02", true},
		{"01/31/2024", true},
		{"31/01/2024", true},
		{"\"11/23/2025\"", true},
		{"bad-date", false},
	}
	for _, c := range cases {
		_, err := util.ParseDateFlexible(c.input)
		if c.ok && err != nil {
			t.Fatalf("expected success for %q got error %v", c.input, err)
		}
		if !c.ok && err == nil {
			t.Fatalf("expected failure for %q", c.input)
		}
	}
}

package opencode

import "testing"

func TestParseConfigDir(t *testing.T) {
	out := []byte("home       /Users/me\nconfig     /Users/me/Library/Application Support/opencode\ndata       /Users/me/.local/share/opencode\n")
	got, err := ParseConfigDir(out)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	want := "/Users/me/Library/Application Support/opencode"
	if got != want {
		t.Fatalf("config dir = %q, want %q", got, want)
	}
}

func TestParseConfigDirWindowsPath(t *testing.T) {
	out := []byte("config     C:\\Users\\me\\AppData\\Roaming\\opencode\r\n")
	got, err := ParseConfigDir(out)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	want := "C:\\Users\\me\\AppData\\Roaming\\opencode"
	if got != want {
		t.Fatalf("config dir = %q, want %q", got, want)
	}
}

func TestParseConfigDirMissing(t *testing.T) {
	if _, err := ParseConfigDir([]byte("home /Users/me\n")); err == nil {
		t.Fatal("expected an error when config is missing")
	}
}

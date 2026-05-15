package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderFileDefaultHeadAndTail(t *testing.T) {
	path := writeTempFile(t, 12)

	out, err := renderFile(path, options{head: 10, tail: 10, theme: "plain"})
	if err != nil {
		t.Fatalf("renderFile returned error: %v", err)
	}

	if !strings.HasPrefix(out, path+"\n") {
		t.Fatalf("output should start with filename, got:\n%s", out)
	}
	if !strings.Contains(out, "\n"+separator+"\n") {
		t.Fatalf("output should contain separator, got:\n%s", out)
	}
	if !strings.Contains(out, "     1 | line 01\n") {
		t.Fatalf("output should include numbered head lines, got:\n%s", out)
	}
	if !strings.Contains(out, "    12 | line 12\n") {
		t.Fatalf("output should include numbered tail lines, got:\n%s", out)
	}
}

func TestRenderFileCustomHeadAndTail(t *testing.T) {
	path := writeTempFile(t, 8)

	out, err := renderFile(path, options{head: 2, tail: 3, theme: "plain"})
	if err != nil {
		t.Fatalf("renderFile returned error: %v", err)
	}

	want := path + "\n" +
		"     1 | line 01\n" +
		"     2 | line 02\n" +
		separator + "\n" +
		"     6 | line 06\n" +
		"     7 | line 07\n" +
		"     8 | line 08\n"

	if out != want {
		t.Fatalf("unexpected output\nwant:\n%s\ngot:\n%s", want, out)
	}
}

func TestRenderFileWithColor(t *testing.T) {
	path := writeTempFile(t, 1)

	out, err := renderFile(path, options{head: 1, tail: 1, theme: "default", useColor: true})
	if err != nil {
		t.Fatalf("renderFile returned error: %v", err)
	}

	if !strings.Contains(out, "\033[1;36m"+path+"\033[0m") {
		t.Fatalf("filename should be colorized, got:\n%q", out)
	}
	if !strings.Contains(out, "\033[2;32m     1\033[0m | line 01") {
		t.Fatalf("line number should be colorized, got:\n%q", out)
	}
}

func TestParseArgs(t *testing.T) {
	opts, files, err := parseArgs([]string{"-lines", "4", "-tail", "2", "-color", "never", "a.txt"}, os.Stdout)
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}

	if opts.head != 4 || opts.tail != 2 || opts.useColor {
		t.Fatalf("unexpected options: %+v", opts)
	}
	if len(files) != 1 || files[0] != "a.txt" {
		t.Fatalf("unexpected files: %#v", files)
	}
}

func TestParseShortArgs(t *testing.T) {
	opts, files, err := parseArgs([]string{"-l", "6", "-h", "3", "-t", "2", "-c", "always", "-T", "plain", "a.txt"}, os.Stdout)
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}

	if opts.head != 3 || opts.tail != 2 || opts.colorMode != "always" || opts.theme != "plain" {
		t.Fatalf("unexpected options: %+v", opts)
	}
	if opts.useColor {
		t.Fatalf("plain theme should disable color: %+v", opts)
	}
	if len(files) != 1 || files[0] != "a.txt" {
		t.Fatalf("unexpected files: %#v", files)
	}
}

func TestParseHelp(t *testing.T) {
	_, _, err := parseArgs([]string{"-help"}, os.Stdout)
	if err == nil {
		t.Fatal("expected help error")
	}
	if err != flag.ErrHelp {
		t.Fatalf("expected flag.ErrHelp, got %v", err)
	}
	if !strings.Contains(helpText(), "Usage:\n  ,ht [flags] FILE [FILE...]") {
		t.Fatalf("help text should contain usage, got:\n%s", helpText())
	}
}

func writeTempFile(t *testing.T, lines int) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "sample.txt")
	var b strings.Builder
	for i := 1; i <= lines; i++ {
		fmt.Fprintf(&b, "line %02d\n", i)
	}

	if err := os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

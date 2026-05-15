package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const separator = "--------------------"

type options struct {
	head      int
	tail      int
	colorMode string
	theme     string
	useColor  bool
}

type theme struct {
	reset  string
	file   string
	lineNo string
	sep    string
}

var themes = map[string]theme{
	"default": {
		reset:  "\033[0m",
		file:   "\033[1;36m",
		lineNo: "\033[2;32m",
		sep:    "\033[2;33m",
	},
	"plain": {},
}

func main() {
	if len(os.Args) == 1 {
		fmt.Print(helpText())
		return
	}

	opts, files, err := parseArgs(os.Args[1:], os.Stdout)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Print(helpText())
			return
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if len(files) == 0 {
		fmt.Print(helpText())
		return
	}

	hadError := false
	wrote := false
	for _, file := range files {
		out, err := renderFile(file, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, ",ht: %s: %v\n", file, err)
			hadError = true
			continue
		}

		if wrote {
			fmt.Println()
		}
		fmt.Print(out)
		wrote = true
	}

	if hadError {
		os.Exit(1)
	}
}

func parseArgs(args []string, stdout *os.File) (options, []string, error) {
	fs := flag.NewFlagSet(",ht", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	lines := 10
	head := -1
	tail := -1
	colorMode := "auto"
	themeName := "default"

	fs.IntVar(&lines, "lines", lines, "number of lines to show for both head and tail")
	fs.IntVar(&lines, "l", lines, "number of lines to show for both head and tail")
	fs.IntVar(&head, "head", head, "number of head lines to show")
	fs.IntVar(&head, "h", head, "number of head lines to show")
	fs.IntVar(&tail, "tail", tail, "number of tail lines to show")
	fs.IntVar(&tail, "t", tail, "number of tail lines to show")
	fs.StringVar(&colorMode, "color", colorMode, "color output: auto, always, never")
	fs.StringVar(&colorMode, "c", colorMode, "color output: auto, always, never")
	fs.StringVar(&themeName, "theme", themeName, "theme name: default, plain")
	fs.StringVar(&themeName, "T", themeName, "theme name: default, plain")

	if err := fs.Parse(args); err != nil {
		return options{}, nil, err
	}

	if lines < 0 {
		return options{}, nil, errors.New("-lines must be non-negative")
	}
	if head < -1 {
		return options{}, nil, errors.New("-head must be non-negative")
	}
	if tail < -1 {
		return options{}, nil, errors.New("-tail must be non-negative")
	}
	if _, ok := themes[themeName]; !ok {
		return options{}, nil, fmt.Errorf("unknown theme %q", themeName)
	}

	opts := options{
		head:      lines,
		tail:      lines,
		colorMode: colorMode,
		theme:     themeName,
	}
	if head >= 0 {
		opts.head = head
	}
	if tail >= 0 {
		opts.tail = tail
	}

	switch colorMode {
	case "auto":
		opts.useColor = isTerminal(stdout)
	case "always":
		opts.useColor = true
	case "never":
		opts.useColor = false
	default:
		return options{}, nil, fmt.Errorf("unknown color mode %q", colorMode)
	}
	if themeName == "plain" {
		opts.useColor = false
	}

	return opts, fs.Args(), nil
}

func renderFile(path string, opts options) (string, error) {
	lines, err := readLines(path)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	t := themes[opts.theme]

	writeColor(&b, opts, t.file, path)
	b.WriteByte('\n')

	headEnd := min(opts.head, len(lines))
	for i := 0; i < headEnd; i++ {
		writeLine(&b, opts, i+1, lines[i])
	}

	writeColor(&b, opts, t.sep, separator)
	b.WriteByte('\n')

	tailStart := max(len(lines)-opts.tail, 0)
	for i := tailStart; i < len(lines); i++ {
		writeLine(&b, opts, i+1, lines[i])
	}

	return b.String(), nil
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func writeLine(b *strings.Builder, opts options, number int, text string) {
	t := themes[opts.theme]
	if opts.useColor {
		fmt.Fprintf(b, "%s%6d%s | %s\n", t.lineNo, number, t.reset, text)
		return
	}
	fmt.Fprintf(b, "%6d | %s\n", number, text)
}

func writeColor(b *strings.Builder, opts options, code string, text string) {
	t := themes[opts.theme]
	if opts.useColor && code != "" {
		b.WriteString(code)
		b.WriteString(text)
		b.WriteString(t.reset)
		return
	}
	b.WriteString(text)
}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func helpText() string {
	return `Usage:
  ,ht [flags] FILE [FILE...]

Print the first and last lines of each file with line numbers.
Each file section starts with the filename and separates head/tail output
with a 20-dash line.

Flags:
  -lines, -l N        number of lines to show for both head and tail (default 10)
  -head,  -h N        number of head lines to show, overriding -lines
  -tail,  -t N        number of tail lines to show, overriding -lines
  -color, -c MODE     color output: auto, always, never (default auto)
  -theme, -T NAME     theme name: default, plain (default default)
  -help               show this help output

Examples:
  ,ht notes.txt
  ,ht -l 5 notes.txt report.txt
  ,ht -h 20 -t 5 -c always app.log
  ,ht -T plain app.log
`
}

# ,ht

`,ht` prints the first and last lines of each file with line numbers, a 20-dash separator, and optional ANSI color.

## Usage

```sh
,ht [flags] FILE [FILE...]
```

Multiple files are separated by a blank line. Each file section starts with the filename.
Calling `,ht` with no arguments, or with `-help`, prints help output.

## Flags

```text
-lines, -l N        number of lines to show for both head and tail (default 10)
-head, -h N         number of head lines to show, overriding -lines
-tail, -t N         number of tail lines to show, overriding -lines
-color, -c MODE     color output: auto, always, never (default auto)
-theme, -T NAME     theme name: default, plain (default default)
```

Examples:

```sh
,ht notes.txt
,ht -l 5 notes.txt report.txt
,ht -h 20 -t 5 -c always app.log
,ht -T plain app.log
```

## Build

Build the native binary for your current platform:

```sh
make build
```

The binary is written to:

```text
bin/,ht
```

On macOS, you can put that directory on your `PATH`, or copy/link `bin/,ht` into a personal bin directory that is already on your `PATH`.

## Test

Run the test suite:

```sh
make test
```

## Development

Run without building:

```sh
make run ARGS="-lines 3 README.md"
```

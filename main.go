package main

import (
"bufio"
"flag"
"fmt"
"io"
"os"
"regexp"

"sigs.k8s.io/yaml"
)

// flatten converts a YAML structure into a flat list of key-value paths
// similar to how gron works for JSON
func flatten(data interface{}, prefix string, result *[]string) {
switch v := data.(type) {
case map[string]interface{}:
for key, val := range v {
newPrefix := key
if prefix != "" {
newPrefix = prefix + "." + key
}
flatten(val, newPrefix, result)
}
case []interface{}:
for i, val := range v {
newPrefix := fmt.Sprintf("%s[%d]", prefix, i)
flatten(val, newPrefix, result)
}
default:
line := fmt.Sprintf("%s = %v", prefix, v)
*result = append(*result, line)
}
}

func main() {
// Define command-line flags
filterPattern := flag.String("filter", "", "regex pattern to filter output")
flag.Usage = func() {
fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [FILE]\n", os.Args[0])
fmt.Fprintf(os.Stderr, "\nkygrep - KYAML grep tool\n")
fmt.Fprintf(os.Stderr, "Converts YAML to flat key-value format and optionally filters with regex\n\n")
fmt.Fprintf(os.Stderr, "Options:\n")
flag.PrintDefaults()
fmt.Fprintf(os.Stderr, "\nIf FILE is not specified, reads from stdin\n")
fmt.Fprintf(os.Stderr, "\nExamples:\n")
fmt.Fprintf(os.Stderr, "  kygrep config.yaml\n")
fmt.Fprintf(os.Stderr, "  kygrep -filter 'api.*key' config.yaml\n")
fmt.Fprintf(os.Stderr, "  cat config.yaml | kygrep -filter 'port'\n")
}
flag.Parse()

// Determine input source
var reader io.Reader
if flag.NArg() > 0 {
file, err := os.Open(flag.Arg(0))
if err != nil {
fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
os.Exit(1)
}
defer file.Close()
reader = file
} else {
reader = os.Stdin
}

// Read input
input, err := io.ReadAll(reader)
if err != nil {
fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
os.Exit(1)
}

// Parse YAML
var data interface{}
err = yaml.Unmarshal(input, &data)
if err != nil {
fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
os.Exit(1)
}

// Flatten the structure
var lines []string
flatten(data, "", &lines)

// Compile regex if filter is provided
var filterRegex *regexp.Regexp
if *filterPattern != "" {
filterRegex, err = regexp.Compile(*filterPattern)
if err != nil {
fmt.Fprintf(os.Stderr, "Error compiling regex: %v\n", err)
os.Exit(1)
}
}

// Output lines, optionally filtered
writer := bufio.NewWriter(os.Stdout)
defer writer.Flush()

for _, line := range lines {
if filterRegex == nil || filterRegex.MatchString(line) {
fmt.Fprintln(writer, line)
}
}
}

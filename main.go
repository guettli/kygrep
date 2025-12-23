package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

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

// unflatten converts flattened lines back to a YAML structure
func unflatten(lines []string) (interface{}, error) {
	root := make(map[string]interface{})

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse line: "path = value"
		parts := strings.SplitN(line, " = ", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}

		path := parts[0]
		value := parseValue(parts[1])

		// Parse path and set value
		setValueAtPath(root, path, value)
	}

	return root, nil
}

// parseValue converts string value to appropriate type
func parseValue(s string) interface{} {
	// Try to parse as int
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	// Try to parse as float
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	// Try to parse as bool
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	// Return as string
	return s
}

// setValueAtPath sets a value in a nested map/slice structure based on path
func setValueAtPath(root map[string]interface{}, path string, value interface{}) {
	// Parse path into tokens
	tokens := parsePath(path)
	if len(tokens) == 0 {
		return
	}

	current := interface{}(root)

	for i := 0; i < len(tokens)-1; i++ {
		token := tokens[i]
		nextToken := tokens[i+1]

		switch t := token.(type) {
		case string:
			// Current is a map
			m, ok := current.(map[string]interface{})
			if !ok {
				return
			}

			// Determine if next level should be map or array
			if _, isIndex := nextToken.(int); isIndex {
				// Next is array
				if _, exists := m[t]; !exists {
					m[t] = []interface{}{}
				}
				current = m[t]
			} else {
				// Next is map
				if _, exists := m[t]; !exists {
					m[t] = make(map[string]interface{})
				}
				current = m[t]
			}

		case int:
			// Current is an array
			arr, ok := current.([]interface{})
			if !ok {
				return
			}

			// Expand array if needed
			for len(arr) <= t {
				arr = append(arr, nil)
			}

			// Determine if next level should be map or array
			if _, isIndex := nextToken.(int); isIndex {
				if arr[t] == nil {
					arr[t] = []interface{}{}
				}
			} else {
				if arr[t] == nil {
					arr[t] = make(map[string]interface{})
				}
			}

			// Update the parent reference
			updateArrayInParent(root, tokens[:i], arr)
			current = arr[t]
		}
	}

	// Set final value
	lastToken := tokens[len(tokens)-1]
	switch t := lastToken.(type) {
	case string:
		if m, ok := current.(map[string]interface{}); ok {
			m[t] = value
		}
	case int:
		if arr, ok := current.([]interface{}); ok {
			for len(arr) <= t {
				arr = append(arr, nil)
			}
			arr[t] = value
			updateArrayInParent(root, tokens[:len(tokens)-1], arr)
		}
	}
}

// updateArrayInParent updates an array value in the parent structure
func updateArrayInParent(root map[string]interface{}, tokens []interface{}, arr []interface{}) {
	if len(tokens) == 0 {
		return
	}

	current := interface{}(root)
	for i := 0; i < len(tokens)-1; i++ {
		token := tokens[i]
		switch t := token.(type) {
		case string:
			if m, ok := current.(map[string]interface{}); ok {
				current = m[t]
			}
		case int:
			if a, ok := current.([]interface{}); ok {
				current = a[t]
			}
		}
	}

	lastToken := tokens[len(tokens)-1]
	switch t := lastToken.(type) {
	case string:
		if m, ok := current.(map[string]interface{}); ok {
			m[t] = arr
		}
	case int:
		if a, ok := current.([]interface{}); ok {
			a[t] = arr
		}
	}
}

// parsePath converts a path string into tokens (strings and ints)
func parsePath(path string) []interface{} {
	var tokens []interface{}
	var current strings.Builder

	i := 0
	for i < len(path) {
		if path[i] == '[' {
			// Flush current token
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

			// Find closing bracket
			j := i + 1
			for j < len(path) && path[j] != ']' {
				j++
			}

			if j < len(path) {
				// Parse index
				indexStr := path[i+1 : j]
				if index, err := strconv.Atoi(indexStr); err == nil {
					tokens = append(tokens, index)
				}
				i = j + 1
			} else {
				i++
			}
		} else if path[i] == '.' {
			// Flush current token
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			i++
		} else {
			current.WriteByte(path[i])
			i++
		}
	}

	// Flush final token
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func main() {
	// Define command-line flags
	unflattenMode := flag.Bool("u", false, "unflatten mode: convert flattened format back to YAML")
	filterPattern := flag.String("filter", "", "regex pattern to filter output (flatten mode only)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [FILE]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nkygrep - KYAML grep tool\n")
		fmt.Fprintf(os.Stderr, "Converts YAML to flat key-value format and optionally filters with regex\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nIf FILE is not specified, reads from stdin\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  kygrep config.yaml              # Flatten YAML\n")
		fmt.Fprintf(os.Stderr, "  kygrep -filter 'port' config.yaml | kygrep -u  # Filter and unflatten\n")
		fmt.Fprintf(os.Stderr, "  cat config.yaml | kygrep | grep metadata | kygrep -u\n")
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

	if *unflattenMode {
		// Unflatten mode: read flattened format and output YAML
		scanner := bufio.NewScanner(reader)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		data, err := unflatten(lines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unflattening: %v\n", err)
			os.Exit(1)
		}

		output, err := yaml.Marshal(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting to YAML: %v\n", err)
			os.Exit(1)
		}

		fmt.Print(string(output))
	} else {
		// Flatten mode: read YAML and output flattened format
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
}

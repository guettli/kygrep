# kygrep

KYAML based grep - A tool to flatten YAML files and filter them with regex patterns, similar to [gron](https://github.com/tomnomnom/gron) for JSON.

## Overview

`kygrep` reads YAML files and converts them into a flat key-value format, making it easy to search and filter specific paths using regular expressions. It uses the [kyaml](https://pkg.go.dev/sigs.k8s.io/yaml) library for YAML processing.

## Installation

### Build from source

```bash
./build.sh
```

This will create the `kygrep` binary in the current directory.

## Usage

```bash
kygrep [OPTIONS] [FILE]
```

If FILE is not specified, kygrep reads from stdin.

### Options

- `-filter string`: Regular expression pattern to filter output

### Examples

#### Basic usage - flatten YAML file

```bash
kygrep config.yaml
```

Example output:
```
apiVersion = v1
kind = Service
metadata.name = my-service
metadata.namespace = default
spec.ports[0].port = 80
spec.ports[0].protocol = TCP
spec.ports[0].targetPort = 9376
spec.selector.app = MyApp
spec.type = LoadBalancer
```

#### Filter by regex pattern

```bash
kygrep -filter 'port' config.yaml
```

Output:
```
spec.ports[0].port = 80
spec.ports[0].protocol = TCP
spec.ports[0].targetPort = 9376
```

#### Filter nested paths

```bash
kygrep -filter 'metadata\.' config.yaml
```

Output:
```
metadata.name = my-service
metadata.namespace = default
```

#### Use with stdin

```bash
cat config.yaml | kygrep -filter 'selector'
```

Output:
```
spec.selector.app = MyApp
```

#### Complex regex patterns

```bash
kygrep -filter '^spec\.ports.*port.*=' config.yaml
```

## Testing

Run the test suite:

```bash
./test.sh
```

Or directly with Go:

```bash
go test -v
```

## How it works

1. Reads YAML input (from file or stdin)
2. Parses YAML using the kyaml library
3. Recursively flattens the structure into key-value pairs
4. Optionally filters lines matching the provided regex pattern
5. Outputs the results

## Comparison to gron

Like `gron` does for JSON, `kygrep` makes YAML greppable by converting it into a flat format. This makes it easy to:

- Search for specific configuration values
- Filter complex nested structures
- Pipe output to other Unix tools
- Quickly understand YAML structure

## License

See LICENSE file for details.

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

- `-filter string`: Regular expression pattern to filter output (flatten mode only)
- `-u`: Unflatten mode - convert flattened format back to YAML

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

#### Unflatten back to YAML

Like `gron -u`, you can convert the flattened format back to YAML:

```bash
kygrep config.yaml | kygrep -u
```

This will output the original YAML structure.

#### Filter and unflatten (gron-style workflow)

The most powerful usage is the gron-style pipeline: flatten → filter → unflatten

```bash
# Filter only metadata fields and convert back to YAML
kygrep config.yaml | grep metadata | kygrep -u
```

Output:
```yaml
metadata:
  name: my-service
  namespace: default
```

```bash
# Filter with built-in regex and unflatten
kygrep -filter 'port' config.yaml | kygrep -u
```

Output:
```yaml
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 9376
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
cat config.yaml | kygrep | grep selector | kygrep -u
```

Output:
```yaml
spec:
  selector:
    app: MyApp
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

### Flatten mode (default)

1. Reads YAML input (from file or stdin)
2. Parses YAML using the kyaml library
3. Recursively flattens the structure into key-value pairs
4. Optionally filters lines matching the provided regex pattern
5. Outputs the flattened format

### Unflatten mode (`-u` flag)

1. Reads flattened key-value format (from file or stdin)
2. Parses each line and reconstructs the nested structure
3. Converts back to YAML format
4. Outputs the YAML

This allows the gron-style workflow: `kygrep config.yaml | grep pattern | kygrep -u`

## Comparison to gron

Like `gron` does for JSON, `kygrep` makes YAML greppable by converting it into a flat format. This makes it easy to:

- Search for specific configuration values
- Filter complex nested structures with standard tools like `grep`
- Use the gron-style pipeline: flatten → filter → unflatten
- Pipe output to other Unix tools
- Quickly understand YAML structure

**Key difference**: While gron uses JavaScript assignment syntax, kygrep uses a simple `path = value` format.

## License

See LICENSE file for details.

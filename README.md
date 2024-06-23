<img src="assets/logo.svg" alt="Render Kit Logo" width="150px" align="right" />

**Render Kit** is a versatile and powerful command-line interface (CLI) tool designed for comprehensive template rendering. It supports multiple template engines and data sources, providing both flexibility and efficiency.

[![CodeQL](https://github.com/orellazri/renderkit/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/orellazri/renderkit/actions/workflows/github-code-scanning/codeql)
[![Test](https://github.com/orellazri/renderkit/actions/workflows/test.yml/badge.svg)](https://github.com/orellazri/renderkit/actions/workflows/test.yml)

## Features

- 🛠️ Supports multiple template engines
- 🌐 Integrates with various data sources
- 🎛️ Customizable rendering options
- ⚡ Lightweight and fast
- 🌍 Cross-platform compatibility
- 📦 Single binary

### Supported Engines

- Envsubst
- Go Templates (including [Sprig Functions](http://masterminds.github.io/sprig/))
- Handlebars
- Jet
- Jinja
- Mustache

### Supported Datasources

- Environment variables
- YAML
- JSON
- TOML

## Usage

To use Render Kit, you have multiple options:

1. **Download the latest release binary**:

   - Visit the [releases page](https://github.com/orellazri/renderkit/releases) and download the latest binary for your operating system. It's recommended to move the binary to a directory in your `PATH` to make it easier to run such as `/usr/local/bin`.

1. **Install via Go**:

   - Ensure that you have Go installed on your machine.
   - Run the following command:
     ```bash
     go install github.com/orellazri/renderkit@latest
     ```

1. **Run the Docker image**:

   - If you prefer using Docker, you can run the `reaperberri/renderkit` Docker image.
   - Make sure you have Docker installed on your machine.
   - Run the following command:
     ```bash
     docker run --rm reaperberri/renderkit <args>
     ```

You need to run the `renderkit` command with the following arguments as either command-line flags, or as a YAML configuration file passed via `--config`.

| Name                   | Description                                                                    | Type   |
| ---------------------- | ------------------------------------------------------------------------------ | ------ |
| `config`               | Load configuration from YAML file                                              | string |
| `input`                | Template string to render                                                      | string |
| `input-file`           | Template input file to render                                                  | string |
| `input-dir`            | Template input directory to render                                             | string |
| `exclude`              | Glob patterns for files to exclude from rendering                              | list   |
| `output`               | Output directory to write to                                                   | string |
| `datasource`           | Datasource to use for rendering (scheme://path)                                | list   |
| `data`                 | Data to use for rendering. Can be used to provide data directly                | list   |
| `engine`               | Templating engine to use for rendering (Go Templates by default)               | string |
| `allow-duplicate-keys` | Allow duplicate keys in datasources. If set, the last value found will be used | bool   |

Below are practical examples demonstrating the usage of `renderkit`:
```bash

# Using a specific env var as a datasource
$ cat ds.yml
FN: "Doe"
$ echo 'Hello {{.FN}} {{.LN}}' | renderkit -ds env://LN -ds ds.yml
Hello John Doe

# Using a template string and envsubst engine
$ export LN="Doe"
$ echo 'Hello $FN $LN' | renderkit -i 'Hello $FN $LN' -e envsubst --data "FN=John"
Hello John Doe

# Using a template file with data from a JSON file
$ cat data.json
{
  "names": {
    "FN": "John",
    "LN": "Doe"
  }
}
$ cat file.tpl
Hello {{ lower .names.FN }} {{ upper .names.LN }}
$ renderkit -f file.tpl -ds data.json
Hello john DOE

# Render input directory [1.tpl, 2.tpl, 3.tpl] to output directory
$ renderkit --input-dir in/ --exclude 'in/[1-2].tpl' --output out/ --datasource data.yml --data myKey=myValue --engine jinja
# Output directory will contain [3.tpl] rendered files

```

### Example YAML Configuration File

```yaml
input-dir: input/
output: output/
exclude:
  - input/exclude[1-2].tpl
  - input/other_*.tpl
datasource:
  - data.yaml
  - data2.json
engine: gotemplates
allow-duplicate-keys: true
```

## Development

### Prerequisites

- [Task](https://taskfile.dev/)
- [pre-commit](https://pre-commit.com/)
- [golangci-lint](https://github.com/golangci/golangci-lint)

### Running locally

1.  Fork and clone the repository
1.  Install pre-commit hooks:

    ```bash
    pre-commit install
    ```

1.  Run with:

    ```bash
    go run .
    ```

### Running tests

```bash
task test             # Run all tests (including integration)
task test SHORT=true  # Run only unit tests
```

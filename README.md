<div align="center">

<img src="assets/logo.svg" alt="Render Kit Logo" width="150px" align="right" />

<h1>Render Kit</h1>

A versatile and powerful command-line interface (CLI) tool designed for comprehensive template rendering. It supports multiple template engines and data sources, providing both flexibility and efficiency.

</div>

## Features

- 🛠️ Supports multiple template engines
- 🌐 Integrates with various data sources
- 🎛️ Customizable rendering options
- ⚡ Lightweight and fast
- 🌍 Cross-platform compatibility

### Supported Engines

- Jinja
- Mustache
- Jet

### Supported Datasources

- YAML
- JSON

## Usage

Download the latest release from the [releases page](https://github.com/orellazri/renderkit/releases).

Run the `renderkit` command with the following arguments as either command-line flags, or as a YAML configuration file passed with `--config`.

| Name                   | Description                                                                    | Type   |
| ---------------------- | ------------------------------------------------------------------------------ | ------ |
| `config`               | Load configuration from YAML file                                              | string |
| `input`                | The input file to render                                                       | list   |
| `output`               | The output file to write to                                                    | string |
| `input-dir`            | The input directory to render                                                  | string |
| `output-dir`           | The output directory to write to                                               | string |
| `datasource`           | The datasource to use for rendering (scheme://path)                            | list   |
| `engine`               | The templating engine to use for rendering                                     | string |
| `allow-duplicate-keys` | Allow duplicate keys in datasources. If set, the last value found will be used | bool   |

This tool supports 3 modes:

#### File to file

Render one input file to one output file.

```bash
renderkit --input input.tpl --output output.txt --datasource data.yaml --engine jinja
```

#### Files to directory

Render multiple input files to one output directory.

```bash
renderkit --input input1.tpl --input input2.tpl --output-dir output --datasource data.yaml --engine jinja
```

#### Directory to directory

Render all files in an input directory to an output directory.

```bash
renderkit --input-dir input --output-dir output --datasource data.yaml --engine jinja
```

## Development

### Prerequisites

- [Task](https://taskfile.dev/)
- [pre-commit](https://pre-commit.com/)
- [golangci-lint](https://github.com/golangci/golangci-lint)

```bash
pre-commit install
go run cmd/renderkit/main.go
```

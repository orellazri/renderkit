<img src="assets/logo.svg" alt="Render Kit Logo" width="150px" align="right" />

**Render Kit** is a versatile and powerful command-line interface (CLI) tool designed for comprehensive template rendering. It supports multiple template engines and data sources, providing both flexibility and efficiency.

## Features

- 🛠️ Supports multiple template engines
- 🌐 Integrates with various data sources
- 🎛️ Customizable rendering options
- ⚡ Lightweight and fast
- 🌍 Cross-platform compatibility
- 📦 Single binary

### Supported Engines

- Go Templates (with [Sprig Functions](http://masterminds.github.io/sprig/))
- Jinja
- Handlebars
- Mustache
- Jet

### Supported Datasources

- Environment variables
- YAML
- JSON
- TOML

## Usage

Download the latest release from the [releases page](https://github.com/orellazri/renderkit/releases) or use the `reaperberri/renderkit` Docker image.

Run the `renderkit` command with the following arguments as either command-line flags, or as a YAML configuration file passed with `--config`.

| Name                   | Description                                                                    | Type   | Required |
| ---------------------- | ------------------------------------------------------------------------------ | ------ | -------- |
| `config`               | Load configuration from YAML file                                              | string | False    |
| `input`                | The input glob to render                                                       | string | True     |
| `exclude`              | The glob pattern for files to exclude from rendering                           | string | False    |
| `output`               | The output directory to write to                                               | string | True     |
| `datasource`           | The datasource to use for rendering (scheme://path)                            | list   | True     |
| `data`                 | The data to use for rendering. Can be used to provide data directly            | list   | False    |
| `engine`               | The templating engine to use for rendering                                     | string | True     |
| `allow-duplicate-keys` | Allow duplicate keys in datasources. If set, the last value found will be used | bool   | False    |

```bash
renderkit --input in/*.tpl --output output/ --datasource data.yaml --data myKey=myValue --engine jinja
```

### Example YAML Configuration File

```yaml
input: in/*.tpl
output: output/
datasource:
  - data.yaml
  - data2.json
engine: jinja
allow-duplicate-keys: true
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

# logslice

A fast CLI tool for filtering and slicing structured log files by time range or field patterns.

## Installation

```bash
go install github.com/yourname/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logslice.git
cd logslice
go build -o logslice .
```

## Usage

```bash
# Filter logs by time range
logslice --from "2024-01-15T10:00:00Z" --to "2024-01-15T11:00:00Z" app.log

# Filter by field pattern
logslice --field "level=error" app.log

# Combine time range and field filter
logslice --from "2024-01-15T10:00:00Z" --field "service=api" app.log

# Read from stdin
cat app.log | logslice --field "status=500"
```

### Flags

| Flag | Description |
|------|-------------|
| `--from` | Start of time range (RFC3339) |
| `--to` | End of time range (RFC3339) |
| `--field` | Field pattern in `key=value` format |
| `--format` | Log format: `json`, `logfmt` (default: `json`) |
| `--output` | Output file (default: stdout) |

## Requirements

- Go 1.21+

## License

MIT © [yourname](https://github.com/yourname)
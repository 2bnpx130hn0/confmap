# confmap

A Go library for merging layered config files (YAML/TOML/env) with override precedence and schema validation.

---

## Installation

```bash
go get github.com/yourorg/confmap
```

---

## Usage

```go
package main

import (
    "fmt"
    "github.com/yourorg/confmap"
)

func main() {
    cfg, err := confmap.Load(
        confmap.WithFile("config.yaml"),
        confmap.WithFile("config.local.toml"),
        confmap.WithEnv("APP_"),
    )
    if err != nil {
        panic(err)
    }

    host := cfg.GetString("server.host")
    port := cfg.GetInt("server.port")

    fmt.Printf("Listening on %s:%d\n", host, port)
}
```

Layers are merged in order — later sources take precedence. Environment variables (prefixed with `APP_`) always win.

### Schema Validation

```go
cfg, err := confmap.Load(
    confmap.WithFile("config.yaml"),
    confmap.WithSchema("schema.json"),
)
```

Validation runs automatically after merging. An error is returned if the final config does not satisfy the schema.

---

## Supported Formats

| Format | Extension |
|--------|-----------|
| YAML   | `.yaml`, `.yml` |
| TOML   | `.toml` |
| Env    | prefix-based |

---

## License

MIT © yourorg
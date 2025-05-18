# Configuration Guide

Dora can be configured using a YAML configuration file. This guide explains the available options.

## Configuration File Location

By default, Dora looks for a configuration file at `config/dora.config.yml`. You can specify a different location using the `--config` flag.

## Main Configuration Sections

### Server Configuration

```yaml
server:
  host: "0.0.0.0"  # Bind address
  port: 8080         # Port to listen on
  tls: false         # Enable HTTPS
  tls_cert: ""       # Path to TLS certificate
  tls_key: ""        # Path to TLS private key
```

### Database Configuration

```yaml
database:
  host: "localhost"
  port: 5432
  name: "dora"
  user: "dora"
  password: ""  # Consider using environment variables for secrets
  sslmode: "disable"
```

### Cache Configuration

```yaml
redis:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

### Ethereum Node Configuration

```yaml
eth1:
  endpoint: "http://localhost:8545"
  ws_endpoint: "ws://localhost:8546"

eth2:
  endpoint: "http://localhost:5052"
```

## Environment Variables

All configuration options can be set using environment variables. For example:

```
DORA_SERVER_PORT=8080
DORA_DATABASE_HOST=localhost
```

## Configuration Precedence

1. Command-line flags
2. Environment variables
3. Configuration file
4. Default values

## Next Steps

- [Installation Guide](Installation)
- [Quick Start Guide](Quick-Start)

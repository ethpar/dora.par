# Installation Guide

This guide will help you install Dora Block Explorer on your system.

## Prerequisites

- Go 1.20 or later
- PostgreSQL 13+
- Redis 6.0+
- Git

## Installation Steps

### 1. Clone the Repository

```bash
git clone https://github.com/ethpar/dora.par.git
cd dora.par
```

### 2. Build from Source

```bash
make build
```

### 3. Install Dependencies

Install PostgreSQL and Redis following their respective documentation.

### 4. Configure Dora

Copy the example configuration file:

```bash
cp config/default.config.yml config/dora.config.yml
```

Edit the configuration file to match your environment.

### 5. Initialize the Database

```bash
./dora db migrate
```

### 6. Start Dora

```bash
./dora start
```

## Next Steps

- [Configuration](Configuration)
- [Quick Start Guide](Quick-Start)

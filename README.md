# vaultpipe

> Sync secrets from HashiCorp Vault into local `.env` files with namespace filtering.

---

## Installation

```bash
go install github.com/yourname/vaultpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/vaultpipe.git
cd vaultpipe
go build -o vaultpipe .
```

---

## Usage

```bash
vaultpipe --addr https://vault.example.com \
          --token s.yourVaultToken \
          --namespace secret/myapp \
          --output .env
```

This will pull all secrets under the `secret/myapp` namespace and write them to a `.env` file in `KEY=VALUE` format.

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--addr` | Vault server address | `http://127.0.0.1:8200` |
| `--token` | Vault authentication token | `$VAULT_TOKEN` |
| `--namespace` | Secret namespace/path to sync | required |
| `--output` | Output `.env` file path | `.env` |
| `--overwrite` | Overwrite existing `.env` file | `false` |

### Example Output

```env
DB_HOST=postgres.internal
DB_PASSWORD=supersecret
API_KEY=abc123
```

---

## Requirements

- Go 1.21+
- HashiCorp Vault with KV v2 secrets engine

---

## License

MIT © [yourname](https://github.com/yourname)
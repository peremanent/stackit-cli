## stackit server backup disable

Disables Server Backup service

### Synopsis

Disables Server Backup service.

```
stackit server backup disable [flags]
```

### Examples

```
  Disable Server Backup functionality for your server.
  $ stackit server backup disable --server-id=zzz
```

### Options

```
  -h, --help               Help for "stackit server backup disable"
  -s, --server-id string   Server ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit server backup](./stackit_server_backup.md)	 - Provides functionality for server backups


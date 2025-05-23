## stackit server os-update disable

Disables server os-update service

### Synopsis

Disables server os-update service.

```
stackit server os-update disable [flags]
```

### Examples

```
  Disable os-update functionality for your server.
  $ stackit server os-update disable --server-id=zzz
```

### Options

```
  -h, --help               Help for "stackit server os-update disable"
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

* [stackit server os-update](./stackit_server_os-update.md)	 - Provides functionality for managed server updates


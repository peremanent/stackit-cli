## stackit redis instance describe

Shows details  of a Redis instance

### Synopsis

Shows details  of a Redis instance.

```
stackit redis instance describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of a Redis instance with ID "xxx"
  $ stackit redis instance describe xxx

  Get details of a Redis instance with ID "xxx" in JSON format
  $ stackit redis instance describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit redis instance describe"
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

* [stackit redis instance](./stackit_redis_instance.md)	 - Provides functionality for Redis instances


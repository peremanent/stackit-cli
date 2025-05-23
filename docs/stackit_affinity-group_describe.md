## stackit affinity-group describe

Show details of an affinity group

### Synopsis

Show details of an affinity group.

```
stackit affinity-group describe AFFINITY_GROUP_ID [flags]
```

### Examples

```
  Get details about an affinity group with the ID "xxx"
  $ stackit affinity-group describe xxx
```

### Options

```
  -h, --help   Help for "stackit affinity-group describe"
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

* [stackit affinity-group](./stackit_affinity-group.md)	 - Manage server affinity groups


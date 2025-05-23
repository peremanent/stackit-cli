## stackit network-area network-range describe

Shows details of a network range in a STACKIT Network Area (SNA)

### Synopsis

Shows details of a network range in a STACKIT Network Area (SNA).

```
stackit network-area network-range describe NETWORK_RANGE_ID [flags]
```

### Examples

```
  Show details of a network range with ID "xxx" in a STACKIT Network Area with ID "yyy" in organization with ID "zzz"
  $ stackit network-area network-range describe xxx --network-area-id yyy --organization-id zzz
```

### Options

```
  -h, --help                     Help for "stackit network-area network-range describe"
      --network-area-id string   STACKIT Network Area (SNA) ID
      --organization-id string   Organization ID
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

* [stackit network-area network-range](./stackit_network-area_network-range.md)	 - Provides functionality for network ranges in STACKIT Network Areas


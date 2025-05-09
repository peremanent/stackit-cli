## stackit public-ip disassociate

Disassociates a Public IP from a network interface or a virtual IP

### Synopsis

Disassociates a Public IP from a network interface or a virtual IP.

```
stackit public-ip disassociate PUBLIC_IP_ID [flags]
```

### Examples

```
  Disassociate public IP with ID "xxx" from a resource (network interface or virtual IP)
  $ stackit public-ip disassociate xxx
```

### Options

```
  -h, --help   Help for "stackit public-ip disassociate"
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

* [stackit public-ip](./stackit_public-ip.md)	 - Provides functionality for public IPs


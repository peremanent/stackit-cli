## stackit organization role list

Lists roles and permissions of an organization

### Synopsis

Lists roles and permissions of an organization.

```
stackit organization role list [flags]
```

### Examples

```
  List all roles and permissions of an organization
  $ stackit organization role list --organization-id xxx

  List all roles and permissions of an organization in JSON format
  $ stackit organization role list --organization-id xxx --output-format json

  List up to 10 roles and permissions of an organization
  $ stackit organization role list --organization-id xxx --limit 10
```

### Options

```
  -h, --help                     Help for "stackit organization role list"
      --limit int                Maximum number of entries to list
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

* [stackit organization role](./stackit_organization_role.md)	 - Manages organization roles


## stackit image delete

Deletes an image

### Synopsis

Deletes an image by its internal ID.

```
stackit image delete IMAGE_ID [flags]
```

### Examples

```
  Delete an image with ID "xxx"
  $ stackit image delete xxx
```

### Options

```
  -h, --help   Help for "stackit image delete"
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

* [stackit image](./stackit_image.md)	 - Manage server images


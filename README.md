# go-with-envs-from-files

Simple tool which loads environment variable values from files

## Usage

```
go-with-envs-from-files command [args...]
```

Runs command passing environment variables and given args. If evironment variables with `_FILE` suffix exist - will treat those variables as filepaths. Will additionaly pass to the child process enviroment variables with names without the `_FILE` suffix and values read from the given filepaths.

System signals are passed from the parent to the child process.

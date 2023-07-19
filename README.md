# xk6-csv

A k6 extension enables k6 tests to copy and move files.

Built for [k6](https://go.k6.io/k6) using [xk6](https://github.com/grafana/xk6).

## Usage

Import an entire module's contents:
```JavaScript
import ingest from "k6/x/ingest";
```

## API

Functions:
- ingest.copy(src,dst) # copies the file from src to dst
- ingest.rename(oldName,newName) # renames the file
- ingest.makeDir(dir) # makes the dir
- ingest.makeDirAll(dir) # makes the dir and the sub dirs recursively

For complete API documentation click [here](docs/README.md)!

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Install `xk6`:
  ```bash
  $ go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```bash
  $ xk6 build --with github.com/szkiba/xk6-csv@latest
  ```
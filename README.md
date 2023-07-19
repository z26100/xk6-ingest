# xk6-ingest

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
- git
- docker
- GNU Makefile

```bash
$ make build
```

This command will compile k6 cli tool including the extension for Linux and Windows stored in the bin folder.


You can also build the latest version by hand

```bash
$ docker run --rm -it -u "$(id -u):$(id -g)" -v "${PWD}/bin:/xk6" grafana/xk6 build latest --with github.com/z26100/xk6-ingest@latest
```

This command will compile k6 for your current OS and saves it as k6 in your local directory.
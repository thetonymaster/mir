# framework

## Usage

### Setup Project

```bash
$ git clone https://github.com/thetonymaster/framework.git
$ make setup

- or -

$ go get github.com/thetonymaster/framework
 
```

### Run the demo

```bash
make demo

```

## make Reference

- `setup`: sets the project without having to setup the GOPATH
- `test`: run the tests 
- `cover`: aggregates the coverage of all tests over all packages
- `format`: runs goimports

### CI Mode

CI mode is enabled if the environment variable `CI` is set to `1`.

The `make test` full verbose output is sent to stdout/stderr.


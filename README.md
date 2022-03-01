# Heighliner

## Quickstart

Build client binary:

```shell
make hln
export PATH="$PWD/bin:$PATH"
```

Make priject dir

```shell
mkdir hlnstack && cd hlnstack
```

List all heighliner stacks:

```shell
hln stack list
```

Choose a stack and create a project

```shell
hln new -s=sample
```

List all input values:

```shell
hln input list
```

Spin up your application

```shell
hln up
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

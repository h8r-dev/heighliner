# Heighliner

## Quickstart

Build client binary:

```shell
make hln
export PATH="$PWD/bin:$PATH"
```

cd to a workdir

List all heighliner stacks:

```shell
hln stack list
```

Choose a stack and create a project

```shell
hln new -s=<stack>
```

List all input values:

```shell
hln input list
```

Input a value:
```shell
hln input <type> <name> <value>
```

Up

```shell
hln up
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

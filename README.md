# Heighliner

Heighliner is a cloud native application development platform.
It encapsulates low-level infrastructure details and let developers focus on writing business code.
It provides great developer experience and all the advantanges of cloud-native technologies:
platform-agnostic, multi-cloud architecture, fast evolving community.

It is also built in a modular approach and you can extend it with more developer services.

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

Config a input value

```shell
hln input text hello.message "Hello heighliner"
```

Spin up your application

```shell
hln up
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

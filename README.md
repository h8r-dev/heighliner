# Heighliner

Heighliner is a cloud native application development platform.
It encapsulates low-level infrastructure details and let developers focus on writing business code.
It provides great developer experience and all the advantanges of cloud-native technologies:
platform-agnostic, multi-cloud architecture, fast evolving community.

It is also built in a modular approach and you can extend it with more developer services.

## Build

```shell
make hln
export PATH="$PWD/bin:$PATH"
```

## Getting Started

Check out the [documentation](https://heighliner.dev/docs/getting_started/first_app) on how to start using heighliner.

## Test stacks

```shell
hln -s /path/to/your/stack -p ./relative/path/to/your/plan test
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

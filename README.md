# Heighliner

Heighliner is a cloud native application development platform.
It encapsulates low-level infrastructure details and let developers focus on writing business code.
It provides great developer experience and all the advantanges of cloud-native technologies:
platform-agnostic, multi-cloud architecture, fast evolving community.

It is also built in a modular approach and you can extend it with more developer services.

## Getting Started

Check out the [documentation](https://heighliner.dev/docs/getting_started/installation) on how to start using heighliner.

## Build from source

We recomend install the stable [releases](https://github.com/h8r-dev/heighliner/releases) of heighliner. But if you want to build heighliner from source code:

```shell
git clone git@github.com:h8r-dev/heighliner.git && cd heighliner
make hln
```

Then check the version:
```
export PATH="$PWD/bin:$PATH"
hln version
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

# Heighliner

Heighliner(/’haɪlaɪnər/) is a modern developer tool that deliver your application stack as code. You can codify low level details into human-readable configuration files that you can version, reuse, and share. You can even import existing stacks to build more advanced stacks.

We provide and maintain official stacks to provide out-of-the-box experience for common use cases. Your development environment can be spinned up in one click. This will help you build apps easily and quickly using state-of-the-art cloud native stacks.

Watch "Heighliner Introduction" vedio:

[![IMAGE ALT TEXT](https://heighliner.dev/img/homepage/video-poster.png)](https://www.youtube.com/watch?v=74KZT-WW-lk&ab_channel=Heighliner "Heighliner Introduction")

## Why Heighliner
**Stack as Code (SaC)**: Your entire application stack can be codified. You can version, reuse, and share your stacks. You can even import existing stacks to build more advanced stacks. You can compose it in a way that optimizes for your environments, including Helm chart, CI/CD pipelines, logging and monitoring, security and access control, etc. We also provide official stacks to provide cloud native best practice out of the box.

**Seamless workflow**: Without Heighliner, we have seen people install and configure various tools (e.g. Argocd, Grafana, Nocalhost, API Gateway) on Kubernetes over and over again. It fragments their development time and makes them painful to connect the dots. With Heighliner, you can enjoy the seamless workflow for developing your apps, integrated with open source tooling. You can do everything on a single platform: writing code, building and testing, managing CI/CD pipelines, viewing logs and metrics.

**Declarative program**: Traditional tools ask you to program workflow step by step. This method doesn't work at scale. Developers often get lost in an overwhelming amount of code. We need a new solution to meet the growing business requirements -- a declarative system to describe the desired goals. You can just compose the application architecture in high level and Heighliner will handle the heavy-lifting.

**Multi-cloud and no vendor lock-in**: Heighliner is open source, vendor neutral, cloud agnostic. With a multi-cloud, pluggable architecture, Heighliner can adapt your apps to any cloud platforms. Your code remains the same across cloud providers (AWS, Azure, etc.) while Heighliner integrates with them intelligently. You can truly build once and run anywhere.

## Getting Started

Check out the [documentation](https://heighliner.dev/docs/getting_started/installation) on how to start using heighliner.

## Community
Join us at [Discord](https://discord.gg/anRxH5uk)

## Development Status

Heighliner is in Alpha stage and being actively developed.

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

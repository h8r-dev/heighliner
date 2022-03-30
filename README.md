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

List all heighliner stacks:

```shell
hln list stack
```

Output:

```shell
NAME          VERSION  DESCRIPTION
sample        1.0.0    Sample is a light-weight stack mainly used for test
go-gin-stack  1.0.0    go-gin-stack helps you configure many cloud native components including prometheus, grafana, nocalhost, etc.
gin-vue       1.0.0    gin-vue is a new version of go-gin-stack
```

Choose a stack and create a project

```shell
hln new -s=sample
```

> note: Please prepare a [kubeconfig file](https://rancher.com/docs/rke/latest/en/kubeconfig/) and apply for a [github access token](https://github.com/settings/tokens) to continue.

Spin up your application

```shell
hln up -i
```

Input the value one by one according to the promt. Your application will be set up automatically.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

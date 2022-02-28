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

Config an environment:

```shell
hln env new demo --stack=sample
hln config list
hln config set app -f ./examples/sample/app.yaml
hln config set push.target ghcr.io/hongchaodeng/my-app
hln config set push.auth.username hongchaodeng
hln secret list
hln secret set push.auth.secret $GITHUB_TOKEN
hln secret set kubeconfig -f $KUBECONFIG
```

Create an application:

```shell
hln up
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

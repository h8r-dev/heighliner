# Heighliner

## Quickstart

Build client binary:

```shell
make hln
export PATH="$PWD/bin:$PATH"
```

cd to a workdir

Init heighliner:

```shell
hln init
```

Pull stack (name could be 'sample'):

```shell
hln stack pull <name>
```

Show stack:

```shell
hln stack show
```

Init stack:

```shell
hln stack init
```

List input values:

```shell
hln stack list
```

Input values:
```shell
hln stack input <type> <name> <value>
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

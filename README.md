# Heighliner

## Quickstart

Build client binary:

```shell
make hln
export PATH="$PWD/bin:$PATH"
```

Init heighliner

```shell
hln init
```

Pull stack:

```shell
hln stack pull sample
```

Show stack:

```shell
hln stack show
```

Init stack:

```shell
hln stack init
```

List input:

```shell
hln stack list
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

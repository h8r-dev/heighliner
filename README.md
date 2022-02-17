# Heighliner

## Quickstart

Build client binary:

```shell
make hln
export PATH="$PWD/bin:$PATH"
```

List stack:

```shell
hln stack list
```

Pull stack:

```shell
hln stack pull sample
```

Show stack:

```shell
hln stack show sample
```

Config an environment:

```shell
hln env new demo --stack=sample
hln config list
hln config set push.target ghcr.io/hongchaodeng/my-app
hln config set push.auth.username hongchaodeng
hln secret list
hln secret set push.auth.secret $GITHUB_TOKEN
hln secret set kubeconfig -f $KUBECONFIG
```

Create an application:

```shell
hln up -f ./examples/sample/app.yaml
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)
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

> Requirements:
> - Install a [Kubernetes Cluster](https://heighliner.cloud/docs/getting_started/installation#2-install-kubernetes-cluster) and export the kubeconfig file.
> - Create a [GitHub personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) with the following scopes selected.
>   - `repo`
>   - `workflow`
>   - `write:packages`
>   - `delete:packages`
>   - `admin:org`
>   - `user`
>   - `delete_repo`

## Choose a stack

List all heighliner stacks:

```shell
hln list stack
```

Output:

```shell
NAME          VERSION  DESCRIPTION
gin-vue       1.0.0    gin-vue helps you configure many cloud native components including prometheus, grafana, nocalhost, etc.
```

At present `gin-vue` is the only available stack.

## Create a project

Use the stack you have chosen previously to create a heighliner project.

```shell
hln new -s=gin-vue
```

> If you have an existed project, then `hln` will prevent you from creating a new one. You can run `hln drop` to drop the old one.

## Spin up your application

Spin up your application interactively:

```shell
hln up -i
```

Input the values one by one according to the promt and your application will be set up automatically.

> If the progress of spinning up your application failed accidentally, feel free to **rerun** the `hln up -i` command and input the same values again. `hln` has the ability of caching.

The output will look like the following:

```shell
application:
  domain: -production.dkvmnr.go-gin.h8r.app
  ingress: 101.32.16.242
repository:
  backend: |
    https://github.com/h8r-dev/docs2.git
  frontend: |
    https://github.com/h8r-dev/docs2-front.git
  deploy: |
    https://github.com/h8r-dev/docs2.git
infra:
  - type: prometheus
    url: dkvmnr.prom.stack.h8r.io
  - type: grafana
    url: dkvmnr.grafana.stack.h8r.io
    username: admin
    password: |
      prom-operator
  - type: alertManager
    url: dkvmnr.alert.stack.h8r.io
  - type: argoCD
    url: dkvmnr.argocd.stack.h8r.io
    username: admin
    password: |
      LTXbYLAVFkPQ-Z9o
  - type: nocalhost
    url: dkvmnr.nocalhost.stack.h8r.io
    username: admin@admin.com
    password: "123456"

```

Congrats! You have initialized your application with `hln` successfully. Click the github url to start developing your new application. All of the Cloud-Native infrastructure and CI/CD pipelines have been set up properly. Feel free to click these links and input the account and password of each component to see its dashboard.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

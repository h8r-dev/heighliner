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

Make project dir:

```shell
mkdir helloapp && cd helloapp
```

List all heighliner stacks:

```shell
hln stack list
```

Output:

```shell
NAME          VERSION  DESCRIPTION
sample        1.0.0    Sample is a light-weight stack mainly used for test
go-gin-stack  1.0.0    Go-gin-stack helps you configure many cloud native components including prometheus, grafana, nocalhost, etc.
```

Choose a stack and create a project

```shell
hln new -s=sample
```

List all input values:

```shell
hln input list
```

Output:

```
Input          Value   Set by user  Description
hello.message  string  false        -
```

Config a input value

```shell
hln input text hello.message "Hello heighliner"
```

Spin up your application

```shell
hln up
```

Output:

```
[✔] hello.createContainer                              0.0s
[✔] hello.createFile.from                              0.0s
[✔] hello.createFile.contents                          0.0s
[✔] hello.outputMessage                                0.0s
Output               Value  Description
hello.outputMessage  """\n  Hello heighliner\n\n  """  -
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

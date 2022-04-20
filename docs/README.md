# Generate docs

Clone the repository:

```
git clone git@github.com:h8r-dev/heighliner.git
cd heighliner
```

Generate the docs:

```
mkdir docs/commands
go run ./docs/gen.go docs/commands
```

You will see the docs generated in `docs/commands`

## Generate docs of Commands for [heighliner](https://heighliner.dev/docs/cli/hln/commands/hln)

1. Go to the `h8r-dev` directory that contains these two projects:

```shell
$ tree -L 1
.
├── heighliner
└── heighliner-website
```

2. Clean up old docs

```shell
rm heighliner-website/docs/07-cli/hln/commands/*.md
```

3. Run generate command in heighliner root dir.

```shell
cd heighliner
go run ./docs/gen.go
```

4. Then you can check the difference in heighliner-website.

```shell
cd ../heighliner-website
git status
git add .
git commit -m "doc: update hln commands"
```
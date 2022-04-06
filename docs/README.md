# heighliner.cloud CLI Commands Doc

1. step up these two projects in the same folder.

```shell
$ tree -L 1
.
├── heighliner
└── heighliner-website
```

2. Run generate command in heighliner root dir.

```shell
cd heighliner/
go run ./docs/gen.go
```

3. Then you can check the difference in heighliner-website.

```shell
cd ../heighliner-website
git status
```
# Generate docs

Clone the repository:

```
git clone git@github.com:h8r-dev/heighliner.git
cd heighliner
```

Generate the docs:

```
go run ./docs/gen.go
```

You will see the docs generated in `docs/commands`

Each release will trigger the docgen workflow and generate documentations in the `website-docs' directory of 'website-docs' branch.

[Heighliner website](https://heighliner.dev/docs/cli/hln/commands/hln) will pull docs from branch website-docs every day

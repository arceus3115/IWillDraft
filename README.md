# IWillDraft

## Local Git setup

Enable the versioned pre-commit hook from the repository root after cloning:

```sh
git config --local core.hooksPath .githooks
```

The hook checks staged files for credential filenames, common token formats, private keys, and literal `GRID_API_KEY` values. It also checks for whitespace errors and conflict markers. Secret findings show filenames only. The GitHub Actions workflow runs the same secret scan over the checked-out files on pushes and pull requests.

## GRID SDK

The [GRID Go client](https://github.com/arceus3115/grid-framework) is maintained in its own repository. Add it as a dependency from your Go module directory.

```sh
go get github.com/arceus3115/grid-framework@v0.1.0
```

Keep `GRID_API_KEY` in your server environment or an ignored local `.env` file. The [environment template](.env.example) lists the variable, and the [SDK setup guide](https://github.com/arceus3115/grid-framework#prerequisites) shows how to export it for Go.

## Hackathon Requirements 
### Macro Scouting & Draft Management System
- A competitive League of Legends coaching platform that eliminates spreadsheet chaos, elevates macro strategy, or accelerates pre-match preparation
- GRID data
### Technological Implementation
Does the project reflect genuine effort and a working, non-trivial implementation?
### Design
Does the project deliver a complete, coherent product experience, not just a technical proof of concept?
### Potential Impact
Does the project make a credible, specific case for solving a real problem for a real audience, and does the solution actually address that problem based on what's demonstrated?
### Quality of the Idea
How creative and novel is the concept and does the project differ from existing concepts?

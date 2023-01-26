
# Gogit

Gogit is a mini version control system like Git that can track changes across a set of files or directories. Although this project's feature set isn't very complete and could use some more work, it currently supports the following:
- Commits
- Commit history
- Branches
- Searching
- Specific file's history



## Run Locally

This assumes that you have Go installed in your system.

Clone the project

```bash
  git clone https://github.com/rushitote/gogit
```

Go to the project directory

```bash
  cd gogit
```

Build the project

```bash
  go build
```

Run using the `gogit` command:
```bash
  ./gogit init
```


## Usage

For information regarding the commands, use the 'help' command:
```bash
./gogit help
```

```bash
A simple VCS

Usage:
  gogit [command]

Available Commands:
  before      Show commit logs before some time
  branch      Create/Delete/Rename a branch
  cb          Checkout a branch
  checkout    Checkout a commit
  commit      Commit changes to the repository
  completion  Generate the autocompletion script for the specified shell
  fh          File history
  gc          Garbage collection
  help        Help about any command
  log         Show commit logs
  mb          Merges two branches
  merge       Merges two commits
  play        Move across commits
  search      Search for a commit

Flags:
  -h, --help   help for gogit

Use "gogit [command] --help" for more information about a command.
```

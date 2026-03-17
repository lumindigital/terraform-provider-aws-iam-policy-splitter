# Contributing to Juno

This document outlines the process of committing changes to the Juno project. In order to ensure that your contributions are accepted, please open an issue, and get the issue assigned to you before opening a PR.


## Code of Conduct

You must read and agree to follow the [Code of Conduct](./CODE_OF_CONDUCT.md)

### Task.yml
This repository uses [Task](https://taskfile.dev/) instead of a makefile. Tasks can be installed by following the [Installation Guide](https://taskfile.dev/docs/installation)

### Commit Hooks
This repository uses [pre-commit](https://pre-commit.com) hooks to ensure code quality and consistency. It allows you to automatically run linters, formatters, and other tools before committing code.
These checks also run when you open a PR.

To install pre-commit follow the instructions [here](https://pre-commit.com/#install)

To install the hooks in this repo, run the following command from the root of the repository.
```bash
pre-commit install
```
## General workflow

Before beginning work, please make sure the issue you are working on has been assigned to you.
This workflow uses the [Conventional Commits](https://www.conventionalcommits.org) specification for commit messages.

The Workflow is as follows
1. Fork the repository under your GitHub username.
2. Create a branch in your forked repository.
3. Commit changes required to the branch you created on step 2.
4. Verify that tests and pre-commit passes
5. Push your branch to your forked repository.
6. Create a Pull Request from your remote fork pointing to the HEAD branch (usually `master` branch) of the target repository.
7. Check the GitHub build and ensure that all checks are green.

## Testing
If your adding new functionality or modifying existing functionality and there is currently not tests, please add them as part of your pull request.
Tests can be run with the following command:
```bash
task local:test && task local:testacc
```

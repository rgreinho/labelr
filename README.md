# labeler

Manage your GitHub labels efficiently.

With `labeler`, managing your GitHub labels becomes effortless. `labeler` will attempt to detect all the information required to apply the labels wherever you need them to be.

## Infered values and environment variables

`labeler` will automatically detect the owner or organization and the repostiory from the directory where you are running the command. It will also look automatically for a file named `labels.yml`.

The following environment variables are used by `labeler`:

* GITHUB_ORGANIZATION
* GITHUB_REPOSITORY
* GITHUB_USER
* GITHUB_TOKEN

### Precedence

`labeler` looks for information in this order:

1. Infered information from current directory
2. environment variables
3. CLI arguments

## Existing labels

For existing labels, description and color will be updated to match the content of `the labels.yml` file.

However, **labels cannot be renamed**. this is due to the fact that the tool does not keep track of the existing configuration. If the name of a label gets changed, a new label will be created instead.

## Usage examples

### From a .github repository

This is the main and most common use case.

It is common for organizations to store their templates and GitHub configuration files in a `.github` repository. For instance,
here is an example of `.github` [repository](https://github.com/aura-atx/.github) from AURA.

The content looks like this:

```bash
.
├── ISSUE_TEMPLATE
│   ├── bug_report.md
│   └── feature_request.md
├── labels.yml
└── stale.yml
```

A  `git remote -vv` provides the following information:

```bash
origin	https://github.com/aura-atx/.github (fetch)
origin	https://github.com/aura-atx/.github (push)
```

From this setup, and unless specified otherwise from the CLI, `labeler` will use the following information:

* `owner` and `organization` will be `aura-atx`
* `repository` will be `.github`
* `token`  will be read from the `${GITHUB_TOKEN}` environment variable
* the new labels will be read from the `labels.yml` file

#### Apply the new labels to all the repositories in the organization

```bash
labeler apply --org
```

#### Remove existing labels and apply the new ones to all the repositories in the organization

```bash
labeler apply --sync --org
```

#### Apply labels to the current repository

```bash
labeler apply [--sync]
```

#### Apply labels to a specific repository from the same owner or organization

```bash
labeler apply [--sync] --repository other-repository
```

#### Apply labels from a random directory

You may want to run `labeler` from anywhere on your system. Not a problem, simply pass all the values to the CLI:

```bash
labeler apply [--sync] --owner myself --repository my-repository --token ${OTHER_GITHUB_TOKEN} /tmp/my-label-file.yml
```

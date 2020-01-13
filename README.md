# labelr

![](https://github.com/rgreinho/labelr/workflows/ci/badge.svg)

Manage your GitHub labels efficiently.

With `labelr`, managing your GitHub labels becomes effortless. `labelr` will attempt to detect all the information required to apply the labels wherever you need them to be.

## Infered values and environment variables

`labelr` will automatically detect the owner or organization and the repostiory from the directory where you are running the command. It will also look automatically for a file named `labels.yml`.

The following environment variables are used by `labelr`:

* GITHUB_ORGANIZATION
* GITHUB_REPOSITORY
* GITHUB_USER
* GITHUB_TOKEN

### Precedence

`labelr` looks for information in this order:

1. Infered information from current directory
2. environment variables
3. CLI arguments

## Existing labels

For existing labels, description and color will be updated to match the content of `the labels.yml` file.

However, **labels cannot be renamed**. This is due to the fact that the tool does not keep track of the existing configuration. If the name of a label gets changed, a new label will be created.

## labels.yml

The `labels.yml` file has a simple format:

```yml
---
labels:
  - name: "kind/bug"
    color: "#D73A4A"
    description: "Something isn't working"
```

The top level key `labels` is used to group the labels together. Each label then becomes an entry under this key.

Each label entry is composed of the following fields:

* `name` (required)
* `color` (required)
* `description` (optional)

For a complete example, have a look at the labels used
[for this project](https://github.com/rgreinho/labelr/blob/master/.github/labels.yml).

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

From this setup, and unless specified otherwise from the CLI, `labelr` will use the following information:

* `owner` and `organization` will be `aura-atx`
* `repository` will be `.github`
* `token`  will be read from the `${GITHUB_TOKEN}` environment variable
* the new labels will be read from the `labels.yml` file

#### Apply the new labels to all the repositories in the organization

```bash
labelr apply --org
```

#### Remove existing labels and apply the new ones to all the repositories in the organization

```bash
labelr apply --sync --org
```

> *NOTE: without the `--sync` option, existing labels are updated and new labels added. Nothing gets removed.*

#### Apply labels to the current repository

```bash
labelr apply [--sync]
```

#### Apply labels to a specific repository from the same owner or organization

```bash
labelr apply [--sync] --repository other-repository
```

### Apply labels from a random directory

You may want to run `labelr` from anywhere on your system. Not a problem, simply pass all the values to the CLI:

```bash
labelr apply [--sync] \
  --owner myself \
  --repository my-repository \
  --token ${OTHER_GITHUB_TOKEN} \
  /tmp/my-label-file.yml
```

# Contributing

Thank you for your interest in contributing!

We accept contributions via GitHub pull requests.

## How to run this project locally

To run this project locally, clone the repository and make sure that all necessary tools defined in mise.toml are installed.
We recommend that you use mise to install these (run `mise install` in the current directory) but you do don't have to.

Also make sure to copy the example config to the root folder:

```shell
cp examples/config.yaml config.yaml
```

Make sure to have an instance of Dex running locally.

Start this project with:

```shell
mise run serve
```
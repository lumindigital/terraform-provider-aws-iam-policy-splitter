[![Dependabot Updates](https://github.com/lumindigital/juno/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/lumindigital/juno/actions/workflows/dependabot/dependabot-updates)
[![release](https://github.com/lumindigital/juno/actions/workflows/release.yml/badge.svg)](https://github.com/lumindigital/juno/actions/workflows/release.yml)
[![Release Version](https://img.shields.io/github/v/release/lumindigital/juno)](https://github.com/lumindigital/juno/releases/latest)

## Terraform AWS IAM Policy Splitter Provider
This provider contains a single data source for terraform. The data source takes as input an array of stringified json policy documents. It performs bin packing using the best-fit decreasing bin packing algorithm on said policy documents and returns an array of arrays containing the json policy documents.

### Documentation
[View the docs](docs/index.md)

# Pod Chaos Monkey

## Introduction

This app serves as an example on how to introduce some chaos in a Kubernetes namespace. It will delete a pod at random in a particular namespace. Since this is done periodically, the tool doesn't need to be continuously running, as such it will be deployed into Kubernetes as a CronJob and we let the orchestrator take care of the scheduling.

# Project Structure

    .
    ├── cmd                     # Executables directory
    ├── internal                # Private application and library code.
    ├── pkg                     # Library code that is OK to be used by external applications
    ├── test                    # Automated tests (alternatively `spec` or `tests`)
    ├── tools                   # Tools and utilities
    ├── LICENSE
    └── readme.md

# How does it work?

If we're using this from outside a Kubernetes Cluster (i.e. as a Standalone Binary or while developing) we can set the `KUBECONFIG` environment variable with a filepath. Otherwise the program will default to use the default in-cluster config.
Example:

```bash
KUBECONFIG=~/.kube/config
```

# Pod Chaos Monkey

## Introduction

This app serves as an example on how to introduce some chaos in a Kubernetes namespace. It will delete a pod at random in a particular namespace when it matches an a label selector. The default label selector is ""

## Project Structure

```shell
    .
    ├── cmd                     # Executables directory
    ├── internal                # Private application and library code.
    ├── pkg                     # Library code that is OK to be used by external applications
    ├── test                    # Automated tests (alternatively `spec` or `tests`)
    ├── tools                   # Tools and utilities
    ├── LICENSE
    └── readme.md
```

## How does it work?

If we're using this from outside a Kubernetes Cluster (i.e. as a Standalone Binary or while developing) we can set the `KUBECONFIG` environment variable with a filepath. Otherwise the program will default to use the default in-cluster config.
Example:

```shell
KUBECONFIG=~/.kube/config
```

## Configuration

The application can be configured with the following environment variables:

```shell

| Variable             | Default                     | Description                                                                                                                   |
|----------------------|-----------------------------|-------------------------------------------------------------------------------------------------------------------------------|
| KUBECONFIG           | none (optional)             | Set the filepath for the KUBECONFIG file you wish to use.                                                                       |
| NAMESPACE            | workloads                   | Namespace that the chaos-monkey-app uses to kill pods                                                                         |
| INTERVAL_SECONDS     | 10                          | Time in seconds between pod deletes                                                                                           |
| GRACE_PERIOD_SECONDS | 5                           | Grace period between the start of the a pp and the moment it starts deleting pods                                             |
| LABELS               | podchaosmonkey=true         | This label needs to be set on the workloads on the namespace to mark them as elegible for the podchaosmonkey application      |
```

The `KUBECONFIG` variable can be set to a filepath. If nothing explicit is supplied the application will attempt to use the `in-cluster` config.

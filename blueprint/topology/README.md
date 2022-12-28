# topology

## Description
topology controller

## Usage

### Fetch the package
`kpt pkg get REPO_URI[.git]/PKG_PATH[@VERSION] topology`
Details: https://kpt.dev/reference/cli/pkg/get/

### View package content
`kpt pkg tree topology`
Details: https://kpt.dev/reference/cli/pkg/tree/

### Apply the package
```
kpt live init topology
kpt live apply topology --reconcile-timeout=2m --output=table
```
Details: https://kpt.dev/reference/cli/live/

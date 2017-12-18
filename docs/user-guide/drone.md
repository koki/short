## Drone Plugin for Koki Short

Define all your Kubernetes manifests using Koki Short, and use this drone plugin to automatically convert the Short manifests into K8s manifests on the fly

### Use Koki Short in your Projects

```
workspace:
  base: /go
  path: src/github.com/koki/short-drone-plugin

  pipeline:
    koki-short:
      image: kokster/short-drone-plugin:0.3.0
      files:
      - deployment.yaml
```

This will convert `deployment.yaml` to `kube_deployment.yaml`, which is the kubernetes manifest equivalent to `deployment.yaml`. The next steps in the pipeline will now be able to use `kube_deployment.yaml`

### Configuration Options

| Option | Type | Description | 
|--------|------|-------------|
| files  | []string | Input files relative to root of the project which is being built using drone |
| overwrite | bool | Set to `true` to allow output files to be overwritten. (default `false`) |
| in_place | bool | Set to `true` to translate files in place. (default `false`). Should always be used with `overwrite: true` |
| prefix | string | The prefix of the output file created. (default `kube_`) |

# Koki Short

<img src="https://codebuild.us-east-1.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiaFdaQUJleFNQRkI3NXdKcllWSVJ0THREYkl2WEJBTkZCRkhrWERoeUYxbjFsUTZNSEZ4WGxWZnNmcWxvenlMallELytyNTI0VSsxZ0JEV3FrOS9JYzVzPSIsIml2UGFyYW1ldGVyU3BlYyI6InFKYXBITkJ0Z3NsTkhWN0UiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master"/>

Manageable Kubernetes manifests through composable, reusable syntax

## Motivation
The description format for Kubernetes manifests, as it stands today, is verbose and unintuitive. Anecdotally, it has been:

 - Time consuming to write
 - Error-prone, hard to get right without referring to documentation
 - Difficult to maintain, read, and reuse

e.g. In order to create a simple nginx pod that runs on any host in region `us-east1` or `us-east2`, here is the Kubernetes native syntax:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  containers:
  - name: nginx_container
    image: nginx:latest
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
      	nodeSelectorTerms:
        - matchExpressions:
          - key: k8s.io/failure-domain
            operator: In
            values: 
            - us-east1
        - matchExpressions:
          - key: k8s.io/failure-domain
            operator: In
            value:
            - us-east2
```

The Short format is designed to be user friendly, intuitive, reusable, and maintainable. The same pod in Short syntax looks like

```yaml
pod:
  name: nginx
  labels:
    app: nginx
  containers:
  - name: nginx
    image: nginx:latest
  affinity:
  - node: k8s.io/failure-domain=us-east1,us-east2
```

Our approach is to reframe Kubernetes manifests in an operator-friendly syntax without sacrificing expressiveness.

Koki Short can transform Kubernetes syntax into Short and Short syntax back into Kubernetes. No information is lost in either direction.

For more information on Koki Short transformations, please refer to [Resources.](https://docs.koki.io/short/resources)

## Modular and Reusable

Koki Short introduces the concept of modules, which are reusable collections of Short resources. Any resource can be reused multiple times in other resources and linked resources can be managed as a single unit on the Koki platform. 

Any valid koki resource object can be reused. This includes subtypes of top-level resource types. For example, here's module called `affinity_east1.yaml`:

```yaml
affinity:
- node: k8s.io/failure-domain=us-east-1
```

This affinity value can be reused in any pod spec:

```yaml
imports:
- affinity: affinity_east1.yaml
pod:
  name: nginx
  labels:
    app: nginx
  containers:
  - name: nginx
    image: nginx-latest
  affinity: ${affinity}  # re-use the affinity resource here
```

For more information on Koki Modules, please refer to [Modules.](https://docs.koki.io/short/modules)

## Getting started

In order to start using Short, simply download the binary from the [releases page](https://github.com/koki/short/releases).

```sh

#start with any existing Kubernetes manifest file
$$ cat kube_manifest.yaml
apiVersion: v1
kind: Pod
metadata:
  name: podName
  namespace: podNamespace
spec:
  hostAliases:
   - ip: 127.0.0.1
     hostnames:
      - localhost
      - myMachine
  affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
     - matchExpressions:
        - key: k8s.io/failure-domain
          operator: In
          value: us-east1
  containers:
   - name: container
     image: busybox:latest
     ports:
      - containerPort: 6379
        hostPort: 8080
        name: service
        protocol: TCP
     resources:
        limits:
          cpu: "2"
          memory: 1024m
        requests:
          cpu: "1"
          memory: 512m

     
#convert Kubernetes manifest into a Short syntax representation
$$ short -f kube_manifest.yaml
pod:
  name: podName
  namespace: podNamespace
  affinity:
   - node: k8s.io/failure-domain=us-east1
  host_aliases:
   - 127.0.0.1 localhost myMachine
  containers:
   - image: busybox
     cpu:
       min: "1"
       max: "2"
     mem:
       min: "512m"
       max: "1024m"
     expose:
      - port_map: "8080:6379"
        name: service


#input can be json or yaml
$$ short -f kube_manifest.json -f kube_manifest2.yaml -f kube_multi_manifest.yaml

#stream input
$$ cat kube_manifest.yaml | short -

#revert to kubernetes type
$$ short -k -f koki_spec.yaml

#-k flag denotes that it should output Kubernetes manifest

```

For more information, refer to our [getting started guide.](https://docs.koki.io/user-guide#getting-started)

## Contribute
Koki is completely open source community driven, including the roadmaps, planning, and implementation. We encourage everyone to help us make Kubernetes manifests more manageable. We welcome Issues, Pull Requests and participation in our [weekly meetings]().

If you'd like to get started with contributing to Koki Short, read our [Roadmap](https://github.com/koki/short/projects) and start with any issue labelled `help-wanted` or `good-first-issue`

## Important Links

- [Releases](https://github.com/koki/short/releases)
- [Docs](https://docs.koki.io/short)
- [Roadmap](https://github.com/koki/short/tree/master/releases/ROADMAP.md)
- [Issues](https://github.com/koki/short/issues)

## LICENSE
[Apache v2.0](https://github.com/koki/short/blob/master/LICENSE)

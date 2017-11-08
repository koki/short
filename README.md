# Koki Short
API friendly Kubernetes resources to user-friendly syntax

## Motivation
The description format for Kubernetes manifests, as it stands today, is verbose and unintuitive. It has anecdotally been 
 - Time consuming to use it to create Kubernetes resources
 - Error prone to get it right the first time
 - Requires constant referral to documentation
 - Difficult to maintain, read or reuse.  

For eg. denoting that a pod runs on host with label `k8s.io/failure-domain=us-east-1`, here is the current syntax:

```yaml
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
     - matchExpressions:
        - key: k8s.io/failure-domain
          operator: In
          value: us-east1
```
The Short format is designed to be user friendly, intuitive, reusable and maintainable. The same affinity syntax in Short looks like

```yaml
affinity:
 - node: k8s.io/failure-domain=us-east1
```
The approach we have followed behind this opinionated syntax is to reframe the Kubernetes syntax to one that is operator friendly without losing any information. 

Since we do not throw away any information during the transformations, users can freely round trip back and forth between Short syntax and Kubernetes syntax.

For more information on Koki Short transformations, please refer to [docs](http://docs.koki.io/short)

## Modular and Reusable

Koki Short introduces the concept of modules, which are reusable collections of Short resources. Any resource can be reused multiple times in other resources and linked resources can be managed as one unit on the Koki platform. 

More on this will be available as soon as it is implemented. This is in the roadmap right now. 

## Getting started

In order to start using Short, simply download the binary from the [releases page](https://github.com/koki/short/releases), and then you can start using it.

```sh

#start with an existing Kubernetes manifest file
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

     
#convert an existing Kubernetes manifest into a Short syntax representation, just run
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


#specify json input, or multi document yaml 
$$ short -f kube_manifest.json -f kube_manifest2.yaml -f kube_multi_manifest.yaml

#stream input
$$ cat kube_manifest.yaml | short -

#revert to kubernetes type
$$ short -k -f koki_spec.yaml

#-k flag denotes that it should output Kubernetes manifest

#find out how koki transforms a particular resource type
short man v1/pod

```

## Contribute
Koki is completely open source community driven, including the roadmaps, planning, and implementation. We encourage everyone to help us make Kubernetes manifests more manageable. We welcome Issues, Pull Requests and participation in our [weekly meetings]().

If you'd like to get started with contributing to Koki Short, read our [Roadmap](https://github.com/koki/short/projects) and start with any issue labelled `help-wanted` or `good-first-issue`

## Important Links

- [Releases](https://github.com/koki/short/releases) 
- [Docs]()
- [Roadmap](https://github.com/koki/short/projects)
- [Issues](https://github.com/koki/short/issues)

## LICENSE
[Apache v2.0](LICENSE)

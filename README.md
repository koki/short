# short(hand)

Kubernetes resource manifests should be easier to use.
We're making them more human-friendly with `short`.

`short` is simplified file format and a set of tools and docs to help you use Kubernetes.

It's implemented as a set of carefully crafted Go types that directly translate to and from the native k8s API.

## Example

*Native k8s API*

```
TODO: before
```

*koki short*
```
TODO: after
```

# Convert your existing k8s files

```
# Convert from a file:
short -f pod.yaml
short -f pod.json

# Convert from a stream:
cat pod.yaml | short -

# Convert from a url:
short -f http://spec.com/pod.yaml

# Output to a file:
short -f pod.yaml -o pod_short.yaml
``` 

# Writing short files

Coming soon!

# How to contribute

Given the breadth (and ever-evolving state) of the Kubernetes API,
this project will be most successful as a community effort. There are lots of ways to help:

* Send us your manifest files! We'll add them to our automated tests. They'll help our documentation, too. (How to send them? See the next point.)

* File a new Issue. If you find a bug, have an idea for an improvement, or anything else, file an Issue, and we'll see it.

* Take on an existing Issue. Be on the lookout for the `good first issue` and `help wanted` tags.

We're starting out with a `short` implementation of just a few core resources, but we'll increase coverage over time. Help us get them all!

# License

[Apache v2.0](LICENSE)

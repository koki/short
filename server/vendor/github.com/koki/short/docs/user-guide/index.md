# Getting started with Koki Short

Koki Short is a tool for writing composable and reusable Kubernetes manifest files. In order to start using it, download the [latest release.](https://github.com/koki/short/releases/latest)

You can also build it from source.
```sh
  # get the short project
  $$ go get github.com/koki/short
```

It is now ready to be used to create Kubernetes manifest files, as well as to convert existing Kubernetes manifests.

# Creating Kubernetes manifests

**We encourage users to use the short syntax to specify their Kubernetes manifests**. This format is easier to write, read, and reuse. 

 - Start with a simple pod resource in the Short syntax

```sh
   # create a simple pod called test that runs nginx
   $$ cat << EOF > kube-manifest.short.yaml
      pod:
       name: test
       labels:
        app: test 
       containers:
       - name: nginx
         image: nginx
      EOF
```

   Notice that the boiler plate that is generally required to specify kubernetes resources (such as `apiVersion` or `metadata`) need not be specified.

 - Convert it into a Kubernetes resource

```sh
   # pass the pod specified in short syntax into short
   $$ short -k -f kube-manifest.short.yaml # -k flag instructs Short to convert to Kubernetes native type
   kind: Pod
   apiVersion: v1
   metadata:
     creationTimestamp: null
     labels:
       app: test
     name: test
   spec:
     containers:
     - image: nginx
       name: nginx
     resources: {}
   status: {}
```

The `-k` flag is used here to instruct Short to convert to Kubernetes native type. Additionally, Note that the boiler plate has been automatically added to the generated Kubernetes manifest. 

This syntax **does not** drop any information from the Kubernetes resource. It can be used to create any Kubernetes resource type.   

For complete information on Short syntax for various types, refer our [Resources section.](../resources/index.md)

# Reading Kubernetes manifests

Short can be used to make your existing Kubernetes manifests more readable. It can convert existing Kubernetes manifest files to short format. 

```sh
   # save the earlier generated file
   $$ short -f kube-manifest.short.yaml > kube-manifest.yaml

   # convert it to short format
   $$ short -f kube-manifest.yaml
   pod:
    containers:
    - image: nginx
      name: nginx
    labels:
      app: test
    name: test
    version: v1
```

For complete information on how the conversion from one type to another is performed, refer our [Resources section.](../resources/index.md)

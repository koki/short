# Command Line Reference

Short converts API-friendly Kubernetes syntax into ops-friendly syntax.

```sh
Usage:
  short [flags]
  short [command]

Available Commands:
  help        Help about any command
  man         Reference and Examples for resources and conversions
  version     Prints the version of short

Flags:
      --alsologtostderr                  log to standard error as well as files
  -f, --filenames strings                path or url to input files to read manifests
  -h, --help                             help for short
  -k, --kube-native                      convert to kube-native syntax
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files (default false)
  -o, --output string                    output format (yaml*|json) (default "yaml")
  -s, --silent                           silence output to stdout
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging

Use "short [command] --help" for more information about a command.
```
Short can convert to and from Short syntax and Kubernetes syntax. Short supports levelled logging for debugging purposes. The log level can be set by using the `-v` flag with a number between 1 and 20.

# Output format

The output from Short can be represented into valid YAML or valid JSON. The user can choose the desired format by using the `-o` flag to denote the output type. 

Valid values for the `-o` flag are `yaml` or `json` (case-insensitive)

```sh
# start with a pod spec in short syntax
$$ cat kube-manifest.short.yaml
pod:
  name: test
  labels:
    app: test 
  containers:
  - name: nginx
   image: nginx

# output yaml
$$ short -k -f kube-manifest.short.yaml -o yaml  # default value for -o is yaml. 
apiVersion: v1
kind: Pod
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


# output json
$$ short -k -f kube-manifest.short.yaml -o json
{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "creationTimestamp": null,
        "labels": {
            "app": "test"
        },
        "name": "test"
    },
    "spec": {
        "containers": [
            {
                "image": "nginx",
                "name": "nginx",
                "resources": {}
            }
        ]
    },
    "status": {}
}
```

# Multiple inputs

Short can read in multiple input files and convert them to the desired format. In order to specify multiple files to short, the `-f` flag can be used. The `-f` flag can be specified multiple times with each of the multiple files corresponding to one of the `-f` flag values. 

```sh
# start with two files
$$ cat kube-manifest-1.yaml
pod:
  containers:
  - expose:
    - 80
    image: nginx
    name: nginx
    readiness_probe:
      delay: 5
      timeout: 5
  labels:
    name: nginx
  name: nginx
  version: v1

$$ cat kube-manifest-2.yaml
pod:
  name: test
  labels:
    app: test 
  containers:
  - name: nginx
    image: nginx

# convert both files at once to Kubernetes syntax
$$ short -k -f kube-manifest-1.yaml -f kube-manifest-2.yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    name: nginx
  name: nginx
spec:
  containers:
  - image: nginx
    name: nginx
    ports:
    - containerPort: 80
      protocol: TCP
    readinessProbe:
      initialDelaySeconds: 5
      timeoutSeconds: 5
    resources: {}
status: {}
---
apiVersion: v1
kind: Pod
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

# Streaming in files

Short can also stream in files through the `|` pipe operator. In order to activate the reading of input from a stream, specify an `-` at the end of the command. 

```sh
$$ cat kube-manifest.short.yaml
pod:
  name: test
  labels:
    app: test 
  containers:
  - name: nginx
   image: nginx

# output yaml
$$ cat kube-manifest.short.yaml | short -k - 
apiVersion: v1
kind: Pod
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

*Note that if you stream in a file as well as specify `-f`, only the file provided via `-f` will be used.*

# Version

Short follows Semver. You can find the version of the running short using the `version` command.

```sh
# short version
$$ short version
koki/short: v0.0.1
```

# Getting Help

Each Short command and sub-command has a help flag. (`--help` or `-h`)

More information about getting help is provided in the next sub-section. 

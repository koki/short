# Current Features:

Short syntax for:
- Pod
- Service
- Deployment, Replica Set, Replication Controller
- Stateful Sets
- Persistent Volumes, Volumes
- Persistent Volume Claims
- Jobs, Cron Jobs
- Config Maps
- Secret
- Storage Class
- Endpoints, Ingress

Conversion to and from Kubernetes-native syntax.

Validation errors show the location of the error in the input file using `$.pod.containers.1.env` path syntax.
Catches typos by looking for extraneous fields in the input file.

Parameterization and templating support for Short manifests.

Chrome plugin and backend server for converting Kubernetes YAMLs on GitHub: https://github.com/koki/kubernetes-viewer

Drone plugin for KubeCI: https://github.com/koki/short-drone-plugin

# Planned Features:

* Support every Kubernetes resource type.
* Support older versions of Kubernetes resource types.
* Support GitHub Gists in the Chrome plugin.
* Support Helm charts.



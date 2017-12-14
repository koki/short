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

Conversion to and from Kubernetes-native syntax.

Validation errors show the location of the error in the input file using `$.pod.containers.1.env` path syntax.
Catches typos by looking for extraneous fields in the input file.

Parameterization and templating support for Short manifests.

Chrome plugin and backend server for converting Kubernetes YAMLs on GitHub.

# Planned Features:

Support every Kubernetes resource type.
- CSI persistent volume

Support older versions of Kubernetes resource types.

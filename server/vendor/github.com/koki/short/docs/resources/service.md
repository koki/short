# Introduction

A Service is an umbrella over a group of pods and policies that define how to access them.

| API group | Resource | Kube Skeleton                                   |
|:---------:|:--------:|:-----------------------------------------------:|
| core/v1   | Service  |  [skel](../skel/service.kube.skel.yaml)         |

Here's an example Kubernetes Service spec:
```yaml
kind: Service
apiVersion: v1
metadata:
  name: my-service
spec:
  selector:
    app: MyApp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 9376
```

The following sections contain detailed information about each field in Short syntax, including how the field translates to and from Kubernetes syntax.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.` `clusterName` | The name of the cluster on which this Service is running |
|name | `string` | `metadata.name`| The name of the Service | 
|namespace | `string` | `metadata.` `namespace` | The K8s namespace this Service will be a member of | 
|labels | `string` | `metadata.labels`| Metadata about the Service, including identifying information | 
|annotations| `string` | `metadata.` `annotations`| Non-identifying information about the Service | 
|cname | `string` | `externalName` | This service will return a CNAME that is set by this field. No proxying will be performed|
|type | `string` | `type` | The type of the service. Can be omitted (for `cname` services) or set to "cluster-ip", "node-port" or "load-balancer"|
|selector| `map[string]` `string` | `selector` | A set of key-value pairs that match the labels of pods that should be proxied to|
|external_ips| `[]string` | `externalIPs` | A set of ip addresses for which nodes in the cluster will accept traffic|
|port | `string` | `ports` | Unnamed port mapping of format `$PROTOCOL://$SVC_PORT:$CONTAINER_PORT`. More details below|
|node_port| `int32` | `ports` | Request specific node port for a node-port service | 
|ports | `[]NamedPort`| `ports` | A list of named ports to expose. See [Named Port Overview](#named-port-overview)|
|cluster_ip| `string` | `clusterIP`| Request specific cluster ip for the service|
|unready_ endpoints| `bool` | `publishNot` `ReadyAddresses` | Publish addresses before backends are ready|
|route_policy| `string` | `externalTraffic` `Policy` | Policy for routing external traffic. Can be "node-local" or "cluster-wide" |
|stickiness | `int` or `bool` | `sessionAffinity` and `sessionAffinity` `Config` | Stickiness Policy for the service. More information below |
|lb_ip | `string` | `loadBalancerIP` | Request specific IP address for the created LB service|
|lb_client_ips | `[]string` | `loadBalancer` `SourceRanges` | (for LB service) IP addresses to allow traffic from. Can specify CIDR here |
|healthcheck_ port | `int32` | `healthCheck` `NodePort`  | Port for health check|

The following fields are status fields, and cannot be set

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| endpoints | `[]string` | `status.loadBalancer.ingress` | A list of hostnames or IP addresses that reflect the ip address/hostname of the created LB |

The `port` and `node_port` keys are used to specify the unnamed port for the Service. Kubernetes only allows one unnamed port. Therefore, there is only one field for unnamed port. The unnamed `port` should be of the format

`$PROTOCOL://$SERVICE_PORT:$CONTAINER_PORT`

The `$PROTOCOL` is defaulted to `tcp` if not specified. 

If only one port value is specified, then `$SERVICE_PORT` and `$CONTAINER_PORT` are considered to be the same value.

The following are valid values for port

```yaml
port: tcp://8080:80 # Protocol: TCP ServicePort: 8080 ContainerPort: 80
port: 8080:80 # Protocol: TCP ServicePort: 8080 ContainerPort: 80
port: udp://8080 # Protocol: UDP ServicePort: 8080 ContainerPort: 8080
```
The stickiness values can be set to `true` or a valid integer value. If set to `true`, it configures the session affinity of the service based on client ip addresses.

If set to a number value, then, along with configuring the session affinity, it also configures the number of seconds the session affinity sticks to a client ip address, before it expires. 

#### Named Port Overview

Named ports contain a struct that contains two keys

 - `node_port` - an `int32` type field that denotes a named node_port of the service
 - `$name` - a `string` type field of the format `$PROTOCOL://$SERVICE_PORT:$CONTAINER_PORT`

The `$name` key is prefixed with a `$` to denote that it can be any valid Kubernetes name, and that will be set as the name of the port. The semantics of the value follow the same behavior as the `port` field in Service (explained above). 

Here are valid named port examples
```yaml
ports:
- web: 8080:80  # port named "web" with service-port: 8080 and container-port: 80
- db: 8080      # port named "db" with node-port set to 5432
  node_port: 5432  
- dns: udp://53:53
  node_port: 53
```

# Examples 

 - A ClusterIP service with stickiness 

```yaml
service:
  name: web-service
  ports:
  - http: 80:8080 # service-port:container-port
  selector:
    app: web
  stickiness: true
  type: cluster-ip # default value is cluster-ip
  version: v1
```

 - A load balancer service with specific endpoints and ingress IPs

```yaml
service:
  endpoints:
  - aws.elb.jngfdgdgpkdfgk484989485.amazonaws.com
  - 45.54.67.79
  name: lb-service
  ports:
  - http: 80:8080
  route_policy: node-local
  selector:
    app: web
  type: load-balancer
  version: v1
```

 - A node port service with a node port chosen 

```yaml
service:
  name: web
  ports:
  - http: 80  # service-port, container-port is same as service-port
    node_port: 32123
  selector:
    app: web
  type: node-port
  version: v1

```

# Skeleton

| Short Type           | Skeleton                                       |
|:--------------------:|:----------------------------------------------:|
| Service                  | [skel](../skel/service.short.skel.yaml)            |

Here's a starter skeleton of a Short Service.
```yaml
service:
  name: app-service
  port: 80:8080  # service-port:container-port
  selector:
    app: web
  type: cluster-ip
  version: v1
```

# Modules

Modules facilitate the reuse and composition of Koki resources.


```yaml
imports:
- env: ./env.short.yaml
params:
- name: for selector/labels and name
  default: koki-short-server
deployment:
  name: ${name}
  labels:
    app: ${name}
  replicas: 3
  selector:
    app: ${name}
  affinity:
  - anti_pod: app=${name}
    topology: kubernetes.io/hostname
  containers:
  - image: ublubu/testserver:${env.server_version}
    name: ${name}
    pull: always
    expose:
    - 8080
    env:
    - COOKIE_AUTH_KEY=${env.cookie_auth_key}
    - STRIPE_KEY=${env.stripe_key}
    - GITHUB_CLIENT_ID=${env.github_client_id}
    - GITHUB_CLIENT_SECRET=${env.github_client_secret}
```

A koki module consists of three sections: `imports`, `params`, and the resource itself.

## Imports

Imports let you incorporate existing files into your module.


```yaml
imports:
- local_name_for_imported_module: ./path/to/imported/module.yaml
  params:
    param_foo: value_for_param_foo
    param_bar: ${template_expansion}
- local_name_for_another_imported_module: ./another/path.yaml
...
```

## Params

Params let you define how your module can be customized.


```yaml
params:
- param_foo: the description for a param named param_foo
  default: default_value_for_param_foo
- param_bar: the description for a param named param_bar. it doesn't have a default value.
```

When another module imports this one, it must provide values for the parameters defined here.

## Template Expansion

Koki supports logic-less text templating using `${identifier}`.
These template expansions can be used in two places:

 * the value of any param on an import

```yaml
imports:
- foo: ./foo.yaml
  params:
    x: ${bar} # the value of 'x' is the value of 'bar'
    y:
    - list item 0
    - ${baz...} # 'baz' is a list. append its items to 'y'.
    z: ${qux.a.x.0} # index into 'qux' to get the value for 'z'.
```

 * the contents of the defined resource--e.g. a Pod

```yaml
pod:
  x: ${bar} # the value of 'x' is the value of 'bar'
  y: 
  - list item 0
  - ${baz...} # 'baz' is a list. insert its items into 'y'.
  z: ${qux.a.x.0} # index into 'qux' to get the value for 'z'

# OR

pod: ${foo} # use 'foo' as 'pod'.
```

### Supported Templating Syntax

* Use a param named `foo` with `${foo}`.
* Use the resource imported as `foo` with `${foo}`
* Index into a list using `.n` for the item at index n.
    - e.g. `${foo.2}` for the third item of `foo`(index 2).
* Index into a dictionary using `.x` for the item at key x.
    - e.g. `${foo.containers}` for the `containers` field of `foo`.
* Indexes can be chained together.
    - e.g. `${foo.containers.1}` for the second container in `foo`.

If a template expansion contains a list, you can either use it as a list OR merge its items into an existing list:

```
list: ${identifier_for_list}

# OR

list:
- list item A
- ${identifier_for_list...}
- list item X
```

### Examples

For more examples, have a look at [these](https://github.com/koki/short/tree/master/testdata/imports).


```yaml
imports:
- env: ./env.short.yaml
- server_deployment: ./deployment_template.short.yaml
  params:
    name: my-koki-short-server
    env: ${env}
deployment: ${server_deployment}
```

```yaml
# deployment_template.short.yaml
imports:
- default_env: ./dummy_env.short.yaml
params:
- name: for selector/labels and name
  default: koki-short-server
- env: env-specific config
  default: ${default_env}
deployment:
  name: ${name}
  labels:
    app: ${name}
  replicas: 3
  selector:
    app: ${name}
  affinity:
  - anti_pod: app=${name}
    topology: kubernetes.io/hostname
  containers:
  - image: ublubu/testserver:${env.server_version}
    name: ${name}
    pull: always
    expose:
    - 8080
    env:
    - COOKIE_AUTH_KEY=${env.cookie_auth_key}
    - STRIPE_KEY=${env.stripe_key}
    - GITHUB_CLIENT_ID=${env.github_client_id}
    - GITHUB_CLIENT_SECRET=${env.github_client_secret}
```

```yaml
# env.short.yaml / dummy_env.short.yaml
env:
  server_version: v1.0
  cookie_auth_key: COOKIEAUTHKEY12345
  stripe_key: STRIPEKEY12345
  github_client_id: 1234567
  github_client_secret: GITHUBSECRET1234567890
```

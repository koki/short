# Modules

Modules facilitate the reuse and composition of Koki resources.


```yaml
imports:
- env: ./env.short.yaml
params:
- name: for selector/labels and name
  default: koki-short-server
exports:
- default: a koki short server deployment
  value:
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
        ports:
        - 8080
        env:
        - COOKIE_AUTH_KEY=${env.cookie_auth_key}
        - STRIPE_KEY=${env.stripe_key}
        - GITHUB_CLIENT_ID=${env.github_client_id}
        - GITHUB_CLIENT_SECRET=${env.github_client_secret}
```

A koki module consists of three sections: `imports`, `params`, and `exports`.

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

## Exports

Exports let you define the resource manifests created by your module.

```yaml
exports:
- default: description of the 'default' export. if this module is imported as 'foo', use 'default' as ${foo}.
  value: ...
```

NOTE: The `default` export only gets special treatment if your module exports one resource.
  If there are multiple exports, use `default` as `${foo.default}`.

```yaml
exports:
- bar: description of the 'bar' export. if the module is imported as 'foo', use 'bar' as ${foo.bar}.
  value: ...
```

## Template Expansion

Koki supports logic-less text templating using `${identifier}`.
These template expansions can be used in two places:

 * the value of any param on an import

```yaml
imports:
- foo: ./foo.yaml
  params:
    x: ${bar}
    y:
    - list item 0
    - ${baz...}
    z:
      key0: item0
      key1: ${qux.export0}
      key2: string with ${qux.export1} inside
```

 * the `value` field of any export

```yaml
exports:
- x: export 'bar' as 'x'
  value: ${bar}
- y: export a list of items that includes the items of 'baz'
  value:
    - list item 0
    - ${baz...}
- z: export a map that contains values from 'qux'
  value:
    key0: item0
    key1: ${qux.export0}
    key2: string with ${qux.export1} inside
```

### Supported Templating Syntax

* Use a param named `foo` with `${foo}`.

* If an import named `foo` contains only a single `default` resource, use `${foo}`

* If an import named `foo` contains multiple resources, use `${foo.export_name}`.

If a template expansion contains a list, you can either use it as a list OR merge its items into an existing list:

```
list: ${identifier_for_list}

# OR

list:
- list item A
- ${identifier_for_list...}
- list item X
```

#### Coming soon:

Index into a value to access specific fields: `${foo.export_name.containers[0]}`

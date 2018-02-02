# Modules

Modular systems are composed of separable pieces (called _modules_) that can evolve independently and be recombined in different ways for new applications.

## Motivation

Why do modules matter?

First, modules reduce duplication. Rather than duplicating code in multiple places, put it in a single shared module.
Wherever you'd copy-paste, simply reference the module.
This approach often yields more maintainable code—instead of changing the same code in multiple files, just change it in the shared module.

Second, modules establish boundaries between different parts of a system.
This allows different members of a team to focus only on the parts that matter to them.
For example, one member of the team might be responsible for a sidecar container module that is used in another team member's Pod templates.

Third, existing modules can be reused and recombined to build new applications.
It's much more costly to create a new application from scratch.

## Overview

A module is meant to be reusable in different contexts.
Depending on the context, the module may need to be configured in different ways.

For example, suppose Application X exposes a service at port `1234`, and Application Y exposes the service at port `5678`.
The same module might be configured with `1234` in Application X and `5678` in Application Y.

This illuminates two concerns for modules:

* __Parameterization__: Configuring a module for different environments
* __Composition__: Incorporating a module into an application

### Parameterization

A parameterized module is essentially a template.

In a template, there are _"holes"_ that need to be filled with actual values.
Modules are the same.
The "holes" are defined as _parameters_.
When the module is used (i.e. _imported_), these parameters are given values—filling the holes in the template.

### Composition

The first step is to _import_ the module you want to use.
If the module has any parameters, assign values to them.
(You have to fill the template before you can use it.)

After importing a module, its contents are available to use.
For example, if the module defines a _container_ resource, you can use that container in a Pod template.

## Koki Short Modules

A Koki Short module contains three sections:

* `imports` - Configure and import __other__ modules for use inside this module.
* `params` - Allow __this__ module to be configured for use in different contexts.
* A section for the resource defined by the module

Before we go deeper, note the three sections in the example file below:

* _Imports_ configuration values from a module called `env.short.yaml`
* Accepts a _parameter_ called `name` to configure its selector, labels, and name
* Defines a `deployment` resource
    - Uses the `name` param for `name`, `selector`, and `labels` fields (and others)
    - Uses the module imported as `env` to configure the Pod's container


```yaml
imports:
- env: ./env.short.yaml
params:
- name: "the value of this param is used for selector, labels, and name"
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

## Params

The `params` section defines the _parameters_ expected by the module.
Parameters are like the _"holes"_ in a template.

Each parameter is defined with a _name_ and a _description_.
If a parameter is optional, its definition also includes a _default value_.


```yaml
params:
- param_foo: "the description for a param named param_foo"
  default: default_value_for_param_foo
- param_bar: "the description for a param named param_bar. it doesn't have a default value."
```

Note that `param_bar` does not have a `default` value.
When this module is imported, `param_bar` must be given a value.

On the other hand, `param_foo` _does_ have a `default` value.
If no value is provided for `param_foo` when this module is imported, `param_foo` will have the value `default_value_for_param_foo`.
However, if a value (e.g. `12`) is provided, then `param_foo` will have that value instead (e.g. `12`).

### Further Info

See the [Imports](#imports) section below to learn how to _set_ parameter values while importing a module.

See the [Templating Syntax](#supported-templating-syntax) section to learn how to _use_ parameters within the module that defines them.

## Imports

The `imports` section allows a module to configure and use other modules.
Each import is defined with a _name_ and a _path_.

* The imported module is loaded from the file at the given _path_.
* The imported module's contents are available under the given _name_. (See the [Templating](#templating) section for more details.)

If the imported module expects parameters, the `import` statement supplies values for these parameters in its `params` field. Note the example below:


```yaml
imports:
- local_name_for_imported_module: ./path/to/imported/module.yaml
  params:
    param_foo: value_for_param_foo
    param_bar: ${interpolation}
- local_name_for_another_imported_module: ./another/path.yaml
```

In this example, two modules are imported.
The first is configured with values for its `param_foo` and `param_bar` parameters.
The second isn't passed any parameters, so it uses default values for any parameters it has.

Both imports in this example load modules using relative paths.
The paths are relative to the directory containing current module (the module that contains the import statements).
Koki Short currently only supports importing from relative paths.

_For information about the `${interpolation}` in the example, see the [Templating](#templating) section._

## Templating

Koki supports logic-free text templating using this pattern: `${some_identifier_here}`

The templating system replaces each instance of the pattern with the value that corresponds to the _identifier_ embedded in the pattern. This process is called _interpolation_.

(See [Supported Templating Syntax](#supported-templating-syntax) for details about syntax.)

Interpolations can be used in two places:


 * the value of any param in an *import* statement (e.g. `x`, `y`, `z` below)

```yaml
imports:
- foo: ./foo.yaml
  params:
    x: ${bar} # "the value of 'x' is the value of 'bar'"
    y:
    - "list item 0"
    - ${baz...} # "'baz' is a list. append its items to 'y'."
    z: ${qux.a.x.0} # "index into 'qux' to get the value for 'z'."
```

 * the contents of the resource defined by the module—e.g. a Pod

```yaml
pod:
  x: ${bar} # "the value of 'pod.x' is the value of 'bar'"
  y: 
  - "list item 0"
  - ${baz...} # "'baz' is a list. insert its items into 'pod.y'."
  z: ${qux.a.x.0} # index into 'qux' to get the value for 'pod.z'

# OR

pod: ${foo} # use 'foo' as 'pod'.
```

### Supported Templating Syntax

Each parameter is defined with a _name_.
Each imported module is also given a _name_.
Koki Short's templating syntax lets us combine these names with operators to select exactly what data is being used in an interpolation.

#### Use the resource (i.e. module) imported as `foo` with `${foo}`:

```yaml
# some_pod.short.yaml

pod:
  name: foo
```

```yaml
# "module that imports some_pod.short.yaml"

imports:
- foo: ./some_pod.short.yaml
pod: ${foo}

---

# after processing:

pod:
  name: "foo"
```

Note that `${foo}` is replaced with `name: "foo"`, not `pod: name: "foo"`.
The top-level key of the resource (e.g. `pod`) is always removed.

Every import name (in a module's `imports` section) can be used for templating in this way.

#### Use a parameter named `foo` with `${foo}`:

```yaml
# "example.short.yaml"

params:
- foo: "description of a parameter named 'foo'"
pod:
  name: ${foo}
```

```yaml
# "module that imports example.short.yaml"

imports:
- example: "./example.short.yaml"
  params:
    foo: bar
pod: ${example}

---

# after processing:

pod:
  name: bar

```

Note that `example.short.yaml` is imported with the parameter `foo: bar`.
This populates the field `name: ${foo}` in `example.short.yaml`, which yields `name: bar`.
(Remember that the top-level key of the resource (e.g. `pod`) is always removed.)

Every parameter name (in a module's `params` section) can be used for templating in this way.

#### Index into `foo` with `${foo.x}`, `${foo.list.0}`, etc:

* Index into a list using `.n` for the item at index n.
    - e.g. `${foo.2}` for the third item of `foo` (index 2).
* Index into a dictionary using `.x` for the item at key x.
    - e.g. `${foo.containers}` for the `containers` field of `foo`.
* Indexes can be chained together.
    - e.g. `${foo.containers.1}` for the second container in `foo`.
    
#### Merge lists using the _spread_ operator: `${foo...}`:

If an identifier represents a list, you can either use it as a list OR merge its items into an existing list. To merge into an existing list, append three periods (called the _spread_ operator) to the identifier:

```yaml
list1:
- list item A
- ${identifier_for_list}
- list item X

# OR

list2:
- list item A
- ${identifier_for_list...}
- list item X
```

```yaml
# if 'identifier_for_list' has this value:

identifier_for_list:
- a
- b
- c

# then after processing, the examples above become:

list1:
- list item A
- # a nested list!
  - a
  - b
  - c
- list item X

list2:
- list item A
- a
- b
- c
- list item X
```

Note that when we don't use the _spread_ operator (in `list1`), so the entire list `[a, b, c]` is added as a single item.
When we use the _spread_ operator (in `list2`), `[a, b, c]` is merged into the list as three separate items.

### Examples

For even more examples, have a look at [these reference examples](https://github.com/koki/short/tree/master/testdata/imports).


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
params:
- name: for selector/labels and name
  default: koki-short-server
- env: env-specific config
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
# env.short.yaml
env:
  server_version: v1.0
  cookie_auth_key: COOKIEAUTHKEY12345
  stripe_key: STRIPEKEY12345
  github_client_id: 1234567
  github_client_secret: GITHUBSECRET1234567890
```

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

Reuse comes from the ability to configure a module by supplying parameter values.
Composition comes from the ability to import other modules into a module.

A Koki Short module contains three sections:

* `imports` - Use other modules inside this module.
* `params` - Allow this module to be configured for use in different contexts.
* A section for the resource defined by the module

Before we go deeper, note the three sections in the example file below.
The module defines a `deployment` resource.
It _imports_ configuration values from a module called `env.short.yaml`
  and accepts a _parameter_ to configure its `name`.


```yaml
imports:
- env: ./env.short.yaml
params:
- name: "used for selector/labels and name"
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

## Params

The `params` section defines the _parameters_ expected by the module.
Each parameter is defined with a _name_ and a _description_.
If a parameter is optional, its definition also includes a _default value_.
(If the user of the module does not provide a value for a parameter, the parameter's default value is used.)


```yaml
params:
- param_foo: "the description for a param named param_foo"
  default: default_value_for_param_foo
- param_bar: "the description for a param named param_bar. it doesn't have a default value."
```

When another module imports this one, it must provide values for the parameters defined in the `params` section.
See the [Imports](#imports) section below for details.

To learn about how parameter values are used within a module, see the [Templating](#templating) section.

## Imports

The `imports` section allows a module to configure and use other modules.
Each import is defined with a _name_ and a _path_.
The imported module is loaded from the file at the given _path_.
Its contents are available under the given _name_. (See the [Templating](#templating) section for more details.)

An import statement may also define its own `params` section,
which configures the imported module by supplying values for its _parameters_.


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
The second isn't passed any parameters, so it uses its default parameter values.

## Templating

Koki supports logic-free text templating using this pattern: `${some_identifier_here}`

The templating system replaces each instance of the pattern with the value that corresponds to the _identifier_ embedded in the pattern. This process is called _interpolation_.

(See [Supported Templating Syntax](#supported-templating-syntax) for details about syntax.)

Interpolations can be used in two places:


 * the value of any param on an *import*

```yaml
imports:
- foo: ./foo.yaml
  params:
    x: ${bar} # "the value of 'x' is the value of 'bar'"
    y
    - "list item 0"
    - ${baz...} # "'baz' is a list. append its items to 'y'."
    z: ${qux.a.x.0} # "index into 'qux' to get the value for 'z'."
```

 * the contents of the defined resource—e.g. a Pod

```yaml
pod:
  x: ${bar} # "the value of 'x' is the value of 'bar'"
  y: 
  - list item 0
  - ${baz...} # "'baz' is a list. insert its items into 'y'."
  z: ${qux.a.x.0} # index into 'qux' to get the value for 'z'

# OR

pod: ${foo} # use 'foo' as 'pod'.
```

### Supported Templating Syntax

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

Every import name (in a module's `imports` section) is a valid _identifier_ for templating.

#### Use a parameter named `foo` with `${foo}`:

```yaml
# "example.short.yaml"

params:
- foo: "a parameter named 'foo'"
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

Note that `example.short.yaml` is imported with the parameter `foo: bar`. This populates the field `name: ${foo}`.

Every parameter name (in a module's `params` section) is a valid _identifier_ for templating.

#### Index into `foo` with `${foo.x}`, `${foo.list.0}`, etc:

* Index into a list using `.n` for the item at index n.
    - e.g. `${foo.2}` for the third item of `foo`(index 2).
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

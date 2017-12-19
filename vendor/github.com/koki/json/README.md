# Koki's JSON library

It's a fork that makes `encoding/json` better for users.

## Features

* Annotates decoder errors with the path (e.g. `$.pod.container.0.env`) where the error occurred. _The golang built-in implementation doesn't contextualize errors, so it's hard to track down where they came from._

* Check for extraneous fields in your JSON, which is useful for detecting misspelled keys and other typos.

* `omitempty` actually omits zero-valued structs.

* Miscellaneous utilities for dealing with JSON represented as `map[string]interface{}`, including an implementation of path-based indexing.

## How to use it

`github.com/koki/json` is a drop-in replacement for `encoding/json`. The decoder implementation is different, but the library API remains unchanged. Use this package's `json.Unmarshal` if you want paths for your decoder errors. Use `json.Marshal` if you want to omit zero-valued structs.

```go
// Deserialize
foo := Foo{}
err := json.Unmarshal(data, &foo)
...

// Serialize
bar := Bar{}
data, err := json.Marshal(bar)
...
```

`github.com/koki/json/jsonutil` is a brand-new package. Use `jsonutil.ExtraneousFieldPaths` to get a list of paths that weren't decoded into your Go object. i.e. potential typos.

```go
// 1. Parse.
obj := map[string]interface{}
parsedObj := Foo{}
err := json.Unmarshal(data, &obj)
...
err = jsonutil.UnmarshalMap(obj, &parsedObj)
...

// 2. Check for unparsed fields--potential typos.
extraneousPaths, err := jsonutil.ExtraneousFieldPaths(obj, parsedObj)
...
if len(extraneousPaths) > 0 {
    return nil, &jsonutil.ExtraneousFieldsError{Paths: extraneousPaths}
}
...
```

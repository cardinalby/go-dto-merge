[![test](https://github.com/cardinalby/go-dto-merge/actions/workflows/test.yml/badge.svg)](https://github.com/cardinalby/go-dto-merge/actions/workflows/test.yml)
[![list](https://github.com/cardinalby/go-dto-merge/actions/workflows/list.yml/badge.svg)](https://github.com/cardinalby/go-dto-merge/actions/workflows/list.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cardinalby/go-dto-merge.svg)](https://pkg.go.dev/github.com/cardinalby/go-dto-merge)

# dtomerge package

This package provides a way to deep merge two structs of the same type.

It's useful for merging **default** values with **user-provided** values (e.g. configs).
Values that are not present in the user-provided struct will be taken from the default struct.

**Only exported** fields are merged.

## Example
```
go get github.com/cardinalby/go-dto-merge
```

Example `Config` struct has a nested `UserConfig` struct. 
```go

import "github.com/cardinalby/go-dto-merge"

type UserConfig struct {
    Role string      // for non-pointer fields zero value indicates it's not specified
    Name string
}

type Config struct {
    Verbose *bool    // it is a pointer to distinguish between "not specified" and false
    User UserConfig
}
```

Given defaults, we can merge them with user-provided values:
```go
// ptr is some helper function to create a pointer to a value
defaults := Config{     // it's called "src"
    Verbose: ptr(true),
    User: UserConfig{
        Role: "admin",
        Name: "John",
    },
}
userProvided := Config{  // it's called "patch"
    User: UserConfig{
        Name: "Jane",
    },
}
res, err := dtomerge.Merge(defaults, userProvided) // (src, patch)

// res == Config{
//     Verbose: (*bool) true,
//     User: UserConfig{
//         Role: "admin",
//         Name: "Jane",
//     },
// }
```

## Pointers
Pointers can be used to distinguish between "not specified" and "explicit zero value" fields.
- If `patch` field contains a nil pointer, it will not override `defaults.Verbose`
- If `patch` field contains a pointer to zero value, it will override `src` field
  only in case `src` pointer field is nil (value will be copied)
- If `patch` field contains a pointer tp **non**-zero value, it will override `src` field (value will be copied)

Use `dtomerge.OptDeRefPointers(false)` option to handle pointers as regular fields.

## Slices and maps
Setting additional option you can merge slices and maps as well.

```go

import (
    "github.com/cardinalby/go-dto-merge"
    "github.com/cardinalby/go-dto-merge/opt"
)

type Config struct {    
    Roles []string
    Permissions map[string]bool
}

defaults := Config{
    Roles: []string{"admin", "user"},
    Permissions: map[string]bool{
        "read": true,
        "write": false,
    },
}

userProvided := Config{
    Roles: []string{"user", "guest"},
    Permissions: map[string]bool{
        "write": true,
    },
}

res, err := dtomerge.Merge(defaults, userProvided,
    // merge map keys
    dtomerge.OptIterateMaps(true),
    // merge slices as unique sets
    dtomerge.OptMergeSlices(dtomerge.SlicesMergeStrategyUnique), 		
)

// res == Config{
//     Roles: []string{"admin", "user", "guest"},
//     Permissions: map[string]bool{
//         "read": true,
//         "write": true,
//     },
// }
```

Possible `MergeSlices` strategies:
- `dtomerge.SlicesMergeStrategyUnique`: `[1, 2, 3] + [4, 2, 1] → [1, 2, 3, 4]`
- `dtomerge.SlicesMergeStrategyByIndex`: `[1, 2, 3] + [11, 12] → [11, 12, 3]`
- `dtomerge.SlicesMergeStrategyAtomic`: merge as a whole, default

## Other options
You can: 
- specify a custom merge function for a specific type.
- specify a custom merge options for a specific type.

See [options godoc](https://pkg.go.dev/github.com/cardinalby/go-dto-merge#Options) for more details.



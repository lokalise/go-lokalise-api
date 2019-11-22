# Lokalise API v2 official Golang client library
[![GoDoc](https://godoc.org/github.com/lokalise/go-lokalise-api?status.svg)](https://godoc.org/github.com/lokalise/go-lokalise-api)
![Build status](https://github.com/lokalise/go-lokalise-api/workflows/tests/badge.svg)
[![Test Coverage](https://codecov.io/gh/lokalise/go-lokalise-api/branch/feature%2Ftests/graph/badge.svg)](https://codecov.io/gh/lokalise/go-lokalise-api)

# Index

* [Getting started](#getting-started)
  + [Installation and Usage](#installation-and-usage)
  + [Initializing the Client](#initializing-the-client)
  + [General options](#general-options)
  + [Objects and models](#objects-and-models)
  + [Request options and pagination](#Request-options-and-pagination)
* [Available Resources](#available-resources)
  + [Comments](#comments)
  + [Contributors](#contributors)
  
# Getting Started

## Usage

```go
import "github.com/lokalise/go-lokalise-api/v2"	// with go modules enabled (GO111MODULE=on or outside GOPATH)
import "github.com/lokalise/go-lokalise-api" // with go modules disabled
```

## Initializing the client

```go

token := os.Getenv("lokalise_token")
client, err := lokalise.New(token)

```

## General options

You can set global API parameters with the ClientOption functions during the initialization. The following functions are available:

* WithBaseURL
* WithRetryCount
* WithRetryTimeout
* WithConnectionTimeout
* WithDebug
* WithPageLimit 

Usage:

```go
Api, err := lokalise.New(
    "token-string",
    lokalise.WithDebug(true),
    ...
)
```

## Objects and models

Individual objects are represented as instances of according structs. Different objects are used for creating and updating in most cases.

Here are some object types:

* Create/Update request objects, i.e. NewKey, NewContributor etc

* Response objects: single/multiple, i.e. KeyResponse/KeysResponse and special , i.e. DeleteKeyResponse. There is no separate ErrorResponse - errors are encapsulated into concrete method's response. 

* List options that are used for sending certain options and pagination, i.e. KeyListOptions.

## Request options and pagination

Some resources, such as Projects, Keys, Files, Tasks, Screenshots, Translations have optional parameters for List method (Keys also have an option for Retrieve). These parameters should be set before calling.

All request options could be set inline and separately:

```go

// separately:
keys := client.Keys()
keys.SetListOptions(lokalise.KeyListOptions{
    IncludeTranslations: 1,
    IncludeComments: 1,
})

resp, err := keys.List("{PROJECT_ID}")

// inline:
client.Keys().WithListOptions(lokalise.KeyListOptions{Limit: 3}).List("{PROJECT_ID}")

```

There are two parameters, used for pagination: Limit and Page.

```go
t := Api.Teams()
t.SetPageOptions(lokalise.PageOptions{
    Page:  2,
    Limit: 10,
})

resp, err := t.List()
```

# Available resources

## Comments

### List project comments

```go
projectId := "aaaabbbb.cccc"
c := Api.Comments()
c.SetPageOptions(lokalise.PageOptions{Page: 1, Limit: 20})
resp, err := c.ListProject(projectId)

```

### List key comments

```go
projectId := "aaaabbbb.cccc"
keyId := 26835183
c := Api.Comments()
c.SetPageOptions(lokalise.PageOptions{Page: 1, Limit: 20})
resp, err := c.ListByKey(projectId, keyId)
```

### Create

```go
c := lokalise.NewComment{Comment: "My new comment"}
resp, err := Api.Comments().Create(projectId, keyId, []lokalise.NewComment{c})
```

### Retrieve

```go
...
commentId := 26835183
resp, err := Api.Comments().Retrieve(projectId, keyId, commentId)
```

### Delete

```go
...
resp, err := Api.Comments().Delete(projectId, keyId, commentId)
```

## Contributors

### List all contributors

```go
projectId := "aaaabbbb.cccc"
pageOpts := lokalise.PageOptions{Page: 1, Limit: 20}

c := Api.Contributors()
c.SetPageOptions(pageOpts)
c.List(projectId)
```

### Create contributors

```go
contributorCreate := lokalise.NewContributor{
    Email:    "some@ema.il",
    Fullname: "New contributor",
    Permission: lokalise.Permission{
        IsAdmin:     true,
        IsReviewer:  true,
        Languages:   []lokalise.Language{{LangISO: "en", IsWritable: true}},
        AdminRights: []string{"upload", "download"},
    },
}
resp, err := Api.Contributors().Create(projectId, []lokalise.NewContributor{contributorCreate})
```

### Retrieve contributor

```go
userId := 47913 
resp, err := Api.Contributors().Retrieve(projectId, userId)
```

### Update contributor

```go
...
permissionUpdate := lokalise.Permission{
    IsReviewer: true,
    IsAdmin: false,
    AdminRights: []string{"keys", "upload", "download"},
}
resp, err := Api.Contributors().Update(projectId, userId, permissionUpdate)
```

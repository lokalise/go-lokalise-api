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
  + [Rate limits](#rate-limits)
* [Available Resources](#available-resources)
  + [Comments](#comments)
  + [Contributors](#contributors)
  
# Getting Started

## Usage

```go
import "github.com/lokalise/go-lokalise-api/v5"	// with go modules enabled (GO111MODULE=on or outside GOPATH)
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

### Cursor pagination

The [List Keys](https://developers.lokalise.com/reference/list-all-keys) and [List Translations](https://developers.lokalise.com/reference/list-all-translations) endpoints support cursor pagination, which is recommended for its faster performance compared to traditional "offset" pagination. By default, "offset" pagination is used, so you must explicitly set `pagination` to `"cursor"` to use cursor pagination.

```go
// This approach is also applicable for `client.Translations()`
keys := Api.Keys()
keys.SetPageOptions(lokalise.PageOptions{
  Pagination: "cursor",
  Cursor: "eyIxIjo1MjcyNjU2MTd9"
})

resp, err := keys.List()
```

After retrieving data from the Lokalise API, you can check for the availability of the next cursor and proceed accordingly:

```go
cursor := ""

for {
	keys := client.Keys()
	keys.SetListOptions(KeyListOptions{
		Pagination: "cursor",
		Cursor:     cursor,
	})
	resp, _ := keys.List(projectId)
	
	// Do something with the response
	
	if !resp.Paged.HasNextCursor() {
		// no more keys
		break
	}
	cursor = resp.Paged.Cursor
}
```

## Queued Processes
Some resource actions, such as Files.upload, are subject to intensive processing before request fulfills. 
These processes got optimised by becoming asynchronous.
The initial request only queues the data for processing and retrieves to queued process identifier.
Additional request to QueuedProcesses resource could be executed to obtain the current processing result.

Example with Files.upload:
```go
projectId := "aaaabbbb.cccc"
uploadOpts := lokalise.FileUpload{
    Filename: "test.html",
    LangISO:  "en"
}

f := Api.Files()
resp, err = f.Upload(projectId, uploadOpts)
```

The successful response will contain process ID, which can be used to obtain the final result:

```go
projectId := "aaaabbbb.cccc"
processId := "ddddeeeeeffff"

q := Api.QueuedProcesses()
resp, err := q.Retrieve(projectId, processId)
```

## Rate limits
[Access to all endpoints is limited](https://app.lokalise.com/api2docs/curl/#resource-rate-limits) to 6 requests per second from 14 September, 2021. This limit is applied per API token and per IP address. If you exceed the limit, a 429 HTTP status code will be returned and the corresponding exception will be raised that you should handle properly. To handle such errors, we recommend an exponential backoff mechanism with a limited number of retries.

Only one concurrent request per token is allowed.


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

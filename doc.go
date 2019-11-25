/*
Package lokalise provides functions to access the Lokalise web API.

Usage:

	import "github.com/lokalise/go-lokalise-api/v2"	// with go modules enabled (GO111MODULE=on or outside GOPATH)
	import "github.com/lokalise/go-lokalise-api" // with go modules disabled

Initializing the client

	token := os.Getenv("lokalise_token")
	client, err := lokalise.New(token)

General options

You can set global API parameters with the ClientOption functions during the initialization. The following functions are available:

* WithBaseURL

* WithRetryCount

* WithRetryTimeout

* WithConnectionTimeout

* WithDebug

* WithPageLimit

Usage:

	Api, err := lokalise.New(
		"token-string",
		lokalise.WithDebug(true),
		...
	)

Objects and models

Individual objects are represented as instances of according structs. Different objects are used for creating and updating in most cases.

Here are some object types:

* Create/Update request objects, i.e. NewKey, NewContributor etc

* Response objects: single/multiple, i.e. KeyResponse/KeysResponse and special , i.e. DeleteKeyResponse. There is no separate ErrorResponse - errors are encapsulated into concrete method's response.

* List options that are used for sending certain options and pagination, i.e. KeyListOptions.

Request options and pagination

Some resources, such as Projects, Keys, Files, Tasks, Screenshots, Translations have optional parameters for List method (Keys also have an option for Retrieve). These parameters should be set before calling.

All request options could be set inline and separately:

	// separately:
	keys := client.Keys()
	keys.SetListOptions(lokalise.KeyListOptions{
		IncludeTranslations: 1,
		IncludeComments: 1,
	})

	resp, err := keys.List("{PROJECT_ID}")

	// inline:
	client.Keys().WithListOptions(lokalise.KeyListOptions{Limit: 3}).List("{PROJECT_ID}")

There are two parameters, used for pagination: Limit and Page.

	t := Api.Teams()
	t.SetPageOptions(lokalise.PageOptions{
		Page:  2,
		Limit: 10,
	})

	resp, err := t.List()

*/
package lokalise

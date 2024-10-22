# Goforce

Goforce is a easy to use Salesforce Client written in Go and modeled after the Go Elasticsearch library. Some inspiration was drawn from the simpleforce library here: https://github.com/simpleforce/simpleforce

## Features

- Execute SOQL queries
    - Raw Queries
    - Structured Queries (ORM)
- Get records by Type & Id
- Execute SOSL Parameterized Search

Most of the implementation referenced Salesforce documentation here: https://developer.salesforce.com/docs/atlas.en-us.214.0.api_rest.meta/api_rest/intro_what_is_rest_api.htm

## Installation

`goforce` can be acquired as any other Go libraries via `go get`:

```
go get github.com/0xArch3r/goforce
```

## Quick Start

### Setup the Client

In order to utilize goforce, you need to instatiate the client. It is as simple as calling the goforce.NewClient() function. You can alter client config by using the functional options provided in options.go. Before you can use the client, you must "Login" using the client.LoginPassword function and providing it credentials.  

```go
package main

import "github.com/0xArch3r/goforce"

var (
	sfURL      = "Custom or instance URL, for example, 'https://my.salesforce.com'"
	sfUser     = "Username of the Salesforce account."
	sfPassword = "Password of the Salesforce account."
	sfToken    = "Security token, could be omitted if Trusted IP is configured."
)

func main() {
    client, err := goforce.NewClient(
        goforce.WithUrl(sfURL),
    )
    if err != nil {
        panic(err)
    }

    err = client.LoginPassword(sfUser, sfPassword, "")
    if err != nil {
        panic(err)
    }
}
```

### Fetch a User Record by Id

```go
package main

import "github.com/0xArch3r/goforce"

var (
	sfURL      = "Custom or instance URL, for example, 'https://my.salesforce.com'"
	sfUser     = "Username of the Salesforce account."
	sfPassword = "Password of the Salesforce account."
	sfToken    = "Security token, could be omitted if Trusted IP is configured."
)

func main() {
    client, err := goforce.NewClient(
        goforce.WithUrl(sfURL),
    )
    if err != nil {
        panic(err)
    }

    err = client.LoginPassword(sfUser, sfPassword, "")
    if err != nil {
        panic(err)
    }

    user, err := client.Get("User", "SomeID")
}
```

### Execute a SELECT SOQL Query

The `client` provides mutliple ways to perform a SOQL. For Basic queries, you can utilize the Select Query method.

```go

results, err := client.Query.Select(
    "User", 
    c.Search.SObjects(userObj), c.Search.WithContext(context.Background())
)

```

### Execute a Raw SOQL Query

The `client` provides mutliple ways to perform a SOQL. For advanced users who have familiarity with SOQL, you can perform a Raw Query.

```go

results, err := client.Query.Raw("SELECT Id, Name FROM User WHERE Name = 'John Doe'")

```

### Execute a SOSL Parameterized Search

The `client` provides a way to perform a parameterized search using the Search method.

```go

userObj := api.SearchObject{
    Name: "User",
}

results, err := client.Search(
    "John Doe", 
    c.Search.Fields("Id", "Name"),
    c.Search.SObjects(userObj),
)

```

## License and Acknowledgement

This package is released under BSD license. Part of the code referenced the simpleforce
(https://github.com/simpleorce/simpleforce) project.

Contributions are welcome!

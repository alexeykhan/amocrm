![Go Package for amoCRM API](logo.png?raw=true)

# Go Package for amoCRM API 

<p>
    <a href="https://pkg.go.dev/github.com/alexeykhan/amocrm">
        <img src="https://img.shields.io/badge/pkg.go.dev-reference-00ADD8?logo=go&logoColor=white" alt="GoDoc Reference">
    </a>
    <a href="https://github.com/alexeykhan/amocrm">
        <img src="https://img.shields.io/badge/codecov-98%25-success?logo=codecov&logoColor=white" alt="Code Coverage">
    </a>
    <a href="https://github.com/alexeykhan/amocrm">
        <img src="https://img.shields.io/badge/build-passes-success?logo=travis-ci&logoColor=white" alt="Build Status">
    </a>
    <a href="https://pkg.go.dev/github.com/alexeykhan/amocrm">
        <img src="https://img.shields.io/badge/licence-MIT-success" alt="License">
    </a>
</p> 

This package provides a Golang client for amoCRM API.


## Disclaimer

This code is in no way affiliated with, authorized, maintained, sponsored 
or endorsed by amoCRM or any of its affiliates or subsidiaries. This is an 
independent and unofficial API client. Use at your own risk.

## Installation

`go get -u github.com/alexeykhan/amocrm`

## Quick Start

**Step №1: Redirect user to the authorization page.**

User grants access to their account and is redirected back 
with `referer` and `code` GET-parameters attached.

```go
package main

import (
    "fmt"

    "github.com/alexeykhan/amocrm"
)

func main() {
    amoCRM := amocrm.New("clientID", "clientSecret", "redirectURL")
    
    state := amocrm.RandomState()  // store this state as a session identifier
    mode := amocrm.PostMessageMode // options: PostMessageMode, PopupMode
    
    authURL, err := amoCRM.AuthorizeURL(state, mode)
    if err != nil {
        fmt.Println("failed to get auth url:", err)
        return
    }
    
    fmt.Println("Redirect user to this URL:")
    fmt.Println(authURL)
}
```

**Step №2: Exchange authorization code for token.**

Use received `referer` and `code` parameters as account domain and
authorization code respectively to make a handshake with amoCRM and
Get a fresh set of `access_token`, `refresh_token` and token meta data. 

```go
package main

import (
    "fmt"

    "github.com/alexeykhan/amocrm"
)

func main() {
    amoCRM := amocrm.New("clientID", "clientSecret", "redirectURL")
    
    if err := amoCRM.SetDomain("example.amocrm.ru"); err != nil {
        fmt.Println("set domain:", err)
        return
    }
    
    token, err := amoCRM.TokenByCode("authorizationCode")
    if err != nil {
        fmt.Println("get token by code:", err)
        return
    }
    
    fmt.Println("access_token:", token.AccessToken())
    fmt.Println("refresh_token:", token.RefreshToken())
    fmt.Println("token_type:", token.TokenType())
    fmt.Println("expires_at:", token.ExpiresAt().Unix())
}
```

**Step №3: Make your first API request.**

Set amoCRM accounts domain and token to authorize your requests.

```go
package main

import (
    "fmt"
    "time"

    "github.com/alexeykhan/amocrm"
)

func main() {
    amoCRM := amocrm.New("clientID", "clientSecret", "redirectURL")
    
    if err := amoCRM.SetDomain("example.amocrm.ru"); err != nil {
        fmt.Println("set domain:", err)
        return
    }
    
    token := amocrm.NewToken("accessToken", "refreshToken", "tokenType", time.Now())
    if err := amoCRM.SetToken(token); err != nil {
        fmt.Println("set token:", err)
        return
    }
    
    cfg := amocrm.AccountsConfig{
        Relations: []string{
            amocrm.WithUUID, 
            amocrm.WithVersion, 
            amocrm.WithAmojoID,
            amocrm.WithTaskTypes,
            amocrm.WithUserGroups,
            amocrm.WithAmojoRights,
            amocrm.WithDatetimeSettings,
        }, 
    }

    account, err := amoCRM.Accounts().Current(cfg)
    if err != nil {
        fmt.Println("fetch current accounts:", err)
        return
    }
    
    fmt.Println("current accounts:", account)
}
```

## Development Status: In Progress

This package is under development so any methods, constants or types may be changed 
in newer versions without backward compatibility with previous ones. Use it at your
own risk and feel free to fork it anytime.

<hr>

Released under the [MIT License](LICENSE.md).
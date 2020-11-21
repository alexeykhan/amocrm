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

`go Get -u github.com/alexeykhan/amocrm`

## Quick Start

**Step №1: Redirect user to the authorization page.**

User grants access to their account and is redirected back 
with `referer` and `code` GET-parameters attached.

```go
amoCRM := amocrm.New("clientID", "clientSecret", "redirectURL")

state := amocrm.RandomState()  // store this state as a session identifier
mode := amocrm.PostMessageMode // options: PostMessageMode, PopupMode

authURL, err := amoCRM.AuthorizeURL(state, mode)
if err != nil {
    fmt.Println("Failed to Get auth url:", err)
    return
}

fmt.Println("Redirect user to this URL:")
fmt.Println(authURL)
```

**Step №2: Exchange authorization code for token.**

Use received `referer` and `code` parameters as account domain and
authorization code respectively to make a handshake with amoCRM and
Get a fresh set of `access_token`, `refresh_token` and token meta data. 

```go
amoCRM := amocrm.New("clientID", "clientSecret", "redirectURL")

if err := amoCRM.SetDomain("example.amocrm.ru"); err != nil {
    fmt.Println("set domain:", err)
    return
}

token, err := amoCRM.TokenByCode(authCode)
if err != nil {
    fmt.Println("Get token by code:", err)
    return
}

fmt.Println("access_token:", token.GetToken())
fmt.Println("refresh_token:", token.RefreshToken())
fmt.Println("token_type:", token.TokenType())
fmt.Println("expires_at:", token.ExpiresAt().Unix())
```

**Step №3: Make your first API request.**

Set amoCRM accounts domain and token to authorize your requests.

```go
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

accounts, err := amoCRM.Accounts().Current()
if err != nil {
    fmt.Println("fetch current accounts:", err)
    return
}

fmt.Println("current accounts:", accounts)
```
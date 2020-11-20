![Go Package for amoCRM API](.github/logo.png?raw=true)

# Go Package for amoCRM API 

<p>
    <a href="https://pkg.go.dev/github.com/alexeykhan/amocrm">
        <img src="https://img.shields.io/badge/pkg.go.dev-reference-00ADD8?logo=go&logoColor=white" alt="GoDoc Reference" style="max-width:100%;">
    </a>
    <a href="https://pkg.go.dev/github.com/alexeykhan/amocrm">
        <img src="https://img.shields.io/badge/codecov-98%25-success?logo=codecov&logoColor=white" alt="Code Coverage">
    </a>
    <a href="https://pkg.go.dev/github.com/alexeykhan/amocrm">
        <img src="https://img.shields.io/badge/build-passes-success?logo=travis-ci&logoColor=white" alt="Build Status">
    </a>
    <a href="https://pkg.go.dev/github.com/alexeykhan/amocrm">
        <img src="https://img.shields.io/badge/licence-MIT-success" alt="License">
    </a>
</p> 

This package provides a Golang client for amoCRM API.

## Installation

`go get -u github.com/alexeykhan/amocrm`

## Quick Start

*Step №1: Redirect user to the authorization page.* 
User grants access to their account and is redirected back to 
`redirect_url` with `referer` and `code` GET-parametetes attached.

```go
// Initialize amoCRM API Client.
amoCRM := amocrm.New(env.clientID, env.clientSecret, env.redirectURL)

// Save this random state as a session identifier to verify
// user identity when they are redirected back with code.
// Set required mode parameter: "post_message" or "popup".
state := amocrm.RandomState()
mode := "post_message"

// Redirect user to authorization URL.
authURL, err := amoCRM.AuthorizeURL(state, mode)
if err != nil {
    fmt.Println("Failed to get auth url:", err)
    return
}

fmt.Println("Redirect user to this URL:")
fmt.Println(authURL)
```

*Step №2: Exchange authorization code for token.* 
Use received `referer` and `code` parameters as account domain and
authorization code respectively to make a handshake with amoCRM and
get a fresh set of `access_token`, `refresh_token` and token meta data. 

```go
// Initialize amoCRM API Client.
amoCRM := amocrm.New(env.clientID, env.clientSecret, env.redirectURL)

// Use the account domain and authorization code that are
// pushed to the redirect URL as "referer" and "code GET
// parameters respectively. AccessTokenByCode will do the
// handshake to retrieve tokens.
domain := "example.amocrm.ru"
authCode := "def502000ba3e1724cac79...92146f93b70fd4ca31"

// Set amoCRM API account domain.
if err := amoCRM.SetDomain(domain); err != nil {
    fmt.Println("set domain:", err)
    return
}

// Exchange authorization code for token.
token, err := amoCRM.TokenByCode(authCode)
if err != nil {
    fmt.Println("get token by code:", err)
    return
}

// Store received token.
fmt.Println("access_token:", token.AccessToken())
fmt.Println("refresh_token:", token.RefreshToken())
fmt.Println("token_type:", token.Type())
fmt.Println("expires_at:", token.ExpiresAt().Unix())
```

*Step №3: Make your first request to amoCRM API.* 
Set amoCRM account domain and token to authorize your requests.

```go
// Initialize amoCRM API Client.
amoCRM := amocrm.New(env.clientID, env.clientSecret, env.redirectURL)

// Retrieve domain from storage.
if err := amoCRM.SetDomain(storage.domain); err != nil {
    fmt.Println("set domain:", err)
    return
}

// Retrieve token from storage.
token := amocrm.NewToken(storage.accessToken, storage.refreshToken, storage.tokenType, storage.expiresAt)
if err := amoCRM.SetToken(token); err != nil {
    fmt.Println("set token:", err)
    return
}

// Fetch current account from API.
account, err := amoCRM.Account().Current()
if err != nil {
    fmt.Println("fetch current account:", err)
    return
}

fmt.Println("current account:", account)
```
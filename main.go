// Copyright (c) 2020 Alexey Khan
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/alexeykhan/amocrm/oauth2"
)

type APIClient struct {
	oAuth2Client oauth2.Client
}

func (c APIClient) Oauth2Client() oauth2.Client {
	return c.oAuth2Client
}

func New(clientID, clientSecret, redirectURL string) APIClient {
	return APIClient{
		oAuth2Client: oauth2.New(clientID, clientSecret, redirectURL),
	}
}

func main() {
	ctx := context.Background()

	api := New(
		os.Getenv("CLIENT_ID"),
		os.Getenv("CLIENT_SECRET"),
		os.Getenv("REDIRECT_URI"))

	client := api.Oauth2Client()

	url, _ := client.AuthorizeURL(oauth2.GenerateState(), oauth2.PostMessageMode)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// НА САМОМ ДЕЛЕ ЗДЕСЬ ЕЩЕ ПРИДЕТСЯ ВЫТАЩИТЬ ИЗ РЕДИРЕКТ УРЛА
	// REFERRER, ЧТОБЫ СДЕЛАТЬ ИЗ НЕГО НОВЫЙ DOMAIN
	client.SetAccountDomain("getmetrics.amocrm.ru")

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	token, err := client.AccessTokenByCode(ctx, code)
	if err != nil {
		fmt.Println("error: ", err)
	}

	fmt.Printf("token: %+v", *token)
}

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

package amocrm_test

import (
	"fmt"
	"time"

	"github.com/alexeykhan/amocrm"
	"github.com/alexeykhan/amocrm/oauth2"
)

var (
	clientID     = "CLIENT_ID"
	clientSecret = "CLIENT_SECRET"
	redirectURL  = "REDIRECT_URI"
)

func ExampleAmoCRM_Account() {
	// Create amoCRM API client.
	api := amocrm.New(clientID, clientSecret, redirectURL)

	token := oauth2.NewToken("access", "refresh", "BEARER", time.Now())

	api.SetToken(token)

	a, err := api.Account()
	if err != nil {
		fmt.Println("get account:", err)
		return
	}

	fmt.Println("account:", a)

	// Output:
	// jhkjhj
}
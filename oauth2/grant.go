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

package oauth2

type GrantType interface {
	Name() string
	Parameters() []string
}

type Grant struct {
	name       string
	parameters []string
}

func (g Grant) Name() string {
	return g.name
}

func (g Grant) Parameters() []string {
	return g.parameters
}

func AuthorizationCodeGrant() Grant {
	return Grant{
		name:       "authorization_code",
		parameters: []string{"code"},
	}
}

func RefreshTokenGrant() Grant {
	return Grant{
		name:       "refresh_token",
		parameters: []string{"refresh_token"},
	}
}

func ClientCredentialsGrant() Grant {
	return Grant{
		name:       "client_credentials",
		parameters: []string{},
	}
}

func PasswordGrant() Grant {
	return Grant{
		name:       "password",
		parameters: []string{"username", "password"},
	}
}

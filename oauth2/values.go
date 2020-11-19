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

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
)

const (
	PostMessageMode = "post_message"
	PopupMode       = "popup"

	amoCRM = "amocrm"
)

type Domain struct {
	account string
	zone    string
}

func (d Domain) String() string {
	return d.account + "." + amoCRM + "." + d.zone
}

type RedirectMode struct {
	value string
}

func (m RedirectMode) Value() string {
	return m.value
}

func NewRedirectMode(name string) (*RedirectMode, error) {
	if name != PostMessageMode && name != PopupMode {
		return nil, errors.New("unexpected redirect mode")
	}

	return &RedirectMode{value: name}, nil
}

func NewDomain(domain string) (*Domain, error) {
	if domain == "" {
		return nil, errors.New("empty domain")
	}

	parts := strings.Split(domain, ".")
	if len(parts) != 3 || parts[0] == "" || parts[1] != amoCRM || !zone(parts[2]) {
		return nil, errors.New("invalid domain format")
	}

	return &Domain{account: parts[0], zone: parts[2]}, nil
}

func NewState() string {
	// Converting bytes to hex will always double length. Hence, we can reduce
	// the amount of bytes by half to produce the correct length of 32 characters.
	key := make([]byte, 16)

	// https://golang.org/pkg/math/rand/#Rand.Read
	// Ignore errors as it always returns a nil error.
	_, _ = rand.Read(key)

	return hex.EncodeToString(key)
}

func zone(name string) bool {
	zones := []string{"ru", "com"}
	for _, z := range zones {
		if z == name {
			return true
		}
	}
	return false
}

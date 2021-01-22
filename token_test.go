// The MIT License (MIT)
//
// Copyright (c) 2021 Alexey Khan
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
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alexeykhan/amocrm"
)

var (
	accessToken  = "access_token"
	refreshToken = "refresh_token"
	tokenType    = "bearer"
	expiresAt    = time.Now()
)

func TestNewToken(t *testing.T) {
	token := amocrm.NewToken(accessToken, refreshToken, tokenType, expiresAt)
	require.Implements(t, (*amocrm.Token)(nil), token)
}

func TestTokenSource_AccessToken(t *testing.T) {
	token := amocrm.NewToken(accessToken, refreshToken, tokenType, expiresAt)
	require.Exactly(t, accessToken, token.AccessToken())
}

func TestTokenSource_RefreshToken(t *testing.T) {
	token := amocrm.NewToken(accessToken, refreshToken, tokenType, expiresAt)
	require.Exactly(t, refreshToken, token.RefreshToken())
}

func TestTokenSource_TokenType(t *testing.T) {
	cases := []struct {
		typeCode  string
		typeValue string
	}{
		{typeCode: "bearer", typeValue: "Bearer"},
		{typeCode: "mac", typeValue: "MAC"},
		{typeCode: "basic", typeValue: "Basic"},
		{typeCode: "example", typeValue: "example"},
	}
	for _, tc := range cases {
		token := amocrm.NewToken(accessToken, refreshToken, tc.typeCode, expiresAt)
		require.Exactly(t, tc.typeValue, token.TokenType())
	}
}

func TestTokenSource_ExpiresAt(t *testing.T) {
	token := amocrm.NewToken(accessToken, refreshToken, tokenType, expiresAt)
	require.Exactly(t, expiresAt, token.ExpiresAt())
}

func TestTokenSource_Expired(t *testing.T) {
	token := amocrm.NewToken(accessToken, refreshToken, tokenType, expiresAt)
	require.True(t, token.Expired())
}

func TestTokenSource_Expired_Limitless(t *testing.T) {
	token := amocrm.NewToken(accessToken, refreshToken, tokenType, time.Time{})
	require.False(t, token.Expired())
}

func TestTokenSource_Expired_EmptyAccessToken(t *testing.T) {
	token := amocrm.NewToken("", refreshToken, tokenType, expiresAt)
	require.True(t, token.Expired())
}

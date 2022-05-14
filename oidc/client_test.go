package oidc

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOidcClient_AuthUrl(t *testing.T) {
	patterns := []struct {
		desc        string
		client      IClient
		respType    string
		scopes      []string
		redirectUrl string
		state       string
		expected    string
	}{
		{
			"",
			NewGoogleOidcClient(),
			"code",
			[]string{"openid", "email", "profile"},
			"http://localhost:8000/auth/google/sign_up/callback",
			"12345678",
			fmt.Sprintf(
				"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&response_type=%s&scope=%s&redirect_uri=%s&state=%s",
				"",
				"code",
				"openid%20email%20profile",
				"http://localhost:8000/auth/google/sign_up/callback",
				"12345678",
			),
		},
		{
			"scopeが一個の時にエラーにならないか",
			NewGoogleOidcClient(),
			"code",
			[]string{"profile"},
			"http://localhost:8000/auth/google/sign_up/callback",
			"12345678",
			fmt.Sprintf(
				"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&response_type=%s&scope=%s&redirect_uri=%s&state=%s",
				"",
				"code",
				"profile",
				"http://localhost:8000/auth/google/sign_up/callback",
				"12345678",
			),
		},
	}

	for _, pattern := range patterns {
		actual := pattern.client.AuthUrl(
			pattern.respType,
			pattern.scopes,
			pattern.redirectUrl,
			pattern.state,
		)
		assert.Equal(t, pattern.expected, actual)
	}
}

func TestOidcClient_PostTokenEndpoint(t *testing.T) {
	client := NewGoogleOidcClient()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://oauth2.googleapis.com/token",
		httpmock.NewStringResponder(
			200,
			`{
     "access_token": "DummyAccessToken",
     "expires_in": 3566,
     "scope": "openid https://www.googleapis.com/auth/userinfo.email",
     "token_type": "Bearer",
     "id_token": "DummyIdToken"
 }
`,
		),
	)

	actual, _ := client.PostTokenEndpoint("", "", "")
	expected := tokenResponse{
		AccessToken: "DummyAccessToken",
		ExpiresIn:   3566,
		Scope:       "openid https://www.googleapis.com/auth/userinfo.email",
		TokenType:   "Bearer",
		IdToken:     "DummyIdToken",
	}

	assert.Equal(t, expected, actual)
}

func TestClient_PostUserInfoEndpoint(t *testing.T) {
	c := NewYahooOidcClient()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://userinfo.yahooapis.jp/yconnect/v2/attribute",
		httpmock.NewStringResponder(
			200,
			` {
   "sub": "qHbYdL0GOBOZRrjW1mB1Y6YqTjjlfaPo0-CRsbm9l511M8GoQM9G00xre",
   "name": "山田太郎",
   "given_name": "太郎",
   "given_name#ja-Kana-JP": "タロウ",
   "given_name#ja-Hani-JP": "太郎",
   "family_name": "山田",
   "family_name#ja-Kana-JP": "山田",
   "family_name#ja-Hani-JP": "ヤマダ",
   "gender": "male",
   "locale": "ja-JP",
   "email": "aaa@example.com",
   "email_verified": true,
   "address": {
     "country": "jp",
     "postal_code": "1680080"
   },
   "birthdate": "1996",
   "zoneinfo": "Asia/Tokyo",
   "nickname": "wtl********",
   "picture": ""
 }
`,
		),
	)

	actual, err := c.PostUserInfoEndpoint("")

	assert.Nil(t, err)
	assert.Equal(
		t,
		userInfoResponse{Email: "aaa@example.com"},
		actual,
	)
}

func TestRandomState(t *testing.T) {
	state, err := RandomState()

	const expectedLength = 10
	assert.Nil(t, err)
	assert.Equal(t, expectedLength, len(state))
}

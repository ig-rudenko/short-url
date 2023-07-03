package tests

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"short-url/internal/http-server/handlers"
	"testing"
)

const host = "localhost:8000"

func Test_SaveRedirect(t *testing.T) {
	hostUrl := url.URL{
		Scheme: "http",
		Host:   host,
	}
	sameAlias := gofakeit.Word() + gofakeit.Word()
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "First same alias",
			url:   gofakeit.URL(),
			alias: sameAlias,
		},
		{
			name:  "Second same alias",
			url:   gofakeit.URL(),
			alias: sameAlias,
			error: "URL под таким именем уже существует",
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "Неверный URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
		// TODO: add more test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			e := httpexpect.Default(t, hostUrl.String())

			// Save

			resp := e.POST("/url").
				WithJSON(handlers.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")
				resp.Value("error").String().IsEqual(tc.error)
				return

			} else {
				resp.
					ContainsKey("alias").
					ContainsKey("status").
					ContainsValue("OK")
			}

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				tc.alias = resp.Value("alias").String().Raw()
			}

			// Redirect
			testRedirect(t, tc.alias, tc.url)

		})
	}

	// Delete Aliases
	for _, tc := range testCases {
		if tc.error != "" {
			t.Run(tc.name, func(t *testing.T) {
				tc := tc
				t.Parallel()
				e := httpexpect.Default(t, hostUrl.String())
				e.DELETE("/" + tc.alias).Expect().Status(http.StatusOK)
			})
		}
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // stop after 1st redirect
		},
	}

	resp, err := client.Get(u.String())
	require.NoError(t, err)

	require.Equal(t, resp.StatusCode, http.StatusFound)

	defer func() { _ = resp.Body.Close() }()

	require.Equal(t, urlToRedirect, resp.Header.Get("Location"))
}

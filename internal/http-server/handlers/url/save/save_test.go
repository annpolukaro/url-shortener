package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/annpolukaro/url-shortener/internal/http-server/handlers/url/save"
	"github.com/annpolukaro/url-shortener/internal/http-server/handlers/url/save/mocks"
	"github.com/annpolukaro/url-shortener/internal/lib/logger/handlers/slogdiscard"

	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name        string
		alias       string
		url         string
		respError   string
		mockError   error
		returnAlias string
		statusCode  int
	}{
		{
			name:        "Success",
			alias:       "test_alias",
			url:         "https://google.com",
			returnAlias: "test_alias",
			statusCode:  http.StatusCreated,
		},
		{
			name:        "Empty alias",
			alias:       "",
			url:         "https://google.com",
			returnAlias: "generated_alias",
			statusCode:  http.StatusCreated,
		},
		{
			name:       "Empty URL",
			alias:      "test_alias",
			url:        "",
			respError:  "field URL is a required field",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Invalid URL",
			alias:      "test_alias",
			url:        "some invalid URL",
			respError:  "field URL is not a valid URL",
			statusCode: http.StatusBadRequest,
		},
		{
			name:        "SaveURL Error",
			alias:       "test_alias",
			url:         "https://google.com",
			respError:   "failed to add url",
			mockError:   errors.New("unexpected error"),
			returnAlias: "",
			statusCode:  http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.
					On("SaveURL", tc.url, tc.alias).
					Return(tc.returnAlias, tc.mockError).
					Once()
			}

			handler := save.New(
				slogdiscard.NewDiscardLogger(),
				urlSaverMock,
			)

			input := fmt.Sprintf(
				`{"url": "%s", "alias": "%s"}`,
				tc.url,
				tc.alias,
			)

			req, err := http.NewRequest(
				http.MethodPost,
				"/save",
				bytes.NewReader([]byte(input)),
			)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.statusCode, rr.Code)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(
				t,
				json.Unmarshal([]byte(body), &resp),
			)

			require.Equal(t, tc.respError, resp.Error)

			if tc.respError == "" {
				require.Equal(t, tc.returnAlias, resp.Alias)
			}
		})
	}
}

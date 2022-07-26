package handler

import (
	"compress/gzip"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, request testRequestStruct) *http.Response {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(request.method, request.target, request.body)
	if request.acceptEncoding != nil {
		req.Header.Set(headerAcceptEncoding, *request.acceptEncoding)
	}
	if request.cookie != nil {
		req.AddCookie(request.cookie)
	}

	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	if request.acceptEncoding == nil {
		resp.Header.Del(`Content-Encoding`)
	}
	return resp
}

func getResponseReader(t *testing.T, resp *http.Response) io.Reader {
	var reader io.Reader
	if resp.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)
		reader = gz
		defer gz.Close()
	} else {
		reader = resp.Body
	}

	return reader
}

func TestHandler_ServeHTTP(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockService := createMockService(ctl)
	handler := NewHandler(mockService)
	ts := httptest.NewServer(handler)

	tests := getTestsDataList(t, ts)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := testRequest(t, tt.request)
			assert.Equal(t, tt.result.code, resp.StatusCode)

			if tt.result.contentEncoding != nil {
				assert.Equal(t, *tt.result.contentEncoding, resp.Header.Get(headerContentEncoding))
			}

			defer resp.Body.Close()

			reader := getResponseReader(t, resp)

			respBody, err := ioutil.ReadAll(reader)
			require.NoError(t, err)

			if tt.result.body != "" {
				assert.Equal(t, tt.result.body, string(respBody))
			}

			if tt.result.contentType != "" {
				assert.Equal(t, tt.result.contentType, resp.Header.Get(headerContentType))
			}
		})
	}
}

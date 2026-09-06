package ginutils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type testBindRequest struct {
	Int      int64  `binding:"required,min=1"     uri:"int"`
	String   string `binding:"required,max=128"   form:"string"   json:"string"`
	Optional string `binding:"omitempty,max=1024" form:"optional" json:"optional"`
}

func TestBindAll(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		method      string
		route       string
		uri         string
		body        string
		contentType string
		opts        BindOptions
		valid       bool
		expected    testBindRequest
	}{
		{
			name:        "uri and form valid",
			method:      http.MethodPost,
			route:       "/prompts/:int",
			uri:         "/prompts/42",
			body:        "string=test&optional=value",
			contentType: "application/x-www-form-urlencoded",
			opts: BindOptions{
				Uri:  true,
				Form: true,
			},
			valid: true,
			expected: testBindRequest{
				Int:      42,
				String:   "test",
				Optional: "value",
			},
		},
		{
			name:        "missing required form field",
			method:      http.MethodPost,
			route:       "/prompts/:int",
			uri:         "/prompts/42",
			body:        "optional=value",
			contentType: "application/x-www-form-urlencoded",
			opts: BindOptions{
				Uri:  true,
				Form: true,
			},
			expected: testBindRequest{
				Int:      42,
				Optional: "value",
			},
		},
		{
			name:        "invalid uri value",
			method:      http.MethodPost,
			route:       "/prompts/:int",
			uri:         "/prompts/abc",
			body:        "string=test",
			contentType: "application/x-www-form-urlencoded",
			opts: BindOptions{
				Uri:  true,
				Form: true,
			},
		},
		{
			name:        "missing uri param",
			method:      http.MethodPost,
			route:       "/prompts",
			uri:         "/prompts",
			body:        "string=test",
			contentType: "application/x-www-form-urlencoded",
			opts: BindOptions{
				Uri:  true,
				Form: true,
			},
			expected: testBindRequest{
				String: "test",
			},
		},
		{
			name:        "malformed form body",
			method:      http.MethodPost,
			route:       "/prompts/:int",
			uri:         "/prompts/42",
			body:        "%zz",
			contentType: "application/x-www-form-urlencoded",
			opts: BindOptions{
				Form: true,
			},
		},
		{
			name:   "query string",
			method: http.MethodGet,
			route:  "/prompts/:int",
			uri:    "/prompts/42?string=test",
			opts: BindOptions{
				Uri:  true,
				Form: true,
			},
			valid: true,
			expected: testBindRequest{
				Int:    42,
				String: "test",
			},
		},
		{
			name:   "multipart valid",
			method: http.MethodPost,
			route:  "/prompts/:int",
			uri:    "/prompts/42",
			body: "--BOUNDARY\r\n" +
				"Content-Disposition: form-data; name=\"string\"\r\n" +
				"\r\n" +
				"test\r\n" +
				"--BOUNDARY\r\n" +
				"Content-Disposition: form-data; name=\"optional\"\r\n" +
				"\r\n" +
				"value\r\n" +
				"--BOUNDARY--\r\n",
			contentType: "multipart/form-data; boundary=BOUNDARY",
			opts: BindOptions{
				Uri:  true,
				Form: true,
			},
			valid: true,
			expected: testBindRequest{
				Int:      42,
				String:   "test",
				Optional: "value",
			},
		},
		{
			name:        "no sources enabled",
			method:      http.MethodPost,
			route:       "/prompts/:int",
			uri:         "/prompts/42",
			body:        `{"string":"test"}`,
			contentType: "application/json",
			opts:        BindOptions{},
		},
		{
			name:        "json valid",
			method:      http.MethodPost,
			route:       "/prompts/:int",
			uri:         "/prompts/42",
			body:        `{"string":"test","optional":"value"}`,
			contentType: "application/json",
			opts: BindOptions{
				Uri:  true,
				JSON: true,
			},
			valid: true,
			expected: testBindRequest{
				Int:      42,
				String:   "test",
				Optional: "value",
			},
		},
		{
			name:        "json malformed body",
			method:      http.MethodPost,
			route:       "/prompts/:int",
			uri:         "/prompts/42",
			body:        "string=test",
			contentType: "application/json",
			opts: BindOptions{
				JSON: true,
			},
		},
		{
			name:        "json body with form query",
			method:      http.MethodPost,
			route:       "/prompts/:int",
			uri:         "/prompts/42?optional=query",
			body:        `{"string":"body"}`,
			contentType: "application/json",
			opts: BindOptions{
				Uri:  true,
				Form: true,
				JSON: true,
			},
			valid: true,
			expected: testBindRequest{
				Int:      42,
				String:   "body",
				Optional: "query",
			},
		},
		{
			name:        "form body with json",
			method:      http.MethodPost,
			route:       "/prompts/:int",
			uri:         "/prompts/42",
			body:        "string=test",
			contentType: "application/x-www-form-urlencoded",
			opts: BindOptions{
				Uri:  true,
				Form: true,
				JSON: true,
			},
			expected: testBindRequest{
				Int:    42,
				String: "test",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var (
				got     testBindRequest
				valid   bool
				bindErr int
			)

			router := gin.New()
			router.Handle(test.method, test.route, func(c *gin.Context) {
				valid = BindAll(c, &got, test.opts)
				bindErr = len(c.Errors.ByType(gin.ErrorTypeBind))
			})

			req := httptest.NewRequestWithContext(
				t.Context(),
				test.method,
				test.uri,
				strings.NewReader(test.body),
			)
			if test.contentType != "" {
				req.Header.Set("Content-Type", test.contentType)
			}

			router.ServeHTTP(httptest.NewRecorder(), req)

			require.Equal(t, test.valid, valid)
			require.Equal(t, test.expected, got)

			if test.valid {
				require.Zero(t, bindErr)
			} else {
				require.Equal(t, 1, bindErr)
			}
		})
	}
}

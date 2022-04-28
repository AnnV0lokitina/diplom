package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testRequestStruct struct {
	method         string
	target         string
	body           io.Reader
	acceptEncoding *string
	cookie         *http.Cookie
}

type config struct {
	RunAddress string `env:"RUN_ADDRESS"  envDefault:"localhost:8080"`
}

type testResultStruct struct {
	body            string
	headerLocation  string
	code            int
	contentType     string
	contentEncoding *string
}

type testStruct struct {
	name        string
	setURLError bool
	request     testRequestStruct
	result      testResultStruct
}

func newStringPtr(x string) *string {
	return &x
}

func newBoolPtr(x bool) *bool {
	return &x
}

//func createJSONEncodedResponse(t *testing.T, responseURL string) string {
//	jsonResponse := JSONResponse{
//		Result: responseURL,
//	}
//
//	jsonEncodedResponse, err := json.Marshal(jsonResponse)
//	require.NoError(t, err)
//
//	jsonEncodedResponse = append(jsonEncodedResponse, '\n')
//
//	return string(jsonEncodedResponse)
//}

func getTestsDataList(_ *testing.T, ts *httptest.Server) []testStruct {
	return []testStruct{
		{
			name: "test register user incorrect 400",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/register",
				body:   strings.NewReader(""),
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name: "test register user incorrect 400",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/register",
				body:   strings.NewReader("{\"login\":\"login\",\"password\":\"password\"},"),
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name: "test register user correct 200",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/register",
				body:   strings.NewReader("{\"login\":\"login\",\"password\":\"password\"}"),
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusOK,
				contentType: "",
			},
		},
		{
			name: "test register user correct 200",
			request: testRequestStruct{
				method:         http.MethodPost,
				target:         ts.URL + "/api/user/register",
				body:           strings.NewReader("{\"login\":\"login\",\"password\":\"password\"}"),
				acceptEncoding: newStringPtr(encoding),
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusOK,
				contentType: "",
			},
		},
		{
			name: "test register user incorrect 500",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/register",
				body:   strings.NewReader("{\"login\":\"login\",\"password\":\"password\"}"),
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusInternalServerError,
				contentType: "",
			},
		},
		{
			name: "test register user incorrect 409",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/register",
				body:   strings.NewReader("{\"login\":\"login\",\"password\":\"password\"}"),
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusConflict,
				contentType: "",
			},
		},
		{
			name: "test login user incorrect 400",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/login",
				body:   strings.NewReader("{\"login\":\"login\",\"password\":\"password\"},"),
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name: "test login user incorrect 401",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/login",
				body:   strings.NewReader("{\"login\":\"login\",\"password\":\"password\"}"),
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusUnauthorized,
				contentType: "",
			},
		},
		{
			name: "test login user incorrect 500",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/login",
				body:   strings.NewReader("{\"login\":\"login\",\"password\":\"password\"}"),
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusInternalServerError,
				contentType: "",
			},
		},
		{
			name: "test login user correct 200",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/login",
				body:   strings.NewReader("{\"login\":\"login\",\"password\":\"password\"}"),
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusOK,
				contentType: "",
			},
		},
		{
			name: "test post order incorrect 401",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/orders",
				body:   strings.NewReader("74528626868518"),
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusUnauthorized,
				contentType: "",
			},
		},
		{
			name: "test post order incorrect 401",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/orders",
				body:   strings.NewReader("74528626868518"),
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusUnauthorized,
				contentType: "",
			},
		},
		{
			name: "test post order incorrect 400",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/orders",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name: "test post order incorrect 422",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/orders",
				body:   strings.NewReader("74528626868518"),
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusUnprocessableEntity,
				contentType: "",
			},
		},
		{
			name: "test post order correct 200",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/orders",
				body:   strings.NewReader("74528626868518"),
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusOK,
				contentType: "",
			},
		},
		{
			name: "test post order incorrect 409",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/orders",
				body:   strings.NewReader("74528626868518"),
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusConflict,
				contentType: "",
			},
		},
		{
			name: "test post order incorrect 500",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/orders",
				body:   strings.NewReader("74528626868518"),
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusInternalServerError,
				contentType: "",
			},
		},
		{
			name: "test post order correct 202",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/orders",
				body:   strings.NewReader("74528626868518"),
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusAccepted,
				contentType: "",
			},
		},
		{
			name: "test get orders correct 200",
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/api/user/orders",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body: "[{\"number\":\"9278923470\",\"status\":\"PROCESSED\",\"accrual\":500," +
					"\"uploaded_at\":\"2020-12-10T15:15:45+03:00\"},{\"number\":\"12345678903\"," +
					"\"status\":\"PROCESSING\",\"uploaded_at\":\"2020-12-10T15:12:01+03:00\"}," +
					"{\"number\":\"346436439\",\"status\":\"INVALID\",\"uploaded_at\":\"2020-12-09T16:09:53+03:00\"}]\n",
				code:        http.StatusOK,
				contentType: jsonContentType,
			},
		},
		{
			name: "test get orders incorrect 204",
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/api/user/orders",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusNoContent,
				contentType: "",
			},
		},
		{
			name: "test get orders incorrect 401",
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/api/user/orders",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusUnauthorized,
				contentType: "",
			},
		},
		{
			name: "test get orders incorrect 500",
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/api/user/orders",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusInternalServerError,
				contentType: "",
			},
		},
		{
			name: "test get balance correct 200",
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/api/user/balance",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "{\"current\":500.5,\"withdrawn\":42}\n",
				code:        http.StatusOK,
				contentType: jsonContentType,
			},
		},
		{
			name: "test get balance incorrect 401",
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/api/user/balance",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusUnauthorized,
				contentType: "",
			},
		},
		{
			name: "test get balance incorrect 500",
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/api/user/balance",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusInternalServerError,
				contentType: "",
			},
		},
		{
			name: "test withdraw correct 200",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/balance/withdraw",
				body:   strings.NewReader("{\"order\":\"2377225624\",\"sum\":751}"),
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusOK,
				contentType: "",
			},
		},
		{
			name: "test withdraw incorrect 401",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/balance/withdraw",
				body:   strings.NewReader("{\"order\":\"2377225624\",\"sum\":751}"),
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusUnauthorized,
				contentType: "",
			},
		},
		{
			name: "test withdraw incorrect 402",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/balance/withdraw",
				body:   strings.NewReader("{\"order\":\"2377225624\",\"sum\":751}"),
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusPaymentRequired,
				contentType: "",
			},
		},
		{
			name: "test withdraw incorrect 422",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/balance/withdraw",
				body:   strings.NewReader("{\"order\":\"2377225624\",\"sum\":751}"),
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusUnprocessableEntity,
				contentType: "",
			},
		},
		{
			name: "test withdraw incorrect 500",
			request: testRequestStruct{
				method: http.MethodPost,
				target: ts.URL + "/api/user/balance/withdraw",
				body:   strings.NewReader("{\"order\":\"2377225624\",\"sum\":751}"),
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusInternalServerError,
				contentType: "",
			},
		},
		{
			name: "test get withdrawals correct 200",
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/api/user/balance/withdrawals",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "[{\"order\":\"2377225624\",\"sum\":500,\"processed_at\":\"2020-12-09T16:09:57+03:00\"}]\n",
				code:        http.StatusOK,
				contentType: jsonContentType,
			},
		},
		{
			name: "test get withdrawals correct 204",
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/api/user/balance/withdrawals",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusNoContent,
				contentType: "",
			},
		},
		{
			name: "test get withdrawals incorrect 401",
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/api/user/balance/withdrawals",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusUnauthorized,
				contentType: "",
			},
		},
		{
			name: "test get withdrawals incorrect 500",
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/api/user/balance/withdrawals",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "session",
					Value:    "session_id",
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:        "",
				code:        http.StatusInternalServerError,
				contentType: "",
			},
		},
	}
}

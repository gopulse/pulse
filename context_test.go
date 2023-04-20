package pulse

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestContext_Write(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := NewContext(w, nil)

	message := "Hello, world!"

	n, err := ctx.Write([]byte(message))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if n != len(message) {
		t.Errorf("Expected %d bytes written, got %d", len(message), n)
	}

	if w.Body.String() != message {
		t.Errorf("Response body does not match expected value. Expected: %s, got: %s", message, w.Body.String())
	}
}

func TestContext_WithParams(t *testing.T) {
	ctx := NewContext(nil, nil)
	ctx.WithParams(map[string]string{"id": "1"})
	if ctx.Params["id"] != "1" {
		t.Errorf("Expected id to be 1, got %s", ctx.Params["id"])
	}
}

func TestContext_Param(t *testing.T) {
	ctx := NewContext(nil, nil)
	ctx.WithParams(map[string]string{"id": "1"})
	if ctx.Param("id") != "1" {
		t.Errorf("Expected id to be 1, got %s", ctx.Param("id"))
	}
}

func TestContext_Query(t *testing.T) {
	ctx := NewContext(nil, nil)
	ctx.Request = &http.Request{
		URL: &url.URL{
			RawQuery: "id=1",
		},
	}
	if ctx.Query("id") != "1" {
		t.Errorf("Expected id to be 1, got %s", ctx.Query("id"))
	}
}

func TestContext_Abort(t *testing.T) {
	ctx := NewContext(nil, nil)
	ctx.Abort()
}

func TestContext_String(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := NewContext(w, nil)

	message := "Hello, world!"
	ctx.String(message)

	if w.Body.String() != message {
		t.Errorf("Response body does not match expected value. Expected: %s, got: %s", message, w.Body.String())
	}
}

func TestContext_Cookie(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	ctx := NewContext(w, r)

	cookieName := "test-cookie"
	cookieValue := "test-value"

	// Set a cookie using the SetCookie method.
	ctx.SetCookie(&Cookie{
		Name:  cookieName,
		Value: cookieValue,
	})

	// Verify that the cookie was set correctly by checking the response header.
	cookies := w.Header().Get("Set-Cookie")
	if !strings.Contains(cookies, cookieName) {
		t.Errorf("Expected response header to contain cookie name '%s'", cookieName)
	}
	if !strings.Contains(cookies, cookieValue) {
		t.Errorf("Expected response header to contain cookie value '%s'", cookieValue)
	}

	// Get the value of the cookie using the GetCookie method.
	retrievedValue := ctx.GetCookie(cookieName)
	fmt.Println(retrievedValue)
	if retrievedValue != cookieValue {
		t.Errorf("Expected retrieved cookie value to be '%s', but got '%s'", cookieValue, retrievedValue)
	}

	// Clear the cookie using the ClearCookie method.
	ctx.ClearCookie(cookieName)
}

func TestContext_GetCookie(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		cookie := ctx.GetCookie("test")
		ctx.String(cookie)
		return nil
	})

}

func TestContext_ClearCookie(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		ctx.ClearCookie("test")
		return nil
	})
}

func TestContext_Header(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	ctx := NewContext(w, r)

	headerKey := "test-header"
	headerValue := "test-value"
	ctx.SetResponseHeader(headerKey, headerValue)

	retrievedHeaderValue := ctx.GetResponseHeader(headerKey)
	if retrievedHeaderValue != headerValue {
		t.Errorf("Expected response header value to be '%s', but got '%s'", headerValue, retrievedHeaderValue)
	}

	reqHeaderKey := "test-request-header"
	reqHeaderValue := "test-request-value"
	ctx.SetRequestHeader(reqHeaderKey, reqHeaderValue)

	retrievedReqHeaderValue := ctx.GetRequestHeader(reqHeaderKey)
	if retrievedReqHeaderValue != reqHeaderValue {
		t.Errorf("Expected request header value to be '%s', but got '%s'", reqHeaderValue, retrievedReqHeaderValue)
	}

	anotherHeaderKey := "another-test-header"
	anotherHeaderValue := "another-test-value"
	ctx.SetResponseHeader(anotherHeaderKey, anotherHeaderValue)

	retrievedAnotherHeaderValue := w.Header().Get(anotherHeaderKey)
	if retrievedAnotherHeaderValue != anotherHeaderValue {
		t.Errorf("Expected response header value to be '%s', but got '%s'", anotherHeaderValue, retrievedAnotherHeaderValue)
	}
}

func TestContext_SetData(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := NewContext(w, nil)

	key := "custom-key"
	value := "custom-value"

	ctx.SetData(key, value)

	if w.Header().Get(key) != value {
		t.Errorf("Response header does not match expected value. Expected: %s, got: %s", value, w.Header().Get(key))
	}
}

func TestContext_GetData(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("custom-key", "custom-value")

	w := httptest.NewRecorder()
	ctx := NewContext(w, r)

	key := "custom-key"

	if ctx.GetData(key) != "custom-value" {
		t.Errorf("Request header does not match expected value. Expected: custom-value, got: %s", ctx.GetData(key))
	}
}

func TestContext_Next(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := NewContext(w, nil)

	ctx.Next()
}

func TestContext_Reset(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	ctx := NewContext(w, r)

	ctx.handlerIdx = 1
	ctx.Reset()

	if ctx.handlerIdx != -1 {
		t.Errorf("Expected handler index to be -1 after calling Reset(), but got: %d", ctx.handlerIdx)
	}
}

func TestContext_Status(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := NewContext(w, nil)

	ctx.Status(200)
}

func TestContext_JSON(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := NewContext(w, nil)

	ctx.JSON(200, map[string]string{"test": "test"})
}

func TestContext_SetContentType(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := NewContext(w, r)

	ctx.SetContentType("application/json")
}

func TestContext_Accepts(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	ctx := NewContext(w, r)

	// Set the Accept header to a single value of type "text/html".
	ctx.SetRequestHeader("Accept", "text/html")

	// Verify that the Accept header was set correctly.
	retrievedHeaderValue := ctx.GetRequestHeader("Accept")
	if retrievedHeaderValue != "text/html" {
		t.Errorf("Expected request header value to be 'text/html', but got '%s'", retrievedHeaderValue)
	}

	// Test with a single acceptable media type.
	acceptedType := ctx.Accepts("text/html")

	if acceptedType != "text/html" {
		t.Errorf("Expected acceptable media type to be 'text/html', but got '%s'", acceptedType)
	}

	// Test with multiple acceptable media types, including one that is not present in the Accept header.
	acceptedType = ctx.Accepts("text/plain", "text/html")

	if acceptedType != "text/html" {
		t.Errorf("Expected acceptable media type to be 'text/html', but got '%s'", acceptedType)
	}

	// Test with multiple acceptable media types, none of which are present in the Accept header.
	acceptedType = ctx.Accepts("application/json", "text/plain")

	if acceptedType != "" {
		t.Errorf("Expected no acceptable media type, but got '%s'", acceptedType)
	}
}

func TestContext_BodyParser(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "John", "age": 30}`))

	ctx := NewContext(w, r)

	// Define a struct to decode the request body into.
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Parse the request body using the BodyParser method.
	var person Person
	err := ctx.BodyParser(&person)

	if err != nil {
		t.Errorf("Expected no error, but got '%s'", err.Error())
	} else {
		// Verify that the request body was parsed correctly.
		if person.Name != "John" || person.Age != 30 {
			t.Errorf("Expected person object to have name 'John' and age 30, but got %+v", person)
		}
	}
}

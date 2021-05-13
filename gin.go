package test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
)

const (
	COOKIE_SESSION = "session_id"
)

var (
	defaultContext = newTestContext()
	// client is safe for concurrent use
	client = http.Client{}
)

type testContext struct {
	mu sync.RWMutex
	data map[string]interface{}
    sessionId string
}

func newTestContext() *testContext {
	return &testContext{
		data: make(map[string]interface{}),
	}
}

func (c *testContext) do(t *testing.T, url string, withSession, setSession bool, method string, statusCode int, query, req, resp interface{}) {
	// 这里考虑对req进行反射，拼接url
	if query != nil {
		url = JointUrl(url, query)
	}
	var body io.Reader
	if req != nil {
		bodyBytes, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
        body = bytes.NewBuffer(bodyBytes)
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}

	if withSession {
		c.mu.RLock()
		sessionId := c.sessionId
		c.mu.RUnlock()
		request.AddCookie(&http.Cookie{
			Name:       COOKIE_SESSION,
			Value:      sessionId,
		})
	}

	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != statusCode {
		t.Fatalf("invoke %s fail, expect %d, get %d", url, statusCode, response.StatusCode)
	}
	if setSession {
		for _, cookie := range response.Cookies() {
			if cookie.Name == COOKIE_SESSION {
				c.mu.Lock()
				c.sessionId = cookie.Value
				c.mu.Unlock()
				break
			}
		}
	}
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(respBytes) > 0 && resp != nil {
		err = json.Unmarshal(respBytes, resp)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (c *testContext) get(t *testing.T, url string, withSession, setSession bool, statusCode int, query, req, resp interface{}) {
	c.do(t, url, withSession, setSession, "GET", statusCode, query, req, resp)
}

func (c *testContext) post(t *testing.T, url string, withSession, setSession bool, statusCode int, query, req, resp interface{}) {
	c.do(t, url, withSession, setSession, "POST", statusCode, query, req, resp)
}

func (c *testContext) put(t *testing.T, url string, withSession bool, statusCode int, query, req, resp interface{}) {
	c.do(t, url, withSession,false, "PUT", statusCode, query, req, resp)
}

func (c *testContext) delete(t *testing.T, url string, withSession bool, statusCode int, query, req, resp interface{}) {
	c.do(t, url, withSession,false, "DELETE", statusCode, query, req, resp)
}

func (c *testContext) setValue(key string, value interface{}) {
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
}

func (c *testContext) getValue(key string) (value interface{}, flag bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, flag = c.data[key]
	return
}

func Get(t *testing.T, url string, withSession, setSession bool, statusCode int, query, req, resp interface{}) {
	defaultContext.get(t, url, withSession, setSession, statusCode, query, req, resp)
}

func Post(t *testing.T, url string, withSession, setSession bool, statusCode int, query, req, resp interface{}) {
	defaultContext.post(t, url, withSession, setSession, statusCode, query, req, resp)
}

func Put(t *testing.T, url string, withSession, setSession bool, statusCode int, query, req, resp interface{}) {
	defaultContext.put(t, url, withSession, statusCode, query, req, resp)
}

func Delete(t *testing.T, url string, withSession, setSession bool, statusCode int, query, req, resp interface{}) {
	defaultContext.delete(t, url, withSession, statusCode, query, req, resp)
}

func SetValue(key string, value interface{}) {
	defaultContext.setValue(key, value)
}

func GetValue(key string) (value interface{}, flag bool) {
	return defaultContext.getValue(key)
}
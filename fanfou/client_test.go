package fanfou

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/mrjones/oauth"
	httpmock "gopkg.in/jarcoal/httpmock.v1"

	"github.com/stretchr/testify/assert"
)

func TestClients(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func()
	}{
		{
			desc: "test NewClientWithOAuth",
			testFunc: func() {
				var err error
				client := NewClientWithOAuth(mockConsumerKey, mockConsumerSecret)
				accessToken := oauth.AccessToken{
					Token:  mockAccessToken,
					Secret: mockAccessSecret,
				}
				client.http, err = client.OAuthConsumer.MakeHttpClient(&accessToken)

				assert.Nil(t, err)
				assert.NotNil(t, client.http)
			},
		},
		{
			desc: "test successful responses (200)",
			testFunc: func() {
				var err error
				// clear mocks
				httpmock.Reset()

				client := NewClientWithOAuth(mockConsumerKey, mockConsumerSecret)
				accessToken := oauth.AccessToken{
					Token:  mockAccessToken,
					Secret: mockAccessSecret,
				}
				client.http, err = client.OAuthConsumer.MakeHttpClient(&accessToken)

				assert.Nil(t, err)

				// register new mocks
				for key, mep := range mockEndpoints {
					httpmock.RegisterResponder(mep.Method, mep.URL, func(req *http.Request) (*http.Response, error) {
						return httpmock.NewStringResponse(mep.Result200.Code, mep.Result200.Body), nil
					})

					v := reflect.ValueOf(client).MethodByName(key).Call([]reflect.Value{reflect.ValueOf(&ReqParams{Photo: "./def.go", Image: "./def.go"})})

					assert.NotNil(t, v[0])
					assert.Equal(t, reflect.TypeOf([]uint8{}), v[1].Type())
					assert.True(t, reflect.Value(v[2]).IsNil())

					if !reflect.Value(v[2]).IsNil() {
						fmt.Printf("[%s] error: %+v\n", key, v[2])
					}
				}
			},
		},
		{
			desc: "test client-error responses (400)",
			testFunc: func() {
				var err error
				// clear mocks
				httpmock.Reset()

				client := NewClientWithOAuth(mockConsumerKey, mockConsumerSecret)
				accessToken := oauth.AccessToken{
					Token:  mockAccessToken,
					Secret: mockAccessSecret,
				}
				client.http, err = client.OAuthConsumer.MakeHttpClient(&accessToken)

				assert.Nil(t, err)

				// register new mocks
				for key, mep := range mockEndpoints {
					httpmock.RegisterResponder(mep.Method, mep.URL, func(req *http.Request) (*http.Response, error) {
						return httpmock.NewStringResponse(mep.Result400.Code, mep.Result400.Body), nil
					})

					v := reflect.ValueOf(client).MethodByName(key).Call([]reflect.Value{reflect.ValueOf(&ReqParams{Photo: "./def.go", Image: "./def.go"})})

					// check for errors first
					assert.NotNil(t, v[2])
					if reflect.Value(v[2]).IsNil() {
						fmt.Printf("[%s] error: %+v\n", key, v[2])
					}

					switch key {
					case "FriendshipsExists":
						assert.False(t, reflect.Value(v[0]).Bool())
					default:
						assert.True(t, reflect.Value(v[0]).IsNil())
					}

					assert.True(t, reflect.Value(v[1]).IsNil())
				}
			},
		},
		{
			desc: "test server-error responses (500)",
			testFunc: func() {
				var err error
				// clear mocks
				httpmock.Reset()

				client := NewClientWithOAuth(mockConsumerKey, mockConsumerSecret)
				accessToken := oauth.AccessToken{
					Token:  mockAccessToken,
					Secret: mockAccessSecret,
				}
				client.http, err = client.OAuthConsumer.MakeHttpClient(&accessToken)

				assert.Nil(t, err)

				// register new mocks
				for key, mep := range mockEndpoints {
					httpmock.RegisterResponder(mep.Method, mep.URL, func(req *http.Request) (*http.Response, error) {
						return httpmock.NewStringResponse(mep.Result500.Code, mep.Result500.Body), nil
					})

					v := reflect.ValueOf(client).MethodByName(key).Call([]reflect.Value{reflect.ValueOf(&ReqParams{Photo: "./def.go", Image: "./def.go"})})

					// check for errors first
					assert.NotNil(t, v[2])
					if reflect.Value(v[2]).IsNil() {
						fmt.Printf("[%s] error: %+v\n", key, v[2])
					}

					switch key {
					case "FriendshipsExists":
						assert.False(t, reflect.Value(v[0]).Bool())
					default:
						assert.True(t, reflect.Value(v[0]).IsNil())
					}

					assert.True(t, reflect.Value(v[1]).IsNil())
				}
			},
		},
		{
			desc: "test chaos-error responses (200 with malformed body)",
			testFunc: func() {
				var err error
				// clear mocks
				httpmock.Reset()

				client := NewClientWithOAuth(mockConsumerKey, mockConsumerSecret)
				accessToken := oauth.AccessToken{
					Token:  mockAccessToken,
					Secret: mockAccessSecret,
				}
				client.http, err = client.OAuthConsumer.MakeHttpClient(&accessToken)

				assert.Nil(t, err)

				// register new mocks
				for key, mep := range mockEndpoints {
					httpmock.RegisterResponder(mep.Method, mep.URL, func(req *http.Request) (*http.Response, error) {
						return httpmock.NewStringResponse(mep.ResultChaos.Code, mep.ResultChaos.Body), nil
					})

					v := reflect.ValueOf(client).MethodByName(key).Call([]reflect.Value{reflect.ValueOf(&ReqParams{Photo: "./def.go", Image: "./def.go"})})

					// check for errors first
					assert.NotNil(t, v[2])
					if reflect.Value(v[2]).IsNil() {
						fmt.Printf("[%s] error: %+v\n", key, v[2])
					}

					assert.False(t, reflect.Value(v[1]).IsNil())
				}
			},
		},
	}

	for _, testCase := range testCases {
		testCase.testFunc()
	}
}

package graph_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"firebase.google.com/go/auth"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/graph"
	"gitlab.slade360emr.com/go/profile/graph/profile"
)

const (
	testHTTPClientTimeout = 180
)

var allowedOrigins = []string{
	"https://healthcloud.co.ke",
	"https://bewell.healthcloud.co.ke",
	"http://localhost:5000",
	"https://api-gateway-test.healthcloud.co.ke",
	"https://api-gateway-prod.healthcloud.co.ke",
	"https://profile-testing-uyajqt434q-ew.a.run.app",
	"https://profile-prod-uyajqt434q-ew.a.run.app",
}

// these are set up once in TestMain and used by all the acceptance tests in
// this package
var srv *http.Server
var baseURL string
var serverErr error

func TestMain(m *testing.M) {
	// setup
	ctx := context.Background()
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("ROOT_COLLECTION_SUFFIX", "onboarding_testing")
	s := profile.NewService()
	srv, baseURL, serverErr = base.StartTestServer(ctx, graph.PrepareServer, allowedOrigins) // set the globals
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
	}

	// run the tests
	log.Printf("about to run tests")
	code := m.Run()
	log.Printf("finished running tests")

	fc := &base.FirebaseClient{}
	fa, err := fc.InitFirebase()
	if err != nil {
		log.Printf("can't initialize Firebase app: %s", err)
	}
	firestore, err := fa.Firestore(context.Background())
	if err != nil {
		log.Printf("can't initialize Firestore client: %s", err)
	}
	collections := []string{
		s.GetPINCollectionName(),
		s.GetUserProfileCollectionName(),
		s.GetPractitionerCollectionName(),
	}
	for _, collection := range collections {
		ref := firestore.Collection(collection)
		base.DeleteCollection(ctx, firestore, ref, 10)
	}
	// cleanup here
	defer func() {
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("test server shutdown error: %s", err)
		}
	}()
	os.Exit(code)
}

func mapToJSONReader(m map[string]interface{}) (io.Reader, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
	}

	buf := bytes.NewBuffer(bs)
	return buf, nil
}

func TestGetProfileAttributesHandler(t *testing.T) {
	client := http.DefaultClient
	_, emailUserAuthToken := base.GetAuthenticatedContextAndToken(t)
	if emailUserAuthToken == nil {
		t.Errorf("can't get test auth token")
		return
	}

	uids := profile.UserUIDs{
		UIDs: []string{
			emailUserAuthToken.UID,
		},
	}
	bs, err := json.Marshal(uids)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "successful get confirmed email addresses",
			args: args{
				url: fmt.Sprintf(
					"%s/internal/contactdetails/%s/",
					baseURL,
					"emails",
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(bs),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "failed get confirmed email addresses",
			args: args{
				url: fmt.Sprintf(
					"%s/internal/contactdetails/%s/",
					baseURL,
					"emails",
				),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
		{
			name: "successful get confirmed phone numbers",
			args: args{
				url: fmt.Sprintf(
					"%s/internal/contactdetails/%s/",
					baseURL,
					"phonenumbers",
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(bs),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "failed get confirmed phone numbers",
			args: args{
				url: fmt.Sprintf(
					"%s/internal/contactdetails/%s/",
					baseURL,
					"phonenumbers",
				),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
		{
			name: "successful get FCM tokens",
			args: args{
				url: fmt.Sprintf(
					"%s/internal/contactdetails/%s/",
					baseURL,
					"tokens",
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(bs),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "failed get FCN tokens",
			args: args{
				url: fmt.Sprintf(
					"%s/internal/contactdetails/%s/",
					baseURL,
					"tokens",
				),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range base.GetDefaultHeaders(t, baseURL, "profile") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}
		})
	}
}

func TestGraphQLRequestPinReset(t *testing.T) {
	ctx := base.GetPhoneNumberAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("nil context")
		return
	}
	gql := map[string]interface{}{}
	gql["query"] = `
	query requestPinReset{
		requestPinReset(msisdn: "+254711223344")
	}
	`

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGraphQLUpdatePin(t *testing.T) {
	ctx := base.GetPhoneNumberAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("nil context")
		return
	}
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation updateUserPIN{
		updateUserPIN(msisdn: "+254711223344", pin: "1234")
	}
	`

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}

}

func TestUpdatePinHandler(t *testing.T) {
	client := http.DefaultClient
	srv := profile.NewService()
	assert.NotNil(t, srv, "service is nil")

	ctx, _ := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}
	set, err := srv.SetUserPIN(ctx, base.TestUserPhoneNumber, "1234")
	if !set {
		t.Errorf("setting a pin for test user failed. It returned false")
	}
	if err != nil {
		t.Errorf("setting a pin for test user failed: %v", err)
	}
	if !set {
		t.Errorf("setting a pin for test user failed. It returned false")
	}
	pinRecovery := profile.PinRecovery{
		MSISDN:    base.TestUserPhoneNumber,
		PINNumber: "4565",
		OTP:       strconv.Itoa(rand.Int()),
	}
	bs, err := json.Marshal(pinRecovery)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "successful update pin",
			args: args{
				url:        fmt.Sprintf("%s/update_pin", baseURL),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "failed generate and send otp",
			args: args{
				url:        fmt.Sprintf("%s/update_pin", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range base.GetDefaultHeaders(t, baseURL, "profile") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}
		})
	}
}

func TestRequestPinResetHandler(t *testing.T) {
	client := http.DefaultClient
	srv := profile.NewService()
	assert.NotNil(t, srv, "service is nil")

	ctx, _ := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}
	set, err := srv.SetUserPIN(ctx, base.TestUserPhoneNumber, "1234")
	if !set {
		t.Errorf("setting a pin for test user failed. It returned false")
	}
	if err != nil {
		t.Errorf("setting a pin for test user failed: %v", err)
	}
	pinRecovery := profile.PinRecovery{
		MSISDN: base.TestUserPhoneNumber,
	}
	bs, err := json.Marshal(pinRecovery)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid pin reset request",
			args: args{
				url:        fmt.Sprintf("%s/request_pin_reset", baseURL),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "failed generate and send otp",
			args: args{
				url:        fmt.Sprintf("%s/request_pin_reset", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range base.GetDefaultHeaders(t, baseURL, "profile") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}
		})
	}
}

func TestResetPinHandler(t *testing.T) {
	client := http.DefaultClient
	srv := profile.NewService()
	assert.NotNil(t, srv, "service is nil")

	ctx, _ := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}
	// prepare payload for user with PIN
	pinRecovery := profile.PinRecovery{
		MSISDN:    base.TestUserPhoneNumberWithPin,
		PINNumber: "4565",
		OTP:       strconv.Itoa(rand.Int()),
	}
	bs, err := json.Marshal(pinRecovery)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payloadUserWithPIN := bytes.NewBuffer(bs)

	// prepare payload for user without PIN
	pinRecoveryNoPIN := profile.PinRecovery{
		MSISDN:    base.TestUserPhoneNumber,
		PINNumber: "7895",
		OTP:       strconv.Itoa(rand.Int()),
	}
	bs, err = json.Marshal(pinRecoveryNoPIN)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payloadUserNoPIN := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "successful reset pin for user with an existing PIN",
			args: args{
				url:        fmt.Sprintf("%s/reset_pin", baseURL),
				httpMethod: http.MethodPost,
				body:       payloadUserWithPIN,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "successful reset pin for user with non existent PIN",
			args: args{
				url:        fmt.Sprintf("%s/reset_pin", baseURL),
				httpMethod: http.MethodPost,
				body:       payloadUserNoPIN,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "failed generate and send otp",
			args: args{
				url:        fmt.Sprintf("%s/reset_pin", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range base.GetDefaultHeaders(t, baseURL, "profile") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}
		})
	}
}

func TestRetrieveUserProfileFirebaseDocSnapshotHandler(t *testing.T) {

	ctx := base.GetAuthenticatedContext(t)
	assert.NotNil(t, ctx)
	auth := ctx.Value(base.AuthTokenContextKey).(*auth.Token)
	assert.NotNil(t, auth)
	profileUid := &profile.BusinessPartnerUID{
		UID: auth.UID,
	}
	assert.NotNil(t, profileUid)
	srv := profile.NewService()
	assert.NotNil(t, srv)
	handler := graph.RetrieveUserProfileHandler(ctx, srv)

	assert.NotNil(t, handler)

	uidJson, err := json.Marshal(profileUid)
	assert.NotNil(t, uidJson)
	assert.Nil(t, err)

	validRequest := httptest.NewRequest(http.MethodPost, "/", nil)
	validRequest.Body = ioutil.NopCloser(bytes.NewReader(uidJson))

	type args struct {
		rw http.ResponseWriter
		r  *http.Request
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "valid case",
			args: args{
				rw: httptest.NewRecorder(),
				r:  validRequest,
			},
			want: http.StatusOK,
		},

		{
			name: "invalid case",
			args: args{
				rw: httptest.NewRecorder(),
				r:  httptest.NewRequest(http.MethodPost, "/", ioutil.NopCloser(bytes.NewReader([]byte{}))),
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler(tt.args.rw, tt.args.r)

			response, ok := tt.args.rw.(*httptest.ResponseRecorder)
			assert.True(t, ok)
			assert.NotNil(t, response)

			assert.Equal(t, tt.want, response.Code)
		})
	}
}

func TestSaveMemberCoverToFirestoreHandler(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)
	if ctx == nil {
		t.Error("nil context")
		return
	}

	aut := ctx.Value(base.AuthTokenContextKey).(*auth.Token)
	if aut == nil {
		t.Errorf("nil auth token")
		return
	}

	srv := profile.NewService()
	if srv == nil {
		t.Errorf("nil profile service")
		return
	}

	handler := graph.SaveMemberCoverHandler(ctx, srv)

	type Payload struct {
		PayerName      string `json:"payerName"`
		MemberName     string `json:"memberName"`
		MemberNumber   string `json:"memberNumber"`
		PayerSladeCode int    `json:"payerSladeCode"`
		UID            string `json:"uid"`
	}

	type args struct {
		payload Payload
		rw      http.ResponseWriter
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "valid case",
			args: args{
				payload: Payload{
					PayerName:      "UAP",
					MemberName:     "Jakaya",
					MemberNumber:   "133",
					PayerSladeCode: 457,
					UID:            aut.UID,
				},
				rw: httptest.NewRecorder(),
			},
			want: http.StatusOK,
		},

		{
			name: "invalid case",
			args: args{
				payload: Payload{
					MemberName:     "Jak",
					MemberNumber:   "132",
					PayerName:      "APA",
					PayerSladeCode: 111,
					UID:            "",
				},
				rw: httptest.NewRecorder(),
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadJson, err := json.Marshal(tt.args.payload)
			if err != nil {
				t.Errorf("can't marshal payload to JSON")
				return
			}
			if payloadJson == nil {
				t.Errorf("nil JSON payload")
				return
			}

			request := httptest.NewRequest(http.MethodPost, "/", nil)
			request.Body = ioutil.NopCloser(bytes.NewReader(payloadJson))
			handler(tt.args.rw, request)

			response, ok := tt.args.rw.(*httptest.ResponseRecorder)
			if response == nil {
				t.Errorf("nil response")
				return
			}
			if !ok {
				t.Errorf(
					"expected response to be a *httptest.ResponseRecorder")
				return
			}

			if response.Code != tt.want {
				t.Errorf(
					"expected status code %d, got %d", tt.want, response.Code)

				data, err := ioutil.ReadAll(response.Body)
				if err != nil {
					t.Errorf("can't read response body")
					return
				}

				log.Printf("raw response data: \n%s\n", string(data))

				return
			}
		})
	}
}

func TestSendRetryOTPHandler(t *testing.T) {
	client := http.DefaultClient
	sendOTPRetry := profile.SendRetryOTP{
		Msisdn:    base.TestUserPhoneNumber,
		RetryStep: 1,
	}
	bs, err := json.Marshal(sendOTPRetry)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid generate and send retry OTPs request",
			args: args{
				url:        fmt.Sprintf("%s/send_retry_otp", baseURL),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid generate and send retry OTPs request",
			args: args{
				url:        fmt.Sprintf("%s/send_retry_otp", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range base.GetDefaultHeaders(t, baseURL, "profile") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}
		})
	}
}

func TestRetrieveUserProfileHandler(t *testing.T) {
	client := http.DefaultClient

	_, authToken := base.GetAuthenticatedContextAndToken(t)
	if authToken == nil {
		t.Errorf("nil auth token")
		return
	}

	bpUID := &profile.BusinessPartnerUID{
		UID: authToken.UID,
	}
	bs, err := json.Marshal(bpUID)
	if err != nil {
		t.Errorf("unable to marshal BP UID payload to JSON: %s", err)
		return
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid user profile retrieve request - valid UID",
			args: args{
				url: fmt.Sprintf(
					"%s/internal/retrieve_user_profile", baseURL),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid user profile retrieve request - nil body",
			args: args{
				url: fmt.Sprintf(
					"%s/internal/retrieve_user_profile", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range base.GetDefaultHeaders(t, baseURL, "profile") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf(
					"expected status %d, got %d and response %s",
					tt.wantStatus,
					resp.StatusCode,
					string(data),
				)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}
		})
	}
}

func TestPhoneSignIn(t *testing.T) {
	client := http.DefaultClient
	srv := profile.NewService()
	assert.NotNil(t, srv, "service is nil")

	ctx, _ := base.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}
	set, err := srv.SetUserPIN(ctx, base.TestUserPhoneNumber, base.TestUserPin)
	if !set {
		t.Errorf("setting a pin for test user failed. It returned false")
	}
	if err != nil {
		t.Errorf("setting a pin for test user failed: %v", err)
	}
	if !set {
		t.Errorf("setting a pin for test user failed. It returned false")
	}
	signIn := profile.PhoneSignInInput{
		PhoneNumber: base.TestUserPhoneNumber,
		Pin:         base.TestUserPin,
	}
	bs, err := json.Marshal(signIn)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	wrongPin := profile.PhoneSignInInput{
		PhoneNumber: base.TestUserPhoneNumber,
		Pin:         "4567",
	}
	w, err := json.Marshal(wrongPin)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	wrongPinPayload := bytes.NewBuffer(w)

	wrongPhone := profile.PhoneSignInInput{
		PhoneNumber: base.TestUserPhoneNumber,
		Pin:         "4567",
	}
	p, err := json.Marshal(wrongPhone)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	wrongPhonePayload := bytes.NewBuffer(p)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "successful sign in with phone number and PIN",
			args: args{
				url:        fmt.Sprintf("%s/msisdn_login", baseURL),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "unsuccessful sign in: nil data given",
			args: args{
				url:        fmt.Sprintf("%s/msisdn_login", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "unsuccessful sign in: wrong pin supplied",
			args: args{
				url:        fmt.Sprintf("%s/msisdn_login", baseURL),
				httpMethod: http.MethodPost,
				body:       wrongPinPayload,
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name: "unsuccessful sign in: wrong phone number supplied",
			args: args{
				url:        fmt.Sprintf("%s/msisdn_login", baseURL),
				httpMethod: http.MethodPost,
				body:       wrongPhonePayload,
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range base.GetDefaultHeaders(t, baseURL, "profile") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantErr && tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}
		})
	}
}

func TestSaveCoverPayloadHandler(t *testing.T) {
	client := http.DefaultClient

	emptyPayload := profile.SaveMemberCoverPayload{}
	emptyPayloadJSONBytes, err := json.Marshal(emptyPayload)
	if err != nil {
		t.Errorf("can't marshal empty save cover payload to JSON: %v", err)
		return
	}

	_, authToken := base.GetAuthenticatedContextAndToken(t)
	if authToken == nil {
		t.Errorf("nil auth token")
		return
	}

	coverRequest := &profile.SaveMemberCoverPayload{
		PayerName:      "Resolution Insurance Company Limited",
		MemberName:     "Daniel Ngure Nyaga",
		MemberNumber:   "1464409",
		PayerSladeCode: 458,
		UID:            authToken.UID,
	}
	coverRequestJSONBytes, err := json.Marshal(coverRequest)
	if err != nil {
		t.Errorf("unable to marshal cover request payload to JSON: %s", err)
		return
	}

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid save cover request - valid UID that exists",
			args: args{
				url: fmt.Sprintf(
					"%s/internal/save_cover", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(coverRequestJSONBytes),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid save cover retrieve request - nil body",
			args: args{
				url: fmt.Sprintf(
					"%s/internal/save_cover", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid save cover retrieve request - no UID",
			args: args{
				url: fmt.Sprintf(
					"%s/internal/save_cover", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(emptyPayloadJSONBytes),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range base.GetDefaultHeaders(t, baseURL, "profile") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf(
					"expected status %d, got %d and response %s",
					tt.wantStatus,
					resp.StatusCode,
					string(data),
				)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if !tt.wantErr {
				//  check response payload format
				var respPayload profile.SaveResponsePayload
				err = json.Unmarshal(data, &respPayload)
				if err != nil {
					log.Print(string(data))
					t.Errorf(
						"can't unmarshal save cover resp payload: %v", err)
					return
				}
				if !respPayload.SuccessfullySaved {
					t.Errorf(
						"expected successfullySaved to be true in the response")
					return
				}
			}
		})
	}
}

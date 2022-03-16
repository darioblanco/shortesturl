package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darioblanco/shortesturl/app/internal/cache"
	"github.com/darioblanco/shortesturl/app/internal/config"
	"github.com/darioblanco/shortesturl/app/internal/logging"
	"github.com/stretchr/testify/assert"
)

func createRequest(method, path string, input interface{}) *http.Request {
	var body bytes.Buffer
	if input != nil {
		reqBody, _ := json.Marshal(input)
		body = *bytes.NewBuffer(reqBody)
	}
	req, _ := http.NewRequest(method, path, &body)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func testRequest(
	t *testing.T,
	handler http.Handler,
	method,
	path string,
	input interface{},
	expectedStatusCode int,
	expectedResponse interface{},
) {
	ts := httptest.NewServer(handler)
	defer ts.Close()

	req := createRequest(method, ts.URL+path, input)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer resp.Body.Close()
	respBody := string(bytes.TrimRight(respBodyBytes, "\n"))

	assert.Equal(t, expectedStatusCode, resp.StatusCode)
	if expectedResponse != nil {
		inputBytes, err := json.Marshal(expectedResponse)
		assert.NoError(t, err)
		assert.Equal(t, string(inputBytes), respBody)
	}
}

func TestEncodeAndDecode(t *testing.T) {
	r, _ := NewRouter(
		context.Background(),
		&config.Values{
			HttpScheme:           "http",
			HttpHost:             "localhost",
			HttpPort:             80,
			UrlLength:            6,
			UrlExpirationInHours: 0,
		},
		logging.NewTest(t),
		cache.NewTest(),
	)
	t.Parallel()
	tests := []struct {
		longUrl  string
		shortUrl string
	}{
		{
			"https://github.com/darioblanco",
			"http://localhost/64fc5e",
		},
		{
			"https://github.com/darioblanco",
			"http://localhost/64fc5e",
		}, // It should always return the same shortened url
		{
			"https://darioblanco.com/?url=https%3A%2F%2Fgithub.com%2Fdarioblanco",
			"http://localhost/104b16",
		},
		{
			"https://darioblanco.com/?url=https://github.com/darioblanco",
			"http://localhost/eaf668",
		},
		{
			"http://münchen.darío.localhost",
			"http://localhost/973335",
		},
		{
			"https://longurlmaker.com/go?id=5loftyv0SimURLrangy70URLcut3EzURL0faraway31lanky6URL0b10running3runningUlimit02040fstringyewTinyURL14elongate94stretchingdistantc1lx01a5lastingyocURLvi0A2N206DigBig14stowering35spunZout301URLbnBeamZtosustainedG8L067URLcut0PiURL0111e4100expanded355EzURLprotracted95NutshellURL98741bbjgreat786c06longish513d6yXilDwarfurl91SnipURL2SHurlenduringc6s1118e6URLcfaraway0lankyEasyURLaURlZiej6p2z279oYepItTraceURLIsZgd1extensiveCanURL17great02TightURLremote5rangy7hlongish6h0continuedcremoteb3dbRedirxtall2highdeep7e2longishrangy20ygwShrinkrDigBig1w6Redirx016distant4TinyLinkXZseucd0n0prolongedenlargedeUlimitelongated11SimURLdrawnZoutrunningeySmallr5lingering20elongated07B65continuedyShrtnd8FwdURLjvht0c0XZsegreat6URLZcoZuk54stretchedxelongated0EzURL10lnkZin40fURLZcoZuk46spunZoutegdstretching00f5farZoff0Smallrk1sustainedShortURL1TightURLc036farawaytowering9remote9511p1TinyLink5321yexpanded610h4lengthened2runningTraceURL00DigBig61251outstretchedtowering3GetShortyMiniliencMyURL6211ShortURL03URLCutterm069NotLonge1g510ufarZreachingrangyb0toutstretched1FhURL2XZse1507x56814stretched1m10a1dDecentURL4b910farZoffCanURLextensiver314zcf70stretchedrangylengthened0Shim0EzURLfaraway9021U76outstretchedspreadZoutstretchTraceURL2NutshellURL9Metamark73e1335rXZsefarZoff92c6continuedlasting046U767zstringy1dvShortURL53remotea126URLbdistanttXilFhURLfarawayWapURL83091340deep6g150farZreachinggangling10Ne1d1tallURLPieSHurlURLHawk95URLvir04y601e8t08SimURLNe1f1b1383c95farZoff0EzURL0eURLHawk65enduring1fx5lastingv09d12poutstretched86fcontinuedeURLPie65longishlengthened337t651a3prolonged1YepIt440wspunZout2781URLHawk6329c0NotLongexpandedvFwdURLURLCutter301URL310301URLBeamZto70295651Shrtnd8extensive9benduringtall00stretchedv880f1041002nh6CanURL114NutshellURLDwarfurlrbShortenURLDoiop266distantim304ShrinkURL1high8008Ne1SimURLuflankyXZselankyTinyURLd04URlZie0greatoutstretched7TraceURL1v45n7Fly2farawayzstretchURLCutter70lanky0lingering0A2N0Shrtndrunningb390tc0Sitelutions05BeamZtodistantkfarZreachingcb04lingering3vcb0lankynShim3b5ShredURL4n0ShortURL0B6551071f1spreadZo",
			"http://localhost/5f7e9f",
		},
	}
	for _, tt := range tests {
		tt := tt // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.longUrl, func(t *testing.T) {
			t.Parallel()
			testRequest(t, r,
				http.MethodPost,
				"/encode",
				URLPayload{URL: tt.longUrl},
				http.StatusOK,
				URLPayload{URL: tt.shortUrl},
			)
			testRequest(t, r,
				http.MethodPost,
				"/decode",
				URLPayload{URL: tt.shortUrl},
				http.StatusOK,
				URLPayload{URL: tt.longUrl},
			)
		})
	}
}

func TestEncode_OKWithCollisions(t *testing.T) {
	mr, client := cache.NewMiniredis()
	mr.Set("64fc5e4dfd2a", "val1")
	mr.Set("4fc5e4dfd2a0", "val2")
	mr.Set("fc5e4dfd2a0f", "val3")
	mr.Set("c5e4dfd2a0ff", "val4")
	mr.Set("5e4dfd2a0ff5", "val5")
	r, _ := NewRouter(
		context.Background(),
		&config.Values{
			HttpScheme:           "http",
			HttpHost:             "localhost",
			HttpPort:             3000,
			UrlLength:            12,
			UrlExpirationInHours: 0,
		},
		logging.NewTest(t),
		client,
	)
	testRequest(t, r,
		http.MethodPost,
		"/encode",
		URLPayload{URL: "https://github.com/darioblanco"},
		http.StatusOK,
		URLPayload{URL: "http://localhost:3000/e4dfd2a0ff5e"},
	)
}

func TestEncode_BadRequest(t *testing.T) {
	r, _ := NewRouter(
		context.Background(),
		&config.Values{},
		logging.NewTest(t),
		cache.NewTest(),
	)
	testRequest(t, r,
		http.MethodPost,
		"/encode",
		URLPayload{URL: "wrong-url"},
		http.StatusBadRequest,
		ErrHTTPResponse{
			StatusText: http.StatusText(http.StatusBadRequest),
			ErrorText:  "invalid http/https url format",
		},
	)
}

func TestEncode_InternalServerError(t *testing.T) {
	mr, client := cache.NewMiniredis()
	mr.SetError("mock error")
	r, _ := NewRouter(
		context.Background(),
		&config.Values{},
		logging.NewTest(t),
		client,
	)
	testRequest(t, r,
		http.MethodPost,
		"/encode",
		URLPayload{URL: "https://github.com/darioblanco"},
		http.StatusInternalServerError,
		ErrHTTPResponse{
			StatusText: http.StatusText(http.StatusInternalServerError),
			ErrorText:  "oops, something went wrong in our side",
		},
	)
}

func TestDecode_NotFound(t *testing.T) {
	r, _ := NewRouter(
		context.Background(),
		&config.Values{
			HttpScheme:           "https",
			HttpHost:             "localhost",
			HttpPort:             443,
			UrlLength:            6,
			UrlExpirationInHours: 0,
		},
		logging.NewTest(t),
		cache.NewTest(),
	)
	testRequest(t, r,
		http.MethodPost,
		"/decode",
		URLPayload{URL: "http://localhost:3000/abcdef"},
		http.StatusNotFound,
		ErrHTTPResponse{
			StatusText: http.StatusText(http.StatusNotFound),
			ErrorText:  "long url not found",
		},
	)
}

func TestDecode_BadRequest(t *testing.T) {
	r, _ := NewRouter(
		context.Background(),
		&config.Values{},
		logging.NewTest(t),
		cache.NewTest(),
	)
	testRequest(t, r,
		http.MethodPost,
		"/decode",
		URLPayload{URL: "wrong-url"},
		http.StatusBadRequest,
		ErrHTTPResponse{
			StatusText: http.StatusText(http.StatusBadRequest),
			ErrorText:  "invalid http/https url format",
		},
	)
}

func TestDecode_InternalServerError(t *testing.T) {
	mr, client := cache.NewMiniredis()
	mr.SetError("mock error")
	r, _ := NewRouter(
		context.Background(),
		&config.Values{},
		logging.NewTest(t),
		client,
	)
	testRequest(t, r,
		http.MethodPost,
		"/decode",
		URLPayload{URL: "http://localhost:3000/64fc5e"},
		http.StatusInternalServerError,
		ErrHTTPResponse{
			StatusText: http.StatusText(http.StatusInternalServerError),
			ErrorText:  "oops, something went wrong in our side",
		},
	)
}

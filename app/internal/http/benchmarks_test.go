package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darioblanco/shortesturl/app/internal/cache"
	"github.com/darioblanco/shortesturl/app/internal/config"
	"github.com/darioblanco/shortesturl/app/internal/logging"
)

func createAPI(b *testing.B, urlLength int) api {
	return api{
		cache: cache.NewTest(),
		config: &config.Values{
			HttpScheme:           "http",
			HttpHost:             "localhost",
			HttpPort:             80,
			UrlLength:            urlLength,
			UrlExpirationInHours: 0,
		},
		logger: logging.NewTest(b),
	}
}

func benchmarkEncode(b *testing.B, url string, urlLength int) {
	a := createAPI(b, urlLength)
	w := httptest.NewRecorder()
	r := createRequest(
		http.MethodPost,
		"/encode",
		URLPayload{URL: url},
	)
	b.ResetTimer()
	a.Encode(w, r)
}

func BenchmarkEncode_5(b *testing.B) {
	benchmarkEncode(b, "https://github.com/darioblanco", 5)
}

func BenchmarkEncode_6(b *testing.B) {
	benchmarkEncode(b, "https://github.com/darioblanco", 6)
}

func BenchmarkEncode_6LongURL(b *testing.B) {
	benchmarkEncode(b, "https://longurlmaker.com/go?id=5loftyv0SimURLrangy70URLcut3EzURL0faraway31lanky6URL0b10running3runningUlimit02040fstringyewTinyURL14elongate94stretchingdistantc1lx01a5lastingyocURLvi0A2N206DigBig14stowering35spunZout301URLbnBeamZtosustainedG8L067URLcut0PiURL0111e4100expanded355EzURLprotracted95NutshellURL98741bbjgreat786c06longish513d6yXilDwarfurl91SnipURL2SHurlenduringc6s1118e6URLcfaraway0lankyEasyURLaURlZiej6p2z279oYepItTraceURLIsZgd1extensiveCanURL17great02TightURLremote5rangy7hlongish6h0continuedcremoteb3dbRedirxtall2highdeep7e2longishrangy20ygwShrinkrDigBig1w6Redirx016distant4TinyLinkXZseucd0n0prolongedenlargedeUlimitelongated11SimURLdrawnZoutrunningeySmallr5lingering20elongated07B65continuedyShrtnd8FwdURLjvht0c0XZsegreat6URLZcoZuk54stretchedxelongated0EzURL10lnkZin40fURLZcoZuk46spunZoutegdstretching00f5farZoff0Smallrk1sustainedShortURL1TightURLc036farawaytowering9remote9511p1TinyLink5321yexpanded610h4lengthened2runningTraceURL00DigBig61251outstretchedtowering3GetShortyMiniliencMyURL6211ShortURL03URLCutterm069NotLonge1g510ufarZreachingrangyb0toutstretched1FhURL2XZse1507x56814stretched1m10a1dDecentURL4b910farZoffCanURLextensiver314zcf70stretchedrangylengthened0Shim0EzURLfaraway9021U76outstretchedspreadZoutstretchTraceURL2NutshellURL9Metamark73e1335rXZsefarZoff92c6continuedlasting046U767zstringy1dvShortURL53remotea126URLbdistanttXilFhURLfarawayWapURL83091340deep6g150farZreachinggangling10Ne1d1tallURLPieSHurlURLHawk95URLvir04y601e8t08SimURLNe1f1b1383c95farZoff0EzURL0eURLHawk65enduring1fx5lastingv09d12poutstretched86fcontinuedeURLPie65longishlengthened337t651a3prolonged1YepIt440wspunZout2781URLHawk6329c0NotLongexpandedvFwdURLURLCutter301URL310301URLBeamZto70295651Shrtnd8extensive9benduringtall00stretchedv880f1041002nh6CanURL114NutshellURLDwarfurlrbShortenURLDoiop266distantim304ShrinkURL1high8008Ne1SimURLuflankyXZselankyTinyURLd04URlZie0greatoutstretched7TraceURL1v45n7Fly2farawayzstretchURLCutter70lanky0lingering0A2N0Shrtndrunningb390tc0Sitelutions05BeamZtodistantkfarZreachingcb04lingering3vcb0lankynShim3b5ShredURL4n0ShortURL0B6551071f1spreadZo", 8)
}

func BenchmarkEncode_7(b *testing.B) {
	benchmarkEncode(b, "https://github.com/darioblanco", 7)
}

func BenchmarkEncode_8(b *testing.B) {
	benchmarkEncode(b, "https://github.com/darioblanco", 8)
}

func BenchmarkDecode_6(b *testing.B) {
	a := createAPI(b, 6)
	w := httptest.NewRecorder()
	r := createRequest(
		http.MethodPost,
		"/decode",
		URLPayload{URL: "http://localhost/64fc5e"},
	)
	b.ResetTimer()
	a.Decode(w, r)
}

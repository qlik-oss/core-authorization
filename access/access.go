package access

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/qlik-oss/enigma-go"
)

const (
	trafficDumpFile = "traffic-dump.json"
)

var (
	secret = []byte("secret")
	ctx    = context.Background()
)

// connect connects to an engine at the proviced URL and claims.
func connect(engineURL string, jwtClaims jwt.MapClaims) (*enigma.Global, error) {
	ctx := context.Background()
	headers := make(http.Header)

	if jwtClaims != nil {
		signedJwt, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims).SignedString(secret)
		headers.Set("Authorization", fmt.Sprintf("Bearer %s", signedJwt))
	}

	return enigma.Dialer{TrafficDumpFile: trafficDumpFile}.Dial(ctx, engineURL, headers)
}

// createApp creates an app using the provided Global object and app name.
func createApp(global *enigma.Global, appName string) (*enigma.Doc, error) {
	return global.CreateDocEx(ctx, appName, "", "", "", "Main")
}

// openApp opens an app using the providad Global object and app name.
func openApp(global *enigma.Global, appName string) (*enigma.Doc, error) {
	return global.OpenDoc(ctx, appName, "", "", "", false)
}

// createSheet creates a sheet object using the provided Doc objecgt and sheet name.
func createSheet(doc *enigma.Doc, name string) (*enigma.GenericObject, error) {
	msg := json.RawMessage(`{
		"title": "/title",
		"description": "/description",
		"meta": "/meta",
		"order": "/order",
		"type": "/qInfo/qType",
		"id": "/qInfo/qId",
		"lb": "/qListObjectDef",
		"hc": "/qHyperCubeDef"
	}`)

	props := enigma.GenericObjectProperties{
		Info: &enigma.NxInfo{Id: name, Type: "sheet"},
		ChildListDef: &enigma.ChildListDef{
			Data: msg,
		},
	}

	return doc.CreateObject(ctx, &props)
}

// getSheetObject returns a sheet in the proviced Doc object with the provided name.
func getSheetObject(doc *enigma.Doc, name string) (*enigma.GenericObject, error) {
	return doc.GetObject(ctx, name)
}

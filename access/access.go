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

func connect(engineURL string, jwtClaims jwt.MapClaims) (*enigma.Global, error) {
	ctx := context.Background()
	headers := make(http.Header)

	if jwtClaims != nil {
		signedJwt, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims).SignedString(secret)
		headers.Set("Authorization", fmt.Sprintf("Bearer %s", signedJwt))
	}

	return enigma.Dialer{TrafficDumpFile: trafficDumpFile}.Dial(ctx, engineURL, headers)
}

func createApp(global *enigma.Global, appName string) (string, error) {
	_, appID, err := global.CreateApp(ctx, appName, "Main")
	return appID, err
}

func openApp(global *enigma.Global, appName string) (*enigma.Doc, error) {
	return global.OpenDoc(ctx, appName, "", "", "", false)
}

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

func readSheetLayout(doc *enigma.Doc, name string) (*enigma.GenericObject, *enigma.GenericObjectLayout, error) {
	object, err := doc.GetObject(ctx, name)
	if err != nil {
		return nil, nil, err
	}

	layout, err := object.GetLayout(ctx)
	if err != nil {
		return object, nil, err
	}

	return object, layout, err
}

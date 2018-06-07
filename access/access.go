package access

import (
	"context"
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

func accessEngineAppObject(engineURL string, docName string, objectID string, jwtClaims jwt.MapClaims) error {
	headers := make(http.Header)

	if jwtClaims != nil {
		signedJwt, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims).SignedString(secret)
		headers.Set("Authorization", fmt.Sprintf("Bearer %s", signedJwt))
	}

	global, err := enigma.Dialer{TrafficDumpFile: trafficDumpFile}.Dial(ctx, engineURL, headers)
	if err != nil {
		return err
	}
	defer global.DisconnectFromServer()

	doc, err := global.OpenDoc(ctx, docName, "", "", "", false)
	if err != nil {
		return err
	}

	countryListObject, err := doc.GetObject(ctx, objectID)
	if err != nil {
		return err
	}

	_, err = countryListObject.GetLayout(ctx)
	return err
}

func updateEngineAppObject(engineURL string, docName string, objectID string, jwtClaims jwt.MapClaims) error {
	headers := make(http.Header)

	if jwtClaims != nil {
		signedJwt, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims).SignedString(secret)
		headers.Set("Authorization", fmt.Sprintf("Bearer %s", signedJwt))
	}

	global, err := enigma.Dialer{TrafficDumpFile: trafficDumpFile}.Dial(ctx, engineURL, headers)
	if err != nil {
		return err
	}
	defer global.DisconnectFromServer()

	doc, err := global.OpenDoc(ctx, docName, "", "", "", false)
	if err != nil {
		return err
	}

	countryListObject, err := doc.GetObject(ctx, objectID)
	if err != nil {
		return err
	}

	_, err = countryListObject.GetLayout(ctx)
	return err
}

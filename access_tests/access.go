package access_tests

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/qlik-oss/enigma-go"
	"context"
	"net/http"
	"fmt"
)

var (
	secret = []byte("secret")
	ctx    = context.Background()
)

func accessEngineAppObjectWithJwt(engine string, app string, object string, jwtClaims jwt.MapClaims) error {

	headers := make(http.Header)
	if  jwtClaims != nil {
		signedJwt, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims).SignedString(secret)
		headers.Set("Authorization", fmt.Sprintf("Bearer %s", signedJwt))
	}
	global, err := enigma.Dialer{TrafficDumpFile: "./traffic.txt"}.Dial(ctx, engine, headers)
	if err != nil {
		return err
	}
	defer global.DisconnectFromServer()

	doc, err := global.OpenDoc(ctx, app, "", "", "", false)
	if err != nil {
		return err
	}

	countryListObject, err := doc.GetObject(ctx, object)
	if err != nil {
		return err
	}
	_, err = countryListObject.GetLayout(ctx)

	if err != nil {
		return err
	}
	return err
}

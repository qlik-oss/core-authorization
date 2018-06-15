package access

import (
	"context"
	"encoding/json"
	"errors"
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
func openApp(global *enigma.Global, appID string) (*enigma.Doc, error) {
	return global.OpenDoc(ctx, appID, "", "", "", false)
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

// readObject returns ...
func readObject(doc *enigma.Doc, name string) (*enigma.GenericObject, error) {
	return doc.GetObject(ctx, name)
}

// reloadMoviesData ...
func reloadMoviesData(doc *enigma.Doc) error {
	const script = "LOAD * FROM '/data/movies.csv' (txt, embedded labels, delimiter is ',', no quotes);"
	var (
		err      error
		reloadOk bool
	)

	err = doc.SetScript(ctx, script)
	if err != nil {
		return err
	}

	reloadOk, err = doc.DoReload(ctx, 0, false, false)
	if err != nil {
		return err
	}
	if !reloadOk {
		return errors.New("reload failed")
	}

	return nil
}

// createMoviesObject ...
func createMoviesObject(doc *enigma.Doc, name string) (*enigma.GenericObject, error) {
	props := &enigma.GenericObjectProperties{
		Info: &enigma.NxInfo{Id: name, Type: "movies"},
		HyperCubeDef: &enigma.HyperCubeDef{
			Dimensions: []*enigma.NxDimension{
				{Def: &enigma.NxInlineDimensionDef{FieldDefs: []string{"Movie"}}},
				{Def: &enigma.NxInlineDimensionDef{FieldDefs: []string{"Year"}}},
				{Def: &enigma.NxInlineDimensionDef{FieldDefs: []string{"Adjusted Cost"}}},
				{Def: &enigma.NxInlineDimensionDef{FieldDefs: []string{"Image"}}},
			},
			InitialDataFetch: []*enigma.NxPage{{Height: 50, Width: 4}},
		},
	}

	obj, err := doc.CreateObject(ctx, props)
	if err != nil {
		return nil, err
	}

	err = doc.SaveObjects(ctx)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func readMoviesData(obj *enigma.GenericObject) (titles []string, years []float64, err error) {
	var layout *enigma.GenericObjectLayout

	layout, err = obj.GetLayout(ctx)
	if err != nil {
		return
	}

	titles = make([]string, 0, 50)
	years = make([]float64, 0, 50)
	movies := layout.HyperCube.DataPages[0]

	for _, m := range movies.Matrix {
		titles = append(titles, m[0].Text)
		years = append(years, float64(m[1].Num))
	}

	return
}

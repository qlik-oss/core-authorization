package access

import (
	"strings"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	enigma "github.com/qlik-oss/enigma-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	engineURL = "ws://localhost:9176"
	sheetID   = "SHEET01"
)

var (
	adminClaims = jwt.MapClaims{
		"sub":   "someAdminUser",
		"roles": []string{"admin"},
	}
	nonAdminClaims = jwt.MapClaims{
		"sub": "someNonAdminUser",
	}
	global *enigma.Global
	app    *enigma.Doc
	sheet  *enigma.GenericObject
	info   *enigma.NxInfo
	id     string
	err    error
)

// Tests that user with admin rights can crate app and app object (sheet).
func TestAdminAccess(t *testing.T) {
	// Connnect to engine with admin rights.
	global, err = connect(engineURL, adminClaims)
	require.NoError(t, err)
	defer global.DisconnectFromServer()

	list, err := global.GetDocList(ctx)
	assert.NotNil(t, list)

	// Admin should be able to create an app.
	app, err = createApp(global, t.Name())
	require.NoError(t, err)

	// Admin should be able to open an app and create a sheet.
	app, err = openApp(global, app.GenericId)
	require.NoError(t, err)
	sheet, err = createSheet(app, sheetID)
	require.NoError(t, err)
	info, err = sheet.GetInfo(ctx)
	require.NoError(t, err)
	assert.Equal(t, sheetID, info.Id)

	// Save the objects for future access.
	app.SaveObjects(ctx)
}

// Tests that non-admin user can read app and app object (sheet), but cannot create app or app object.
func TestNonAdminAccess(t *testing.T) {
	// First create the app and app object with admin rights.
	app, sheet, err = adminCreateAppWithSheetHelper(t.Name())
	require.NoError(t, err)

	// Connect to engine with non-admin rights.
	global, err = connect(engineURL, nonAdminClaims)
	require.NoError(t, err)
	defer global.DisconnectFromServer()

	// Non-admin should be able to open an app.
	require.NoError(t, err)
	app, err = openApp(global, app.GenericId)
	require.NoError(t, err)

	// Non-admin should be able to get (read) a sheet object.
	sheet, err = getSheetObject(app, sheetID)
	require.NoError(t, err)
	info, err = sheet.GetInfo(ctx)
	require.NoError(t, err)
	assert.Equal(t, sheetID, info.Id)

	// Non-admin should not be able to create a sheet.
	_, err = createSheet(app, sheetID)
	require.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "access denied")

	// Non-admin should not be able to create an app.
	_, err = createApp(global, "CannotCreateThisApp")
	require.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "access denied")
}

func adminCreateAppWithSheetHelper(appName string) (*enigma.Doc, *enigma.GenericObject, error) {
	global, err = connect(engineURL, adminClaims)
	if err != nil {
		return nil, nil, err
	}

	defer global.DisconnectFromServer()

	app, err := createApp(global, appName)
	if err != nil {
		return nil, nil, err
	}

	sheet, err = createSheet(app, sheetID)
	if err != nil {
		return nil, nil, err
	}

	err = app.SaveObjects(ctx)
	if err != nil {
		return nil, nil, err
	}

	return app, sheet, nil
}

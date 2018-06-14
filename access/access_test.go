package access

import (
	"strings"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	enigma "github.com/qlik-oss/enigma-go"
	"github.com/stretchr/testify/require"
)

const (
	emptyEngineURL  = "ws://localhost:9176"
	reloadEngineURL = "ws://localhost:9276"
	sheetID         = "SHEET01"
	script          = "LOAD * FROM '/data/movies.csv' (txt, embedded labels, delimiter is ',', no quotes);"
)

var (
	adminClaims = jwt.MapClaims{
		"sub":   "someAdminUser",
		"roles": []string{"admin"},
	}
	nonAdminClaims = jwt.MapClaims{
		"sub": "someNonAdminUser",
	}
	createClaims = jwt.MapClaims{
		"sub":         "someViewUser",
		"allowCreate": true,
	}
	reloadClaims = jwt.MapClaims{
		"sub":         "someViewUser",
		"allowReload": true,
	}
	viewClaims = jwt.MapClaims{
		"sub":       "someViewUser",
		"allowView": true,
	}
	global *enigma.Global
	app    *enigma.Doc
	sheet  *enigma.GenericObject
	info   *enigma.NxInfo
	err    error
)

// ...
func TestAccessEmptyEngine(t *testing.T) {
	func() {
		// Connnect to engine with admin access.
		global, err = connect(emptyEngineURL, adminClaims)
		require.NoError(t, err)
		defer global.DisconnectFromServer()

		// Run test cases for admin user.
		t.Run("admin user should be able to create app", testCanCreateApp)
		t.Run("admin user should be able to open app", testCanOpenApp)
		t.Run("admin user should be able to create sheet", testCanCreateSheet)
	}()

	func() {
		// Connect to engine with non-admin access.
		global, err = connect(emptyEngineURL, nonAdminClaims)
		require.NoError(t, err)
		defer global.DisconnectFromServer()

		// Run test cases for non-admin user.
		t.Run("non-admin user should be able to open app", testCanOpenApp)
		t.Run("non-admin user should be able to read sheet", testCanReadSheet)
		t.Run("non-admin user should not be able to create sheet", testCannotCreateSheet)
		t.Run("non-admin user should bot be able to create app", testCannotCreateApp)
	}()
}

func TestAccessReloadEngine(t *testing.T) {
	func() {
		// Connnect to engine with create access.
		global, err = connect(reloadEngineURL, createClaims)
		require.NoError(t, err)
		defer global.DisconnectFromServer()

		// Run test cases for user with create access.
		t.Run("user with create access should be able to create app", testCanCreateApp)
		t.Run("user with create access should not be able to reload app", testCannotReloadApp)
	}()

	func() {
		// Connect to engine with reload access.
		global, err = connect(reloadEngineURL, reloadClaims)
		require.NoError(t, err)
		defer global.DisconnectFromServer()

		// Run test cases for user with reload access.
		t.Run("user with reload access should not be able to create app", testCannotCreateApp)
		t.Run("user with reload access should be able to open app", testCanOpenApp)
		t.Run("user with reload access should be able to reload app", testCanReloadApp)
	}()

	func() {
		// Connect to engine with view access.
		global, err = connect(reloadEngineURL, viewClaims)
		require.NoError(t, err)
		defer global.DisconnectFromServer()

		// Run test cases for user with view access.
		t.Run("user with view access should not be able to create app", testCannotCreateApp)
		t.Run("user with view access should be able to open app", testCanOpenApp)
		t.Run("user with view access should be able to ...", testCanReadObject)
		t.Run("user with view access should not be able to create sheet", testCannotCreateSheet)
		t.Run("user with view access should not be able to reload app", testCannotReloadApp)
	}()
}

//
// Test functions for "empty-enigne" rules.
//

func testCanCreateApp(t *testing.T) {
	app, err = createApp(global, t.Name())
	require.NoError(t, err)
}

func testCannotCreateApp(t *testing.T) {
	_, err = createApp(global, "CannotCreateThisApp")
	require.Error(t, err)
	require.Contains(t, strings.ToLower(err.Error()), "access denied")
}

func testCanOpenApp(t *testing.T) {
	app, err = openApp(global, app.GenericId)
	require.NoError(t, err)
}

func testCanCreateSheet(t *testing.T) {
	sheet, err = createSheet(app, sheetID)
	require.NoError(t, err)

	info, err = sheet.GetInfo(ctx)
	require.NoError(t, err)
	require.Equal(t, sheetID, info.Id)

	err = app.SaveObjects(ctx)
	require.NoError(t, err)
}

func testCannotCreateSheet(t *testing.T) {
	_, err = createSheet(app, sheetID)
	require.Error(t, err)
	require.Contains(t, strings.ToLower(err.Error()), "access denied")
}

func testCanReadSheet(t *testing.T) {
	sheet, err = getSheetObject(app, sheetID)
	require.NoError(t, err)

	info, err = sheet.GetInfo(ctx)
	require.NoError(t, err)
	require.Equal(t, sheetID, info.Id)
}

func testCanReadObject(t *testing.T) {
	field, err := app.GetField(ctx, "Movie", "")
	info, err := field.GetInfo(ctx)
	require.NoError(t, err)
	require.NotNil(t, info)
}

func testCanReloadApp(t *testing.T) {
	err = app.SetScript(ctx, script)
	require.NoError(t, err)

	reloadOk, err := app.DoReload(ctx, 0, false, false)
	require.NoError(t, err)
	require.True(t, reloadOk)
}

func testCannotReloadApp(t *testing.T) {
	err = app.SetScript(ctx, script)
	require.Error(t, err)
	require.Contains(t, strings.ToLower(err.Error()), "access denied")
}

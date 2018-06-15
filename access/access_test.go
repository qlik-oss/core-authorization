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
	appName         = "APP01"
	sheetID         = "SHEET01"
	moviesID        = "MOVIES01"
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
		"sub":         "someCreateUser",
		"allowCreate": true,
	}
	reloadClaims = jwt.MapClaims{
		"sub":         "someReloadUser",
		"allowReload": true,
	}
	viewClaims = jwt.MapClaims{
		"sub":       "someViewUser",
		"allowView": true,
	}
	err error
)

var state struct {
	global *enigma.Global
	app    *enigma.Doc
	sheet  *enigma.GenericObject
	movies *enigma.GenericObject
	appID  string
}

// ...
func TestAccessEmptyEngine(t *testing.T) {
	func() {
		// Connnect to engine with admin access.
		state.global, err = connect(emptyEngineURL, adminClaims)
		require.NoError(t, err)
		defer state.global.DisconnectFromServer()

		// Run test cases for admin user.
		t.Run("admin user should be able to create app", testCanCreateApp)
		t.Run("admin user should be able to open app", testCanOpenApp)
		t.Run("admin user should be able to create sheet", testCanCreateSheet)
	}()

	func() {
		// Connect to engine with non-admin access.
		state.global, err = connect(emptyEngineURL, nonAdminClaims)
		require.NoError(t, err)
		defer state.global.DisconnectFromServer()

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
		state.global, err = connect(reloadEngineURL, createClaims)
		require.NoError(t, err)
		defer state.global.DisconnectFromServer()

		// Run test cases for user with create access.
		t.Run("user with create access should be able to create app", testCanCreateApp)
		t.Run("user with create access should be able to create movies object", testCanCreateMoviesObject)
		t.Run("user with create access should not be able to reload app", testCannotReloadApp)
	}()

	func() {
		// Connect to engine with reload access.
		state.global, err = connect(reloadEngineURL, reloadClaims)
		require.NoError(t, err)
		defer state.global.DisconnectFromServer()

		// Run test cases for user with reload access.
		t.Run("user with reload access should not be able to create app", testCannotCreateApp)
		t.Run("user with reload access should be able to open app", testCanOpenApp)
		t.Run("user with reload access should be able to reload app", testCanReloadApp)
		t.Run("user with reload access should be able to save app", testCanSaveApp)
	}()

	func() {
		// Connect to engine with view access.
		state.global, err = connect(reloadEngineURL, viewClaims)
		require.NoError(t, err)
		defer state.global.DisconnectFromServer()

		// Run test cases for user with view access.
		t.Run("user with view access should not be able to create app", testCannotCreateApp)
		t.Run("user with view access should be able to open app", testCanOpenApp)
		t.Run("user with view access should be able to read movies object", testCanReadMoviesObject)
		t.Run("user with view access should be able to read movies data", testCanReadMoviesData)
		t.Run("user with view access should not be able to create sheet", testCannotCreateSheet)
		t.Run("user with view access should not be able to reload app", testCannotReloadApp)
	}()
}

func testCanCreateApp(t *testing.T) {
	state.app, err = createApp(state.global, appName)
	require.NoError(t, err)
	state.appID = state.app.GenericId
}

func testCannotCreateApp(t *testing.T) {
	_, err = createApp(state.global, "CannotCreateThisApp")
	require.Error(t, err)
	require.Contains(t, strings.ToLower(err.Error()), "access denied")
}

func testCanOpenApp(t *testing.T) {
	state.app, err = openApp(state.global, state.appID)
	require.NoError(t, err)
}

func testCanSaveApp(t *testing.T) {
	err = state.app.DoSave(ctx, "")
	require.NoError(t, err)
}

func testCanCreateSheet(t *testing.T) {
	state.sheet, err = createSheet(state.app, sheetID)
	require.NoError(t, err)

	err = state.app.SaveObjects(ctx)
	require.NoError(t, err)
}

func testCannotCreateSheet(t *testing.T) {
	_, err = createSheet(state.app, sheetID)
	require.Error(t, err)
	require.Contains(t, strings.ToLower(err.Error()), "access denied")
}

func testCanReadSheet(t *testing.T) {
	state.sheet, err = readObject(state.app, sheetID)
	require.NoError(t, err)
}

func testCanCreateMoviesObject(t *testing.T) {
	state.movies, err = createMoviesObject(state.app, moviesID)
	require.NoError(t, err)
}

func testCanReadMoviesObject(t *testing.T) {
	state.movies, err = readObject(state.app, moviesID)
	require.NoError(t, err)
}

func testCanReadMoviesData(t *testing.T) {
	var (
		titles []string
		years  []float64
	)
	titles, years, err = readMoviesData(state.movies)
	require.NoError(t, err)
	require.Equal(t, "Armageddon", titles[1])
	require.Equal(t, float64(1998), years[1])
}

func testCanReloadApp(t *testing.T) {
	err = reloadMoviesData(state.app)
	require.NoError(t, err)
}

func testCannotReloadApp(t *testing.T) {
	err = reloadMoviesData(state.app)
	require.Error(t, err)
	require.Contains(t, strings.ToLower(err.Error()), "access denied")
}

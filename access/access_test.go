// This file contains two test suites that execute tests on the "empty-engine" and "reload-engine" instances.
// Each test suite consists of calls to different access functions. These functions are expected to succeed or fail
// depending on which access rights users are given and on the rules provided to each engine instance. Each test case
// verifies that access is granted or denied properly.

package access

import (
	"strings"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	enigma "github.com/qlik-oss/enigma-go"
	"github.com/stretchr/testify/require"
)

const (
	emptyEngineURL  = "ws://empty-engine:9076"
	reloadEngineURL = "ws://reload-engine:9076"
	appName         = "APP01"
	sheetID         = "SHEET01"
	moviesID        = "MOVIES01"
)

// Claims for the "empty-engine" instance.
var (
	adminClaims = jwt.MapClaims{
		"sub":   "someAdminUser",
		"roles": []string{"admin"},
	}
	nonAdminClaims = jwt.MapClaims{
		"sub": "someNonAdminUser",
	}
)

// Claims for the "reload-engine" instance.
var (
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

// testState holds the global test state.
var testState struct {
	global *enigma.Global
	app    *enigma.Doc
	sheet  *enigma.GenericObject
	movies *enigma.GenericObject
	appID  string
}

// TestAccessEmptyEngine executes the test suite against the "empty-engine" engine.
// Test cases depend on the order of execution.
// Test cases make assumptions on the contents of the global testState variable.
func TestAccessEmptyEngine(t *testing.T) {
	func() {
		// Connnect to engine with admin access.
		testState.global, err = connect(emptyEngineURL, adminClaims)
		require.NoError(t, err)
		defer testState.global.DisconnectFromServer()

		// Run test cases for admin user.
		t.Run("admin user should be able to create app", testCanCreateApp)
		t.Run("admin user should be able to open app", testCanOpenApp)
		t.Run("admin user should be able to create sheet", testCanCreateSheet)
	}()

	func() {
		// Connect to engine with non-admin access.
		testState.global, err = connect(emptyEngineURL, nonAdminClaims)
		require.NoError(t, err)
		defer testState.global.DisconnectFromServer()

		// Run test cases for non-admin user.
		t.Run("non-admin user should be able to open app", testCanOpenApp)
		t.Run("non-admin user should be able to read sheet", testCanReadSheet)
		t.Run("non-admin user should not be able to create sheet", testCannotCreateSheet)
		t.Run("non-admin user should bot be able to create app", testCannotCreateApp)
	}()
}

// TestAccessEmptyEngine executes the test suite against the "reload-engine" engine.
// Test cases depend on the order of execution.
// Test cases make assumptions on the contents of the global testState variable.
func TestAccessReloadEngine(t *testing.T) {
	func() {
		// Connnect to engine with create access.
		testState.global, err = connect(reloadEngineURL, createClaims)
		require.NoError(t, err)
		defer testState.global.DisconnectFromServer()

		// Run test cases for user with create access.
		t.Run("user with create access should be able to create app", testCanCreateApp)
		t.Run("user with create access should be able to create movies object", testCanCreateMoviesObject)
		t.Run("user with create access should not be able to reload app", testCannotReloadApp)
	}()

	func() {
		// Connect to engine with reload access.
		testState.global, err = connect(reloadEngineURL, reloadClaims)
		require.NoError(t, err)
		defer testState.global.DisconnectFromServer()

		// Run test cases for user with reload access.
		t.Run("user with reload access should not be able to create app", testCannotCreateApp)
		t.Run("user with reload access should be able to open app", testCanOpenApp)
		t.Run("user with reload access should be able to reload app", testCanReloadApp)
		t.Run("user with reload access should be able to save app", testCanSaveApp)
	}()

	func() {
		// Connect to engine with view access.
		testState.global, err = connect(reloadEngineURL, viewClaims)
		require.NoError(t, err)
		defer testState.global.DisconnectFromServer()

		// Run test cases for user with view access.
		t.Run("user with view access should not be able to create app", testCannotCreateApp)
		t.Run("user with view access should be able to open app", testCanOpenApp)
		t.Run("user with view access should be able to read movies object", testCanReadMoviesObject)
		t.Run("user with view access should be able to read movies data", testCanReadMoviesData)
		t.Run("user with view access should not be able to create sheet", testCannotCreateSheet)
		t.Run("user with view access should not be able to reload app", testCannotReloadApp)
	}()
}

// testCanCreateApp verifies that the current user can create an app.
func testCanCreateApp(t *testing.T) {
	testState.app, err = createApp(testState.global, appName)
	require.NoError(t, err)
	testState.appID = testState.app.GenericId
}

// testCannotCreateApp verifies that the current user is not allowed to create an app.
func testCannotCreateApp(t *testing.T) {
	_, err = createApp(testState.global, "CannotCreateThisApp")
	require.Error(t, err)
	require.Contains(t, strings.ToLower(err.Error()), "access denied")
}

// testCanOpenApp verifies that the current user can open an app.
func testCanOpenApp(t *testing.T) {
	testState.app, err = openApp(testState.global, testState.appID)
	require.NoError(t, err)
}

// testCanSaveApp verifies that the current user can save an app.
func testCanSaveApp(t *testing.T) {
	err = testState.app.DoSave(ctx, "")
	require.NoError(t, err)
}

// testCanCreateSheet verifies that the current user can create a sheet object.
func testCanCreateSheet(t *testing.T) {
	testState.sheet, err = createSheet(testState.app, sheetID)
	require.NoError(t, err)

	err = testState.app.SaveObjects(ctx)
	require.NoError(t, err)
}

// testCanCreateSheet verifies that the current user is not allowed to create a sheet object.
func testCannotCreateSheet(t *testing.T) {
	_, err = createSheet(testState.app, sheetID)
	require.Error(t, err)
	require.Contains(t, strings.ToLower(err.Error()), "access denied")
}

// testCanReadSheet verifies that the current user can read a sheet object.
func testCanReadSheet(t *testing.T) {
	testState.sheet, err = readObject(testState.app, sheetID)
	require.NoError(t, err)
}

// testCanCreateMoviesObject verifies that the current user can create a custom movies object.
func testCanCreateMoviesObject(t *testing.T) {
	testState.movies, err = createMoviesObject(testState.app, moviesID)
	require.NoError(t, err)
}

// testCanReadMoviesObject verifies that the current user can read a custom movies object.
func testCanReadMoviesObject(t *testing.T) {
	testState.movies, err = readObject(testState.app, moviesID)
	require.NoError(t, err)
}

// testCanReadMoviesData verifies that the current user can read data in a custom movies object.
func testCanReadMoviesData(t *testing.T) {
	var (
		titles []string
		years  []float64
	)
	titles, years, err = readMoviesData(testState.movies)
	require.NoError(t, err)
	require.Equal(t, "Armageddon", titles[1])
	require.Equal(t, float64(1998), years[1])
}

// testCanReloadApp verifies that the current user can reload data into an app.
func testCanReloadApp(t *testing.T) {
	err = reloadMoviesData(testState.app)
	require.NoError(t, err)
}

// testCannotReloadApp verifies that the current user is not allowed to reload data into an app.
func testCannotReloadApp(t *testing.T) {
	err = reloadMoviesData(testState.app)
	require.Error(t, err)
	require.Contains(t, strings.ToLower(err.Error()), "access denied")
}

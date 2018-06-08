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
	engineURL              = "ws://localhost:9176"
	countryListBoxObjectID = "eVTNPmv"
)

func TestAdminCanCreateAppAndSheet(t *testing.T) {
	const (
		sheetID = "SHEET01"
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
		layout *enigma.GenericObjectLayout
		info   *enigma.NxInfo
		err    error
		id     string
	)

	func() {
		global, err = connect(engineURL, adminClaims)
		require.NoError(t, err)
		defer global.DisconnectFromServer()

		// Admin should be able to create an app.
		id, err = createApp(global, t.Name())
		require.NoError(t, err)

		// Admin should be able to open an app and create a sheet.
		app, err = openApp(global, id)
		require.NoError(t, err)
		sheet, err = createSheet(app, sheetID)
		require.NoError(t, err)
		info, err = sheet.GetInfo(ctx)
		require.NoError(t, err)
		assert.Equal(t, sheetID, info.Id)
	}()

	func() {
		global, err = connect(engineURL, nonAdminClaims)
		require.NoError(t, err)
		defer global.DisconnectFromServer()

		// Non-admin should not be able to create an app.
		_, err = createApp(global, "CannotCreateThisApp")
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "access denied")

		// Non-admin should be able to open an app.
		app, err = openApp(global, id)
		require.NoError(t, err)

		// Non-admin should not be able to create a sheet.
		_, err = createSheet(app, sheetID)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "access denied")

		// Non-admin should be able to read a sheet.
		sheet, layout, err = readSheetLayout(app, sheetID)
		require.NoError(t, err)
	}()
}

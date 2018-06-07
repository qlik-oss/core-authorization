package access

import (
	"testing"

	"github.com/dgrijalva/jwt-go"

	"github.com/stretchr/testify/assert"
)

const (
	engineURL              = "ws://localhost:9176"
	countryListBoxObjectID = "eVTNPmv"
)

func TestThatAccessIsDeniedWithoutJWT(t *testing.T) {
	err := accessEngineAppObject(engineURL, "Shared-Africa-Urbanization.qvf", "eVTNPmv", nil)
	assert.Error(t, err)
}

func TestThatAnyoneCanReadAppObject(t *testing.T) {
	jwtClaims := jwt.MapClaims{
		"sub": "someUser",
	}
	err := accessEngineAppObject(engineURL, "Shared-Africa-Urbanization.qvf", "eVTNPmv", jwtClaims)
	assert.NoError(t, err)
}

// func TestThatAdminCanXXX(t *testing.T) {
// 	err := accessEngineAppObjectWithJwt("ws://localhost:9071", "Shared-Africa-Urbanization.qvf", "eVTNPmv", nil)
// }

// func TestThatTheCountryListBoxIsBlockedWithoutJwt(t *testing.T) {
// 	err := accessEngineAppObjectWithJwt("ws://localhost:9071", "Shared-Africa-Urbanization.qvf", "eVTNPmv", nil)
// 	assert.Error(t, err)
// }

// func TestThatTheCountryListBoxIsAccessibleInEngine1WithJwt(t *testing.T) {

// 	jwtClaims := jwt.MapClaims{
// 		"sub":    "user1",
// 		"groups": []string{"A", "B"},
// 	}
// 	err := accessEngineAppObjectWithJwt("ws://localhost:9071", "Shared-Africa-Urbanization.qvf", "eVTNPmv", jwtClaims)
// 	assert.Nil(t, err)
// }

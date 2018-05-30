package access_tests

import (
	"testing"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)


func TestThatTheCountryListBoxIsBlockedWithoutJwt(t *testing.T) {
	err := accessEngineAppObjectWithJwt("ws://localhost:9071", "Shared-Africa-Urbanization.qvf", "eVTNPmv", nil)
	assert.Error(t, err)
}

func TestThatTheCountryListBoxIsAccessibleInEngine1WithJwt(t *testing.T) {

	jwtClaims := jwt.MapClaims{
		"sub":    "user1",
		"groups": []string{"A", "B"},
	}
	err := accessEngineAppObjectWithJwt("ws://localhost:9071", "Shared-Africa-Urbanization.qvf", "eVTNPmv", jwtClaims)
	assert.Nil(t, err)
}

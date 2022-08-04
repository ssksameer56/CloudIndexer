package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAccessToken(t *testing.T) {
	err := LoadConfig()
	require.NoError(t, err)
	require.NotEmpty(t, Config.AccessToken)

}

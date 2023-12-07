package authentication

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAuthToken(t *testing.T) {
	app := App()

	req, err := http.NewRequest(http.MethodGet, "/token", nil)
	require.Nil(t, err)

	res, err := app.Test(req)
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	var resData map[string]interface{}
	err = json.Unmarshal(resBody, &resData)
	require.Nil(t, err)

	assert.Equal(t, 1, len(resData))

	token, ok := resData["token"]
	assert.True(t, ok)
	assert.NotEmpty(t, token)
}

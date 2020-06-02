package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckAccess(t *testing.T) {
	adminHash := "ZehL4zUy-3hMSBKWdfnv86aCsnFowOp0Syz1juAjN8U="
	password := "qwerty"

	sauth := NewSimpleAuth(adminHash)

	res, err := sauth.CheckAccess(password)
	require.NoError(t, err)
	assert.Equal(t, true, res)
}

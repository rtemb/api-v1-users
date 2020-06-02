package auth

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/pkg/errors"
)

type SimpleAuth struct {
	adminHash string
}

func NewSimpleAuth(adminHash string) *SimpleAuth {
	return &SimpleAuth{adminHash: adminHash}
}

// CheckAccess check access to admin functions
func (s *SimpleAuth) CheckAccess(password string) (bool, error) {
	bv := []byte(password)
	hasher := sha256.New()

	_, err := hasher.Write(bv)
	if err != nil {
		return false, errors.Wrap(err, "Unable to prepare Hasher")
	}

	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return sha == s.adminHash, nil
}

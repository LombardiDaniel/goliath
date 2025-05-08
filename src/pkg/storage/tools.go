package storage

import (
	"net/url"
	"path"

	"github.com/LombardiDaniel/goliath/src/pkg/constants"
)

type storageDir string

const (
	UserAvatars storageDir = "user-avatars"
)

func GetFullObjUrl(objPath string) (string, error) {
	return url.JoinPath(constants.S3Endpoint, constants.S3Bucket, objPath)
}

func GetPublicPath(p storageDir, filename string) string {
	return path.Join("public", string(p), filename)
}

func GetPrivatePath(p storageDir, filename string) string {
	return path.Join("private", string(p), filename)
}

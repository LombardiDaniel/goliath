package storage

import (
	"net/url"
	"path"

	"github.com/LombardiDaniel/gopherbase/common"
)

type storageDir string

const (
	USER_AVATARS storageDir = "user-avatars"
)

func GetFullObjUrl(objPath string) (string, error) {
	return url.JoinPath(common.S3Endpoint, common.S3Bucket, objPath)
}

func GetPublicPath(p storageDir, filename string) string {
	return path.Join("public", string(p), filename)
}

func GetPrivatePath(p storageDir, filename string) string {
	return path.Join("private", string(p), filename)
}

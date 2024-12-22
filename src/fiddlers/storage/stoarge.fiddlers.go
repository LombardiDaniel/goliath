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
	return url.JoinPath(common.S3_ENDPOINT, common.S3_BUCKET, objPath)
}

func GetPublicPath(p storageDir, filename string) string {
	return path.Join("public", string(p), filename)
}

func GetPrivatePath(p storageDir, filename string) string {
	return path.Join("private", string(p), filename)
}

package storage

import (
	"net/url"
	"path"

	"github.com/LombardiDaniel/gopherbase/common"
)

type storagePrefix string

const (
	USER_AVATARS storagePrefix = "user-avatars"
)

func GetFullObjUrl(objPath string) (string, error) {
	prefix := "https://"
	if !common.S3_SECURE {
		prefix = "http://"
	}
	return url.JoinPath(prefix+common.S3_ENDPOINT, common.S3_BUCKET, objPath)
}

func GetPublicPath(p storagePrefix, filename string) string {
	return path.Join("public", string(p), filename)
}

func GetPrivatePath(p storagePrefix, filename string) string {
	return path.Join("private", string(p), filename)
}

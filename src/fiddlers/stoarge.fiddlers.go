package fiddlers

import (
	"net/url"

	"github.com/LombardiDaniel/gopherbase/common"
)

func GetFullObjStorageUrl(objPath string) (string, error) {
	prefix := "https://"
	if !common.S3_SECURE {
		prefix = "http://"
	}
	return url.JoinPath(prefix+common.S3_ENDPOINT, common.DEFAULT_BUCKET, objPath)
}

package constants

import (
	"time"

	"github.com/LombardiDaniel/gopherbase/src/pkg/common"
)

const (
	TimestampStrFormat       string = time.RFC3339 // "yyyy-mm-ddThh:mm:ssZhh:mm" and "2006-01-02T15:04:05-07:00"
	DefaultTimzone           string = "GMT-3"
	GinCtxJwtClaimKeyName    string = "jwtClaims"
	JwtTimeoutSecs           int    = 30 * 60
	OptLen                   int    = 128
	OrgInviteTimeoutDays     int    = 15
	PasswordResetTimeoutDays int    = 1
	MaxRequestSize           int64  = 5 * 1024 * 1024 // 5MB default
)

var (
	ProjectName                       string = common.GetEnvVarDefault("PROJECT_NAME", "gopherbase")
	NoreplyEmail                      string = common.GetEnvVarDefault("NO_REPLY_EMAIL", "no-reply@example.com")
	AppHostUrl                        string = common.GetEnvVarDefault("APP_HOST_URL", "http://127.0.0.1:8080/")
	ApiHostUrl                        string = common.GetEnvVarDefault("API_HOST_URL", "http://127.0.0.1:8080/")
	JwtCookieName                     string = ProjectName + "_jwt"
	PasswordResetTimeoutJwtCookieName string = ProjectName + "_pwreset_jwt"
	S3Endpoint                        string = common.GetEnvVarDefault("S3_ENDPOINT", "https://br-se1.magaluobjects.com")
	S3Region                          string = common.GetEnvVarDefault("S3_REGION", "br-se1")
	S3Bucket                          string = common.GetEnvVarDefault("S3_BUCKET", ProjectName+"-gopherbase")
)

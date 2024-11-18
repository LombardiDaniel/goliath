package common

import (
	"time"
)

const (
	PROJECT_NAME string = "PROJECT_NAME"
	// TIMESTAMP_STR_FORMAT string = "yyyy-mm-ddThh:mm:ss"
	TIMESTAMP_STR_FORMAT       string = time.RFC3339
	COOKIE_NAME                string = PROJECT_NAME + "_JWT"
	GIN_CTX_JWT_CLAIM_KEY_NAME string = "jwtClaims"
	JWT_TIMEOUT_SECS           int64  = 30 * 60
	NOREPLY_EMAIL              string = "noreply@patos.dev"
)

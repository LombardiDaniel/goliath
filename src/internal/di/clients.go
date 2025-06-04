package di

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/LombardiDaniel/goliath/src/pkg/common"
	"github.com/LombardiDaniel/goliath/src/pkg/constants"
	"github.com/LombardiDaniel/goliath/src/pkg/it"
	"github.com/LombardiDaniel/goliath/src/pkg/oauth"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/resendlabs/resend-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type clients struct {
	MongoClient     *mongo.Client
	MetricsMongoCol *mongo.Collection
	EventsMongoCol  *mongo.Collection

	PostgressConn *sql.DB

	MinioClient *minio.Client

	OauthConfigMap map[string]oauth.Provider

	ResendClient *resend.Client
}

func newClients() *clients {
	ctx := context.Background()

	// PSQL
	pgConnStr := common.GetEnvVarDefault("POSTGRES_URI", "postgres://user:password@localhost:5432/db?sslmode=disable")
	db := it.Must(sql.Open("postgres", pgConnStr))
	if err := db.Ping(); err != nil {
		panic(errors.Join(err, errors.New("could not ping pgsql")))
	}
	pgIdleConns := it.Must(strconv.Atoi(common.GetEnvVarDefault("POSTGRES_IDLE_CONNS", "2")))
	pgOpenConns := it.Must(strconv.Atoi(common.GetEnvVarDefault("POSTGRES_OPEN_CONNS", "10")))
	db.SetMaxIdleConns(pgIdleConns)
	db.SetMaxOpenConns(pgOpenConns)
	it.Must(db.Exec(fmt.Sprintf("SET TIME ZONE '%s';", constants.DefaultTimzone)))

	// Mongo
	mongoConn := options.Client().ApplyURI(
		common.GetEnvVarDefault("MONGO_URI", "mongodb://localhost:27017"),
	)
	mongoClient := it.Must(mongo.Connect(ctx, mongoConn))
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(errors.Join(err, errors.New("could not ping mongodb")))
	}
	tsIdxModel := mongo.IndexModel{
		Keys:    bson.M{"ts": 1},
		Options: options.Index(),
	}
	metricsCol := mongoClient.Database("telemetry").Collection("metrics")
	eventsCol := mongoClient.Database("telemetry").Collection("events")
	it.Must(metricsCol.Indexes().CreateOne(ctx, tsIdxModel))
	it.Must(eventsCol.Indexes().CreateOne(ctx, tsIdxModel))

	// Minio
	s3Host := it.Must(common.ExtractHostFromUrl(constants.S3Endpoint))
	s3Secure := it.Must(common.UrlIsSecure(constants.S3Endpoint))
	minioClient := it.Must(minio.New(
		s3Host,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				os.Getenv("S3_ACCESS_KEY_ID"),
				os.Getenv("S3_SECRET_ACCESS_KEY"),
				"",
			),
			Region: constants.S3Region,
			Secure: s3Secure,
		},
	))

	// OAuth
	oauthBaseCallback := constants.ApiHostUrl + "v1/auth/%s/callback"
	oauthConfigMap := make(map[string]oauth.Provider)
	oauthConfigMap[oauth.GOOGLE_PROVIDER] = oauth.NewGoogleProvider(&oauth2.Config{
		RedirectURL:  fmt.Sprintf(oauthBaseCallback, oauth.GOOGLE_PROVIDER),
		ClientID:     os.Getenv("OAUTH_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_GOOGLE_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	})
	oauthConfigMap[oauth.GITHUB_PROVIDER] = oauth.NewGithubProvider(&oauth2.Config{
		RedirectURL:  fmt.Sprintf(oauthBaseCallback, oauth.GITHUB_PROVIDER),
		ClientID:     os.Getenv("OAUTH_GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
		Scopes: []string{
			"read:user",
		},
		Endpoint: github.Endpoint,
	})

	resendClient := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	return &clients{
		MongoClient:     mongoClient,
		MetricsMongoCol: metricsCol,
		EventsMongoCol:  eventsCol,
		PostgressConn:   db,
		MinioClient:     minioClient,
		OauthConfigMap:  oauthConfigMap,
		ResendClient:    resendClient,
	}
}

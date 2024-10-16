package firestore

import (
	"context"
	"image_filter_server/pkg/errors"
	"image_filter_server/pkg/logging"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var firestoreClient *firestore.Client
var firestoreOnce sync.Once

// Code to initialise firestore
func InitialiaseFirestore() *firestore.Client {
	ctx := context.Background()

	firestoreOnce.Do(func() {
		projectID := viper.GetString("FIRESTORE_PROJECT_ID")
		if projectID == "" {
			logging.Error(&gin.Context{}, errors.ErrFirestoreProjectIDMissing.Error(), zap.Error(errors.ErrFirestoreProjectIDMissing))
			panic(errors.ErrFirestoreProjectIDMissing)
		}
		// Set up a new client
		client, err := firestore.NewClient(ctx, projectID)
		if err != nil {
			logging.Error(&gin.Context{}, errors.ErrFirestoreConnectionFailed.Error(), zap.Error(err))
			panic(err)
		}
		firestoreClient = client
	})
	return firestoreClient
}

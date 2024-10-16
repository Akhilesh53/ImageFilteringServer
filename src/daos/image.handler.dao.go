package daos

import (
	"sync"

	firestoreDB "image_filter_server/pkg/firestore"
	"image_filter_server/src/models"
	"image_filter_server/src/utils"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var imageHandlerDao *ImageHandlerDao
var imageHandlerDaoOnce sync.Once

type ImageHandlerDao struct {
	fireWorkClient *firestore.Client
}

// GetImageHandlerDao is a function that returns a singleton instance of ImageHandlerDao
func GetImageHandlerDao() *ImageHandlerDao {
	imageHandlerDaoOnce.Do(func() {
		imageHandlerDao = &ImageHandlerDao{
			fireWorkClient: firestoreDB.InitialiaseFirestore(),
		}
	})
	return imageHandlerDao
}

// get firestore client
func (dao *ImageHandlerDao) GetFirestoreClient() *firestore.Client {
	return dao.fireWorkClient
}

// func to check whether image url is present in the collection or not
func (dao *ImageHandlerDao) IsDocPresent(ctx *gin.Context, imageURL string) (bool, error) {
	doc, err := dao.GetFirestoreClient().Collection(viper.GetString("FIRESTORE_COLLECTION_NAME")).Doc(utils.ModifyURL(imageURL)).Get(ctx)
	if err != nil {
		if !doc.Exists() {
			return false, nil
		}
		return false, errors.WithStack(errors.WithMessage(err, " : error while fetching image url response from collection"))
	}

	return doc.Exists(), nil
}

// func to get image response from collection
func (dao *ImageHandlerDao) GetDocResponse(ctx *gin.Context, imageURL string) (*models.FirebaseCollectionResult, error) {
	doc, err := dao.GetFirestoreClient().Collection(viper.GetString("FIRESTORE_COLLECTION_NAME")).Doc(utils.ModifyURL(imageURL)).Get(ctx)
	if err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, " error while fetching image url response from collection"))
	}

	if doc.Exists() {
		var result models.FirebaseCollectionResult
		err := doc.DataTo(&result)
		if err != nil {
			return nil, errors.WithStack(errors.WithMessage(err, " error while converting firestore data to struct"))
		}
		return &result, nil
	}

	return nil, errors.WithStack(errors.New("image url not found in collection"))
}

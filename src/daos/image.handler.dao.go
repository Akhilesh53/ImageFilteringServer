package daos

import (
	"sync"

	firestoreDB "image_filter_server/pkg/firestore"

	"cloud.google.com/go/firestore"
)

var imageHandlerDao *ImageHandlerDao
var imageHandlerDaoOnce sync.Once

type ImageHandlerDao struct {
	fireWorkCliet *firestore.Client
}

// GetImageHandlerDao is a function that returns a singleton instance of ImageHandlerDao
func GetImageHandlerDao() *ImageHandlerDao {
	imageHandlerDaoOnce.Do(func() {
		imageHandlerDao = &ImageHandlerDao{
			fireWorkCliet: firestoreDB.InitialiaseFirestore(),
		}
	})
	return imageHandlerDao
}

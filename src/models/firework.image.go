package models

/*
Fields keys and values are:

	String image_url < image url >
	Map googleVisionResult < Json here >
	String gemini1.5Flash8bFlashResult < string response here >
	String conclusion < “open” or “block” >
*/
type FirebaseCollectionResult struct {
	ImageURL                   string                 `firestore:"image_url"`
	GoogleVisionResult         map[string]interface{} `firestore:"googleVisionResult"`
	Gemini15Flash8bFlashResult string                 `firestore:"gemini1.5Flash8bFlashResult"`
	Conclusion                 string                 `firestore:"conclusion"`
}

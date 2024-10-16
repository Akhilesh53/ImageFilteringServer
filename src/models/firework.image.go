package models

/*
Fields keys and values are:

	String image_url < image url >
	Map googleVisionResult < Json here >
	String gemini1.5Flash8bFlashResult < string response here >
	String conclusion < “open” or “block” >
*/
type FirebaseCollectionResult struct {
	ImageURL                   string         `firestore:"image_url"`
	GoogleVisionResult         *GCPSafeSearch `firestore:"googleVisionResult"`
	Gemini15Flash8bFlashResult string         `firestore:"gemini1.5Flash8bFlashResult"`
	Conclusion                 string         `firestore:"conclusion"`
}

// func to create a new instance of FirebaseCollectionResult
func NewFirebaseCollectionResult(imageURL string, googleVisionResult *GCPSafeSearch, gemini15Flash8bFlashResult string, conclusion string) *FirebaseCollectionResult {
	return &FirebaseCollectionResult{
		ImageURL:                   imageURL,
		GoogleVisionResult:         googleVisionResult,
		Gemini15Flash8bFlashResult: gemini15Flash8bFlashResult,
		Conclusion:                 conclusion,
	}
}

// default constructor
func NewDefaultFirebaseCollectionResult() *FirebaseCollectionResult {
	return &FirebaseCollectionResult{}
}

// func to set the image url
func (f *FirebaseCollectionResult) SetImageURL(imageURL string) *FirebaseCollectionResult {
	f.ImageURL = imageURL
	return f
}

// func to get the image url
func (f *FirebaseCollectionResult) GetImageURL() string {
	return f.ImageURL
}

// func to set the google vision result
func (f *FirebaseCollectionResult) SetGoogleVisionResult(googleVisionResult *GCPSafeSearch) *FirebaseCollectionResult {
	f.GoogleVisionResult = googleVisionResult
	return f
}

// func to get the google vision result
func (f *FirebaseCollectionResult) GetGoogleVisionResult() *GCPSafeSearch {
	return f.GoogleVisionResult
}

// func to set the gemini 1.5 flash 8b flash result
func (f *FirebaseCollectionResult) SetGemini15Flash8bFlashResult(gemini15Flash8bFlashResult string) *FirebaseCollectionResult {
	f.Gemini15Flash8bFlashResult = gemini15Flash8bFlashResult
	return f
}

// func to get the gemini 1.5 flash 8b flash result
func (f *FirebaseCollectionResult) GetGemini15Flash8bFlashResult() string {
	return f.Gemini15Flash8bFlashResult
}

// func to set the conclusion
func (f *FirebaseCollectionResult) SetConclusion(conclusion string) *FirebaseCollectionResult {
	f.Conclusion = conclusion
	return f
}

// func to get the conclusion
func (f *FirebaseCollectionResult) GetConclusion() string {
	return f.Conclusion
}

// func to get the firestore map of FirebaseCollectionResult

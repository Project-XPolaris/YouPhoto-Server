package task

type Signal struct {
}

const (
	TypeScanLibrary         = "ScanLibrary"
	TypeRemove              = "RemoveLibrary"
	TypeRemoveNotExistImage = "RemoveNotExistImage"
	TypeScanImageFile       = "ScanImageFile"
	TypeCreateImage         = "CreateImage"
	TypeGenerateThumbnail   = "GenerateThumbnail"
	TypeReadImage           = "ReadImage"
	TypeImageClassify       = "ImageClassify"
	TypeNSFWCheck           = "NSFWCheck"
	TypeDeepdanbooru        = "Deepdanbooru"
	TypePreprocess          = "Preprocess"
	TypeLoraTrain           = "LoraTrain"
	TypeTagger              = "Tagger"
)

//const (
//	TaskStatusRunning = "Running"
//	TaskStatusDone    = "Done"
//	TaskStatusError   = "Error"
//	TaskStatusStop    = "Stop"
//)

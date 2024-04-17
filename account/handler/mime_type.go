package handler

var validImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png": true,
}


func isAllowedImageType(mimeType string) bool{
	_, exists := validImageTypes[mimeType]

	return exists
}
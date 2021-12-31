package filter

var defaultPage = 1
var defaultSize = 20

type Filter struct {
	Page int               `json:"page"`
	Size int               `json:"size"`
	Sort map[string]string `json:"sort"`
}

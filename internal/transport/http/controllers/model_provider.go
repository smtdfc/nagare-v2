package controllers

type ModelProviderController struct{}

// @Injectable
func NewModelProviderController() *ModelProviderController {
	return &ModelProviderController{}
}

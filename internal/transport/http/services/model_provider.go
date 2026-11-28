package services

type ModelProviderService struct{}

// @Injectable
func NewModelProviderService() *ModelProviderService {
	return &ModelProviderService{}
}

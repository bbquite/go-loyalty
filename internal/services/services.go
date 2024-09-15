package services

type StorageRepo interface {
	RegisterUser()
}

type AppService struct {
	store StorageRepo
}

func NewAppService(store StorageRepo) *AppService {
	return &AppService{
		store: store,
	}
}

func (h *AppService) RegisterUser() error {
	h.store.RegisterUser()
	return nil
}

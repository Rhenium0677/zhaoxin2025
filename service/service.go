package service

type Service struct {
	Admin
}

func New() *Service {
	service := &Service{}
	return service
}

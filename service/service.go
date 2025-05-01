package service

type Service struct {
	Admin
	Stu
	Que
}

func New() *Service {
	service := &Service{}
	return service
}

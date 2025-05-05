package service

type Service struct {
	Admin
	Stu
	Interv
	Que
}

func New() *Service {
	service := &Service{}
	return service
}

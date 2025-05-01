package controller

type Controller struct {
	Admin
	Que
	Stu
}

func New() *Controller {
	Controller := &Controller{}
	return Controller
}

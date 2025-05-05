package controller

type Controller struct {
	Admin
	Stu
	Interv
	Que
}

func New() *Controller {
	Controller := &Controller{}
	return Controller
}

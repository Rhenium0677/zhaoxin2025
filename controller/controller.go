package controller

type Controller struct {
	Admin
}

func New() *Controller {
	Controller := &Controller{}
	return Controller
}

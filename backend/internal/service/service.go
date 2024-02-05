package service

type Service interface {
	AuthenticateUser(email, password string) (string, error)
}

type serviceImpl struct {
	// TODO: add db interface here

}

func New() Service {
	return &serviceImpl{}
}

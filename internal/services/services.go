package services

//go:generate mockery --name ServiceInterface --output ./mocks --filename services.go
type ServiceInterface interface {
}

type Service struct {
}

func NewService() (*Service, error) {
	return &Service{}, nil
}

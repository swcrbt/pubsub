package service

type Service interface {
	Run () error
	GetName() string
}
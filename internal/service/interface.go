package service

type IURLService interface {
	Shorten(string) (string, error)
	Resolve(string) (string, error)
}

package repository

type IURLRepository interface {
	Save(string, string) error
	Get(string) (string, error)
}

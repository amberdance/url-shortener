package storage

type Storage interface {
	Save(shortID, url string) error
	Get(shortID string) (string, bool)
}

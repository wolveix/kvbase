package kvbase

type Backend interface {
	Count(bucket string) (int, error)
	Create(bucket string, key string, model interface{}) error
	Delete(bucket string, key string) error
	Get(bucket string, model interface{}) (*map[string]interface{}, error)
	Read(bucket string, key string, model interface{}) error
	Update(bucket string, key string, model interface{}) error
}

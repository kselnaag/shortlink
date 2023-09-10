package i7e

type ICtrlHTTP interface {
	AllPairs() (string, error)
	Long(body []byte) (string, error)
	Short(body []byte) (string, error)
	Save(body []byte) (string, error)
	Hash(hash string) (string, error)
}

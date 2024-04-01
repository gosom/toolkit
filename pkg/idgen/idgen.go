package idgen

import "github.com/sqids/sqids-go"

var (
	defaultProvider = New()
)

type IDProvider interface {
	Encode(int64) string
	Decode(string) (int64, error)
}

func New() IDProvider {
	const length = 10

	s, err := sqids.New(sqids.Options{
		MinLength: length,
	})

	if err != nil {
		panic(err)
	}

	ans := sqidProvider{
		s: s,
	}

	return &ans
}

func Encode(id int64) string {
	return defaultProvider.Encode(id)
}

func Decode(s string) (int64, error) {
	return defaultProvider.Decode(s)
}

type sqidProvider struct {
	s *sqids.Sqids
}

func (p *sqidProvider) Encode(id int64) string {
	ans, _ := p.s.Encode([]uint64{uint64(id)})

	return ans
}

func (p *sqidProvider) Decode(s string) (int64, error) {
	ans := p.s.Decode(s)

	if len(ans) == 0 {
		return 0, nil
	}

	return int64(ans[0]), nil
}

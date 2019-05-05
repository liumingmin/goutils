package cbk

var Impls = make(map[string]CircuitBreaker)

const (
	SIMPLE = "simple"
)

type CircuitBreaker interface {
	Check(key string) error
	Succeed(key string)
	Failed(key string)
}

type Error struct {
	error
	Msg string
}

func (c Error) Error() string {
	if c.Msg != "" {
		return c.Msg
	}
	return "CircuitBreaker is break"
}

func init() {
	simple := &CircuitBreakerSimple{}
	simple.Init()
	Impls[SIMPLE] = simple
}

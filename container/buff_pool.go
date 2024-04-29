package container

type Buffer128K []byte
type Buffer4M []byte

var PoolBuffer128K = NewPool(func() Buffer128K {
	return Buffer128K(make([]byte, 128*1024))
})

var PoolBuffer4M = NewPool(func() Buffer4M {
	return Buffer4M(make([]byte, 4*1024*1024))
})

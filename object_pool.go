package main

// Currently unused

const objectBlockSize = 1 << 10

type ObjectPool struct {
	_blocks [][objectBlockSize]interface{}
	_size   sizeType
}

func NewObjectPool() *ObjectPool {
	return &ObjectPool{
		_blocks: [][objectBlockSize]interface{}{},
		_size:   0,
	}
}

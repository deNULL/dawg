package main

const bitBlockSize = 1 << 10

type BitPool struct {
	words []uint64
	size  sizeType
}

func NewBitPool() *BitPool {
	return &BitPool{
		words: []uint64{},
		size:  0,
	}
}

func (pool *BitPool) clear() {
	pool.words = pool.words[:0]
	pool.size = 0
}

func (pool *BitPool) allocate() {
	if pool.size%64 == 0 {
		pool.words = append(pool.words, 0)
	}
	pool.size += 1
}

func (pool *BitPool) get(index baseType) bool {
	return pool.words[index>>6]&(1<<uint8(index&63)) != 0
}

func (pool *BitPool) set(index baseType, bit bool) {
	if bit {
		pool.words[index>>6] |= 1 << uint8(index&63)
		return
	}
	pool.words[index>>6] &^= 1 << uint8(index&63)
}

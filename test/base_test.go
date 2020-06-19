package test

import "testing"

func TestT(t *testing.T) {
	t.Log("aa")
}

func TestAssign(t *testing.T) {
	var (
		a int = 1
		b     = 2
	)
	a, b = b, a // 交换
	t.Log(a, b)
	const (
		c = iota + 1
		d
		e
		f
	)
	const g = 5
	t.Log(c, d, e, f, g)
	const (
		h = 1 << iota
		i
		j
		k = 3 << iota
		l
		m
	)
	t.Log(h, i, j, k, l, m)
	const (
		n = 2
		o
		p
	)
	// 不可用连续赋值
	var (
		q = 3
		r int
		s uint32
	)
	t.Log(n, o, p, q, r, s)
	const (
		u string = "aa"
		v
		w = 2
		x
	)
	t.Log(u, v, w, x)
}

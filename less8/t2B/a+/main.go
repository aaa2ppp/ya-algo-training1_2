package main

import (
	"bufio"
	"cmp"
	"io"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type TreeNode[T cmp.Ordered] struct {
	Left  *TreeNode[T]
	Right *TreeNode[T]
	Value T
}

// Add добавляет узел в дерево и возвращает true.
// Если узел уже существовал, возвращает false.
func (p *TreeNode[T]) Add(v T) bool {
	if v < p.Value {
		if p.Left == nil {
			p.Left = &TreeNode[T]{Value: v}
			return true
		}
		return p.Left.Add(v)
	} else if v > p.Value {
		if p.Right == nil {
			p.Right = &TreeNode[T]{Value: v}
			return true
		}
		return p.Right.Add(v)
	}
	return false
}

// Search проверяет наличе узла в дереве.
func (p *TreeNode[T]) Search(v T) bool {
	if v < p.Value {
		if p.Left == nil {
			return false
		}
		return p.Left.Search(v)
	} else if v > p.Value {
		if p.Right == nil {
			return false
		}
		return p.Right.Search(v)
	}
	return true
}

func (p *TreeNode[T]) LNR(h int, f func(h int, v T)) {
	if p.Left != nil {
		p.Left.LNR(h+1, f)
	}
	f(h, p.Value)
	if p.Right != nil {
		p.Right.LNR(h+1, f)
	}
}

type Tree[T cmp.Ordered] struct {
	root *TreeNode[T]
}

func (t *Tree[T]) Add(v T) bool {
	if t.root == nil {
		t.root = &TreeNode[T]{Value: v}
		return true
	}
	return t.root.Add(v)
}

func (t *Tree[T]) Search(v T) bool {
	if t.root == nil {
		return false
	}
	return t.root.Search(v)
}

func (t *Tree[T]) LNR(f func(h int, v T)) {
	if t.root != nil {
		t.root.LNR(0, f)
	}
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var tree Tree[int]

	for {
		op, err := scanWord(sc)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		switch op {
		case "ADD":
			v, err := scanInt(sc)
			if err != nil {
				panic(err)
			}
			if tree.Add(v) {
				bw.WriteString("DONE\n")
			} else {
				bw.WriteString("ALREADY\n")
			}
		case "SEARCH":
			v, err := scanInt(sc)
			if err != nil {
				panic(err)
			}
			if tree.Search(v) {
				bw.WriteString("YES\n")
			} else {
				bw.WriteString("NO\n")
			}
		case "PRINTTREE":
			tree.LNR(func(h int, v int) {
				bw.WriteString(strings.Repeat(".", h))
				writeInt(bw, v, writeOpts{end: '\n'})
			})
		}
	}
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout)
}

// ----------------------------------------------------------------------------

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func scanWord(sc *bufio.Scanner) (string, error) {
	if sc.Scan() {
		return sc.Text(), nil
	}
	return "", io.EOF
}

func scanInt(sc *bufio.Scanner) (int, error)                  { return scanIntX[int](sc) }
func scanTwoInt(sc *bufio.Scanner) (_, _ int, _ error)        { return scanTwoIntX[int](sc) }
func scanThreeInt(sc *bufio.Scanner) (_, _, _ int, _ error)   { return scanThreeIntX[int](sc) }
func scanFourInt(sc *bufio.Scanner) (_, _, _, _ int, _ error) { return scanFourIntX[int](sc) }

func scanIntX[T Int](sc *bufio.Scanner) (res T, err error) {
	sc.Scan()
	v, err := strconv.ParseInt(unsafeString(sc.Bytes()), 0, int(unsafe.Sizeof(res))<<3)
	return T(v), err
}

func scanTwoIntX[T Int](sc *bufio.Scanner) (v1, v2 T, err error) {
	v1, err = scanIntX[T](sc)
	if err == nil {
		v2, err = scanIntX[T](sc)
	}
	return v1, v2, err
}

func scanThreeIntX[T Int](sc *bufio.Scanner) (v1, v2, v3 T, err error) {
	v1, err = scanIntX[T](sc)
	if err == nil {
		v2, err = scanIntX[T](sc)
	}
	if err == nil {
		v3, err = scanIntX[T](sc)
	}
	return v1, v2, v3, err
}

func scanFourIntX[T Int](sc *bufio.Scanner) (v1, v2, v3, v4 T, err error) {
	v1, err = scanIntX[T](sc)
	if err == nil {
		v2, err = scanIntX[T](sc)
	}
	if err == nil {
		v3, err = scanIntX[T](sc)
	}
	if err == nil {
		v4, err = scanIntX[T](sc)
	}
	return v1, v2, v3, v4, err
}

func scanInts[T Int](sc *bufio.Scanner, a []T) error {
	for i := range a {
		v, err := scanIntX[T](sc)
		if err != nil {
			return err
		}
		a[i] = v
	}
	return nil
}

type Int interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

type Number interface {
	Int | ~float32 | ~float64
}

type writeOpts struct {
	sep   byte
	begin byte
	end   byte
}

func defaultWriteOpts() writeOpts {
	return writeOpts{sep: ' ', end: '\n'}
}

func writeInt[I Int](bw *bufio.Writer, v I, opts writeOpts) error {
	var buf [32]byte

	var err error
	if opts.begin != 0 {
		err = bw.WriteByte(opts.begin)
	}

	if err == nil {
		_, err = bw.Write(strconv.AppendInt(buf[:0], int64(v), 10))
	}

	if err == nil && opts.end != 0 {
		err = bw.WriteByte(opts.end)
	}

	return err
}

func writeInts[I Int](bw *bufio.Writer, a []I, opts writeOpts) error {
	var err error
	if opts.begin != 0 {
		err = bw.WriteByte(opts.begin)
	}

	if len(a) != 0 {
		var buf [32]byte

		if opts.sep == 0 {
			opts.sep = ' '
		}

		_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[0]), 10))

		for i := 1; err == nil && i < len(a); i++ {
			err = bw.WriteByte(opts.sep)
			if err == nil {
				_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[i]), 10))
			}
		}
	}

	if err == nil && opts.end != 0 {
		err = bw.WriteByte(opts.end)
	}

	return err
}

// ----------------------------------------------------------------------------

func gcd[I Int](a, b I) I {
	if a > b {
		a, b = b, a
	}
	for a > 0 {
		a, b = b%a, a
	}
	return b
}

func gcdx(a, b int, x, y *int) int {
	if a == 0 {
		*x = 0
		*y = 1
		return b
	}
	var x1, y1 int
	d := gcdx(b%a, a, &x1, &y1)
	*x = y1 - (b/a)*x1
	*y = x1
	return d
}

func abs[N Number](a N) N {
	if a < 0 {
		return -a
	}
	return a
}

func sign[N Number](a N) N {
	if a < 0 {
		return -1
	} else if a > 0 {
		return 1
	}
	return 0
}

type Ordered interface {
	Number | ~string
}

func max[T Ordered](a, b T) T {
	if a < b {
		return b
	}
	return a
}

func min[T Ordered](a, b T) T {
	if a > b {
		return b
	}
	return a
}

// ----------------------------------------------------------------------------

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}
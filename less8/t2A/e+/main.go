package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"unsafe"
)

const (
	OR    = 0
	AND   = 1
	TRUE  = 1
	FALSE = 0
)

func and(a, b int) int {
	if a+b < 2 {
		return 0
	}
	return 1
}

type item [5]int

func calc(tree []item, i int) {
	if i*2 >= len(tree) {
		return
	}

	calc(tree, i*2)
	calc(tree, i*2+1)

	a := tree[i*2]
	b := tree[i*2+1]
	s := a[0] + b[0]

	changeOne := func() (int, int) {
		if a[4] == FALSE {
			return b[3], b[4]
		}
		if b[4] == FALSE {
			return a[3], a[4]
		}
		return min(a[3], b[3]), TRUE
	}

	changeBoth := func() (int, int) {
		return a[3] + b[3], and(a[4], b[4])
	}

	changeBothOrOne := func() (int, int) { // case OR (1 1) or AND (0 0)
		v13, v14 := changeBoth()

		if tree[i][2] == FALSE {
			return v13, v14
		}

		v23, v24 := changeOne()
		v23++ // change operation

		if v14 == FALSE {
			return v23, v24
		}
		if v24 == FALSE {
			return v13, v14
		}
		return min(v13, v23), TRUE
	}

	switch tree[i][1] {
	case OR:
		switch s {
		case 0:
			tree[i][0] = 0
			tree[i][3], tree[i][4] = changeOne()
		case 1:
			tree[i][0] = 1
			if tree[i][2] == TRUE {
				tree[i][3] = 1
				tree[i][4] = 1
				return
			}
			if a[0] == 1 {
				tree[i][3] = a[3]
				tree[i][4] = a[4]
				return
			}
			if b[0] == 1 {
				tree[i][3] = b[3]
				tree[i][4] = b[4]
				return
			}
		case 2:
			tree[i][0] = 1
			tree[i][3], tree[i][4] = changeBothOrOne()
		}
	case AND:
		switch s {
		case 0:
			tree[i][0] = 0
			tree[i][3], tree[i][4] = changeBothOrOne()
		case 1:
			tree[i][0] = 0
			if tree[i][2] == TRUE {
				tree[i][3] = 1
				tree[i][4] = 1
				return
			}
			if a[0] == 0 {
				tree[i][3] = a[3]
				tree[i][4] = a[4]
				return
			}
			if b[0] == 0 {
				tree[i][3] = b[3]
				tree[i][4] = b[4]
				return
			}
		case 2:
			tree[i][0] = 1
			tree[i][3], tree[i][4] = changeOne()
		}
	}
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, v, err := scanTwoInt(sc)
	if err != nil {
		panic(err)
	}

	tree := make([]item, n+1)
	i := 1
	for ; i <= n/2; i++ {
		g, c, err := scanTwoInt(sc)
		if err != nil {
			panic(err)
		}
		tree[i] = item{-1, g, c}
	}
	for ; i <= n; i++ {
		v, err := scanInt(sc)
		if err != nil {
			panic(err)
		}
		tree[i] = item{v}
	}

	calc(tree, 1)

	if tree[1][0] == v {
		bw.WriteString("0\n")
	} else if tree[1][4] == 0 {
		bw.WriteString("IMPOSSIBLE\n")
	} else {
		writeInt(bw, tree[1][3], defaultWriteOpts())
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

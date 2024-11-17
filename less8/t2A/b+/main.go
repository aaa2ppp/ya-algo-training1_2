package main

import (
	"bufio"
	"io"
	"math"
	"os"
	"strconv"
	"unsafe"
)

// solve ожидает, что
// граф это дерево, индексация узлов начинается с 1 (graph[0] не используется),
// в графе должен быть покрайней мере 1 узел (т.е len(graph) >= 2)
func solve(graph [][]int) (int, []int) {

	var dfs func(node, prev int) (int, int)
	dir := make([][][2]int, len(graph))

	dfs = func(node, prev int) (int, int) {
		nodeDir := make([][2]int, len(graph[node])) // {sumH, count}

		// расчитываем направления за исключением prev
		for i, neig := range graph[node] {
			if neig != prev {
				sumH, count := dfs(neig, node)
				count++
				sumH += count
				nodeDir[i] = [2]int{sumH, count}
			}
		}

		dir[node] = nodeDir

		// расчитываем сумму за исключением prev
		sumH, count := 0, 0
		for _, v := range nodeDir {
			sumH += v[0]
			count += v[1]
		}

		// возвращаем prev, что бы prev мог расчитать направление
		return sumH, count
	}

	var dfs2 func(node, prev int, sumH, count int)
	sumHH := make([]int, len(graph))

	dfs2 = func(node, prev int, sumH, count int) {
		nodeDir := dir[node]

		// заполняем пробел в направлениях для prev
		for i, neig := range graph[node] {
			if neig == prev {
				count++
				sumH += count
				nodeDir[i] = [2]int{sumH, count}
				break
			}
		}

		// расчитываем полную сумму для node
		sumH, count = 0, 0
		for _, v := range nodeDir {
			sumH += v[0]
			count += v[1]
		}

		sumHH[node] = sumH

		// передаем соседям за исключением prev
		for i, neig := range graph[node] {
			if neig != prev {
				sumH := sumH - nodeDir[i][0]
				count := count - nodeDir[i][1]
				dfs2(neig, node, sumH, count)
			}
		}
	}

	dfs(1, -1)
	dfs2(1, -1, 0, -1)

	minVal := math.MaxInt - 1
	sumHH[0] = math.MaxInt // (!) что бы 0 не попал в выборку

	for _, sumH := range sumHH {
		minVal = min(minVal, sumH)
	}

	var res []int
	for node, sumH := range sumHH {
		if sumH == minVal {
			res = append(res, node)
		}
	}

	return minVal, res
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	graph := make([][]int, n+1)

	for i := 0; i < n-1; i++ {
		a, b, err := scanTwoInt(sc)
		if err != nil {
			panic(err)
		}
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
	}

	minH, nodes := solve(graph)

	writeInts(bw, []int{minH, len(nodes)}, writeOpts{sep: ' ', end: ' '})
	writeInts(bw, nodes, writeOpts{sep: ' ', end: '\n'})
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

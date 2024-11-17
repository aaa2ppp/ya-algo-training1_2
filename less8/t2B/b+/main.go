package main

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type TreeNode struct {
	key     string
	parent  *TreeNode
	childs  []*TreeNode
	count   int
	height  int
	isChild bool
}

func (node *TreeNode) String() string {
	childs := make([]string, 0, len(node.childs))
	for _, child := range node.childs {
		childs = append(childs, child.key)
	}
	var parentKey string
	if node.parent != nil {
		parentKey = node.parent.key
	}
	return fmt.Sprintf("{%v %v %d %d}", parentKey, childs, node.DescendantCount(), node.Height())
}

func (node *TreeNode) DescendantCount() int {
	return node.count
}

func (node *TreeNode) calcDescendantCount() int {
	node.count = len(node.childs)
	for _, child := range node.childs {
		node.count += child.calcDescendantCount()
	}
	return node.count
}

func (node *TreeNode) Height() int {
	return node.height
}

func (node *TreeNode) setHeight(h int) {
	node.height = h
	for _, child := range node.childs {
		child.setHeight(h + 1)
	}
}

func (node *TreeNode) IsDescendantOf(ancestor *TreeNode) bool {
	if node == nil {
		return false
	}
	return node == ancestor || node.parent.IsDescendantOf(ancestor)
}

func (node *TreeNode) Add(child *TreeNode) {
	node.childs = append(node.childs, child)
}

type Tree map[string]*TreeNode

func NewTree(cap int) Tree {
	return make(Tree, cap)
}

func (t Tree) Keys() []string {
	return maps_Keys(t)
}

func (t Tree) Get(key string) *TreeNode {
	node := t[key]
	if node == nil {
		node = &TreeNode{key: key}
		t[key] = node
	}
	return node
}

func (t Tree) AddChild(childKey, parentKey string) {
	child := t.Get(childKey)
	parent := t.Get(parentKey)

	child.isChild = true
	child.parent = parent
	parent.Add(child)
}

func (t Tree) Calc() {
	for _, node := range t {
		if !node.isChild {
			node.setHeight(0)
			node.calcDescendantCount()
		}
	}
}

func maps_Keys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
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

	tree := NewTree(n)
	for i := 0; i < n-1; i++ {
		childName, err := scanWord(sc)
		if err != nil {
			panic(err)
		}

		parentName, err := scanWord(sc)
		if err != nil {
			panic(err)
		}

		tree.AddChild(childName, parentName)
	}

	tree.Calc()

	if debugEnable {
		log.Println("tree:", tree)
	}

	for {
		aName, err := scanWord(sc)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		bName, err := scanWord(sc)
		if err != nil {
			panic(err)
		}

		a := tree.Get(aName)
		b := tree.Get(bName)

		res := 0
		if a.Height() < b.Height() && b.IsDescendantOf(a) {
			res = 1
		} else if a.Height() > b.Height() && a.IsDescendantOf(b) {
			res = 2
		}

		writeInt(bw, res, writeOpts{end: ' '})
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

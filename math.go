package util

import (
	"fmt"
	"slices"
	"strings"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type SignedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Abs returns the absolute value of integer n.
func Abs[T Integer](n T) T {
	if n < 0 {
		return -n
	}
	return n
}

// Sign returns -1 if n < 0, 1 if n > 0, and 0 if n == 0.
func Sign[T SignedInteger](n T) T {
	if n < 0 {
		return -1
	}
	if n > 0 {
		return 1
	}
	return 0
}

func gcd[T Integer](a, b T) T {
	if a < b {
		a, b = b, a
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm[T Integer](a, b T) T {
	return a / GCD(a, b) * b
}

// GCD returns the greatest common divisor of all arguments.
func GCD[T Integer](a, b T, more ...T) T {
	g := gcd(a, b)
	for _, n := range more {
		g = gcd(g, n)
	}
	return g
}

// LCM returns the least common multiple of all arguments.
func LCM[T Integer](a, b T, more ...T) T {
	l := lcm(a, b)
	for _, n := range more {
		l = lcm(l, n)
	}
	return l
}

// ToBalancedQuinary converts integer n to a balanced quinary string.
//
// In returned string, digit -2 is represented by '=', digit -1 by '-', digit 0, 1, 2 by '0', '1', '2' respectively.
func ToBalancedQuinary(n int64) string {
	if n == 0 {
		return "0"
	}

	digits := []int{}
	sign := n > 0
	n = Abs(n)

	for n > 0 {
		digits = append(digits, int(n%5))
		n /= 5
	}
	digits = append(digits, 0)

	res := make([]byte, 0, len(digits))
	for i := 0; i < len(digits); i++ {
		switch digits[i] {
		case 0, 1, 2:
			res = append(res, '0'+byte(digits[i]))
		case 3:
			res = append(res, '=')
			digits[i+1]++
		case 4:
			res = append(res, '-')
			digits[i+1]++
		case 5:
			res = append(res, '0')
			digits[i+1]++
		}
	}
	slices.Reverse(res)
	if res[0] == '0' {
		res = res[1:]
	}

	if !sign {
		return NegateBalancedQuinary(string(res))
	}
	return string(res)
}

// FromBalancedQuinary converts balanced quinary string s to integer.
//
// In s, digit -2 is represented by '=', digit -1 by '-', digit 0, 1, 2 by '0', '1', '2' respectively.
// If other characters are present, an error is returned.
func FromBalancedQuinary(s string) (int64, error) {
	val := int64(0)
	base := int64(1)

	for i := len(s) - 1; i >= 0; i-- {
		switch s[i] {
		case '=':
			val += base * -2
		case '-':
			val += base * -1
		case '0':
			// noop
		case '1':
			val += base * 1
		case '2':
			val += base * 2
		default:
			return 0, fmt.Errorf("invalid balanced quinary digit %c at index %d", s[i], i)
		}
		base *= 5
	}

	return val, nil
}

// NegateBalancedQuinary returns the nagative value of balanced quinary string s.
//
// In s, digit -2 is represented by '=', digit -1 by '-', digit 0, 1, 2 by '0', '1', '2' respectively.
// If other characters are present, this function panics.
func NegateBalancedQuinary(s string) string {
	b := strings.Builder{}
	b.Grow(len(s))
	for _, c := range s {
		switch c {
		case '0':
			b.WriteByte('0')
		case '=':
			b.WriteByte('2')
		case '-':
			b.WriteByte('1')
		case '1':
			b.WriteByte('-')
		case '2':
			b.WriteByte('=')
		default:
			panic(fmt.Errorf("invalid balanced quinary digit %c", c))
		}
	}
	return b.String()
}

// Point2I[T] represents a point in 2D integer space.
type Point2I[T Integer] struct {
	X, Y T
}

// NewPoint2I converts [2]Integer to Point2I[T].
func NewPoint2I[T Integer](cord [2]T) Point2I[T] {
	return Point2I[T]{cord[0], cord[1]}
}

// Add returns the sum of two points.
func (p Point2I[T]) Add(q Point2I[T]) Point2I[T] {
	return Point2I[T]{p.X + q.X, p.Y + q.Y}
}

func (p Point2I[T]) AddCord(q [2]T) Point2I[T] {
	return Point2I[T]{p.X + q[0], p.Y + q[1]}
}

// Sub returns the difference of two points.
func (p Point2I[T]) Sub(q Point2I[T]) Point2I[T] {
	return Point2I[T]{p.X - q.X, p.Y - q.Y}
}

// Mul returns the product of a point and an integer.
func (p Point2I[T]) Mul(n T) Point2I[T] {
	return Point2I[T]{p.X * n, p.Y * n}
}

// Div returns the quotient of a point and an integer.
func (p Point2I[T]) Div(n T) Point2I[T] {
	return Point2I[T]{p.X / n, p.Y / n}
}

// Neg returns the negative value of a point.
func (p Point2I[T]) Neg() Point2I[T] {
	return Point2I[T]{-p.X, -p.Y}
}

// Less returns true if p is lexicographically less than q.
func (p Point2I[T]) Less(rhs Point2I[T]) bool {
	if p.X == rhs.X {
		return p.Y < rhs.Y
	}
	return p.X < rhs.X
}

type ByIndex[T Integer] []Point2I[T]

func (a ByIndex[T]) Len() int           { return len(a) }
func (a ByIndex[T]) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIndex[T]) Less(i, j int) bool { return a[i].Less(a[j]) }

type Point3I[T Integer] struct {
	X, Y, Z T
}

func NewPoint3I[T Integer](cord [3]T) Point3I[T] {
	return Point3I[T]{cord[0], cord[1], cord[2]}
}

func (p Point3I[T]) Sub(r Point3I[T]) Point3I[T] {
	return Point3I[T]{p.X - r.X, p.Y - r.Y, p.Z - r.Z}
}

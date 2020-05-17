package tpl_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/qiniu/text/tpl"
	"github.com/qiniu/text/tpl/qlang"
	_ "github.com/qiniu/text/tpl/qlang/lib/builtin" // 导入 builtin 包
)

// -----------------------------------------------------------------------------

const maxprime = `
primes = [2, 3]
n = 1
limit = 9

isPrime = fn(v) {
	for i = 0; i < n; i++ {
		if v % primes[i] == 0 {
			return false
		}
	}
	return true
}

listPrimes = fn(max) {
	v = 5
	for {
		for v < limit {
			if isPrime(v) {
				primes = append(primes, v)
				if v * v >= max {
					return
				}
			}
			v += 2
		}
		v += 2
		n; n++
		limit = primes[n] * primes[n]
	}
}

maxPrimeOf = fn(max) {
	if max % 2 == 0 {
		max--
	}

	listPrimes(max)
	n; n = len(primes)

	for {
		if isPrime(max) {
			return max
		}
		max -= 2
	}
}

v = maxPrimeOf(10000)
`

func TestMaxPrime(t *testing.T) {
	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal(err)
	}

	err = lang.SafeEval(maxprime)
	if err != nil {
		t.Fatal(err)
	}

	v, _ := lang.Var("v")
	fmt.Println(v)
	if v != 9973 {
		t.Fatal("MaxPrime ret:", v)
	}
}

// -----------------------------------------------------------------------------

func TestCompiler(t *testing.T) {
	const grammar = `

term = factor *(('*' factor)/mul | ('/' factor)/div)

doc = term *(('+' term)/add | ('-' term)/sub)
// this is a comment

factor =
	FLOAT/push |
	('-' factor)/neg |
	'(' doc ')' |
	(IDENT '(' doc % ','/arity ')')/call |
	'+' factor
`
	var stk []string
	var args tpl.Grammar

	scanner := new(tpl.Scanner)
	push := func(tokens []tpl.Token, g tpl.Grammar) {
		v := tokens[0].Literal
		if v == "" {
			v = scanner.Ttol(tokens[0].Kind)
		} else if tokens[0].Kind == tpl.IDENT {
			v += "/" + strconv.Itoa(args.Len())
		}
		stk = append(stk, v)
	}

	marker := func(g tpl.Grammar, mark string) tpl.Grammar {
		if mark == "arity" {
			args = g
			return g
		}
		return tpl.Action(g, push)
	}

	compiler := &tpl.Compiler{
		Grammar: []byte(grammar),
		Marker:  marker,
	}
	m, err := compiler.Cl()
	if err != nil {
		t.Fatal("compiler.Cl failed:", err)
	}

	err = m.MatchExactly([]byte(`max(1.2 + sin(.3) * 2, cos(3), pow(5, 6), 7)`), "")
	if err != nil {
		t.Fatal("MatchExactly failed:\n", err)
	}
	text := strings.Join(stk, " ")
	if text != "1.2 .3 sin/1 2 '*' '+' 3 cos/1 5 6 pow/2 7 max/4" {
		t.Fatal("MatchExactly failed:", text)
	}
}

// -----------------------------------------------------------------------------

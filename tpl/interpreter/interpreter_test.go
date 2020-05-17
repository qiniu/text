package interpreter_test

import (
	"fmt"
	"testing"

	"github.com/qiniu/text/tpl/interpreter"
	"github.com/qiniu/text/tpl/qlang"
	_ "github.com/qiniu/text/tpl/qlang/lib/builtin" // 导入 builtin 包
)

// -----------------------------------------------------------------------------

const codeIf = `
today = 3

if today == 1 {
	today = "Mon"
} elif today == 2 {
	today = "Tue"
} elif today == 3 {
	today = "Wed"
}

println(today)
`

func TestIf(t *testing.T) {
	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal(err)
	}

	err = lang.SafeEval(codeIf)
	if err != nil {
		t.Fatal(err)
	}

	v, _ := lang.Var("today")
	fmt.Println(v)
	if v != "Wed" {
		t.Fatal("MaxPrime ret:", v)
	}
}

// -----------------------------------------------------------------------------

func TestParseInt(t *testing.T) {
	if v, err := interpreter.ParseInt("0"); err != nil || v != 0 {
		t.Fatal(`ParseInt("0")`, v, err)
	}

	if v, err := interpreter.ParseInt("012"); err != nil || v != 012 {
		t.Fatal(`ParseInt("012")`, v, err)
	}

	if v, err := interpreter.ParseInt("0x12"); err != nil || v != 18 {
		t.Fatal(`ParseInt("0x12")`, v, err)
	}
}

func TestParseFloat(t *testing.T) {
	if v, err := interpreter.ParseFloat("0"); err != nil || v != 0 {
		t.Fatal(`ParseFloat("0")`, v, err)
	}

	if v, err := interpreter.ParseFloat("012"); err != nil || v != 012 {
		t.Fatal(`ParseFloat("012")`, v, err)
	}

	if v, err := interpreter.ParseFloat("012.34"); err != nil || v != 012.34 {
		t.Fatal(`ParseFloat("012.34")`, v, err)
	}

	if v, err := interpreter.ParseFloat("0x12"); err != nil || v != 0x12 {
		t.Fatal(`ParseFloat("0x12")`, v, err)
	}
}

// -----------------------------------------------------------------------------

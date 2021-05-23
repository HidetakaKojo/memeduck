package memeduck_test

import (
	"testing"

	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/genkami/memeduck"
	"github.com/stretchr/testify/assert"
)

type whereCond interface {
	ToASTWhere() (*ast.Where, error)
}

func test(t *testing.T, cond whereCond, expected string) {
	w, err := cond.ToASTWhere()
	assert.Nil(t, err, expected)
	assert.Equal(t, expected, w.Expr.SQL())
}

func TestBool(t *testing.T) {
	test(t, memeduck.Bool(true), `TRUE`)
	test(t, memeduck.Bool(false), `FALSE`)
}

func TestOp(t *testing.T) {
	test(t, memeduck.Op(1, memeduck.EQ, 1), `1 = 1`)
	test(t, memeduck.Op("hoge", memeduck.NE, "fuga"), `"hoge" != "fuga"`)
	test(t, memeduck.Op(1.23, memeduck.LT, 4.56), `1.23e+00 < 4.56e+00`)
	test(t, memeduck.Op(4.56, memeduck.GT, 1.23), `4.56e+00 > 1.23e+00`)
	test(t, memeduck.Op(1, memeduck.LE, 2), `1 <= 2`)
	test(t, memeduck.Op(2, memeduck.GE, 1), `2 >= 1`)

	test(t, memeduck.Eq(1, 1), `1 = 1`)
	test(t, memeduck.Ne("hoge", "fuga"), `"hoge" != "fuga"`)
	test(t, memeduck.Lt(1.23, 4.56), `1.23e+00 < 4.56e+00`)
	test(t, memeduck.Gt(4.56, 1.23), `4.56e+00 > 1.23e+00`)
	test(t, memeduck.Le(1, 2), `1 <= 2`)
	test(t, memeduck.Ge(2, 1), `2 >= 1`)
}

func TestAnd(t *testing.T) {
	_, err := memeduck.And().ToASTWhere()
	assert.Error(t, err, "empty AND")
	test(t,
		memeduck.And(
			memeduck.Op(1, memeduck.EQ, 1),
		),
		`1 = 1`,
	)
	test(t,
		memeduck.And(
			memeduck.Op(1, memeduck.EQ, 1),
			memeduck.Op("hoge", memeduck.EQ, "hoge"),
		),
		`1 = 1 AND "hoge" = "hoge"`,
	)
	test(t,
		memeduck.And(
			memeduck.Op(1, memeduck.EQ, 1),
			memeduck.Op("hoge", memeduck.EQ, "hoge"),
			memeduck.Op(true, memeduck.EQ, true),
		),
		`1 = 1 AND "hoge" = "hoge" AND TRUE = TRUE`,
	)
}

func TestOr(t *testing.T) {
	_, err := memeduck.Or().ToASTWhere()
	assert.Error(t, err, "empty Or")
	test(t,
		memeduck.Or(
			memeduck.Op(1, memeduck.EQ, 1),
		),
		`1 = 1`,
	)
	test(t,
		memeduck.Or(
			memeduck.Op(1, memeduck.EQ, 1),
			memeduck.Op("hoge", memeduck.EQ, "hoge"),
		),
		`1 = 1 OR "hoge" = "hoge"`,
	)
	test(t,
		memeduck.Or(
			memeduck.Op(1, memeduck.EQ, 1),
			memeduck.Op("hoge", memeduck.EQ, "hoge"),
			memeduck.Op(true, memeduck.EQ, true),
		),
		`1 = 1 OR "hoge" = "hoge" OR TRUE = TRUE`,
	)

	// TODO: this shoud pass
	// test(t,
	// 	memeduck.And(
	// 		memeduck.Eq(1, 1),
	// 		memeduck.Or(
	// 			memeduck.Eq(2, 2),
	// 			memeduck.Eq(3, 3),
	// 		),
	// 	),
	// 	`1 = 1 AND (2 = 2 OR 3 = 3)`,
	// )
}

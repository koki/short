package expressions

import (
	"strings"

	"github.com/koki/short/util"
)

// Expr is the generic AST format of a koki NodeSelectorRequirement or LabelSelectorRequirement
type Expr struct {
	Key    string
	Op     string
	Values []string
}

func ParseExpr(s string, ops []string) (*Expr, error) {
	for _, op := range ops {
		x, err := ParseOp(s, op)
		if err != nil {
			return nil, err
		}

		if x != nil {
			return x, nil
		}
	}

	return nil, nil
}

func ParseOp(s string, op string) (*Expr, error) {
	if strings.Contains(s, op) {
		segs := strings.Split(s, op)
		if len(segs) != 2 {
			return nil, util.InvalidValueErrorf(s, "not a valid expression with operator (%s)", op)
		}

		return &Expr{
			Key:    segs[0],
			Op:     op,
			Values: ParseValues(segs[1]),
		}, nil
	}

	return nil, nil
}

func ParseValues(s string) []string {
	return strings.Split(s, ",")
}

package assert

import (
	"fmt"

	"github.com/eddycharly/json-kyverno/pkg/engine/template"
	reflectutils "github.com/eddycharly/json-kyverno/pkg/utils/reflect"
	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// scalarNode is a terminal type of assertion.
// it receives a value and compares it with an expected value.
// the expected value can be the result of an expression.
type scalarNode struct {
	rhs interface{}
}

func (n *scalarNode) assert(path *field.Path, value interface{}, bindings binding.Bindings) (field.ErrorList, error) {
	rhs := n.rhs
	expression := parseExpression(rhs)
	if expression != nil {
		// we only project if the expression uses the engine syntax
		// this is to avoid the case where the value is a map and the RHS is a string
		// TODO: we need a way to escape projection
		if expression.engine != "" {
			if expression.foreach {
				return nil, field.Invalid(path, rhs, "foreach is not supported on the RHS")
			}
			// TODO: something is wrong here, we should jump into a more standard matching at this point
			projected, err := template.Execute(expression.statement, value, bindings)
			if err != nil {
				return nil, field.InternalError(path, err)
			}
			rhs = projected
		}
	}
	var errs field.ErrorList
	if match, err := reflectutils.MatchScalar(rhs, value); err != nil {
		return nil, field.InternalError(path, err)
	} else if !match {
		errs = append(errs, field.Invalid(path, value, fmt.Sprint("Expected value:", rhs)))
	}
	return errs, nil
}

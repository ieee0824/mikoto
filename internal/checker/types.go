package checker

import "go/types"

// typeTerms expands the structural type restrictions of a constraint. A nil
// slice is unrestricted; a non-nil empty slice is an empty intersection.
// Interface embeddings intersect restrictions, while union terms combine them.
// Method requirements do not introduce additional underlying types.
func typeTerms(t types.Type) []*types.Term {
	if t == nil {
		return nil
	}
	t = types.Unalias(t)
	switch t := t.(type) {
	case *types.TypeParam:
		return typeTerms(t.Constraint())
	case *types.Union:
		terms := []*types.Term{}
		for i := 0; i < t.Len(); i++ {
			term := t.Term(i)
			part := []*types.Term{term}
			if !term.Tilde() {
				part = typeTerms(term.Type())
			}
			if part == nil {
				return nil
			}
			for _, term := range part {
				terms = addTerm(terms, term)
			}
		}
		return terms
	default:
		if iface, ok := t.Underlying().(*types.Interface); ok {
			var terms []*types.Term
			for i := 0; i < iface.NumEmbeddeds(); i++ {
				terms = intersectTerms(terms, typeTerms(iface.EmbeddedType(i)))
			}
			return terms
		}
		return []*types.Term{types.NewTerm(false, t)}
	}
}

func intersectTerms(a, b []*types.Term) []*types.Term {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	terms := []*types.Term{}
	for _, left := range a {
		for _, right := range b {
			switch {
			case types.Identical(left.Type(), right.Type()):
				if left.Tilde() {
					terms = addTerm(terms, right)
				} else {
					terms = addTerm(terms, left)
				}
			case left.Tilde() && types.Identical(left.Type(), right.Type().Underlying()):
				terms = addTerm(terms, right)
			case right.Tilde() && types.Identical(left.Type().Underlying(), right.Type()):
				terms = addTerm(terms, left)
			}
		}
	}
	return terms
}

func addTerm(terms []*types.Term, term *types.Term) []*types.Term {
	for _, existing := range terms {
		if existing.Tilde() == term.Tilde() && types.Identical(existing.Type(), term.Type()) {
			return terms
		}
	}
	return append(terms, term)
}

func arrayPointer(t types.Type) bool {
	for _, term := range typeTerms(t) {
		if ptr, ok := term.Type().Underlying().(*types.Pointer); ok {
			if _, ok := ptr.Elem().Underlying().(*types.Array); ok {
				return true
			}
		}
	}
	return false
}

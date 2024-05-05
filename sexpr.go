package sexpr

import (
	"errors"
	"strings"
)

type Sexpr struct {
	name   string
	params []*SexprParam

	parent *Sexpr
	line   int
	col    int
}

type FindPredicate func(sexpr *Sexpr, depth int) bool

func NewSexpr(name string) *Sexpr {
	return &Sexpr{name: name}
}

func (s *Sexpr) Name() string {
	return s.name
}

func (s *Sexpr) SetName(name string) {
	s.name = name
}

func (s *Sexpr) Params() []*SexprParam {
	return s.params
}

func (s *Sexpr) AddParam(idx int, param *SexprParam) error {
	if idx < 0 || idx > len(s.params) {
		return errors.New("index out of range")
	}
	if idx == len(s.params) {
		s.params = append(s.params, param)
	} else {
		s.params = append(s.params[:idx+1], s.params[idx:]...)
		s.params[idx] = param
	}
	return nil
}

func (s *Sexpr) SetParam(idx int, param *SexprParam) error {
	if idx < 0 || idx >= len(s.params) {
		return errors.New("index out of range")
	}
	s.params[idx] = param
	param.SetParent(s)
	return nil
}

func (s *Sexpr) SetParamSexpr(idx int, sexpr *Sexpr) error {
	if idx < 0 || idx >= len(s.params) {
		return errors.New("index out of range")
	}
	param, err := NewSexprParam(sexpr)
	if err != nil {
		return err
	}
	return s.SetParam(idx, param)
}

func (s *Sexpr) SetParamString(idx int, sexprStr *SexprString) error {
	if idx < 0 || idx >= len(s.params) {
		return errors.New("index out of range")
	}
	param, err := NewSexprParam(sexprStr)
	if err != nil {
		return err
	}
	return s.SetParam(idx, param)
}

func (s *Sexpr) RemoveParam(idx int, param *SexprParam) error {
	if idx < 0 || idx >= len(s.params) {
		return errors.New("index out of range")
	}
	s.params = append(s.params[:idx], s.params[idx+1:]...)
	return nil
}

func (s *Sexpr) Parent() *Sexpr {
	return s.parent
}

func (s *Sexpr) SetParent(parent *Sexpr) {
	s.parent = parent
}

func (s *Sexpr) Location() (int, int) {
	return s.line, s.col
}

func (s *Sexpr) SetLocation(line int, col int) {
	s.line = line
	s.col = col
}

func (s *Sexpr) String() string {
	var sb strings.Builder
	s.string_(&sb, 0)
	return sb.String()
}

func (s *Sexpr) string_(acc *strings.Builder, level int) {
	wasSexpr := false
	indent := strings.Repeat("\t", level)
	params := s.params
	acc.WriteString(indent)
	acc.WriteString("(")
	acc.WriteString(s.Name())
	if len(s.Params()) == 0 {
		acc.WriteString(")\n")
	} else {
		for _, param := range params {
			paramv := param.Value()
			switch spv := paramv.(type) {
			case *Sexpr:
				acc.WriteString("\n")
				spv.string_(acc, level+1)
				wasSexpr = true
			case *SexprString:
				acc.WriteString(" ")
				acc.WriteString(spv.String())
				wasSexpr = false
			}
		}
		if wasSexpr {
			acc.WriteString("\n")
			acc.WriteString(indent)
			acc.WriteString(")")
		} else {
			acc.WriteString(")")
		}
	}
}

func (s *Sexpr) FindChild(fp FindPredicate, maxDepth int) *Sexpr {
	queue := NewSexprQueue()
	queue.Enqueue(s)
	depth := 1

	for {
		sexpr := queue.Dequeue()
		if sexpr == nil {
			return nil
		}
		for _, param := range sexpr.params {
			if sexpr, ok := param.Value().(*Sexpr); ok {
				if found := fp(sexpr, depth); found {
					return sexpr
				}
				if maxDepth == -1 || depth < maxDepth {
					queue.Enqueue(sexpr)
				}
			}
		}
		depth += 1
	}
}

func (s *Sexpr) FindChildren(fp FindPredicate, maxDepth int) []*Sexpr {
	children := []*Sexpr{}

	queue := NewSexprQueue()
	queue.Enqueue(s)
	depth := 1

	for {
		sexpr := queue.Dequeue()
		if sexpr == nil {
			return children
		}
		for _, param := range sexpr.params {
			if sexpr, ok := param.Value().(*Sexpr); ok {
				if found := fp(sexpr, depth); found {
					children = append(children, sexpr)
				}
				if maxDepth == -1 || depth < maxDepth {
					queue.Enqueue(sexpr)
				}
			}
		}
		depth += 1
	}
}

func (s *Sexpr) FindChildByName(name string, maxDepth int) *Sexpr {
	return s.FindChild(func(s *Sexpr, d int) bool {
		return s.Name() == name
	}, maxDepth)
}

func (s *Sexpr) FindChildrenByName(name string, maxDepth int) []*Sexpr {
	return s.FindChildren(func(s *Sexpr, d int) bool {
		return s.Name() == name
	}, maxDepth)
}

func (s *Sexpr) FindDirectChildByName(name string) *Sexpr {
	return s.FindChildByName(name, 1)
}

func (s *Sexpr) FindDirectChildrenByName(name string) []*Sexpr {
	return s.FindChildrenByName(name, 1)
}

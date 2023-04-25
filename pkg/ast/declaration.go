package ast

import (
	"bytes"
	"fmt"
)

type Declaration interface {
	Node
	declarationNode()
	Type() Declaration
	SetType(Declaration)
}

type Pointer struct {
	PointsTo Declaration
	Const    bool
	Volatile bool
}

func (p *Pointer) declarationNode() {}
func (p *Pointer) String() string {
	var ret string
	if p.PointsTo != nil {
		ret = fmt.Sprintf("(%s) *", p.PointsTo.String())
	} else {
		ret = "(nil) *" // this should not happen
	}

	if p.Const {
		ret = ret + " const"
	}

	if p.Volatile {
		ret = ret + " volatile"
	}

	return ret
}

func (p *Pointer) Type() Declaration {
	return p.PointsTo
}

func (p *Pointer) SetType(decl Declaration) {
	p.PointsTo = decl
}

type Array struct {
	ArrayOf   Declaration
	ArraySize Expression
}

func (a *Array) declarationNode() {}

func (a *Array) String() string {
	if a.ArrayOf != nil {
		return fmt.Sprintf("(%s)[]", a.ArrayOf.String())
	} else {
		return fmt.Sprintf("(nil)[]") // shouldn't happen
	}
}

func (a *Array) Type() Declaration {
	return a.ArrayOf
}

func (a *Array) SetType(decl Declaration) {
	a.ArrayOf = decl
}

type TypeSpecification struct {
	Name     string // "int" or "long long int" or custom typedef identifier "uint8_t" etc..
	Const    bool
	Volatile bool
}

func (t *TypeSpecification) declarationNode() {}

func (t *TypeSpecification) String() string {
	var out bytes.Buffer

	if t.Const {
		out.WriteString("const ")
	}

	if t.Volatile {
		out.WriteString("volatile ")
	}

	out.WriteString(t.Name)

	return out.String()
}

func (t *TypeSpecification) Type() Declaration {
	return nil
}

func (t *TypeSpecification) SetType(d Declaration) {} // no op

type StructOrUnionSpecification struct{} //TODO

type VariableDeclaration struct {
	StorageClass string
	TypeSpec     Declaration
	Name         string
	Definition   Expression
}

func (v *VariableDeclaration) declarationNode() {}
func (v *VariableDeclaration) String() string   { return "" }
func (v *VariableDeclaration) Type() Declaration {
	return v.TypeSpec
}

func (v *VariableDeclaration) SetType(d Declaration) {
	v.TypeSpec = d
}

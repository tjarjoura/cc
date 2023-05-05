package ast

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
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
	ptr, ok := decl.(*Pointer)
	if ok && ptr == p {
		panic("creating self-referential pointer!!")
	}

	p.PointsTo = decl
}

type Array struct {
	ArrayOf   Declaration
	ArraySize Expression
}

func (a *Array) declarationNode() {}

func (a *Array) String() string {
	var expr = ""
	if a.ArraySize != nil {
		expr = a.ArraySize.String()
	}

	if a.ArrayOf != nil {
		return fmt.Sprintf("(%s)[%s]", a.ArrayOf.String(), expr)
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

type BaseType struct {
	Name     string // "int" or "long long int" or custom typedef identifier "uint8_t" etc..
	Const    bool
	Volatile bool
	Signed   bool
}

func (t *BaseType) declarationNode() {}

func (t *BaseType) String() string {
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

func (t *BaseType) Type() Declaration {
	return nil
}

func (t *BaseType) SetType(d Declaration) {} // no op

type StructOrUnionSpecification struct{} //TODO

type FunctionDeclaration struct {
	Name         string
	StorageClass string
	ReturnType   Declaration
	Parameters   []Declaration
	Body         *BlockStatement
}

func (f *FunctionDeclaration) declarationNode() {}
func (f *FunctionDeclaration) String() string {
	paramTypes := []string{}
	for _, param := range f.Parameters {
		paramTypes = append(paramTypes, param.String())
	}

	return fmt.Sprintf("%v %s(%s)", f.ReturnType.String(), f.Name,
		strings.Join(paramTypes, ", "))
}
func (f *FunctionDeclaration) Type() Declaration     { return f.ReturnType }
func (f *FunctionDeclaration) SetType(d Declaration) { f.ReturnType = d }

type VariableDeclaration struct {
	Name         string
	StorageClass string
	VarType      Declaration
	Definition   Expression
}

func (v *VariableDeclaration) declarationNode() {}
func (v *VariableDeclaration) String() string {
	ret := fmt.Sprintf("%s %s", v.VarType, v.Name)
	if v.Definition != nil {
		return fmt.Sprintf("%s = %s", ret, v.Definition.String())
	}

	return ret
}

func (v *VariableDeclaration) Type() Declaration     { return v.VarType }
func (v *VariableDeclaration) SetType(d Declaration) { v.VarType = d }

func ConvertError(to Declaration, from Declaration) bool {
	validCombinations := [][]string{
		[]string{"Array", "Pointer"},
		[]string{"Array", "BaseType"},
		[]string{"BaseType", "Pointer"},
	}

	toType, fromType := reflect.TypeOf(to).Name(), reflect.TypeOf(from).Name()

	if toType == fromType {
		// TODO handle struct, union, and typedef types
		return false
	}

	for _, combo := range validCombinations {
		if (toType == combo[0] && fromType == combo[1]) ||
			(toType == combo[1] && fromType == combo[0]) {
			return true
		}
	}

	return false
}

func ConvertWarn(to Declaration, from Declaration) bool {
	// TODO implement
	return false
}

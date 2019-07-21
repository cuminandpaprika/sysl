package relgom

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	sysl "github.com/anz-bank/sysl/src/proto"
	. "github.com/anz-bank/sysl/sysl2/codegen/golang" //nolint:golint,stylecheck
	"github.com/anz-bank/sysl/sysl2/language/go/pkg/codegen"
)

type commonModules struct {
	json      func(name string) Expr
	relgomlib func(name string) Expr
	seq       func(name string) Expr
}

func newCommonModules(g *sourceGenerator) *commonModules {
	return &commonModules{
		json:      g.imported("encoding/json"),
		relgomlib: g.imported("github.com/anz-bank/sysl/sysl2/language/go/pkg/relgom/relgomlib"),
		seq:       g.imported("github.com/mediocregopher/seq"),
	}
}

type sourceGenerator struct {
	fsw             codegen.FileSystemWriter
	packageName     string
	model           *sysl.Application
	builtinImports  map[string]struct{}
	externalImports map[string]struct{}
}

func newSourceGenerator(fsw codegen.FileSystemWriter, packageName string, model *sysl.Application) *sourceGenerator {
	return &sourceGenerator{
		fsw:             fsw,
		packageName:     packageName,
		model:           model,
		builtinImports:  map[string]struct{}{},
		externalImports: map[string]struct{}{},
	}
}

func (g *sourceGenerator) genSourceForDecls(
	basepath string, decls ...Decl,
) error {
	return g.genSourceForFile(basepath, &File{Decls: decls})
}

func (g *sourceGenerator) genSourceForFile(basepath string, file *File) error {
	{
		file := *file
		file.Doc = codegen.Prelude()
		file.Name = *I(g.packageName)
		file.Imports = ImportGroups(
			Imports(sortedSetElements(g.builtinImports)...),
			Imports(sortedSetElements(g.externalImports)...),
		)

		var buf bytes.Buffer
		fmt.Fprintln(&buf, &file)

		final, err := format.Source(buf.Bytes())
		if err != nil {
			var lines bytes.Buffer
			for i, line := range strings.Split(buf.String(), "\n") {
				fmt.Fprintf(&lines, "%3d: %s\n", 1+i, line)
			}
			logrus.Errorf("Error formatting %#v:\n%s", basepath+".go", lines.String())
			return errors.Wrap(err, "gofmt")
		}

		f, err := g.fsw.Create(basepath + ".go")
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(final)
		return err
	}
}

func (g *sourceGenerator) imported(imp string, extra ...string) func(name string) Expr { //nolint:unparam
	var pkgName string
	switch len(extra) {
	case 0:
		pkgName = imp[strings.LastIndex(imp, "/")+1:]
	case 1:
		pkgName = extra[0]
	default:
		panic("Too many parameters to imported()")
	}

	var imports map[string]struct{}
	if strings.Contains(imp, ".") {
		imports = g.externalImports
	} else {
		imports = g.builtinImports
	}

	return func(name string) Expr {
		imports[imp] = struct{}{}
		return Dot(I(pkgName), name)
	}
}

type typeInfo struct {
	final Expr
	param Expr
	opt   bool
	fkey  *sysl.ScopedRef
}

func (g *sourceGenerator) typeInfoForSyslType(t *sysl.Type) typeInfo {
	var ti typeInfo
	switch t := t.Type.(type) {
	case *sysl.Type_Primitive_:
		switch t.Primitive {
		case sysl.Type_BOOL:
			ti.final = I("bool")
		case sysl.Type_INT:
			ti.final = I("int64")
		case sysl.Type_FLOAT:
			ti.final = I("float64")
		case sysl.Type_DECIMAL:
			ti.final = g.imported("github.com/anz-bank/decimal")("Decimal64")
		case sysl.Type_DATE, sysl.Type_DATETIME:
			ti.final = g.imported("time")("Time")
		case sysl.Type_STRING, sysl.Type_STRING_8:
			ti.final = I("string")
		default:
			panic(fmt.Errorf("type: %#v", t))
		}
	case *sysl.Type_TypeRef:
		fkInfo := g.typeInfoForSyslType(g.typeForScopedRef(t.TypeRef))
		ti = typeInfo{final: fkInfo.final, fkey: t.TypeRef}
	default:
		panic(fmt.Errorf("type: %#v", t))
	}
	ti.param = ti.final
	if t.Opt {
		ti.opt = true
		ti.final = Star(ti.final)
	}
	return ti
}

func (g *sourceGenerator) typeForScopedRef(t *sysl.ScopedRef) *sysl.Type {
	if len(t.GetRef().GetAppname().GetPart()) > 0 {
		panic(fmt.Errorf("non-local refs not implemented: %#v", t.Ref))
	}
	if len(t.Ref.Path) != 2 {
		panic(fmt.Errorf("ScopedRef path must be length 2: %#v", t.Ref)) //nolint:golint
	}
	return g.model.Types[t.Ref.Path[0]].Type.(*sysl.Type_Relation_).Relation.AttrDefs[t.Ref.Path[1]]
}

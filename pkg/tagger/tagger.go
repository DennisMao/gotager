package tagger

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"strings"

	"github.com/DennisMao/strcase"
)

const (
	STYLE_RAW         = "raw"
	STYLE_LOWER       = "lower"
	STYLE_CAMEL       = "camel"
	STYLE_CAMEL_LOWER = "camel_lower"
	STYLE_SNAKE       = "snake"
)

type Tag struct {
	opt *TagOpt
}

type TagOpt struct {
	Overwrite bool   // Tag overwrite setting.If true,tagger will overwrite the tag that exist.
	Style     string // Tag style,tag name is generated by converting field name on setting style.
}

var (
	GFset *token.FileSet
	Print = true
)

func New(opt *TagOpt) *Tag {
	if opt == nil {
		opt = &TagOpt{false, "snake"}
	}

	return &Tag{opt}
}

func (this *Tag) Tag(src []byte, output io.Writer, re, prefix string) error {

	debugf("[New Tag] re:%s prefix:%s style=%s", re, prefix, this.opt.Style)

	fset := token.NewFileSet() // positions are relative to fset

	GFset = fset
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	tagFunc := func(n ast.Node) bool {
		s, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if re != "*" && !strings.HasPrefix(s.Name.Name, re) {
			return true
		}

		debugf("Find struct '%s' on pos:%s", s.Name.Name, GFset.Position(s.Pos()))

		t, ok := (s.Type).(*ast.StructType)
		if !ok {
			return true
		}

		// Ignore the structure without any fields
		if t.Fields.NumFields() == 0 {
			return true
		}

		for i, _ := range t.Fields.List {

			if len(t.Fields.List[i].Names) == 0 {
				continue
			}

			for j, _ := range t.Fields.List[i].Names {

				// If tag not exists on currenct field,allocate memory for Tag struct.
				if t.Fields.List[i].Tag == nil {
					t.Fields.List[i].Tag = new(ast.BasicLit)
				}

				debugf("Find field '%s' on pos:%s \n", t.Fields.List[i].Names[j].Name, GFset.Position(t.Fields.List[i].Pos()))

				// Replace the old tag
				t.Fields.List[i].Tag.Value = tagGenerate(
					t.Fields.List[i].Tag.Value,
					t.Fields.List[i].Names[j].Name,
					prefix,
					this.opt.Style,
					this.opt.Overwrite,
				)

			}

		}
		return true
	}

	ast.Inspect(f, tagFunc)

	err = printer.Fprint(output, fset, f)
	if err != nil {
		return err
	}

	return nil
}

// match a tag by prefix on an old tag
func tagSearch(rawTag, prefix string) string {
	startIdx := strings.Index(rawTag, prefix)
	if startIdx == -1 {
		return ""
	}

	for i := startIdx; i < len(rawTag); i++ {
		if rawTag[i] == ' ' {
			return rawTag[startIdx:i]
		}
	}

	return ""
}

// tagGenerate returns a new tag generated by rules
func tagGenerate(rawTag, fieldName, prefix, style string, isOverwrite bool) string {
	idxTagExist := strings.Index(rawTag, prefix)
	if idxTagExist != -1 {
		// If overwrite setting is not enabled ,just skip this field.
		if !isOverwrite {
			return rawTag
		}

		// If not,remove this tag before appending.
		idxTagExistEnd := idxTagExist
		lenRawTag := len(rawTag)
		isFindTagStart := false
		for idxTagExistEnd < lenRawTag-1 {
			if rawTag[idxTagExistEnd] == '"' {
				if isFindTagStart {
					// To find the end index of current tag and remote it
					rawTag = rawTag[:idxTagExist] + rawTag[idxTagExistEnd+1:]
					break
				}
				isFindTagStart = true
			}
			idxTagExistEnd++
		}
	}

	if len(rawTag) > 2 {
		rawTag = strings.TrimRight(rawTag, "`")
		rawTag += " "
	} else {
		rawTag = "`"
	}

	return fmt.Sprintf("%s%s`", rawTag, tagConvert(fieldName, prefix, style))

}

// tagConvert returns a tag convert by
func tagConvert(fieldName, prefix, style string) string {

	if style == "" {
		style = STYLE_CAMEL_LOWER
	}

	switch style {
	case STYLE_RAW:
	case STYLE_CAMEL:
		fieldName = strcase.ToCamel(fieldName)
	case STYLE_CAMEL_LOWER:
		fieldName = strcase.ToLowerCamel(fieldName)
	case STYLE_SNAKE:
		fieldName = strcase.ToSnake(fieldName)
	case STYLE_LOWER:
		fieldName = strings.ToLower(fieldName)
	default:
		fieldName = strcase.ToLowerCamel(fieldName)
	}

	return fmt.Sprintf("%s:\"%s\"", prefix, fieldName)
}

func debugf(fmt string, args ...interface{}) {
	if !Print {
		return
	}

	if !strings.HasSuffix(fmt, "\n") {
		fmt = fmt + "\n"
	}

	log.Printf(fmt, args...)
}

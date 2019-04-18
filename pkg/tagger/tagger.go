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
)

type Tag struct {
	opt *TagOpt
}

type TagOpt struct {
	Overwrite bool
	Style     string
}

var (
	GFset *token.FileSet
	Debug = true
)

func New(opt *TagOpt) *Tag {
	if opt == nil {
		opt = &TagOpt{false, "snake"}
	}

	return &Tag{opt}
}

func (this *Tag) Tag(src []byte, output io.Writer, re, prefix string) error {

	debugf("[New Tag] re:%s prefix:%s style=%s", re, prefix, this.opt.Style)

	// 创建AST树
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

		debugf("Find struct '%s' pos:%s", s.Name.Name, GFset.Position(s.Pos()))

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

				// If nod tag
				if t.Fields.List[i].Tag == nil {
					t.Fields.List[i].Tag = new(ast.BasicLit)
				}

				debugf("Find field '%s' pos:%s \n", t.Fields.List[i].Names[j].Name, GFset.Position(t.Fields.List[i].Pos()))

				lenOldTag := len(t.Fields.List[i].Tag.Value)

				if lenOldTag < 2 {
					t.Fields.List[i].Tag.Value = fmt.Sprintf("`%s`", tagGenerate(
						t.Fields.List[i].Tag.Value,
						t.Fields.List[i].Names[j].Name,
						prefix,
						this.opt.Style,
						this.opt.Overwrite,
					))
				} else {
					t.Fields.List[i].Tag.Value = fmt.Sprintf("%s %s`", t.Fields.List[i].Tag.Value[:lenOldTag-1], tagGenerate(
						t.Fields.List[i].Tag.Value,
						t.Fields.List[i].Names[j].Name,
						prefix,
						this.opt.Style,
						this.opt.Overwrite,
					))
				}
			}

		}
		return true
	}

	ast.Inspect(f, tagFunc) // 遍历ast树并且找到我们需要修改的地方

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

// generate the new tag
func tagGenerate(rawTag, fieldName, prefix, style string, isOverwrite bool) string {

	if style == "" {
		style = "sl"
	}

	if strings.Contains(style, "l") {
		fieldName = strings.ToLower(fieldName)
	}
	if strings.Contains(style, "u") {
		fieldName = strings.ToUpper(fieldName)
	}

	switch style[0] {
	case 'c':
		return fmt.Sprintf("%s:\"%s\"", prefix, camelString(fieldName))
	case 's':
		return fmt.Sprintf("%s:\"%s\"", prefix, snakeString(fieldName))
	}

	return ""
}

// snake string, XxYy to xx_yy , XxYY to xx_y_y
func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// camel string, xx_yy to XxYy
func camelString(s string) string {
	data := make([]byte, 0, len(s))
	flag, num := true, len(s)-1
	for i := 0; i <= num; i++ {
		d := s[i]
		if d == '_' {
			flag = true
			continue
		} else if flag {
			if d >= 'a' && d <= 'z' {
				d = d - 32
			}
			flag = false
		}
		data = append(data, d)
	}
	return string(data[:])
}

func debugf(fmt string, args ...interface{}) {
	if !Debug {
		return
	}

	if !strings.HasSuffix(fmt, "\n") {
		fmt = fmt + "\n"
	}

	log.Printf(fmt, args...)
}

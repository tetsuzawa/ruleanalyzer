package ruleanalyzer

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"strings"
	"text/template"
)

func generate(tplCfg *TplConfig) ([]byte, error) {
	funcMap := template.FuncMap{
		"toLower":  strings.ToLower,
		"buildObj": buildObj,
	}

	t, err := template.New("template.go.txt").Funcs(funcMap).Parse(tpl)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, tplCfg); err != nil {
		return nil, err
	}
	return format.Source(buf.Bytes())
}

func buildObj(q MilestoneQueue) string {
	buf := &bytes.Buffer{}
	for i, v := range q {
		switch obj := v.(type) {
		case *types.Func:
			if obj.Type() != nil {
				sig := obj.Type().(*types.Signature)
				if recv := sig.Recv(); recv != nil {
					if _, ok := recv.Type().(*types.Interface); ok {
						//TODO interface
					} else {
						fmt.Fprintf(buf, "obj%d := analysisutil.MethodOf(analysisutil.TypeOf(pass, \"%s\", \"%s\"), \"%s\")\n",
							i, obj.Pkg().Name(), removePkgName(recv.Type().String()), obj.Name())
					}
				} else if obj.Pkg() != nil {
					fmt.Fprintf(buf, "obj%d := analysisutil.ObjectOf(pass, \"%s\", \"%s\")\n",
						i, obj.Pkg().Name(), obj.Name())
				}
			}
		default:
		}
	}
	return buf.String()
}

func removePkgName(s string) string {
	ss := strings.Split(s, ".")
	if len(ss) < 2 {
		return s
	}
	return ss[1]
}

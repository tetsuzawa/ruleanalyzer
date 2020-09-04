package ruleanalyzer

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/gostaticanalysis/comment"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

const rulePrefix = "Rule"

func Run() error {
	qs, err := analyze()
	if err != nil {
		return err
	}
	for _, q := range qs {
		if err := generate(q); err != nil {
			return err
		}
	}
	return nil
}

func analyze() (qs []MilestoneQueue, err error) {
	flag.Parse()
	mode := packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedImports |
		packages.NeedTypes |
		packages.NeedTypesSizes |
		packages.NeedSyntax |
		packages.NeedTypesInfo |
		packages.NeedDeps

	cfg := &packages.Config{Mode: mode}
	pkgs, err := packages.Load(cfg, flag.Args()...)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, errors.New("some errors occurred")
	}
	var maps comment.Maps
	for _, pkg := range pkgs {
		maps = comment.New(pkg.Fset, pkg.Syntax)
	}
	_, ssaPkgs := ssautil.AllPackages(pkgs, 0)

	for pkgIdx, ssaPkg := range ssaPkgs {
		ssaPkg.Build()
		for _, member := range (*ssaPkg).Members {

			// Init queue
			q := MilestoneQueue{}

			if !types.Identical(member.Type(), &types.Signature{}) {
				continue
			}
			// whether func name is `RuleXxx` or not
			if !strings.HasPrefix(member.Name(), rulePrefix) {
				continue
			}
			f, ok := member.(*ssa.Function)
			if !ok {
				continue
			}

			// analyze rule
			var ruleComments []string
			for _, b := range f.Blocks {

				commentProcessed := map[*ast.CommentGroup]bool{}
				for _, instr := range b.Instrs {

					// process the step comment. eg: // step: call xx.Xxx function
					for _, c := range maps.CommentsByPosLine(pkgs[pkgIdx].Fset, instr.Pos()) {
						if !commentProcessed[c] {
							if hasStepCheck(c.Text()) {
								ruleComments = append(ruleComments, c.Text())
							}
							commentProcessed[c] = true
							break
						}
					}
					if len(ruleComments) == len(q) {
						continue
					}

					switch instr := instr.(type) {
					case *ssa.Alloc:
						p, ok := instr.Type().(*types.Pointer)
						if !ok {
							continue
						}
						named, ok := p.Elem().(*types.Named)
						q.Push(named.Obj())
					case *ssa.Call:
						callCommon := instr.Common()
						if callCommon == nil {
							continue
						}
						var fn *types.Func
						if callCommon.Method == nil {
							callee := callCommon.StaticCallee()
							if callee == nil {
								continue
							}
							fn, ok = callee.Object().(*types.Func)
							if !ok {
								continue
							}
						} else {
							fn = callCommon.Method
						}
						q.Push(fn)
					case *ssa.Defer:
						callCommon := instr.Common()
						if callCommon == nil {
							continue
						}
						var fn *types.Func
						if callCommon.Method == nil {
							callee := callCommon.StaticCallee()
							if callee == nil {
								continue
							}
							fn, ok = callee.Object().(*types.Func)
							if !ok {
								continue
							}
						} else {
							fn = callCommon.Method
						}
						q.Push(fn)

					default:
						continue
					}
				}
			}
			if len(ruleComments) != len(q) {
				fmt.Println("length of ruleComments and q does not match. it may caused by tool compatibility")
				continue
			}
			// add milestone queue to generate analyzers
			qs = append(qs, q)
		}
	}
	return qs, nil
}

func hasStepCheck(s string) bool {
	txt := strings.Split(s, " ")
	if txt[0] != "step:" {
		return false
	}
	return true
}

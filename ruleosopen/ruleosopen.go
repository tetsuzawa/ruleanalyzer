package osopen

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"

	"github.com/gostaticanalysis/analysisutil"
	"github.com/tetsuzawa/ruleanalyzer"
)

const doc = "OsOpen detects the rule violation"

var Analyzer = &analysis.Analyzer{
	Name: "OsOpen",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

const ruleName = "OsOpen"

func run(pass *analysis.Pass) (interface{}, error) {
	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs

	obj0 := analysisutil.ObjectOf(pass, "os", "Open")
	obj1 := analysisutil.MethodOf(analysisutil.TypeOf(pass, "os", "File"), "Close")

	for _, f := range funcs {
		mq := ruleanalyzer.MilestoneQueue{

			obj0,

			obj1,
		}
		initMqLen := mq.Len()

		for _, b := range f.Blocks {
			for _, instr := range b.Instrs {
				switch instr := instr.(type) {
				case *ssa.Alloc:
					typeName, ok := mq.Head().(*types.TypeName)
					if !ok {
						continue
					}
					if ruleanalyzer.Alloc(instr, typeName) {
						mq.Pop()
						continue
					}
				case *ssa.Call:
					fn, ok := mq.Head().(*types.Func)
					if !ok {
						continue
					}
					if ruleanalyzer.Func(instr, nil, fn) {
						mq.Pop()
						continue
					}
				case *ssa.Defer:
					fn, ok := mq.Head().(*types.Func)
					if !ok {
						continue
					}
					if ruleanalyzer.Defer(instr, nil, fn) {
						mq.Pop()
						continue
					}
				default:
					continue
				}
			}
		}
		if initMqLen > mq.Len() && mq.Len() > 0 {
			pass.Reportf(f.Pos(), "this function does not match the rule: %s", ruleName)
		}
	}
	return nil, nil
}

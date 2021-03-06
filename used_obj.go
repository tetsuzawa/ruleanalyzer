package ruleanalyzer

import (
	"go/types"
	"golang.org/x/tools/go/ssa"
)

func Func(instr ssa.CallInstruction, recv ssa.Value, f *types.Func) bool {
	common := instr.Common()
	if common == nil {
		return false
	}

	callee := common.StaticCallee()
	if callee == nil {
		return false
	}

	fn, ok := callee.Object().(*types.Func)
	if !ok {
		return false
	}

	if recv != nil &&
		common.Signature().Recv() != nil &&
		(len(common.Args) == 0 && recv != nil || common.Args[0] != recv &&
			!referrer(recv, common.Args[0])) {
		return false
	}

	return fn == f
}

func referrer(a, b ssa.Value) bool {
	return isReferrerOf(a, b) || isReferrerOf(b, a)
}

func isReferrerOf(a, b ssa.Value) bool {
	if a == nil || b == nil {
		return false
	}
	if b.Referrers() != nil {
		brs := *b.Referrers()

		for _, br := range brs {
			brv, ok := br.(ssa.Value)
			if !ok {
				continue
			}
			if brv == a {
				return true
			}
		}
	}
	return false
}

func Alloc(instr *ssa.Alloc, f *types.TypeName) bool {
	p, ok := instr.Type().(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := p.Elem().(*types.Named)
	return named.Obj() == f
}

func BinOp(instr *ssa.BinOp, recv ssa.Value, f *types.Func) bool {
	//TODO
	return false
}

func Defer(instr *ssa.Defer, recv ssa.Value, f *types.Func) bool {
	callCommon := instr.Common()

	var fn *types.Func
	var ok bool
	if callCommon.Method == nil {
		callee := callCommon.StaticCallee()
		if callee == nil {
			return false
		}
		fn, ok = callee.Object().(*types.Func)
		if !ok {
			return false
		}
	} else {
		fn = callCommon.Method
	}

	if recv != nil &&
		callCommon.Signature().Recv() != nil &&
		(len(callCommon.Args) == 0 && recv != nil || callCommon.Args[0] != recv &&
			!referrer(recv, callCommon.Args[0])) {
		return false
	}

	return fn == f
}

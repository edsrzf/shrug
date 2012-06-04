package main

var builtins = map[string]builtinCmd{
	"if": ifCmd,
	"set": setCmd,
}

func ifCmd(args []val, ctx *context) int {
	if len(args) < 1 {
		ctx.stderr.Write([]byte("if: usage: set cond [ iftrue ] [ iffalse ]"))
		return 2
	}
	cond, ok := args[0].(cmd)
	if !ok {
	}
	ret := cond.exec(nil, ctx)
	if ret == 0 {
		if len(args) >= 2 {
			iftrue, ok := args[1].(cmd)
			if ok {
				iftrue.exec(nil, ctx)
			}
		}
		return 0
	} else if len(args) >= 3 {
		iffalse, ok := args[2].(cmd)
		if ok {
			iffalse.exec(nil, ctx)
		}
	}
	return 1
}

func setCmd(args []val, ctx *context) int {
	if len(args) < 2 {
		ctx.stderr.Write([]byte("set: usage: set variable value"))
		return 1
	}
	varname, ok := args[0].(word)
	if !ok {
		ctx.stderr.Write([]byte("set: invalid variable name"))
		return 1
	}
	ctx.set(string(varname), args[1])
	return 0
}

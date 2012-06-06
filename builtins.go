package main

var builtins = []*builtinCmd {
	{"for", forCmd},
	{"if", ifCmd},
	{"result", resultCmd},
	{"set", setCmd},
}

func ifCmd(args []val, ctx *context) val {
	if len(args) < 1 {
		ctx.stderr.Write([]byte("if: usage: if cond [ iftrue ] [ iffalse ]"))
		return intVal(1)
	}
	cond, ok := args[0].(cmd)
	if !ok {
	}
	ret := cond.exec(nil, ctx)
	if ret.bool() {
		if len(args) >= 2 {
			iftrue, ok := args[1].(cmd)
			if ok {
				iftrue.exec(nil, ctx)
			}
		}
	} else if len(args) >= 3 {
		iffalse, ok := args[2].(cmd)
		if ok {
			iffalse.exec(nil, ctx)
		}
	}
	return nilVal{}
}

func forCmd(args []val, ctx *context) val {
	if len(args) < 2 {
		ctx.stderr.Write([]byte("for: usage: for variable [ list ... ] body"))
		return intVal(1)
	}
	varname, ok := args[0].(word)
	if !ok {
		ctx.stderr.Write([]byte("for: invalid variable name"))
		return intVal(1)
	}
	body, ok := args[len(args)-1].(cmd)
	if !ok {
		ctx.stderr.Write([]byte("for: invalid body"))
		return intVal(1)
	}
	for _, item := range args[1:len(args)-1] {
		ctx.set(string(varname), item)
		body.exec(nil, ctx)
	}
	return nilVal{}
}

func resultCmd(args []val, ctx *context) val {
	switch len(args) {
	case 0:
		return nilVal{}
	case 1:
		return args[0]
	}
	return list(args)
}

func setCmd(args []val, ctx *context) val {
	if len(args) < 2 {
		ctx.stderr.Write([]byte("set: usage: set variable value"))
		return intVal(1)
	}
	varname, ok := args[0].(word)
	if !ok {
		ctx.stderr.Write([]byte("set: invalid variable name"))
		return intVal(1)
	}
	ctx.set(string(varname), args[1])
	return nilVal{}
}

package main

var builtins = map[string]builtinCmd{
	"set": setCmd,
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

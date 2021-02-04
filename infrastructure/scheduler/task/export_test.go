package task

func (c *Cmd) CmdFunc() interface{} {
	return c.f.Interface()
}

func (c *Cmd) Args() []interface{} {
	args := make([]interface{}, len(c.args))
	for i, arg := range c.args {
		args[i] = arg.Interface()
	}
	return args
}

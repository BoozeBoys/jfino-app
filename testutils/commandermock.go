package testutils

type CommandDescription struct {
	Name string
	Args []interface{}
}

type CommanderMock struct {
	Commands []CommandDescription
	errors   map[string]error
}

func (mc *CommanderMock) Power(on bool) error {
	descr := CommandDescription{
		Name: "power",
		Args: []interface{}{on},
	}
	mc.Commands = append(mc.Commands, descr)

	if mc.errors != nil && mc.errors["power"] != nil {
		return mc.errors["power"]
	}

	return nil
}

func (mc *CommanderMock) Speed(speed1, speed2 int) error {
	descr := CommandDescription{
		Name: "speed",
		Args: []interface{}{speed1, speed2},
	}
	mc.Commands = append(mc.Commands, descr)

	if mc.errors != nil && mc.errors["speed"] != nil {
		return mc.errors["speed"]
	}

	return nil
}

func (mc *CommanderMock) SetError(command string, err error) {
	if mc.errors == nil {
		mc.errors = make(map[string]error)
	}
	mc.errors[command] = err
}

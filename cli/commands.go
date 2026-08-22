package cli

import "fmt"

type Commands struct {
	commandNames map[string]func(*State, Command) error
}

func (c *Commands) run(s *State, cmd Command) error {
	v, exists := c.commandNames[cmd.name]

	if !exists {
		return fmt.Errorf("Error, command %s does not exist", cmd.name)
	}

	return v(s, cmd)
}

func (c *Commands) register(name string, f func(*State, Command) error) error {
	_, exist := c.commandNames[name]

	if exist {
		return fmt.Errorf("Error, command %s already exists.", name)
	}

	c.commandNames[name] = f

	return nil
}

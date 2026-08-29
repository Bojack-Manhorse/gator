package cli

import "fmt"

type Commands struct {
	CommandMap map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	v, exists := c.CommandMap[cmd.Name]

	if !exists {
		return fmt.Errorf("Error, command %s does not exist", cmd.Name)
	}

	return v(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) error {
	_, exist := c.CommandMap[name]

	if exist {
		return fmt.Errorf("Error, command %s already exists.", name)
	}

	c.CommandMap[name] = f

	return nil
}

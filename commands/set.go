package commands

import "fmt"

type Set map[string]Command

func (s Set) Execute(command string, args []string) error {
	cmd, ok := s[command]
	if !ok {
		return fmt.Errorf("unknown command: %s", command)
	}

	err := cmd.Execute(args)
	if err != nil {
		return fmt.Errorf("could not execute %q: %s", command, err)
	}

	return nil
}

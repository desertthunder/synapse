package main

import (
	"flag"
)

type FlagType string

type Command struct {
	Name    string
	flagSet *flag.FlagSet
	Logger  *Logger
	flags   map[string]Flag
}

type Flag struct {
	Name         string
	UsageText    string
	DefaultValue interface{}
	Type         FlagType
}

const (
	StringFlag FlagType = "string"
	BoolFlag   FlagType = "bool"
	IntFlag    FlagType = "int"
	FloatFlag  FlagType = "float"
)

func (f Flag) WithType(t FlagType) Flag {
	f.Type = t
	return f
}

func NewCommand(name string, logger *Logger) Command {
	return Command{
		Name:    name,
		flagSet: flag.NewFlagSet(name, flag.ExitOnError),
		Logger:  logger,
		flags:   make(map[string]Flag, 0),
	}
}

func (f Flag) Parse(val any) any {
	switch f.Type {
	case BoolFlag:
		return val.(bool)
	case FloatFlag:
		return val.(float64)
	case IntFlag:
		return val.(int)
	default:
		return val.(string)
	}
}

func (c *Command) AddFlag(f Flag) {
	switch f.DefaultValue.(type) {
	case string:
		c.flags[f.Name] = f.WithType(StringFlag)
	case bool:
		c.flags[f.Name] = f.WithType(BoolFlag)
	case int, int8, int16, int32, int64:
		c.flags[f.Name] = f.WithType(IntFlag)
	case float32, float64:
		c.flags[f.Name] = f.WithType(FloatFlag)
	default:
		c.Logger.Warn("type not supported")
	}
}

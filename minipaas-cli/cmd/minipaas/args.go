package main

type BaseArgs struct {
	Env     string `arg:"positional,required" help:"Directory for MiniPaaS environment to use"`
	Verbose bool   `arg:"-v,--verbose" help:"Verbose output" default:"false"`
}

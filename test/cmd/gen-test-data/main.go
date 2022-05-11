package main

import (
	"github.com/alexflint/go-arg"
)

const InitialDataN = 10_000_000
const RequestsN = 10_000_000
const UpperBoundRand = 20

type cmdArgs struct {
	Init       *cmdInit       `arg:"subcommand"`
	Ammo       *cmdAmmo       `arg:"subcommand"`
	ServerInit *cmdServerInit `arg:"subcommand"`
}

func main() {
	var args cmdArgs
	arg.MustParse(&args)

	switch {
	case args.Ammo != nil:
		genAmmoData(args.Ammo)
	case args.Init != nil:
		getInitial(args.Init)
	case args.ServerInit != nil:
		serverInit(args.ServerInit)
	}
}

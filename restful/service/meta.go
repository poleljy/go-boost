package service

import (
	"bufio"
	"flag"
	"io"
	"os"

	"github.com/mitchellh/cli"
)

// Meta contains the meta-options and functionality that nearly every
// command inherits.
type Meta struct {
	Ui               cli.Ui
	ServiceName      string
	ServiceNamespace string

	// These are set by the command line flags.
	ServiceAddress string
}

func DefaultMeta(svcName string) *Meta {
	meta := new(Meta)
	meta.ServiceName = svcName

	if meta.Ui == nil {
		meta.Ui = &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		}
	}
	return meta
}

// FlagSetFlags is an enum to define what flags are present in the
// default FlagSet returned by Meta.FlagSet.
type FlagSetFlags uint

const (
	FlagSetNone    FlagSetFlags = 0
	FlagSetClient  FlagSetFlags = 1 << iota
	FlagSetDefault              = FlagSetClient
)

func (m *Meta) FlagSet(n string, fs FlagSetFlags) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)
	if fs&FlagSetClient != 0 {
		f.StringVar(&m.ServiceAddress, "server_address", ":8080", "server address, env: SERVER_ADDRESS")
	}

	errR, errW := io.Pipe()
	errScanner := bufio.NewScanner(errR)
	go func() {
		for errScanner.Scan() {
			m.Ui.Error(errScanner.Text())
		}
	}()
	f.SetOutput(errW)

	return f
}

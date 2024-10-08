package main

import (
	"os"

	"github.com/jagottsicher/termcolor"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	inputPolicyFile  string
	outputPolicyFile string
	jcAPIKey         string
	logger           *pterm.Logger
	noColor          bool
	stripEmailDomain bool

	cliCmd = &cobra.Command{
		Use:   "jc2headscale",
		Short: "jc2headscale - fills groups in policy based on Jumpcloud user groups",
		Long: `Collects information about Jumpcloud groups, group members
               and prepare a group list for Headscale policy.`,
	}
)

func init() {
	// Add command persistent flag
	cliCmd.PersistentFlags().StringVar(&inputPolicyFile, "input-policy", "./policy.hjson", "Headscale policy file template")
	cliCmd.PersistentFlags().StringVar(&outputPolicyFile, "output-policy", "./current.json", "Headscale prepared policy file")
	cliCmd.PersistentFlags().StringVar(&jcAPIKey, "jc-api-key", os.Getenv("JC_API_KEY"), "The Jumpcloud API key (can use env var JC_API_KEY)")
	cliCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")
	cliCmd.PersistentFlags().BoolVar(&stripEmailDomain, "strip-email-domain", true, "Strip e-mail domain")

	// Disable colors if terminal doesn't support or user set flag --no-color
	if !checkTerminalColorSupport() || noColor {
		pterm.DisableColor()
	}

	logger = pterm.DefaultLogger.WithLevel(pterm.LogLevelInfo)
}

func checkTerminalColorSupport() bool {
	var termColorSupport bool
	switch l := termcolor.SupportLevel(os.Stderr); l {
	case termcolor.Level16M:
		termColorSupport = true
	case termcolor.Level256:
		termColorSupport = true
	case termcolor.LevelBasic:
		termColorSupport = true
	case termcolor.LevelNone:
		termColorSupport = false
	default:
		termColorSupport = false
	}

	return termColorSupport
}

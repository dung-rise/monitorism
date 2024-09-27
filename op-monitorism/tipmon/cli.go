package tipmon

import (
	opservice "github.com/ethereum-optimism/optimism/op-service"

	"github.com/urfave/cli/v2"
)

const (
	NodeURLFlagName = "node.url"
)

type CLIConfig struct {
	NodeUrl string
}

func ReadCLIFlags(ctx *cli.Context) (CLIConfig, error) {
	cfg := CLIConfig{NodeUrl: ctx.String(NodeURLFlagName)}

	return cfg, nil
}

func CLIFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    NodeURLFlagName,
			Usage:   "Node URL of a peer",
			Value:   "127.0.0.1:8545",
			EnvVars: opservice.PrefixEnvVar(envPrefix, "NODE_URL"),
		},
	}
}

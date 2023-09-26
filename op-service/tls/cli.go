// This file contains CLI and env TLS configurations that can be used by clients or servers
package tls

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	opservice "github.com/ethereum-optimism/optimism/op-service"
)

const (
	TLSCaCertFlagName = "tls.ca"
	TLSCertFlagName   = "tls.cert"
	TLSKeyFlagName    = "tls.key"
)

// CLIFlags returns flags with env var envPrefix
// This should be used for server TLS configs, or when client and server tls configs are the same
func CLIFlags(envPrefix string) []cli.Flag {
	return CLIFlagsWithFlagPrefix(envPrefix, "")
}

var (
	defaultTLSCaCert = "tls/ca.crt"
	defaultTLSCert   = "tls/tls.crt"
	defaultTLSKey    = "tls/tls.key"
)

// CLIFlagsWithFlagPrefix returns flags with env var and cli flag prefixes
// Should be used for client TLS configs when different from server on the same process
func CLIFlagsWithFlagPrefix(envPrefix string, flagPrefix string) []cli.Flag {
	prefixFunc := func(flagName string) string {
		return strings.Trim(fmt.Sprintf("%s.%s", flagPrefix, flagName), ".")
	}
	prefixEnvVars := func(name string) []string {
		return opservice.PrefixEnvVar(envPrefix, name)
	}
	return []cli.Flag{
		&cli.StringFlag{
			Name:    prefixFunc(TLSCaCertFlagName),
			Usage:   "tls ca cert path",
			Value:   defaultTLSCaCert,
			EnvVars: prefixEnvVars("TLS_CA"),
		},
		&cli.StringFlag{
			Name:    prefixFunc(TLSCertFlagName),
			Usage:   "tls cert path",
			Value:   defaultTLSCert,
			EnvVars: prefixEnvVars("TLS_CERT"),
		},
		&cli.StringFlag{
			Name:    prefixFunc(TLSKeyFlagName),
			Usage:   "tls key",
			Value:   defaultTLSKey,
			EnvVars: prefixEnvVars("TLS_KEY"),
		},
	}
}

type CLIConfig struct {
	TLSCaCert string
	TLSCert   string
	TLSKey    string
}

func NewCLIConfig() CLIConfig {
	return CLIConfig{
		TLSCaCert: defaultTLSCaCert,
		TLSCert:   defaultTLSCert,
		TLSKey:    defaultTLSKey,
	}
}

func (c CLIConfig) Check() error {
	if c.TLSEnabled() && (c.TLSCaCert == "" || c.TLSCert == "" || c.TLSKey == "") {
		return errors.New("all tls flags must be set if at least one is set")
	}

	return nil
}

func (c CLIConfig) TLSEnabled() bool {
	return !(c.TLSCaCert == "" && c.TLSCert == "" && c.TLSKey == "")
}

// ReadCLIConfig reads tls cli configs
// This should be used for server TLS configs, or when client and server tls configs are the same
func ReadCLIConfig(ctx *cli.Context) CLIConfig {
	return CLIConfig{
		TLSCaCert: ctx.String(TLSCaCertFlagName),
		TLSCert:   ctx.String(TLSCertFlagName),
		TLSKey:    ctx.String(TLSKeyFlagName),
	}
}

// ReadCLIConfigWithPrefix reads tls cli configs with flag prefix
// Should be used for client TLS configs when different from server on the same process
func ReadCLIConfigWithPrefix(ctx *cli.Context, flagPrefix string) CLIConfig {
	prefixFunc := func(flagName string) string {
		return strings.Trim(fmt.Sprintf("%s.%s", flagPrefix, flagName), ".")
	}
	return CLIConfig{
		TLSCaCert: ctx.String(prefixFunc(TLSCaCertFlagName)),
		TLSCert:   ctx.String(prefixFunc(TLSCertFlagName)),
		TLSKey:    ctx.String(prefixFunc(TLSKeyFlagName)),
	}
}

const (
	EndpointFlagName = "signer.endpoint"
	AddressFlagName  = "signer.address"
)

func SignerCLIFlags(envPrefix string) []cli.Flag {
	envPrefix += "_SIGNER"
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    EndpointFlagName,
			Usage:   "Signer endpoint the client will connect to",
			EnvVars: opservice.PrefixEnvVar(envPrefix, "ENDPOINT"),
		},
		&cli.StringFlag{
			Name:    AddressFlagName,
			Usage:   "Address the signer is signing transactions for",
			EnvVars: opservice.PrefixEnvVar(envPrefix, "ADDRESS"),
		},
	}
	flags = append(flags, CLIFlagsWithFlagPrefix(envPrefix, "signer")...)
	return flags
}

type SignerCLIConfig struct {
	Endpoint  string
	Address   string
	TLSConfig CLIConfig
}

func NewSignerCLIConfig() SignerCLIConfig {
	return SignerCLIConfig{
		TLSConfig: NewCLIConfig(),
	}
}

func (c SignerCLIConfig) Check() error {
	if err := c.TLSConfig.Check(); err != nil {
		return err
	}
	if !((c.Endpoint == "" && c.Address == "") || (c.Endpoint != "" && c.Address != "")) {
		return errors.New("signer endpoint and address must both be set or not set")
	}
	return nil
}

func (c SignerCLIConfig) Enabled() bool {
	if c.Endpoint != "" && c.Address != "" {
		return true
	}
	return false
}

func ReadSignerCLIConfig(ctx *cli.Context) SignerCLIConfig {
	cfg := SignerCLIConfig{
		Endpoint:  ctx.String(EndpointFlagName),
		Address:   ctx.String(AddressFlagName),
		TLSConfig: ReadCLIConfigWithPrefix(ctx, "signer"),
	}
	return cfg
}

package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/skycoin/dmsg/cipher"
	coinCipher "github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/spf13/cobra"

	"github.com/skycoin/skywire/pkg/visor/visorconfig"
)

func init() {
	RootCmd.AddCommand(updateConfigCmd)
}

var (
	addOutput              string
	addInput               string
	environment            string
	addHypervisorPKs       string
	resetHypervisor        bool
	setVPNClientKillswitch string
	addVPNClientSrv        string
	addVPNClientPasscode   string
	resetVPNclient         bool
	addVPNServerPasscode   string
	setVPNServerSecure     string
	resetVPNServer         bool
	addSkysocksClientSrv   string
	resetSkysocksClient    bool
	skysocksPasscode       string
	resetSkysocks          bool
	setPublicAutoconnect   string
	minhops                int
)

func init() {
	updateConfigCmd.Flags().StringVarP(&addOutput, "output", "o", "skywire-config.json", "path of output config file.")
	updateConfigCmd.Flags().StringVarP(&addInput, "input", "i", "skywire-config.json", "path of input config file.")
	updateConfigCmd.Flags().StringVarP(&environment, "environment", "e", "production", "desired environment (values production or testing)")
	updateConfigCmd.Flags().StringVar(&addHypervisorPKs, "add-hypervisor-pks", "", "public keys of hypervisors that should be added to this visor")
	updateConfigCmd.Flags().BoolVar(&resetHypervisor, "reset-hypervisor-pks", false, "resets hypervisor`s configuration")

	updateConfigCmd.Flags().StringVar(&setVPNClientKillswitch, "vpn-client-killswitch", "", "public keys of hypervisors that should be added to this visor")
	updateConfigCmd.Flags().StringVar(&addVPNClientSrv, "add-vpn-client-server", "", "public keys of hypervisors that should be added to this visor")
	updateConfigCmd.Flags().StringVar(&addVPNClientPasscode, "add-vpn-client-passcode", "", "public keys of hypervisors that should be added to this visor")
	updateConfigCmd.Flags().BoolVar(&resetVPNclient, "reset-vpn-client", false, "public keys of hypervisors that should be added to this visor")

	updateConfigCmd.Flags().StringVar(&addVPNServerPasscode, "add-vpn-server-passcode", "", "public keys of hypervisors that should be added to this visor")
	updateConfigCmd.Flags().StringVar(&setVPNServerSecure, "vpn-server-secure", "", "public keys of hypervisors that should be added to this visor")
	updateConfigCmd.Flags().BoolVar(&resetVPNServer, "reset-vpn-server", false, "public keys of hypervisors that should be added to this visor")

	updateConfigCmd.Flags().StringVar(&addSkysocksClientSrv, "add-skysocks-client-server", "", "public keys of hypervisors that should be added to this visor")
	updateConfigCmd.Flags().BoolVar(&resetSkysocksClient, "reset-skysocks-client", false, "public keys of hypervisors that should be added to this visor")

	updateConfigCmd.Flags().StringVar(&skysocksPasscode, "add-skysocks-passcode", "", "public keys of hypervisors that should be added to this visor")
	updateConfigCmd.Flags().BoolVar(&resetSkysocks, "reset-skysocks", false, "public keys of hypervisors that should be added to this visor")

	updateConfigCmd.Flags().StringVar(&setPublicAutoconnect, "set-public-autoconnect", "", "public keys of hypervisors that should be added to this visor")

	updateConfigCmd.Flags().IntVar(&minhops, "set-minhop", -1, "public keys of hypervisors that should be added to this visor")
}

var updateConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates a config file",
	PreRun: func(_ *cobra.Command, _ []string) {
		var err error
		if output, err = filepath.Abs(addOutput); err != nil {
			logger.WithError(err).Fatal("Invalid config output.")
		}
	},
	Run: func(_ *cobra.Command, _ []string) {
		mLog := logging.NewMasterLogger()
		mLog.SetLevel(logrus.InfoLevel)
		f, err := os.Open(addInput) // nolint: gosec
		if err != nil {
			mLog.WithError(err).
				WithField("filepath", addInput).
				Fatal("Failed to read config file.")
		}

		raw, err := ioutil.ReadAll(f)
		if err != nil {
			mLog.WithError(err).Fatal("Failed to read config.")
		}

		conf, ok := visorconfig.Parse(mLog, addInput, raw)
		if ok != nil {
			mLog.WithError(err).Fatal("Failed to parse config.")
		}

		if addHypervisorPKs != "" {
			keys := strings.Split(addHypervisorPKs, ",")
			for _, key := range keys {
				keyParsed, err := coinCipher.PubKeyFromHex(strings.TrimSpace(key))
				if err != nil {
					logger.WithError(err).Fatalf("Failed to parse hypervisor private key: %s.", key)
				}
				conf.Hypervisors = append(conf.Hypervisors, cipher.PubKey(keyParsed))
			}
		}

		switch environment {
		case "production":
			visorconfig.SetDefaultProductionValues(conf)
		case "testing":
			visorconfig.SetDefaultTestingValues(conf)
		default:
			logger.Fatal("Unrecognized environment value: ", environment)
		}

		if resetHypervisor {
			conf.Hypervisors = []cipher.PubKey{}
		}

		switch setVPNClientKillswitch {
		case "true":
			changeAppsConfig(conf, "vpn-client", "--killswitch", setVPNClientKillswitch)
		case "false":
			changeAppsConfig(conf, "vpn-client", "--killswitch", setVPNClientKillswitch)
		}

		if addVPNClientSrv != "" {
			keyParsed, err := coinCipher.PubKeyFromHex(strings.TrimSpace(addVPNClientSrv))
			if err != nil {
				logger.WithError(err).Fatalf("Failed to parse hypervisor private key: %s.", addVPNClientSrv)
			}
			changeAppsConfig(conf, "vpn-client", "--srv", keyParsed.Hex())
		}

		if addVPNClientPasscode != "" {
			changeAppsConfig(conf, "vpn-client", "--passcode", addVPNClientPasscode)
		}

		if resetVPNclient {
			resetAppsConfig(conf, "vpn-client")
		}

		if addVPNServerPasscode != "" {
			changeAppsConfig(conf, "vpn-server", "--passcode", addVPNClientPasscode)
		}

		switch setVPNServerSecure {
		case "true":
			changeAppsConfig(conf, "vpn-server", "--secure", setVPNClientKillswitch)
		case "false":
			changeAppsConfig(conf, "vpn-server", "--secure", setVPNClientKillswitch)
		}

		if resetVPNServer {
			resetAppsConfig(conf, "vpn-server")
		}

		if addSkysocksClientSrv != "" {
			keyParsed, err := coinCipher.PubKeyFromHex(strings.TrimSpace(addSkysocksClientSrv))
			if err != nil {
				logger.WithError(err).Fatalf("Failed to parse hypervisor private key: %s.", addSkysocksClientSrv)
			}
			changeAppsConfig(conf, "skysocks-client", "--srv", keyParsed.Hex())
		}

		if resetSkysocksClient {
			resetAppsConfig(conf, "skysocks-client")
		}

		if skysocksPasscode != "" {
			changeAppsConfig(conf, "skysocks", "--passcode", skysocksPasscode)
		}

		if resetSkysocks {
			resetAppsConfig(conf, "skysocks")
		}

		switch setPublicAutoconnect {
		case "true":
			conf.Transport.PublicAutoconnect = true
		case "false":
			conf.Transport.PublicAutoconnect = false
		}

		if minhops >= 0 {
			conf.Routing.MinHops = uint16(minhops)
		}

		// Save config to file.
		if err := conf.Flush(); err != nil {
			logger.WithError(err).Fatal("Failed to flush config to file.")
		}

		// Print results.
		j, err := json.MarshalIndent(conf, "", "\t")
		if err != nil {
			logger.WithError(err).Fatal("An unexpected error occurred. Please contact a developer.")
		}
		logger.Infof("Updated file '%s' to: %s", output, j)
	},
}

func changeAppsConfig(conf *visorconfig.V1, appName string, argName string, argValue string) {
	apps := conf.Launcher.Apps
	for index := range apps {
		if apps[index].Name != appName {
			continue
		}
		updated := false
		for index, arg := range apps[index].Args {
			if arg == argName {
				apps[index].Args[index+1] = argValue
				updated = true
			}
		}
		if !updated {
			apps[index].Args = append(apps[index].Args, argName, argValue)
		}
	}
}

func resetAppsConfig(conf *visorconfig.V1, appName string) {
	apps := conf.Launcher.Apps
	for index := range apps {
		if apps[index].Name == appName {
			apps[index].Args = []string{}
		}
	}
}

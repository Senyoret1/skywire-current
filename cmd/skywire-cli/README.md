# CLI Documentation

skywire command line interface

<!-- MarkdownTOC autolink="true" bracket="round" levels="1,2,3" -->

- [Install](#install)
- [skywire-cli usage](#skywire-cli-usage)
    - [config usage](#config-usage)
        - [gen](#config-gen)
        - [update](#config-update)
    - [visor usage](#visor-usage)
        - [exec](#exec)
        - [pk](#pk)
        - [hv](#hv)
        - [info](#summary)
        - [version](#version)
        - [app](#app)
            - [ls](#app-ls)
            - [autostart](#app-autostart)
            - [start](#app-start)
            - [stop](#app-stop)
        - [route](#route)
            - [ls rules](#route-ls-rules)
            - [rule](#route-rule)
            - [add rule](#route-add-rule)
            - [rm rule](#route-rm-rule)
        - [tp](#tp)
            - [type](#tp-type)
            - [disc](#tp-disc)
            - [id](#tp-id)
            - [ls](#tp-ls)
            - [add](#tp-add)
            - [rm](#tp-rm)
        - [vpn](#vpn)
            - [ui](#vpn-ui)
            - [url](#vpn-url)
        - [update](#visor-update)
    - [rtfind usage](#rtfind-usage)
    - [mdisc usage](#mdisc-usage)
        - [servers](#servers)
        - [entry](#entry)
    - [completion usage](#completion-usage)


<!-- /MarkdownTOC -->

## Install

```bash
$ cd $GOPATH/src/github.com/skycoin/skywire/cmd/skywire-cli
$ go install ./...
```

## skywire-cli usage

After the installation, you can run `skywire-cli` to see the usage:

```
$ skywire-cli

	┌─┐┬┌─┬ ┬┬ ┬┬┬─┐┌─┐  ┌─┐┬  ┬
	└─┐├┴┐└┬┘││││├┬┘├┤───│  │  │
	└─┘┴ ┴ ┴ └┴┘┴┴└─└─┘  └─┘┴─┘┴

Usage:
  skywire-cli [command]

Available Commands:
  config       Generate or update a skywire config
  visor        Query the Skywire Visor
  rtfind       Query the Route Finder
  mdisc        Query remote DMSG Discovery
  completion   Generate completion script
  help         Help about any command

Flags:
  -h, --help   help for skywire-cli

Use "skywire-cli [command] --help" for more information about a command.
```

### config usage

```
skywire-cli config -h
Generate or update a skywire config

Usage:
  skywire-cli config [command]

Available Commands:
  gen         generate a config file
  update      update a config file

Flags:
  -h, --help   help for config

Use "skywire-cli config [command] --help" for more information about a command.
```

#### config gen

```
$ skywire-cli config gen --help
generate a config file

Usage:
  skywire-cli config gen [flags]

Flags:
  -b, --bestproto      best protocol (dmsg | direct) based on location
  -i, --ishv           local hypervisor configuration
  -j, --hvpks string   list of public keys to use as hypervisor
  -o, --out string     output config default:skywire-config.json
  -p, --package        use paths for package /opt/skywire
  -r, --regen          re-generate existing config & retain keys
      --all            show all flags
  -h, --help           help for gen

$ skywire-cli config gen --all
generate a config file

Usage:
  skywire-cli config gen [flags]

Flags:
  -a, --url string           services conf (default "conf.skywire.skycoin.com")
  -b, --bestproto            best protocol (dmsg | direct) based on location
  -c, --noauth               disable authentication for hypervisor UI
  -d, --dmsghttp             use dmsg connection to skywire services
  -e, --auth                 enable auth on hypervisor UI
  -f, --force                remove pre-existing config
  -g, --disableapps string   comma separated list of apps to disable
  -i, --ishv                 local hypervisor configuration
  -j, --hvpks string         list of public keys to use as hypervisor
  -k, --os string            (linux / mac / win) paths (default "linux")
  -n, --stdout               write config to stdout
  -o, --out string           output config: skywire-config.json
  -p, --pkg                  use path for package: /opt/skywire
  -u, --user                 use paths for user space: /home/user
  -q, --publicrpc            allow rpc requests from LAN
  -r, --regen                re-generate existing config & retain keys
  -s, --sk cipher.SecKey     a random key is generated if unspecified
 (default 0000000000000000000000000000000000000000000000000000000000000000)
  -t, --testenv              use test deployment conf.skywire.dev
  -v, --servevpn             enable vpn server
  -w, --hide                 dont print the config to the terminal
  -x, --retainhv             retain existing hypervisors with regen
      --version string       custom version testing override
      --binpath string       set bin_path
  -h, --help                 help for gen
```

##### Example defaults

The default visor config generation assumes the command is run from the root of the cloned repository.
The example further assumes the compiled binary is available in the executable PATH and that GOPATH is set

<details>

<summary>
cd $GOPATH/src/github.com/skycoin/skywire && skywire-cli config gen
</summary>

```
$ cd $GOPATH/src/github.com/skycoin/skywire
$ skywire-cli config gen
[user@linux skywire]$ go run cmd/skywire-cli/skywire-cli.go config gen
[2022-05-26T17:17:01-05:00] INFO []: Fetched service endpoints from 'http://conf.skywire.skycoin.com'
[2022-05-26T17:17:01-05:00] INFO [visor:config]: Flushing config to file. config_version="v1.0.0" filepath="/home/user/go/src/github.com/skycoin/skywire/skywire-config.json"
[2022-05-26T17:17:01-05:00] INFO [skywire-cli]: Updated file 'skywire-config.json' to:
{
	"version": "v1.0.0",
	"sk": "7df293067dc1942e5e1484a945120da4729819166ad8e02341b13dfd4d5469bd",
	"pk": "03caa7c5ea048c8200538ff81f53ffad08cc289f280f2cfce9c261ea8b03405097",
	"dmsg": {
		"discovery": "http://dmsgd.skywire.skycoin.com",
		"sessions_count": 1,
		"servers": []
	},
	"dmsgpty": {
		"dmsg_port": 22,
		"cli_network": "unix",
		"cli_address": "/tmp/dmsgpty.sock"
	},
	"skywire-tcp": {
		"pk_table": null,
		"listening_address": ":7777"
	},
	"transport": {
		"discovery": "http://tpd.skywire.skycoin.com",
		"address_resolver": "http://ar.skywire.skycoin.com",
		"public_autoconnect": true,
		"transport_setup_nodes": null
	},
	"routing": {
		"setup_nodes": [
			"0324579f003e6b4048bae2def4365e634d8e0e3054a20fc7af49daf2a179658557"
		],
		"route_finder": "http://rf.skywire.skycoin.com",
		"route_finder_timeout": "10s",
		"min_hops": 0
	},
	"uptime_tracker": {
		"addr": "http://ut.skywire.skycoin.com"
	},
	"launcher": {
		"service_discovery": "http://sd.skycoin.com",
		"apps": [
			{
				"name": "vpn-client",
				"auto_start": false,
				"port": 43
			},
			{
				"name": "skychat",
				"args": [
					"-addr",
					":8001"
				],
				"auto_start": true,
				"port": 1
			},
			{
				"name": "skysocks",
				"auto_start": true,
				"port": 3
			},
			{
				"name": "skysocks-client",
				"auto_start": false,
				"port": 13
			},
			{
				"name": "vpn-server",
				"auto_start": false,
				"port": 44
			}
		],
		"server_addr": "localhost:5505",
		"bin_path": "./apps"
	},
	"hypervisors": [],
	"cli_addr": "localhost:3435",
	"log_level": "info",
	"local_path": "./local",
	"stun_servers": [
		"172.104.188.139:3478",
		"172.104.59.235:3478",
		"172.104.183.187:3478",
		"139.162.54.63:3478",
		"172.105.115.97:3478",
		"172.104.188.39:3478",
		"172.104.188.140:3478",
		"172.104.40.88:3478"
	],
	"shutdown_timeout": "10s",
	"restart_check_delay": "1s",
	"is_public": false,
	"persistent_transports": null
}

```
</details>

##### Example hypervisor defaults

The default configuration is for a visor only. To generate a configuration which provides the hypervisor web interface,
the `-i` or `--is-hypervisor` flag should be specified.

<details>

<summary>
skywire-cli config gen -i
</summary>

```
$ skywire-cli config gen -i
[2022-05-26T17:21:48-05:00] INFO []: Fetched service endpoints from 'http://conf.skywire.skycoin.com'
[2022-05-26T17:21:48-05:00] INFO [visor:config]: Flushing config to file. config_version="v1.0.0" filepath="/home/user/go/src/github.com/skycoin/skywire/skywire-config.json"
[2022-05-26T17:21:48-05:00] INFO [skywire-cli]: Updated file 'skywire-config.json' to:
{
	"version": "v1.0.0",
	"sk": "053a98a9e69a80d7359297991f46a6d1b57c20f4d86a2f6dac60fbd9ef342a24",
	"pk": "02aebbac3f4c0cb2a5606ab6d6d63ea88dd8e53255d3413c733c5588ff5acf2305",
	"dmsg": {
		"discovery": "http://dmsgd.skywire.skycoin.com",
		"sessions_count": 1,
		"servers": []
	},
	"dmsgpty": {
		"dmsg_port": 22,
		"cli_network": "unix",
		"cli_address": "/tmp/dmsgpty.sock"
	},
	"skywire-tcp": {
		"pk_table": null,
		"listening_address": ":7777"
	},
	"transport": {
		"discovery": "http://tpd.skywire.skycoin.com",
		"address_resolver": "http://ar.skywire.skycoin.com",
		"public_autoconnect": true,
		"transport_setup_nodes": null
	},
	"routing": {
		"setup_nodes": [
			"0324579f003e6b4048bae2def4365e634d8e0e3054a20fc7af49daf2a179658557"
		],
		"route_finder": "http://rf.skywire.skycoin.com",
		"route_finder_timeout": "10s",
		"min_hops": 0
	},
	"uptime_tracker": {
		"addr": "http://ut.skywire.skycoin.com"
	},
	"launcher": {
		"service_discovery": "http://sd.skycoin.com",
		"apps": [
			{
				"name": "vpn-client",
				"auto_start": false,
				"port": 43
			},
			{
				"name": "skychat",
				"args": [
					"-addr",
					":8001"
				],
				"auto_start": true,
				"port": 1
			},
			{
				"name": "skysocks",
				"auto_start": true,
				"port": 3
			},
			{
				"name": "skysocks-client",
				"auto_start": false,
				"port": 13
			},
			{
				"name": "vpn-server",
				"auto_start": false,
				"port": 44
			}
		],
		"server_addr": "localhost:5505",
		"bin_path": "./apps"
	},
	"hypervisors": [],
	"cli_addr": "localhost:3435",
	"log_level": "info",
	"local_path": "./local",
	"stun_servers": [
		"172.104.188.139:3478",
		"172.104.59.235:3478",
		"172.104.183.187:3478",
		"139.162.54.63:3478",
		"172.105.115.97:3478",
		"172.104.188.39:3478",
		"172.104.188.140:3478",
		"172.104.40.88:3478"
	],
	"shutdown_timeout": "10s",
	"restart_check_delay": "1s",
	"is_public": false,
	"persistent_transports": null,
	"hypervisor": {
		"db_path": "/home/user/go/src/github.com/skycoin/skywire/users.db",
		"enable_auth": false,
		"cookies": {
			"hash_key": "80b9c82cff3280592ca91e9f3d276f4b7c900551a24ef4f311567d21ad19d63dbdda75d3cbcce24863124151ba6d07846def16cb300cd2bd87117e802c743ac4",
			"block_key": "d113be28b2a4b256dafde1598651b0f402b27912fb38b8449ba0b9a6c947775c",
			"expires_duration": 43200000000000,
			"path": "/",
			"domain": ""
		},
		"dmsg_port": 46,
		"http_addr": ":8000",
		"enable_tls": false,
		"tls_cert_file": "./ssl/cert.pem",
		"tls_key_file": "./ssl/key.pem"
	}
}
```
</details>


Note that it is possible to start the visor with the hypervisor interface explicitly; regardless of how the config was generated; using the -i flag

```
skywire-visor -i
```

##### Example dmsghttp defaults

Using dmsghttp routes http traffic used for connecting to the skywire services through dmsg.

The dmsghttp-config.json file must be present to generate a config with dmsghttp

The `-b` or `--bestproto` flag will automatically determine if dmsghttp should be used based on region

The `-d` or `--dmsghttp` flag explicitly creates the config with dmsghttp

It is recommended to use the `-b` flag for config file generation.

The example below uses `-d` to create a dmsghttp config

<details>

<summary>
skywire-cli config gen -d
</summary>

```
[2022-05-26T17:30:06-05:00] INFO [skywire-cli]: Found Dmsghttp config: dmsghttp-config.json
[2022-05-26T17:30:07-05:00] INFO []: Fetched service endpoints from 'http://conf.skywire.skycoin.com'
[2022-05-26T17:30:07-05:00] INFO [visor:config]: Flushing config to file. config_version="v1.0.0" filepath="/home/user/go/src/github.com/skycoin/skywire/skywire-config.json"
[2022-05-26T17:30:07-05:00] INFO [skywire-cli]: Updated file 'skywire-config.json' to:
{
	"version": "v1.0.0",
	"sk": "4aaa3cc266cb8c4ec9689a0d88c49802646517ae2d2240315c480d461c6af70e",
	"pk": "0291b61ab823eebe79575d6d0f8e122e84cd17dae1cab5c7dfb9043f1ee4f0a206",
	"dmsg": {
		"discovery": "dmsg://022e607e0914d6e7ccda7587f95790c09e126bbd506cc476a1eda852325aadd1aa:80",
		"sessions_count": 1,
		"servers": [
			{
				"version": "",
				"sequence": 0,
				"timestamp": 0,
				"static": "02a2d4c346dabd165fd555dfdba4a7f4d18786fe7e055e562397cd5102bdd7f8dd",
				"server": {
					"address": "dmsg.server02a2d4c3.skywire.skycoin.com:30081",
					"availableSessions": 0
				}
			},
			{
				"version": "",
				"sequence": 0,
				"timestamp": 0,
				"static": "03717576ada5b1744e395c66c2bb11cea73b0e23d0dcd54422139b1a7f12e962c4",
				"server": {
					"address": "dmsg.server03717576.skywire.skycoin.com:30082",
					"availableSessions": 0
				}
			},
			{
				"version": "",
				"sequence": 0,
				"timestamp": 0,
				"static": "0228af3fd99c8d86a882495c8e0202bdd4da78c69e013065d8634286dd4a0ac098",
				"server": {
					"address": "45.118.133.242:30084",
					"availableSessions": 0
				}
			},
			{
				"version": "",
				"sequence": 0,
				"timestamp": 0,
				"static": "03d5b55d1133b26485c664cf8b95cff6746d1e321c34e48c9fed293eff0d6d49e5",
				"server": {
					"address": "dmsg.server03d5b55d.skywire.skycoin.com:30083",
					"availableSessions": 0
				}
			},
			{
				"version": "",
				"sequence": 0,
				"timestamp": 0,
				"static": "0281a102c82820e811368c8d028cf11b1a985043b726b1bcdb8fce89b27384b2cb",
				"server": {
					"address": "192.53.114.142:30085",
					"availableSessions": 0
				}
			},
			{
				"version": "",
				"sequence": 0,
				"timestamp": 0,
				"static": "02a49bc0aa1b5b78f638e9189be4ed095bac5d6839c828465a8350f80ac07629c0",
				"server": {
					"address": "dmsg.server02a4.skywire.skycoin.com:30089",
					"availableSessions": 0
				}
			}
		]
	},
	"dmsgpty": {
		"dmsg_port": 22,
		"cli_network": "unix",
		"cli_address": "/tmp/dmsgpty.sock"
	},
	"skywire-tcp": {
		"pk_table": null,
		"listening_address": ":7777"
	},
	"transport": {
		"discovery": "dmsg://02b307aee5c8ce1666c63891f8af25ad2f0a47a243914c963942b3ba35b9d095ae:80",
		"address_resolver": "dmsg://03234b2ee4128d1f78c180d06911102906c80795dfe41bd6253f2619c8b6252a02:80",
		"public_autoconnect": true,
		"transport_setup_nodes": null
	},
	"routing": {
		"setup_nodes": [
			"0324579f003e6b4048bae2def4365e634d8e0e3054a20fc7af49daf2a179658557"
		],
		"route_finder": "dmsg://039d89c5eedfda4a28b0c58b0b643eff949f08e4f68c8357278081d26f5a592d74:80",
		"route_finder_timeout": "10s",
		"min_hops": 0
	},
	"uptime_tracker": {
		"addr": "dmsg://022c424caa6239ba7d1d9d8f7dab56cd5ec6ae2ea9ad97bb94ad4b48f62a540d3f:80"
	},
	"launcher": {
		"service_discovery": "dmsg://0204890f9def4f9a5448c2e824c6a4afc85fd1f877322320898fafdf407cc6fef7:80",
		"apps": [
			{
				"name": "vpn-client",
				"auto_start": false,
				"port": 43
			},
			{
				"name": "skychat",
				"args": [
					"-addr",
					":8001"
				],
				"auto_start": true,
				"port": 1
			},
			{
				"name": "skysocks",
				"auto_start": true,
				"port": 3
			},
			{
				"name": "skysocks-client",
				"auto_start": false,
				"port": 13
			},
			{
				"name": "vpn-server",
				"auto_start": false,
				"port": 44
			}
		],
		"server_addr": "localhost:5505",
		"bin_path": "./apps"
	},
	"hypervisors": [],
	"cli_addr": "localhost:3435",
	"log_level": "info",
	"local_path": "./local",
	"stun_servers": [
		"172.104.188.139:3478",
		"172.104.59.235:3478",
		"172.104.183.187:3478",
		"139.162.54.63:3478",
		"172.105.115.97:3478",
		"172.104.188.39:3478",
		"172.104.188.140:3478",
		"172.104.40.88:3478"
	],
	"shutdown_timeout": "10s",
	"restart_check_delay": "1s",
	"is_public": false,
	"persistent_transports": null
}
```
</details>


##### Example package based installation defaults

This assumes the skywire linux installation is at `/opt/skywire` with binaries and apps in their own subdirectories.
The `-p` flag default paths are provided by the skywire linux / mac packages or windows .msi installer and generate the skywire config within the install dir.


<details>

<summary>
sudo skywire-cli config gen -bipr
</summary>

```
$ sudo skywire-cli config gen -bipr
[sudo] password for user:
[2022-05-26T17:43:42-05:00] INFO []: Fetched service endpoints from 'http://conf.skywire.skycoin.com'
[2022-05-26T17:43:42-05:00] INFO [visor:config]: Flushing config to file. config_version="1.0.0" filepath="/opt/skywire/skywire.json"
[2022-05-26T17:43:42-05:00] INFO [skywire-cli]: Updated file '/opt/skywire/skywire.json' to:
{
	"version": "1.0.0",
	"sk": "4aaa3cc266cb8c4ec9689a0d88c49802646517ae2d2240315c480d461c6af70e",
	"pk": "0291b61ab823eebe79575d6d0f8e122e84cd17dae1cab5c7dfb9043f1ee4f0a206",
	"dmsg": {
		"discovery": "http://dmsgd.skywire.skycoin.com",
		"sessions_count": 1,
		"servers": []
	},
	"dmsgpty": {
		"dmsg_port": 22,
		"cli_network": "unix",
		"cli_address": "/tmp/dmsgpty.sock"
	},
	"skywire-tcp": {
		"pk_table": null,
		"listening_address": ":7777"
	},
	"transport": {
		"discovery": "http://tpd.skywire.skycoin.com",
		"address_resolver": "http://ar.skywire.skycoin.com",
		"public_autoconnect": true,
		"transport_setup_nodes": null
	},
	"routing": {
		"setup_nodes": [
			"0324579f003e6b4048bae2def4365e634d8e0e3054a20fc7af49daf2a179658557"
		],
		"route_finder": "http://rf.skywire.skycoin.com",
		"route_finder_timeout": "10s",
		"min_hops": 0
	},
	"uptime_tracker": {
		"addr": "http://ut.skywire.skycoin.com"
	},
	"launcher": {
		"service_discovery": "http://sd.skycoin.com",
		"apps": [
			{
				"name": "vpn-client",
				"auto_start": false,
				"port": 43
			},
			{
				"name": "skychat",
				"args": [
					"-addr",
					":8001"
				],
				"auto_start": true,
				"port": 1
			},
			{
				"name": "skysocks",
				"auto_start": true,
				"port": 3
			},
			{
				"name": "skysocks-client",
				"auto_start": false,
				"port": 13
			},
			{
				"name": "vpn-server",
				"auto_start": false,
				"port": 44
			}
		],
		"server_addr": "localhost:5505",
		"bin_path": "/opt/skywire/apps"
	},
	"hypervisors": [],
	"cli_addr": "localhost:3435",
	"log_level": "info",
	"local_path": "/opt/skywire/local",
	"stun_servers": [
		"172.104.188.139:3478",
		"172.104.59.235:3478",
		"172.104.183.187:3478",
		"139.162.54.63:3478",
		"172.105.115.97:3478",
		"172.104.188.39:3478",
		"172.104.188.140:3478",
		"172.104.40.88:3478"
	],
	"shutdown_timeout": "10s",
	"restart_check_delay": "1s",
	"is_public": false,
	"persistent_transports": null,
	"hypervisor": {
		"db_path": "/opt/skywire/users.db",
		"enable_auth": true,
		"cookies": {
			"hash_key": "80b9c82cff3280592ca91e9f3d276f4b7c900551a24ef4f311567d21ad19d63dbdda75d3cbcce24863124151ba6d07846def16cb300cd2bd87117e802c743ac4",
			"block_key": "d113be28b2a4b256dafde1598651b0f402b27912fb38b8449ba0b9a6c947775c",
			"expires_duration": 43200000000000,
			"path": "/",
			"domain": ""
		},
		"dmsg_port": 46,
		"http_addr": ":8000",
		"enable_tls": false,
		"tls_cert_file": "./ssl/cert.pem",
		"tls_key_file": "./ssl/key.pem"
	}
}
```
</details>

The configuration is written (or rewritten)

##### Example remote hypervisor configuration for package based installation

The typical arrangement uses a remote hypervisor if a local instance is not started.

The hypervisor public key can be determined by running the following command on the running hypervisor:

```
$ skywire-cli visor pk
```

configure the visor to use the public key of the remote hypervisor:

```
# skywire-cli config gen -bprj <hypervisor-public-key>
```

The configuration is regenerated

#### config update

```
$ skywire-cli config update --help
update a config file

Usage:
  skywire-cli config update [flags]
  skywire-cli config update [command]

Available Commands:
  hv           update hypervisor config
  sc           update skysocks-client config
  ss           update skysocks-server config
  vpnc         update vpn-client config
  vpns         update vpn-server config

Flags:
  -a, --endpoints                update server endpoints
  -b, --url string               service config URL: conf.skywire.skycoin.com
  -t, --testenv                  use test deployment: conf.skywire.dev
      --public-autoconn string   change public autoconnect configuration
      --set-minhop int           change min hops value (default -1)
  -i, --input string             path of input config file.
  -o, --output string            config file to output
  -u, --user                     update config at: $HOME/skywire-config.json
  -h, --help                     help for update

Use "skywire-cli config update [command] --help" for more information about a command.
```

##### Example

<details>

<summary>
skywire-cli config update
</summary>


```
$ skywire-cli config update
[2022-05-26T17:58:16-05:00] INFO [visor:config]: Flushing config to file. config_version="v1.0.0" filepath="/home/user/go/src/github.com/skycoin/skywire/skywire-config.json"
[2022-05-26T17:58:16-05:00] INFO [skywire-cli]: Updated file '/home/user/go/src/github.com/skycoin/skywire/skywire-config.json' to: {
	"version": "v1.0.0",
	"sk": "754ee10a20f8bc824c8c2ac9c6b26ab10112e0a339681bc759ccf87d714195a0",
	"pk": "02bc601c45633f98da260946936b409ee098609f1192ff2b88986308f100edb721",
	"dmsg": {
		"discovery": "http://dmsgd.skywire.skycoin.com",
		"sessions_count": 1,
		"servers": []
	},
	"dmsgpty": {
		"dmsg_port": 22,
		"cli_network": "unix",
		"cli_address": "/tmp/dmsgpty.sock"
	},
	"skywire-tcp": {
		"pk_table": null,
		"listening_address": ":7777"
	},
	"transport": {
		"discovery": "http://tpd.skywire.skycoin.com",
		"address_resolver": "http://ar.skywire.skycoin.com",
		"public_autoconnect": true,
		"transport_setup_nodes": null
	},
	"routing": {
		"setup_nodes": [
			"0324579f003e6b4048bae2def4365e634d8e0e3054a20fc7af49daf2a179658557"
		],
		"route_finder": "http://rf.skywire.skycoin.com",
		"route_finder_timeout": "10s",
		"min_hops": 0
	},
	"uptime_tracker": {
		"addr": "http://ut.skywire.skycoin.com"
	},
	"launcher": {
		"service_discovery": "http://sd.skycoin.com",
		"apps": [
			{
				"name": "vpn-client",
				"auto_start": false,
				"port": 43
			},
			{
				"name": "skychat",
				"args": [
					"-addr",
					":8001"
				],
				"auto_start": true,
				"port": 1
			},
			{
				"name": "skysocks",
				"auto_start": true,
				"port": 3
			},
			{
				"name": "skysocks-client",
				"auto_start": false,
				"port": 13
			},
			{
				"name": "vpn-server",
				"auto_start": false,
				"port": 44
			}
		],
		"server_addr": "localhost:5505",
		"bin_path": "./apps"
	},
	"hypervisors": [],
	"cli_addr": "localhost:3435",
	"log_level": "info",
	"local_path": "./local",
	"stun_servers": [
		"172.104.188.139:3478",
		"172.104.59.235:3478",
		"172.104.183.187:3478",
		"139.162.54.63:3478",
		"172.105.115.97:3478",
		"172.104.188.39:3478",
		"172.104.188.140:3478",
		"172.104.40.88:3478"
	],
	"shutdown_timeout": "10s",
	"restart_check_delay": "1s",
	"is_public": false,
	"persistent_transports": null
}
```
</details>

### visor usage

```
$ skywire-cli visor -h
Query the Skywire Visor

Usage:
  skywire-cli visor [command]

Available Commands:
  exec         execute a command
  pk           Public key of the visor
  hvpk         Public key of hypervisor this visor is using
  info         summary of visor info
  version      version and build info
  app          app settings
  route        view and set rules
  tp           view and set transports
  vpn          vpn interface
  update       update the local visor

Flags:
  -h, --help         help for visor
      --rpc string   RPC server address (default "localhost:3435")

Use "skywire-cli visor [command] --help" for more information about a command.
```

#### exec

execute a given command

```
$ skywire-cli visor exec <command>
```

##### Example

ls

```
$ skywire-cli visor exec ls
bin
boot
dev
efi
etc
home
lib
lib64
media
mnt
opt
proc
root
run
sbin
share
srv
sys
tmp
usr
var
```

echo

```
$ skywire-cli visor exec echo "hello world"
hello world
```

escape a flag

```
$skywire-cli visor exec echo -- "-a"
-a
```

#### pk

public key of the visor

```
$ skywire-cli visor pk
```

```
Flags:
  -i, --input string   path of input config file.
```

##### Example

```
$ skywire-cli visor pk                                                          
0359f02198933550ad5b41a21470a0bbe0f73c0eb6e93d7d279133a0d5bffc645c   
```

#### hv

show hypervisor(s)

```
$ skywire-cli visor hvpk
```

```
Flags:
  -i, --input string   path of input config file.
```

##### Example

```
$ skywire-cli visor hvpk
[0359f02198933550ad5b41a21470a0bbe0f73c0eb6e93d7d279133a0d5bffc645c]
```

#### info

summary of visor info

```
$ skywire-cli visor info
```

##### Example

```
$ skywire-cli visor info
.:: Visor Summary ::.
Public key: "03caa7c5ea048c8200538ff81f53ffad08cc289f280f2cfce9c261ea8b03405097"
Symmetric NAT: false
IP: 192.168.0.2
DMSG Server: "0371ab4bcff7b121f4b91f6856d6740c6f9dc1fe716977850aeb5d84378b300a13"
Ping: "447.072621ms"
Visor Version: 1.0.0
Skybian Version:
Uptime Tracker: healthy
Time Online: 6082.307143 seconds
Build Tag:
```


#### version

version and build info

```
$ skywire-cli visor version
```

##### Example

```
$ skywire-cli visor version
Version "v1.0.0" built on "2022-05-26T18:18:39Z" against commit "668d5ad8"
```

#### app

```
$ skywire-cli visor app
app settings

Usage:
  skywire-cli visor app [command]

Available Commands:
  ls          list apps
  start       launch app
  stop        halt app
  autostart   set autostart flag for app
  log         logs from app since RFC3339Nano-formated timestamp.
                    "beginning" is a special timestamp to fetch all the logs

Flags:
  -h, --help   help for app

Global Flags:
      --rpc string   RPC server address (default "localhost:3435")

Use "skywire-cli visor app [command] --help" for more information about a command.
```

##### app ls

list apps

```
$ skywire-cli visor app ls
```

##### Example

```
$ skywire-cli visor app ls
app                 ports     auto_start     status
skychat             1         true           running
skysocks            3         true           running
skysocks-client     13        false          stopped
vpn-server          44        false          stopped
vpn-client          43        false          stopped
```

#### app start

start application

```
$ skywire-cli visor app start <name>
```

##### Example

```
$ skywire-cli visor app start vpn-server
OK
```

#### app stop

stop application

```
$ skywire-cli visor app stop <name>
```

##### Example

```
$ skywire-cli visor app stop skychat
OK
```


#### app autostart

set autostart flag for app

```
$ skywire-cli visor app autostart <name> (on|off)
```

##### Example

```
$ skywire-cli visor app autostart vpn-server on
OK
```

#### app logs

logs from app since RFC3339Nano-formated timestamp.
                    "beginning" is a special timestamp to fetch all the logs

```
$ skywire-cli visor app logs <name> <timestamp>
```

##### Example

```
$ skywire-cli visor app log skysocks beginning
 [2022-03-11T21:15:55-06:00] INFO [public_autoconnect]: Fetching public visors
 [2022-03-11T21:16:06-06:00] INFO [public_autoconnect]: Fetching public visors
 [2022-03-11T21:16:09-06:00] INFO [dmsgC]: Session stopped. error="failed to serve dialed session to 0371ab4bcff7b121f4b91f6856d6740c6f9dc1fe716977850aeb5d84378b300a13: EOF"
 [2022-03-11T21:16:09-06:00] WARN [dmsgC]: Stopped accepting streams. error="EOF" session=0371ab4bcff7b121f4b91f6856d6740c6f9dc1fe716977850aeb5d84378b300a13
 [2022-03-11T21:16:10-06:00] INFO [dmsgC]: Dialing session... remote_pk=0281a102c82820e811368c8d028cf11b1a985043b726b1bcdb8fce89b27384b2cb
 [2022-03-11T21:16:14-06:00] INFO [dmsgC]: Serving session. remote_pk=0281a102c82820e811368c8d028cf11b1a985043b726b1bcdb8fce89b27384b2cb
```

#### route

```
$ skywire-cli visor route
view and set rules

Usage:
  skywire-cli visor route [command]

Available Commands:
  ls-rules    list routing rules
  rule        return routing rule by route ID key
  rm-rule     remove routing rule
  add-rule    add routing rule

Flags:
  -h, --help   help for route

Global Flags:
      --rpc string   RPC server address (default "localhost:3435")

Use "skywire-cli visor route [command] --help" for more information about a command.
```

#### route add-rule

```
$ skywire-cli visor route add-rule (app <route-id> <local-pk> <local-port> <remote-pk> <remote-port> | fwd <next-route-id> <next-transport-id>) [flags]
```

##### Example

```
$ skywire-cli visor route add-rule -h
add routing rule

Usage:
  skywire-cli visor route add-rule (app <route-id> <local-pk> <local-port> <remote-pk> <remote-port> | fwd <next-route-id> <next-transport-id>) [flags]

Flags:
  -h, --help                  help for add-rule
      --keep-alive duration   duration after which routing rule will expire if no activity is present (default 30s)

```

#### route rm-rule

Removes a routing rule

```
$ skywire-cli visor route rm-rule <route-id>
```

##### Example

```
$ skywire-cli visor route rm-rule -h
Removes a routing rule via route ID key

Usage:
  skywire-cli visor rm-rule <route-id> [flags]

```

#### route ls-rules

list routing rules

```
$ skywire-cli visor route ls-rules
```

#### route rule

```
$ skywire-cli visor route rule <route-id>
```

##### Example

```
$ skywire-cli visor route rule -h
Returns a routing rule via route ID key

Usage:
  skywire-cli visor route rule <route-id> [flags]

```

#### tp

```
view and set transports

Usage:
  skywire-cli visor tp [command]

Available Commands:
  disc        discover transport(s) by ID or public key
  type        transport types used by the local visor
  ls          available transports
  id          transport summary by id
  add         add a transport
  rm          remove transport(s) by id

Flags:
  -h, --help   help for tp

Global Flags:
      --rpc string   RPC server address (default "localhost:3435")

Use "skywire-cli visor tp [command] --help" for more information about a command.

```

#### tp add

add transport

```
$ skywire-cli visor tp add <remote-public-key> [flags]
```

##### Example

```
$ skywire-cli visor tp add -h
Adds a new transport

Usage:
  skywire-cli visor add-tp <remote-public-key> [flags]

Flags:
  -h, --help               help for add-tp
      --public             whether to make the transport public (deprecated)
  -t, --timeout duration   if specified, sets an operation timeout
      --type string        type of transport to add; if unspecified, cli will attempt to establish a transport in the following order: stcp, stcpr, sudph, dmsg
```

#### tp disc

discover transport(s) by ID or public key

```
$ skywire-cli visor tp disc (--id=<transport-id> | --pk=<edge-public-key>)
```

##### Example

```
$ skywire-cli visor tp disc -h
discover transport(s) by ID or public key

Usage:
  skywire-cli visor tp disc (--id=<transport-id> | --pk=<edge-public-key>) [flags]

Flags:
  -h, --help               help for disc-tp
      --id transportID     if specified, obtains a single transport of given ID (default 00000000-0000-0000-0000-000000000000)
      --pk cipher.PubKey   if specified, obtains transports associated with given public key (default 000000000000000000000000000000000000000000000000000000000000000000)

```

#### tp id

transport summary by id

```
$ skywire-cli visor tp id <transport-id>
```

##### Example

```
$ skywire-cli visor tp id -h
transport summary by id

Usage:
skywire-cli visor tp <transport-id> [flags]
```

#### tp ls

list transports

```
$ skywire-cli visor tp ls
```

##### Example

```
$ skywire-cli visor tp ls
type     id     remote     mode     is_up
```

#### tp type

Lists transport types used by the local visor

```
$ skywire-cli visor tp type
```

##### Example

```
$ skywire-cli visor tp type
dmsg
stcp
stcpr
sudph
```

#### tp rm

remove transport

```
$ skywire-cli visor tp rm <transport-id>
```

##### Example

```
$ skywire-cli visor tp rm -h
Removes transport with given id

Usage:
  skywire-cli visor tp rm <transport-id> [flags]

```


#### vpn

vpn interface


```
$   skywire-cli visor vpn [command]
```

##### Example

```
$ skywire-cli visor vpn -h

vpn interface

Usage:
  skywire-cli visor vpn [command]

Available Commands:
  ui          Open VPN UI in default browser
  url         Show VPN UI URL

```

#### vpn ui

Open VPN UI in default browser

```
$   skywire-cli visor vpn ui
```


##### Example

```
$ skywire-cli visor vpn ui
```

the VPN user interface is opened in the default browser

#### vpn url

Show VPN UI URL

```
$   skywire-cli visor vpn url
```


##### Example

```
$ skywire-cli visor vpn url
http://127.0.0.1:8000/#/vpn/027087fe40d97f7f0be4a0dc768462ddbb371d4b9e7679d4f11f117d757b9856ed/
```


#### update

update

```
$ skywire-cli visor update
```

### rtfind usage

```
skywire-cli rtfind <public-key-visor-1> <public-key-visor-2>
```

##### Example

```
$ skywire-cli rtfind -h

Query the Route Finder

Usage:
  skywire-cli rtfind <public-key-visor-1> <public-key-visor-2> [flags]

Flags:
  -a, --addr string        route finder service address (default "http://rf.skywire.skycoin.com")
  -h, --help               help for rtfind
  -x, --max-hops uint16    maximum hops (default 1000)
  -n, --min-hops uint16    minimum hops (default 1)
  -t, --timeout duration   request timeout (default 10s)
```

### mdisc usage

```
Query remote DMSG Discovery

Usage:
  skywire-cli mdisc [command]

Available Commands:
  entry       fetch an entry
  servers     fetch available servers

Flags:
      --addr string   address of DMSG discovery server
                       (default "http://dmsgd.skywire.skycoin.com")
  -h, --help          help for mdisc

  Use "skywire-cli mdisc [command] --help" for more information about a command.
```

#### servers

```
$ skywire-cli mdisc servers
```

```
Flags:
      --addr string   address of DMSG discovery server
      (default "http://dmsgd.skywire.skycoin.com")
```

##### Example

```
$ skywire-cli mdisc server
[2022-03-13T21:10:44-05:00] DEBUG disc.NewHTTP [mdisc:disc]: Created HTTP client. addr="http://dmsgd.skywire.skycoin.com"
version     registered              public-key                                                             address                                           available-sessions
0.0.1       1647224020460616235     02347729662a901d03f1a1ab6c189a173349fa11e79fe82117cca0f8d0e4d64a31     192.53.115.181:8082                               2582
0.0.1       1647224015059832662     02e4660279c83bc6ca0122d3a78c0cb3f3564e03e04876ae7fa30b4e0a63217425     192.53.115.181:8081                               1299
0.0.1       1647224018690620887     02a2d4c346dabd165fd555dfdba4a7f4d18786fe7e055e562397cd5102bdd7f8dd     dmsg.server02a2d4c3.skywire.skycoin.com:30081     1109
0.0.1       1647224019967944735     0371ab4bcff7b121f4b91f6856d6740c6f9dc1fe716977850aeb5d84378b300a13     192.53.114.142:30086                              582
0.0.1       1647224016544544252     0228af3fd99c8d86a882495c8e0202bdd4da78c69e013065d8634286dd4a0ac098     45.118.133.242:30084                              48
0.0.1       1647224021047139719     03717576ada5b1744e395c66c2bb11cea73b0e23d0dcd54422139b1a7f12e962c4     dmsg.server03717576.skywire.skycoin.com:30082     31
0.0.1       1647224018229901714     0281a102c82820e811368c8d028cf11b1a985043b726b1bcdb8fce89b27384b2cb     192.53.114.142:30085                              19
0.0.1       1647224017051283856     02a49bc0aa1b5b78f638e9189be4ed095bac5d6839c828465a8350f80ac07629c0     dmsg.server02a4.skywire.skycoin.com:30089         1

```

#### entry

```
$ skywire-cli mdisc entry <visor-public-key>
```

```
Flags:
      --addr string   address of DMSG discovery server
```

##### Example

```
$ skywire-cli mdisc entry 034b68c4d8ec6d934d3ecb28595fea7e89a8de2048f0f857759c5018cb8e2f9525
[2022-03-13T21:17:11-05:00] DEBUG disc.NewHTTP [mdisc:disc]: Created HTTP client. addr="http://dmsgd.skywire.skycoin.com"
	version: 0.0.1
	sequence: 4
	registered at: 1647205336195743639
	static public key: 034b68c4d8ec6d934d3ecb28595fea7e89a8de2048f0f857759c5018cb8e2f9525
	signature: 7a7cee456a17b13207a8eba6dd60102505e0d5b3b98f047225da8bfc8e963a557c75fbbba5c7654835230c9372d6faae2f7570bb71b1af9d36cbdc4da195b74701
	entry is registered as client. Related info:
		delegated servers:
			0371ab4bcff7b121f4b91f6856d6740c6f9dc1fe716977850aeb5d84378b300a13
```


### completion usage

```
#skywire-cli completion
```

```
To load completions

Usage:
  skywire-cli completion [bash|zsh|fish|powershell]

Flags:
  -h, --help   help for completion
```

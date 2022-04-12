/*
proxy server app for skywire visor
*/
package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	cc "github.com/ivanpirog/coloredcobra"
	ipc "github.com/james-barrow/golang-ipc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/skycoin/skywire-utilities/pkg/buildinfo"
	"github.com/skycoin/skywire/internal/skysocks"
	"github.com/skycoin/skywire/pkg/app"
	"github.com/skycoin/skywire/pkg/app/appnet"
	"github.com/skycoin/skywire/pkg/routing"
	"github.com/skycoin/skywire/pkg/skyenv"
)

const (
	netType              = appnet.TypeSkynet
	port    routing.Port = 3
)

var (
	log      = logrus.New()
	passcode string
)

func init() {
	rootCmd.Flags().SortFlags = false

	rootCmd.Flags().StringVarP(&passcode, "passcode", "r", "", "Passcode to authenticate connection")
}

var rootCmd = &cobra.Command{
	Use:   "skysocks",
	Short: "Skywire SOCKS5 Proxy Server",
	Long: `
	┌─┐┬┌─┬ ┬┌─┐┌─┐┌─┐┬┌─┌─┐
	└─┐├┴┐└┬┘└─┐│ ││  ├┴┐└─┐
	└─┘┴ ┴ ┴ └─┘└─┘└─┘┴ ┴└─┘`,
	Run: func(_ *cobra.Command, _ []string) {
		appC := app.NewClient(nil)
		defer appC.Close()

		skysocks.Log = log

		if _, err := buildinfo.Get().WriteTo(os.Stdout); err != nil {
			fmt.Printf("Failed to output build info: %v", err)
		}

		srv, err := skysocks.NewServer(passcode, log)
		if err != nil {
			log.Fatal("Failed to create a new server: ", err)
		}

		l, err := appC.Listen(netType, port)
		if err != nil {
			log.Fatalf("Error listening network %v on port %d: %v\n", netType, port, err)
		}

		fmt.Println("Starting serving proxy server")

		if runtime.GOOS == "windows" {
			ipcClient, err := ipc.StartClient(skyenv.VPNClientName, nil)
			if err != nil {
				fmt.Printf("Error creating ipc server for VPN client: %v\n", err)
				os.Exit(1)
			}
			go srv.ListenIPC(ipcClient)
		} else {
			termCh := make(chan os.Signal, 1)
			signal.Notify(termCh, os.Interrupt)

			go func() {
				<-termCh

				if err := srv.Close(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}()

		}

		if err := srv.Serve(l); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Execute executes root CLI command.
func Execute() {
	cc.Init(&cc.Config{
		RootCmd:         rootCmd,
		Headings:        cc.HiBlue + cc.Bold,
		Commands:        cc.HiBlue + cc.Bold,
		CmdShortDescr:   cc.HiBlue,
		Example:         cc.HiBlue + cc.Italic,
		ExecName:        cc.HiBlue + cc.Bold,
		Flags:           cc.HiBlue + cc.Bold,
		FlagsDescr:      cc.HiBlue,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	Execute()
}

/*
Copyright © 2021 Mario Fritschen <mario@fritschen.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Mario-F/cert-manager-selfservice/internal/cleaner"
	"github.com/Mario-F/cert-manager-selfservice/internal/kube"
	"github.com/Mario-F/cert-manager-selfservice/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	serverPort    int
	metricsPort   int
	managerId     string
	kubeNamespace string
	issuerKind    string
	issuerName    string
	cleanupHours  int64
	noCleanup     bool
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start webserver to provide certificates from cert-manager",
	Run: func(cmd *cobra.Command, args []string) {
		kube.SetNamespace(kubeNamespace)
		kube.SetManagerId(managerId)

		// TODO: check for kubernetes resources and error out with proper error message

		cleaner := cleaner.Cleaner{}
		if !noCleanup {
			err := cleaner.Start(cleanupHours)
			if err != nil {
				log.Errorf("Failed to start cleaner: %+v", err)
				return
			}
		}

		go server.StartMetricsExporter(metricsPort)
		go server.Start(serverPort, issuerKind, issuerName)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		signal.Notify(quit, syscall.SIGTERM)
		<-quit
		log.Println("Sutting down gracefully.")
		cleaner.Stop()
		server.Stop()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serverCmd.Flags().IntVar(&serverPort, "port", 8030, "Port for the webserver to use")
	serverCmd.Flags().IntVar(&metricsPort, "metrics-port", 8040, "Port for exporting prometheus metrics")
	serverCmd.Flags().StringVar(&managerId, "manager-id", "default", "A uniq id to mark certificates managed by this instance")
	serverCmd.Flags().StringVar(&kubeNamespace, "kube-namespace", "default", "Kubernetes namespace to use")
	serverCmd.Flags().StringVar(&issuerKind, "issuer-kind", "ClusterIssuer", "Cert Manager issuer to use")
	serverCmd.Flags().StringVar(&issuerName, "issuer-name", "", "Cert Manager issuer instance to use")
	serverCmd.Flags().Int64Var(&cleanupHours, "cleanup-hours", 72, "Cleanup certificates not accessed after hours")
	serverCmd.Flags().BoolVar(&noCleanup, "no-cleanup", false, "Disable cleanup of unused certificates")
	err := serverCmd.MarkFlagRequired("issuer-name")
	if err != nil {
		log.Error(err)
	}
}

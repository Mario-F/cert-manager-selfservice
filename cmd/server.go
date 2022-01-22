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
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start webserver to provide certificates from cert-manager",
	Run: func(cmd *cobra.Command, args []string) {
		kube.SetNamespace(kubeNamespace)
		kube.SetManagerId(managerId)

		go server.StartMetricsExporter(metricsPort)

		cleaner := cleaner.Cleaner{}
		cleaner.Start(cleanupHours)

		server.Start(serverPort, issuerKind, issuerName)
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
	err := serverCmd.MarkFlagRequired("issuer-name")
	if err != nil {
		log.Error(err)
	}
}

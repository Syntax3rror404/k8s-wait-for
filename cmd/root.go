/*
Copyright © 2025 Marcel Zapf
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	namespace string
	label     string
	timer     int32
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "waitfor",
	Short: "This tool waits for kubernetes pods or jobs and SQL databases to be ready",
	Long: `This tool waits for kubernetes pods or jobs and SQL databases to be ready

a common usecase is to use it as init container to wait for other pods to be ready before starting the main application.
For example waiting for a database to be ready before starting the app to prevent errors.

Example:
  waitfor pod -n vault -l app.kubernetes.io/instance=vault
  waitfor job -n snipeit -l job=generate-app-key
  waitfor sql -u root -p mysecretpassword -s mariadb.mydatabase.cluster.local -d mydb
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Flags
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Namespace to use")
	rootCmd.PersistentFlags().StringVarP(&label, "label", "l", "", "Label to filter (required)")
	rootCmd.PersistentFlags().Int32VarP(&timer, "timer", "t", 3, "Wait time between checks")
}

// Create and return a Kubernetes clientset
// 1. If running inside Kubernetes, use in-cluster config
// 2. Otherwise, fall back to ~/.kube/config
func getKubeClient() (*kubernetes.Clientset, error) {
	// Try to load cluster config
	config, err := rest.InClusterConfig()
	fmt.Println("Info: Trying in-cluster config...")
	if err != nil {
		fmt.Println("Info: Error, falling back to kubeconfig...")
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			home, _ := os.UserHomeDir()
			kubeconfig = home + "/.kube/config"
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	// Create Kubernetes API client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset, nil
}

func validateFlags(cmd *cobra.Command) error {
	nsFlag := cmd.Flags().Lookup("namespace")
	labelFlag := cmd.Flags().Lookup("label")

	if nsFlag != nil && nsFlag.Changed {
		fmt.Printf("Info: Namespace set by user: %q\n", namespace)
	} else {
		fmt.Printf("Info: No namespace set by user... Using default namespace: %q\n", namespace)
	}

	if labelFlag != nil && labelFlag.Changed {
		fmt.Printf("Info: Label set by user: %q\n", label)
	} else {
		return fmt.Errorf("you must provide a label")
	}

	return nil
}

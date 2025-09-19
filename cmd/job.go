/*
Copyright © 2025 Marcel Zapf
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// jobCmd represents the job command
var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "wait for a job to complete",
	Long:  `wait for a pod to be ready`,
	Run: func(cmd *cobra.Command, args []string) {

		// 1. Check flags
		if err := validateFlags(cmd); err != nil {
			fmt.Println("Error:", err)
			return
		}

		// 2. Create Kubernetes client
		clientset, err := getKubeClient()
		if err != nil {
			fmt.Printf("Error: Failed to create Kubernetes client: %v\n", err)
			return
		}

		// 3. Wait for jobs to be completed
		waitForJobsCompleted(clientset, namespace, label)
	},
}

func init() {
	rootCmd.AddCommand(jobCmd)
}

// Poll until all matching Jobs have completed
func waitForJobsCompleted(clientset *kubernetes.Clientset, namespace, selector string) {
	fmt.Printf("Info: Waiting for Jobs in namespace %s with selector '%s'...\n", namespace, selector)
	for {
		// List jobs by namespace + label
		jobs, err := clientset.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: selector,
		})
		if err != nil {
			panic(err.Error())
		}

		if len(jobs.Items) == 0 {
			fmt.Println("Info: No Jobs found yet...")
		} else {
			allDone := true
			for _, job := range jobs.Items {
				// Determine how many completions are expected
				want := 1
				if job.Spec.Completions != nil {
					want = int(*job.Spec.Completions)
				}
				done := int(job.Status.Succeeded)
				fmt.Printf("State: Job %s --> %d/%d Completed\n", job.Name, done, want)
				if done < want {
					allDone = false
				}
			}
			if allDone {
				fmt.Println("Info: All Jobs are Completed!")
				return
			}
		}

		time.Sleep(time.Duration(timer) * time.Second)
	}
}

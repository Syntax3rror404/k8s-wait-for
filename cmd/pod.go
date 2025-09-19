/*
Copyright © 2025 Marcel Zapf
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// podCmd represents the pod command
var podCmd = &cobra.Command{
	Use:   "pod",
	Short: "wait for a pod to be ready",
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

		// 3. Wait for pods to be ready
		waitForPodsReady(clientset, namespace, label)
	},
}

func init() {
	rootCmd.AddCommand(podCmd)
}

// Check if a Pod has the Ready condition set to True
func isPodReady(pod v1.Pod) bool {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == v1.PodReady && cond.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

// Poll until all matching Pods are Ready
func waitForPodsReady(clientset *kubernetes.Clientset, namespace, selector string) {
	fmt.Printf("Info: Waiting for Pods in namespace %s with selector '%s'...\n", namespace, selector)
	for {
		// List pods by namespace + label
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: selector,
		})
		if err != nil {
			panic(err.Error())
		}

		if len(pods.Items) == 0 {
			fmt.Println("Info: No Pods found yet...")
		} else {
			allReady := true
			fmt.Printf("=== %s ===========\n", time.Now().Format(time.RFC1123))
			for _, pod := range pods.Items {
				ready := isPodReady(pod)

				fmt.Printf("State: Pod %s --> Ready: %v\n", pod.Name, ready)
				if !ready {
					allReady = false
				}

			}

			if allReady {
				fmt.Printf("==============================================\n")
				fmt.Println("Info: All pods are ready!")
				return
			}
		}

		time.Sleep(time.Duration(timer) * time.Second)
	}
}

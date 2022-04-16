package gencore

import "k8s.io/klog/v2"

func FatalfOnErr(err error) {
	if err != nil {
		klog.Fatalf("Error: %v", err)
	}
}

package gencore

import "k8s.io/klog/v2"

func FatalfOnErr(err error) {
	if err != nil {
		klog.Fatalf("error: %v", err)
	}
}

func Must1(err error) {
	if err != nil {
		klog.Fatalf("error: %v", err)
	}
	return
}

func Must[T any](t T, err error) T {
	Must1(err)
	return t
}

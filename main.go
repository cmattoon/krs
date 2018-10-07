package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mhausenblas/kubecuddler"
)

const (
	// ScrapeDelayInSec defines how long to wait between
	// sampling events, called the observation period:
	ScrapeDelayInSec = 5
)

func main() {
	ns := "krs"
	for {
		res := fromFirehose(ns)
		metrics := toOpenMetrics(ns, res)
		if metrics != "" {
			store(os.Stdout, metrics)
		}
		time.Sleep(ScrapeDelayInSec * time.Second)
	}
}

// fromFirehose uses kubectl to query for resources
// and returns them as a JSON format list string.
func fromFirehose(namespace string) string {
	if namespace == "" {
		namespace = "default"
	}
	res, err := kubecuddler.Kubectl(false, false, "", "get", "--namespace="+namespace, "all", "--output=json")
	if err != nil {
		log(err)
	}
	return res
}

// store takes OpenMetrics lines as input and stores it in the target file
// which could be, for example, stdout
func store(target io.Writer, metrics string) {
	fmt.Fprintf(target, "%v", metrics)
}

func log(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "\x1b[91m%v\x1b[0m\n", err)
}

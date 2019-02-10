package main

import (
	"log"
	"time"

	"github.com/jumincorp/constrictor"
	"github.com/jumincorp/micrometric"

	"github.com/btcsuite/btcd/rpcclient"
)

const programName = "equibit-core-metrics"

var (
	label             = constrictor.StringVar("label", "l", "default", "Label to identify this miner's data")
	node              = constrictor.AddressPortVar("node", "n", ":18331", "Address:Port of the node's RPC port")
	username          = constrictor.StringVar("user", "u", "default", "Node username")
	password          = constrictor.StringVar("pass", "", "default", "Node password")
	prometheusAddress = constrictor.AddressPortVar("prometheus", "p", ":40012", "Address:Port to expose to Prometheus")
	queryDelay        = constrictor.TimeDurationVar("time", "t", "30", "Delay between RPC calls to the miner")

	exporter micrometric.Exporter
	app      = constrictor.App("equibit-core-metrics", "Some Core Equibit Metrics", "Gaze lovingly into your Equibits", run)
)

func gather() error {
	log.Printf("gather\n")

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         node(),
		User:         username(),
		Pass:         password(),
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	//var wallet btcjson.InfoWalletResult

	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return err
	}
	defer client.Shutdown()

	// Get the current block count.
	accounts, err := client.ListAccounts()
	if err != nil {
		return err
	}

	log.Printf("acc %v", accounts)
	var metrics = make([]micrometric.Metric, len(accounts))

	i := 0
	total := 0

	for name, amount := range accounts {
		sanitizedAccountName := name
		if sanitizedAccountName == "" {
			sanitizedAccountName = "\\\"\\\""
		}
		eqb := amount.ToBTC()

		labels := make(map[string]string)
		labels["namespace"] = programName
		labels["account"] = sanitizedAccountName

		metrics[i] = micrometric.Metric{Labels: labels, Name: "equibit_balance", Value: eqb}
		i++

		// Check transactions:

		if l, err := client.ListTransactionsCountFrom(name, 10000, 0); err == nil {
			//log.Printf("len %v l %v\n", len(l), l)
			confirmed := 0
			abandoned := 0
			generated := 0
			watchOnly := 0
			trusted := 0

			for _, item := range l {
				if item.Confirmations > 2 {
					confirmed++
				}
				if item.Abandoned == true {
					abandoned++
				}
				if item.Generated == true {
					generated++
				}
				if item.InvolvesWatchOnly == true {
					watchOnly++
				}
				if item.Trusted == true {
					trusted++
				}
			}
			log.Printf("len %v confirmed %v abandoned %v generated %v watchOnly %v trusted %v\n", len(l), confirmed, abandoned, generated, watchOnly, trusted)
			total += len(l)
		}
	}
	log.Printf("total %v\n", total)

	log.Printf("metrics %v\n", metrics)
	exporter.Export(metrics)

	return nil
}

func run([]string) error {
	log.Printf("node %s u/p %s/%s prometheus %s\n", node(), username(), password(), prometheusAddress())
	exporter = micrometric.NewPrometheusExporter(prometheusAddress())
	go func() {
		for {
			if err := gather(); err != nil {
				log.Printf("Gather Error: %v\n", err)
			}
			time.Sleep(time.Second * queryDelay())
		}
	}()
	exporter.Setup()
	return nil
}

func main() {
	app.Execute()
}

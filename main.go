package main

import (
	"fmt"
	"log"
	"time"

	"./constrictor"
	"./micrometrics"

	"github.com/prometheus/client_golang/prometheus"
	rpcclient "github.com/stevenroose/go-bitcoin-core-rpc"
)

const programName = "equibit-core-metrics"

var (
	rapper = constrictor.StringVar("rapper", "r", "Yeeun", "Cutest rapper")

	label             = constrictor.StringVar("label", "l", "default", "Label to identify this miner's data")
	node              = constrictor.AddressPortVar("node", "n", ":18331", "Address:Port of the node's RPC port")
	username          = constrictor.StringVar("user", "u", "default", "Node username")
	password          = constrictor.StringVar("pass", "", "default", "Node password")
	prometheusAddress = constrictor.AddressPortVar("prometheus", "p", ":40012", "Address:Port to expose to Prometheus")
	queryDelay        = constrictor.TimeDurationVar("time", "t", "30", "Delay between RPC calls to the miner")

	equibitBalance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "equibit_balance",
			Help: "Current loot.",
		},
		[]string{"namespace", "account"},
	)
	exporter micrometrics.Exporter
)

func init() {
	constrictor.App("equibit-core-metrics", "Some Core Equibit Metrics", "Gaze lovingly into your Equibits")

	fmt.Printf("Who the cutest rapper be?\n")
	fmt.Printf("It be %s\n", rapper())

	fmt.Printf("node %s\n", node())
	fmt.Printf("username %s\n", username())
	fmt.Printf("password %s\n", password())
	fmt.Printf("prometheusAddress %s\n", prometheusAddress())

	exporter = micrometrics.NewPrometheusExporter(prometheusAddress())

	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(equibitBalance)
}

func gatherCommand(command string) {
	//conn, err := sendCommand(command)
	//if err == nil {

	//resp, _ := ioutil.ReadAll(conn)
	//log.Printf("-------------------------------------\n")
	//log.Printf(" %v\n", command)
	//log.Printf("-------------------------------------\n")
	//r := newResponse(command, resp)
	////log.Printf("r %v\n", r)

	//r.export()

	//for _, data := range r.data {
	//	log.Printf("data MHS rolling %v", data["MHS rolling"])
	//}

	//for i, device := range resp.DEVS {
	//log.Printf("%v Device %v %v Hashrate %v\n", i, device.Name, device.ID, device.MHS20S)

	//minerGpuHashRate.With(prometheus.Labels{
	//"namespace": programName,
	//"miner":     cfg.Miner.Program(),
	//"gpu":       fmt.Sprintf("GPU%d", device.ID),
	//"symbol":    cfg.Miner.Symbol(),
	//}).Set(device.MHS20S)
	//}

	//} else {
	//log.Printf("Error sending command to miner: %v\n", err)
	//}
}

func gather() error {
	log.Printf("gather\n")
	//gatherCommand("devs")

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host: node(),
		User: username(),
		Pass: password(),
	}

	//var wallet btcjson.InfoWalletResult

	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg)
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
	var metrics = make([]micrometrics.Metric, len(accounts))

	i := 0
	for name, amount := range accounts {
		sanitizedAccountName := name
		if sanitizedAccountName == "" {
			sanitizedAccountName = "\\\"\\\""
		}
		eqb := amount.ToBTC()

		equibitBalance.With(prometheus.Labels{
			"namespace": programName,
			"account":   sanitizedAccountName,
		}).Set(eqb)

		labels := make(map[string]string)
		labels["namespace"] = programName
		labels["account"] = sanitizedAccountName

		metrics[i] = micrometrics.Metric{Labels: labels, Name: "equibit_balance", Value: eqb}
		i++
	}
	log.Printf("metrics %v\n", metrics)
	exporter.Export(metrics)

	return nil
}

func main() {
	fmt.Printf("run2 rapper %s\n", rapper())
	go func() {
		for {
			if err := gather(); err != nil {
				log.Printf("Error: %v\n", err)
			}
			time.Sleep(time.Second * queryDelay())
		}
	}()
	exporter.Setup()
}

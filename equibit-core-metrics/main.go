package main

import (
	"log"
	"time"

	"../config"
	"../export"
	"github.com/prometheus/client_golang/prometheus"
	rpcclient "github.com/stevenroose/go-bitcoin-core-rpc"
)

const programName = "equibit-core-metrics"

var (
	equibitBalance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "equibit_balance",
			Help: "Current loot.",
		},
		[]string{"namespace", "account"},
	)

	cfg      *config.Config
	exporter export.Exporter
)

func init() {
	cfg = config.NewConfig(programName)
	exporter = export.NewPrometheus(cfg.Prometheus.Address())

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

func gather() {

	//gatherCommand("devs")

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host: cfg.Node.Host(),
		User: cfg.Node.User(),
		Pass: cfg.Node.Password(),
	}

	//var wallet btcjson.InfoWalletResult

	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	// Get the current block count.
	accounts, err := client.ListAccounts()
	if err != nil {
		log.Fatal(err)
	}
	var sum float64
	for name, amount := range accounts {
		sanitizedAccountName := name
		if sanitizedAccountName == "" {
			sanitizedAccountName = "\"\""
		}
		eqb := amount.ToBTC()
		sum += eqb

		equibitBalance.With(prometheus.Labels{
			"namespace": programName,
			"account":   sanitizedAccountName,
		}).Set(eqb)
	}
}

func main() {
	go func() {
		for {
			gather()
			time.Sleep(time.Second * cfg.QueryDelay())
		}
	}()
	exporter.Setup()
}

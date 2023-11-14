package main

import (
	"context"
	"flag"
	"github.com/dsrvlabs/aleo_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"

	"k8s.io/klog/v2"
)

const (
	httpTimeout = 5 * time.Second
)

var (
	rpcAddr = flag.String("rpcURI", "", "Aleo RPC URI")
	addr    = flag.String("addr", ":8080", "Listen address")
)

func init() {
	klog.InitFlags(nil)
}

type aleoCollector struct {
	rpcClient *rpc.RPCClient

	latestHeight *prometheus.Desc
	peerCount    *prometheus.Desc
	networkId    *prometheus.Desc
	latestRound  *prometheus.Desc
	cumulativeWeight *prometheus.Desc
	cumulativeProofTarget *prometheus.Desc
	coinbaseTarget *prometheus.Desc
	proofTarget *prometheus.Desc
	lastCoinbaseTarget *prometheus.Desc
	lastCoinbaseTimestamp *prometheus.Desc
	timestamp *prometheus.Desc
}

func NewAleoCollector(rpcAddr string) *aleoCollector {
	return &aleoCollector{
		rpcClient: rpc.NewRPCClient(rpcAddr),
		latestHeight: prometheus.NewDesc(
			"aleo_latest_height",
			"The latest block height",
			[]string{"account"}, nil),
		peerCount: prometheus.NewDesc(
			"aleo_peer_count",
			"The number of peers connected to the node",
			[]string{"account"}, nil),
		networkId: prometheus.NewDesc(
			"aleo_network_id",
			"The network ID of the block",
			[]string{"account"}, nil),
		latestRound: prometheus.NewDesc(
			"aleo_latest_round",
			"The round that produced this block",
			[]string{"account"}, nil),
		cumulativeWeight: prometheus.NewDesc(
			"aleo_cumulative_weight",
			"The cumulative weight for this block",
			[]string{"account"}, nil),
		cumulativeProofTarget: prometheus.NewDesc(
			"aleo_cumulative_proof_target",
			"The cumulative proof target for this block",
			[]string{"account"}, nil),
		coinbaseTarget: prometheus.NewDesc(
			"aleo_coinbase_target",
			"The coinbase target for this block",
			[]string{"account"}, nil),
		proofTarget: prometheus.NewDesc(
			"aleo_proof_target",
			"The proof target for this block",
			[]string{"account"}, nil),
		lastCoinbaseTarget: prometheus.NewDesc(
			"aleo_last_coinbase_target",
			"The coinbase target for the last coinbase",
			[]string{"account"}, nil),
		lastCoinbaseTimestamp: prometheus.NewDesc(
			"aleo_last_coinbase_timestamp",
			"The unix timestamp(UTC) for the last coinbase",
			[]string{"account"}, nil),
		timestamp: prometheus.NewDesc(
			"aleo_timestamp",
			"The unix timestamp(UTC) for this block",
			[]string{"account"}, nil),
	}
}

func (c *aleoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.latestHeight
	ch <- c.peerCount
	ch <- c.networkId
	ch <- c.latestRound
	ch <- c.cumulativeWeight
	ch <- c.cumulativeProofTarget
	ch <- c.coinbaseTarget
	ch <- c.proofTarget
	ch <- c.lastCoinbaseTarget
	ch <- c.lastCoinbaseTimestamp
	ch <- c.timestamp
}

func (c *aleoCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	height, err := c.rpcClient.GetLatestHeight(ctx)
	if err != nil {
		klog.Errorf("failed to get latest height, err :%v", err)
		ch <- prometheus.NewInvalidMetric(c.latestHeight, err)
	} else {
		ch <- prometheus.MustNewConstMetric(c.latestHeight, prometheus.GaugeValue, float64(height), "account")
	}

	peers, err := c.rpcClient.GetPeerCount(ctx)
	if err != nil {
		klog.Errorf("failed to get peer count, err :%v", err)
		ch <- prometheus.NewInvalidMetric(c.peerCount, err)
	} else {
		ch <- prometheus.MustNewConstMetric(c.peerCount, prometheus.GaugeValue, float64(peers), "account")
	}

	block, err := c.rpcClient.GetLatestBlock(ctx)
	if err != nil {
		klog.Errorf("failed to get latest block err :%v", err)
		ch <- prometheus.NewInvalidMetric(c.networkId, err)
		ch <- prometheus.NewInvalidMetric(c.latestRound, err)
		ch <- prometheus.NewInvalidMetric(c.cumulativeWeight, err)
		ch <- prometheus.NewInvalidMetric(c.cumulativeProofTarget, err)
		ch <- prometheus.NewInvalidMetric(c.coinbaseTarget, err)
		ch <- prometheus.NewInvalidMetric(c.proofTarget, err)
		ch <- prometheus.NewInvalidMetric(c.lastCoinbaseTarget, err)
		ch <- prometheus.NewInvalidMetric(c.lastCoinbaseTimestamp, err)
		ch <- prometheus.NewInvalidMetric(c.timestamp, err)
	} else {
		ch <- prometheus.MustNewConstMetric(c.networkId, prometheus.GaugeValue, float64(block.Network), "account")
		ch <- prometheus.MustNewConstMetric(c.latestRound, prometheus.GaugeValue, float64(block.Round), "account")
		ch <- prometheus.MustNewConstMetric(c.cumulativeWeight, prometheus.GaugeValue, float64(block.CumulativeWeight), "account")
		ch <- prometheus.MustNewConstMetric(c.cumulativeProofTarget, prometheus.GaugeValue, float64(block.CumulativeProofTarget), "account")
		ch <- prometheus.MustNewConstMetric(c.coinbaseTarget, prometheus.GaugeValue, float64(block.CoinbaseTarget), "account")
		ch <- prometheus.MustNewConstMetric(c.proofTarget, prometheus.GaugeValue, float64(block.ProofTarget), "account")
		ch <- prometheus.MustNewConstMetric(c.lastCoinbaseTarget, prometheus.GaugeValue, float64(block.LastCoinbaseTarget), "account")
		ch <- prometheus.MustNewConstMetric(c.lastCoinbaseTimestamp, prometheus.GaugeValue, float64(block.LastCoinbaseTimestamp), "account")
		ch <- prometheus.MustNewConstMetric(c.timestamp, prometheus.GaugeValue, float64(block.Timestamp), "account")
	}
}

func main() {
	flag.Parse()

	if *rpcAddr == "" {
		klog.Fatal("Please specify -rpcURI")
	}

	collector := NewAleoCollector(*rpcAddr)

	prometheus.MustRegister(collector)
	http.Handle("/metrics", promhttp.Handler())

	klog.Infof("listening on %s", *addr)
	klog.Fatal(http.ListenAndServe(*addr, nil))
}

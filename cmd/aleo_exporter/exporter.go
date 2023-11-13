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
	}
}

func (c *aleoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.latestHeight
	ch <- c.peerCount
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

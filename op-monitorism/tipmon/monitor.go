package tipmon

import (
	"context"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/metrics"
	"github.com/ethereum-optimism/optimism/op-service/sources"

	"github.com/ethereum/go-ethereum/log"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricsNamespace = "tip_mon"
)

type Monitor struct {
	log log.Logger

	rpc client.RPC

	// metrics
	laggingDistance     *prometheus.GaugeVec
	unexpectedRpcErrors *prometheus.CounterVec
}

func NewMonitor(ctx context.Context, log log.Logger, m metrics.Factory, cfg CLIConfig) (*Monitor, error) {
	log.Info("creating tip time monitor")
	rpc, err := client.NewRPC(ctx, log, cfg.NodeUrl)
	if err != nil {
		return nil, err
	}

	return &Monitor{
		log: log,
		rpc: rpc,

		laggingDistance: m.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: MetricsNamespace,
			Name:      "lagging",
			Help:      "lagging distance between tip and real time",
		}, []string{"type"}),
		unexpectedRpcErrors: m.NewCounterVec(prometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Name:      "unexpectedRpcErrors",
			Help:      "number of unexpcted rpc errors",
		}, []string{"section", "name"}),
	}, nil
}

func (m *Monitor) Run(ctx context.Context) {
	m.log.Info("querying tip...")
	result := new(sources.RPCHeader)
	if err := m.rpc.CallContext(ctx, result, "eth_getBlockByNumber", "latest", false); err != nil {
		m.log.Error("failed eth_getBlockByNumber request", "err", err)
		m.unexpectedRpcErrors.WithLabelValues("laggingDistance", "eth_getBlockByNumber").Inc()
		return
	}

	lag := time.Now().UTC().Unix() - int64(result.Time)
	m.laggingDistance.WithLabelValues("latest").Set(float64(lag))
	m.log.Info("set lagging distance", "type", "latest", "lag", lag)
}

func (m *Monitor) Close(_ context.Context) error {
	m.rpc.Close()
	return nil
}

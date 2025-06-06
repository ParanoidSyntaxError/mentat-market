package ccip

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/internal/gethwrappers2/generated/offchainaggregator"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/batchreader"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/ccipdataprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

const OffchainAggregator = "OffchainAggregator"
const DecimalsMethodName = "decimals"
const LatestRoundDataMethodName = "latestRoundData"

func GenericAddrToEvm(addr ccip.Address) (common.Address, error) {
	return ccipcalc.GenericAddrToEvm(addr)
}

func EvmAddrToGeneric(addr common.Address) ccip.Address {
	return ccipcalc.EvmAddrToGeneric(addr)
}

func NewEvmPriceRegistry(lp logpoller.LogPoller, ec client.Client, lggr logger.Logger, pluginLabel string) *ccipdataprovider.EvmPriceRegistry {
	return ccipdataprovider.NewEvmPriceRegistry(lp, ec, lggr, pluginLabel)
}

type VersionFinder = factory.VersionFinder

func NewCommitStoreReader(ctx context.Context, lggr logger.Logger, versionFinder VersionFinder, address ccip.Address, ec client.Client, lp logpoller.LogPoller, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) (ccipdata.CommitStoreReader, error) {
	return factory.NewCommitStoreReader(ctx, lggr, versionFinder, address, ec, lp, feeEstimatorConfig)
}

func CloseCommitStoreReader(ctx context.Context, lggr logger.Logger, versionFinder VersionFinder, address ccip.Address, ec client.Client, lp logpoller.LogPoller, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) error {
	return factory.CloseCommitStoreReader(ctx, lggr, versionFinder, address, ec, lp, feeEstimatorConfig)
}

func NewOffRampReader(ctx context.Context, lggr logger.Logger, versionFinder VersionFinder, addr ccip.Address, destClient client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator, destMaxGasPrice *big.Int, registerFilters bool, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) (ccipdata.OffRampReader, error) {
	return factory.NewOffRampReader(ctx, lggr, versionFinder, addr, destClient, lp, estimator, destMaxGasPrice, registerFilters, feeEstimatorConfig)
}

func CloseOffRampReader(ctx context.Context, lggr logger.Logger, versionFinder VersionFinder, addr ccip.Address, destClient client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator, destMaxGasPrice *big.Int, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) error {
	return factory.CloseOffRampReader(ctx, lggr, versionFinder, addr, destClient, lp, estimator, destMaxGasPrice, feeEstimatorConfig)
}

func NewEvmVersionFinder() factory.EvmVersionFinder {
	return factory.NewEvmVersionFinder()
}

func NewOnRampReader(ctx context.Context, lggr logger.Logger, versionFinder VersionFinder, sourceSelector, destSelector uint64, onRampAddress ccip.Address, sourceLP logpoller.LogPoller, source client.Client) (ccipdata.OnRampReader, error) {
	return factory.NewOnRampReader(ctx, lggr, versionFinder, sourceSelector, destSelector, onRampAddress, sourceLP, source)
}

func CloseOnRampReader(ctx context.Context, lggr logger.Logger, versionFinder VersionFinder, sourceSelector, destSelector uint64, onRampAddress ccip.Address, sourceLP logpoller.LogPoller, source client.Client) error {
	return factory.CloseOnRampReader(ctx, lggr, versionFinder, sourceSelector, destSelector, onRampAddress, sourceLP, source)
}

type OffRampReader = ccipdata.OffRampReader

type DynamicPriceGetterClient = pricegetter.DynamicPriceGetterClient

type DynamicPriceGetter = pricegetter.DynamicPriceGetter

type AllTokensPriceGetter = pricegetter.AllTokensPriceGetter

func NewPipelineGetter(
	source string,
	runner pipeline.Runner,
	jobID int32,
	externalJobID uuid.UUID,
	name string,
	lggr logger.Logger,
	sourceNativeTokenAddr ccip.Address,
	sourceChainSelector uint64,
	destChainSelector uint64,
) (*pricegetter.PipelineGetter, error) {
	return pricegetter.NewPipelineGetter(source, runner, jobID, externalJobID, name, lggr,
		sourceNativeTokenAddr, sourceChainSelector, destChainSelector)
}

func NewDynamicPriceGetterClient(batchCaller rpclib.EvmBatchCaller) DynamicPriceGetterClient {
	return pricegetter.NewDynamicPriceGetterClient(batchCaller)
}

func NewDynamicPriceGetter(cfg config.DynamicPriceGetterConfig, contractReaders map[uint64]types.ContractReader) (*DynamicPriceGetter, error) {
	return pricegetter.NewDynamicPriceGetter(cfg, contractReaders)
}

func NewDynamicLimitedBatchCaller(
	lggr logger.Logger, batchSender rpclib.BatchSender, batchSizeLimit, backOffMultiplier, parallelRpcCallsLimit uint,
) *rpclib.DynamicLimitedBatchCaller {
	return rpclib.NewDynamicLimitedBatchCaller(lggr, batchSender, batchSizeLimit, backOffMultiplier, parallelRpcCallsLimit)
}

func NewUSDCReader(ctx context.Context, lggr logger.Logger, jobID string, transmitter common.Address, lp logpoller.LogPoller, registerFilters bool) (*ccipdata.USDCReaderImpl, error) {
	return ccipdata.NewUSDCReader(ctx, lggr, jobID, transmitter, lp, registerFilters)
}

func CloseUSDCReader(ctx context.Context, lggr logger.Logger, jobID string, transmitter common.Address, lp logpoller.LogPoller) error {
	return ccipdata.CloseUSDCReader(ctx, lggr, jobID, transmitter, lp)
}

type USDCReaderImpl = ccipdata.USDCReaderImpl

var DefaultRpcBatchSizeLimit = rpclib.DefaultRpcBatchSizeLimit
var DefaultRpcBatchBackOffMultiplier = rpclib.DefaultRpcBatchBackOffMultiplier
var DefaultMaxParallelRpcCalls = rpclib.DefaultMaxParallelRpcCalls

func NewEVMTokenPoolBatchedReader(lggr logger.Logger, remoteChainSelector uint64, offRampAddress ccip.Address, evmBatchCaller rpclib.EvmBatchCaller) (*batchreader.EVMTokenPoolBatchedReader, error) {
	return batchreader.NewEVMTokenPoolBatchedReader(lggr, remoteChainSelector, offRampAddress, evmBatchCaller)
}

type ChainAgnosticPriceRegistry struct {
	p ChainAgnosticPriceRegistryFactory
}

// [ChainAgnosticPriceRegistryFactory] is satisfied by [commontypes.CCIPCommitProvider] and [commontypes.CCIPExecProvider]
type ChainAgnosticPriceRegistryFactory interface {
	NewPriceRegistryReader(ctx context.Context, addr ccip.Address) (ccip.PriceRegistryReader, error)
}

func (c *ChainAgnosticPriceRegistry) NewPriceRegistryReader(ctx context.Context, addr ccip.Address) (ccip.PriceRegistryReader, error) {
	return c.p.NewPriceRegistryReader(ctx, addr)
}

func NewChainAgnosticPriceRegistry(provider ChainAgnosticPriceRegistryFactory) *ChainAgnosticPriceRegistry {
	return &ChainAgnosticPriceRegistry{provider}
}

type JSONCommitOffchainConfigV1_2_0 = v1_2_0.JSONCommitOffchainConfig
type CommitOnchainConfig = ccipdata.CommitOnchainConfig

func NewCommitOffchainConfig(
	gasPriceDeviationPPB uint32,
	gasPriceHeartBeat time.Duration,
	tokenPriceDeviationPPB uint32,
	tokenPriceHeartBeat time.Duration,
	inflightCacheExpiry time.Duration,
	priceReportingDisabled bool,
) ccip.CommitOffchainConfig {
	return ccipdata.NewCommitOffchainConfig(gasPriceDeviationPPB, gasPriceHeartBeat, tokenPriceDeviationPPB, tokenPriceHeartBeat, inflightCacheExpiry, priceReportingDisabled)
}

const OffChainAggregatorABI = offchainaggregator.OffchainAggregatorABI

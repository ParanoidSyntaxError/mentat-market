package blockheaderfeeder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestValidate(t *testing.T) {
	v1Coordinator := types.EIP55Address("0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139")
	v2Coordinator := types.EIP55Address("0x2be990eE17832b59E0086534c5ea2459Aa75E38F")
	v2PlusCoordinator := types.EIP55Address("0x92B5e28Ac583812874e4271380c7d070C5FB6E6b")
	fromAddresses := []types.EIP55Address{("0x469aA2CD13e037DC5236320783dCfd0e641c0559")}

	var tests = []struct {
		name      string
		toml      string
		assertion func(t *testing.T, os job.Job, err error)
	}{
		{
			name: "valid",
			toml: `
type = "blockheaderfeeder"
name = "valid-test"
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
coordinatorV2Address = "0x2be990eE17832b59E0086534c5ea2459Aa75E38F"
coordinatorV2PlusAddress = "0x92B5e28Ac583812874e4271380c7d070C5FB6E6b"
lookbackBlocks = 2000
waitBlocks = 500
blockhashStoreAddress = "0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"
batchBlockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
pollPeriod = "23s"
runTimeout = "7s"
evmChainID = "4"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]
getBlockhashesBatchSize = 20
storeBlockhashesBatchSize = 10
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				require.Equal(t, job.BlockHeaderFeeder, os.Type)
				require.Equal(t, "valid-test", os.Name.String)
				require.Equal(t, &v1Coordinator,
					os.BlockHeaderFeederSpec.CoordinatorV1Address)
				require.Equal(t, &v2Coordinator,
					os.BlockHeaderFeederSpec.CoordinatorV2Address)
				require.Equal(t, &v2PlusCoordinator,
					os.BlockHeaderFeederSpec.CoordinatorV2PlusAddress)
				require.Equal(t, int32(2000), os.BlockHeaderFeederSpec.LookbackBlocks)
				require.Equal(t, int32(500), os.BlockHeaderFeederSpec.WaitBlocks)
				require.Equal(t, types.EIP55Address("0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"),
					os.BlockHeaderFeederSpec.BlockhashStoreAddress)
				require.Equal(t, types.EIP55Address("0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"),
					os.BlockHeaderFeederSpec.BatchBlockhashStoreAddress)
				require.Equal(t, 23*time.Second, os.BlockHeaderFeederSpec.PollPeriod)
				require.Equal(t, 7*time.Second, os.BlockHeaderFeederSpec.RunTimeout)
				require.Equal(t, big.NewI(4), os.BlockHeaderFeederSpec.EVMChainID)
				require.Equal(t, fromAddresses,
					os.BlockHeaderFeederSpec.FromAddresses)
				require.Equal(t, uint16(20),
					os.BlockHeaderFeederSpec.GetBlockhashesBatchSize)
				require.Equal(t, uint16(10),
					os.BlockHeaderFeederSpec.StoreBlockhashesBatchSize)
			},
		},
		{
			name: "defaults-test",
			toml: `
type = "blockheaderfeeder"
name = "defaults-test"
evmChainID = "4"
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
coordinatorV2Address = "0x2be990eE17832b59E0086534c5ea2459Aa75E38F"
blockhashStoreAddress = "0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"
batchBlockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				require.Equal(t, int32(1000), os.BlockHeaderFeederSpec.LookbackBlocks)
				require.Equal(t, int32(256), os.BlockHeaderFeederSpec.WaitBlocks)
				require.Equal(t, 15*time.Second, os.BlockHeaderFeederSpec.PollPeriod)
				require.Equal(t, 30*time.Second, os.BlockHeaderFeederSpec.RunTimeout)
				require.Equal(t, big.NewI(4), os.BlockHeaderFeederSpec.EVMChainID)
				require.Equal(t, fromAddresses,
					os.BlockHeaderFeederSpec.FromAddresses)
				require.Equal(t, uint16(100),
					os.BlockHeaderFeederSpec.GetBlockhashesBatchSize)
				require.Equal(t, uint16(10),
					os.BlockHeaderFeederSpec.StoreBlockhashesBatchSize)
			},
		},
		{
			name: "invalid-job-type",
			toml: `
type = "invalidjob"
name = "invalid-job-type"
lookbackBlocks = 2000
waitBlocks = 500
blockhashStoreAddress = "0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"
batchBlockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
pollPeriod = "23s"
runTimeout = "7s"
evmChainID = "4"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]
getBlockhashesBatchSize = 20
storeBlockhashesBatchSize = 10
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.EqualError(t, err, "unsupported type invalidjob")
			},
		},
		{
			name: "missing-coordinators",
			toml: `
type = "blockheaderfeeder"
name = "missing-coordinators"
lookbackBlocks = 2000
waitBlocks = 500
blockhashStoreAddress = "0x3e20Cef636EdA7ba135bCbA4fe6177Bd3cE0aB17"
batchBlockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
pollPeriod = "23s"
runTimeout = "7s"
evmChainID = "4"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]
getBlockhashesBatchSize = 20
storeBlockhashesBatchSize = 10
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Equal(t, `at least one of "coordinatorV1Address", "coordinatorV2Address" and "coordinatorV2PlusAddress" must be set`, err.Error())
			},
		},
		{
			name: "missing blockhash store address",
			toml: `
type = "blockheaderfeeder"
name = "missing blockhash store address"
lookbackBlocks = 2000
waitBlocks = 500
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
batchBlockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
pollPeriod = "23s"
runTimeout = "7s"
evmChainID = "4"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]
getBlockhashesBatchSize = 20
storeBlockhashesBatchSize = 10
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Equal(t, `"blockhashStoreAddress" must be set`, err.Error())
			},
		},
		{
			name: "missing batch blockhash store address",
			toml: `
type = "blockheaderfeeder"
name = "missing batch blockhash store address"
lookbackBlocks = 2000
waitBlocks = 500
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
blockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
pollPeriod = "23s"
runTimeout = "7s"
evmChainID = "4"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]
getBlockhashesBatchSize = 20
storeBlockhashesBatchSize = 10
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Equal(t, `"batchBlockhashStoreAddress" must be set`, err.Error())
			},
		},
		{
			name: "missing evmChainID",
			toml: `
type = "blockheaderfeeder"
name = "missing evmChainID"
lookbackBlocks = 2000
waitBlocks = 500
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
blockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
batchBlockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
pollPeriod = "23s"
runTimeout = "7s"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]
getBlockhashesBatchSize = 20
storeBlockhashesBatchSize = 10
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Equal(t, `"evmChainID" must be set`, err.Error())
			},
		},
		{
			name: "wait block lower than 256 blocks",
			toml: `
type = "blockheaderfeeder"
name = "wait block lower than 256 blocks"
lookbackBlocks = 2000
waitBlocks = 255
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
blockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
batchBlockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
pollPeriod = "23s"
runTimeout = "7s"
evmChainID = "4"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]
getBlockhashesBatchSize = 20
storeBlockhashesBatchSize = 10
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Equal(t, `"waitBlocks" must be greater than or equal to 256`, err.Error())
			},
		},
		{
			name: "lookback block lower than 256 blocks",
			toml: `
type = "blockheaderfeeder"
name = "lookback block lower than 256 blocks"
lookbackBlocks = 255
waitBlocks = 256
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
blockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
batchBlockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
pollPeriod = "23s"
runTimeout = "7s"
evmChainID = "4"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]
getBlockhashesBatchSize = 20
storeBlockhashesBatchSize = 10
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Equal(t, `"lookbackBlocks" must be greater than 256`, err.Error())
			},
		},
		{
			name: "lookback blocks lower than wait blocks",
			toml: `
type = "blockheaderfeeder"
name = "lookback blocks lower than wait blocks"
lookbackBlocks = 300
waitBlocks = 500
coordinatorV1Address = "0x1F72B4A5DCf7CC6d2E38423bF2f4BFA7db97d139"
blockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
batchBlockhashStoreAddress = "0xD04E5b2ea4e55AEbe6f7522bc2A69Ec6639bfc63"
pollPeriod = "23s"
runTimeout = "7s"
evmChainID = "4"
fromAddresses = ["0x469aA2CD13e037DC5236320783dCfd0e641c0559"]
getBlockhashesBatchSize = 20
storeBlockhashesBatchSize = 10
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Equal(t, `"lookbackBlocks" must be greater than "waitBlocks"`, err.Error())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, err := ValidatedSpec(test.toml)
			test.assertion(t, s, err)
		})
	}
}

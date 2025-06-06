package resolver

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	clnull "github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// Specs are only embedded on the job and are not fetchable by it's own id, so
// we test the spec resolvers by fetching a job by id.

func TestResolver_CronSpec(t *testing.T) {
	var (
		id = int32(1)
	)

	testCases := []GQLTestCase{
		{
			name:          "cron spec success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.Cron,
					CronSpec: &job.CronSpec{
						CronSchedule: "CRON_TZ=UTC 0 0 1 1 *",
						EVMChainID:   ubig.NewI(42),
						CreatedAt:    f.Timestamp(),
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on CronSpec {
									schedule
									evmChainID
									createdAt
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "CronSpec",
							"schedule": "CRON_TZ=UTC 0 0 1 1 *",
							"evmChainID": "42",
							"createdAt": "2021-01-01T00:00:00Z"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_DirectRequestSpec(t *testing.T) {
	var (
		id               = int32(1)
		requesterAddress = common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	)
	contractAddress, err := evmtypes.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		{
			name:          "direct request spec success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.DirectRequest,
					DirectRequestSpec: &job.DirectRequestSpec{
						ContractAddress:          contractAddress,
						CreatedAt:                f.Timestamp(),
						EVMChainID:               ubig.NewI(42),
						MinIncomingConfirmations: clnull.NewUint32(1, true),
						MinContractPayment:       commonassets.NewLinkFromJuels(1000),
						Requesters:               models.AddressCollection{requesterAddress},
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on DirectRequestSpec {
									contractAddress
									createdAt
									evmChainID
									minIncomingConfirmations
									minContractPaymentLinkJuels
									requesters
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "DirectRequestSpec",
							"contractAddress": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"createdAt": "2021-01-01T00:00:00Z",
							"evmChainID": "42",
							"minIncomingConfirmations": 1,
							"minContractPaymentLinkJuels": "1000",
							"requesters": ["0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"]
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_FluxMonitorSpec(t *testing.T) {
	var (
		id = int32(1)
	)
	contractAddress, err := evmtypes.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		{
			name:          "flux monitor spec with standard timers",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.FluxMonitor,
					FluxMonitorSpec: &job.FluxMonitorSpec{
						ContractAddress:   contractAddress,
						CreatedAt:         f.Timestamp(),
						EVMChainID:        ubig.NewI(42),
						DrumbeatEnabled:   false,
						IdleTimerDisabled: false,
						IdleTimerPeriod:   1 * time.Hour,
						MinPayment:        commonassets.NewLinkFromJuels(1000),
						PollTimerDisabled: false,
						PollTimerPeriod:   1 * time.Minute,
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on FluxMonitorSpec {
									absoluteThreshold
									contractAddress
									createdAt
									drumbeatEnabled
									drumbeatRandomDelay
									drumbeatSchedule
									evmChainID
									idleTimerDisabled
									idleTimerPeriod
									minPayment
									pollTimerDisabled
									pollTimerPeriod
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "FluxMonitorSpec",
							"absoluteThreshold": 0,
							"contractAddress": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"createdAt": "2021-01-01T00:00:00Z",
							"drumbeatEnabled": false,
							"drumbeatRandomDelay": null,
							"drumbeatSchedule": null,
							"evmChainID": "42",
							"idleTimerDisabled": false,
							"idleTimerPeriod": "1h0m0s",
							"minPayment": "1000",
							"pollTimerDisabled": false,
							"pollTimerPeriod": "1m0s"
						}
					}
				}
			`,
		},
		{
			name:          "flux monitor spec with drumbeat",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.FluxMonitor,
					FluxMonitorSpec: &job.FluxMonitorSpec{
						ContractAddress:     contractAddress,
						CreatedAt:           f.Timestamp(),
						EVMChainID:          ubig.NewI(42),
						DrumbeatEnabled:     true,
						DrumbeatRandomDelay: 1 * time.Second,
						DrumbeatSchedule:    "CRON_TZ=UTC 0 0 1 1 *",
						IdleTimerDisabled:   true,
						IdleTimerPeriod:     1 * time.Hour,
						MinPayment:          commonassets.NewLinkFromJuels(1000),
						PollTimerDisabled:   true,
						PollTimerPeriod:     1 * time.Minute,
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on FluxMonitorSpec {
									absoluteThreshold
									contractAddress
									createdAt
									drumbeatEnabled
									drumbeatRandomDelay
									drumbeatSchedule
									evmChainID
									idleTimerDisabled
									idleTimerPeriod
									minPayment
									pollTimerDisabled
									pollTimerPeriod
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "FluxMonitorSpec",
							"absoluteThreshold": 0,
							"contractAddress": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"createdAt": "2021-01-01T00:00:00Z",
							"drumbeatEnabled": true,
							"drumbeatRandomDelay": "1s",
							"drumbeatSchedule": "CRON_TZ=UTC 0 0 1 1 *",
							"evmChainID": "42",
							"idleTimerDisabled": true,
							"idleTimerPeriod": "1h0m0s",
							"minPayment": "1000",
							"pollTimerDisabled": true,
							"pollTimerPeriod": "1m0s"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_KeeperSpec(t *testing.T) {
	var (
		id          = int32(1)
		fromAddress = common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	)
	contractAddress, err := evmtypes.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		{
			name:          "keeper spec",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.Keeper,
					KeeperSpec: &job.KeeperSpec{
						ContractAddress: contractAddress,
						CreatedAt:       f.Timestamp(),
						EVMChainID:      ubig.NewI(42),
						FromAddress:     evmtypes.EIP55AddressFromAddress(fromAddress),
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on KeeperSpec {
									contractAddress
									createdAt
									evmChainID
									fromAddress
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "KeeperSpec",
							"contractAddress": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"createdAt": "2021-01-01T00:00:00Z",
							"evmChainID": "42",
							"fromAddress": "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_OCRSpec(t *testing.T) {
	var (
		id = int32(1)
	)
	contractAddress, err := evmtypes.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	transmitterAddress, err := evmtypes.NewEIP55Address("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	require.NoError(t, err)

	keyBundleID := models.MustSha256HashFromHex("f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5")

	testCases := []GQLTestCase{
		{
			name:          "OCR spec",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.OffchainReporting,
					OCROracleSpec: &job.OCROracleSpec{
						BlockchainTimeout:                      models.Interval(1 * time.Minute),
						ContractAddress:                        contractAddress,
						ContractConfigConfirmations:            1,
						ContractConfigTrackerPollInterval:      models.Interval(1 * time.Minute),
						ContractConfigTrackerSubscribeInterval: models.Interval(2 * time.Minute),
						DatabaseTimeout:                        models.NewInterval(3 * time.Second),
						ObservationGracePeriod:                 models.NewInterval(4 * time.Second),
						ContractTransmitterTransmitTimeout:     models.NewInterval(555 * time.Millisecond),
						CreatedAt:                              f.Timestamp(),
						EVMChainID:                             ubig.NewI(42),
						IsBootstrapPeer:                        false,
						EncryptedOCRKeyBundleID:                &keyBundleID,
						ObservationTimeout:                     models.Interval(2 * time.Minute),
						P2PV2Bootstrappers:                     pq.StringArray{"12D3KooWL3XJ9EMCyZvmmGXL2LMiVBtrVa2BuESsJiXkSj7333Jw@localhost:5001"},
						TransmitterAddress:                     &transmitterAddress,
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on OCRSpec {
									blockchainTimeout
									contractAddress
									contractConfigConfirmations
									contractConfigTrackerPollInterval
									contractConfigTrackerSubscribeInterval
									databaseTimeout
									observationGracePeriod
									contractTransmitterTransmitTimeout
									createdAt
									evmChainID
									isBootstrapPeer
									keyBundleID
									observationTimeout
									p2pv2Bootstrappers
									transmitterAddress
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "OCRSpec",
							"blockchainTimeout": "1m0s",
							"contractAddress": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"contractConfigConfirmations": 1,
							"contractConfigTrackerPollInterval": "1m0s",
							"contractConfigTrackerSubscribeInterval": "2m0s",
							"databaseTimeout": "3s",
							"observationGracePeriod": "4s",
							"contractTransmitterTransmitTimeout": "555ms",
							"createdAt": "2021-01-01T00:00:00Z",
							"evmChainID": "42",
							"isBootstrapPeer": false,
							"keyBundleID": "f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5",
							"observationTimeout": "2m0s",
							"p2pv2Bootstrappers": ["12D3KooWL3XJ9EMCyZvmmGXL2LMiVBtrVa2BuESsJiXkSj7333Jw@localhost:5001"],
							"transmitterAddress": "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_OCR2Spec(t *testing.T) {
	var (
		id = int32(1)
	)
	contractAddress, err := evmtypes.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	transmitterAddress, err := evmtypes.NewEIP55Address("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	require.NoError(t, err)

	keyBundleID := models.MustSha256HashFromHex("f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5")

	relayConfig := map[string]interface{}{
		"chainID": 1337,
	}
	pluginConfig := map[string]interface{}{
		"juelsPerFeeCoinSource": 100000000,
	}
	onchainSigningStrategy := map[string]interface{}{
		"strategyName": "multi-chain",
		"config": map[string]any{
			"evm": "b3df4d8748b67731a1112e8b45a764941974f5590c93672eebbc4f3504dd10ed",
		},
	}
	require.NoError(t, err)

	testCases := []GQLTestCase{
		{
			name:          "OCR 2 spec",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.OffchainReporting2,
					OCR2OracleSpec: &job.OCR2OracleSpec{
						BlockchainTimeout:                 models.Interval(1 * time.Minute),
						ContractID:                        contractAddress.String(),
						ContractConfigConfirmations:       1,
						ContractConfigTrackerPollInterval: models.Interval(1 * time.Minute),
						OnchainSigningStrategy:            onchainSigningStrategy,
						CreatedAt:                         f.Timestamp(),
						OCRKeyBundleID:                    null.StringFrom(keyBundleID.String()),
						MonitoringEndpoint:                null.StringFrom("https://monitor.endpoint"),
						P2PV2Bootstrappers:                pq.StringArray{"12D3KooWL3XJ9EMCyZvmmGXL2LMiVBtrVa2BuESsJiXkSj7333Jw@localhost:5001"},
						Relay:                             relay.NetworkEVM,
						RelayConfig:                       relayConfig,
						TransmitterID:                     null.StringFrom(transmitterAddress.String()),
						PluginType:                        types.Median,
						PluginConfig:                      pluginConfig,
						AllowNoBootstrappers:              true,
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on OCR2Spec {
									blockchainTimeout
									contractID
									contractConfigConfirmations
									contractConfigTrackerPollInterval
									onchainSigningStrategy
									createdAt
									ocrKeyBundleID
									monitoringEndpoint
									p2pv2Bootstrappers
									relay
									relayConfig
									transmitterID
									pluginType
									pluginConfig
									allowNoBootstrappers
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "OCR2Spec",
							"blockchainTimeout": "1m0s",
							"contractID": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"contractConfigConfirmations": 1,
							"contractConfigTrackerPollInterval": "1m0s",
							"onchainSigningStrategy": {
								"strategyName": "multi-chain",
								"config": {
									"evm": "b3df4d8748b67731a1112e8b45a764941974f5590c93672eebbc4f3504dd10ed"
								}
							},
							"createdAt": "2021-01-01T00:00:00Z",
							"ocrKeyBundleID": "f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5",
							"monitoringEndpoint": "https://monitor.endpoint",
							"p2pv2Bootstrappers": ["12D3KooWL3XJ9EMCyZvmmGXL2LMiVBtrVa2BuESsJiXkSj7333Jw@localhost:5001"],
							"relay": "evm",
							"relayConfig": {
								"chainID": 1337
							},
							"transmitterID": "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42",
							"pluginType": "median",
							"pluginConfig": {
								"juelsPerFeeCoinSource": 100000000
							},
							"allowNoBootstrappers": true
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_VRFSpec(t *testing.T) {
	var (
		id = int32(1)
	)
	coordinatorAddress, err := evmtypes.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	batchCoordinatorAddress, err := evmtypes.NewEIP55Address("0x0ad9FE7a58216242a8475ca92F222b0640E26B63")
	require.NoError(t, err)

	fromAddress1, err := evmtypes.NewEIP55Address("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	require.NoError(t, err)

	fromAddress2, err := evmtypes.NewEIP55Address("0x2301958F1BFbC9A068C2aC9c6166Bf483b95864C")
	require.NoError(t, err)

	pubKey, err := secp256k1.NewPublicKeyFromHex("0x9dc09a0f898f3b5e8047204e7ce7e44b587920932f08431e29c9bf6923b8450a01")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		{
			name:          "vrf spec",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.VRF,
					VRFSpec: &job.VRFSpec{
						BatchCoordinatorAddress:       &batchCoordinatorAddress,
						BatchFulfillmentEnabled:       true,
						CustomRevertsPipelineEnabled:  true,
						MinIncomingConfirmations:      1,
						CoordinatorAddress:            coordinatorAddress,
						CreatedAt:                     f.Timestamp(),
						EVMChainID:                    ubig.NewI(42),
						FromAddresses:                 []evmtypes.EIP55Address{fromAddress1, fromAddress2},
						PollPeriod:                    1 * time.Minute,
						PublicKey:                     pubKey,
						RequestedConfsDelay:           10,
						RequestTimeout:                24 * time.Hour,
						ChunkSize:                     25,
						BatchFulfillmentGasMultiplier: 1,
						BackoffInitialDelay:           time.Minute,
						BackoffMaxDelay:               time.Hour,
						GasLanePrice:                  assets.GWei(200),
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on VRFSpec {
									coordinatorAddress
									createdAt
									evmChainID
									fromAddresses
									minIncomingConfirmations
									pollPeriod
									publicKey
									requestedConfsDelay
									requestTimeout
									batchCoordinatorAddress
									batchFulfillmentEnabled
									batchFulfillmentGasMultiplier
									customRevertsPipelineEnabled
									chunkSize
									backoffInitialDelay
									backoffMaxDelay
									gasLanePrice
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "VRFSpec",
							"coordinatorAddress": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"createdAt": "2021-01-01T00:00:00Z",
							"evmChainID": "42",
							"fromAddresses": ["0x3cCad4715152693fE3BC4460591e3D3Fbd071b42", "0x2301958F1BFbC9A068C2aC9c6166Bf483b95864C"],
							"minIncomingConfirmations": 1,
							"pollPeriod": "1m0s",
							"publicKey": "0x9dc09a0f898f3b5e8047204e7ce7e44b587920932f08431e29c9bf6923b8450a01",
							"requestedConfsDelay": 10,
							"requestTimeout": "24h0m0s",
							"batchCoordinatorAddress": "0x0ad9FE7a58216242a8475ca92F222b0640E26B63",
							"batchFulfillmentEnabled": true,
							"batchFulfillmentGasMultiplier": 1,
							"customRevertsPipelineEnabled": true,
							"chunkSize": 25,
							"backoffInitialDelay": "1m0s",
							"backoffMaxDelay": "1h0m0s",
							"gasLanePrice": "200 gwei"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_WebhookSpec(t *testing.T) {
	var (
		id = int32(1)
	)

	testCases := []GQLTestCase{
		{
			name:          "webhook spec",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.Webhook,
					WebhookSpec: &job.WebhookSpec{
						CreatedAt: f.Timestamp(),
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on WebhookSpec {
									createdAt
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "WebhookSpec",
							"createdAt": "2021-01-01T00:00:00Z"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_BlockhashStoreSpec(t *testing.T) {
	var (
		id = int32(1)
	)
	coordinatorV1Address, err := evmtypes.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	coordinatorV2Address, err := evmtypes.NewEIP55Address("0x2fcA960AF066cAc46085588a66dA2D614c7Cd337")
	require.NoError(t, err)

	coordinatorV2PlusAddress, err := evmtypes.NewEIP55Address("0x92B5e28Ac583812874e4271380c7d070C5FB6E6b")
	require.NoError(t, err)

	fromAddress1, err := evmtypes.NewEIP55Address("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	require.NoError(t, err)

	fromAddress2, err := evmtypes.NewEIP55Address("0xD479d7c994D298cA05bF270136ED9627b7E684D3")
	require.NoError(t, err)

	blockhashStoreAddress, err := evmtypes.NewEIP55Address("0xb26A6829D454336818477B946f03Fb21c9706f3A")
	require.NoError(t, err)

	trustedBlockhashStoreAddress, err := evmtypes.NewEIP55Address("0x0ad9FE7a58216242a8475ca92F222b0640E26B63")
	require.NoError(t, err)
	trustedBlockhashStoreBatchSize := int32(20)

	testCases := []GQLTestCase{
		{
			name:          "blockhash store spec",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.BlockhashStore,
					BlockhashStoreSpec: &job.BlockhashStoreSpec{
						CoordinatorV1Address:           &coordinatorV1Address,
						CoordinatorV2Address:           &coordinatorV2Address,
						CoordinatorV2PlusAddress:       &coordinatorV2PlusAddress,
						CreatedAt:                      f.Timestamp(),
						EVMChainID:                     ubig.NewI(42),
						FromAddresses:                  []evmtypes.EIP55Address{fromAddress1, fromAddress2},
						PollPeriod:                     1 * time.Minute,
						RunTimeout:                     37 * time.Second,
						WaitBlocks:                     100,
						LookbackBlocks:                 200,
						HeartbeatPeriod:                450 * time.Second,
						BlockhashStoreAddress:          blockhashStoreAddress,
						TrustedBlockhashStoreAddress:   &trustedBlockhashStoreAddress,
						TrustedBlockhashStoreBatchSize: trustedBlockhashStoreBatchSize,
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on BlockhashStoreSpec {
									coordinatorV1Address
									coordinatorV2Address
									coordinatorV2PlusAddress
									createdAt
									evmChainID
									fromAddresses
									pollPeriod
									runTimeout
									waitBlocks
									lookbackBlocks
									blockhashStoreAddress
									trustedBlockhashStoreAddress
									trustedBlockhashStoreBatchSize
									heartbeatPeriod
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "BlockhashStoreSpec",
							"coordinatorV1Address": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"coordinatorV2Address": "0x2fcA960AF066cAc46085588a66dA2D614c7Cd337",
							"coordinatorV2PlusAddress": "0x92B5e28Ac583812874e4271380c7d070C5FB6E6b",
							"createdAt": "2021-01-01T00:00:00Z",
							"evmChainID": "42",
							"fromAddresses": ["0x3cCad4715152693fE3BC4460591e3D3Fbd071b42", "0xD479d7c994D298cA05bF270136ED9627b7E684D3"],
							"pollPeriod": "1m0s",
							"runTimeout": "37s",
							"waitBlocks": 100,
							"lookbackBlocks": 200,
							"blockhashStoreAddress": "0xb26A6829D454336818477B946f03Fb21c9706f3A",
							"trustedBlockhashStoreAddress": "0x0ad9FE7a58216242a8475ca92F222b0640E26B63",
							"trustedBlockhashStoreBatchSize": 20,
							"heartbeatPeriod": "7m30s"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_BlockHeaderFeederSpec(t *testing.T) {
	var (
		id = int32(1)
	)
	coordinatorV1Address, err := evmtypes.NewEIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C")
	require.NoError(t, err)

	coordinatorV2Address, err := evmtypes.NewEIP55Address("0x2fcA960AF066cAc46085588a66dA2D614c7Cd337")
	require.NoError(t, err)

	coordinatorV2PlusAddress, err := evmtypes.NewEIP55Address("0x92B5e28Ac583812874e4271380c7d070C5FB6E6b")
	require.NoError(t, err)

	fromAddress, err := evmtypes.NewEIP55Address("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	require.NoError(t, err)

	blockhashStoreAddress, err := evmtypes.NewEIP55Address("0xb26A6829D454336818477B946f03Fb21c9706f3A")
	require.NoError(t, err)

	batchBHSAddress, err := evmtypes.NewEIP55Address("0xd23BAE30019853Caf1D08b4C03291b10AD7743Df")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		{
			name:          "block header feeder spec",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.BlockHeaderFeeder,
					BlockHeaderFeederSpec: &job.BlockHeaderFeederSpec{
						CoordinatorV1Address:       &coordinatorV1Address,
						CoordinatorV2Address:       &coordinatorV2Address,
						CoordinatorV2PlusAddress:   &coordinatorV2PlusAddress,
						CreatedAt:                  f.Timestamp(),
						EVMChainID:                 ubig.NewI(42),
						FromAddresses:              []evmtypes.EIP55Address{fromAddress},
						PollPeriod:                 1 * time.Minute,
						RunTimeout:                 37 * time.Second,
						WaitBlocks:                 100,
						LookbackBlocks:             200,
						BlockhashStoreAddress:      blockhashStoreAddress,
						BatchBlockhashStoreAddress: batchBHSAddress,
						GetBlockhashesBatchSize:    5,
						StoreBlockhashesBatchSize:  3,
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on BlockHeaderFeederSpec {
									coordinatorV1Address
									coordinatorV2Address
									coordinatorV2PlusAddress
									createdAt
									evmChainID
									fromAddresses
									pollPeriod
									runTimeout
									waitBlocks
									lookbackBlocks
									blockhashStoreAddress
									batchBlockhashStoreAddress
									getBlockhashesBatchSize
									storeBlockhashesBatchSize
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "BlockHeaderFeederSpec",
							"coordinatorV1Address": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"coordinatorV2Address": "0x2fcA960AF066cAc46085588a66dA2D614c7Cd337",
							"coordinatorV2PlusAddress": "0x92B5e28Ac583812874e4271380c7d070C5FB6E6b",
							"createdAt": "2021-01-01T00:00:00Z",
							"evmChainID": "42",
							"fromAddresses": ["0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"],
							"pollPeriod": "1m0s",
							"runTimeout": "37s",
							"waitBlocks": 100,
							"lookbackBlocks": 200,
							"blockhashStoreAddress": "0xb26A6829D454336818477B946f03Fb21c9706f3A",
							"batchBlockhashStoreAddress": "0xd23BAE30019853Caf1D08b4C03291b10AD7743Df",
							"getBlockhashesBatchSize": 5,
							"storeBlockhashesBatchSize": 3
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_BootstrapSpec(t *testing.T) {
	var (
		id = int32(1)
	)

	testCases := []GQLTestCase{
		{
			name:          "Bootstrap spec",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.Bootstrap,
					BootstrapSpec: &job.BootstrapSpec{
						ID:                                id,
						ContractID:                        "0x613a38AC1659769640aaE063C651F48E0250454C",
						Relay:                             "evm",
						RelayConfig:                       map[string]interface{}{},
						MonitoringEndpoint:                null.String{},
						BlockchainTimeout:                 models.Interval(2 * time.Minute),
						ContractConfigTrackerPollInterval: models.Interval(2 * time.Minute),
						ContractConfigConfirmations:       100,
						CreatedAt:                         f.Timestamp(),
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on BootstrapSpec {
									id
									contractID
									relay
									relayConfig
									monitoringEndpoint
									blockchainTimeout
									contractConfigTrackerPollInterval
									contractConfigConfirmations
									createdAt
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "BootstrapSpec",
							"id": "1",
							"contractID": "0x613a38AC1659769640aaE063C651F48E0250454C",
							"relay": "evm",
							"relayConfig": {},
							"monitoringEndpoint": null,
							"blockchainTimeout": "2m0s",
							"contractConfigTrackerPollInterval": "2m0s",
							"contractConfigConfirmations": 100,
							"createdAt": "2021-01-01T00:00:00Z"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_WorkflowSpec(t *testing.T) {
	var (
		id = int32(1)
	)

	testCases := []GQLTestCase{
		{
			name:          "Workflow spec",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.Workflow,
					WorkflowSpec: &job.WorkflowSpec{
						ID:         id,
						WorkflowID: "<test workflow id>",
						Workflow:   "<test workflow spec>",
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on WorkflowSpec {
									id
									workflowID
									workflow
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "WorkflowSpec",
							"id": "1",
							"workflowID": "<test workflow id>",
							"workflow": "<test workflow spec>"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_GatewaySpec(t *testing.T) {
	var (
		id = int32(1)
	)

	gatewayConfig := map[string]interface{}{
		"NodeServerConfig": map[string]interface{}{},
	}

	testCases := []GQLTestCase{
		{
			name:          "Gateway spec",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.Gateway,
					GatewaySpec: &job.GatewaySpec{
						ID:            id,
						GatewayConfig: gatewayConfig,
						CreatedAt:     f.Timestamp(),
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on GatewaySpec {
									id
									gatewayConfig
									createdAt
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "GatewaySpec",
							"id": "1",
							"gatewayConfig": {"NodeServerConfig": {}},
							"createdAt": "2021-01-01T00:00:00Z"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_StandardCapabilitiesSpec(t *testing.T) {
	var (
		id = int32(1)
	)

	testCases := []GQLTestCase{
		{
			name:          "StandardCapabilities spec",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					Type: job.StandardCapabilities,
					StandardCapabilitiesSpec: &job.StandardCapabilitiesSpec{
						ID:        id,
						CreatedAt: f.Timestamp(),
						Command:   "testcommand",
						Config:    "testconfig",
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							spec {
								__typename
								... on StandardCapabilitiesSpec {
									id
									createdAt
									command
									config
								}
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"spec": {
							"__typename": "StandardCapabilitiesSpec",
							"id": "1",
							"createdAt": "2021-01-01T00:00:00Z",
							"command": "testcommand",
							"config": "testconfig"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_StreamSpec(t *testing.T) {
	var (
		id1      = int32(1)
		id2      = int32(2)
		streamID = uint32(3)
	)

	testCases := []GQLTestCase{
		{
			name:          "stream spec with stream ID",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id1).Return(job.Job{
					Type:     job.Stream,
					StreamID: &streamID,
				}, nil)
			},
			query: fmt.Sprintf(`
				query GetJob {
					job(id: "%d") {
						... on Job {
							spec {
								__typename
								... on StreamSpec {
									streamID
								}
							}
						}
					}
				}
			`, id1),
			result: fmt.Sprintf(`
				{
					"job": {
						"spec": {
							"__typename": "StreamSpec",
							"streamID": "%d"
						}
					}
				}
			`, streamID),
		},
		{
			name:          "stream spec without stream ID",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id2).Return(job.Job{
					Type: job.Stream,
				}, nil)
			},
			query: fmt.Sprintf(`
				query GetJob {
					job(id: "%d") {
						... on Job {
							spec {
								__typename
								... on StreamSpec {
									streamID
								}
							}
						}
					}
				}
			`, id2),
			result: `
				{
					"job": {
						"spec": {
							"__typename": "StreamSpec",
							"streamID": null
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

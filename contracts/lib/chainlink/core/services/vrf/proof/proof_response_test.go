package proof_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/solidity_vrf_verifier_wrapper"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	proof2 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
)

func TestMarshaledProof(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	key := cltest.DefaultVRFKey
	require.NoError(t, keyStore.VRF().Add(ctx, key))
	blockHash := common.Hash{}
	blockNum := 0
	preSeed := big.NewInt(1)
	s := proof2.TestXXXSeedData(t, preSeed, blockHash, blockNum)
	proofResponse, err := proof2.GenerateProofResponse(keyStore.VRF(), key.ID(), s)
	require.NoError(t, err)
	goProof, err := proof2.UnmarshalProofResponse(proofResponse)
	require.NoError(t, err)
	actualProof, err := goProof.CryptoProof(s)
	require.NoError(t, err)
	proof, err := proof2.MarshalForSolidityVerifier(&actualProof)
	require.NoError(t, err)
	// NB: For changes to the VRF solidity code to be reflected here, "go generate"
	// must be run in core/services/vrf.
	ethereumKey, _ := crypto.GenerateKey()
	auth, err := bind.NewKeyedTransactorWithChainID(ethereumKey, big.NewInt(1337))
	require.NoError(t, err)
	genesisData := gethtypes.GenesisAlloc{auth.From: {Balance: assets.Ether(100).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	_, _, verifier, err := solidity_vrf_verifier_wrapper.DeployVRFTestHelper(auth, backend.Client())
	if err != nil {
		panic(errors.Wrapf(err, "while initializing EVM contract wrapper"))
	}
	backend.Commit()
	_, err = verifier.RandomValueFromVRFProof(nil, proof[:])
	require.NoError(t, err)
}

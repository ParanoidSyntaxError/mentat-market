package mercury

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func Test_OffchainConfigDigester_ConfigDigest(t *testing.T) {
	ctx := t.Context()
	// ChainID and ContractAddress are taken into account for computation
	cd1, err := OffchainConfigDigester{ChainID: big.NewInt(0)}.ConfigDigest(ctx, types.ContractConfig{})
	require.NoError(t, err)
	cd2, err := OffchainConfigDigester{ChainID: big.NewInt(0)}.ConfigDigest(ctx, types.ContractConfig{})
	require.NoError(t, err)
	cd3, err := OffchainConfigDigester{ChainID: big.NewInt(1)}.ConfigDigest(ctx, types.ContractConfig{})
	require.NoError(t, err)
	cd4, err := OffchainConfigDigester{ChainID: big.NewInt(1), ContractAddress: common.Address{1}}.ConfigDigest(ctx, types.ContractConfig{})
	require.NoError(t, err)

	require.Equal(t, cd1, cd2)
	require.NotEqual(t, cd2, cd3)
	require.NotEqual(t, cd2, cd4)
	require.NotEqual(t, cd3, cd4)

	// malformed signers
	_, err = OffchainConfigDigester{}.ConfigDigest(ctx, types.ContractConfig{
		Signers: []types.OnchainPublicKey{{1, 2}},
	})
	require.Error(t, err)

	// malformed transmitters
	_, err = OffchainConfigDigester{}.ConfigDigest(ctx, types.ContractConfig{
		Transmitters: []types.Account{"0x"},
	})
	require.Error(t, err)

	_, err = OffchainConfigDigester{}.ConfigDigest(ctx, types.ContractConfig{
		Transmitters: []types.Account{"7343581f55146951b0f678dc6cfa8fd360e2f353"},
	})
	require.Error(t, err)

	_, err = OffchainConfigDigester{}.ConfigDigest(ctx, types.ContractConfig{
		Transmitters: []types.Account{"7343581f55146951b0f678dc6cfa8fd360e2f353aabbccddeeffaaccddeeffaz"},
	})
	require.Error(t, err)

	// well-formed transmitters
	_, err = OffchainConfigDigester{ChainID: big.NewInt(0)}.ConfigDigest(ctx, types.ContractConfig{
		Transmitters: []types.Account{"7343581f55146951b0f678dc6cfa8fd360e2f353aabbccddeeffaaccddeeffaa"},
	})
	require.NoError(t, err)
}

package ocr2key

import (
	"bytes"
	cryptorand "crypto/rand"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestEVMKeyring_SignVerify(t *testing.T) {
	kr1, err := newEVMKeyring(cryptorand.Reader)
	require.NoError(t, err)
	kr2, err := newEVMKeyring(cryptorand.Reader)
	require.NoError(t, err)

	ctx := ocrtypes.ReportContext{}

	t.Run("can verify", func(t *testing.T) {
		report := ocrtypes.Report{}
		sig, err := kr1.Sign(ctx, report)
		require.NoError(t, err)
		t.Log(len(sig))
		result := kr2.Verify(kr1.PublicKey(), ctx, report, sig)
		assert.True(t, result)
	})

	t.Run("invalid sig", func(t *testing.T) {
		report := ocrtypes.Report{}
		result := kr2.Verify(kr1.PublicKey(), ctx, report, []byte{0x01})
		assert.False(t, result)
	})

	t.Run("invalid pubkey", func(t *testing.T) {
		report := ocrtypes.Report{}
		sig, err := kr1.Sign(ctx, report)
		require.NoError(t, err)
		result := kr2.Verify([]byte{0x01}, ctx, report, sig)
		assert.False(t, result)
	})
}

func TestEVMKeyring_Sign3Verify3(t *testing.T) {
	kr1, err := newEVMKeyring(cryptorand.Reader)
	require.NoError(t, err)
	kr2, err := newEVMKeyring(cryptorand.Reader)
	require.NoError(t, err)

	digest, err := types.BytesToConfigDigest(testutils.MustRandBytes(32))
	require.NoError(t, err)
	seqNr := rand.Uint64()
	r := ocrtypes.Report(testutils.MustRandBytes(rand.Intn(1024)))

	t.Run("can verify", func(t *testing.T) {
		sig, err := kr1.Sign3(digest, seqNr, r)
		require.NoError(t, err)
		t.Log(len(sig))
		result := kr2.Verify3(kr1.PublicKey(), digest, seqNr, r, sig)
		assert.True(t, result)
	})

	t.Run("invalid sig", func(t *testing.T) {
		result := kr2.Verify3(kr1.PublicKey(), digest, seqNr, r, []byte{0x01})
		assert.False(t, result)
	})

	t.Run("invalid pubkey", func(t *testing.T) {
		sig, err := kr1.Sign3(digest, seqNr, r)
		require.NoError(t, err)
		result := kr2.Verify3([]byte{0x01}, digest, seqNr, r, sig)
		assert.False(t, result)
	})
}

func TestEVMKeyring_Marshalling(t *testing.T) {
	kr1, err := newEVMKeyring(cryptorand.Reader)
	require.NoError(t, err)

	m, err := kr1.Marshal()
	require.NoError(t, err)

	kr2 := evmKeyring{}
	err = kr2.Unmarshal(m)
	require.NoError(t, err)

	assert.True(t, bytes.Equal(kr1.PublicKey(), kr2.PublicKey()))
	assert.True(t, bytes.Equal(kr1.privateKey().D.Bytes(), kr2.privateKey().D.Bytes()))

	// Invalid seed size should error
	assert.Error(t, kr2.Unmarshal([]byte{0x01}))
}

func TestRawReportContext3(t *testing.T) {
	testCases := []struct {
		name     string
		digest   [32]byte
		seqNr    uint64
		expected [2][32]byte
	}{
		{
			name:   "zero values",
			digest: [32]byte{},
			seqNr:  0,
			expected: [2][32]byte{
				{},
				{},
			},
		},
		{
			name:   "some digest",
			digest: [32]byte{1, 2, 3},
			seqNr:  0,
			expected: [2][32]byte{
				{1, 2, 3},
				{},
			},
		},
		{
			name:   "sequence number set to 1",
			digest: [32]byte{},
			seqNr:  1,
			expected: [2][32]byte{
				{},
				{
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 1,
				},
			},
		},
		{
			name:   "sequence number set to max uint64",
			digest: [32]byte{1, 2, 3},
			seqNr:  math.MaxUint64,
			expected: [2][32]byte{
				{1, 2, 3},
				{
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := RawReportContext3(tc.digest, tc.seqNr)
			assert.Equal(t, tc.expected, actual, "unexpected result")
		})
	}
}

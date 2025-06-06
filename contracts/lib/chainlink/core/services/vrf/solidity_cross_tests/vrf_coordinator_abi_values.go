package solidity_cross_tests

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/solidity_vrf_coordinator_interface"
)

// VRFRandomnessRequestLogTopic returns the signature of the RandomnessRequest log
// emitted by the VRFCoordinator contract
func VRFRandomnessRequestLogTopic() common.Hash {
	return coordinatorABIValues().randomnessRequestLogTopic
}

// randomnessRequestRawDataArgs returns a list of the arguments to the
// RandomnessRequest log emitted by the VRFCoordinator contract
func randomnessRequestRawDataArgs() abi.Arguments {
	return coordinatorABIValues().randomnessRequestRawDataArgs
}

var fulfillMethodName = "fulfillRandomnessRequest"

// abiValues is a singleton carrying information parsed once from the
// VRFCoordinator abi string
type abiValues struct {
	// CoordinatorABI is the ABI of the VRFCoordinator
	coordinatorABI  abi.ABI
	fulfillSelector string
	fulfillMethod   abi.Method
	// RandomnessRequestLogTopic is the signature of the RandomnessRequest log
	randomnessRequestLogTopic    common.Hash
	randomnessRequestRawDataArgs abi.Arguments
}

var coordinatorABIValues = sync.OnceValue(func() (v *abiValues) {
	v = new(abiValues)
	var err error
	v.coordinatorABI, err = abi.JSON(strings.NewReader(
		solidity_vrf_coordinator_interface.VRFCoordinatorABI))
	if err != nil {
		panic(err)
	}
	var found bool
	v.fulfillMethod, found = v.coordinatorABI.Methods[fulfillMethodName]
	if !found {
		panic(fmt.Errorf("could not find method %s in VRFCoordinator ABI", fulfillMethodName))
	}
	v.fulfillSelector = hexutil.Encode(v.fulfillMethod.ID)
	randomnessRequestABI := v.coordinatorABI.Events["RandomnessRequest"]
	v.randomnessRequestLogTopic = randomnessRequestABI.ID
	for _, arg := range randomnessRequestABI.Inputs {
		if !arg.Indexed {
			v.randomnessRequestRawDataArgs = append(v.randomnessRequestRawDataArgs, arg)
		}
	}
	return
})

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";

import {FunctionsClient} from "@chainlink/contracts/src/v0.8/functions/v1_0_0/FunctionsClient.sol";

import {IMentatMarket} from "./IMentatMarket.sol";
import {MM20, IMM20} from "./MM20.sol";

contract MentatMarket is FunctionsClient, IMentatMarket {
    address public immutable mm20Implementation;

    mapping(bytes32 => Market) private _markets;

    constructor(address functionsRouter) FunctionsClient(functionsRouter) {
        mm20Implementation = address(new MM20());
    }

    function createMarket(
        bytes32 id,
        string memory name,
        string memory symbol,
        MarketParams memory params
    ) external override {
        if (_marketExists(id)) {
            revert MarketExists(id);
        }

        if (params.expiry <= block.timestamp) {
            revert MarketExpired(params.expiry, block.timestamp);
        }

        address yes = Clones.clone(mm20Implementation);
        IMM20(yes).initialize(
            string.concat(name, " (yes)"),
            string.concat(symbol, "yes"),
            address(this)
        );

        address no = Clones.clone(mm20Implementation);
        IMM20(no).initialize(
            string.concat(name, " (no)"),
            string.concat(symbol, "no"),
            address(this)
        );

        _markets[id] = Market(yes, no, params);
    }

    function fulfillRequest(
        bytes32 requestId,
        bytes memory response,
        bytes memory err
    ) internal override {
        // TODO: Implement
    }

    function _marketExists(bytes32 id) internal view returns (bool) {
        return _markets[id].yes != address(0);
    }
}

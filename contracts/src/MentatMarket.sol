// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";

import {FunctionsClient} from "@chainlink/contracts/src/v0.8/functions/v1_0_0/FunctionsClient.sol";

import {IMentatMarket} from "./IMentatMarket.sol";
import {MM20, IMM20} from "./MM20.sol";
import {FunctionsConsumer} from "./FunctionsConsumer.sol";

contract MentatMarket is FunctionsClient, FunctionsConsumer, IMentatMarket {
    address public immutable mm20Implementation;

    address public immutable liquidityToken;

    mapping(bytes32 => Market) private _markets;

    constructor(
        address functionsRouter,
        address liquidity
    ) FunctionsClient(functionsRouter) FunctionsConsumer(functionsRouter) {
        mm20Implementation = address(new MM20());
        liquidityToken = liquidity;
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

        _markets[id] = Market(yes, no, false, false, params);
    }

    function mint(
        bytes32 id,
        bool outcome,
        uint256 liquidityAmount,
        uint256 minMintAmount,
        address receiver
    ) external {
        if (!_marketExists(id)) {
            revert MarketDoesNotExist(id);
        }

        IERC20(liquidityToken).transferFrom(
            msg.sender,
            address(this),
            liquidityAmount
        );

        address token = outcome ? _markets[id].yes : _markets[id].no;
        uint256 mintAmount = _mintAmount(token, liquidityAmount);

        if (mintAmount < minMintAmount) {
            revert SlippageExceeded(mintAmount, minMintAmount);
        }

        IMM20(token).mint(receiver, mintAmount);
    }

    function request(
        bytes32 id,
        uint64 subscriptionId,
        uint32 callbackGasLimit,
        bytes32 donId
    )
        external
        onlySubscriptionConsumer(subscriptionId, msg.sender)
        returns (bytes32)
    {
        if (!_marketExists(id)) {
            revert MarketDoesNotExist(id);
        }

        if (_marketResolved(id)) {
            revert MarketResolved(id);
        }

        if (_markets[id].params.expiry < block.timestamp) {
            _markets[id].resolved = true;
            return bytes32(0);
        }

        bytes32 requestId = _sendRequest(
            _markets[id].params.request,
            subscriptionId,
            callbackGasLimit,
            donId
        );

        return requestId;
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

    function _marketResolved(bytes32 id) internal view returns (bool) {
        return _markets[id].resolved;
    }

    function _mintAmount(
        address token,
        uint256 amount
    ) internal view returns (uint256) {
        // TODO: Implement
        return 1;
    }
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {IFunctionsConsumer} from "./IFunctionsConsumer.sol";

interface IMentatMarket is IFunctionsConsumer {
    error MarketExists(bytes32 id);
    error MarketDoesNotExist(bytes32 id);
    error MarketExpired(uint256 expiry, uint256 blockTimestamp);
    error MarketResolved(bytes32 id);
    error SlippageExceeded(uint256 mintAmount, uint256 minOutAmount);

    struct MarketParams {
        bytes request;
        string metadataUri;
        address creator;
        uint256 expiry;
    }

    struct Market {
        address yes;
        address no;
        bool resolved;
        bool outcome;
        MarketParams params;
    }

    function createMarket(
        bytes32 id,
        string memory name,
        string memory symbol,
        MarketParams memory params
    ) external;
}

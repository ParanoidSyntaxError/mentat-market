// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

interface IMentatMarket {
    error MarketExists(bytes32 id);
    error MarketExpired(uint256 expiry, uint256 blockTimestamp);

    struct MarketParams {
        string metadataUri;
        address creator;
        uint256 expiry;
    }

    struct Market {
        address yes;
        address no;
        MarketParams params;
    }

    function createMarket(
        bytes32 id,
        string memory name,
        string memory symbol,
        MarketParams memory params
    ) external;
}

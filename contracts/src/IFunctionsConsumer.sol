// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

interface IFunctionsConsumer {
    function isSubscriptionPublic(
        uint64 subscriptionId
    ) external view returns (bool);

    function isSubscriptionConsumer(
        uint64 subscriptionId,
        address account
    ) external view returns (bool);

    function setSubscriptionConsumers(
        uint64 subscriptionId,
        address[] memory accounts,
        bool[] memory values
    ) external;
}

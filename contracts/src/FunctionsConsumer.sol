// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {IFunctionsSubscriptions} from "@chainlink/contracts/src/v0.8/functions/v1_0_0/interfaces/IFunctionsSubscriptions.sol";

import {IFunctionsConsumer} from "./IFunctionsConsumer.sol";

contract FunctionsConsumer is IFunctionsConsumer {
    address public immutable functionsRouter;

    mapping(uint64 => mapping(address => bool)) private _consumers;

    address public constant PUBLIC_SUBSCRIPTION = address(0);

    constructor(address router) {
        functionsRouter = router;
    }

    modifier onlySubscriptionConsumer(uint64 subscriptionId, address account) {
        if (!isSubscriptionConsumer(subscriptionId, account)) {
            revert();
        }
        _;
    }

    function isSubscriptionPublic(
        uint64 subscriptionId
    ) public view override returns (bool) {
        return _consumers[subscriptionId][PUBLIC_SUBSCRIPTION] == true;
    }

    function isSubscriptionConsumer(
        uint64 subscriptionId,
        address account
    ) public view override returns (bool) {
        return
            _consumers[subscriptionId][account] == true ||
            isSubscriptionPublic(subscriptionId);
    }

    function setSubscriptionConsumers(
        uint64 subscriptionId,
        address[] memory accounts,
        bool[] memory values
    ) external override {
        IFunctionsSubscriptions.Subscription
            memory subscription = IFunctionsSubscriptions(functionsRouter)
                .getSubscription(subscriptionId);

        if (subscription.owner != msg.sender) {
            revert();
        }

        if (accounts.length != values.length) {
            revert();
        }

        for (uint256 i = 0; i < accounts.length; i++) {
            _consumers[subscriptionId][accounts[i]] = values[i];
        }
    }
}

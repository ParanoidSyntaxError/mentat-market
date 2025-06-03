// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {IERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol";
import {IERC1363} from "@openzeppelin/contracts/interfaces/IERC1363.sol";

interface IMM20 is IERC20, IERC20Metadata, IERC20Permit, IERC1363 {
    function initialize(
        string memory initName,
        string memory initSymbol,
        address initOwner
    ) external;

    function mint(address to, uint256 amount) external;

    function burn(uint256 amount) external;
    function burnFrom(address from, uint256 amount) external;
}

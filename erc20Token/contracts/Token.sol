pragma solidity ^0.4.24;

import "./ERC20.sol";
import "./Owned.sol";

contract Token is ERC20Interface, Owned {

  mapping (address => uint256) public balances;
  mapping (address => mapping (address => uint256)) public allowed;
  string public name;
  uint8 public decimals;
  string public symbol;
  uint256 constant private MAX_UINT256 = 2**256 - 1;
  mapping(address => bool) private accountGroup;
  
  // Event to notify clients for AccountGroup operation.
  // Permit is false means freeze the account.
  event AccountGroupOp(address target, bool permit);
  
  constructor(
	      uint256 _initialAmount,
	      string _tokenName,
	      uint8 _decimals,
	      string _tokenSymbol
	      ) public {
    totalSupply = _initialAmount;
    balances[msg.sender] = totalSupply; // All init to owner.
    accountGroup[msg.sender] = true;
    name = _tokenName;
    decimals = _decimals;
    symbol = _tokenSymbol;
  }

  
  function transfer(address _to, uint256 _value) public returns (bool success) {
    require(balances[msg.sender] >= _value);
    require(accountGroup[msg.sender] && accountGroup[_to]);
    balances[msg.sender] -= _value;
    balances[_to] += _value;
    emit Transfer(msg.sender, _to, _value);
    return true;
  }


  function transferFrom(address _from, address _to, uint256 _value) public returns (bool success) {
    uint256 allowance = allowed[_from][msg.sender];
    require(balances[_from] >= _value && allowance >= _value);
    require(accountGroup[_from] && accountGroup[_to]);
    balances[_to] += _value;
    balances[_from] -= _value;
    if (allowance < MAX_UINT256) {
      allowed[_from][msg.sender] -= _value;
    }
    emit Transfer(_from, _to, _value);
    return true;
  }

  
  function balanceOf(address _owner) public view returns (uint256 balance) {
    return balances[_owner];
  }

  
  function approve(address _spender, uint256 _value) public returns (bool success) {
    require(accountGroup[msg.sender] && accountGroup[_spender]);
    allowed[msg.sender][_spender] = _value;
    emit Approval(msg.sender, _spender, _value);
    return true;
  }

  
  function allowance(address _owner, address _spender) public view returns (uint256 remaining) {
    return allowed[_owner][_spender];
  }

  
  function permit(address target) onlyOwner public {
    accountGroup[target] = true;
    emit AccountGroupOp(target, true);
  }
  
  function freeze(address target) onlyOwner public {
    accountGroup[target] = false;
    emit AccountGroupOp(target, false);
  }
  
}

pragma solidity ^0.4.24;

contract MyTokenERC20 {
    // Public variables of the token
    string public name;
    string public symbol;
    uint8 public decimals;
    uint256 public totalSupply;

    address public owner;

    // This creates an array with all balances
    mapping (address => uint256) public balanceOf;

    modifier onlyOwner { require(msg.sender == owner);_; }

    // This generates a public event on the blockchain that will notify clients
    event Transfer(address indexed from, address indexed to, uint256 value);

    // This notifies clients about the amount burnt
    event Burn(address indexed from, uint256 value);

    /**
     * Constrctor function
     *
     * Initializes contract with initial supply tokens to the creator of the contract
     */
    //function MyTokenERC20 (
    constructor (
        uint256 initialSupply,
        string tokenName,
        string tokenSymbol,
		uint8 tokendecimals
    ) public {
		
		owner = msg.sender;
	
        totalSupply = initialSupply ;  						// Update total supply
        balanceOf[owner] = totalSupply;                // Give the creator all initial tokens
        name = tokenName;                                   // Set the name for display purposes
        symbol = tokenSymbol;                               // Set the symbol for display purposes
        decimals = tokendecimals;
    }

    /**
     * Internal transfer, only can be called by this contract
     */
    function _transfer(address _from, address _to, uint _value) internal {
        // Prevent transfer to 0x0 address. Use burn() instead
        require(_to != 0x0);
        // Check if the sender has enough
        require(balanceOf[_from] >= _value);
        // Check for overflows uint256类型 2^256
        require(balanceOf[_to] + _value > balanceOf[_to]);
        // Save this for an assertion in the future
        uint previousBalances = balanceOf[_from] + balanceOf[_to];
        // Subtract from the sender
        balanceOf[_from] -= _value;
        // Add the same to the recipient
        balanceOf[_to] += _value;
        emit Transfer(_from, _to, _value);
        // Asserts are used to use static analysis to find bugs in your code. They should never fail
        assert(balanceOf[_from] + balanceOf[_to] == previousBalances);
    }

    function transfer(address _to, uint256 _value) public returns (bool success){
        _transfer(owner, _to, _value);
		return true;
    }

    function transferFrom(address _from, address _to, uint256 _value) public returns (bool success){
        _transfer(_from, _to, _value);
		return true;
    }

    function burn(uint256 _value) onlyOwner public returns (bool success) {
        require(balanceOf[owner] >= _value);   // Check if the sender has enough 万一owner余额不足，此函数就失去了意义
        balanceOf[owner] -= _value;            // Subtract from the sender        是否需要指定地址进行销毁？？
        totalSupply -= _value;                      // Updates totalSupply
        emit Burn(owner, _value);
        return true;
    }

    function mintToken(uint256 mintedAmount) onlyOwner public returns (bool success) {
        balanceOf[owner] += mintedAmount;
        totalSupply += mintedAmount;
        emit Transfer(0, owner, mintedAmount);
		return true;
    }
}

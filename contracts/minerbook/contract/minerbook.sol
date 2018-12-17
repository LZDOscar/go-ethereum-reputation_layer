pragma solidity 0.5.1;

contract minerbook {
    event MinerRegistered(
        bytes32 indexed hashedPubkey,
        address indexed withdrawalAddressbytes32,
        bytes32 indexed randaoCommitment
    );

    //TODO: address => (register_name, register_ID)
    mapping (bytes32 => bool) public usedHashedPubkey;

    uint public constant MINER_ADMISSION = 32 ether;
    //TODO：we assume the information(register_name, register_ID) that the registers sent are valid
    //TODO: because the contract checks this by database API which government offerd, but the function has not been achieved now.
    //TODO：one register can register miner with one address,so the function must check.
    function register(
        bytes memory _pubkey,
        address  _withdrawalAddressbytes32,
        bytes32  _randaoCommitment
    )
        public
        payable
    {
        require(
            msg.value == MINER_ADMISSION,
            "Incorrect miner admission"
        );
        require(
            _pubkey.length == 48,
            "Public key is not 48 bytes"
        );

        bytes32 hashedPubkey = keccak256(abi.encodePacked(_pubkey));
        //one address must be registerd once.
        require(
            !usedHashedPubkey[hashedPubkey],
            "Public key already used"
        );

        //TODO：check the register's info whether it is used

        usedHashedPubkey[hashedPubkey] = true;

        emit MinerRegistered(hashedPubkey, _withdrawalAddressbytes32, _randaoCommitment);
    }

}
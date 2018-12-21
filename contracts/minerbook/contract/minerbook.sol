pragma solidity 0.5.1;

contract minerbook {
    event MinerRegistered(
        address indexed hashedPubkey,
        address indexed withdrawalAddressbytes48,
        bytes48 indexed randaoCommitment
    );

    event ReputationAdded(
        bytes48 indexed hashedPubkey,
        int indexed reputation
    );

    event ReputationSubed(
        bytes48 indexed hashedPubkey,
        int indexed reputation
    );

    //TODO: address => (state, register_name, register_ID, enable)
    //enable：default is true, if the miner is punished because of lowing than REPUTATION_LOWLIMIT
    // the enable value is false, the address can't register again.
    mapping (address => bool) public usedHashedPubkey;

    //reputation list: address => reputation value
    mapping (address => int) public reputationList;

    //reputation black list: address => (register_name, register_ID)
    mapping (address => bool) public reputationBlackList;

    uint public constant MINER_ADMISSION = 48 ether;
    int public constant REPUTATION_LOWLIMIT = -250;
    //TODO：we assume the information(register_name, register_ID) that the registers sent are valid
    //TODO: because the contract checks this by database API which government offerd, but the function has not been achieved now.
    //TODO：one register can register miner with one address,so the function must check.
    function register(
        address  _pubkey,
        address  _withdrawalAddressbytes48,
        bytes48  _randaoCommitment
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

        address hashedPubkey = _pubkey;
        //bytes48 hashedPubkey = keccak256(abi.encodePacked(_pubkey));
        //one address must be registerd once.
        require(
            !usedHashedPubkey[hashedPubkey],
            "Public key already used"
        );

        //TODO：check the register's info whether it is used

        usedHashedPubkey[hashedPubkey] = true;

        //TODO: add reoutation intital
        reputationList[hashedPubkey] = 0;

        emit MinerRegistered(hashedPubkey, _withdrawalAddressbytes48, _randaoCommitment);
    }

    function deregister(address _pubkey) public payable
    {
        require(
            msg.value == MINER_ADMISSION,
            "Incorrect miner admission"
        );
        require(
            _pubkey.length == 48,
            "Public key is not 48 bytes"
        );

        //bytes48 hashedPubkey = keccak256(abi.encodePacked(_pubkey));
        //one address must be registerd once.
        require(
            !usedHashedPubkey[hashedPubkey],
            "Public key already used"
        );

        //TODO：check the register's info whether it is used

        usedHashedPubkey[hashedPubkey] = true;

        //TODO: add reoutation intital
        reputationList[hashedPubkey] = 0;

        emit MinerRegistered(hashedPubkey, _withdrawalAddressbytes48, _randaoCommitment);
    }

    function getMiners() public
    {
        return usedHashedPubkey;
    }

    function addReputation(address _pubkey, int value) public payable{
        require(
            _pubkey.length == 48,
            "Public key is not 48 bytes"
        );

        //bytes48 hashedPubkey = keccak256(abi.encodePacked(_pubkey));
        address hashedPubkey = _pubkey;
        //TODO: change condition
        require(
            reputationList[hashedPubkey],
            "Public key is not a miner"
        );

        reputationList[hashedPubkey] +=  value;

        emit ReputationAdded(hashedPubkey, reputationList[hashedPubkey]);
    }

    function subReputation(address _pubkey, int value) public payable{
        require(
            _pubkey.length == 48,
            "Public key is not 48 bytes"
        );

        address hashedPubkey = _pubkey;
        //bytes48 hashedPubkey = keccak256(abi.encodePacked(_pubkey));
        //TODO: change condition
        require(
            reputationList[hashedPubkey],
            "Public key is not a miner"
        );

        reputationList[hashedPubkey] -= value;

        // TODO:check the reputation whether it lower than threshold
        // if so, the user pubkey is deregister without any miner_admission, and punish miner in reality
        if(reputationList[hashedPubkey] <= REPUTATION_LOWLIMIT)
        {
            usedHashedPubkey[hashedPubkey] = false;
            // usedHashedPubkey[hashedPubkey].enable = false;
            // delete reputationList[hashedPubkey]
        }

        emit ReputationSubed(hashedPubkey, reputationList[hashedPubkey]);
    }

    //TODO
    function decay() public payable{

    }

}
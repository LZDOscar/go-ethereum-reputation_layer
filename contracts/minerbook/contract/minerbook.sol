pragma solidity 0.5.2;

library IterableMapping
{
  struct itmap
  {
    mapping(address => IndexValue) data;
    KeyFlag[] keys;
    uint size;
  }
  struct IndexValue { uint keyIndex; bool value; }
  struct KeyFlag { address key; bool deleted; }
  function insert(itmap storage self, address key, bool value) public returns (bool replaced)
  {
    uint keyIndex = self.data[key].keyIndex;
    self.data[key].value = value;
    if (keyIndex > 0)
      return true;
    else
    {
      keyIndex = self.keys.length++;
      self.data[key].keyIndex = keyIndex + 1;
      self.keys[keyIndex].key = key;
      self.size++;
      return false;
    }
  }
  function remove(itmap storage self, address key) public returns (bool success)
  {
    uint keyIndex = self.data[key].keyIndex;
    if (keyIndex == 0)
      return false;
    delete self.data[key];
    self.keys[keyIndex - 1].deleted = true;
    self.size --;
  }
  function contains(itmap storage self, address key)public  returns (bool)
  {
    return self.data[key].keyIndex > 0;
  }
  function iterate_start(itmap storage self)public returns (uint keyIndex)
  {
    return iterate_next(self, uint(-1));
  }
  function iterate_valid(itmap storage self, uint keyIndex)public returns (bool)
  {
    return keyIndex < self.keys.length;
  }
  function iterate_next(itmap storage self, uint keyIndex)public returns (uint r_keyIndex)
  {
    keyIndex++;
    while (keyIndex < self.keys.length && self.keys[keyIndex].deleted)
      keyIndex++;
    return keyIndex;
  }
  function iterate_get(itmap storage self, uint keyIndex)public returns (address key, bool value)
  {
    key = self.keys[keyIndex].key;
    value = self.data[key].value;
  }
}


contract minerbook {
    event MinerRegistered(
        address indexed hashedPubkey,
        address indexed withdrawalAddressbytes48
        // bytes48 indexed randaoCommitment
    );

    event MinerDeRegistered(
        address indexed hashedPubkey
        // address indexed withdrawalAddressbytes48
        // bytes48 indexed randaoCommitment
    );

    event ReputationAdded(
        address indexed hashedPubkey,
        uint indexed reputation
    );

    event ReputationSubed(
        address indexed hashedPubkey,
        uint indexed reputation
    );

    //TODO: address => (state, register_name, register_ID, enable)
    //enable：default is true, if the miner is punished because of lowing than REPUTATION_LOWLIMIT
    // the enable value is false, the address can't register again.
    // IterableMapping.itmap public usedHashedPubkey;
    mapping (address => bool) public usedHashedPubkey;
    address[] public regedAddrs;
    uint public regedAddrsLen = 0;
    mapping (address => address) public withdrawAddrs;

    //reputation list: address => reputation value
    mapping (address => uint) public reputationList;

    //reputation black list: address => (register_name, register_ID)
    mapping (address => bool) public reputationBlackList;


    uint public constant MINER_ADMISSION = 0 ether;
    uint public constant REPUTATION_LOWLIMIT = 0;
    uint public constant REPUTATION_HIGHLIMIT = 2000;
    uint public constant REPUTATION_INIT = 0;
    //TODO：we assume the information(register_name, register_ID) that the registers sent are valid
    //TODO: because the contract checks this by database API which government offerd, but the function has not been achieved now.
    //TODO：one register can register miner with one address,so the function must check.
    function register(
        address  _pubkey,
        address  _withdrawalAddressbytes48
        // bytes48  _randaoCommitment
    )
        public
        payable
    {
        require(
            msg.value == MINER_ADMISSION,
            "Incorrect miner admission"
        );
        // require(
        //     _pubkey  typeof(address),
        //     "Public key is not 48 bytes"
        // );

        address hashedPubkey = _pubkey;
        //bytes48 hashedPubkey = keccak256(abi.encodePacked(_pubkey));
        //one address must be registerd once.
        require(
            // !IterableMapping.contains(usedHashedPubkey,hashedPubkey),
            !usedHashedPubkey[hashedPubkey],
            "Public key already used"
        );

        //TODO：check the register's info whether it is used

        // IterableMapping.insert(usedHashedPubkey,hashedPubkey,true);
        usedHashedPubkey[hashedPubkey] = true;

        if(regedAddrsLen == regedAddrs.length){
            regedAddrs.push(hashedPubkey);
            regedAddrsLen ++;
        }
        else{
            regedAddrs[regedAddrsLen] = hashedPubkey;
            regedAddrsLen ++;
        }


        withdrawAddrs[hashedPubkey] = _withdrawalAddressbytes48;
        //TODO: add reoutation intital
        reputationList[hashedPubkey] = 0;

        emit MinerRegistered(hashedPubkey, _withdrawalAddressbytes48);
    }

    function deregister(address _pubkey) public payable
    {
        // require(
        //     _pubkey.length == 48,
        //     "Public key is not 48 bytes"
        // );
        require(
            msg.sender == _pubkey,
            "Incorrect miner admission"
        );
        address hashedPubkey = _pubkey;
        //bytes48 hashedPubkey = keccak256(abi.encodePacked(_pubkey));
        //one address must be registerd once.
        require(
            // IterableMapping.contains(usedHashedPubkey,hashedPubkey),
            usedHashedPubkey[hashedPubkey],
            "Public key is not used"
        );

        //TODO：check the register's info whether it is used
        // IterableMapping.remove(usedHashedPubkey,hashedPubkey);
        delete usedHashedPubkey[hashedPubkey];

        for(uint i = 0; i < regedAddrsLen;i++){
            if (regedAddrs[i] == hashedPubkey){
                if(regedAddrsLen == 1){
                    delete regedAddrs[i];
                    regedAddrsLen = 0;
                }
                else{
                    regedAddrs[i] = regedAddrs[regedAddrsLen-1];
                    delete regedAddrs[regedAddrsLen-1];
                    regedAddrsLen -= 1;
                }

                break;
            }
        }
        //TODO: add reoutation intital
        delete reputationList[hashedPubkey];

        emit MinerDeRegistered(hashedPubkey);
    }

    function getMiners()  public view  returns(address[] memory)
    {
        // address [] T = [];
    //   address[] T = [];
        // for (uint i = IterableMapping.iterate_start(usedHashedPubkey); IterableMapping.iterate_valid(usedHashedPubkey, i); i = IterableMapping.iterate_next(usedHashedPubkey, i))
        // {
        //     address key;
        //     bool value;
        //     (key, value) =IterableMapping.iterate_get(usedHashedPubkey, i);
        //     miners.push(key);
        // }
        // miners = regedAddrs;
        address[] memory miners = new address[](regedAddrsLen);
        for(uint i = 0; i < regedAddrsLen; i++){
            miners[i] = (regedAddrs[i]);
            // miners.length += 1;
        }
        return (miners);
     }

    // function addReputation(address _pubkey, uint value) public payable{
    //     require(
    //         _pubkey.length == 48,
    //         "Public key is not 48 bytes"
    //     );

    //     //bytes48 hashedPubkey = keccak256(abi.encodePacked(_pubkey));
    //     address hashedPubkey = _pubkey;
    //     //TODO: change condition
    //     require(
    //         reputationList[hashedPubkey],
    //         "Public key is not a miner"
    //     );

    //     reputationList[hashedPubkey] +=  value;

    //     emit ReputationAdded(hashedPubkey, reputationList[hashedPubkey]);
    // }

    // function subReputation(address _pubkey, uint value) public payable{
    //     require(
    //         _pubkey.length == 48,
    //         "Public key is not 48 bytes"
    //     );

    //     address hashedPubkey = _pubkey;
    //     //bytes48 hashedPubkey = keccak256(abi.encodePacked(_pubkey));
    //     //TODO: change condition
    //     require(
    //         reputationList[hashedPubkey],
    //         "Public key is not a miner"
    //     );

    //     reputationList[hashedPubkey] -= value;

    //     // TODO:check the reputation whether it lower than threshold
    //     // if so, the user pubkey is deregister without any miner_admission, and punish miner in reality
    //     if(reputationList[hashedPubkey] <= REPUTATION_LOWLIMIT)
    //     {
    //         usedHashedPubkey[hashedPubkey] = false;
    //         // usedHashedPubkey[hashedPubkey].enable = false;
    //         // delete reputationList[hashedPubkey]
    //     }

    //     emit ReputationSubed(hashedPubkey, reputationList[hashedPubkey]);
    // }

    // //TODO
    // function decay() public payable{

    // }

}



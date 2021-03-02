# This is a doc that shows you how to update token limits on mainnet

If you are trying to whitelist and update token limits many addresses at once, you will need to use the bulkSetTokenLockBurnLimit.js script. If you need to only update the lock or burn limit and the token has already been whitelisted, use the setTokenLockBurnLimit script.

Before running the following script go to the data folder and open the `limitWhitelistUpdate.json`. Edit the address value to the address you want to whitelist. Then, update the limit to the desired limit you would like to allow a user to move through the bridge in a single tx. If you want to set this to the max allowable value in solidity, set the limit value to `115792089237316195423570985008687907853269984665640564039457584007913129639935` which is the UINT_MAX in solidity.

Make sure the private key in your .env file is the operator address, ensure the infura id is set correctly as well. Get the bridgebank address and set it in the env var when running the script. To bulk update the whitelist and limits for each token, use bulkSetTokenLockBurnLimit.js like so:
```
BRIDGEBANK_ADDRESS='0x30753E4A8aad7F8597332E813735Def5dD395028' truffle exec scripts/bulkSetTokenLockBurnLimit.js --network develop ../data/limitWhitelistUpdate.json
```

To update the limit on the amount of tokens an already whitelisted smart contract can be locked or burned in the ethereum smart contracts, use this command:
```
UPDATE_ADDRESS="0x0d8cc4b8d15D4c3eF1d70af0071376fb26B5669b" truffle exec scripts/setTokenLockBurnLimit.js --network develop 1000000000000000000000
```

To update the ethereum lock limit in the smart contract, use the previous script with address 0 like so:
```
UPDATE_ADDRESS="0x0000000000000000000000000000000000000000" truffle exec scripts/setTokenLockBurnLimit.js --network develop 1000000000000000000000
```

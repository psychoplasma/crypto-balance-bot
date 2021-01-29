# crypto-balance-bot

`crypto-balance-bot` is a subscription system to inform you about any movements in the accounts subscribed. The publishing is handled by a set of publisher services. Currently only Telegram is implemented. As a blockchain backend, there are serveral 3rd-party service implementations like Blockchain.com, Etherscan.io, Trezor's blockbook. If you want to run a complete independent service as-a-whole, you can install and run your own blockbook as a backend. See [this repo](https://github.com/psychoplasma/blockbook-dockerized) how to install and run blockbook.

I've tried to implement domain-driven-design as much as I can do in this porject. However it's not at its best. This is my first try and I constantly try to imporve. Any comments and help in this aspect would be much appreciated!

## Supported blockchains

* Bitcoin
* Ethereum

and more is coming...

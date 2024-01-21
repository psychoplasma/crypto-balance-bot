# crypto-balance-bot

`crypto-balance-bot` is a subscription system to inform you about any movements in the accounts subscribed. The publishing is handled by a set of publisher services. Currently only Telegram is implemented. As a blockchain backend, there are serveral 3rd-party service implementations like Blockchain.com, Etherscan.io, Trezor's blockbook. If you want to run a complete independent service as-a-whole, you can install and run your own blockbook as a backend. See [this repo](https://github.com/psychoplasma/blockbook-dockerized) how to install and run blockbook.

I've tried to implement domain-driven-design as much as I can do in this porject. However it's not at its best. This is my first try and I constantly try to imporve. Any comments and help in this aspect would be much appreciated!

## Supported blockchains

* Bitcoin
* Ethereum

and more is coming...

## Building and running

```bash
# Copy docker/.example.env file to docker/.env
# And modify the environment file according to your system settings
cp ./docker/.example.env ./docker/.env

# Copy ./example.config.yaml ./config.yaml
# And modify the configuration file.
# You will need a Telegram bot created and a bot token.
# See here https://core.telegram.org/bots/tutorial#obtain-your-bot-token
cp ./example.config.yaml ./config.yaml

# Build the images
docker-compose -f docker/docker-compose.yml build

# Run all the containers of the project
docker-compose -f docker/docker-compose.yml up -d

```

## TODO

* [ ] Make Publisher selectable through configuration file
* [ ] Implement persistance for postgresql

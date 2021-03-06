const fetch = require('node-fetch');

const jsonfile = require('./utils/jsonfile.js');
const blockchain = require('./utils/blockchain-helpers.js');
const configLocation = './config.json';
const config = require(configLocation);

blockchain.deploy(config.peers[0].endpoint,
    "https://github.com/Capgemini-AIE/blockchain-supplychain/chaincode",
    config.peers[0].user,
    ["hi there"])
    
    .then(result => {
        // Read existing config return in an object along with api call result
        return jsonfile.read(configLocation).then(json => ({
            file: json,
            hash: result
        }));
    })
    .then(result => {
        let newConfigData = Object.assign({}, result.file, {chaincodeHash: result.hash});
        return jsonfile.write(configLocation, newConfigData);
    }).catch(error => {
      console.error("[Blockchain-Deploy] Chaincode deployment failed with error:");
      console.error(error);  
    });

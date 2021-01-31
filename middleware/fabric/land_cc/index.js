const CreateLand = require("./createLand");
const GetLand = require("./getLand");
const TransferLand = require("./updateLand");
const GetLands = require("./getLands");
const payload = {
    CreateLand,
    GetLand,
    TransferLand,
    GetLands,
};

module.exports = payload;

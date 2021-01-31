const CreateTransferRequest = require("./createTransferRequest");
const GetTransferRequest = require("./getTransferRequest");
const Transfer2RegistryOfficer = require("./transfer2RegistryOfficer");
const Transfer2BLRO = require("./transfer2BLRO");
const ApproveTransferRequest = require("./approveTransferRequest");
const payload = {
    CreateTransferRequest,
    GetTransferRequest,
    Transfer2RegistryOfficer,
    Transfer2BLRO,
    ApproveTransferRequest,
};

module.exports = payload;

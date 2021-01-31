const express = require("express");
const md5 = require("md5");
const JWTmiddleware = require("../helpers/jwtVerifyMiddleware");
const TransferRequestCC = require("../../fabric/transfer_cc");

const router = new express.Router();

router.get("/api/main/transfer/get/:id", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    const ID = req.params.id;
    try {
        let data = await TransferRequestCC.GetTransferRequest(req.user, ID);
        res.status(200).send(data);
    } catch (error) {
        console.log(error);
        res.status(404).send({ message: "TransferRequest NOT found!" });
    }
});

router.post("/api/main/transfer/createTransferRequest", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    try {
        transferData = req.body.payload;
        await TransferRequestCC.CreateTransferRequest(req.user, transferData);
        res.status(200).send({
            message: "TransferRequest has been successfully added!",
            id: transferData.ID,
        });
    } catch (error) {
        console.log(error);
        res.status(500).send({ message: "Error! TransferRequest NOT Added!" });
    }
});

router.post("/api/main/transfer/transfer2RegistryOfficer", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    try {
        transferData = req.body.payload;
        await TransferRequestCC.Transfer2RegistryOfficer(req.user, transferData);
        res.status(200).send({
            message: "TransferRequest has been successfully Transfered!",
            id: transferData.ID,
        });
    } catch (error) {
        console.log(error);
        res.status(500).send({ message: "Error! TransferRequest NOT Transfered!" });
    }
});

router.post("/api/main/transfer/transfer2BLRO", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    try {
        transferData = req.body.payload;
        await TransferRequestCC.Transfer2BLRO(req.user, transferData);
        res.status(200).send({
            message: "TransferRequest has been successfully Transfered to BLRO!",
            id: transferData.ID,
        });
    } catch (error) {
        console.log(error);
        res.status(500).send({ message: "Error! TransferRequest NOT Transfered to BLRO!" });
    }
});

router.post("/api/main/transfer/approveTransferRequest", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    try {
        transferData = req.body.payload;
        await TransferRequestCC.ApproveTransferRequest(req.user, transferData);
        res.status(200).send({
            message: "TransferRequest has been successfully Approved!",
            id: transferData.ID,
        });
    } catch (error) {
        console.log(error);
        res.status(500).send({ message: "Error! TransferRequest NOT Approved!" });
    }
});

module.exports = router;

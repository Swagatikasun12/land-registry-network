const express = require("express");
const md5 = require("md5");
const JWTmiddleware = require("../helpers/jwtVerifyMiddleware");
const BLROCC = require("../../fabric/blro_cc");

const router = new express.Router();

router.get("/api/main/blro/get/:id", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    const ID = req.params.id;
    try {
        let data = await BLROCC.GetBLRO(req.user, ID);
        res.status(200).send(data);
    } catch (error) {
        console.log(error);
        res.status(404).send({ message: "BLRO NOT found!" });
    }
});

router.post("/api/main/blro/createBLRO", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    try {
        blroData = req.body.payload;
        await BLROCC.CreateBLRO(req.user, blroData);
        res.status(200).send({
            message: "BLRO has been successfully added!",
            id: blroData.ID,
        });
    } catch (error) {
        console.log(error);
        res.status(500).send({ message: "Error! BLRO NOT Added!" });
    }
});

module.exports = router;

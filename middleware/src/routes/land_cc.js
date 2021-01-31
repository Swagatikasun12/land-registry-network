const express = require("express");
const md5 = require("md5");
const JWTmiddleware = require("../helpers/jwtVerifyMiddleware");
const LandCC = require("../../fabric/land_cc");

const router = new express.Router();

router.get("/api/main/land/get/:id", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    const ID = req.params.id;
    try {
        let data = await LandCC.GetLand(req.user, ID);
        res.status(200).send(data);
    } catch (error) {
        console.log(error);
        res.status(404).send({ message: "Land NOT found!" });
    }
});

router.get("/api/main/land/queryOwner/:owner", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    const Owner = req.params.owner;
    try {
        let data = await LandCC.GetLands(req.user, Owner);
        res.status(200).send(data);
    } catch (error) {
        console.log(error);
        res.status(404).send({ message: "Query Error!" });
    }
});

router.post("/api/main/land/createLand", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    try {
        landData = req.body.payload;
        await LandCC.CreateLand(req.user, landData);
        res.status(200).send({
            message: "Land has been successfully added!",
            id: landData.ID,
        });
    } catch (error) {
        console.log(error);
        res.status(500).send({ message: "Error! Land NOT Added!" });
    }
});

router.post("/api/main/land/transferLand", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    try {
        landData = req.body.payload;
        await LandCC.TransferLand(req.user, landData);
        res.status(200).send({ message: `Land has been Successfully Transferd. ID: ${landData.ID}` });
    } catch (error) {
        console.log(error);
        res.status(500).send({ message: `Land NOT transfered! ID: ${landData.ID}` });
    }
});

module.exports = router;

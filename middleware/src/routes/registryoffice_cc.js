const express = require("express");
const md5 = require("md5");
const JWTmiddleware = require("../helpers/jwtVerifyMiddleware");
const RegistryOfficeCC = require("../../fabric/registryoffice_cc");

const router = new express.Router();

router.get("/api/main/registryoffice/get/:id", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    const ID = req.params.id;
    try {
        let data = await RegistryOfficeCC.GetRegistryOfficer(req.user, ID);
        res.status(200).send(data);
    } catch (error) {
        console.log(error);
        res.status(404).send({ message: "RegistryOfficer NOT found!" });
    }
});

router.post("/api/main/registryoffice/createRegistryOfficer", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    try {
        registryofficeData = req.body.payload;
        await RegistryOfficeCC.CreateRegistryOfficer(req.user, registryofficeData);
        res.status(200).send({
            message: "RegistryOfficer has been successfully added!",
            id: registryofficeData.ID,
        });
    } catch (error) {
        console.log(error);
        res.status(500).send({ message: "Error! RegistryOfficer NOT Added!" });
    }
});

module.exports = router;

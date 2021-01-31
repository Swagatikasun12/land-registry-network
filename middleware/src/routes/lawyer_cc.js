const express = require("express");
const md5 = require("md5");
const JWTmiddleware = require("../helpers/jwtVerifyMiddleware");
const LawyerCC = require("../../fabric/lawyer_cc");

const router = new express.Router();

router.get("/api/main/lawyer/get/:id", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    const ID = req.params.id;
    try {
        let data = await LawyerCC.GetLawyer(req.user, ID);
        res.status(200).send(data);
    } catch (error) {
        console.log(error);
        res.status(404).send({ message: "Lawyer NOT found!" });
    }
});

router.post("/api/main/lawyer/createLawyer", JWTmiddleware, async (req, res) => {
    res.setHeader("Access-Control-Allow-Origin", "*");

    try {
        lawyerData = req.body.payload;
        await LawyerCC.CreateLawyer(req.user, lawyerData);
        res.status(200).send({
            message: "Lawyer has been successfully added!",
            id: lawyerData.ID,
        });
    } catch (error) {
        console.log(error);
        res.status(500).send({ message: "Error! Lawyer NOT Added!" });
    }
});

module.exports = router;

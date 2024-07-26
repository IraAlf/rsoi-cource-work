import React from "react";
import { useNavigate, useParams } from "react-router-dom";
import FlightInfoPage from "./FlightInfoPage";

const FlightInfo = () => {
    let navigate = useNavigate();
    return (
        <FlightInfoPage match={useParams()} navigate={navigate}/>
    )
}

export default FlightInfo;

import React from "react";
import { useNavigate } from "react-router-dom";

import LoginPage from "./LoginPage";

const Login = () => {
    let navigate = useNavigate();
    return (
        <LoginPage navigate={navigate}/>
    )
}

export default Login;

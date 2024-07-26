import React from "react";
import { useNavigate } from "react-router-dom";

import SignUpPage from "./SignupPage";

const SignUp = () => {
    let navigate = useNavigate();
    return ( 
        <SignUpPage navigate={navigate}/>
    )
}

export default SignUp;
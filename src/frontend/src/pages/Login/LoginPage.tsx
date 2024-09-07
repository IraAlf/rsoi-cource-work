import React from "react";

import { Box, Link } from "@chakra-ui/react";
import { NavigateFunction } from "react-router-dom";

import Input from "components/Input";
import RoundButton from "components/RoundButton";

import { AuthRequest } from "types/Account"
import { Login as LoginQuery } from "postAPI/accounts/Login";

import styles from "./LoginPage.module.scss";

type LoginProps = {
    navigate: NavigateFunction
}


class LoginPage extends React.Component<LoginProps> {
    acc: AuthRequest = {
        login: "",
        password: "",
    }

    setLogin(val: string) {
        this.acc.login = val
    }
    setPassword(val: string) {
        this.acc.password = val
    }

    submit(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
        let button = e.currentTarget
        button.disabled = true
        LoginQuery(this.acc).then(data => {
            button.disabled = false
            if (data.status === 200) {
                window.location.href = '/flights';
            } else {
                var title = document.getElementById("undertitle")
                if (title)
                    title.innerText = "Ошибка авторизации!"
            }
        });
    }

    render() {
        return <Box className={styles.login_page}>
            <Box className={styles.input_div}>
                <Input name="login" placeholder="Введите логин"
                onInput={event => this.setLogin(event.currentTarget.value)}/>
                <Input name="password" type="password" placeholder="Введите пароль"
                onInput={event => this.setPassword(event.currentTarget.value)}/>
            </Box>

            <Box className={styles.button_div}>
                <RoundButton onClick={ (event) => this.submit(event) }> Войти </RoundButton>
            </Box>
        </Box>
    }
}

export default LoginPage;

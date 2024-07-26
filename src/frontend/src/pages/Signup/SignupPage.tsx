import React from "react";
import theme from "styles/extendTheme";

import { Box } from "@chakra-ui/react";
import { NavigateFunction } from "react-router-dom";

import { NewUserRequest } from "types/Account";
import { Create as CreateQuery } from "postAPI/accounts/Create";

import Input from "components/Input";
import RoundButton from "components/RoundButton";

import styles from "./SignupPage.module.scss";

type SignUpProps = {
    navigate: NavigateFunction
}


class SignUpPage extends React.Component<SignUpProps> {
    acc: NewUserRequest = {
        profile: {
            firstName: "",
            lastName: "",
            email: "",
            login: "",
            mobilePhone: "555-415-1337",
        },
        credentials: {
            password: {
                value: ""
            }
        }
    }
    password: string = ""
    repPassword: string = ""

    setLogin(val: string) {
        this.acc.profile.email = val
        this.acc.profile.login = val
    }
    setFirstName(val: string) {this.acc.profile.firstName = val}
    setLastName(val: string) {this.acc.profile.lastName = val}
    setPassword(val: string) {this.password = val}
    setRepPassword(val: string) {this.repPassword = val}

    highlightNotMatch() {
        let node1 = document.getElementsByName("password")[0]
        let node2 = document.getElementsByName("rep-password")[0]

        if (node1.parentElement && node2.parentElement) {
            node1.parentElement.style.borderColor = theme.colors["title"]
            node2.parentElement.style.borderColor = theme.colors["title"]
        }

        var title = document.getElementById("undertitle")
        if (title)
            title.innerText = "Пароли не совпадают!"
    }

    async submit(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
        if (this.password !== this.repPassword)
            return this.highlightNotMatch()
        else
            this.acc.credentials.password.value = this.password

        // e.currentTarget.disabled = true
        var data = await CreateQuery(this.acc)
        // e.currentTarget.disabled = false
        var title = document.getElementById("undertitle")
        if (title) {
            if (data.status === 200) {
                title.innerText = `Аккаунт ${this.acc.profile.login} успешно создан`
            } else {
                title.innerText = "Ошибка создания аккаунта!"
            }
        }
    }

    render() {
        return <Box className={styles.login_page}>
            <Box className={styles.input_div}>
                <Input name="firstName" placeholder="Введите имя"
                onInput={event => this.setFirstName(event.currentTarget.value)}/>
                <Input name="lastName" placeholder="Введите фамилию" 
                onInput={event => this.setLastName(event.currentTarget.value)}/>
                <Input name="login" placeholder="Введите логин" 
                onInput={event => this.setLogin(event.currentTarget.value)}/>
                <Input name="password" type="password" placeholder="Введите пароль"
                onInput={event => this.setPassword(event.currentTarget.value)}/>
                <Input name="rep-password" type="password" placeholder="Повторите пароль"
                onInput={event => this.setRepPassword(event.currentTarget.value)}/>
            </Box>

            <Box className={styles.input_div}>
                <RoundButton type="submit" onClick={event => this.submit(event)}>
                    Создать аккаунт
                </RoundButton>
            </Box>
        </Box>
    }
}

export default SignUpPage;
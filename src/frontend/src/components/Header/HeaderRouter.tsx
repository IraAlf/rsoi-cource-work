import React from "react";
import { Routes, Route } from "react-router";
import { BrowserRouter } from "react-router-dom";
import Header from ".";
import CategoryHeader from "./Category";
import FlightHeader from "./Recipe";
import SearchHeader from "./Search";
import UserHeader from "./User";


export const HeaderRouter: React.FC<{}> = () => {
    return <BrowserRouter>
        <Routes>
            <Route path="/" element={<Header title="" />} />
            <Route path="/auth/signin" element={<Header title="Вход" />} />
            <Route path="/auth/signup" element={<Header title="Регистрация" undertitle="Зарегистрируйте нового пользователя" />} />
            <Route path="/flights/:fligtNumber" element={<FlightHeader title="" />} />
            <Route path="/tickets" element={<Header title="Личный кабинет" />} />
            <Route path="/statistics" element={<Header title="Статистика сервиса" />} />
            <Route path="/flights" element={<Header title="Доступные рейсы" />} />
            <Route path="*" element={<Header title="Страница не найдена" />} />
        </Routes>
    </BrowserRouter>
}

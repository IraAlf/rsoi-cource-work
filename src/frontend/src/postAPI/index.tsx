import { Category } from "types/Categories";
import { Flight } from "types/Flight";
import { Account } from "types/Account";
import axios from "axios";

export const backUrl = "http://localhost:8080/api/v1";
// export const backUrl = "http://gateway-service:8080/api/v1";

const axiosBackend = () => {
    if (localStorage.getItem('authToken') !== null) {
        console.log('authToken is present in localStorage');
    } else {
        console.log('authToken is not present in localStorage');
    }
    let instance = axios.create({
        baseURL: backUrl
    });

    instance.interceptors.request.use(function (config) {
        const token = localStorage.getItem("authToken");
        if (config.headers && token) {
            config.headers.Authorization = 'Bearer ' + token;
            console.log('Authorization header added:', config.headers.Authorization); // Логирование для проверки
        } else {
            console.log('No token found or headers missing');
        }

        return config;
    });

    return instance;
};
export default axiosBackend();

export type AllFilghtsResp = {
    page: number,
    pageSize: number, 
    totalElements: number,
    items: Flight[]
}

export type AllCategoriesResp = {
    status: number,
    content: Category[]
}

export type AllUsersResp = {
    status: number,
    content: Account[]
}

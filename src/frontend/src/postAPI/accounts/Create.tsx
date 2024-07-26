import axiosBackend from "..";
import { NewUserRequest } from "types/Account";

interface resp {
    status: number
}

export const Create = async function(data: NewUserRequest): Promise<resp> {
    console.log('Create function called');
    const response = await axiosBackend
        .post(`/register`, data)
        .catch((error) => {
            return {
                status: error.response?.status,
            };
        });

    return {
        status: response?.status,
    };
}

import axiosBackend, { AllFilghtsResp } from "..";

const GetFlights = async function(page: number, size: number): Promise<AllFilghtsResp> {
    const response = await axiosBackend
        .get(`/flights?page=${page}&size=${size}`);
    return  response.data;
}

export default GetFlights

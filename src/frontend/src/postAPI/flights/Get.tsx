import axiosBackend from "..";
import { Flight } from "types/Flight";

interface resp {
    status: number
    content: Flight
}

const GetFlight = async function(flightNumber: string): Promise<resp> {
    const response = await axiosBackend
        .get(`/flights/${flightNumber}`);
    return {
        status: response.status,
        content: response.data as Flight
    };
}

export default GetFlight

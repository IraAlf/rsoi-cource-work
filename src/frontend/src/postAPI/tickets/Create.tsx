import axiosBackend from "..";

interface resp {
    status: number
}

type Request = {
    flightNumber: string
    price: number
    paidFromBalance: boolean
}

const CreateTicket = async function(flight: string, price: number, fromBalance: boolean): Promise<resp> {
    let data: Request =  {
        flightNumber: flight,
        price: price,
        paidFromBalance: fromBalance,
    }
    const response = await axiosBackend.post(`/tickets`, data);
    return {
        status: response.status
    };
}
export default CreateTicket;

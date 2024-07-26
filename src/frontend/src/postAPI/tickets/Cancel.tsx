import axiosBackend from "..";

interface resp {
    status: number
}

const CancelTicket = async function(ticketUid: string): Promise<resp> {
    const response = await axiosBackend.delete(`/tickets/${ticketUid}`);
    return {status: response.status};
}
export default CancelTicket;

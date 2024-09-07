import { VStack, Text, Box, HStack } from "@chakra-ui/react"
import React, { useEffect, useState } from "react";
import { Ticket } from "types/Ticket";
import { Privilege } from "types/Privilege";
import ProfileInfo from "postAPI/tickets/Me";
import { subtitleError } from "functions";
import styles from "./AllTickets.module.scss";
import RoundButton from "components/RoundButton";
import CancelTicket from "postAPI/tickets/Cancel";


const PrivilegeCard = ({ privilege }) => {
    const { balance, status } = privilege;

    let cardColor = '';
    switch (status) {
        case 'GOLD':
            cardColor = 'gold';
            break;
        case 'SILVER':
            cardColor = 'silver';
            break;
        case 'BRONZE':
            cardColor = 'bronze';
            break;
        default:
            cardColor = 'default';
            break;
    }

    return (
        <div className={`${styles['privilege-card']} ${styles[cardColor]}`}>
            <h2 className={styles.balance}>{balance}</h2>
            <Text className={styles.notion}>бонусная карта</Text>
            <Text className={`${styles.status} ${styles[cardColor]}`}>{status}</Text>
        </div>
    );
};

const TicketBox = (ticket: Ticket) => {
    const cancel = async (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
        let button = e.currentTarget;
        button.disabled = true;
        let resp = await CancelTicket(ticket.ticketUid);
        if (resp.status === 204) {
            ticket.status = "CANCELED"
        }
        button.disabled = false;
    }

    return (
        <Box className={styles.ticket_box}>
            <HStack>
                <Text>{ticket.flightNumber}</Text>
            </HStack>
            <HStack className={styles.flight_row}>
                <Text>{ticket.fromAirport}</Text>

                <div style={{ margin: "0px" }} >➤</div>
                <div className={styles.arrow}></div>
                <Text>{ticket.date}</Text>
                <div className={styles.arrow}></div>
                <div style={{ margin: "0px" }} >➤</div>

                <Text>{ticket.toAirport}</Text>
            </HStack>

            {ticket.status === "PAID" &&
                <RoundButton className={styles.cancel_button} type="submit" onClick={event => cancel(event)}>
                    Отменить бронирование
                </RoundButton>
            }

            {ticket.status === "CANCELED" &&
                <RoundButton className={styles.canceled_button}>
                    Бронирование отменено
                </RoundButton>
            }
        </Box>
    )
}

const AllTicketsPage = () => {
    const [tickets, setTickets] = useState<Array<Ticket>>([]);
    const [privilege, setPrivilege] = useState<Privilege>();

    useEffect(() => { getProfileInfo() }, []);

    const getProfileInfo = async () => {
        let response = await ProfileInfo();
        if (response.status !== 200) {
            subtitleError("Ошибка получения информации о профиле");
        }
        let profile = response.data;
        setTickets(profile.tickets);
        setPrivilege(profile.privilege);
        console.log(tickets)
    };

    return (
        <VStack style={{width: "100%"}}>
            { privilege && <PrivilegeCard privilege={privilege}/> }
            <Box className={styles.tickets_box}>
                {tickets.map(item => TicketBox(item))}
            </Box>
        </VStack>
    )
}

export default AllTicketsPage;

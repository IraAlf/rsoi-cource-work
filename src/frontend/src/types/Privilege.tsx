import { Ticket } from "./Ticket";

export interface Privilege {
    balance: number;
    status: string;
};

export interface UserInfo {
    tickets: Ticket[];
    privilege: Privilege;
}
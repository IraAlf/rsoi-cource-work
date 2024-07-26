export interface RequestStat {
    path: string;
    responseCode: number;
    method: string;
    startedAt: Date;
    finishedAt: Date;
    duration: number; // в 10^-9 
    userName?: string;
}
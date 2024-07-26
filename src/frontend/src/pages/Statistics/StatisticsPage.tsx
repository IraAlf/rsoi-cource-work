import React, { useCallback, useEffect, useState } from "react";
import { Input } from '@chakra-ui/react';
import ListRequests from "postAPI/statistics";
import RoundButton from "components/RoundButton";
import styles from "./StatisticsPage.module.scss";
import { RequestStat } from "types/Statistics";
import PopularFlights from "./PopularFlights";
import ActiveUsers from "./ActiveUsers";
import { formatDateTime, RequestStatsWidget } from "./formatDateTime";

const StatisticsPage = () => {
    const yesterday = new Date();
    yesterday.setDate(yesterday.getDate() - 1);
    const future = new Date((new Date()).getTime() + 60000);

    const [startDate, setStartDate] = useState(formatDateTime(yesterday));
    const [endDate, setEndDate] = useState(formatDateTime(future));
    const [requests, setRequests] = useState<Array<RequestStat>>([]);

    const handleStartDateChange = (e) => {
        setStartDate(e.target.value);
    };

    const handleEndDateChange = (e) => {
        setEndDate(e.target.value);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        await loadData();
    };
    
    const loadData = useCallback(async () => {
        console.log('Selected time range:', startDate, '-', endDate);
        let resp = await ListRequests(new Date(startDate), new Date(endDate));
        setRequests(resp.requests);
    }, [startDate, endDate]);

    useEffect(() => { loadData() }, [loadData]);

    return (
        <div className={styles.main_div}>
            <form className={styles.date_form} onSubmit={handleSubmit}>
                <div>
                    <Input type="datetime-local" value={startDate} onChange={handleStartDateChange} />
                    -
                    <Input type="datetime-local" value={endDate} onChange={handleEndDateChange} />
                </div>
                <RoundButton type="submit">Запросить статистику</RoundButton>
            </form>

            {(requests.length > 0) &&
            <div className={styles.stats_div}>
                <PopularFlights requestStats={requests}/>
                <ActiveUsers requestStats={requests}/>
                <RequestStatsWidget requestStats={requests} />
            </div>}
        </div>
    );
};


export default StatisticsPage;
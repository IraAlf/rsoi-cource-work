import React, { useState } from "react";
import RoundButton from "components/RoundButton";
import styles from "./StatisticsPage.module.scss";
import { RequestStat } from "types/Statistics";

const PageSize = 20;

// Функция для форматирования даты и времени в строку
export const formatDateTime = (dateTime) => {
    const year = dateTime.getFullYear();
    const month = `${dateTime.getMonth() + 1}`.padStart(2, '0');
    const day = `${dateTime.getDate()}`.padStart(2, '0');
    const hours = `${dateTime.getHours()}`.padStart(2, '0');
    const minutes = `${dateTime.getMinutes()}`.padStart(2, '0');

    return `${year}-${month}-${day}T${hours}:${minutes}`;
};

interface RequestStatsWidgetProps {
    requestStats: RequestStat[];
}

export const RequestStatsWidget: React.FC<RequestStatsWidgetProps> = ({ requestStats }) => {
    const [startIndex, setStartIndex] = useState(0);

    const handleNextPage = () => {
        setStartIndex(startIndex + PageSize);
    };

    const handlePreviousPage = () => {
        setStartIndex(startIndex - PageSize);
    };

    const renderTableRows = () => {
        const slicedRequestStats = requestStats.slice(startIndex, startIndex + PageSize);
        return slicedRequestStats.map((stat, index) => (
            <tr key={index} className={styles.table_row}>
                <td>{stat.path}</td>
                <td>{stat.responseCode}</td>
                <td>{stat.method}</td>
                <td title={stat.startedAt.toDateString()}>
                    {stat.startedAt.toLocaleTimeString()} [{(stat.duration / 1000000000).toFixed(3)} s]
                </td>
                <td>{stat.userName}</td>
            </tr>
        ));
    };

    return (
        <div className={styles.requests_div}>
            <h2>Запросы</h2>

            <table className={styles.table}>
                <tbody>{renderTableRows()}</tbody>
            </table>
            <div className={styles.table_buttons}>
                <RoundButton onClick={handlePreviousPage} disabled={startIndex === 0}>
                    Назад
                </RoundButton>
                <p>{startIndex+1} - {Math.min(startIndex+PageSize, requestStats.length)} (из {requestStats.length})</p>
                <RoundButton onClick={handleNextPage} disabled={startIndex + PageSize >= requestStats.length}>
                    Дальше
                </RoundButton>
            </div>
        </div>
    );
};

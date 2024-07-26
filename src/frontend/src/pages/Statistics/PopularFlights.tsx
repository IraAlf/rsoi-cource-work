import React from "react";
import { RequestStat } from "types/Statistics";

interface Props {
  requestStats: RequestStat[];
}

const PopularFlights: React.FC<Props> = ({ requestStats }) => {
  const countRequestsPerFlight = (requestStats: RequestStat[]): Record<string, number> => {
    const flightCounts: Record<string, number> = {};

    requestStats.forEach((stat) => {
      const match = stat.path.match(/\/api\/v1\/flights\/(\w+)/);
      if (match) {
        const flightNumber = match[1];
        if (flightCounts[flightNumber]) {
          flightCounts[flightNumber]++;
        } else {
          flightCounts[flightNumber] = 1;
        }
      }
    });

    return flightCounts;
  };

  const flightCounts = countRequestsPerFlight(requestStats);

  // Получаем массив рейсов в формате [{ flightNumber: string, count: number }]
  const flights = Object.entries(flightCounts).map(([flightNumber, count]) => ({
    flightNumber,
    count,
  }));

  // Сортируем рейсы по убыванию количества обращений
  const sortedFlights = flights.sort((a, b) => b.count - a.count);

  // Ограничиваем список только тремя самыми популярными рейсами
  const topFlights = sortedFlights.slice(0, 5);

  return (
    <div>
      <h2>Самые популярные рейсы</h2>
      <ul>
        {topFlights.map((flight) => (
          <li key={flight.flightNumber}>
            Рейс {flight.flightNumber}: {flight.count} просмотров
          </li>
        ))}
      </ul>
    </div>
  );
};

export default PopularFlights;

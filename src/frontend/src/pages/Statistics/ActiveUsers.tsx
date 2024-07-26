import React from "react";
import { RequestStat } from "types/Statistics";

interface Props {
  requestStats: RequestStat[];
}

const ActiveUsers: React.FC<Props> = ({ requestStats }) => {
  const countRequestsPerUser = (requestStats: RequestStat[]): Record<string, number> => {
    const userCounts: Record<string, number> = {};

    requestStats.forEach((stat) => {
      const userName = stat.userName || "";
      if (userCounts[userName]) {
        userCounts[userName]++;
      } else {
        userCounts[userName] = 1;
      }
    });

    return userCounts;
  };

  const userCounts = countRequestsPerUser(requestStats);

  // Получаем массив пользователей в формате [{ userName: string, count: number }]
  const users = Object.entries(userCounts).map(([userName, count]) => ({
    userName,
    count,
  }));

  // Сортируем пользователей по убыванию количества обращений
  const sortedUsers = users.sort((a, b) => b.count - a.count);

  // Ограничиваем список только тремя самыми активными пользователями
  const topUsers = sortedUsers.slice(0, 3);

  return (
    <div>
      <h2>Самые активные пользователи</h2>
      <ul>
        {topUsers.map((user) => (
          <li key={user.userName}>
            {user.userName ? (
              <span>
                Пользователь {user.userName}: {user.count} запроов
              </span>
            ) : (
              <span>
                {user.count} неавторизированных запросов
              </span>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default ActiveUsers;

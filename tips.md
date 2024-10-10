# about

## 同日で被っているチケット

userIdごとに、同日で有効な(支払済)被るチケットが何枚あるか(ticket_count)を出してくれるクエリ
```sql
SELECT
    u.userId,
    MAX(t.eventName) AS eventName,
    DATE(t.eventDate) AS eventDate,
    COUNT(*) AS ticket_count
FROM
    user_tickets u
INNER JOIN tickets t ON u.ticketId = t.ticketId
WHERE
    (u.isPaid = 1) OR (u.isPaid = 0 AND u.payLimitDate > CURDATE())
GROUP BY
    u.userId, DATE(t.eventDate)
HAVING
    COUNT(*) > 1;
```

```
+---------------+----------------+------------+--------------+
| userId        | eventName      | eventDate  | ticket_count |
+---------------+----------------+------------+--------------+
| hoge@huga.com | Example Event1 | 2024-10-11 |            2 |
| foo@bar.com   | Example Event1 | 2024-10-11 |            2 |
| foo@bar.com   | Example Event2 | 2024-12-10 |            3 |
+---------------+----------------+------------+--------------+
3 rows in set (0.00 sec)
```

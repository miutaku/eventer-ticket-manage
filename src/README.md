# 概要

以下のように投げてやることで、

```
(*>△<)< curl -X POST -H "Content-Type: application/json" -d '{
          "ticketService": "exampleService",
          "ticketRegistDate": "2024-02-03T13:34:56Z",
          "eventDate": "2024-12-10T11:34:56Z",
          "eventPlace": "Tokyo Dome",
          "eventName": "Example Event2",
          "ticketCount": 2,
          "isReserve": true,
          "payLimitDate": "2011-04-30T23:59:59Z","isPaid": false,"userId": "foo@bar.com"
        }' http://localhost:8080/insert

Data inserted successfully
```

DBにこのように登録される。

```
mysql> select * from tickets;
+----------+----------------+---------------------+----------------+---------------------+------------+
| ticketId | ticketService  | ticketRegistDate    | eventName      | eventDate           | eventPlace |
+----------+----------------+---------------------+----------------+---------------------+------------+
|        1 | exampleService | 2024-02-03 13:34:56 | Example Event2 | 2024-12-10 11:34:56 | Tokyo Dome |
+----------+----------------+---------------------+----------------+---------------------+------------+
1 row in set (0.00 sec)

mysql> select * from user_tickets;
+-------------+----------+-------------+-----------+---------------------+--------+
| userId      | ticketId | ticketCount | isReserve | payLimitDate        | isPaid |
+-------------+----------+-------------+-----------+---------------------+--------+
| foo@bar.com |        1 |           2 |         1 | 2011-04-30 23:59:59 |      1 |
| foo@bar.com |        1 |           2 |         1 | 2011-04-30 23:59:59 |      0 |
+-------------+----------+-------------+-----------+---------------------+--------+
2 rows in set (0.00 sec)```

# tips
チケットはそもそもやらかして2重に登録することもあるし、user_ticketsテーブルは主キーを定義していない。

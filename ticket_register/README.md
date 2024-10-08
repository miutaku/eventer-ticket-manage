# 概要

以下のように投げてやることで、

```
(*>△<)< curl -X POST -H "Content-Type: application/json" -d '{
          "ticketService": "livepocket",
          "ticketRegistDate": "2024-02-03T13:34:56Z",
          "eventDate": "2024-10-11T16:30:00Z",
          "eventPlace": "Tokyo Dome",
          "eventName": "Example  Event",
          "ticketCount": 2,
          "isReserve": true,
          "payLimitDate": "2024-04-30T23:59:59Z",
          "userId": "foo@bar.com"
        }' http://localhost:8080/insert

Data inserted succ
```

DBにこのように登録される。

```
mysql> select * from tickets;
+----------+----------------+---------------------+----------------+---------------------+------------+
| ticketId | ticketService  | ticketRegistDate    | eventName      | eventDate           | eventPlace |
+----------+----------------+---------------------+----------------+---------------------+------------+
|        1 | livepocket     | 2024-02-03 13:34:56 | Example Event  | 2024-10-11 16:30:00 | Tokyo Dome |
|        2 | eplus          | 2023-04-13 23:14:56 | Example Event2 | 2023-09-21 13:00:00 | Tokyo Dome |
+----------+----------------+---------------------+----------------+---------------------+------------+
1 row in set (0.00 sec)

mysql> select * from user_tickets;
+-----------------+----------+-------------+-----------+---------------------+
| userId          | ticketId | ticketCount | isReserve | payLimitDate        |
+-----------------+----------+-------------+-----------+---------------------+
| hoge_2@huga.com |        1 |           2 |         1 | 2011-04-30 23:59:59 |
| foo@bar.com     |        1 |           2 |         1 | 2024-04-30 23:59:59 |
| hoge@huga.com   |        2 |           2 |         1 | 2011-04-30 23:59:59 |
| hoge@huga.com   |        2 |           2 |         1 | 2011-04-30 23:59:59 |
| foo@bar.com     |        2 |           2 |         1 | 2011-04-30 23:59:59 |
+-----------------+----------+-------------+-----------+---------------------+
4 rows in set (0.00 sec)
```

# tips
チケットはそもそもやらかして2重に登録することもあるし、user_ticketsテーブルは主キーを定義していない。

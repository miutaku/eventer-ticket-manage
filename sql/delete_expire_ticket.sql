CREATE EVENT delete_old_tickets
ON SCHEDULE EVERY 1 DAY
DO
  DELETE FROM tickets
  WHERE event_date < CURDATE();

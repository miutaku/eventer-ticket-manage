CREATE DATABASE IF NOT EXISTS ticket_register;
USE ticket_register;

CREATE TABLE tickets (
  ticketId BIGINT AUTO_INCREMENT PRIMARY KEY,
  ticketService VARCHAR(100),
  ticketRegistDate DATETIME,
  eventName VARCHAR(100),
  eventDate DATETIME,
  eventPlace VARCHAR(100),
  UNIQUE INDEX unique_event_index (eventName, eventDate, eventPlace)
);

CREATE TABLE user_tickets (
  userId VARCHAR(100),
  ticketId BIGINT,
  ticketCount INT,
  isReserve BOOLEAN,
  payLimitDate DATETIME,
  isPaid BOOLEAN,
  FOREIGN KEY (ticketId) REFERENCES tickets(ticketId)
);

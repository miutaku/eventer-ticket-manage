CREATE DATABASE IF NOT EXISTS ticket_register;
USE ticket_register;

CREATE TABLE tickets (
  ticketId BIGINT AUTO_INCREMENT PRIMARY KEY,
  ticketService VARCHAR(100),
  ticketRegistDate DATETIME,
  eventName VARCHAR(100),
  eventDate DATETIME,
  eventPlace VARCHAR(100)
);


CREATE TABLE user_tickets (
  userId VARCHAR(100),
  ticketId BIGINT,
  ticketCount INT,
  isReserve BOOLEAN,
  payLimitDate DATETIME,
  PRIMARY KEY (userId, ticketId),
  FOREIGN KEY (ticketId) REFERENCES tickets(ticketId)
);

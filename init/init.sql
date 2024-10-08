CREATE DATABASE IF NOT EXISTS ticket_register;
USE ticket_register;
CREATE TABLE tickets (
  ticketID BIGINT AUTO_INCREMENT PRIMARY KEY,
  ticketService VARCHAR(100),
  eventName VARCHAR(100),
  eventDate DATETIME,
  eventPlace VARCHAR(100)
);

CREATE TABLE user_details (
  userID VARCHAR(100),
  ticketRegistDate DATETIME,
  ticketCount INT,
  isReserve BOOLEAN,
  payLimitDate DATETIME,
  PRIMARY KEY (userID, ticketRegistDate)
);
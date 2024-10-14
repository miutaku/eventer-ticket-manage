CREATE DATABASE IF NOT EXISTS ticket_management;
USE ticket_management;

CREATE TABLE tickets (
  ticketId BIGINT AUTO_INCREMENT PRIMARY KEY,
  ticketService VARCHAR(100),
  ticketRegistDate DATETIME,
  eventName VARCHAR(100),
  eventDate DATETIME,
  eventPlace VARCHAR(100),
  UNIQUE INDEX unique_event_index (eventName, eventDate, eventPlace)
);

CREATE TABLE users { /* auth0 でユーザアカウントは管理する */
  userId VARCHAR(100) PRIMARY KEY,
  userEmail VARCHAR(100),
  auth0Nickname VARCHAR(100), /* auth0 req.nickname */
  registDateTime DATETIME
};

CREATE TABLE user_lawson_ticket_accounts {
  userId VARCHAR(100) /* usersのuserId */,
  lawsonTicketMailAddress VARCHAR(100) PRIMARY KEY,
};

CREATE TABLE user_tickets (
  userId VARCHAR(100),
  ticketId BIGINT,
  ticketCount INT,
  isReserve BOOLEAN,
  payLimitDate DATETIME,
  isPaid BOOLEAN,
  duplicateTicketId VARCHAR(255),
  isDuplicate BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (ticketId) REFERENCES tickets(ticketId)
);

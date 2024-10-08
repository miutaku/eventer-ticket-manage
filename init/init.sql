CREATE DATABASE IF NOT EXISTS ticket_register;
USE ticket_register;
CREATE TABLE ticketList (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ticketService VARCHAR(100) NOT NULL,
    registDate DATETIME NOT NULL,
    eventDate DATETIME NOT NULL,
    eventPlace VARCHAR(255),
    eventName VARCHAR(255),
    ticketCount INT,
    isReserve BOOLEAN,
    payLimitDate DATETIME
);

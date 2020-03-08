CREATE DATABASE IF NOT EXISTS spendbucket;
USE spendbucket;

--  Entities

CREATE TABLE Users (
  userId INT NOT NULL AUTO_INCREMENT,
  username VARCHAR(32) UNIQUE NOT NULL,
  pword CHAR(60) NOT NULL,
  PRIMARY KEY (userId)
);


CREATE TABLE Rooms (
  groupId INT NOT NULL AUTO_INCREMENT,
  groupName VARCHAR(32) NOT NULL,
  inviteCode VARCHAR(32) UNIQUE NOT NULL,
  createdBy INT NOT NULL,
  FOREIGN KEY (createdBy) REFERENCES Users (userId) ON DELETE CASCADE,
  PRIMARY KEY (groupId)
);


CREATE TABLE CanAccess (
  userId INT NOT NULL,
  groupId INT NOT NULL,
  FOREIGN KEY (userId) REFERENCES Users (userId) ON DELETE CASCADE,
  FOREIGN KEY (groupId) REFERENCES Rooms (groupId) ON DELETE CASCADE,
  PRIMARY KEY (userId, groupId)
);


-- Persons are strong entities for convenience.
CREATE TABLE Persons (
  personId INT NOT NULL AUTO_INCREMENT,
  groupId INT NOT NULL,
  personName VARCHAR(32) NOT NULL,
  FOREIGN KEY (groupId) REFERENCES Rooms (groupId) ON DELETE CASCADE,
  PRIMARY KEY (personId)
);


CREATE TABLE Transactions (
  transId INT NOT NULL AUTO_INCREMENT,
  title VARCHAR(32) NOT NULL,
  amount DECIMAL(13,2) NOT NULL,
  transDate DATE NOT NULL,
  groupId INT NOT NULL,
  paidBy INT NOT NULL,
  FOREIGN KEY (groupId) REFERENCES Rooms (groupId) ON DELETE CASCADE,
  FOREIGN KEY (paidBy) REFERENCES Persons (personId) ON DELETE CASCADE,
  PRIMARY KEY (transId)
);

CREATE TABLE SharedBy (
  transId INT NOT NULL,
  personId INT NOT NULL,
  FOREIGN KEY (transId) REFERENCES Transactions (transId) ON DELETE CASCADE,
  FOREIGN KEY (personId) REFERENCES Persons (personId) ON DELETE CASCADE,
  PRIMARY KEY(transId, personId)
);

-- Views

CREATE VIEW Credit AS
  SELECT paidBy AS personId, SUM(amount) AS amount
  FROM Transactions
  GROUP BY personId;

CREATE VIEW SharedCount AS
  SELECT transId, COUNT(*) AS num
  FROM SharedBy
  GROUP BY transId;

CREATE VIEW SharedAmount AS
  SELECT t.transId, amount/num AS amount
  FROM Transactions t
  JOIN SharedCount s
  ON t.transId = s.transId;

CREATE VIEW Debt AS
  SELECT personId, SUM(amount) AS amount
  FROM SharedBy b
  JOIN SharedAmount a
  ON b.transId = a.transId
  GROUP BY personId;

-- union required for full outer join (sans duplicates)
CREATE VIEW Balance AS
  SELECT c.personId, IFNULL(c.amount, 0) - IFNULL(d.amount, 0) as amount
  FROM Credit c
  LEFT OUTER JOIN Debt d
  ON c.personId = d.personId
  UNION
  SELECT d.personId, IFNULL(c.amount, 0) - IFNULL(d.amount, 0) as amount
  FROM Credit c
  RIGHT OUTER JOIN Debt d
  ON c.personId = d.personId;
Id;
redit c
--   LEFT OUTER JOIN Debt d
--   ON c.personId = d.personId
--   UNION
--   SELECT d.personId, IFNULL(c.amount, 0) - IFNULL(d.amount, 0) as amount
--   FROM Credit c
--   RIGHT OUTER JOIN Debt d
--   ON c.personId = d.personId;
--
d;
--

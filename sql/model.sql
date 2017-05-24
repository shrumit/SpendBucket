--  Entities

CREATE TABLE Users (
  userId INT NOT NULL AUTO_INCREMENT,
  username VARCHAR(32) UNIQUE NOT NULL,
  pword VARCHAR(32) NOT NULL,
  PRIMARY KEY (userId)
);

INSERT INTO Users (username, pword) VALUES ('bob', 'pass');

CREATE TABLE Groups (
  groupId INT NOT NULL AUTO_INCREMENT,
  groupName VARCHAR(32) NOT NULL,
  inviteCode VARCHAR(32) UNIQUE NOT NULL,
  PRIMARY KEY (groupId)
);

INSERT INTO Groups (inviteCode, groupName) VALUES ('invite', 'g1');

CREATE TABLE CanAccess (
  userId INT NOT NULL,
  groupId INT NOT NULL,
  FOREIGN KEY (userId) REFERENCES Users (userId),
  FOREIGN KEY (groupId) REFERENCES Groups (groupId),
  PRIMARY KEY (userId, groupId)
);

INSERT INTO CanAccess VALUES (1,1);

-- Persons are strong entities for convenience.
CREATE TABLE Persons (
  personId INT NOT NULL AUTO_INCREMENT,
  groupId INT NOT NULL,
  personName VARCHAR(32) NOT NULL,
  FOREIGN KEY (groupId) REFERENCES Groups (groupId),
  PRIMARY KEY (personId)
);

INSERT INTO Persons VALUES (1,1,'p1');
INSERT INTO Persons VALUES (2,1,'p2');
INSERT INTO Persons VALUES (3,1,'p3');

CREATE TABLE Transactions (
  transId INT NOT NULL AUTO_INCREMENT,
  title VARCHAR(32) NOT NULL,
  amount DECIMAL(13,2) NOT NULL,
  transDate DATE NOT NULL,
  groupId INT NOT NULL,
  paidBy INT NOT NULL,
  FOREIGN KEY (groupId) REFERENCES Groups (groupId),
  FOREIGN KEY (paidBy) REFERENCES Persons (personId),
  PRIMARY KEY (transId)
);

INSERT INTO Transactions VALUES (1,'trans1', 30.00, '2017-01-01', 1, 1);

CREATE TABLE SharedBy (
  transId INT NOT NULL,
  personId INT NOT NULL,
  FOREIGN KEY (transId) REFERENCES Transactions (transId) ON DELETE CASCADE,
  FOREIGN KEY (personId) REFERENCES Persons (personId),
  PRIMARY KEY(transId, personId)
);

INSERT INTO SharedBy VALUES (1, 1);
INSERT INTO SharedBy VALUES (1, 2);

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

-- UNION required for full outer join (sans duplicates)
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



-- CREATE VIEW GetTransactions
--   SELECT transId, date, name, paidBy, t.amount, groupId, GROUP_CONCAT(CONCAT(userId, ":", s.amount) ORDER BY userId ASC SEPARATOR ',') AS share
--   FROM Transactions t
--   JOIN SharedBy s USING (transId)
--   WHERE groupId=1
--   GROUP BY transId,
--   ORDER BY date;

-- CREATE TABLE test (
--   ID INTEGER,
--   NAME VARCHAR (50),
--   VALUE INTEGER
-- );

-- INSERT INTO test VALUES (1, 'A', 4);
-- INSERT INTO test VALUES (1, 'A', 5);
-- INSERT INTO test VALUES (1, 'B', 8);
-- INSERT INTO test VALUES (2, 'C', 9);

-- SELECT ID, GROUP_CONCAT(NAME ORDER BY NAME ASC SEPARATOR ',')
-- FROM (
--   SELECT ID, CONCAT(NAME, ':', GROUP_CONCAT(VALUE ORDER BY VALUE ASC SEPARATOR ',')) AS NAME
--   FROM test
--   GROUP BY ID, NAME
-- ) AS A
-- GROUP BY ID;

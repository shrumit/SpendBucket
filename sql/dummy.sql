-- Initialize

INSERT INTO Users 
    VALUES (1, 'dummy@example.com', '$2a$10$DX.CmOpStKyqxnBXPw1yvOZHS2wWz.95aLV63ZPZr9yu6atFabCZa');

INSERT INTO Groups VALUES (1, 'Smith residence', 'PLkgLbrYqE', 1);
INSERT INTO Groups VALUES (2, 'Rick and Birdperson', 'lVishTJCBJ', 1);

-- Event

    DELIMITER $$
    CREATE EVENT dummy
    ON SCHEDULE EVERY 10 MINUTE STARTS '2017-01-01 00:00:00'
    DO BEGIN
        -- cleanup
        DELETE FROM Groups WHERE createdBy=1 AND groupId!=1 AND groupId!=2;
        DELETE FROM CanAccess WHERE userId=1 OR groupId=1 OR groupId=2;
        DELETE FROM Transactions WHERE groupId=1 OR groupId=2;
        DELETE FROM Persons WHERE groupId=1 OR groupId=2;

        --  group 1
        INSERT INTO CanAccess VALUES (1, 1);
        INSERT INTO Persons VALUES (1,1,'Rick');
        INSERT INTO Persons VALUES (2,1,'Morty');
        INSERT INTO Persons VALUES (3,1,'Summer');

        INSERT INTO Transactions VALUES (1, 'Fancy Eats', 14.50, '2017-05-15', 1, 3);
        INSERT INTO SharedBy VALUES (1,1);
        INSERT INTO SharedBy VALUES (1,2);
        INSERT INTO SharedBy VALUES (1,3);

        INSERT INTO Transactions VALUES (2, 'Jerryboree', 32, '2017-05-10', 1, 2);
        INSERT INTO SharedBy VALUES (2,1);
        INSERT INTO SharedBy VALUES (2,2);

        INSERT INTO Transactions VALUES (3, 'Galactic Sauce Vault', 5.49, '2017-05-16', 1, 1);
        INSERT INTO SharedBy VALUES (3,1);

        INSERT INTO Transactions VALUES (4, 'Blips and Chitz', 27.12, '2017-05-02', 1, 1);
        INSERT INTO SharedBy VALUES (4,1);
        INSERT INTO SharedBy VALUES (4,2);

        INSERT INTO Transactions VALUES (5, 'Cogspot', 48, '2017-05-25', 1, 2);
        INSERT INTO SharedBy VALUES (5,1);

        INSERT INTO Transactions VALUES (6, 'Egan\'s Cinema', 16.90, '2017-06-01', 1, 1);
        INSERT INTO SharedBy VALUES (6,1);
        INSERT INTO SharedBy VALUES (6,2);
        INSERT INTO SharedBy VALUES (6,3);

        -- group 2
        INSERT INTO CanAccess VALUES (1,2);
        INSERT INTO Persons VALUES (4,2,'Rick');
        INSERT INTO Persons VALUES (5,2,'Birdperson');

        INSERT INTO Transactions VALUES (7, 'Plim Plom Tavern', 88.29, '2017-05-15', 2, 4);
        INSERT INTO SharedBy VALUES (7,4);
        INSERT INTO SharedBy VALUES (7,5);

        INSERT INTO Transactions VALUES (8, 'Snorlab', 150, '2017-05-14', 2, 4);
        INSERT INTO SharedBy VALUES (8,4);
        INSERT INTO SharedBy VALUES (8,5);
    END $$
    DELIMITER ;

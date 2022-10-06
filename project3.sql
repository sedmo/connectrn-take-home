CREATE TABLE Users{
    UserID INT PRIMARY KEY AUTO_INCREMENT,
    Firstname VARCHAR(30),
    Lastname VARCHAR(30) NOT NULL,
    City VARCHAR(30),
    Zipcode VARCHAR(10),
    PasswordHistoryID INT NOT NULL
};
CREATE TABLE PasswordHistory{
    PasswordHistoryID INT PRIMARY KEY AUTO_INCREMENT,
    Password VARCHAR(32) NOT NULL,
    Changedate VARCHAR(10) NOT NULL,
    Currentlyactive BOOL NOT NULL DEFAULT true,
    UserID INT NOT NULL,
    FOREIGN KEY (UserID)
        REFERENCES Users(UserID)
        ON DELETE CASCADE
};

SELECT  u.UserId,
        p.Password 
FROM PasswordHistory p 
INNER JOIN Users u 
USING(UserId);

START TRANSACTION;

INSERT INTO PasswordHistory(Password, Changedate, Currentlyactive, UserID)
VALUES('password','2022-10-06',true,1);

SELECT @latestUserID:=MAX(UserID)+1
FROM Users;

UPDATE PasswordHistory
SET Currentlyactive=false
WHERE UserID=@latestUserID;


COMMIT;

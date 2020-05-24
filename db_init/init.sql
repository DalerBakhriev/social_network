CREATE TABLE IF NOT EXISTS users (
      id INT NOT NULL AUTO_INCREMENT,
      email VARCHAR(100) NOT NULL UNIQUE,
      name VARCHAR(100),
      surname VARCHAR(100),
      age INT NOT NULL,
      sex VARCHAR(100),
      interests MEDIUMTEXT,
      city VARCHAR(100),
      encrypted_password VARCHAR(100) NOT NULL,
      PRIMARY KEY (id, email)
      );

CREATE TABLE IF NOT EXISTS friends (
      user_id INT,
      friend_id INT,
      is_accepted BOOLEAN NOT NULL DEFAULT FALSE,
      PRIMARY KEY (user_id, friend_id),
      FOREIGN KEY (user_id)
          REFERENCES users (id)
          ON UPDATE RESTRICT ON DELETE CASCADE,
      FOREIGN KEY (friend_id) 
          REFERENCES users (id)
          ON UPDATE RESTRICT ON DELETE CASCADE
      );

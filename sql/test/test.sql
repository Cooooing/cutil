CREATE TABLE users
(
    id         INT PRIMARY KEY AUTO_INCREMENT,
    name       VARCHAR(50),
    age        INT,
    email      VARCHAR(100),
    created_at DATETIME
);

INSERT INTO users (name, age, email, created_at)
VALUES ('Alice', 25, 'alice@example.com', NOW()),
       ('Bob', 30, 'bob@example.com', NOW()),
       ('Charlie', 22, 'charlie@example.com', NOW());

CREATE TABLE posts
(
    id         INT PRIMARY KEY AUTO_INCREMENT,
    user_id    INT,
    title      VARCHAR(100),
    content    TEXT,
    created_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO posts (user_id, title, content, created_at)
VALUES (1, 'Hello World', 'This is Alice''s first post', NOW()),
       (2, 'My Post', 'This is Bob''s post', NOW()),
       (1, 'Another Post', 'Alice writes again', NOW());

CREATE TABLE comments
(
    id         INT PRIMARY KEY AUTO_INCREMENT,
    post_id    INT,
    user_id    INT,
    content    TEXT,
    created_at DATETIME,
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO comments (post_id, user_id, content, created_at)
VALUES (1, 2, 'Nice post, Alice!', NOW()),
       (1, 3, 'Thanks for sharing!', NOW()),
       (2, 1, 'Good work, Bob!', NOW());

CREATE TABLE tags
(
    id   INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) UNIQUE
);

INSERT INTO tags (name) VALUES ('go'), ('sql'), ('programming');

CREATE TABLE post_tags
(
    post_id INT,
    tag_id  INT,
    PRIMARY KEY (post_id, tag_id),
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (tag_id) REFERENCES tags (id)
);

INSERT INTO post_tags (post_id, tag_id)
VALUES (1, 1),
       (1, 2),
       (2, 3),
       (3, 1);

CREATE TABLE users (
    chat_id INTEGER, -- telegram chatID.
    user_id INTEGER, -- userID инициатора взаимодействия.
    username VARCHAR(255), -- username инициатора взаимодействия.
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- created_at
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP -- updated_at
);

CREATE TABLE notifications (
   id INTEGER  primary key AUTOINCREMENT,
   chat_id INTEGER not null, -- telegram chatID.
   user_id INTEGER not null, -- userID инициатора взаимодействия.
   tag VARCHAR(255) not null, -- username инициатора взаимодействия.
   description text not null, -- описание события.
   notify_at DATETIME not null,
   event_at DATETIME not null,
   created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
   updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notifications_build (
   chat_id INTEGER primary key, -- telegram chatID.
   user_id INTEGER primary key, -- userID инициатора взаимодействия.
   tag VARCHAR(255), -- username инициатора взаимодействия.
   description text, -- описание события.
   notify_at DATETIME,
   event_at DATETIME,
   created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
   updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
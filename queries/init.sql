CREATE TABLE IF NOT EXISTS folders (
                                       id INTEGER PRIMARY KEY AUTOINCREMENT,
                                       path TEXT NOT NULL UNIQUE,
                                       total_size INTEGER DEFAULT 0,
                                       file_count INTEGER DEFAULT 0,
                                       created_time DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS files (
                                     id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     folder_id INTEGER NOT NULL,
                                     name TEXT NOT NULL,
                                     size INTEGER NOT NULL,
                                     created_time DATETIME NOT NULL,
                                     FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE
);


INSERT INTO area (name, description) VALUES 
    ('msk.ru', 'Moscow'), 
    ('spb.ru', 'SPB');

INSERT INTO file (name, sha, size, description) VALUES 
    ('file1.mp4', '123456789', 1024, 'desc file1.mp4'), 
    ('file2.mp4', '678912345', 4096, 'desc file2.mp4');

INSERT INTO server (area_id, name, hostname, description) VALUES 
    (1, 'srv1', 'srv1.msk.ru', 'desc srv1.msk.ru'), 
    (1, 'srv2', 'srv2.msk.ru', 'desc srv2.msk.ru'), 
    (2, 'srv3', 'srv3.spb.ru', 'desc srv3.spb.ru'), 
    (2, 'srv4', 'srv4.spb.ru', 'desc srv4.spb.ru');

INSERT INTO "user" (name, balance) VALUES 
    ('user1', 100), 
    ('user2', 200);

INSERT INTO userfile (user_id, file_id) VALUES 
    (1, 1), 
    (1, 2),
    (2, 2);

INSERT INTO serverfile (server_id, file_id) VALUES 
    (1, 1), 
    (1, 2),
    (2, 1),
    (2, 2);

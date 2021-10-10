


-- список файлов на сервере 'srv1'
SELECT file.name FROM serverfile AS sf left join file on sf.file_id = file.id where sf.removed is null and sf.server_id in (select id from "server" s  where s.name = 'srv1');

-- список файлов в зоне 'msk.ru'
select distinct file.name FROM serverfile AS sf left join file on sf.file_id = file.id where sf.removed is null and sf.server_id in (select id from "server" s  where s.area_id in (select id from area where name = 'msk.ru'));

-- список файлов пользователя 'user1'
select file.name from userfile as uf left join file on uf.file_id = file.id where uf.removed is null and uf.user_id in (select id from "user" where name = 'user1');


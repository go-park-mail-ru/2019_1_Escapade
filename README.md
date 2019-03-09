# 2019_1_Escapade
:heart_eyes: Top backend :heart_eyes:


Локальный запуск 

1. Для подключения к БД необходимо создать переменную окружения

> cd
> nano .bashrc 
в конец дописываем добавим такие строчки:

### Support heroku env DATABASE_URL
export DATABASE_URL="dbname=my_database user=postgres password=my_password sslmode=disable" 

# Сервис поздравления сотрудников с днем рождения

## Задание
Написать сервис для поздравлений с днем рождения:  
• Цель удобное поздравление сотрудников; <br>
• Получение списка сотрудников любым способом(api/ad ldap/прямая регистрация); <br>
• Авторизация; <br>
• Возможность подписаться/отписаться от оповещения о дне рождения; <br>
• Оповещение о ДР того, на кого подписан; <br>
• Внешнее взаимодействие (json арi или фронт или тг бот); <br>
• В случае взаимодействия через телеграм бот (создание группы и добавление в нее всех подписанных); <br>
• В случае взаимодействие через фронт настройка времени оповещения до дня рождения на почту; <br>
• Желательно предоставить юнит тесты; <br>
• Обязательно нужно добавлять в репозиторий файл README с информацией о том, как запускать приложение, и как запускать тесты. А так же описать систему каталогов, что и где у вас расположено. <br>

## Как запустить проект
```
// используем ssh для подключения к удаленному репозиторию
git clone git@github.com:cantylv/rutube-happy-birthday.git
// используем команды, определенные в Makefile
make vendor
make run
```

## API
### Авторизация

`POST /api/v1/signup` - Регистрация пользователя в системе. <br>
<b>Тело запроса (raw json):</b>
```
{
	"full_name": "Medvedev Igor Petrovich",
    "email": "mede2020@mail.ru",
    "password": "qwerty12345",
    "birthday": "23.08.2003"
}
```
<b>Статусы ответов (code statuses): </b>

200 - `{"detail": "you've succesful signed up"}` <br>
400 - `{"error": "incorrect data received, please try again"}` <br>
`{"error: "user already exists"}` <br>
401 - `{"error": "you're already registered"}` <br>
500 - `{"error":"unexpected internal server error, please try again in one minute"}` <br>


<br>`POST /api/v1/signin` - Авторизация пользователя в системе.<br>
<b>Тело запроса (raw json):</b>
```
{
    "email": "mede2020@mail.ru",
    "password": "qwerty12345"
}
```
<b>Статусы ответов (code statuses): </b> <br>

200 - `{"detail": "you've succesful signed in"}` <br>
400 - `{"error": "incorrect data received, please try again"}` <br>
    - `{"error": "incorrect password or login"}` <br>
401 - `{"error": "you're already authenticated"}` <br>
500 - `{"error":"unexpected internal server error, please try again in one minute"}` <br>   

<br> `POST /api/v1/signout` - Деавторизация пользователя в системе. <br>
<b>Статусы ответов (code statuses): </b> <br>

200 - `{"detail": "you've succesful signed out"}` <br>
401 - `{"error": "you're not authenticated"}` <br>
500 - `{"error":"unexpected internal server error, please try again in one minute"}` <br>


#### Важное примечание: после авторизации/регистрации незабудьте использовать значение заголовка X-CSRF-Token для последующих запросах авторизованного пользователя (необходимо создать заголовок `X-CSRF-Token`: value). Таким образом обеспечивается защита от CSRF-аттак.

### Пользователь

<br> `GET /api/v1/user` - Получение данных о пользователе. <br>
<b>Статусы ответов (code statuses): </b> <br>

200 - 
```
    {
        "id": "66c606d34b546306dcf504c2",
        "full_name": "Medvedev Igor Petrovich",
        "birthday": "23.08.2003",
        "email": "mede2020@mail.ru",
        "subs": []
    }
```
401 - `{"error": "you're already authenticated"}` <br>
`{"error": "you're not registered in our system"}` <br>
500 - `{"error":"unexpected internal server error, please try again in one minute"}` <br>   


<br> `PUT /api/v1/user` - Изменение данных пользователя. <br>
<b>Статусы ответов (code statuses): </b> <br>

200 - 
```
    {
        "full_name": "Medvedev Igor Petrovich",
        "birthday": "23.08.2003",
        "password": "qwerty12345",
        "email": "mede2020@mail.ru",
    }
```
400 - `{"error": "incorrect data received, please try again"}` <br>
`{"error": "user with this email already exist, failed"}` <br>
401 - `{"error": "you're already authenticated"}` <br>
`{"error": "you're not registered in our system"}` <br>
500 - `{"error":"unexpected internal server error, please try again in one minute"}` <br>   


### Подписки

<br>`POST /api/v1/sub/{employee_id}` - подписка на пользователя с id == {employee_id}
<b>Статусы ответов (code statuses): </b> <br>

200 - `{"detail": "you've succesful subed on employee"}`
400 - `{"error": "you can't subscribe to a non-existent user"}` <br>
    - `{"error": "you can't subscribe to yourself"}` <br>
401 - `{"error": "you're not authenticated"}` <br>
    - `{"error": "you're not registered in our system"}` <br>
500 - `{"error":"unexpected internal server error, please try again in one minute"}` <br>  

<br>`PUT /api/v1/sub/{employee_id}/new_interval/{interval}` - изменение кол-ва дней до оповещения о дне рождении пользователя
<b>Статусы ответов (code statuses): </b> <br>

200 - `{"detail": "you've succesful set up new interval"}`
400 - `{"error": "provided wrong value of path parameter 'employee_id'"}` <br>
    - `{"error": "you don't have subscription"}` <br>
    - `{"error": "you can't set the interval for non-existent user birthday"}` <br>
    - `{"error": "you can't set an interval for your birthday"}` <br>
    - `{"error": "yyou can't set interval birthday if you are't subscribed on employee"}` <br>
401 - `{"error": "you're not authenticated"}` <br>
    - `{"error": "you're not registered in our system"}` <br>
500 - `{"error":"unexpected internal server error, please try again in one minute"}` <br>  


<br>`POST /api/v1/unsub/{employee_id}` - отписка от пользователя с id == {employee_id}
<b>Статусы ответов (code statuses): </b> <br>

200 - `{"detail": "you've succesful unsubed on employee"}`
400 - `{"error": "you can't unsubscribe to a non-existent user"}` <br>
    - `{"error": "you can't unsubscribe to yourself"}` <br>
401 - `{"error": "you're not authenticated"}` <br>
    - `{"error": "you're not registered in our system"}` <br>
500 - `{"error":"unexpected internal server error, please try again in one minute"}` <br>  


### Особенности
1) В качестве базы данных была выбрана MongoDB - документоориентированная база данных. <br>
2) Отписка от пользователя происходит следующим образом: при корректном вызове обработчика идет проверка на то, есть ли в подписках нужный сотрудник, если есть, то у нее меняется поле `is_followed`, которое означает, подписан ли текущий пользователь на сотрудника. 
Если подписки не было до этого момента, то ничего не происходит. Для удаления устаревших данных есть cron-задача, которая раз в установленное время проходит по записям и удаляет все устаревшие подписки (они же отписки). <br>
3) Для оповещения пользователя о дне рождении сотрудника реализовано следующее:<br>
    - есть ручка, которая меняет интервал оповещение (н-р, можно поставить оповещение о дне рождении сотрудника за любой интервал времени, например за 5 дней).
    - есть cron-задача, которая проходится по подпискам пользователей и отправляет в брокер kafka оповещение, если даты сходятся (т.е. текущая дата + интервал = дата дня рождения сотрудника в этом году). После отправки продьюсером консьюмер получает оповещение и записывает его в `notification.json`. <br>
4) Для гибкой настройки окружения использовался viper. <br>


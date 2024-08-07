Распределенный калькулятор 

- На порту 8080 запускается сервер, который регистрирует пользователей, принимает и отправляет выражения оркестратору, сохраняет их в базе данных.
- Оркестратор разбивает их на простые задачки и помещает в очередь.
- Агенты обращаются к оркестратору с помощью GRPC, подбирают задачки из очереди, вычисляют и отправляют решения обратно.
- Оркестратор собирает новое выражение из полученных результатов.
- ...Повторение до конечного результата.

Для запуска:
  - Скачать и запустить файл main.exe !от имени администратора! (Код писался на windows 10 с версией go 1.22.0)

Для взаимодействия:
  - Открыть командную строку и отправить запросы по 8080 порту

Возможные запросы:
  - curl --header "Content-Type:application/json" --data "{\\"login\\": \\"<логин>\\", \\"password\\": \\"<пароль>\\"}" http://localhost:8080/api/v1/register
     - Возможные ответы:
       - try another login // логин занят
       - user registered // пользователь зарегистрирован
  - curl --header "Content-Type:application/json" --data "{\\"login\\": \\"<логин>\\", \\"password\\": \\"<пароль>\\"}" http://localhost:8080/api/v1/login
     - Возможные ответы:
       - wrong password // неверный пароль
       - ошибки бд
       - <токен>
  - curl --header "Content-Type:application/json" --data "{\\"token\\": \\"<токен>\\"}" http://localhost:8080/api/v1/expressions //для получения списка выражений
     - Возможные ответы:
       - [{"Id":0,"Status":"resolved","Result":0},{"Id":1,"Status":"resolved","Result":-1}] (200) // список выражений
       - ошибки токена
       - ошибки бд
   
  - curl --header "Content-Type:application/json" --data "{\\"token\\": \\"<токен>\\"}" http://localhost:8080/api/v1/expressions/:id //для получения определенного выражения
     - Возможные ответы:  
       - {"Id":0,"Status":"accepted","Result":0} (200) // id выражения, его статус(accepted(принят на вычисление) или resolved(решено)) и результат
       - invalid id (500) // некорректный id
       - no rights (500) // id выражения не относится к данному пользователю
       - ошибки токена
       - ошибки бд
   
  - curl --header "Content-Type:application/json" --data "{\\"expression\\": \\"<выражение>\\", \\"token\\": \\"<токен>\\"}" http://localhost:8080/api/v1/calculate //для отправки выражения на вычисление
     - Возможные ответы:  
       - accepted, id = 0 (201) // выражение принято и его id
       - smth goes wrong (500) // ошибка на сервере
       - invalid data <...> (422) // некорректное выражение и ошибка в выражении
       - ошибки токена
       - ошибки бд

Пример взаимодействия с сервером:
  - curl --header "Content-Type:application/json" --data "{\\"login\\": \\"test\\", \\"password\\": \\"test\\"}" http://localhost:8080/api/v1/register
    - user registered
  - curl --header "Content-Type:application/json" --data "{\\"login\\": \\"test\\", \\"password\\": \\"test\\"}" http://localhost:8080/api/v1/login
    - <ТОКЕН> // вставлять в body всех последующих запросов как в примере
  - curl --header "Content-Type:application/json" --data "{\\"expression\\": \\"10 + -99\\", \\"token\\": \\"<ТОКЕН>\\"}" http://localhost:8080/api/v1/calculate
    - accepted, id = 1
  - curl --header "Content-Type:application/json" --data "{\\"expression\\": \\"(2.5 * 4) + (3.5 * -5)\\", \\"token\\": \\"<ТОКЕН>\\"}" http://localhost:8080/api/v1/calculate
    - accepted, id = 2
  - curl --header "Content-Type:application/json" --data "{\\"token\\": \\"<ТОКЕН>\\"}" http://localhost:8080/api/v1/expressions
    - {"expressions":[{"Id":1,"Status":"resolved","Result":-89},{"Id":2,"Status":"resolved","Result":-7.5}]}
  - curl --header "Content-Type:application/json" --data "{\\"token\\": \\"<ТОКЕН>\\"}" http://localhost:8080/api/v1/expressions/2
    - {"expression":{"Id":2,"Status":"resolved","Result":-7.5}}

Примечание:
  - Выражения могут быть как с целыми числами, так и с дробными
  - Скобки ставятся произвольно
  - !!!ВАЖНО!!! У отрицательного числа минус должен стоять вплотную к числу(без пробела), если будет пробел то минус засчитается не к числу и выражение примется но не будет посчитано(остальные знаки могут стоять как вплотную, так и раздельно)

Переменные среды:
  - Для каждого вида вычислений можно установить время минимального выполнения в миллисекундах через переменные среды(по умолчанию значения равны 100) 
     - TIME_ADDITION_MS // суммирование
     - TIME_SUBTRACTION_MS // вычитание
     - TIME_MULTIPLICATIONS_MS // умножение
     - TIME_DIVISIONS_MS // деление

  - Для количества возможных агентов
     - COMPUTING_POWER // по умолчанию равен 10

Распределенный калькулятор 

- На порту 8080 запускается сервер, который принимает и хранит выражения.
- Оркестратор разбивает их на простые задачки и помещает в очередь.
- Агенты подбирают задачки из очереди, вычисляют и отправляют решения обратно.
- Оркестратор собирает новое выражение из полученных результатов.
- ...Повторение до конечного результата.

Для запуска:
  (Код писался на windows 10 с версией go 1.22.0)
  1.Скачать и запустить файл build.exe

Для взаимодействия:
  1. Открыть командную строку и отправить запросы по 8080 порту

Возможные запросы:
  1. curl http://localhost:8080/api/v1/expressions //для получения списка выражений
     Возможные ответы:  {"expressions":[{"Id":0,"Status":"resolved","Result":0},{"Id":1,"Status":"resolved","Result":-1}]} (200) // список выражений
   
  2. curl http://localhost:8080/api/v1/expressions/:id //для получения определенного выражения
     Возможные ответы:  
       1. {"expression":{"Id":0,"Status":"resolved","Result":0}} (200) // id выражения, его статус(accepted(принят на вычисление) или     resolved(решено)) и результат
       2. bad id (404) // выражения по такому id не существует
       3. invalid id (500) // некорректный id
   
  3. curl --header "Content-Type:application/json" --data "{\"expression\": \"<выражение>\"}" http://localhost:8080/api/v1/calculate //для отправки выражения на вычисление
    Возможные ответы:  
      1. accepted, id = 0 (201) // выражение принято и его id
      2. smth goes wrong (500) // ошибка на сервере
      3. invalid data <...> (422) // некорректное выражение и ошибка в выражении

Пример взаимодействия с сервером:
  - curl --header "Content-Type:application/json" --data "{\"expression\": \"10 + -99\"}" http://localhost:8080/api/v1/calculate
    accepted, id = 0
  - curl --header "Content-Type:application/json" --data "{\"expression\": \"(2.5 * 4) + (3.5 * -5)\"}" http://localhost:8080/api/v1/calculate
    accepted, id = 1
  - curl http://localhost:8080/api/v1/expressions
    {"expressions":[{"Id":0,"Status":"resolved","Result":-89},{"Id":1,"Status":"resolved","Result":-7.5}]}
  - curl http://localhost:8080/api/v1/expressions/1
    {"expression":{"Id":1,"Status":"resolved","Result":-7.5}}

Примечание:
  - Выражения могут быть как с целыми числами, так и с дробными
  - Скобки ставятся произвольно
  - !!!ВАЖНО!!! У отрицательного числа минус должен стоять вплотную к числу(без пробела), если будет пробел то минус засчитается не к числу и выражение примется но не будет посчитано(остальные знаки могут стоять как вплотную, так и раздельно)

Переменные среды:
  Для каждого вида вычислений можно установить время минимального выполнения в миллисекундах через переменные среды(по умолчанию значения равны 100) 
   - TIME_ADDITION_MS // суммирование
   - TIME_SUBTRACTION_MS // вычитание
   - TIME_MULTIPLICATIONS_MS // умножение
   - TIME_DIVISIONS_MS // деление

 Для количества возможных агентов
   - COMPUTING_POWER // по умолчанию равен 10



!если что-то пошло не так! -> telegram: @gduebsdh1

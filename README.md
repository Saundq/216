# Распределенный вычислитель арифметических выражений

Состоит из Front-end и Back-end частей. В свою очередь Back-end состоит из 2 элементов
оркестратора и агента.

## Установка и запуск
Для упрощения установки и запуска можно воспользоваться [docker и docker-compose](https://docs.docker.com/get-docker/).
В корне проекта находится docker-compose файл который развернет postgres и node, иначе для запуска их необходимо развернуть самостоятельно и отредактировать .env для корректного подключения.


- Клонируем проект к себе
    ```
    git clone 
     ```
- Установка frontend
  
    Для установки можно воспользоваться makefile:
    ```
    make frontend-install
    ```
    или

    запустить из папки frontend 

    ```
    npm i
    ```

- В отдельном терминале запуск  оркестратора
    ```
    make run-orc
    ```
  или из папки backend проекта запустить
    ```
    go run cmd/orchestrator/main.go
    ```
-  В отдельном терминале запуск агента
    ```
    make run-agent
    ```
  или из папки backend проекта запустить
    ```
    go run cmd/agent/main.go
    ```
-  В отдельном терминале запуск frontend
   ```
   make run-frontend
   ```
   или из папки frontend проекта запустить
   ```
   npm run dev
   ```

## Принцип работы
По адресу [доступен frontend](http://localhost:5173/)
Позволяющий отправить идемпотентный запрос на сервер с математическим выражением на решение. Идемпотентность достигается добавлением в заголовок запроса X-Request-ID с хешем математического выражения. При повторной отправке выражения сервер вернет 200 или 200 с ответом.
Данные отправляются POST запросом на сервер

По адресу [cтраница со списком выражений](http://localhost:5173/expressions) с информацией о
статусе, дате создания и заверщения вычисления
Страница получает данные GET запросом

По адресу [cтраница со списком операций](http://localhost:5173/operations) с информацией о имени операции и времени его выполнения (доступное для редактирования поле)
Страница получает данные GET запросом

По адресу [страница со списком вычислительных можностей](http://localhost:5173/computing_capabilities) с информацией о имени вычислительного ресурса и выполняемой на нём операции
Страница получает данные GET запросом


При добавлении математического выражения (вводится без пробелов) на вычисление, оркестратор переводит его в постфиксную форму и разбивает на мелкие дочерние мат. выражения. Получаем не зависимые и зависимы выражения.

Агент раз в единицу времени делает опрос оркестратора на задачи вычисления. Оркестратор на запрос задачи агента первоочередно отдает не зависимые выражения. После получения задачи проставляется информация для вычислительного реурса, выполняется вычисление с указанной задержкой и проставляется результат c округлением до второго знака.
Если посчитаны все дочерние мат. выражения, то проставляется результат для всего выражения.

При остановке агента и последующем его запуске все выражения со статусом progress возвращаются в статут wait и сново доступны для вычисления.

Агент раз в единицу времени отправляет временную метку о своей доступности. Если интервал превысит 10 секунд статус вычислителей изменится на не достпун, если 20 секунд произойдет удаление из доступных вычислителей. Кол-во вычислителей регулируется файлом окружения параметром NUMBER_OF_COMPUTERS 

Все endpoint-ы оркестратора описаны в файле api.http

## Схема БД
![img.png](img.png)

## Примеры

| Мат. выражение                      | Результат |
|-------------------------------------|-----------|
| (3+3)*4/1+1-1+(3-2)-8-1             | 16        |
| (3+3)*4/1+1-1-(6+1)+8               | 25        |
| 9/5+6*7+23-3                        | 63.8      |
| 3-1+(4*53)/4-3+((23-1)+(7+3)*2)     | 94        |        
| ((7+3)/3)-4+(15-4)/3+21*4-((8-3)*2) | 77        |        
| 2*2                                 | 4         |         
| 2*2+2                               | 6         |         
| 2*(2+1)                             | 6         |         
| 2*((2-1)*3)                         | 6         |         




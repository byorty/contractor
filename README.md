# Contractor - HTTP API testing tool

Contractor читает спецификацию вашего API чтобы:
1. Создать mock-сервер
2. Валидировать backend-имплементацию

Поддерживаемые форматы спецификации API:
1. [OpenAPI 2](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md) (также известный как Swagger)
2. [OpenAPI 3](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md)

Поддерживаемые расширения файлов спецификации API:
1. [YAML](https://yaml.org)
2. [JSON](https://www.json.org/json-en.html)

Поддерживаемые платформы:
1. macOS
2. Linux

## Установка

С помощью go install:
```shell
go install github.com/byorty/contractor/cmd/...
```

С помощью Homebrew:
```shell
brew install byorty/tap/contractor 
```

## Быстрый старт

1. Опишите ваше API [например в формате OpenAPI 2](https://github.com/byorty/contractor/blob/master/specs/oa2.yml):
   ```yaml
   swagger: "2.0"
   info:
     title: Test API
     version: "1.0"
   schemes:
     - https
   consumes:
     - application/json
   produces:
     - application/json
   paths:
     # ...
     # описание других методов API
     # ...
     /v1/users/{user_id}/news/{news_id}:
       patch:
         operationId: Users_ChangeUserNews
         responses:
           "200":
             description: A successful response.
             schema:
               $ref: '#/definitions/UserNews'
         parameters:
         - name: user_id
           in: path
           required: true
           type: integer
           format: int64
         - name: news_id
           in: path
           required: true
           type: integer
           format: int64
         - name: body
           in: body
           required: true
           schema:
              $ref: '#/definitions/ChangeUserNewsRequestParams'
         tags:
         - Users
   definitions:
     # ...
     # описание сущностей API
     # ...
   ```
2. Добавьте расширение `x-examples`:
   ```yaml
   /v1/users/{user_id}/news/{news_id}:
     get:
       operationId: Users_ChangeUserNews
       # ...
       # описание метода API
       # ...
       x-examples: относительный/путь/до/примеров/запросов/и/ответов.yml
   ```
3. Опишите примеры запросов и ответов метода API:
   ```yaml
   CHANGE_USER_NEWS_SUCCESS:                        # уникальное имя в рамках всей спецификации
     priority: 1                                    # приоритет
     tags:                                          # поддержка тэгов
       - user                                       #
       - develop                                    #
       - production                                 #
     request:                                       # описание запроса
       headers:                                     # описание заголовков запроса
         Authorization: ${VAR_AUTHORIZATION}        # поддержка переменных
       parameters:                                  # описание параметров запроса
         user_id: 472                               # имена задаются так, как описаны в спецификации API
         news_id: 11401                             #
         body:                                      # описание тела запроса для POST, PUT, PATCH запросов
           is_viewed: true                          #
           reaction: USER_NEWS_REACTION_LIKE        #
     response:                                      # описание ответа
       status_code: 200                             # статус обязателен
       body:                                        # описание тела ответа
         id: eq(30676003)
         user:
           id: eq(472)
         news:
           id: eq(11401)
           title: >
             eq('Домашняя работа: как стартап бывшего топ-менеджера Microsoft зарабатывает на покупке жилья за наличные')
           annotation: regex('([\\w\\d\\-\\.\\,\\?\\:]){100}')
           partner: empty()
           content: regex('([\\w\\d\\-\\.])')
           hash: eq('ee43501beb0944c412c9580ae604546f')
           preview_img: eq('977081c0dc136761a13d60c513437dbf')
           tags: empty()
           status: eq('NEWS_STATUS_ACTIVE')
           type: eq('CONTENT_TYPE_NEWS')
           # ...
           # описание свойств ответа
           # ...
   ```
4. Создайте mock-сервер:
   ```shell
   contractor  -m mock \
               -s ./specs/oa2.yml \
               -u http://localhost:8181 \
               -f oa2 \
               -v "VAR_AUTHORIZATION: Bearer some-jwt"
   ```
6. Провалидируйте свой сервер:
   ```shell
   contractor  -m test \
               -s ./specs/oa2.yml \
               -u http://localhost:8181 \
               -f oa2 \
               -t tag1 \
               -t tag2 \
               -v "VAR_AUTHORIZATION: Bearer some-jwt"
   ```

## Использование
```shell
contractor [OPTIONS]
```
### Аргументы командной строки
`-m, --mode`: обязательный, режим работы contractor, возможные значения:
1. test - валидация API по указанной спецификации
2. mock - создание mock-сервера по указанной спецификации

`-s, --spec-filename`: обязательный, путь по спецификации API

`-f, --spec-format`: опциональный, формат спецификации API, возможные значения:
1. oa2 - для [OpenAPI 2](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md), значение по умолчанию
2. oa3 - для [OpenAPI 3](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md)

`-u, --url`: обязательный, хост и порт, на котором работает имплементация вашего API, либо должен работать mock-сервер 

`-v, --variable`: опциональный, переменные доступные в примерах запросов и ответов

`-t, --tags`: опциональный, список тэгов

`-c, --cert-filename`: опциональный, путь до сертификата, если mock-сервер должен работать по https 

`-k, --key-filename`: опциональный, путь до ключа, если mock-сервер должен работать по https

### Список функций
| Название                | Режим тестирования                                                                                                      | Режим mock-сервера                                                                                                             |
|-------------------------|-------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------|
| eq(value)               | Проверяет на строгое соответствие свойство объекта                                                                      | Присваивается указанное значение свойству объекта                                                                              |
| positive()              | Проверяет что свойство объекта положительное число                                                                      | Присваивается свойству объекта положительное число                                                                             |
| min(number)             | Проверяет что свойство объекта > number                                                                                 | Присваивается свойству объекта число > number                                                                                  |
| max(number)             | Проверяет что свойство объекта < number                                                                                 | Присваивается свойству объекта число < number                                                                                  |
| range(number1, number2) | Проверяет что свойство объекта number1 < объекта < number2                                                              | Присваивается свойству объекта number1 < объекта < number2                                                                     |
| regex(exp)              | Проверяет что свойство объекта соответствует регулярному выражению                                                      | Присваивается свойству объекта генерируемую строку на основе регулярного выражения                                             |
| empty()                 | Проверяет что свойство объекта пустое                                                                                   | Присваивается свойству объекта null                                                                                            |
| date(format)            | Проверяет что свойство дата в указанном формате.  Возможные значения format: RFC33392, RFC3339NANO, произвольный формат | Присваивается свойству объекта дату в указанном формате. Возможные значения format: RFC33392, RFC3339NANO, произвольный формат |
| contains(substring)     | Проверяет что подстрока содержится в свойстве объекта                                                                   | Присваивается свойству объекта строку, в которой содержится указанная подстрока                                                |

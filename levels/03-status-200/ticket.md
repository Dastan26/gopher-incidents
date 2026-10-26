# GOPHER-3 · Ошибка валидации приходит со статусом 200

Приоритет: критичный · Темы: http.ResponseWriter, порядок заголовков и тела, коды ответа

## Что видит пользователь

Менеджер создаёт бренд с пустым казахским названием. Фронт показывает
тост «Бренд сохранён», но в списке бренда нет. Ошибка под полем
не появляется.

Для фронта, который проверяет response.ok (как axios в нашем портале),
любой 2xx — это успех. Ошибку он даже не пытается читать.

## Что в DevTools

```
POST /brands
Status: 200 OK
Content-Type: text/plain; charset=utf-8

{"errors":[{"status":"422","detail":"Обязательное поле","source":{"pointer":"/data/attributes/nameKz"}}]}
```

Тело правильное, а статус и Content-Type нет. Успешное создание
тоже приходит с 200 вместо 201.

В логе сервиса на каждый запрос:

```
http: superfluous response.WriteHeader call from ...
```

## Воспроизвести

```
make run
curl -i -X POST localhost:8080/brands -d '{"data":{"type":"brands","attributes":{"nameKz":"","nameRu":"Тест"}}}'
```

## Критерии готовности

- Ошибка валидации: 422 и Content-Type application/vnd.api+json.
- Успешное создание: 201 и тот же Content-Type.
- Список и карточка: 200, а несуществующий бренд — 404.
- В логе нет «superfluous response.WriteHeader call».
- make check зелёный.

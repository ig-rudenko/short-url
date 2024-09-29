# Short URL

[![Used](https://skillicons.dev/icons?i=go,sqlite,docker)](https://skillicons.dev)

## API

#### POST `/url` - Создание ссылки:

```json
{
  "url": "string",
  "alias": "string"
} 
```

#### GET `/{alias}` - Получение оригинала ссылки.

#### DELETE `/{alias}` - Удаление ссылки.
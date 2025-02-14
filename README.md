# Система постов и комментариев

## Описание

Система для добавления и чтения постов и комментариев. Система реализована на языке Go с использованием GraphQL и поддерживает как in-memory, так и PostgreSQL для хранения данных.

### Основные характеристики

#### Система постов:

- Просмотр списка постов.
- Просмотр поста и комментариев под ним.
- Пользователь, написавший пост, может запретить оставлять комментарии под своим постом.

#### Система комментариев:

- Комментарии организованы иерархически, позволяя вложенность без ограничений.
- Длина текста комментария ограничена до 2000 символов.
- Пагинация для получения списка комментариев под постом.

#### Дополнительные требования:

- Поддержка GraphQL Subscriptions для асинхронного получения новых комментариев.

## Запуск

Переменные окружения хранятся в `.env` файле.

### Для запуска с PostgreSQL

```
docker-compose -f docker-compose.postgres.yml up
```

### Для запуска c in-memory хранилищем

```
docker-compose -f docker-compose.inmemory.yml up
```

## Примеры запросов

### Создание поста

```graphql
mutation CreatePost {
  CreatePost(
    input: {
      title: "First post"
      author: "Bob"
      content: "something."
      commentsAllowed: true
    }
  ) {
    id
    title
    author
    content
    createdAt
    commentsAllowed
  }
}
```

### Получение списка постов

```graphql
query GetPosts {
  GetPosts {
    id
    title
    author
    content
    createdAt
    commentsAllowed
  }
}
```

### Получение детальной информации о посте

```graphql
query GetPostByID {
  GetPostByID(id: "1", page: 1, pageSize: 10) {
    id
    title
    author
    content
    createdAt
    commentsAllowed
    comments {
      id
      author
      content
      createdAt
      postID
      parentID
      replies {
        id
        author
        content
        createdAt
        postID
        parentID
      }
    }
  }
}
```

### Создание комментария

```graphql
mutation CreateComment {
  CreateComment(input: { postID: "1", author: "Bob", content: "comment" }) {
    id
    author
    content
    createdAt
    postID
    parentID
  }
}
```

### Подписка на комментарии поста

```graphql
subscription CommentsSubscription {
  CommentAdded(postId: "1") {
    id
    author
    content
    createdAt
    postID
    parentID
  }
}
```

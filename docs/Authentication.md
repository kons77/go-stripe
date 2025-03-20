# Authentication 

(Section 6)

In this section of the course, we will be implementing authentication, a crucial aspect of any secure web application. Authentication ensures that users are valid, active, and have the necessary permissions to access the system. The approach to authentication varies depending on the type of web application.

For traditional websites without an API, session-based authentication is common. The process involves:

- The user filling out a login form.
- The backend validating the provided credentials.
- If successful, setting a session variable (e.g., `authenticated = true`).
- Validating the session variable on every request until the session expires or the user logs out.

For API-based authentication, tokens are generally preferred. Several types of token-based authentication exist, each with its own advantages and challenges:

1. **HTTP Basic Authentication**
   - Simple to implement.
   - Secure if used over HTTPS.
   - Generally slow and inefficient for high-traffic applications.
   - Requires credentials to be sent with every request, increasing security risks.
2. **Simple Tokens**
   - The backend issues a token upon successful authentication.
   - The client includes the token in requests for access validation.
   - The token must be validated on each request.
   - Requires manual handling of token expiration and revocation.
3. **Stateful Tokens**
   - A validated token (or a hash of it) is stored in a database along with an expiry timestamp.
   - Tokens can be invalidated by removing them from the database.
   - Provides control over user sessions and allows selective invalidation.
   - Requires backend storage and management.
4. **JSON Web Tokens (JWTs)**
   - Stateless authentication method.
   - The token itself contains the user’s authentication information and expiry timestamp.
   - Popular due to its self-contained nature and ease of use.
   - Major drawback: Tokens cannot be revoked individually. If a user needs to be deauthorized immediately, all tokens must be invalidated at once.
   - Requires additional client-side logic to refresh and validate tokens.
5. **API Keys**
   - Used primarily for service-to-service authentication.
   - Common in third-party integrations (e.g., GitHub APIs).
   - Not suited for user authentication due to lack of session control.
6. **OAuth 2.0**
   - Allows users to authenticate using third-party services (e.g., Google, Facebook, GitHub).
   - Reduces the need to manage passwords directly.
   - Requires users to have accounts with the authentication provider.
   - Common in applications that support social login.

For this course, we will implement **stateful tokens**, as they offer a balance between security and flexibility. This approach ensures that we can manage user sessions efficiently while retaining control over authentication validity. In the upcoming lessons, we will explore how to implement stateful token authentication step by step.

***

## func GenerateToken - explanation 

#### **1. Разбор `hash[:]`**

```
goCopyEdithash := sha256.Sum256([]byte(token.PlainText))
token.Hash = hash[:]
```

**Что здесь происходит?**

1. `sha256.Sum256([]byte(token.PlainText))` — вычисляет хеш SHA-256 для `token.PlainText`.
   - `sha256.Sum256` возвращает **массив фиксированной длины** `[32]byte` (массив из 32 байт).
   - Это важно! В Go массив `[32]byte` отличается от среза `[]byte`.
2. `hash[:]` — это **срез всего массива** `[32]byte`, который превращается в `[]byte`.
   - `[:]` означает "взять все элементы массива и создать срез".
   - Это необходимо, потому что `token.Hash` имеет тип `[]byte`, а не `[32]byte`.

**Почему это важно?**

В Go массивы `[N]T` и срезы `[]T` — это разные типы. `[:]` помогает преобразовать массив `[32]byte` в `[]byte`, что упрощает работу с данными.

#### **2. Что такое `base32`?**

`base32` — это способ кодирования бинарных данных в строку, используя 32 символа из ASCII.

- Он удобен, потому что символы хорошо читаются человеком и легко передаются в URL.
- В отличие от `base64`, он не использует `+` и `/`, а вместо них — буквы.
- Кодирование делает строку длиннее, но она получается удобной для копирования.

Пример:

```
goCopyEditdata := []byte("hello")
encoded := base32.StdEncoding.EncodeToString(data)
fmt.Println(encoded) // "NBSWY3DP"
```

------

#### **3. Что такое `.WithPadding(base32.NoPadding)`?**

По умолчанию `base32` использует символ `=` для выравнивания длины строки.
Функция `.WithPadding(base32.NoPadding)` убирает эти `=`.

**Пример:**

```
goCopyEditdata := []byte("hello")
encoded := base32.StdEncoding.EncodeToString(data)
fmt.Println(encoded) // "NBSWY3DP===="

encodedNoPad := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(data)
fmt.Println(encodedNoPad) // "NBSWY3DP"
```

Такой вариант удобнее, если строку нужно передавать в URL.

****

## JS as main auth method 

If you choose to use this authentication method, there is nothing inherently wrong with it. To implement it properly, you need to adjust the placement of the `checkAuth` function within the HTML structure.

Currently, in the base layout file, `checkAuth` is located at the bottom. However, if you intend to use this as your primary authentication method, it should be moved into the `<head>` section, right after the `<title>` tag. This ensures that authentication is checked before the rest of the page is rendered.

Instead of calling `checkAuth` at the bottom of the terminal page, it should be placed within a new block in the `<head>`. The implementation would involve wrapping the function in a `<script>` tag within the head section, ensuring that it is executed as soon as the browser processes the document.

The reason for this approach lies in how browsers work. When rendering a page, the browser immediately halts processing upon encountering a script tag in the `<head>`, executes the JavaScript, and then continues loading the rest of the document. By placing `checkAuth` in the head, the user will not see any content until authentication has been verified.

If the authentication check fails—whether because the token is missing, expired, or invalid—the user will be redirected to the login page before any content is displayed. This guarantees that unauthorized users do not gain access to protected parts of the application.

This is one approach to handling authentication, ensuring that users are validated before they interact with the page.

```html
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{block "title" .}}{{end}}</title>

  <script>
     checkAuth()
  </script>
  {{block "in head" .}} {{end}}
</head>
```






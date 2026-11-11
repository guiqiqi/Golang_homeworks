# ДЗ к уроку 5. Ошибки в сервисе уведомлений

Продолжите сервис из Практики №4. Добавьте обработку ошибок и сделайте так,
чтобы сбой одного отправителя не мешал остальным.

Нужен Go 1.22 или новее. Используйте стандартную библиотеку, без горутин.

## Что нужно сделать

### 1. Проверить пустое сообщение

В начале `EmailNotifier.Send` и `PushNotifier.Send` проверяйте `msg == ""`.
Если сообщение пустое, возвращайте объявленную в заготовке ошибку `ErrEmptyMessage`.
Проверяйте это раньше `DeviceID`. Пробелы обрезать не нужно.

При успешной отправке выводите канал, получателя и сообщение, возвращайте `nil`.

### 2. Добавить свою ошибку для push

В заготовке есть структура `SendError` с полями `Channel` и `Reason`.

- Реализуйте `func (e *SendError) Error() string`: верните текст с каналом и причиной.
- В `PushNotifier.Send` после проверки сообщения проверьте `DeviceID == 0`.
- Если ID не задан, верните `&SendError`: канал возьмите из `p.BaseNotifier.Channel`,
  причина: `"не задан device id"`.

### 3. Перехватить панику при отправке

Допишите функцию:

```go
func SafeSend(notifier Notifier, msg string) (err error)
```

Перед вызовом `notifier.Send(msg)` добавьте `defer` с функцией, которая вызывает
`recover()`. Если произошла паника, запишите в `err` ошибку с её описанием.
Здесь можно использовать `%v`.

При обычной отправке возвращайте результат `Send` без изменений.
Перехватывайте панику отдельно для каждого вызова, чтобы рассылка продолжалась.

### 4. Собрать ошибки в диспетчере

Допишите функцию:

```go
func Dispatch(all []Channel, msg string) []error
```

- Вызывайте всех отправителей по порядку через `SafeSend`.
- Каждую ошибку оборачивайте: `fmt.Errorf("канал %s: %w", channel.Name, err)`.
- Собирайте ошибки в срез в порядке вызовов. После ошибки продолжайте рассылку.
- Верните все собранные ошибки. При успехе длина среза должна быть нулевой.
- Для пустого или `nil`-среза каналов возвращайте результат без ошибок.

Имя канала берите из `Channel.Name`. Обрабатывать `nil` внутри `Notifier` не нужно.
Сигнатуры и структуры заготовки не меняйте, тесты не ослабляйте.
Ошибки ввода возвращайте через `error`, без `panic`. Ошибки нужно возвращать,
а не только печатать. Не скрывайте сбои за `nil`.

### 5. Показать три сценария в `main`

| Сценарий | Какой результат ожидается |
|---|---|
| Email, затем `BrokenNotifier`, затем корректный push; сообщение `"Заказ готов"` | Email и push отправили сообщение. Диспетчер вернул одну ошибку с именем `broken` и описанием паники. |
| Email и push; сообщение `""` | Диспетчер вернул две ошибки. Для каждой `errors.Is(err, ErrEmptyMessage)` вернул `true`. Выведено пояснение, что нужно ввести текст. |
| Push с `DeviceID == 0`; непустое сообщение | Диспетчер вернул ошибку. Через `errors.As(err, &sendErr)` получены и выведены поля `Channel` и `Reason`, где `var sendErr *SendError`. |

В сценарий с пустым сообщением `BrokenNotifier` не добавляйте.
Проверяйте каждый элемент `[]error` отдельно. Поля `sendErr` читайте только после
успешного `errors.As`.

Проверьте, что `errors.Is` находит `ErrEmptyMessage` после дополнительной обёртки
через `%w`. Затем временно замените `%w` на `%v` и сравните результат.
Перед сдачей верните `%w`. Не сравнивайте ошибки по тексту.

## Как проверить

Из каталога урока 5 выполните:

```bash
cd homework
gofmt -w main.go
go run .
go test -v ./...
go vet ./...
```

- В выводе `go run .` видны все три сценария из таблицы.
- Все тесты проходят, `go vet` не сообщает об ошибках.
- После ошибки или паники следующий отправитель получает сообщение.

Тесты проверяют ошибки, порядок проверок и вызовов, продолжение рассылки
и пустой список каналов. Вывод успешных отправок и работу `Is`/`As` в `main`
проверьте вручную.

## Доп. вопросики

В README оставьте команды запуска, описание трёх сценариев и короткие ответы:

1. Чем `errors.Is` отличается от `errors.As`?

   `errors.Is` is used to determine whether a specific error exists within an error chain, while `errors.As` is used to locate an error of a specific type within an error chain and assign it to a variable, allowing access to that error's fields.

2. Что потеряется при замене `%w` на `%v`?

   `%v` merely converts the original error into text; it does not preserve the error chain. Consequently, `errors.Is` and `errors.As` can no longer recognize the original error or its type, even though the textual description of the error remains in the message.

3. Почему одного `recover` вокруг цикла недостаточно для продолжения рассылки?

   If the `recover` mechanism is placed outside the entire loop, a panic occurring during any `Send` operation will immediately terminate the loop, with control passing only to the handler located outside the loop. By placing `recover` within `SafeSend`, you can isolate a panic to the specific send operation in which it occurs, allowing the loop to continue and proceed to the next sender.

4. Что означает частичный успех, отменяются ли уже выполненные отправки?

   "Partial success" indicates that some notifications were sent successfully, while errors or panics occurred for others. Notifications that have already been processed will not be cancelled or rolled back; the scheduler will continue with the remaining tasks and return a list of errors encountered for the failed attempts.

package main

import (
	"errors"
	"fmt"
)

var ErrEmptyMessage = errors.New("пустое сообщение")

type Notifier interface {
	Send(msg string) error
}

type BaseNotifier struct {
	Channel string
}

func (b BaseNotifier) Log() {
	fmt.Println("Sending notification via channel:", b.Channel)
}

type EmailNotifier struct {
	BaseNotifier
	Address string
}

func (e EmailNotifier) Send(msg string) error {
	if msg == "" {
		return ErrEmptyMessage
	}
	fmt.Printf("канал: %s, получатель: %s, сообщение: %s\n", e.Channel, e.Address, msg)
	return nil
}

var _ Notifier = EmailNotifier{}

type PushNotifier struct {
	BaseNotifier BaseNotifier
	DeviceID     uint64
}

func (p PushNotifier) Send(msg string) error {
	if msg == "" {
		return ErrEmptyMessage
	}
	if p.DeviceID == 0 {
		return &SendError{
			Channel: p.BaseNotifier.Channel,
			Reason:  "не задан device id",
		}
	}
	fmt.Printf("канал: %s, получатель: устройство %d, сообщение: %s\n", p.BaseNotifier.Channel, p.DeviceID, msg)
	return nil
}

var _ Notifier = PushNotifier{}

type SendError struct {
	Channel string
	Reason  string
}

func (e *SendError) Error() string {
	return fmt.Sprintf("канал %s: %s", e.Channel, e.Reason)
}

var _ error = (*SendError)(nil)

type BrokenNotifier struct{}

func (BrokenNotifier) Send(string) error {
	panic("драйвер канала повреждён")
}

type Channel struct {
	Name   string
	Sender Notifier
}

func SafeSend(notifier Notifier, msg string) (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = fmt.Errorf("паника при отправке: %v", value)
		}
	}()
	return notifier.Send(msg)
}

func Dispatch(all []Channel, msg string) []error {
	failures := make([]error, 0, len(all))
	for _, channel := range all {
		if err := SafeSend(channel.Sender, msg); err != nil {
			failures = append(failures, fmt.Errorf("канал %s: %w", channel.Name, err))
		}
	}
	return failures
}

func main() {
	email := EmailNotifier{
		BaseNotifier: BaseNotifier{Channel: "email"},
		Address:      "example@mail.ru",
	}
	push := PushNotifier{
		BaseNotifier: BaseNotifier{Channel: "push"},
		DeviceID:     123456,
	}

	fmt.Println("Сценарий 1: email, сломанный канал и push")
	channels := []Channel{
		{Name: "email", Sender: email},
		{Name: "broken", Sender: BrokenNotifier{}},
		{Name: "push", Sender: push},
	}

	failures := Dispatch(channels, "Заказ готов")
	for _, err := range failures {
		fmt.Println("ошибка:", err)
	}
	fmt.Println("Всего ошибок:", len(failures))

	fmt.Println()
	fmt.Println("Сценарий 2: пустое сообщение")
	failures = Dispatch([]Channel{
		{Name: "email", Sender: email},
		{Name: "push", Sender: push},
	}, "")
	for _, err := range failures {
		fmt.Println("ошибка:", err)
		if errors.Is(err, ErrEmptyMessage) {
			fmt.Println("errors.Is: пустое сообщение найдено — нужно ввести текст")
		}
	}
	fmt.Println("Всего ошибок:", len(failures))

	fmt.Println()
	fmt.Println("Сценарий 3: push без DeviceID")
	failures = Dispatch([]Channel{
		{Name: "push", Sender: PushNotifier{BaseNotifier: BaseNotifier{Channel: "push"}}},
	}, "Заказ готов")
	for _, err := range failures {
		fmt.Println("ошибка:", err)
		var sendErr *SendError
		if errors.As(err, &sendErr) {
			fmt.Printf("errors.As: Channel=%q, Reason=%q\n", sendErr.Channel, sendErr.Reason)
		}
	}
	fmt.Println("Всего ошибок:", len(failures))

	fmt.Println()
	fmt.Println("Проверка сохранения причины после обёртки")
	wrapped := fmt.Errorf("дополнительная обёртка: %w", ErrEmptyMessage)
	fmt.Println("errors.Is после обёртки через %w:", errors.Is(wrapped, ErrEmptyMessage))
}

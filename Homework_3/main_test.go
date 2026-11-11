package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestEmptyMessage(t *testing.T) {
	notifiers := map[string]Notifier{
		"email":               EmailNotifier{BaseNotifier: BaseNotifier{Channel: "email"}, Address: "example@mail.ru"},
		"push":                PushNotifier{BaseNotifier: BaseNotifier{Channel: "push"}, DeviceID: 123456},
		"push without device": PushNotifier{BaseNotifier: BaseNotifier{Channel: "push"}},
	}
	for name, notifier := range notifiers {
		t.Run(name, func(t *testing.T) {
			err := notifier.Send("")
			if err == nil || !errors.Is(err, ErrEmptyMessage) {
				t.Fatalf("для пустого сообщения нужна ErrEmptyMessage, получено: %v", err)
			}
		})
	}
}

func TestSuccessfulSend(t *testing.T) {
	notifiers := []Notifier{
		EmailNotifier{BaseNotifier: BaseNotifier{Channel: "email"}, Address: "example@mail.ru"},
		PushNotifier{BaseNotifier: BaseNotifier{Channel: "push"}, DeviceID: 123456},
	}
	for _, notifier := range notifiers {
		if err := notifier.Send("Заказ готов"); err != nil {
			t.Errorf("настроенный %T вернул ошибку: %v", notifier, err)
		}
	}
}

func TestSendError(t *testing.T) {
	push := PushNotifier{BaseNotifier: BaseNotifier{Channel: "mobile"}}
	err := push.Send("Заказ готов")
	var target *SendError
	if !errors.As(err, &target) || target == nil {
		t.Fatalf("ожидается *SendError для нулевого DeviceID, получено: %v", err)
	}
	if target.Channel != "mobile" || target.Reason != "не задан device id" {
		t.Fatalf("неверные поля: Channel=%q, Reason=%q", target.Channel, target.Reason)
	}
	if text := target.Error(); !strings.Contains(text, target.Channel) || !strings.Contains(text, target.Reason) {
		t.Fatalf("описание должно содержать канал и причину, получено: %q", text)
	}
}

type stubNotifier struct {
	name       string
	calls      *[]string
	err        error
	panicValue string
}

func (n stubNotifier) Send(msg string) error {
	if n.calls != nil {
		*n.calls = append(*n.calls, n.name+":"+msg)
	}
	if n.panicValue != "" {
		panic(n.panicValue)
	}
	return n.err
}

func TestSafeSendPreservesResult(t *testing.T) {
	want := errors.New("канал недоступен")
	for _, original := range []error{nil, want} {
		if got := SafeSend(stubNotifier{err: original}, "hello"); got != original {
			t.Errorf("обычный результат Send изменился: получено %v, ожидается %v", got, original)
		}
	}
}

func TestSafeSendPanic(t *testing.T) {
	defer func() {
		if value := recover(); value != nil {
			t.Errorf("паника вышла из SafeSend: %v", value)
		}
	}()
	err := SafeSend(BrokenNotifier{}, "hello")
	if err == nil || !strings.Contains(err.Error(), "драйвер канала повреждён") {
		t.Fatalf("нужна ошибка с описанием перехваченной паники, получено: %v", err)
	}
}

func TestDispatchSuccess(t *testing.T) {
	var calls []string
	all := []Channel{
		{Name: "first", Sender: stubNotifier{name: "first", calls: &calls}},
		{Name: "last", Sender: stubNotifier{name: "last", calls: &calls}},
	}
	if failures := Dispatch(all, "hello"); len(failures) != 0 {
		t.Fatalf("успешные каналы вернули ошибки: %v", failures)
	}
	if !slices.Equal(calls, []string{"first:hello", "last:hello"}) {
		t.Fatalf("неверные вызовы: %v", calls)
	}
}

func TestDispatchPreservesAllErrors(t *testing.T) {
	var calls []string
	custom := &SendError{Channel: "push", Reason: "не задан device id"}
	all := []Channel{
		{Name: "email", Sender: stubNotifier{name: "email", calls: &calls, err: ErrEmptyMessage}},
		{Name: "push", Sender: stubNotifier{name: "push", calls: &calls, err: custom}},
		{Name: "last", Sender: stubNotifier{name: "last", calls: &calls}},
	}
	failures := Dispatch(all, "hello")
	if len(failures) != 2 {
		t.Fatalf("ожидаются две ошибки, получено: %v", failures)
	}
	if !slices.Equal(calls, []string{"email:hello", "push:hello", "last:hello"}) {
		t.Fatalf("после ошибки нужно продолжить вызовы по порядку: %v", calls)
	}
	if !errors.Is(failures[0], ErrEmptyMessage) {
		t.Fatalf("потеряна известная причина: %v", failures[0])
	}
	if !errors.Is(fmt.Errorf("рассылка: %w", failures[0]), ErrEmptyMessage) {
		t.Fatal("причина не найдена после дополнительного оборачивания")
	}
	var target *SendError
	if !errors.As(failures[1], &target) || target != custom {
		t.Fatalf("потеряна собственная ошибка: %v", failures[1])
	}
	for i, name := range []string{"email", "push"} {
		if failures[i] == nil || !strings.Contains(failures[i].Error(), "канал "+name+":") {
			t.Errorf("в ошибке нет контекста канала %q: %v", name, failures[i])
		}
	}
}

func TestDispatchEmptyMessage(t *testing.T) {
	all := []Channel{
		{Name: "email", Sender: EmailNotifier{BaseNotifier: BaseNotifier{Channel: "email"}}},
		{Name: "push", Sender: PushNotifier{BaseNotifier: BaseNotifier{Channel: "push"}, DeviceID: 1}},
	}
	failures := Dispatch(all, "")
	if len(failures) != 2 {
		t.Fatalf("каждый канал должен вернуть ошибку пустого сообщения: %v", failures)
	}
	for _, err := range failures {
		if err == nil || !errors.Is(err, ErrEmptyMessage) {
			t.Errorf("потеряна причина пустого сообщения: %v", err)
		}
	}
}

func TestDispatchContinuesAfterPanic(t *testing.T) {
	defer func() {
		if value := recover(); value != nil {
			t.Errorf("паника вышла из Dispatch: %v", value)
		}
	}()
	var calls []string
	wantErr := errors.New("последний канал недоступен")
	all := []Channel{
		{Name: "first", Sender: stubNotifier{name: "first", calls: &calls}},
		{Name: "broken", Sender: stubNotifier{name: "broken", calls: &calls, panicValue: "сломанный драйвер"}},
		{Name: "last", Sender: stubNotifier{name: "last", calls: &calls, err: wantErr}},
	}
	failures := Dispatch(all, "hello")
	if !slices.Equal(calls, []string{"first:hello", "broken:hello", "last:hello"}) {
		t.Fatalf("после паники нужно вызвать последний канал: %v", calls)
	}
	if len(failures) != 2 {
		t.Fatalf("нужны ошибка паники и ошибка последнего канала, получено: %v", failures)
	}
	if failures[0] == nil || !strings.Contains(failures[0].Error(), "канал broken:") ||
		!strings.Contains(failures[0].Error(), "сломанный драйвер") {
		t.Errorf("потеряны имя канала или описание паники: %v", failures[0])
	}
	if !errors.Is(failures[1], wantErr) {
		t.Errorf("потеряна ошибка канала после паники: %v", failures[1])
	}
}

func TestDispatchEmptyList(t *testing.T) {
	for _, all := range [][]Channel{nil, {}} {
		if failures := Dispatch(all, "hello"); len(failures) != 0 {
			t.Errorf("пустой список не должен возвращать ошибки: %v", failures)
		}
	}
}

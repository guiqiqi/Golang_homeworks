package main

import "testing"

func TestStudentScores(t *testing.T) {
	scores := StudentScores()
	if scores == nil {
		t.Fatal("StudentScores вернула nil; ожидается созданная map с баллами")
	}
	for name, want := range map[string]int{"Анна": 5, "Борис": 0} {
		got, ok := scores[name]
		if !ok {
			t.Errorf("в map отсутствует студент %s", name)
		} else if got != want {
			t.Errorf("баллы %s = %d, ожидается %d", name, got, want)
		}
	}
	for _, name := range []string{"Вера", "Глеб"} {
		if score, ok := scores[name]; ok {
			t.Errorf("студент %s не должен присутствовать в map, получены баллы %d", name, score)
		}
	}
	if len(scores) != 2 {
		t.Errorf("после удаления Веры длина map = %d, ожидается 2", len(scores))
	}
}

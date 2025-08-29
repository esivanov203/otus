package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("single element list", func(t *testing.T) {
		l := NewList()
		elem := "последний элемент"
		item := l.PushBack(elem)          // вставка строки в конец
		require.Equal(t, item, l.Back())  // в списке из одного элемента -
		require.Equal(t, item, l.Front()) // front и back - один и тот же элемент

		l.MoveToFront(item)
		require.Equal(t, item, l.Front()) // проверка перемещения

		l.Remove(item)
		require.Equal(t, 0, l.Len()) // Проверка удаления
		// ps. Проверку перемещения и удаления в пустом списке не проверяем -
		// по условию ТЗ - это не нужно реализовывать
	})

	t.Run("difference types elements list", func(t *testing.T) {
		l := NewList()

		tstIntEl := 10
		tstStruct := struct {
			id   int
			name string
		}{id: 1, name: "Мое имя"}

		l.PushFront(tstIntEl)
		tstStructItem := l.PushBack(tstStruct) // вставка структуры
		l.MoveToFront(tstStructItem)           // перемещение структуры

		require.Equal(t, 2, l.Len())
		require.Equal(t, tstIntEl, l.Back().Value)
		require.Equal(t, tstStruct, l.Front().Value)

		l.Remove(l.Front()) // удаление структуры (первого элемента)

		require.Equal(t, 1, l.Len())                // структура удалилась
		require.Equal(t, tstIntEl, l.Front().Value) // второй элемент становится первым
	})

	t.Run("head back remove", func(t *testing.T) {
		l := NewList()
		l.PushFront("head")
		l.PushBack("middle")
		l.PushBack("back")

		l.Remove(l.Front())
		l.Remove(l.Back())
		require.Equal(t, "middle", l.Front().Value)
		require.Equal(t, "middle", l.Back().Value)
	})
}

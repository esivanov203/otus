package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 { // если стадий нет
		return in // возвращаем исходный канал (откуда считаются исходные данные)
	}
	current := in
	for _, stage := range stages {
		if stage == nil { // если стадия пустая, то ее пропускаем
			continue
		}
		// между стадиями врезаем горутину,
		// которая получает выходной канал из предыдущей стадии
		// проверяет сигнал done и если нет сигнала
		// передает данные на следующую стадию
		wrapCh := make(Bi)
		go func(cur In) {
			defer func() {
				close(wrapCh)        // безопасно закрываем вспомогательный канал
				for v := range cur { // канал из предыдущей стадии очищаем чтобы избежать deadlock
					_ = v
				}
			}()
			for {
				select {
				case <-done: // проверяем сигнал перед считыванием из канала с предыдущей стадией
					return
				case v, ok := <-cur:
					if !ok {
						return
					}
					select {
					case <-done: // проверяем сигнал перед записью в канал со след стадией
						return
					case wrapCh <- v:
					}
				}
			}
		}(current)
		// связываем стадии через нашу техническую горутину
		// если сигнал done был бы не нужен,
		// то было бы просто next := stage(in) без технической горутины и канала wrapCh
		next := stage(wrapCh)
		current = next
	}

	return current
}

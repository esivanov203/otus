package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// wrap копирует данные из in в новый канал, реагируя на done.
// Если done закрыт — прекращает работу и закрывает out.
func wrap(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		defer func() {
			for v := range in {
				_ = v
			}
		}()

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	current := in
	for _, stage := range stages {
		if stage == nil {
			continue
		}
		// оборачиваем канал перед каждой стадией
		current = stage(wrap(current, done))
	}

	// и выходной канал - тоже
	return wrap(current, done)
}

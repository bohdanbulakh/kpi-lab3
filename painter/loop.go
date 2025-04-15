package painter

import (
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/colornames"
	"image"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(texture screen.Texture)
}

// Loop реалізує цикл подій для формування текстури, отриманої через виконання операцій із внутрішньої черги.
type Loop struct {
	Receiver Receiver

	currentTexture  screen.Texture // Текстура, яка зараз формується
	previousTexture screen.Texture // Текстура, яка була відправлена останнього разу у Receiver

	messageQueue MessageQueue
	state        State
	doneFunc     func()
}

var defaultSize = image.Pt(800, 800)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (eventLoop *Loop) Start(screenDevice screen.Screen) {
	eventLoop.currentTexture, _ = screenDevice.NewTexture(defaultSize)
	eventLoop.messageQueue = MessageQueue{queue: make(chan Operation)}

	eventLoop.state = State{
		backgroundColor: &Fill{Color: colornames.Lime},
	}

	go func() {
		for {
			operation := eventLoop.messageQueue.Pull()

			switch operation.(type) {
			case Figure, Bgrect, Move, Fill, Reset:
				operation.Update(&eventLoop.state)
			case Update:
				eventLoop.state.backgroundColor.Do(eventLoop.currentTexture)

				if eventLoop.state.backgroundRect != nil {
					eventLoop.state.backgroundRect.Do(eventLoop.currentTexture)
				}

				for _, figure := range eventLoop.state.figureCenters {
					figure.Do(eventLoop.currentTexture)
				}

				eventLoop.previousTexture = eventLoop.currentTexture
				eventLoop.Receiver.Update(eventLoop.currentTexture)

				eventLoop.currentTexture, _ = screenDevice.NewTexture(defaultSize)
			}

			if eventLoop.doneFunc != nil {
				eventLoop.doneFunc()
			}
		}
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (eventLoop *Loop) Post(operations OperationList) {
	for _, operation := range operations {
		eventLoop.messageQueue.Push(operation)
	}
}

// MessageQueue — черга повідомлень
type MessageQueue struct {
	queue chan Operation
}

func (messageQueue *MessageQueue) Push(operation Operation) {
	messageQueue.queue <- operation
}

func (messageQueue *MessageQueue) Pull() Operation {
	return <-messageQueue.queue
}

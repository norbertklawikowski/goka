package tester

import (
	"github.com/Shopify/sarama"
	"github.com/lovoo/goka"
)

// emitHandler abstracts a function that allows to overwrite kafkamock's Emit function to
// simulate producer errors
type emitHandler func(topic string, key string, value []byte, options ...EmitOption) *goka.Promise

type producer struct {
	emitter emitHandler
}

func newProducerMock(emitter emitHandler) *producer {
	return &producer{
		emitter: emitter,
	}
}

func (p *producer) AsyncClose() {

}

func (p *producer) Close() error {
	return nil
}
func (p *producer) Input() chan<- *sarama.ProducerMessage {
	return nil
}

func (p *producer) Successes() <-chan *sarama.ProducerMessage {
	return nil
}

func (p *producer) Errors() <-chan *sarama.ProducerError {
	return nil
}

// // Emit emits messages to arbitrary topics.
// // The mock simply forwards the emit to the KafkaMock which takes care of queueing calls
// // to handled topics or putting the emitted messages in the emitted-messages-list
// func (p *producerMock) EmitWithHeaders(topic string, key string, value []byte, headers goka.Headers) *goka.Promise {
// 	return p.emitter(topic, key, value, WithHeaders(headers))
// }

// // Emit emits messages to arbitrary topics.
// // The mock simply forwards the emit to the KafkaMock which takes care of queueing calls
// // to handled topics or putting the emitted messages in the emitted-messages-list
// func (p *producerMock) Emit(topic string, key string, value []byte) *goka.Promise {
// 	return p.emitter(topic, key, value)
// }

// // Close closes the producer mock
// // No action required in the mock.
// func (p *producerMock) Close() error {
// 	logger.Printf("Closing producer mock")
// 	return nil
// }

// flushingProducer wraps the producer and
// waits for all consumers after the Emit.
type flushingProducer struct {
	tester   *Tester
	producer goka.Producer
}

// Emit using the underlying producer
func (e *flushingProducer) EmitWithHeaders(topic string, key string, value []byte, headers goka.Headers) *goka.Promise {
	prom := e.producer.EmitWithHeaders(topic, key, value, headers)
	e.tester.waitForClients()
	return prom
}

// Emit using the underlying producer
func (e *flushingProducer) Emit(topic string, key string, value []byte) *goka.Promise {
	prom := e.producer.Emit(topic, key, value)
	e.tester.waitForClients()
	return prom
}

// Close using the underlying producer
func (e *flushingProducer) Close() error {
	return e.producer.Close()
}

package usecases

import "github.com/csocoteanu/go-stats-collector/stats_writer/domain"

type statsProducer struct {
	out chan domain.HostStats
}

type statsConsumer struct {
	in chan domain.HostStats
}

func NewStatsProducer(ch chan domain.HostStats) *statsProducer {
	p := statsProducer{
		out: ch,
	}

	return &p
}

func (p *statsProducer) Produce(msg domain.HostStats) {
	p.out <- msg
}

func NewStatsConsumer(ch chan domain.HostStats) *statsConsumer {
	c := statsConsumer{
		in: ch,
	}

	return &c
}

func (p *statsConsumer) Consume() (domain.HostStats, bool) {
	stat, ok := <-p.in
	return stat, ok
}

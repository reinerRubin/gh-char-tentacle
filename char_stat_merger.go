package tentacle

// StatMerger TBD
type StatMerger struct {
	incomeStats <-chan CharStat
	acc         CharStat
}

// NewStatMerger TBD
func NewStatMerger(incomeStats <-chan CharStat) *StatMerger {
	return &StatMerger{
		incomeStats: incomeStats,
		acc:         make(CharStat),
	}
}

// Run TBD
func (sm *StatMerger) Run() CharStat {
	for stat := range sm.incomeStats {
		sm.acc.Merge(stat)
	}

	return sm.acc
}

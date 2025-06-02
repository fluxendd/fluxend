package mapper

import (
	statsDto "fluxend/internal/api/dto/stat"
	statsDomain "fluxend/internal/domain/stats"
)

func ToStatResource(stats *statsDomain.Stat) statsDto.Response {
	return statsDto.Response{
		Id:           stats.Id,
		DatabaseName: stats.DatabaseName,
		TotalSize:    stats.TotalSize,
		IndexSize:    stats.IndexSize,
		UnusedIndex:  stats.UnusedIndex,
		TableCount:   stats.TableCount,
		TableSize:    stats.TableSize,
		CreatedAt:    stats.CreatedAt,
	}
}

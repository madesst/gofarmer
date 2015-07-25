package config

import "strconv"

type GlobalConfig struct {
	Version       string `json:"version"`
	DefaultRegion string `json:"default-region"`
	Quotas        Quotas `json:"quotas"`
}

type FarmConfigs map[string]FarmConfig

type FarmConfig struct {
	Name      string     `json:"name"`
	Status    int        `json:"status"`
	CreatedAt int64      `json:"created-at"`
	AMI       string     `json:"ami"`
	Region    string     `json:region`
	Quotas    FarmQuotas `json:"quotas"`
}

type FarmQuotas struct {
	Quotas     Quotas `json:"quotas"`
	FromGlobal bool   `json:"from-global"`
}

type Quotas struct {
	MaxInstances int     `json:"max-instances"`
	MinInstances int     `json:"min-instances"`
	MaxPrice     float64 `json:"max-price"`
	MaxAmount    float64 `json:"max-amout-per-day"`
}

func FloatToString(InputNum float64) string {
	return strconv.FormatFloat(InputNum, 'f', 4, 64)
}

func (fq FarmQuotas) Merge() Quotas {
	q := fq.Quotas
	if fq.FromGlobal {
		gq := GetGlobal().Quotas

		if q.MaxInstances == -1 {
			q.MaxInstances = gq.MaxInstances
		}
		if q.MinInstances == -1 {
			q.MinInstances = gq.MinInstances
		}
		if q.MaxPrice == -1 {
			q.MaxPrice = gq.MaxPrice
		}
		if q.MaxAmount == -1 {
			q.MaxAmount = gq.MaxAmount
		}
	}

	return q
}

func (q Quotas) String() string {
	return "MaxInstances = " + strconv.Itoa(q.MaxInstances) +
		"; MinInstances = " + strconv.Itoa(q.MinInstances) +
		"; MaxPrice = " + FloatToString(q.MaxPrice) +
		"; MaxAmount = " + FloatToString(q.MaxAmount)
}

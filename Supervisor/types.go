package supervisor

import (
	config "URand/Config"
	rng "URand/RNG"
	"context"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/splode/fname"
)

var checkFillmu = sync.Mutex{}

type NameSpace struct {
	Name      string
	Len       uint
	generator *fname.Generator
}

func (ns *NameSpace) ReserveKey() string {
	return ns.Name + "-res"
}
func (ns *NameSpace) ViodKey() string {
	return ns.Name + "-void"
}
func (ns *NameSpace) TmpKey() string {
	return ns.Name + "-tmp"
}
func (ns *NameSpace) FillReserve(c uint) error {
	ctx := context.Background()
	rnd, err := rng.GenerateRandom(c, ns.generator)
	if err != nil {
		return err
	}
	if err := getRedis().Del(ctx, ns.TmpKey()).Err(); err != nil {
		logrus.WithError(err).Errorf("error clearning %s key", ns.TmpKey())
		return err
	}
	if err := getRedis().SAdd(ctx, ns.TmpKey(), *rnd).Err(); err != nil {
		logrus.WithError(err).Errorf("error filling %s key", ns.TmpKey())
		return err
	}
	avail := getRedis().SDiff(ctx, ns.TmpKey(), ns.ViodKey())
	if avail.Err() != nil {
		logrus.WithError(err).Errorf("error diffing %s with %s", ns.TmpKey(), ns.ViodKey())
		return err
	}
	availStr, _ := avail.Result()
	logrus.WithField("total", len(*rnd)).WithField("available", len(availStr)).Debug("unique results")
	getRedis().SAdd(ctx, ns.ReserveKey(), availStr)
	if len(availStr) < len(*rnd) {
		if err = ns.FillReserve(uint(len(*rnd) - len(availStr))); err != nil {
			return err
		}
	}
	return nil
}
func (ns *NameSpace) Get() (string, error) {
	ctx := context.Background()
	val := getRedis().SPop(ctx, ns.ReserveKey())
	err := val.Err()
	if err != nil {
		logrus.WithError(err).Errorf("error getting random name from %s", ns.ReserveKey())
		return "", err
	}
	v, _ := val.Result()
	if err = getRedis().SAdd(ctx, ns.ViodKey(), v).Err(); err != nil {
		logrus.WithError(err).Errorf("error pushing void to %s", ns.ViodKey())
		return v, err
	}
	go ns.checkFill()
	return v, err
}
func (ns *NameSpace) checkFill() error {
	checkFillmu.Lock()
	defer checkFillmu.Unlock()
	ctx := context.Background()
	avail := getRedis().SCard(ctx, ns.ReserveKey())
	err := avail.Err()
	if err != nil {
		logrus.WithError(err).Errorf("error counting available reserves at %s", ns.ReserveKey())
		return err
	}
	cnt, _ := avail.Result()
	if uint(cnt) < config.Config().ReserveMin {
		logrus.WithField("current", cnt).WithField("expected", config.Config().ReserveMin).Info("filling reserve set")
		ns.FillReserve(config.Config().ReserveCount - uint(cnt))
	}
	return nil
}

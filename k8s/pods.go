package k8s

import (
	"errors"
	"github.com/flyhope/kubetea/comm"
	"github.com/sirupsen/logrus"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sync"
	"time"
)

type podCache struct {
	pods      *v12.PodList  // 缓存的数据
	CreatedAt time.Time     // 数据时间
	livetime  time.Duration // 缓存有效期
	timer     *time.Ticker  // 缓存更新定时
	lock      sync.Mutex    // 缓存更新锁，保证同时只有一个更新
	dataRead  chan bool     // 数据读取标识，如果没有读取，则缓存不会主动再取
}

// Tick 定时拉取更新数据
func (p *podCache) Tick() {
	for {
		<-p.timer.C
		_, err := p.ShowListUpdate()
		if err != nil {
			logrus.Warnln(err)
		}
		comm.Program.Send(comm.MsgPodCache(p.CreatedAt))
		p.dataRead <- true
	}
}

// 获取缓存中的pod数据
func (p *podCache) showCachePods() *v12.PodList {
	select {
	case <-p.dataRead:
	default:
	}
	return p.pods
}

// ShowList 获取Pod的列表
func (p *podCache) ShowList(forceUpdate bool) (*v12.PodList, error) {
	// 缓存没有或失效，更新缓存
	if forceUpdate || p.needUpdate() {
		return p.ShowListUpdate()
	}
	return p.showCachePods(), nil
}

// ShowListUpdate 强制获取列表数据并更新缓存
func (p *podCache) ShowListUpdate() (*v12.PodList, error) {
	var err error
	if p.lock.TryLock() {
		defer p.lock.Unlock()
		p.pods, err = Client().CoreV1().Pods(ShowNamespace()).List(comm.Context.Context, v1.ListOptions{
			Limit: 10000,
		})
		p.CreatedAt = time.Now()

		// 启用定时获取数据
		if p.timer == nil {
			p.timer = time.NewTicker(p.livetime)
			go p.Tick()
		}

	}
	return p.showCachePods(), err
}

// Show 获取一个Pod数据
func (p *podCache) Show(name string, force bool) (*v12.Pod, time.Time, error) {
	if !force {
		// 缓存没有或失效，更新缓存
		if p.needUpdate() {
			if _, err := p.ShowListUpdate(); err != nil {
				return nil, p.CreatedAt, err
			}
		}

		pods := p.showCachePods()
		if pods == nil {
			return nil, time.Now(), errors.New("empty pods")
		}

		for _, pod := range pods.Items {
			if pod.Name == name {
				return &pod, p.CreatedAt, nil
			}
		}
	}

	pod, err := PodGet(name)
	return pod, time.Now(), err
}

// 判断是否需要更新数据
func (p *podCache) needUpdate() bool {
	return p.pods == nil || time.Now().After(p.CreatedAt.Add(p.livetime))
}

var PodCache = sync.OnceValue(func() *podCache {
	c := &podCache{
		livetime: time.Second * time.Duration(comm.ShowKConfig().PodCacheLivetime),
		dataRead: make(chan bool),
	}
	return c
})

// PodGet 获取一个POD的配置
func PodGet(name string) (*v12.Pod, error) {
	return Client().CoreV1().Pods(ShowNamespace()).Get(comm.Context.Context, name, v1.GetOptions{})
}

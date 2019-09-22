package producer

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/you/hello/pkg/entry"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
)

type status int32

const (
	statusInit status = iota
	statusStarted
	statusPendingStop
	statusStopped
)

const maxScanEntries = 10
const maxBufferedResults = 10

type ImageProducer struct {
	statusLock  sync.Mutex
	dirPath     string
	status      status
	entriesChan chan *entry.Entry
	errorsChan  chan error
	doneChan    chan struct{}
}

func NewProducer(dirPath string) *ImageProducer {
	return &ImageProducer{
		dirPath:     dirPath,
		status:      statusInit,
		doneChan:    make(chan struct{}),
		errorsChan:  make(chan error, maxBufferedResults),
		entriesChan: make(chan *entry.Entry, maxBufferedResults),
	}
}

func (p *ImageProducer) Start() error {
	p.statusLock.Lock()
	defer p.statusLock.Unlock()
	currStatus := p.loadStatus()
	if currStatus != statusInit {
		return fmt.Errorf("illegal status: %v", p.status)
	}
	p.storeStatus(statusStarted)
	go p.scanDirectory()
	return nil
}

func (p *ImageProducer) loadStatus() status {
	return status(atomic.LoadInt32((*int32)(&p.status)))
}

func (p *ImageProducer) storeStatus(newStatus status) {
	atomic.StoreInt32((*int32)(&p.status), int32(newStatus))
}

func (p *ImageProducer) scanDirectory() {
	defer func() {
		p.storeStatus(statusStopped)
		close(p.doneChan)
	}()
	srcDir, err := os.Open(p.dirPath)
	if err != nil {
		p.errorsChan <- errors.Wrapf(err, "failed to open directory $q", p.dirPath)
		return
	}
	for p.loadStatus() == statusStarted {
		fileNames, err := srcDir.Readdirnames(maxScanEntries)
		if err != nil {
			if err != io.EOF {
				p.errorsChan <- errors.Wrapf(err, "failed to scan the directory %q", p.dirPath)
			}
			break
		}
		for _, fileName := range fileNames {
			if p.loadStatus() != statusStarted {
				break
			}
			filePath := path.Join(srcDir.Name(), fileName)
			srcImg, err := imaging.Open(filePath)
			if err != nil {
				p.errorsChan <- errors.Wrapf(err, "failed to open %q as image", filePath)
				continue
			}
			p.entriesChan <- &entry.Entry{
				Name:  strings.TrimSuffix(fileName, path.Ext(fileName)),
				Image: srcImg,
			}
		}
	}
}

func (p *ImageProducer) Stop() {
	if p.setStatusPendingStop() == statusPendingStop {
		log.Print("waiting background job to stop...")
		<- p.doneChan
	}
}

func (p *ImageProducer) setStatusPendingStop() status {
	p.statusLock.Lock()
	defer p.statusLock.Unlock()
	currStatus := p.loadStatus()
	newStatus := currStatus
	if currStatus == statusStarted {
		newStatus = statusPendingStop
		p.storeStatus(newStatus)
	}
	return newStatus
}

func (p *ImageProducer) Done() <- chan struct{} {
	return p.doneChan
}

func (p *ImageProducer) Entries() <-chan *entry.Entry {
	return p.entriesChan
}

func (p *ImageProducer) Errors() <-chan error {
	return p.errorsChan
}
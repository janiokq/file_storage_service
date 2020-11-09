package manager

import (
	"filesrv/conf"
	"filesrv/library/log"
	"filesrv/repositoty"
	"go.uber.org/zap"
	"sync"
)

var m *FileManager

type FileManager struct {
	fileItems *sync.Map
	clearItem chan int64
	r         *repositoty.Repository
}

type FileUploadItem struct {
	Fid  int64
	Part int32
	Data []byte
	Md5  string
}

func NewFileManager(r *repositoty.Repository) {
	m = &FileManager{
		fileItems: new(sync.Map),
		clearItem: make(chan int64, 10),
		r:         r,
	}
	go m.run()
	return
}

func GetFileManager() *FileManager {
	return m
}

func (f *FileManager) send(fid int64) {
	if f == nil {
		return
	}
	f.clearItem <- fid
}

func (f *FileManager) NewItem(item *FileItem) {
	item = NewFileItem(item)
	f.fileItems.Store(item.Fid, item)
}

func (f *FileManager) AddItem(upItem *FileUploadItem) error {
	item, ok := f.fileItems.Load(upItem.Fid)
	if !ok {
		return conf.ErrFileIdInvalid
	}
	fItem := item.(*FileItem)
	return fItem.AddItem(upItem)
}

func (f *FileManager) DelItem(fid int64) {
	f.fileItems.Delete(fid)
}

func (f *FileManager) run() {
	for {
		select {
		case fid := <-f.clearItem:
			log.GetLogger().Debug("auto clear file cache", zap.Any("fid", fid))
			f.fileItems.Delete(fid)
		}
	}
}

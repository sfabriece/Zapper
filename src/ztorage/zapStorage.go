package ztorage

import(
	"data"
	"sync"
)

var lock sync.Mutex
type ZapList []data.Zap 
type ViewersList struct {
	M map[string]int
}

func NewZapStore() *ViewersList {
	//zs := make(ZapList, 0)
	M := ViewersList{make(map[string]int)}
	return &M
}

func (vl *ViewersList)StoreZap(z data.Zap){
	lock.Lock()
//	*vl = append(*zs, z)
	vl.M[z.ToCH] ++
	vl.M[z.FromCH] --
	lock.Unlock()
}

func (vl *ViewersList) ComputeViewers(chName string) int {
	lock.Lock()
	viewers := vl.M[chName]
	lock.Unlock()
	//if viewers > 0 {return viewers}
	return viewers
}




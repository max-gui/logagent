package confload

import (
	"context"
	"os"
	"sync"

	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/logagent/pkg/logsets"
	"gopkg.in/yaml.v2"
)

func Load(c context.Context) []byte {

	log := logagent.InstArch(c)
	bytes, err := os.ReadFile(*logsets.Apppath + string(os.PathSeparator) + "application-" + *logsets.DCENV + ".yml")
	if err != nil {
		log.Panic(err)
	}

	return bytes
}

type mutexKV struct {
	sync.RWMutex
	kvs map[string]interface{}
}

var kvmap = mutexKV{kvs: make(map[string]interface{})}

func (v *mutexKV) help(tricky func(map[string]interface{}) (bool, interface{})) (bool, interface{}) {
	v.Lock()
	ok, res := tricky(v.kvs)
	v.Unlock()
	return ok, res
}

func LoadEnv(env string, c context.Context) interface{} {
	if ok, value := kvmap.help(func(kvs map[string]interface{}) (bool, interface{}) {
		if val, ok := kvs[env]; ok {
			return ok, val
		} else {
			return ok, nil
		}
	}); ok {
		// inst = value
		return value
	}

	log := logagent.InstArch(c)
	bytes, err := os.ReadFile(*logsets.Apppath + string(os.PathSeparator) + "application-" + env + ".yml")
	if err != nil {
		log.Panic(err)
	}
	inst := map[string]interface{}{}
	err = yaml.Unmarshal(bytes, &inst)
	if err != nil {
		log.Panic(err)
	}

	kvmap.help(func(kvs map[string]interface{}) (bool, interface{}) {
		kvs[env] = inst
		return true, nil
	})
	return inst
}

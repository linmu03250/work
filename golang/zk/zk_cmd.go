package server

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	ZKPUBPATH   = "/metis_pubs/%s"
	ZKPUBS_NODE = "/metis_pubs"

	ZKFILTERS_PATH = "/metis_filters/%s"
	ZKFILTERS_NODE = "/metis_filters"

	ZKFILTER_PROSECCERS      = "/metis_processers"
	ZKFILTER_PROSECCERS_PATH = "/metis_processers/%s/%s"

	ZKELECTIONPATH = "/metis/election/members"
	ZKHTTP_OFFSETS = "/metis_http_offsets/%s"

	VERSION = -1
)

// 向zk注册processer
// fname： 过滤规则名
// id: 每个processer goroutine 内存维护的唯一数值
func ZKRegisterProcesser(c *zk.Conn, fname string, id int64) error {
	zkPath := fmt.Sprintf(ZKFILTER_PROSECCERS_PATH, fname, strconv.FormatInt(id, 10))
	//create Ephemeral znode
	_, err := c.Create(zkPath, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil && err != zk.ErrNoNode {
		return err
	}
	if err == zk.ErrNoNode {
		mkdirRecursive(c, path.Dir(zkPath))
	}
	return nil
}

func ZKunRegisterProcesser(c *zk.Conn, fname string, id int64) error {
	zkPath := fmt.Sprintf(ZKFILTER_PROSECCERS_PATH, fname, strconv.FormatInt(id, 10))
	err := c.Delete(zkPath, VERSION)
	if err != nil && err != zk.ErrNoNode {
		return err
	}
	return nil
}

// 获取某个filter的processer 数量
func ZKGetProcesserNum(c *zk.Conn, fname string) (int, error) {
	zkPath := ZKFILTER_PROSECCERS + "/" + fname
	isExist, stat, err := c.Exists(zkPath)
	if err != nil {
		return 0, err
	}
	if isExist {
		return int(stat.NumChildren), nil
	}
	return 0, nil
}

func ZKaddFilter(c *zk.Conn, fname string, data []byte) error {
	zkPath := fmt.Sprintf(ZKFILTERS_PATH, fname)
	_, err := c.Set(zkPath, data, VERSION)
	if err != nil && err != zk.ErrNoNode {
		return err
	}
	if err == zk.ErrNoNode {
		mkdirRecursive(c, path.Dir(zkPath))
		_, err = c.Create(zkPath, data, 0, zk.WorldACL(zk.PermAll))
	}
	return err
}

func ZKgetFilter(c *zk.Conn, fname string) (data []byte, err error, ch <-chan zk.Event) {
	zkPath := fmt.Sprintf(ZKFILTERS_PATH, fname)
	data, _, ch, err = c.GetW(zkPath)
	if err != nil {
		return nil, err, ch
	}
	return data, nil, ch
}

func ZKaddOffsetByAccessKey(c *zk.Conn, accesskey string, data []byte) error {
	zkPath := fmt.Sprintf(ZKHTTP_OFFSETS, accesskey)
	_, err := c.Set(zkPath, data, VERSION)
	if err != nil && err != zk.ErrNoNode {
		return err
	}
	if err == zk.ErrNoNode {
		mkdirRecursive(c, path.Dir(zkPath))
		_, err = c.Create(zkPath, data, 0, zk.WorldACL(zk.PermAll))
	}
	return err
}

func ZKgetOffsetByAccessKey(c *zk.Conn, accesskey string) (data []byte, err error, ch <-chan zk.Event) {
	zkPath := fmt.Sprintf(ZKHTTP_OFFSETS, accesskey)
	data, _, ch, err = c.GetW(zkPath)
	if err != nil {
		return nil, err, ch
	}
	return data, nil, ch
}

func ZKexistsFilter(c *zk.Conn, fname string) (bool, error) {
	zkPath := fmt.Sprintf(ZKFILTERS_PATH, fname)
	exists, _, err := c.Exists(zkPath)
	if err != nil {
		return true, err
	}
	if exists {
		return true, nil
	}

	return false, nil
}

func ZKdelFilter(c *zk.Conn, fname string) error {
	zkPath := fmt.Sprintf(ZKFILTERS_PATH, fname)
	err := c.Delete(zkPath, VERSION)
	//if err != nil && err != zk.ErrNoNode {
	if err != nil {
		return err
	}
	return nil
}

func ZKgetallfilters(c *zk.Conn) ([][]byte, error) {
	data, _, err := c.Children(ZKFILTERS_NODE)
	if err != nil {
		return nil, err
	}
	var filters [][]byte
	var filter []byte

	for _, fname := range data {
		filter, _, err = c.Get(fmt.Sprintf(ZKFILTERS_PATH, fname))
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}
	return filters, nil
}

func ZKgetallpubs(c *zk.Conn) ([][]byte, error) {
	data, _, err := c.Children(ZKPUBS_NODE)
	if err != nil {
		return nil, err
	}
	var pubs [][]byte
	var pub []byte

	for _, pubname := range data {
		pub, _, err = c.Get(fmt.Sprintf(ZKPUBPATH, pubname))
		if err != nil {
			return nil, err
		}
		pubs = append(pubs, pub)
	}
	return pubs, nil
}

func ZKgetPub(c *zk.Conn, pubname string) (data []byte, err error, ch <-chan zk.Event) {
	zkPath := fmt.Sprintf(ZKPUBPATH, pubname)
	data, _, ch, err = c.GetW(zkPath)
	if err != nil {
		return nil, err, ch
	}
	return data, nil, ch
}

func ZKexistsPub(c *zk.Conn, pubname string) (bool, error) {
	zkPath := fmt.Sprintf(ZKPUBPATH, pubname)
	isExist, _, err := c.Exists(zkPath)
	if err != nil {
		return true, err
	}
	if isExist {
		return true, nil
	}

	return false, nil
}

func ZKaddPub(c *zk.Conn, pubname string, data []byte) error {
	zkPath := fmt.Sprintf(ZKPUBPATH, pubname)
	_, err := c.Set(zkPath, data, VERSION)
	if err != nil && err != zk.ErrNoNode {
		return err
	}
	if err == zk.ErrNoNode {
		mkdirRecursive(c, path.Dir(zkPath))
		_, err = c.Create(zkPath, data, 0, zk.WorldACL(zk.PermAll))
	}
	return err
}

func ZKdelPub(c *zk.Conn, pubname string) error {
	zkPath := fmt.Sprintf(ZKPUBPATH, pubname)
	err := c.Delete(zkPath, VERSION)
	//if err != nil && err != zk.ErrNoNode {
	if err != nil {
		return err
	}
	return nil
}

func ZKChildrenW(c *zk.Conn, path string) (childs []string, ch <-chan zk.Event, err error) {
	childs, _, ch, err = c.ChildrenW(path)
	if err != nil {
		return nil, nil, err
	}
	return childs, ch, nil
}

func ZKCreateElection(c *zk.Conn) (string, error) {
	r, err := c.Create(ZKELECTIONPATH, nil, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil && err != zk.ErrNoNode {
		return "", err
	}
	if err == zk.ErrNoNode {
		mkdirRecursive(c, path.Dir(ZKELECTIONPATH))
		return ZKCreateElection(c)
	}
	return r, nil
}

func ZKElectionWatch(c *zk.Conn) (<-chan zk.Event, error) {
	_, _, ch, err := c.ChildrenW(path.Dir(ZKELECTIONPATH))
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func IsLeader(p string, c *zk.Conn) (bool, error) {
	if c.State() == zk.StateDisconnected {
		return false, zk.ErrConnectionClosed
	}
	fmt.Println(p, "Children")
	childs, _, err := c.Children(path.Dir(ZKELECTIONPATH))
	// fmt.Println(p, childs, err)
	if err != nil {
		return false, err
	}
	basePath := path.Base(p)
	isLeader := true
	for _, child := range childs {
		if strings.Compare(basePath, child) > 0 {
			isLeader = false
			break
		}
	}
	return isLeader, nil
}

func mkdirRecursive(c *zk.Conn, zkPath string) error {
	var err error
	parent := path.Dir(zkPath)
	if parent != "/" {
		if err = mkdirRecursive(c, parent); err != nil {
			return err
		}
	}

	_, err = c.Create(zkPath, nil, 0, zk.WorldACL(zk.PermAll))
	if err == zk.ErrNodeExists {
		err = nil
	}
	return err
}


package merkledag

import "hash"

const (
	K          = 1 << 10
	BLOCK_SIZE = 256 * K
)

type Link struct {
	Name string
	Hash []byte
	Size int
}

type Object struct {
	Links []Link
	Data  []byte
}

func Add(store KVStore, node Node, h hash.Hash) []byte {
	// 将分片写入到KVStore中，并返回Merkle Root
	switch n := node.(type) {
	case File:
		return StoreFile(store, n, h)
	case Dir:
		return StoreDir(store, n, h)
	}
	return nil
}

func StoreFile(store KVStore, node File, h hash.Hash) []byte {
	// 处理文件存储
	t := []byte("blob")
	if node.Size() > BLOCK_SIZE {
		t = []byte("list")
	}

	data := node.Bytes()

	// 将数据存储到KVStore中
	err := store.Put(h.Sum(nil), data)
	if err != nil {
		// 处理错误
		return nil
	}

	return h.Sum(nil)
}

func StoreDir(store KVStore, dir Dir, h hash.Hash) []byte {
	// 处理文件夹存储
	tree := Object{
		Links: make([]Link, 0),
		Data:  make([]byte, 0),
	}

	It := dir.It()
	for It.Next() {
		node := It.Node()

		// 递归处理子文件/文件夹
		linkHash := Add(store, node, h)
		if linkHash != nil {
			// 创建并添加链接到对象中
			link := Link{
				Name: node.Name(), // 假设Node接口有Name()方法来获取名称
				Hash: linkHash,
				Size: node.Size(),
			}
			tree.Links = append(tree.Links, link)
		}
	}

	// 计算对象的哈希并将其存储到KVStore中
	// 这里省略了计算哈希和存储的部分

	return nil
}

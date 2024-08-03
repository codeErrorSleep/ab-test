package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash/fnv"
)

// HashFunc is a function type that takes a string and returns a uint32 hash value
type HashFunc func(data string) uint32

// DefaultHashFunc is the default hash function using FNV-1a
func DefaultHashFunc(data string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(data))
	return h.Sum32()
}

// MD5HashFunc is a hash function using MD5
func MD5HashFunc(data string) uint32 {
	h := md5.New()
	h.Write([]byte(data))
	sum := h.Sum(nil)
	// 截取MD5散列的前4个字节来获取一个uint32值
	return uint32(sum[0]) | uint32(sum[1])<<8 | uint32(sum[2])<<16 | uint32(sum[3])<<24
}

// SHA1HashFunc is a hash function using SHA-1
func SHA1HashFunc(data string) uint32 {
	h := sha1.New()
	h.Write([]byte(data))
	sum := h.Sum(nil)
	// 截取SHA-1散列的前4个字节来获取一个uint32值
	return uint32(sum[0]) | uint32(sum[1])<<8 | uint32(sum[2])<<16 | uint32(sum[3])<<24
}

// SHA256HashFunc is a hash function using SHA-256
func SHA256HashFunc(data string) uint32 {
	h := sha256.New()
	h.Write([]byte(data))
	sum := h.Sum(nil)
	// 截取SHA-256散列的前4个字节来获取一个uint32值
	return uint32(sum[0]) | uint32(sum[1])<<8 | uint32(sum[2])<<16 | uint32(sum[3])<<24
}

package ceph

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

var cephConn *s3.S3

func GetCephConn() *s3.S3 {
	if cephConn != nil{
		return cephConn
	}
	auth := aws.Auth{
		AccessKey: "",
		SecretKey: "",
	}
	curRegion := aws.Region{
		Name:                 "default",
		EC2Endpoint:          "http://192.168.0.10:9080",
		S3Endpoint:           "http://192.168.0.10:9080",
		S3BucketEndpoint:     "",
		S3LocationConstraint: false,
		S3LowercaseBucket:    false,
		Sign:                 aws.SignV2,
	}
	return s3.New(auth, curRegion)
}

func GetCephBucket(bucket string) *s3.Bucket {
	conn := GetCephConn()
	return conn.Bucket(bucket)
}

// PutObject : 上传文件到ceph集群
func PutObject(bucket string, path string, data []byte) error {
	return GetCephBucket(bucket).Put(path, data, "octet-stream", s3.PublicRead)
}

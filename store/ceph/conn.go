package ceph

import (
	"deck/config"
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

var cephConn *s3.S3

func GetCephConnection() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}
	auth := aws.Auth{
		AccessKey: config.CephAccessKey,
		SecretKey: config.CephSecretKey,
	}
	curRegion := aws.Region{
		Name:                 "default",
		EC2Endpoint:          config.CephGWEndpoint,
		S3Endpoint:           config.CephGWEndpoint,
		S3BucketEndpoint:     "",
		S3LocationConstraint: false,
		S3LowercaseBucket:    false,
		Sign:                 aws.SignV2,
	}
	return s3.New(auth, curRegion)
}
func GetCephBucket(bucket string) *s3.Bucket {
	conn := GetCephConnection()
	return conn.Bucket(bucket)
}
func PutObject(bucket string, path string, data []byte) error {
	return GetCephBucket(bucket).Put(path, data, "octet-stream", s3.PublicRead)
}

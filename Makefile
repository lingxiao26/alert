.PHONY: build

ALL: build

build:
	docker build . -t lbxdrugs.tencentcloudcr.com/xls/member-alert:v4


package main

import (
	"fmt"
	"os"
	"strings"

	//"github.com/kr/pretty"
	"github.com/flant/docker-registry-client/registry"
)

var registryUrl = "https://registry.flant.com"
var username = "oauth2"
var password = "PASTE_YOUR_GITLAB_PRIVATE_TOKEN"
var image = "registry.flant.com/sys/antiopa:play_registry"

func main() {

	hub, err := registry.New(registryUrl, username, password)

	if err != nil {
		fmt.Printf("registry New err: %v", err)
		os.Exit(1)
	}

	imageId, err := GetImageId(hub, image)
	if err != nil {
		fmt.Printf("GetImageId err: %v", err)
		os.Exit(1)
	}

	fmt.Printf("ImageId: %s\n", imageId)

}

func GetImageId(registry *registry.Registry, image string) (string, error) {
	imageName, tag := ParseAntiopaImageName(image)

	antiopaManifest, err := registry.ManifestV2(imageName, tag)

	if err != nil {
		return "", err
	}

	imageId := antiopaManifest.Config.Digest.String()

	return imageId, nil
}

func ParseAntiopaImageName(image string) (repository string, tag string) {
	if strings.Contains(image, ":") {
		res := strings.SplitN(image, ":", 2)
		return res[0], res[1]
	} else {
		return image, "stable"
	}
}

// листинг доступных репозиториев (имён образов без тэгов) в gitlab-е пока не сделан.
//repositories, err := hub.Repositories()

//pretty.Printf("%# v\n", antiopaManifest)
//
//for _, layer := range antiopaManifest.Layers {
//	pretty.Printf("%# v %# v %# v\n", layer.MediaType, layer.Size, layer.Digest)
//}

// Примеры вывода:
/*

pretty.Printf("%# v\n", antiopaManifest)

&schema2.DeserializedManifest{
    Manifest: schema2.Manifest{
        Versioned: manifest.Versioned{SchemaVersion:2, MediaType:"application/vnd.docker.distribution.manifest.v2+json"},
        Config:    distribution.Descriptor{
            MediaType: "application/vnd.docker.container.image.v1+json",
            Size:      8235,
            Digest:    "sha256:242f56fa1511bd7b14843e2ca5b52bfca60ee4b7b6e11ac80cca54b515c336b5",
            URLs:      nil,
        },
        Layers: {
            {
                MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
                Size:      47536248,
                Digest:    "sha256:f5c64a3438f6e850c2d09c6cdd407dff522129d5304fdd3a9dbb397c9488e7aa",
                URLs:      nil,
            },
            {
                MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
                Size:      848,
                Digest:    "sha256:51899d335aae78435d7960bc74ca06c3a1db706bcb13c2dc19e2fcf567153f97",
                URLs:      nil,
            },
            {
                MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
                Size:      621,
                Digest:    "sha256:6ff2b7de3c1376e8f6f0d823d7c15d5a32df338d7bc9013841dff4cefc2ba2ca",
                URLs:      nil,
            },
            {
                MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
                Size:      853,
                Digest:    "sha256:50366c31f7fd3f86c9718ad3724dcd6926ba3f69f25f4bed47be5ca9640a0616",
                URLs:      nil,
            },
            {
                MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
                Size:      169,
                Digest:    "sha256:f441b9a68d131e4aecdfeca0d35f7c984b6a66582f0d4b9e952ffb5c45c196aa",
                URLs:      nil,
            },
            {
                MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
                Size:      147,
                Digest:    "sha256:446a15856c11ef8ec3c9ee74591277a7f155c771308b1692addc84ba55d7a79f",
                URLs:      nil,
            },
            {
                MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
                Size:      196742897,
                Digest:    "sha256:f053c12ab4844d72e2ec90f38dab888d58826ef2a8ace62733f2dfb7052a46d1",
                URLs:      nil,
            },
            {
                MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
                Size:      8604496,
                Digest:    "sha256:e51c1ab60bc39a0a1fba2a78cbfa9aa684b91f1ac56a896765e4dbf6eb0e42cf",
                URLs:      nil,
            },
            {
                MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
                Size:      820,
                Digest:    "sha256:3a19e9c6c66e940d2496c1e312df85b67e9a7e2f63dc25a102f1ae58562ae442",
                URLs:      nil,
            },
            {
                MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
                Size:      678,
                Digest:    "sha256:9a933bf37691be08c526f47a0062c566111144dd12b40e735d0b2739fa4c7868",
                URLs:      nil,
            },
        },
    },
    canonical: {0x7b, 0xa, 0x20, 0x20, 0x20, 0x22, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x20, 0x32, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x6e, 0x64, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x64, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x6d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x32, 0x2b, 0x6a, 0x73, 0x6f, 0x6e, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x22, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x3a, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x6e, 0x64, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2b, 0x6a, 0x73, 0x6f, 0x6e, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x3a, 0x20, 0x38, 0x32, 0x33, 0x35, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x20, 0x22, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x3a, 0x32, 0x34, 0x32, 0x66, 0x35, 0x36, 0x66, 0x61, 0x31, 0x35, 0x31, 0x31, 0x62, 0x64, 0x37, 0x62, 0x31, 0x34, 0x38, 0x34, 0x33, 0x65, 0x32, 0x63, 0x61, 0x35, 0x62, 0x35, 0x32, 0x62, 0x66, 0x63, 0x61, 0x36, 0x30, 0x65, 0x65, 0x34, 0x62, 0x37, 0x62, 0x36, 0x65, 0x31, 0x31, 0x61, 0x63, 0x38, 0x30, 0x63, 0x63, 0x61, 0x35, 0x34, 0x62, 0x35, 0x31, 0x35, 0x63, 0x33, 0x33, 0x36, 0x62, 0x35, 0x22, 0xa, 0x20, 0x20, 0x20, 0x7d, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x22, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x22, 0x3a, 0x20, 0x5b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x6e, 0x64, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x72, 0x6f, 0x6f, 0x74, 0x66, 0x73, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x74, 0x61, 0x72, 0x2e, 0x67, 0x7a, 0x69, 0x70, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x3a, 0x20, 0x34, 0x37, 0x35, 0x33, 0x36, 0x32, 0x34, 0x38, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x20, 0x22, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x3a, 0x66, 0x35, 0x63, 0x36, 0x34, 0x61, 0x33, 0x34, 0x33, 0x38, 0x66, 0x36, 0x65, 0x38, 0x35, 0x30, 0x63, 0x32, 0x64, 0x30, 0x39, 0x63, 0x36, 0x63, 0x64, 0x64, 0x34, 0x30, 0x37, 0x64, 0x66, 0x66, 0x35, 0x32, 0x32, 0x31, 0x32, 0x39, 0x64, 0x35, 0x33, 0x30, 0x34, 0x66, 0x64, 0x64, 0x33, 0x61, 0x39, 0x64, 0x62, 0x62, 0x33, 0x39, 0x37, 0x63, 0x39, 0x34, 0x38, 0x38, 0x65, 0x37, 0x61, 0x61, 0x22, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7d, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x6e, 0x64, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x72, 0x6f, 0x6f, 0x74, 0x66, 0x73, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x74, 0x61, 0x72, 0x2e, 0x67, 0x7a, 0x69, 0x70, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x3a, 0x20, 0x38, 0x34, 0x38, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x20, 0x22, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x3a, 0x35, 0x31, 0x38, 0x39, 0x39, 0x64, 0x33, 0x33, 0x35, 0x61, 0x61, 0x65, 0x37, 0x38, 0x34, 0x33, 0x35, 0x64, 0x37, 0x39, 0x36, 0x30, 0x62, 0x63, 0x37, 0x34, 0x63, 0x61, 0x30, 0x36, 0x63, 0x33, 0x61, 0x31, 0x64, 0x62, 0x37, 0x30, 0x36, 0x62, 0x63, 0x62, 0x31, 0x33, 0x63, 0x32, 0x64, 0x63, 0x31, 0x39, 0x65, 0x32, 0x66, 0x63, 0x66, 0x35, 0x36, 0x37, 0x31, 0x35, 0x33, 0x66, 0x39, 0x37, 0x22, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7d, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x6e, 0x64, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x72, 0x6f, 0x6f, 0x74, 0x66, 0x73, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x74, 0x61, 0x72, 0x2e, 0x67, 0x7a, 0x69, 0x70, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x3a, 0x20, 0x36, 0x32, 0x31, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x20, 0x22, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x3a, 0x36, 0x66, 0x66, 0x32, 0x62, 0x37, 0x64, 0x65, 0x33, 0x63, 0x31, 0x33, 0x37, 0x36, 0x65, 0x38, 0x66, 0x36, 0x66, 0x30, 0x64, 0x38, 0x32, 0x33, 0x64, 0x37, 0x63, 0x31, 0x35, 0x64, 0x35, 0x61, 0x33, 0x32, 0x64, 0x66, 0x33, 0x33, 0x38, 0x64, 0x37, 0x62, 0x63, 0x39, 0x30, 0x31, 0x33, 0x38, 0x34, 0x31, 0x64, 0x66, 0x66, 0x34, 0x63, 0x65, 0x66, 0x63, 0x32, 0x62, 0x61, 0x32, 0x63, 0x61, 0x22, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7d, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x6e, 0x64, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x72, 0x6f, 0x6f, 0x74, 0x66, 0x73, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x74, 0x61, 0x72, 0x2e, 0x67, 0x7a, 0x69, 0x70, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x3a, 0x20, 0x38, 0x35, 0x33, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x20, 0x22, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x3a, 0x35, 0x30, 0x33, 0x36, 0x36, 0x63, 0x33, 0x31, 0x66, 0x37, 0x66, 0x64, 0x33, 0x66, 0x38, 0x36, 0x63, 0x39, 0x37, 0x31, 0x38, 0x61, 0x64, 0x33, 0x37, 0x32, 0x34, 0x64, 0x63, 0x64, 0x36, 0x39, 0x32, 0x36, 0x62, 0x61, 0x33, 0x66, 0x36, 0x39, 0x66, 0x32, 0x35, 0x66, 0x34, 0x62, 0x65, 0x64, 0x34, 0x37, 0x62, 0x65, 0x35, 0x63, 0x61, 0x39, 0x36, 0x34, 0x30, 0x61, 0x30, 0x36, 0x31, 0x36, 0x22, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7d, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x6e, 0x64, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x72, 0x6f, 0x6f, 0x74, 0x66, 0x73, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x74, 0x61, 0x72, 0x2e, 0x67, 0x7a, 0x69, 0x70, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x3a, 0x20, 0x31, 0x36, 0x39, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x20, 0x22, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x3a, 0x66, 0x34, 0x34, 0x31, 0x62, 0x39, 0x61, 0x36, 0x38, 0x64, 0x31, 0x33, 0x31, 0x65, 0x34, 0x61, 0x65, 0x63, 0x64, 0x66, 0x65, 0x63, 0x61, 0x30, 0x64, 0x33, 0x35, 0x66, 0x37, 0x63, 0x39, 0x38, 0x34, 0x62, 0x36, 0x61, 0x36, 0x36, 0x35, 0x38, 0x32, 0x66, 0x30, 0x64, 0x34, 0x62, 0x39, 0x65, 0x39, 0x35, 0x32, 0x66, 0x66, 0x62, 0x35, 0x63, 0x34, 0x35, 0x63, 0x31, 0x39, 0x36, 0x61, 0x61, 0x22, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7d, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x6e, 0x64, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x72, 0x6f, 0x6f, 0x74, 0x66, 0x73, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x74, 0x61, 0x72, 0x2e, 0x67, 0x7a, 0x69, 0x70, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x3a, 0x20, 0x31, 0x34, 0x37, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x20, 0x22, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x3a, 0x34, 0x34, 0x36, 0x61, 0x31, 0x35, 0x38, 0x35, 0x36, 0x63, 0x31, 0x31, 0x65, 0x66, 0x38, 0x65, 0x63, 0x33, 0x63, 0x39, 0x65, 0x65, 0x37, 0x34, 0x35, 0x39, 0x31, 0x32, 0x37, 0x37, 0x61, 0x37, 0x66, 0x31, 0x35, 0x35, 0x63, 0x37, 0x37, 0x31, 0x33, 0x30, 0x38, 0x62, 0x31, 0x36, 0x39, 0x32, 0x61, 0x64, 0x64, 0x63, 0x38, 0x34, 0x62, 0x61, 0x35, 0x35, 0x64, 0x37, 0x61, 0x37, 0x39, 0x66, 0x22, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7d, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x6e, 0x64, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x72, 0x6f, 0x6f, 0x74, 0x66, 0x73, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x74, 0x61, 0x72, 0x2e, 0x67, 0x7a, 0x69, 0x70, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x3a, 0x20, 0x31, 0x39, 0x36, 0x37, 0x34, 0x32, 0x38, 0x39, 0x37, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x20, 0x22, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x3a, 0x66, 0x30, 0x35, 0x33, 0x63, 0x31, 0x32, 0x61, 0x62, 0x34, 0x38, 0x34, 0x34, 0x64, 0x37, 0x32, 0x65, 0x32, 0x65, 0x63, 0x39, 0x30, 0x66, 0x33, 0x38, 0x64, 0x61, 0x62, 0x38, 0x38, 0x38, 0x64, 0x35, 0x38, 0x38, 0x32, 0x36, 0x65, 0x66, 0x32, 0x61, 0x38, 0x61, 0x63, 0x65, 0x36, 0x32, 0x37, 0x33, 0x33, 0x66, 0x32, 0x64, 0x66, 0x62, 0x37, 0x30, 0x35, 0x32, 0x61, 0x34, 0x36, 0x64, 0x31, 0x22, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7d, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x6e, 0x64, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x72, 0x6f, 0x6f, 0x74, 0x66, 0x73, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x74, 0x61, 0x72, 0x2e, 0x67, 0x7a, 0x69, 0x70, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x3a, 0x20, 0x38, 0x36, 0x30, 0x34, 0x34, 0x39, 0x36, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x20, 0x22, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x3a, 0x65, 0x35, 0x31, 0x63, 0x31, 0x61, 0x62, 0x36, 0x30, 0x62, 0x63, 0x33, 0x39, 0x61, 0x30, 0x61, 0x31, 0x66, 0x62, 0x61, 0x32, 0x61, 0x37, 0x38, 0x63, 0x62, 0x66, 0x61, 0x39, 0x61, 0x61, 0x36, 0x38, 0x34, 0x62, 0x39, 0x31, 0x66, 0x31, 0x61, 0x63, 0x35, 0x36, 0x61, 0x38, 0x39, 0x36, 0x37, 0x36, 0x35, 0x65, 0x34, 0x64, 0x62, 0x66, 0x36, 0x65, 0x62, 0x30, 0x65, 0x34, 0x32, 0x63, 0x66, 0x22, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7d, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x6e, 0x64, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x72, 0x6f, 0x6f, 0x74, 0x66, 0x73, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x74, 0x61, 0x72, 0x2e, 0x67, 0x7a, 0x69, 0x70, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x3a, 0x20, 0x38, 0x32, 0x30, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x20, 0x22, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x3a, 0x33, 0x61, 0x31, 0x39, 0x65, 0x39, 0x63, 0x36, 0x63, 0x36, 0x36, 0x65, 0x39, 0x34, 0x30, 0x64, 0x32, 0x34, 0x39, 0x36, 0x63, 0x31, 0x65, 0x33, 0x31, 0x32, 0x64, 0x66, 0x38, 0x35, 0x62, 0x36, 0x37, 0x65, 0x39, 0x61, 0x37, 0x65, 0x32, 0x66, 0x36, 0x33, 0x64, 0x63, 0x32, 0x35, 0x61, 0x31, 0x30, 0x32, 0x66, 0x31, 0x61, 0x65, 0x35, 0x38, 0x35, 0x36, 0x32, 0x61, 0x65, 0x34, 0x34, 0x32, 0x22, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7d, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x6e, 0x64, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x72, 0x6f, 0x6f, 0x74, 0x66, 0x73, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x74, 0x61, 0x72, 0x2e, 0x67, 0x7a, 0x69, 0x70, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x3a, 0x20, 0x36, 0x37, 0x38, 0x2c, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x22, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x20, 0x22, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x3a, 0x39, 0x61, 0x39, 0x33, 0x33, 0x62, 0x66, 0x33, 0x37, 0x36, 0x39, 0x31, 0x62, 0x65, 0x30, 0x38, 0x63, 0x35, 0x32, 0x36, 0x66, 0x34, 0x37, 0x61, 0x30, 0x30, 0x36, 0x32, 0x63, 0x35, 0x36, 0x36, 0x31, 0x31, 0x31, 0x31, 0x34, 0x34, 0x64, 0x64, 0x31, 0x32, 0x62, 0x34, 0x30, 0x65, 0x37, 0x33, 0x35, 0x64, 0x30, 0x62, 0x32, 0x37, 0x33, 0x39, 0x66, 0x61, 0x34, 0x63, 0x37, 0x38, 0x36, 0x38, 0x22, 0xa, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7d, 0xa, 0x20, 0x20, 0x20, 0x5d, 0xa, 0x7d},
}

*/

/*
json ответ от registry
GET https://registry.flant.com/v2/sys/antiopa/manifests/play_registry

{
   "schemaVersion": 2,
   "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
   "config": {
      "mediaType": "application/vnd.docker.container.image.v1+json",
      "size": 8236,
      "digest": "sha256:bd6693bf36f2a3efd4998d4fc719bacbe0e36dcc2899112bdf9ba258fdb90c64"
   },
   "layers": [
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 47536248,
         "digest": "sha256:f5c64a3438f6e850c2d09c6cdd407dff522129d5304fdd3a9dbb397c9488e7aa"
      },
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 848,
         "digest": "sha256:51899d335aae78435d7960bc74ca06c3a1db706bcb13c2dc19e2fcf567153f97"
      },
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 621,
         "digest": "sha256:6ff2b7de3c1376e8f6f0d823d7c15d5a32df338d7bc9013841dff4cefc2ba2ca"
      },
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 853,
         "digest": "sha256:50366c31f7fd3f86c9718ad3724dcd6926ba3f69f25f4bed47be5ca9640a0616"
      },
     {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 169,
         "digest": "sha256:f441b9a68d131e4aecdfeca0d35f7c984b6a66582f0d4b9e952ffb5c45c196aa"
      },
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 147,
         "digest": "sha256:446a15856c11ef8ec3c9ee74591277a7f155c771308b1692addc84ba55d7a79f"
      },
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 196742897,
         "digest": "sha256:f053c12ab4844d72e2ec90f38dab888d58826ef2a8ace62733f2dfb7052a46d1"
      },
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 8604447,
         "digest": "sha256:ae2fd08a532b8e931eb04e0781f48fddd64171615d1c243b5754d662a734d0e6"
      },
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 825,
         "digest": "sha256:9d96be0c98bea0f7aef457b615db9f14613e06548d663fc4ad329b8345790b94"
      },
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 680,
         "digest": "sha256:25da774b242d1ec5491911adfdc382c28959d8985f2147ecc950fbd3fd4f8585"
      }
   ]
}

*/

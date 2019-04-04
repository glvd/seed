# SEED
シードはIPFSを利用して、ヴィデオをアップして、内容の索引を作ってツールです。

JSONファイルのExample:
```
[
 [
   {
     "bangumi": "avvd01312",
     "alias": [
       "Lady Ninja: Aoi kage",
       "极乐女忍者"
     ],
     "slice": true,
     "files": [
       "D:\\video\\Okita Souji.mp4",
       "D:\\video\\HD Epic Sax Gandalf.mp4"
     ],
     "Sharpness": "1080P",
     "poster_path": "D:\\video\\H\\p2526321708.webp",
     "role": [
       "叶加濑麻衣",
       "赤井沙希",
       "鸟肌实",
       "坂口征夫",
       "阿部祐二"
     ],
     "publish": "2018-01-20"
   },
   ...
]
```

処理完了のJSONファイル：
```
[
    {
        "bangumi": "avvd01312",
        "type": "",
        "output": "",
        "vr": "",
        "thumb": "",
        "intro": "",
        "alias": [
            "Lady Ninja: Aoi kage",
            "极乐女忍者"
        ],
        "language": "",
        "caption": "",
        "poster": "QmZPAqK6Hc2pwrmU3mPJLdMUGpTadaBhN8vFwJ9ay2ULBe",
        "role": [
            "叶加濑麻衣",
            "赤井沙希",
            "鸟肌实",
            "坂口征夫",
            "阿部祐二"
        ],
        "director": "",
        "season": "",
        "Episode": "",
        "TotalEpisode": "",
        "sharpness": "",
        "group": "",
        "publish": "2018-01-20",
        "video_group_list": [
            {
                "index": "2c381fcf095cb08450bcb93d6a3ba707dd8f9b43",
                "sharpness": "",
                "sliced": true,
                "hls": {
                    "encrypt": false,
                    "key": "",
                    "m3u8": "media",
                    "segment_file": "media-%05d.ts"
                },
                "object": [
                    {
                        "links": [
                            {
                                "hash": "QmNrSKRarcXrqFsGDUg86aqrPJPWXAj2wJ7wu9QaKqRpdN",
                                "name": "media",
                                "size": 200,
                                "type": 2
                            },
                            {
                                "hash": "Qmct4Ejhg8oYkR1gfFfdwewcT2QmPjngnzymuvsC4pyXsh",
                                "name": "media-00000.ts",
                                "size": 5299903,
                                "type": 2
                            },
                            {
                                "hash": "Qmbkdwohj6nZzyGgoiHJy91y4M1ckBfQbkSFNLRWa8Gd8w",
                                "name": "media-00001.ts",
                                "size": 3500497,
                                "type": 2
                            },
                            {
                                "hash": "QmcDA6tRx3H5oJddWXKLhyDN9SBXiYQ3QmjDupr7KaCqDe",
                                "name": "media-00002.ts",
                                "size": 3838207,
                                "type": 2
                            }
                        ],
                        "hash": "QmYqtCY7fbNp79tfEgu19BwPe9eSppwXJuLhnwgN16vCpi",
                        "name": "D:\\workspace\\goproject\\seed\\tmp\\c8974da7-365c-48f6-8b31-9319115535d2",
                        "size": 12639036,
                        "type": 2
                    },
                    {
                        "links": [
                            {
                                "hash": "QmcZ1CbwsW7NodvQgc4jgYiGMnAhUYCK9xzyjULX9JAoZt",
                                "name": "media",
                                "size": 504,
                                "type": 2
                            },
                            {
                                "hash": "QmWrZfnT1uDmYo9mSkkLYzAY1qJx4FbCk3Raood2im7xbx",
                                "name": "media-00000.ts",
                                "size": 1577945,
                                "type": 2
                            },
                            {
                                "hash": "QmbZ2bDrLYWREVZbhURMfNqScoZpxyK2s52QZNdbeTz5ev",
                                "name": "media-00001.ts",
                                "size": 1535776,
                                "type": 2
                            },
                            {
                                "hash": "QmWmKvx7caP9MsNttEYbaTx2JP75DwJTtDNSUmYwSNcpcS",
                                "name": "media-00002.ts",
                                "size": 1473924,
                                "type": 2
                            },
                            {
                                "hash": "QmYQC44Hx1VyXh3qWxwZpGzuqBGnLnzdQjk3JXi4MmPQBf",
                                "name": "media-00003.ts",
                                "size": 1583773,
                                "type": 2
                            },
                            {
                                "hash": "QmbV4tLYjv3TMG5ecqx1LW2ehrPkyLARnT4VL9DzxLwvN1",
                                "name": "media-00004.ts",
                                "size": 1569616,
                                "type": 2
                            },
                            {
                                "hash": "QmNX5Qx97iueq7U17Q3N1AQBgacZAMWZtvs8k26tDUPaWa",
                                "name": "media-00005.ts",
                                "size": 2724999,
                                "type": 2
                            },
                            {
                                "hash": "QmTEanGcdLk6Y21SYLoH7NuHurSGGk6L8ZwB3PfgBHJNwk",
                                "name": "media-00006.ts",
                                "size": 1550252,
                                "type": 2
                            },
                            {
                                "hash": "QmakcZeTtHG56ofTkDTrzPAWK5G5ALyipn52dCDcRDiUk7",
                                "name": "media-00007.ts",
                                "size": 1491596,
                                "type": 2
                            },
                            {
                                "hash": "QmNxiLb7YgWZhMW4UkHNRS7Cwmr11d2PFTxpma1FJDAAS6",
                                "name": "media-00008.ts",
                                "size": 1176822,
                                "type": 2
                            },
                            {
                                "hash": "QmXMxPLc9Tpn5ZNjSfRFfZWayM8agBNGQUfrigbTr7ZxXD",
                                "name": "media-00009.ts",
                                "size": 1550252,
                                "type": 2
                            },
                            {
                                "hash": "QmSf8nR1NjPSksDfzLts8rs1upbb5EXUEJVcBmuz7XCJWX",
                                "name": "media-00010.ts",
                                "size": 1542732,
                                "type": 2
                            },
                            {
                                "hash": "QmPyUsJPZCAtKoeEUAAiYeBETPdvWv9bz2fK7vpbbWyjTe",
                                "name": "media-00011.ts",
                                "size": 195534,
                                "type": 2
                            }
                        ],
                        "hash": "QmQ2DcWorM3fkyeUnLmHG9Mm18JxttcD9pJYFbZ6z3y67c",
                        "name": "D:\\workspace\\goproject\\seed\\tmp\\b05de57f-d57c-4c5c-b76c-d2e6899335ab",
                        "size": 17974474,
                        "type": 2
                    }
                ]
            }
        ],
        "source_info_list": [
            {
                "id": "QmQwp6XEcXauaiBjt3bUsSwcy7SomtE5ve74YRVfi2vCba",
                "public_key": "CAASpgIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC/voc3C3ZpdTuudhWcZIje7M2gE/V3Av33bUXqhaThmd/MVQZ1CVkzWSJmuYiv7Lq+Uop+pGwF2lKuNQ0xm3YeQrDTOgBkR2Gv8rF69VFdk1olIsFB7XvXwhpREP/l0BZ0u8hYML/gaePzDGpSXRSQi/tzgAbCwirhr18Vd+bh5VuaNcsT/hyCWxJv+uvPCnrrfBrViT7T6w+Oo2fcu8BhbyngBbT01aIhHwFbxuDuMHHDUhi6WvH5k3QPbEp9VNcLXmSiS6Lxvy/5tpUeBOosS99HtiWfswlOOJcjj0jfyMMcPXNgJNLfcoFw+idxW8aUgGnpTS0fixEEa0lk0Kz3AgMBAAE=",
                "addresses": [
                    "/ip4/10.0.75.1/tcp/4001/ipfs/QmQwp6XEcXauaiBjt3bUsSwcy7SomtE5ve74YRVfi2vCba",
                    "/ip4/169.254.112.168/tcp/4001/ipfs/QmQwp6XEcXauaiBjt3bUsSwcy7SomtE5ve74YRVfi2vCba",
                    "/ip4/169.254.206.232/tcp/4001/ipfs/QmQwp6XEcXauaiBjt3bUsSwcy7SomtE5ve74YRVfi2vCba",
                    "/ip4/169.254.248.113/tcp/4001/ipfs/QmQwp6XEcXauaiBjt3bUsSwcy7SomtE5ve74YRVfi2vCba",
                    "/ip4/192.168.1.47/tcp/4001/ipfs/QmQwp6XEcXauaiBjt3bUsSwcy7SomtE5ve74YRVfi2vCba",
                    "/ip4/169.254.55.142/tcp/4001/ipfs/QmQwp6XEcXauaiBjt3bUsSwcy7SomtE5ve74YRVfi2vCba",
                    "/ip4/127.0.0.1/tcp/4001/ipfs/QmQwp6XEcXauaiBjt3bUsSwcy7SomtE5ve74YRVfi2vCba",
                    "/ip4/172.29.42.81/tcp/4001/ipfs/QmQwp6XEcXauaiBjt3bUsSwcy7SomtE5ve74YRVfi2vCba",
                    "/ip6/::1/tcp/4001/ipfs/QmQwp6XEcXauaiBjt3bUsSwcy7SomtE5ve74YRVfi2vCba"
                ],
                "agent_version": "go-ipfs/0.4.19-dev/",
                "protocol_version": "ipfs/0.1.0"
            }
        ]
    },
    ...
]
```
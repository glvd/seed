# SEED
シードはIPFSを利用して、ヴィデオをアップして、内容の索引を作ってツールです。

アテンション：
番号捜索のみ使用すれば、[银河VR共享总部](https://t.me/yinhevr)を入れって、
```
/video 番号 
```
を使って、番号を捜索してください。

或いは、Telegramのローボード@yinhe_botを使いって、
```
/video 番号 
```
を使って、番号を捜索してください。

SEEDの使い方法：
```
SEED process //seed.jsonを読み出して、ipfsネートをアップします。

SEED pin 番号　//番号を探して、IPFSネートにPINをします。番号を入力しなければ、全ての番号PINをします。

SEED transfer //DBからjsonファイルへ転換する。

SEED update //jsonファイルを読み出して、ビデオの紹介を更新します。
```

JSONファイルのExample:
```
[
 [
   {
     "bangumi": "video",
     "alias": [
       "video1"
     ],
     "slice": true,
     "files": [
       "D:\\video\\Okita Souji.mp4",
     ],
     "Sharpness": "1080P",
     "poster_path": "D:\\video\\H\\p2526321708.webp",
     "role": [
       ""
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
        "video_group_list": [],
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

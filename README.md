# iyashi-apex

[iyashi](https://github.com/mix3/iyashi) を API Gateway + Lambda with [Apex](https://github.com/apex/apex) & [ridge](https://github.com/fujiwara/ridge) で動かせるようにしてみた

## Usage

### API Gateway

プロキシリソースを設定する

<img src="https://cloud.githubusercontent.com/assets/36567/19107145/cff2504a-8b25-11e6-9eea-5d508029bcf5.png" width="487">

### Lambda with Apex

各 token を環境変数で設定して deploy

```
apex deploy \
    -s SLACK_BOT_TOKEN="xxxx-00000000000-xxxxxxxxxxxxxxxxxxxxxxxx" \
    -s FLICKR_API_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
    -s TUMBLR_API_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
    iyashi
```

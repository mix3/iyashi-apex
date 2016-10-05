# iyashi-apex

[iyashi](https://github.com/mix3/iyashi) を API Gateway + Lambda with [Apex](https://github.com/apex/apex) & [ridge](https://github.com/fujiwara/ridge) で動かせるようにしてみた

<img width="700" alt="2016-10-05 11 02 16" src="https://cloud.githubusercontent.com/assets/36567/19098803/56f3e378-8aeb-11e6-8f5f-e6ab43c202e5.png">

## Usage

### API Gateway

slack の webhook を API Gateway で受けて POST パラメータを JSON に変換して Lambda に投げる

#### mapping template sample

```
Content-Type: application/x-www-form-urlencoded
```
```
{
    "path" : "$context.resourcePath",
    "queryStringParameters" : {
#foreach( $kvstr in $input.body.split( '&' ) )
#set( $kv = $kvstr.split( '=' ) )
        "$util.urlDecode($kv[0])" : "$util.urlDecode($kv[1])"#if( $foreach.hasNext ),#end
#end
    }
}
```

### Lambda with Apex

各 token を環境変数で設定して deploy

```
apex deploy \
    -s SLACK_BOT_TOKEN="xxxx-00000000000-xxxxxxxxxxxxxxxxxxxxxxxx" \
    -s FLICKR_API_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
    -s TUMBLR_API_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
    iyashi
```

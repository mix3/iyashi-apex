# iyashi-apex

[iyashi](https://github.com/mix3/iyashi) を API Gateway + Lambda with Apex で動かせるようにしてみた


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


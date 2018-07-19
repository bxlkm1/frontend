<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1; maximum-scale=1.0; user-scalable=no;" />
    <title>ChinaJoy</title>
    <link href="/static/css/style.css" rel="stylesheet" type="text/css" />
    <script src="/static/js/jquery.min.js" charset="utf-8"></script>
</head>

<body>
<div class="main">
    <div class="logo"></div>
    <div class="form form_home">
        <p style="margin-top:10px;">您已收到/You received {{.reward}} ELA</p>
        <a target="_blank" href="{{.elaWallet}}" style="margin-top:5px;">进入Elastos钱包<br />
            <span>Enter Elastos Wallet</span> </a>
        <a href="javascript:;">Elastos欢乐竞猜<br />
            <span>Elastos Lottery</span> </a>
        <a href="javascript:;">Elastos游戏中心<br />
            <span>Elastos Game Center</span> </a>
    </div>
    <div class="footer"><a href="/index/?openid={{.openid}}&isModify=1">修改钱包 Edit Wallet</a></div>
</div>
</body>
</html>

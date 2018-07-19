<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width; initial-scale=1; maximum-scale=1.0; user-scalable=no;" />
    <title>ChinaJoy</title>
    <link href="/static/css/style.css" rel="stylesheet" type="text/css" />
    <script src="/static/js/jquery.min.js" charset="utf-8"></script>
</head>

<body>
<div class="main">
    <div class="logo"></div>
    <div class="form"> <a href="javascript:{{.walletAddr}};">注册Elastos钱包<br />
        <span>Register Elastos Wallet</span> </a> <a href="javascript:;" id="input_address">输入钱包地址<br />
        <span>Enter Wallet Address</span> </a> </div>
    <div class="footer"><a href="javascript:;">帮助 Need help</a></div>
</div>
{{ $reg := .registed}}
<div class="dialog" {{if eq 1 $reg}} style="display: none" {{end}}>
    <div class="logo"></div>
    <div class="dialog_bg" id="dialog_bg"></div>
    <div class="dialog_form">
        <input type="text" name="address" id="elaAddr" value="{{.addr}}" />
    </div>
    <a href="javascript:;" id="goto"><img src="/static/images/Arrow.png" height="36" /></a>
</div>
<div class="dialog_loading" {{if eq 1 $reg}} style="display: block" {{end}}>
    <div class="logo"></div>
    <div class="dialog_loading_bg"></div>
    <div class="dialog_loading_main">
        <p class="fz14">您将收到一笔来自亦来云的奖励金，可用于 ChinaJoy 亦来云展位的竞猜、游戏等活动。</p>
        <p class="fz14">You will receive a transfer from Elastos, that can be used in Elastos’ Pavilion at chinaJoy.</p>
        <p class="fz11">*区块链交易分布于各广泛节点，具有匿名化，去中心化等特点，有一定时间延迟。请 2 分钟后刷新界面。</p>
        <p class="fz11">*Blockchain technology is decentralized and anonymous, the transfer is distributed in active blocks worldwide,  so it needs approximate 2 minutes to deposit to your wallet address. Please reload this page later.</p>
    </div>
    <img src="/static/images/Wave.gif" height="100" />
</div>
<script type="text/javascript">
    $("#input_address").click(function(){
        $(".dialog").show();
    });
    $("#dialog_bg").click(function(){
        $(".dialog").hide();
    });
    // $("#goto").click(function(){
    //     $(".dialog").hide();
    //     $(".dialog_loading").show();
    // });
    // function gourl(){
    //     location.href="home.html";
    // };
    $('#goto').click(function(){
        var elaAddr = $("#elaAddr").val();
        if (elaAddr == "" || elaAddr.length != 34 || !elaAddr.startsWith("E")){
            alert("please fill in the right ELA address")
            return
        }
        location.href="/submitAddr/"+elaAddr+"?vldCode="+{{.vldCode}}+"&openid="+{{.openid}}
    })

    console.info({{.registed}},{{.openid}})

    if ({{.registed}} === 1 ){
        console.log("start callback")
        setTimeout(callback,15000)
    }

    function callback(){
        retObj = $.ajax({url:"/registerUser/"+{{.openid}},async:false});
        console.log(retObj,JSON.parse(retObj.responseText))
        if(retObj.status == 200 && JSON.parse(retObj.responseText).Error == 0){
            location.href = "/home/"+{{.openid}}
        }
        setTimeout(callback,15000)
    }

</script>
</body>
</html>
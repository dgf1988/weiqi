/**
 * Created by dgf19 on 2015/10/30.
 */
//获取xml对象
/**
 * @return {boolean}
 */
function ajax_getXmlHttpObject()
{
    var xmlHttp=false;
    try
    {
// Firefox, Opera 8.0+, Safari
        xmlHttp=new XMLHttpRequest();
    }
    catch (e)
    {
// Internet Explorer
        try
        {
            xmlHttp=new ActiveXObject("Msxml2.XMLHTTP");
        }
        catch (e)
        {
            try{
                xmlHttp=new ActiveXObject("Microsoft.XMLHTTP");
            }
            catch(e){
                xmlHttp = false;
            }
        }
    }
    return xmlHttp;
}
function ajax_get(xmlHttp,serverUrl,callback){
    xmlHttp.onreadystatechange = function(){
        if(xmlHttp.readyState == 4 ){
            if( xmlHttp.status == 200 ){
                callback( ajax_jsonToObj( xmlHttp.responseText ) );
            }
        }
    };
    xmlHttp.open('get',serverUrl,true);
    xmlHttp.setRequestHeader("X-Requested-With","XMLHttpRequest");
    xmlHttp.send();
}
function ajax_post(xmlHttp,serverUrl,data,callback){
    xmlHttp.onreadystatechange = function(){
        if(xmlHttp.readyState == 4 ){
            if( xmlHttp.status == 200 ){
                callback( ajax_jsonToObj( xmlHttp.responseText ));
            }
        }
    };
    xmlHttp.open('post',serverUrl,true);
    xmlHttp.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    xmlHttp.setRequestHeader("X-Requested-With","XMLHttpRequest");
    var sends = [];
    for(var key in data) sends.push(key + '=' + data[key]);
    xmlHttp.send(sends.join('&'));
}

function ajax_objToJson(js_obj){
    var item = [];
    for(var key in js_obj) item.push(key + '=' + js_obj[key]);
    return item.join('&');
}



function ajax_jsonToObj(php_json_encode){
    return eval('('+ php_json_encode +')' );
}

// v2.0

function http_objectToParameter(obj){
    var item = [];
    for(var key in obj) item.push(key + '=' + obj[key]);
    return item.join('&');
}

function ajaxSend(action,method,listKeyValue,callback){
    var xmlHttp = ajax_getXmlHttpObject();
    if( !xmlHttp ){
        callback(null,'浏览器不支持。');
    }else if( listKeyValue ){
        xmlHttp.onreadystatechange = function(){
            if(xmlHttp.readyState == 4 ){
                if( xmlHttp.status == 200 ){
                    callback( ajax_jsonToObj( xmlHttp.responseText ) );
                }
            }
        };
        switch (method.toLowerCase()){
            case 'get':
                xmlHttp.open('get',action+'?'+http_objectToParameter(listKeyValue),true);
                xmlHttp.setRequestHeader("X-Requested-With","XMLHttpRequest");
                xmlHttp.send();
                break;
            case 'post':
                xmlHttp.open('post',action,true);
                xmlHttp.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
                xmlHttp.setRequestHeader("X-Requested-With","XMLHttpRequest");
                xmlHttp.send(http_objectToParameter(listKeyValue));
                break;
            default :
                callback(null,'http方法错误。');
                break;
        }
    }else{
        callback(null,'参数为空');
    }
}

function httpSend(action,method,listKeyValue){
    var form = document.createElement('form');
    form.style.display = 'none';
    form.action = action;
    form.method=  method;

    for( var key in listKeyValue){
        var input = document.createElement('input');
        input.name = key;
        input.value = listKeyValue[key];
        form.appendChild(input);
    }

    document.body.appendChild(form);
    form.submit();
}


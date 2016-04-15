/**
 * Created by dgf19 on 2015/11/14.
 */
//  棋盘脚本。
//棋盘初始化
var content = document.getElementById('weiqi-content');
weiqi_appendChessBoard(content,26);
var displaymap= document.getElementById('display');
var handlermap= document.getElementById('handler');
var stepnum = document.getElementById('stepnumber');
var setshow = document.getElementById('setshow');
var lastpoint = go_createPoint();
//数据结构初始化 - 创建一个围棋
var weiqi = api_createWeiqi();
weiqi.showNumber = setshow.checked;
stepnum.value = 0;

var saveweiqi = {};
var savestep = 0;
var saveshow = false;

//清空事件
function weiqi_clear(){
    stepnum.value = 0;
    weiqi_goto(-1);
    weiqi = api_createWeiqi();
}
//模式切换事件
function weiqi_mode(e){
    if(e.innerHTML == '试下'){
        e.innerHTML = '恢复';
        dom_setMapClickEvent(handlermap,weiqi_handler);
        handlermap.onmousedown = null;
        handlermap.oncontextmenu = null;
        saveweiqi = weiqi;
        savestep = parseInt(stepnum.value);
        saveshow = setshow.checked;
        weiqi = api_createWeiqi('','',go_nextPlayer(weiqi.steps[savestep].player) ,weiqi.maps[savestep]);
        weiqi.showNumber = true;
        setshow.checked = true;
    }else{
        e.innerHTML = '试下';
        dom_setMapClickEvent(handlermap,null);
        handlermap.onmousedown = function(event){
            if( event.button != 2 ){
                weiqi_forward();
            }
            return false;
        };
        handlermap.oncontextmenu = function(event){
            weiqi_back();
            return false;
        };
        saveweiqi.lastRefreshMap = weiqi.lastRefreshMap;
        weiqi = saveweiqi;
        weiqi.showNumber = saveshow;
        setshow.checked = saveshow;
        weiqi_goto(savestep);
    }
}

//下棋事件
function weiqi_handler(event){
    var target = dom_getTargetByEvent(event);
    var point = dom_getPointByTarget(target);
    weiqi_step(point);
    weiqi_goto();
}

function weiqi_step(point,player){
    var step = weiqi.tryStep(point,player);
    if( step){
        stepnum.value = step;
    }
    return step;
}

//跳转事件
function weiqi_back(){
    var num = parseInt(stepnum.value)-1;
    weiqi_goto(num);
}
function weiqi_forward(){
    var num = parseInt(stepnum.value)+1;
    weiqi_goto(num);
}
function weiqi_goto(num){
    if( !num )
        num = parseInt(stepnum.value);
    if(num < 0){
        num = 0 ;
    }else if( num > weiqi.getStepLen() ){
        num = weiqi.getStepLen();
    }
    if( weiqi.hasStep(num) ){
        weiqi.refreshQizi(displaymap,num);
        //dom_clearMapBGImg(handlermap);
        if( !go_isEmptyPoint(lastpoint)){
            var clear = dom_getMapCellByPoint(handlermap,lastpoint);
            dom_clearBGImg(clear);
        }
        var point = weiqi.steps[num].point;
        var cell = dom_getMapCellByPoint(handlermap,point);
        if( cell && ! weiqi.showNumber ){
            dom_setBGImg(cell,go_imgs[go_redPng],'35% 35%');
        }
        go_copyPoint(point,lastpoint);
        stepnum.value = num;
    }
}
//键盘输入检测事件
function weiqi_checkKeyDown(event){
    if(!((event.keyCode>=48&&event.keyCode<=57)||(event.keyCode>=96&&event.keyCode<=105)||event.keyCode == 13 ||  event.keyCode == 8 )){
        event.returnValue=false;
        window.event.returnValue = false;
        return false;
    }
}
//输入数字验证事件
function weiqi_checkKeyUp(){
    var  num = parseInt(stepnum.value);
    if( num < 0 ){
        stepnum.value = 0;
    }else if( num > weiqi.getStepLen() ){
        stepnum.value = weiqi.getStepLen();
    }
    weiqi_goto();
}
//显示手数切换事件
function weiqi_numbershow(e){
    weiqi.showNumber = e.checked;

    var num = parseInt(stepnum.value);
    var point = weiqi.steps[num].point;
    dom_refreshMapNumber(weiqi.maps[num],displaymap,weiqi.showNumber);
    var cell = dom_getMapCellByPoint(handlermap,point);
    if( cell ){
        if( weiqi.showNumber){
            dom_clearBGImg(cell);
        }else{
            dom_setBGImg(cell,go_imgs[go_redPng],'35% 35%');
        }
    }
}
!function(){"use strict";var t="optimizeSelectOption";window.top.layui?window.top.layui.use("layer",function(){layui.define(["form"],function(e){e(t,function(u){var e="0.2.0",c=layui.$,t=layui.form,n=layui.layer,o=layui.cache.modules.optimizeSelectOption.substr(0,layui.cache.modules.optimizeSelectOption.lastIndexOf("/"));layui.link(o+"/optimizeSelectOption.css?v"+e);var f=[".layui-table-view",".layui-layer-content",".select_option_to_layer"];if(window.top.layer._indexTemp=window.top.layer._indexTemp||{},!t.render.plugFlag){var s=t.render;t.render=function(l,e,t){var a=this;var r;if(t){layui.each(t,function(e,t){t=c(t);var n=t.parent();var o=n.hasClass("layui-form");var i=n.attr("lay-filter");o?"":n.addClass("layui-form");i?"":n.attr("lay-filter","tablePlug_form_filter_temp_"+(new Date).getTime()+"_"+Math.floor(Math.random()*1e5));r=s.call(a,l,n.attr("lay-filter"));o?"":n.removeClass("layui-form");i?"":n.attr("lay-filter",null)})}else{r=s.call(a,l,e)}return r};t.render.plugFlag=true}var p=function e(){window.top.layer.close(window.top.layer._indexTemp[u])};function d(e,t,n){t=t||window;e=e.length?e.get(0):e;var o={};if(n&&t.top!==t.self){var i=t.frames.frameElement;o=d(i,t.parent,n)}var l=e.getBoundingClientRect();return{top:l.top+(o.top||0),left:l.left+(o.left||0)}}var i={},l=function e(t,s){var n=this;if(i.name){console.warn("针对",t,"的显示优化已经存在，请不要重复渲染！");return}c(document).on("click",f.map(function(e){return e+" "+s.triggerElem}).join(","),function(e){layui.stope(e);p();var t=c(this);var i=t;var l=typeof s.dlElem==="function"?s.dlElem(t):i.next();var n=i.parent().prev();var a=i.parent().hasClass("layui-form-selectup");function o(){var e=d(i,window,true);var t=e.top;var n=e.left;if(a){t=t-l.outerHeight()+i.outerHeight()-parseFloat(l.css("bottom"))}else{t+=parseFloat(l.css("top"))}if(t+l.outerHeight()>window.top.innerHeight&&!a){a=true;t-=l.outerHeight()+(2*parseFloat(l.css("top"))-i.outerHeight())}return{top:t,left:n}}var r=o();i.css({backgroundColor:"transparent"});window.top.layer._indexTemp[u]=window.top.layer.open({type:1,title:false,closeBtn:0,shade:0,anim:-1,fixed:i.closest(".layui-layer-content").length||window.top!==window.self,isOutAnim:false,offset:[r.top+"px",r.left+"px"],area:l.outerWidth()+"px",content:'<div class="layui-unselect layui-form-select layui-form-selected"></div>',skin:"layui-option-layer",success:function e(t,n){l.css({top:0,position:"relative"}).appendTo(t.find(".layui-layer-content").css({overflow:"hidden"}).find(".layui-form-selected"));t.width(i.width());var o=window.top.innerHeight-t.outerHeight()-parseFloat(t.css("top"));a&&t.css({top:"auto",bottom:o+"px"});typeof s.success==="function"&&s.success.call(this,n,t);t.on("mousedown",function(e){layui.stope(e)});setTimeout(function(){i.parentsUntil(f.join(",")).one("scroll",function(e){p()});i.parents(f.join(",")).one("scroll",function(e){p()});var e=window;do{var t=e.$||e.layui.$;if(t){t(e.document).one("click",function(e){p()});t(e.document).one("mousedown",function(e){p()});t(e).one("resize",function(e){p()});t(e.document).one("scroll",function(){if(top!==self&&parent.parent){p()}})}}while(e.self!==e.top?e=e.parent:false)},500)},end:function e(){typeof s.end==="function"&&s.end.call(this,n)}})})};return function e(t,s){var n=this;if(i.name){console.warn("针对",t,"的显示优化已经存在，请不要重复渲染！");return}c(document).on("click",f.map(function(e){return e+" "+s.triggerElem}).join(","),function(e){layui.stope(e);p();var t=c(this);var i=t;var l=typeof s.dlElem==="function"?s.dlElem(t):i.next();var n=i.parent().prev();var a=i.parent().hasClass("layui-form-selectup");function o(){var e=d(i,window,true);var t=e.top;var n=e.left;if(a){t=t-l.outerHeight()+i.outerHeight()-parseFloat(l.css("bottom"))}else{t+=parseFloat(l.css("top"))}if(t+l.outerHeight()>window.top.innerHeight&&!a){a=true;t-=l.outerHeight()+(2*parseFloat(l.css("top"))-i.outerHeight())}return{top:t,left:n}}var r=o();i.css({backgroundColor:"transparent"});window.top.layer._indexTemp[u]=window.top.layer.open({type:1,title:false,closeBtn:0,shade:0,anim:-1,fixed:i.closest(".layui-layer-content").length||window.top!==window.self,isOutAnim:false,offset:[r.top+"px",r.left+"px"],area:l.outerWidth()+"px",content:'<div class="layui-unselect layui-form-select layui-form-selected"></div>',skin:"layui-option-layer",success:function e(t,n){l.css({top:0,position:"relative"}).appendTo(t.find(".layui-layer-content").css({overflow:"hidden"}).find(".layui-form-selected"));t.width(i.width());var o=window.top.innerHeight-t.outerHeight()-parseFloat(t.css("top"));a&&t.css({top:"auto",bottom:o+"px"});typeof s.success==="function"&&s.success.call(this,n,t);t.on("mousedown",function(e){layui.stope(e)});setTimeout(function(){i.parentsUntil(f.join(",")).one("scroll",function(e){p()});i.parents(f.join(",")).one("scroll",function(e){p()});var e=window;do{var t=e.$||e.layui.$;if(t){t(e.document).one("click",function(e){p()});t(e.document).one("mousedown",function(e){p()});t(e).one("resize",function(e){p()});t(e.document).one("scroll",function(){if(top!==self&&parent.parent){p()}})}}while(e.self!==e.top?e=e.parent:false)},500)},end:function e(){typeof s.end==="function"&&s.end.call(this,n)}})})}("layuiSelect",{triggerElem:"div:not(.layui-select-disabled)>.layui-select-title",success:function(e,t){t.find("dl dd").click(function(){p()})},end:function(e){t.render("select",null,e)}}),{version:e,getPosition:d,close:p}}(t))})}):(console.warn("使用插件："+t+"页面顶层窗口必须引入layui"),layui.define(["form"],function(e){e(t,{msg:"使用插件："+t+"页面顶层窗口必须引入layui"})}))}();
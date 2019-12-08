layui.define(function (exports) {
    layui.use(['Swiper', 'form'], async function () {
        const $ = layui.$,Swiper = layui.Swiper, form = layui.form;
        var search_type = 1;
        //导航条
        on_scroll = function () {
            //变量t是滚动条滚动时，距离顶部的距离
            var t = document.documentElement.scrollTop || document.body.scrollTop;
            //当滚动到距离顶部200px时，返回顶部的锚点显示
            if (t >= 10) {
                $('#custom-header').addClass('header_bottom');
                $('#logo').attr('src', '/static/images/logo.png');

            } else {          //恢复正常
                $('#custom-header').removeClass('header_bottom');
                $('#logo').attr('src', '/static/images/logo_top.png');
            }
        };
        window.onscroll = on_scroll;
        $(on_scroll);
        //搜索tab
        $('.tab-search').click(function (data) {
            var parent = this.parentNode.getElementsByClassName('show');
            for (var i = 0; i < parent.length; i++) {
                parent[i].classList.remove('show');
            }
            this.classList.add('show');
            var curr = data.currentTarget.attributes[1].nodeValue;
            if (curr == 'key1') {
                info = '输入商品名称或商品编码';
                search_type = 1;
            } else if (curr == 'key2') {
                info = '输入商品名称/商品编码/CIQ代码';
                search_type = 2;
            }
            $('#search_input').attr('placeholder', info);
        });
        //监听搜索按钮
        form.on('submit(search)', function (data) {
            location.href = "/search?key=" + data.field.key + '&&type=' + search_type;
            return false;
        });
        //自定义验证
        form.verify({
            account: function (value, item) { //value：表单的值、item：表单的DOM对象
                if (!(/^1\d{10}$/.test(value) || /^([a-zA-Z0-9_\.\-])+\@(([a-zA-Z0-9\-])+\.)+([a-zA-Z0-9]{2,4})+$/.test(value))) {
                    return '输入格式不正确';
                }
            }
            , password: [
                /^[\S]{6,16}$/
                , '密码必须6到16位，且不能出现空格'
            ]
        });
    });
    exports('login', {});
});

layui.define(function (exports) {
    layui.use(['Swiper', 'form'], async function () {
        const $ = layui.$,
            Swiper = layui.Swiper,
            form = layui.form,
            admin = layui.admin;

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


        // //新闻数据
        // try {
        //     var data_come_type1 = await admin.get(`/showNews?type=1`);
        //     var data_come_type2 = await admin.get(`/showNews?type=2`);
        // } catch (e) {
        //     console.log('error')
        // }


        // shownews = function (datas, type, url) {
        //     $('#news-ul-' + type).empty();
        //     for (var i = 0; i < datas.data.length; i++) {
        //         $('#news-ul-' + type).append('<li>\n' +
        //             '<a href="/index_lists/' + datas.data[i].id + '">\n' +
        //             '        <p class="date">\n' +
        //             '            <span class="year">' + datas.data[i].year + '</span>\n' +
        //             '            <span class="xian"></span>\n' +
        //             '            <span class="month-day">' + datas.data[i].month + '</span>\n' +
        //             '        </p>\n' +
        //             '        <h3>' + datas.data[i].title + '</h3>\n' +
        //             '        <p class="desc">资讯来源：' + datas.data[i].author + '</p>\n' +
        //             '    </a>\n' +
        //             '</li>')
        //     };
        //     $('#news-ul-' + type).append('<li class="more"><a href="' + url + '">查看更多 ></a></li>');
        // };
        // $(shownews(data_come_type2, 2, '/index_lists?type=2'));


        // //新闻tab
        // $('.tab-news').click(function (data) {
        //     var parent = this.parentNode.getElementsByClassName('show');
        //     for (var i = 0; i < parent.length; i++) {
        //         parent[i].classList.remove('show');
        //     }
        //     this.classList.add('show');
        //     var open = data.currentTarget.attributes[1].nodeValue;
        //     var close = (open == 1) ? 2 : 1;
        //     $('#news-ul-' + open).css('display', 'block');
        //     $('#news-ul-' + close).css('display', 'none');
        //
        //     if (open == 2) {
        //         $(shownews(data_come_type2, 2, '/index_lists?type=2'));
        //     } else {
        //         $(shownews(data_come_type1, 1, '/index_lists?type=1'));
        //     }
        //
        // });

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

        const newsSwiper = new Swiper(".pc-bg-container .home-market-section .swiper-container-news", {
            loop: true,
            autoplay: false,
            navigation: {
                nextEl: ".swiper-news-button-next",
                prevEl: ".swiper-news-button-prev"
            }
        });
        const aboutSwiper = new Swiper(".about_swoper", {
            loop: true,
            autoplay: {
                autoplay: true,
                disableOnInteraction: false
            }
        });

        const mySwiper = new Swiper('.swiper-container', {
            direction: 'horizontal', // 垂直切换选项
            loop: true, // 循环模式选项
            autoplay: {
                autoplay: true,
                disableOnInteraction: false
            }, //自动播放
            // 如果需要分页器
            pagination: {
                el: '.swiper-pagination',
                bulletActiveClass: 'bullet-active'
            },
            navigation: {
                nextEl: '.swiper-button-next',
                prevEl: '.swiper-button-prev'
            }
        })


    });
    exports('login', {});
});

layui.define(function (exports) {
    layui.use(['form'], function () {
        const $ = layui.$,
            form = layui.form,
            admin = layui.admin;
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
        //发送验证码
        admin.sendAuthCode({
            elem: '#auth_reset_getsmscode'
            , elemPhone: '#auth_reset_phone'
            , elemVercode: '#auth_reset_smscode'
            , ajax: {
                url: get_sms_code //实际使用请改成服务端真实接口
            }
        });
    });
    exports('reset', {});
});

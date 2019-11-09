layui.define(function (exports) {
    layui.use(['form', 'admin', 'table', 'AutoComplete', 'laydate', 'laytpl', 'upload'], async function () {
        const { form, admin, table, $, laydate, laytpl, upload } = layui;
        const cusIEFlag = admin.cusIEFlag = 'I';
        let order_i_edit_data;
        /**企业资质数据**/
        let ent_qualif_data = []
            /**历史申报要素关闭层index**/
            , his_dec_index
            /**使用人数据**/
            , dec_users_data = []
            /**特殊业务标识数据**/
            , spec_decl_flag = []
            /**检验检疫签证申报要素数据**/
            , dec_request_certs_data = {}
            /**产品数据**/
            , order_pros_data = []
            /**检验检疫货物规格数据**/
            , goods_spec_data
            /**产品许可证数据**/
            , quas = []
            /**货物属性数据**/
            , goods_attr_data = []
            /**危险货物信息**/
            , dang_data
            /**批量修改--货物属性数据**/
            , goods_attr_batch_data = []
            /**集装箱数据**/
            , order_containers = []
            /**随附单证数据**/
            , order_documents = []
            /**企业承诺**/
            , declaratio_material_code = 1
            /**订单ID**/
            , order_id
            /**附注数据**/
            , order_note_data
            /**获取附件类型**/
            , edoc_code = 'local_2'
            , edoc_code_name = '原始资料'
            /**商品序号**/
            , order_pros_index = null
            , no_weight_save = false;
        /**回车键光标跳转**/
        admin.keydown_input_textarea();
        /**首次进入自动聚焦 and 初始数据赋值**/
        const first_name = $("#custom_master_name").val();
        $("#custom_master_name").val("").focus().val(first_name);
        $("#i_e_date").val(admin.getCurrDate());
        $("#apl_date").val(admin.getCurrDate());
        admin.transModeControl(cusIEFlag);
        /**自动完成方法再次封装**/
        const auto_fn = admin.auto_fn;
        /**产品**/
        admin.getOrderTable(order_pros_data);
        /**光标进入全选**/
        $('body').on('focus', 'input,textarea', function (e) {
            $(this).select();
            let tipText = "";
            for (let item of admin.tipsJson) {
                if (item.id == this.id) {
                    tipText = item.name;
                    break;
                }
            }
            $("#tipsMessagetext").text(tipText);
        });
        /**table点击行选中**/
        admin.table_radio_click();
        /**table左右拖动**/
        admin.table_mousedown();
        /**合同备案号列表搜索**/
        $("body").on('input', '#manual_no_search', async function () {
            const text = $(this).val();
            const data = await admin.get(`/order/account_manual?search=${text}`);
            table.reload('manual_no_list_table', {
                data: data
            });
        });
        /**选择提运单号**/
        $('body').on('click', '#bill_no_open', async function () {
            layer.open({
                type: 1,
                title: '选择提运单号',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['910px', '730px'],
                content: $('#bill_no_list').html()
            });
            table.render({
                elem: '#bill_no_list_table'
                , toolbar: true
                , defaultToolbar: ['filter']
                , colFilterRecord: 'local'
                , url: `/car/lists?sortedBy=desc&orderBy=created_at`
                , response: {
                    countName: 'total'
                }
                , cols: [[
                    { type: 'radio' }
                    , { field: 'code', title: '车辆海关编号' }
                    , { field: 'number', title: '车牌号' }
                    , { field: 'driver', title: '司机姓名' }
                    , { field: 'phone', title: '司机电话' }
                    , { field: 'password', title: '进场密码' }
                ]]
                , limit: 10
                , page: true
                , height: 550
            });
        });
        /**保存提运单号**/
        $('body').on('click', '#bill_no_save', function () {
            const checkStatus = table.checkStatus('bill_no_list_table');
            if (checkStatus.data.length == 0) {
                layer.msg("请选择提运单号");
                return
            }
            $("#bill_no").val("").focus().val(checkStatus.data[0].code);
            layer.closeAll();
        });
        /**提运单号列表搜索**/
        $("body").on('input', '#bill_no_search', async function () {
            const text = $(this).val();
            table.reload('bill_no_list_table', {
                where: {
                    search: text
                },
                page: {
                    curr: 1
                }
            });
        });
        /**选择其他包装**/
        $('body').on('click', '#packageInfoBtn', async function () {
            layer.open({
                type: 1,
                title: '编辑其他包装信息',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['910px', '730px'],
                content: $('#other_packs_list').html(),
                success: function (layero, index) {
                    this.enterEsc = function (event) {
                        if (event.keyCode === 13) {
                            $("#other_packs_save").click();
                            return false;
                        }
                    };
                    $(document).on('keydown', this.enterEsc);
                },
                end: function () {
                    $("#gross_wet").focus();
                    $(document).off('keydown', this.enterEsc);
                }
            });
            const data = await admin.post(`/clearance/no_paginate`,{Type:"包装种类代码"})
            data.forEach(function (item, index) {
                item.index = index
            });
            table.render({
                elem: '#other_packs_list_table'
                , toolbar: true
                , defaultToolbar: ['filter']
                , colFilterRecord: 'local'
                , cols: [[
                    { type: 'checkbox' }
                    , { type: 'numbers', title: '序号' }
                    , { field: 'customs_code', title: '包装材料种类代码' }
                    , { field: 'name', title: '包装材料种类名称' }

                ]]
                , data: data
                , limit: data.length
                , height: 550
            });
            $(`.layui-table-view[lay-id='other_packs_list_table'] .layui-table-body tr input`).each(function (index, el) {
                el.checked = false;
            });
            if ($("#dec_other_packs").val() && $("#dec_other_packs").val() != 'null') {
                const dec_other_packs_data = JSON.parse($("#dec_other_packs").val());
                for (let dec_item of dec_other_packs_data) {
                    for (let data_item of data) {

                        if (dec_item.pack_type == data_item.customs_code) {

                            $(`.layui-table-view[lay-id='other_packs_list_table'] .layui-table-body tr[data-index=${data_item.index}] .layui-form-checkbox`).click();
                        }
                    }
                }
            }
            form.render('checkbox');
        });
        /**保存其他包装**/
        $('body').on('click', '#other_packs_save', function () {
            const checkStatus = table.checkStatus('other_packs_list_table');
            /*if (checkStatus.data.length == 0) {
                layer.msg("请选择包装");
                return
            }*/
            let data = [];
            for (let item of checkStatus.data) {
                data.push({
                    pack_qty: "0",
                    pack_type: item.customs_code,
                    pack_type_name: item.name

                })
            }
            $("#dec_other_packs").val(JSON.stringify(data));
            layer.closeAll();
        });
        /**填写备注**/
        $('body').on('click', '#note_s_open', async function () {
            layer.open({
                type: 1,
                title: '备注',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['800px', 'auto'],
                content: $('#note_s_template').html(),
                success: function (layero, index) {
                    this.enterEsc = function (event) {
                        if (event.keyCode === 13) {
                            $("#note_s_save").click();
                            return false;
                        }
                    };
                    $(document).on('keydown', this.enterEsc);
                },
                end: function () {
                    $(document).off('keydown', this.enterEsc);
                    let value = $("#note_s").val();
                    $("#note_s").val("").focus().val(value);
                }
            });
            $("#note_s_fu").focus();
            $("#note_s_fu").val($("#note_s").val());
        });
        $("body").on("input", "#note_s_fu", function () {
            $(".note_s_fu_number").text($(this).val().length)
        });
        /**保存备注**/
        $('body').on('click', '#note_s_save', function () {
            $("#note_s").val($("#note_s_fu").val());
            layer.closeAll();
        });
        /**标记唛码**/
        $("body").on("input", "#mark_no", function () {
            $(".mark_no_number span").text($(this).val().length)
        });
        /**填写标记唛码**/
        $('body').on('click', '#mark_no_open', async function () {
            layer.open({
                type: 1,
                title: '标记唛码',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['800px', 'auto'],
                content: $('#mark_no_template').html(),
                success: function (layero, index) {
                    this.enterEsc = function (event) {
                        if (event.keyCode === 13) {
                            $("#mark_no_save").click();
                            return false;
                        }
                    };
                    $(document).on('keydown', this.enterEsc);
                },
                end: function () {
                    $(document).off('keydown', this.enterEsc);
                    let value = $("#mark_no").val();
                    $("#mark_no").val("").focus().val(value);
                }
            });
            $("#mark_no_fu").focus();
            $("#mark_no_fu").val($("#mark_no").val());
            $(".mark_no_fu_number").text($("#mark_no").val().length);
        });
        $("body").on("input", "#mark_no_fu", function () {
            $(".mark_no_fu_number").text($(this).val().length)
        });
        /**保存标记唛码**/
        $('body').on('click', '#mark_no_save', function () {
            $("#mark_no").val($("#mark_no_fu").val());
            $(".mark_no_number span").text($("#mark_no_fu").val().length);
            layer.closeAll();
        });
        /**标记唛码回车跳转商品焦点**/
        $("body").on("keyup", "#mark_no", function (event) {
            const eCode = event.keyCode ? event.keyCode : event.which ? event.which : event.charCode;
            if (event.shiftKey != 1 && eCode == 13) {
                if ($("#non-disabled-btn").attr("isdeclistshow") == 0) {
                    if ($("#contr_item").is(":disabled")) {
                        $("#code_t_s").focus()
                    } else {
                        $("#contr_item").focus()
                    }
                }
            }
        });
        /**收缩切换**/
        $("#non-disabled-btn").click(function () {
            let iTag = $(this).find('i');
            if (iTag.hasClass("transform_xia")) {
                $(this).attr("isDecListShow", "1");
                iTag.removeClass("transform_xia");
                $(this).parent().parent().nextAll().show();
                $("#org_code_name").attr("lay-verify", "required");
                $("#vsa_org_code_name").attr("lay-verify", "required");
                $("#desp_date").attr("lay-verify", "required");
                $("#insp_org_name").attr("lay-verify", "required");
                $("#purp_org_name").attr("lay-verify", "required");
            } else {
                $(this).attr("isDecListShow", "0");
                iTag.addClass("transform_xia");
                $(this).parent().parent().nextAll().hide();
                $("#org_code_name").removeAttr("lay-verify");
                $("#vsa_org_code_name").removeAttr("lay-verify");
                $("#desp_date").removeAttr("lay-verify");
                $("#insp_org_name").removeAttr("lay-verify");
                $("#purp_org_name").removeAttr("lay-verify");
            }
        });
        $("#decListShow").click(function () {
            let iTag = $(this).find('i');
            if (iTag.hasClass("transform_xia")) {
                $(this).attr("isDecListShow", "1");
                iTag.removeClass("transform_xia");
                $(this).parent().parent().nextAll().show();
                $("#purpose_name").attr("lay-verify", "required");
            } else {
                $(this).attr("isDecListShow", "0");
                iTag.addClass("transform_xia");
                $(this).parent().parent().nextAll().hide();
                $("#purpose_name").removeAttr("lay-verify");
            }
        });
        /**编辑企业资质信息**/
        let ent_qualif_index = null;
        $('body').on('click', '#enterpriseBtn', async function () {
            layer.open({
                type: 1,
                title: '编辑企业资质信息',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['910px', '530px'],
                content: $('#ent_qualif_template').html()
            });
            $("#ent_qualif_type_name_tem").focus();
            form.render('checkbox');
            if (declaratio_material_code == 0 || declaratio_material_code == "") {
                $("#declaratio_material_code").next('div').click();
            }
            table.render({
                elem: '#ent_qualif_table'
                , toolbar: '#ent_qualif_tool'
                , defaultToolbar: ['filter']
                , colFilterRecord: 'local'
                , cols: [[
                    { type: 'checkbox' }
                    , { field: 'index', title: '序号' }
                    , { field: 'ent_qualif_type_code', title: '企业资质类别代码' }
                    , { field: 'ent_qualif_type_name', title: '企业资质类别名称' }
                    , { field: 'ent_qualif_no', title: '企业资质编号' }
                ]]
                , data: ent_qualif_data
                , limit: ent_qualif_data.length
                , height: 350
            });
            /**自动完成--企业资质类别**/
            auto_fn({
                data: admin.all_complete_data.enterprise_product,
                listDirection: false,
                id: ['#ent_qualif_type_name_tem'],
                after: ['#ent_qualif_type_code_tem']
            });
        });
        /**编辑企业资质信息--点击行反填数据**/
        table.on('row(ent_qualif_table)', function (obj) {
            $("#ent_qualif_type_code_tem").val(obj.data.ent_qualif_type_code);
            $("#ent_qualif_type_name_tem").val(obj.data.ent_qualif_type_name);
            $("#ent_qualif_no_tem").val(obj.data.ent_qualif_no);
            ent_qualif_index = obj.data.index;
        });
        /**编辑企业资质信息--保存**/
        form.on('submit(ent_qualif_submit)', async (data) => {
            delete data.field.layTableCheckbox;
            if (ent_qualif_index) {
                for (let item in data.field) {
                    ent_qualif_data[ent_qualif_index - 1][item] = data.field[item]
                }
            } else {
                ent_qualif_data.push(data.field);
            }
            ent_qualif_data.forEach((value, index) => {
                value.index = index + 1
            });
            table.reload('ent_qualif_table', {
                data: ent_qualif_data,
                limit: ent_qualif_data.length
            });
            $("#ent_qualif_type_code_tem").val("");
            $("#ent_qualif_type_name_tem").val("");
            $("#ent_qualif_no_tem").val("");
            $("#ent_qualif_type_code").val(ent_qualif_data[0].ent_qualif_type_code);
            $("#ent_qualif_type_name").val(ent_qualif_data[0].ent_qualif_type_name);
            $("#ent_qualif_type_name_tem").focus();
            ent_qualif_index = null;
        });
        /**编辑企业资质信息--删除**/
        table.on('toolbar(ent_qualif_table)', function (obj) {
            let checkStatus = table.checkStatus(obj.config.id);
            let checkData = checkStatus.data;
            switch (obj.event) {
                case 'add':
                    $("#ent_qualif_type_code_tem").val("");
                    $("#ent_qualif_type_name_tem").val("");
                    $("#ent_qualif_no_tem").val("");
                    ent_qualif_index = null;
                    $("#ent_qualif_type_name_tem").focus();
                    break;
                case 'delete':
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    layer.confirm('真的删除么', { title: '提示' }, async (index) => {
                        let sup_data = ent_qualif_data.filter(item => {
                            return checkData.every(item2 => {
                                return item.index != item2.index;
                            })
                        });
                        for (let item_apply of checkData) {
                            if (item_apply.id) {
                                admin.order_ent_qualif_delete_ids.push(item_apply.id);
                            }
                        }
                        ent_qualif_data = sup_data;
                        ent_qualif_data.forEach((value, index) => {
                            value.index = index + 1
                        });
                        table.reload('ent_qualif_table', {
                            data: ent_qualif_data,
                            limit: ent_qualif_data.length
                        });
                        if (ent_qualif_data.length === 0) {
                            $("#ent_qualif_type_code").val("");
                            $("#ent_qualif_type_name").val("");
                        } else {
                            $("#ent_qualif_type_code").val(ent_qualif_data[0].ent_qualif_type_code);
                            $("#ent_qualif_type_name").val(ent_qualif_data[0].ent_qualif_type_name);
                        }
                        $("#ent_qualif_type_code_tem").val("");
                        $("#ent_qualif_type_name_tem").val("");
                        $("#ent_qualif_no_tem").val("");
                        ent_qualif_index = null;
                        $("#ent_qualif_type_name_tem").focus();
                        layer.close(index);
                    });
                    break;
            }
        });
        /**编辑企业资质信息--按enter保存**/
        $("body").on("keyup", "#ent_qualif_no_tem", function (event) {
            if (event.keyCode == 13) {
                $("#ent_qualif_submit").click();
                $("#ent_qualif_type_name_tem").focus();
            }
        });
        /**编辑企业资质信息--左移**/
        $("#enterpriseUpBtn").click(function () {
            let entQualifSeq = $("#entQualifSeq").val();
            if (!entQualifSeq) {
                entQualifSeq = 0;
            }
            if (ent_qualif_data.length > 0) {
                if (parseInt(entQualifSeq) < 1) {
                    layer.msg('此为第一条记录');
                } else {
                    $("input[id='entQualifSeq']").val(parseInt(entQualifSeq) - 1);
                    $("input[id='ent_qualif_type_code']").val(ent_qualif_data[parseInt(entQualifSeq) - 1].ent_qualif_type_code);
                    $("input[id='ent_qualif_type_name']").val(ent_qualif_data[parseInt(entQualifSeq) - 1].ent_qualif_type_name);
                }
            } else {
                layer.msg('没有可查看的数据');
            }
        });
        /**编辑企业资质信息--右移**/
        $("#enterpriseDownBtn").click(function () {
            let entQualifSeq = $("#entQualifSeq").val();
            if (!entQualifSeq) {
                entQualifSeq = 0;
            }
            if (ent_qualif_data.length > 0) {
                if (parseInt(entQualifSeq) < ent_qualif_data.length - 1) {
                    $("input[id='entQualifSeq']").val(parseInt(entQualifSeq) + 1);
                    $("input[id='ent_qualif_type_code']").val(ent_qualif_data[parseInt(entQualifSeq) + 1].ent_qualif_type_code);
                    $("input[id='ent_qualif_type_name']").val(ent_qualif_data[parseInt(entQualifSeq) + 1].ent_qualif_type_name);
                } else {
                    layer.msg('此为最后一条记录');
                }
            } else {
                layer.msg('没有可查看的数据');
            }
        });
        /**监听企业承诺**/
        form.on('checkbox(declaratio_material_code)', function (data) {
            if ($(data.elem).is(':checked')) {
                declaratio_material_code = 1
            } else {
                declaratio_material_code = 0
            }
        });
        laydate.render({ elem: '#desp_date', theme: '#1E9FFF', fixed: true });
        /**编辑使用人信息**/
        let dec_users_index = null;
        $('body').on('click', '#dec_users_btn', async function () {
            layer.open({
                type: 1,
                title: '编辑使用人信息',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['910px', '530px'],
                content: $('#dec_users_template').html()
            });
            $("#use_org_person_code").focus();
            table.render({
                elem: '#dec_users_table'
                , toolbar: '#dec_users_tool'
                , defaultToolbar: ['filter']
                , colFilterRecord: 'local'
                , cols: [[
                    { type: 'checkbox' }
                    , { field: 'index', title: '序号' }
                    , { field: 'use_org_person_code', title: '使用单位联系人' }
                    , { field: 'use_org_person_tel', title: '使用单位联系电话' }
                ]]
                , data: dec_users_data
                , limit: dec_users_data.length
                , height: 350
            });
        });
        /**编辑使用人信息--点击行反填数据**/
        table.on('row(dec_users_table)', function (obj) {
            $("#use_org_person_code").val(obj.data.use_org_person_code);
            $("#use_org_person_tel").val(obj.data.use_org_person_tel);
            dec_users_index = obj.data.index;
        });
        /**编辑使用人信息--保存**/
        form.on('submit(dec_users_submit_tem)', async (data) => {
            delete data.field.layTableCheckbox;
            if (dec_users_index) {
                for (let item in data.field) {
                    dec_users_data[dec_users_index - 1][item] = data.field[item]
                }
            } else {
                dec_users_data.push(data.field);
            }
            dec_users_data.forEach((value, index) => {
                value.index = index + 1
            });
            table.reload('dec_users_table', {
                data: dec_users_data,
                limit: dec_users_data.length
            });
            $("#use_org_person_code").val("");
            $("#use_org_person_tel").val("");
            $("#use_org_person_code").focus();
            dec_users_index = null;
        });
        /**编辑使用人信息--删除**/
        table.on('toolbar(dec_users_table)', function (obj) {
            let checkStatus = table.checkStatus(obj.config.id);
            let checkData = checkStatus.data;
            switch (obj.event) {
                case 'add':
                    $("#use_org_person_code").val("");
                    $("#use_org_person_tel").val("");
                    $("#use_org_person_code").focus();
                    dec_users_index = null;
                    break;
                case 'delete':
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    layer.confirm('真的删除么', { title: '提示' }, async (index) => {
                        dec_users_index = null;
                        let sup_data = dec_users_data.filter(item => {
                            return checkData.every(item2 => {
                                return item.index != item2.index;
                            })
                        });
                        dec_users_data = sup_data;
                        dec_users_data.forEach((value, index) => {
                            value.index = index + 1
                        });
                        table.reload('dec_users_table', {
                            data: dec_users_data,
                            limit: dec_users_data.length
                        });
                        $("#use_org_person_code").val("");
                        $("#use_org_person_tel").val("");
                        $("#use_org_person_code").focus();
                        layer.close(index);
                    });
                    break;
            }
        });
        /**编辑使用人信息--按enter保存**/
        $("body").on("keyup", "#use_org_person_tel", function (event) {
            if (event.keyCode == 13) {
                $("#dec_users_submit_tem").click();
                $("#use_org_person_code").focus();
            }
        });
        /**特殊业务标识**/
        $("body").on("click", "#spec_decl_flag_open", function () {
            layer.open({
                type: 1,
                title: '特殊业务标识',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['680px', '300px'],
                content: $('#spec_decl_flag_template').html(),
                success: function (layero, index) {
                    this.enterEsc = function (event) {
                        if (event.keyCode === 13) {
                            layer.closeAll();
                            return false;
                        }
                    };
                    $(document).on('keydown', this.enterEsc);
                },
                end: function () {
                    if ($("#contr_item").is(":disabled")) {
                        $("#code_t_s").focus()
                    } else {
                        $("#contr_item").focus()
                    }
                    $(document).off('keydown', this.enterEsc);
                }
            });
            if (spec_decl_flag.length > 0) {
                $(".tableDiv_spec_td").each(function (index, element) {
                    if (spec_decl_flag[index] == '1') {
                        $(this).addClass("bgcolor");
                    }
                })
            }
        });
        /**特殊业务标识选择**/
        $("body").on("click", ".tableDiv_spec_td", function () {
            if ($(this).hasClass("bgcolor")) {
                $(this).removeClass("bgcolor")
            } else {
                $(this).addClass("bgcolor");
            }
            let data = $("#spec_decl_flag").val().split(",");
            if (data.includes($(this).attr("name"))) {
                data.splice(data.indexOf($(this).attr("name")), 1)
            } else {
                let spec_data = [];
                $(".tableDiv_spec_td").each(function () {
                    if ($(this).hasClass("bgcolor")) {
                        spec_data.push($(this).attr("name"))
                    }
                });
                data = spec_data;
            }
            $("#spec_decl_flag").val(data.join(","));
            spec_decl_flag = [];
            $(".tableDiv_spec_td").each(function () {
                if ($(this).hasClass("bgcolor")) {
                    spec_decl_flag.push(1)
                } else {
                    spec_decl_flag.push(0)
                }
            });
        });
        /**enter关联理由打开特殊业务标识**/
        $("body").on("keyup", "#correlation_reason_flag_name", function (event) {
            if (event.keyCode == 13) {
                $("#spec_decl_flag_open").click();
            }
        });
        /**检验检疫签证申报要素**/
        $('body').on('click', '#dec_request_certs_btn', async function () {
            layer.open({
                type: 1,
                title: '检验检疫签证申报要素',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['1110px', '630px'],
                content: $('#dec_request_certs_template').html()
            });
            let queryDocRep = []
                , check_data = [];
            const data = await admin.post(`/clearance/no_paginate`,{Type:"商检签证申报要素"});
            for (let item of data) {
                queryDocRep.push({
                    app_cert_code: item.customs_code,
                    appl_ori: 1,
                    appl_copy_quan: 2,
                    app_cert_name: item.name
                })
            }
            if (dec_request_certs_data.checkData) {
                check_data = dec_request_certs_data.checkData.map(function (item) {
                    return item.app_cert_code
                });
            }
            table.render({
                elem: '#dec_request_certs_table'
                , toolbar: true
                , defaultToolbar: ['filter']
                , colFilterRecord: 'local'
                , primaryKey: 'app_cert_code'
                , checkStatus: {
                    default: check_data
                }
                , cols: [[
                    { type: 'checkbox' }
                    , { type: 'numbers', title: '序号', width: 120 }
                    , { field: 'app_cert_code', title: '证书代码', width: 240 }
                    , { field: 'app_cert_name', title: '证书名称', width: 300 }
                    , { field: 'appl_ori', title: '正本数量', width: 240 }
                    , { field: 'appl_copy_quan', title: '副本数量', width: 240 }
                ]]
                , data: queryDocRep
                , limit: queryDocRep.length
                , height: 350
            });
            laydate.render({ elem: '#cmpl_dschrg_dt', theme: '#1E9FFF', fixed: true });
            $("#domestic_consignee_ename").val(dec_request_certs_data.domestic_consignee_ename);
            $("#overseas_consignor_cname").val(dec_request_certs_data.overseas_consignor_cname);
            $("#overseas_consignor_addr").val(dec_request_certs_data.overseas_consignor_addr);
            $("#cmpl_dschrg_dt").val(dec_request_certs_data.cmpl_dschrg_dt);
            $("#decl_goods_enames").val(dec_request_certs_data.decl_goods_enames);
        });
        /**保存检验检疫签证申报要素**/
        form.on('submit(dec_request_certs_submit)', function (data) {
            let checkData = table.checkStatus('dec_request_certs_table').data;
            dec_request_certs_data = data.field;
            dec_request_certs_data.checkData = checkData;
            let arr = [];
            for (let item of dec_request_certs_data.checkData) {
                arr.push(item.app_cert_name)
            }
            $("#dec_request_certs").val(arr.join(","));
            layer.closeAll()
        });


       

        /**选择备案序号**/
        $('body').on('click', '#contr_item_open', async function () {
            if (!$("#manual_no").val()) {
                layer.msg("请先填写备案号");
                $("#manual_no").focus();
                return
            }
            layer.open({
                type: 1,
                title: '选择[备案序号]',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['910px', '730px'],
                content: $('#contr_item_list').html(),
                success: function (layero, index) {
                    this.enterEsc = function (event) {
                        if (event.keyCode === 13) {
                            $("#contr_item_save").click();
                            return false;
                        }
                    };
                    $(document).on('keydown', this.enterEsc);
                },
                end: function () {
                    $("#g_qty").focus();
                    $(document).off('keydown', this.enterEsc);
                }
            });
            $("#contr_item_search").val($("#contr_item").val());

            /** 根据监管方式和进出口标志 获取备案手册账册表体  */
            admin.goods_materials_data = admin.get_goods_materials_data(cusIEFlag);

            let filter_data = admin.goods_materials_data.filter(item => {
                return (item.serial).indexOf($("#contr_item").val()) > -1
            });

            table.render({
                elem: '#contr_item_list_table'
                , toolbar: true
                , defaultToolbar: ['filter']
                , colFilterRecord: 'local'
                , primaryKey: 'serial'
                , checkStatus: {
                    default: [filter_data[0] ? filter_data[0].serial : 0]
                }
                , cols: [[
                    { type: 'radio' }
                    , { field: 'serial', title: '备案序号', width: 120 }
                    , { field: 'hs_code', title: '商品编码', width: 160 }
                    , { field: 'name', title: '商品名称', width: 260 }
                    , { field: 'special', title: '规格型号', width: 130 }
                    , { field: 'unit_one', title: '申报单位', width: 130 }
                    , { field: 'unit_two', title: '法定单位', width: 130 }
                    , { field: 'price', title: '单价', width: 90 }
                    , { field: 'moneyunit', title: '币制名称', width: 110 }
                    , { field: 'taxationlx', title: '征免性质', width: 110 }
                    , { field: 'manuplace', title: '产销国', width: 90 }
                ]]
                , data: filter_data
                , limit: 10
                , page: true
                , height: 550
            });
            $(".layui-table-view[lay-id='contr_item_list_table'] .layui-table-body tr[data-index='0'] .layui-form-radio").click();
        });
        /**保存备案序号**/
        $('body').on('click', '#contr_item_save', async function () {
            const data = table.checkStatus('contr_item_list_table').data;
            if (data.length == 0) {
                layer.msg("请选择备案号");
                return
            }
            $("#contr_item").val("").focus().val(data[0].serial);

            $("#code_t_s").val(data[0].hs_code);
            $("#g_name").val(data[0].name);
            $("#g_model").val(data[0].special);

            $("#g_unit").val(data[0].unit_one_code);
            $("#g_unit_name").val(data[0].unit_one);

            $("#decl_price").val(data[0].price);

            $("#trade_curr").val(data[0].moneyunit_code);
            if(!data[0].moneyunit){
                if(order_pros_data.length){
                    $("#trade_curr_name").val(order_pros_data[0].trade_curr_name);

                }
            }else{
                $("#trade_curr_name").val(data[0].moneyunit);
            }

            $("#first_unit").val(data[0].unit_two_code);
            $("#first_unit_name").val(data[0].unit_two);

            $("#origin_country").val(data[0].manuplace_code);
            $("#origin_country_name").val(data[0].manuplace);

            $("#g_qty").focus();
            const code_data = await admin.get(`/hs_code/lists?limit=0&search=${data[0].hs_code}`);
            if (code_data.data.length > 0) {
                if (code_data.data[0].unit2) {
                    $("#second_unit").val(code_data.data[0].unit2);
                    $("#second_unit_name").val(code_data.data[0].unit2_name);
                }
            }
            layer.closeAll();
        });
        /**备案序号列表搜索**/
        $("body").on('input', '#contr_item_search', async function () {
            const text = $(this).val();

            /** 根据监管方式和进出口标志 获取备案手册账册表体  */
            admin.goods_materials_data = admin.get_goods_materials_data(cusIEFlag);

            let filter_data = admin.goods_materials_data.filter(item => {
                return (item.serial).indexOf(text) > -1
            });

            table.reload('contr_item_list_table', {
                checkStatus: {
                    default: [filter_data[0] ? filter_data[0].serial : 0]
                },
                data: filter_data
            });
            $(".layui-table-view[lay-id='contr_item_list_table'] .layui-table-body tr[data-index='0'] .layui-form-radio").click();
        });
        /**保存商品编号**/
        $('body').on('click', '#code_t_s_save', function () {
            const data = table.checkStatus('code_t_s_list_table').data;
            if (data.length == 0) {
                layer.msg("请选择商品");
                return
            }
            admin.declaration_data = data[0].declaration;


            $("#code_t_s").val(data[0].code);

            // if (!($("#g_name").val().trim())) {
            //     $("#g_name").val(data[0].name);
            // }

            $("#first_unit").val(data[0].unit1);
            $("#first_unit_name").val(data[0].unit1_name);

            $("#second_unit").val(data[0].unit2);
            $("#second_unit_name").val(data[0].unit2_name);

            layer.closeAll();
        });
        /**保存规格型号**/
        let declaration_edit = false;
        let declaration_index = null;
        form.on('submit(declaration_save)', function (data) {
            declaration_save_btn(data);
        });
        form.on('submit(declaration_next)', function (data) {
            declaration_save_btn(data);
            next_dec('edit');
        });
        form.on('submit(declaration_last)', function (data) {
            declaration_save_btn(data);
            last_dec('edit');
        });
        $("body").on("click", "#declaration_next_look", async function () {
            next_dec('look');
        });
        $("body").on("click", "#declaration_last_look", async function () {
            last_dec('look');
        });
        function submit_btn(){
            $('#declaration_last').attr('lay-submit','');
            $('#declaration_next').attr('lay-submit','');
            form.render();
        }
        function declaration_save_btn(data) {
            if (declaration_index != null && declaration_edit) {
                order_pros_data[declaration_index].g_model = data.field.elements;
                layer.msg("修改成功！");
                table.reload('order_pros', {
                    data: order_pros_data,
                    limit: order_pros_data.length
                });
            } else {
                $("#g_model").val(data.field.elements);
                $("#g_qty").focus();
                layer.closeAll()
            }
        }
        /**选择检验检疫名称**/
        $('body').on('click', '#ciq_name_open', async function () {
            if (!$("#code_t_s").val()) {
                layer.msg("请先填写商品编号");
                $("#code_t_s").focus();
                return
            }
            const data = await admin.get(`/ciq/lists?limit=0&search=${$("#code_t_s").val()}`);
            layer.open({
                type: 1,
                title: '检验检疫编码列表',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['910px', '450px'],
                content: $('#ciq_name_list').html(),
                success: function (layero, index) {
                    this.enterEsc = function (event) {
                        if (event.keyCode === 13) {
                            $("#ciq_name_save").click();
                            return false;
                        }
                    };
                    $(document).on('keydown', this.enterEsc);
                },
                end: function () {
                    $("#g_qty").focus();
                    $(document).off('keydown', this.enterEsc);
                }
            });
            table.render({
                elem: '#ciq_name_list_table'
                , toolbar: true
                , defaultToolbar: ['filter']
                , colFilterRecord: 'local'
                , primaryKey: 'id'
                , checkStatus: {
                    default: [data.data[0].id]
                }
                , cols: [[
                    { type: 'radio' }
                    , { field: 'hs', title: 'HS代码', width: 150 }
                    , { field: 'name', title: '名称', width: 255 }
                    , { field: 'ciq_code', title: 'CIQ代码', width: 120 }
                    , { field: 'ciq_name', title: 'CIQ代码中文名称' }
                ]]
                , data: data.data
                , limit: data.data.length
                , height: 300
            });
        });
        /**保存检验检疫名称**/
        $('body').on('click', '#ciq_name_save', function () {
            const data = table.checkStatus('ciq_name_list_table').data;
            if (data.length == 0) {
                layer.msg("请选择编码");
                return
            }
            $("#ciq_code").val(data[0].ciq_code);
            $("#ciq_name").val(data[0].name);

            layer.closeAll();
        });
        /**编辑检验检疫货物规格**/
        $('body').on('click', '#goods_spec_open', async function () {
            layer.open({
                type: 1,
                title: '编辑检验检疫货物规格',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['680px', '350px'],
                content: $('#goods_spec_list').html()
            });
            $("#stuff").focus();
            laydate.render({ elem: '#prod_valid_dt', theme: '#1E9FFF', fixed: true });
            laydate.render({ elem: '#produce_date', theme: '#1E9FFF', fixed: true });
            if (goods_spec_data) {
                for (let item in goods_spec_data) {
                    $(`#${item}`).val(goods_spec_data[item])
                }
            }
        });
        /**保存检验检疫货物规格**/
        form.on('submit(goods_spec_save)', function (data) {
            let array = [];
            goods_spec_data = data.field;
            for (let item in goods_spec_data) {
                if (goods_spec_data[item]) {
                    array.push(goods_spec_data[item])
                }
            }
            $("#goods_spec_data").val(array.join(";"));
            layer.closeAll();
            $("#goods_attr_open").click()
        });
        /**编辑产品许可证/审批/备案信息**/
        let dec_goods_index = null;
        $('body').on('click', '#goodsTargetBtn', async function () {
            if (!($("#code_t_s").val())) {
                layer.msg("请填写商品编号");
                $("#code_t_s").focus();
                return
            }
            layer.open({
                type: 1,
                title: '编辑编辑产品许可证/审批/备案信息',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['910px', '530px'],
                content: $('#dec_goods_template').html()
            });
            $("#licenceCodeTs").val($("#code_t_s").val());
            $("#licenceGName").val($("#g_name").val());
            $("#licenceCiqName").val($("#ciq_name").val());
            table.render({
                elem: '#dec_goods_table'
                , toolbar: '#dec_goods_tool'
                , defaultToolbar: ['filter']
                , colFilterRecord: 'local'
                , cols: [[
                    { type: 'checkbox' }
                    , { field: 'goods_no', title: '序号', width: 80 }
                    , { field: 'lic_type_code', title: '许可证类别代码', width: 150 }
                    , { field: 'lic_type_name', title: '证书类别名称', width: 150 }
                    , { field: 'licence_no', title: '许可证编码', width: 130 }
                    , { field: 'lic_wrtof_detail_no', title: '核销货物序号', width: 150 }
                    , { field: 'lic_wrtof_qty', title: '核销数量', width: 110 }
                    , { field: 'lic_wrtof_qty_unit_name', title: '核销数量单位', width: 140 }
                ]]
                , data: quas
                , limit: quas.length
                , height: 350
            });
            $("#lic_type_name").focus();
            /**自动完成--许可证类别**/
            auto_fn({
                data: admin.all_complete_data.enterprise_product,
                listDirection: false,
                id: ['#lic_type_name'],
                after: ['#lic_type_code']
            });
            /**自动完成--核销数量单位**/
            auto_fn({
                data: admin.all_complete_data.unit_measurement,
                listDirection: false,
                id: ['#lic_wrtof_qty_unit_name'],
                after: ['#lic_wrtof_qty_unit']
            });
        });
        /**编辑产品许可证/审批/备案信息--点击行反填数据**/
        table.on('row(dec_goods_table)', function (obj) {
            $("#goodsLimitSeqNo").val(obj.data.goods_no);
            $("#lic_type_code").val(obj.data.lic_type_code);
            $("#lic_type_name").val(obj.data.lic_type_name);
            $("#licence_no").val(obj.data.licence_no);
            $("#lic_wrtof_detail_no").val(obj.data.lic_wrtof_detail_no);
            $("#lic_wrtof_qty").val(obj.data.lic_wrtof_qty);
            $("#lic_wrtof_qty_unit").val(obj.data.lic_wrtof_qty_unit);
            $("#lic_wrtof_qty_unit_name").val(obj.data.lic_wrtof_qty_unit_name);
            $("#lic_type_name").focus();
            dec_goods_index = obj.data.goods_no;
        });
        /**编辑产品许可证/审批/备案信息--保存**/
        form.on('submit(dec_goods_submit)', async (data) => {
            delete data.field.layTableCheckbox;
            if (dec_goods_index) {
                for (let item in data.field) {
                    quas[dec_goods_index - 1][item] = data.field[item]
                }
            } else {
                quas.push(data.field);
            }
            quas.forEach((value, index) => {
                value.goods_no = index + 1
            });
            table.reload('dec_goods_table', {
                data: quas,
                limit: quas.length
            });
            $("#goodsLimitSeqNo").val("");
            $("#lic_type_code").val("");
            $("#lic_type_name").val("");
            $("#licence_no").val("");
            $("#lic_wrtof_detail_no").val("");
            $("#lic_wrtof_qty").val("");
            $("#lic_wrtof_qty_unit").val("");
            $("#lic_wrtof_qty_unit_name").val("");
            $("#lic_type_name").focus();
            dec_goods_index = null;
        });
        /**编辑产品许可证/审批/备案信息--删除**/
        table.on('toolbar(dec_goods_table)', function (obj) {
            let checkStatus = table.checkStatus(obj.config.id);
            let checkData = checkStatus.data;
            switch (obj.event) {
                case 'add':
                    $("#goodsLimitSeqNo").val("");
                    $("#lic_type_code").val("");
                    $("#lic_type_name").val("");
                    $("#licence_no").val("");
                    $("#lic_wrtof_detail_no").val("");
                    $("#lic_wrtof_qty").val("");
                    $("#lic_wrtof_qty_unit").val("");
                    $("#lic_wrtof_qty_unit_name").val("");
                    $("#lic_type_name").focus();
                    dec_goods_index = null;
                    break;
                case 'delete':
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    layer.confirm('真的删除么', { title: '提示' }, async (index) => {
                        dec_goods_index = null;
                        let sup_data = quas.filter(item => {
                            return checkData.every(item2 => {
                                return item.goods_no != item2.goods_no;
                            })
                        });
                        for (let item_apply of checkData) {
                            if (item_apply.id) {
                                admin.order_quas_delete_ids.push({
                                    index: order_pros_index - 1,
                                    id: item_apply.id
                                });
                            }
                        }
                        quas = sup_data;
                        quas.forEach((value, index) => {
                            value.goods_no = index + 1
                        });
                        table.reload('dec_goods_table', {
                            data: quas,
                            limit: quas.length
                        });
                        $("#goodsLimitSeqNo").val("");
                        $("#lic_type_code").val("");
                        $("#lic_type_name").val("");
                        $("#licence_no").val("");
                        $("#lic_wrtof_detail_no").val("");
                        $("#lic_wrtof_qty").val("");
                        $("#lic_wrtof_qty_unit").val("");
                        $("#lic_wrtof_qty_unit_name").val("");
                        $("#lic_type_name").focus();
                        layer.close(index);
                    });
                    break;
            }
        });
        /**编辑产品许可证/审批/备案信息--按enter保存**/
        $("body").on("keyup", "#lic_wrtof_qty_unit_name", function (event) {
            if (event.keyCode == 13) {
                $("#dec_goods_submit").click();
                $("#lic_type_name").focus();
            }
        });
        /**许可证VIN信息**/
        let dec_goods_vin_index = null;
        form.on('submit(dec_goods_vin)', async (data) => {
            delete data.field.layTableCheckbox;
            if (dec_goods_index) {
                for (let item in data.field) {
                    quas[dec_goods_index - 1][item] = data.field[item]
                }
            } else {
                quas.push(data.field);
            }
            quas.forEach((value, index) => {
                value.goods_no = index + 1
            });
            table.reload('dec_goods_table', {
                data: quas,
                limit: quas.length
            });
            layer.open({
                type: 1,
                title: '编辑许可证VIN',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['1002px', '590px'],
                content: $('#dec_goods_vin_template').html()
            });
            let dec_goods_vin_data = quas[dec_goods_index - 1].order_pro_qua_vins || [];
            table.render({
                elem: '#dec_goods_vin_table'
                , toolbar: '#dec_goods_vin_tool'
                , defaultToolbar: ['filter']
                , colFilterRecord: 'local'
                , cols: [[
                    { type: 'checkbox' }
                    , { field: 'vin_no', title: 'VIN序号', width: 100 }
                    , { field: 'bill_lad_date', title: '提/运单日期', width: 140 }
                    , { field: 'quality_qgp', title: '质量保质期', width: 130 }
                    , { field: 'motor_no', title: '发动机号或电机号', width: 200 }
                    , { field: 'vin_code', title: '车辆识别代码(VIN)', width: 200 }
                    , { field: 'chassis_no', title: '底盘(车架)号', width: 200 }
                    , { field: 'invoice_no', title: '发票号', width: 130 }
                    , { field: 'invoice_num', title: '发票所列数量', width: 150 }
                    , { field: 'prod_cnnm', title: '品名(中文名称)', width: 220 }
                    , { field: 'prod_ennm', title: '品名(英文名称)', width: 220 }
                    , { field: 'model_en', title: '型号(英文名称)', width: 220 }
                    , { field: 'price_per_unit', title: '单价', width: 120 }
                ]]
                , data: dec_goods_vin_data
                , limit: dec_goods_vin_data.length
                , height: 350
            });
            laydate.render({ elem: '#bill_lad_date', theme: '#1E9FFF', fixed: true });
            $("#licTypeCodeVinName").val(quas[dec_goods_index - 1].lic_type_name);
            $("#licenceNo").val(quas[dec_goods_index - 1].licence_no);
            $("#vin_no").focus();
        });
        /**许可证VIN信息--点击行反填数据**/
        table.on('row(dec_goods_vin_table)', function (obj) {
            $("#goodsLimitVinSeqNo").val(obj.data.index);
            $("#vin_no").val(obj.data.vin_no);
            $("#bill_lad_date").val(obj.data.bill_lad_date);
            $("#quality_qgp").val(obj.data.quality_qgp);
            $("#motor_no").val(obj.data.motor_no);
            $("#vin_code").val(obj.data.vin_code);
            $("#chassis_no").val(obj.data.chassis_no);
            $("#invoice_no").val(obj.data.invoice_no);
            $("#invoice_num").val(obj.data.invoice_num);
            $("#prod_cnnm").val(obj.data.prod_cnnm);
            $("#prod_ennm").val(obj.data.prod_ennm);
            $("#model_en").val(obj.data.model_en);
            $("#price_per_unit").val(obj.data.price_per_unit);
            $("#vin_no").focus();
            dec_goods_vin_index = obj.data.index;
        });
        /**许可证VIN信息--保存**/
        form.on('submit(dec_goods_vin_submit)', async (data) => {
            delete data.field.layTableCheckbox;
            if (!(quas[dec_goods_index - 1].order_pro_qua_vins)) {
                quas[dec_goods_index - 1].order_pro_qua_vins = []
            }
            if (dec_goods_vin_index) {
                for (let item in data.field) {
                    quas[dec_goods_index - 1].order_pro_qua_vins[dec_goods_vin_index - 1][item] = data.field[item]
                }
            } else {
                quas[dec_goods_index - 1].order_pro_qua_vins.push(data.field);
            }
            quas[dec_goods_index - 1].order_pro_qua_vins.forEach((value, index) => {
                value.index = index + 1
            });
            table.reload('dec_goods_vin_table', {
                data: quas[dec_goods_index - 1].order_pro_qua_vins,
                limit: quas[dec_goods_index - 1].order_pro_qua_vins.length
            });
            $("#goodsLimitVinSeqNo").val("");
            $("#vin_no").val("");
            $("#bill_lad_date").val("");
            $("#quality_qgp").val("");
            $("#motor_no").val("");
            $("#vin_code").val("");
            $("#chassis_no").val("");
            $("#invoice_no").val("");
            $("#invoice_num").val("");
            $("#prod_cnnm").val("");
            $("#prod_ennm").val("");
            $("#model_en").val("");
            $("#price_per_unit").val("");
            $("#vin_no").focus();
            dec_goods_vin_index = null;
        });
        /**许可证VIN信息--删除**/
        table.on('toolbar(dec_goods_vin_table)', function (obj) {
            let checkStatus = table.checkStatus(obj.config.id);
            let checkData = checkStatus.data;
            switch (obj.event) {
                case 'add':
                    $("#goodsLimitVinSeqNo").val("");
                    $("#vin_no").val("");
                    $("#bill_lad_date").val("");
                    $("#quality_qgp").val("");
                    $("#motor_no").val("");
                    $("#vin_code").val("");
                    $("#chassis_no").val("");
                    $("#invoice_no").val("");
                    $("#invoice_num").val("");
                    $("#prod_cnnm").val("");
                    $("#prod_ennm").val("");
                    $("#model_en").val("");
                    $("#price_per_unit").val("");
                    $("#vin_no").focus();
                    dec_goods_vin_index = null;
                    break;
                case 'delete':
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    layer.confirm('真的删除么', { title: '提示' }, async (index) => {
                        dec_goods_vin_index = null;
                        let sup_data = quas[dec_goods_index - 1].order_pro_qua_vins.filter(item => {
                            return checkData.every(item2 => {
                                return item.index != item2.index;
                            })
                        });
                        for (let item_apply of checkData) {
                            if (item_apply.id) {
                                admin.order_quas_vin_delete_ids.push({
                                    pros_index: order_pros_index - 1,
                                    index: dec_goods_index - 1,
                                    id: item_apply.id
                                });
                            }
                        }
                        quas[dec_goods_index - 1].order_pro_qua_vins = sup_data;
                        quas[dec_goods_index - 1].order_pro_qua_vins.forEach((value, index) => {
                            value.index = index + 1
                        });
                        table.reload('dec_goods_vin_table', {
                            data: quas[dec_goods_index - 1].order_pro_qua_vins,
                            limit: quas[dec_goods_index - 1].order_pro_qua_vins.length
                        });
                        $("#goodsLimitVinSeqNo").val("");
                        $("#vin_no").val("");
                        $("#bill_lad_date").val("");
                        $("#quality_qgp").val("");
                        $("#motor_no").val("");
                        $("#vin_code").val("");
                        $("#chassis_no").val("");
                        $("#invoice_no").val("");
                        $("#invoice_num").val("");
                        $("#prod_cnnm").val("");
                        $("#prod_ennm").val("");
                        $("#model_en").val("");
                        $("#price_per_unit").val("");
                        $("#vin_no").focus();
                        layer.close(index);
                    });
                    break;
            }
        });
        /**许可证VIN信息--按enter保存**/
        $("body").on("keyup", "#price_per_unit", function (event) {
            if (event.keyCode == 13) {
                $("#dec_goods_vin_submit").click();
                $("#vin_no").focus();
            }
        });
        /**货物属性**/
        $("body").on("click", "#goods_attr_open", async function () {
            const data = await admin.post(`/clearance/no_paginate`,{Type:"货物属性代码"});
            const len = data.length;
            let result = [];
            const sliceNum = 4;
            for (let i = 0; i < len / sliceNum; i++) {
                result.push(data.slice(i * sliceNum, (i + 1) * sliceNum))
            }
            laytpl($("#goods_attr_data_template").html()).render(result, function (html) {
                $('#goods_attr_data_list').html(html)
            });
            layer.open({
                type: 1,
                title: '货物属性',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['680px', '300px'],
                content: $('#goods_attr_data_list').html(),
                success: function (layero, index) {
                    this.enterEsc = function (event) {
                        if (event.keyCode === 13) {
                            layer.closeAll();
                            return false;
                        }
                    };
                    $(document).on('keydown', this.enterEsc);
                },
                end: function () {
                    $(document).off('keydown', this.enterEsc);
                    $("#purpose_name").focus();
                }
            });
            $(".tableDiv_attr_td").each(function () {
                const result = goods_attr_data.some(item => {
                    if (item.code == $(this).data("id")) {
                        return true
                    }
                });
                if (result) {
                    $(this).addClass("bgcolor");
                }
            })
        });
        /**货物属性选择**/
        $("body").on("click", ".tableDiv_attr_td", function () {
            if ($(this).hasClass("bgcolor")) {
                $(this).removeClass("bgcolor")
            } else {
                $(this).addClass("bgcolor");
            }
            const result = goods_attr_data.some(item => {
                if (item.code == $(this).data("id")) {
                    return true
                }
            });
            if (result) {
                const index = goods_attr_data.findIndex(item => item.code == $(this).data("id"));
                goods_attr_data.splice(index, 1);
            } else {
                goods_attr_data.push({
                    code: $(this).data("id"),
                    name: $(this).attr("name")
                })
            }
            const data = goods_attr_data.map(item => item.name);
            $("#goods_attr").val(data.join(","));
        });
        /**编辑货物危险信息**/
        $('body').on('click', '#dangerBtn', async function () {
            layer.open({
                type: 1,
                title: '编辑货物危险信息',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['680px', '270px'],
                content: $('#dang_list').html()
            });
            /**自动完成--非危险化学品**/
            auto_fn({
                data: admin.chemicals_data,
                listDirection: false,
                id: ['#no_dang_flag_name'],
                after: ['#no_dang_flag']
            });
            /**自动完成--危包类别**/
            auto_fn({
                data: admin.category_data,
                listDirection: false,
                id: ['#dang_pack_type_name'],
                after: ['#dang_pack_type']
            });
            $("#no_dang_flag_name").focus();
            if (dang_data) {
                for (let item in dang_data) {
                    $(`#${item}`).val(dang_data[item])
                }
            }
        });
        /**保存货物危险信息**/
        form.on('submit(dang_save)', function (data) {
            let array = [];
            dang_data = data.field;
            for (let item in dang_data) {
                if (dang_data[item]) {
                    array.push(dang_data[item])
                }
            }
            layer.closeAll()
        });
        /**调用历史申报数据**/
        $("body").on("click", "#his_dec", async function () {
            his_dec_index = layer.open({
                type: 1,
                title: `历史申报数据${$("#code_t_s").val()}-${$("#g_name").val()}`,
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['910px', '730px'],
                content: $('#elements_his_list').html(),
                success: function (layero, index) {
                    this.enterEsc = function (event) {
                        if (event.keyCode === 13) {
                            $("#elements_his_save").click();
                            return false;
                        }
                    };
                    $(document).on('keydown', this.enterEsc);
                },
                end: function () {
                    $("#val0Name").focus();
                    $(document).off('keydown', this.enterEsc);
                }
            });
            let url;
            if ($("#manual_no").val()) {
                url = `/order/history_good_g_model?search=manual_no:${$("#manual_no").val()};trade_code:${$("#trade_code").val()};code_t_s:${$("#code_t_s").val()};i_e_flag:${cusIEFlag}`
            } else {
                url = `/order/history_good_g_model?search=trade_code:${$("#trade_code").val()};code_t_s:${$("#code_t_s").val()};i_e_flag:${cusIEFlag}`
            }
            admin.elements_his_data = await admin.get(url);
            table.render({
                elem: '#elements_his_list_table'
                , toolbar: true
                , defaultToolbar: ['filter']
                , colFilterRecord: 'local'
                , cols: [[
                    { type: 'radio' }
                    , {
                        field: 'apl_date', title: '申报日期', width: 160, templet: function (data) {
                            return data.order.apl_date
                        }
                    }
                    , { field: 'g_model', title: '申报要素' }
                ]]
                , data: admin.elements_his_data
                , limit: 10
                , page: true
                , height: 550
            });
            $(".layui-table-view[lay-id='elements_his_list_table'] .layui-table-body tr[data-index='0'] .layui-form-radio").click();
        });
        /**搜索历史申报要素**/
        $("body").on('input', '#elements_his_search', async function () {
            const text = $(this).val();
            const filter_data = admin.elements_his_data.filter(item => {
                return (item.g_model).indexOf(text) > -1
            });
            table.reload('elements_his_list_table', {
                data: filter_data
            });
            $(".layui-table-view[lay-id='elements_his_list_table'] .layui-table-body tr[data-index='0'] .layui-form-radio").click();
        });
        /**保存历史申报要素**/
        $('body').on('click', '#elements_his_save', function () {
            const data = table.checkStatus('elements_his_list_table').data;
            $("#elements").val(data[0].g_model);
            const arrAnge = data[0].g_model.split("|");
            arrAnge.forEach((value, index) => {
                $(`#val${index}`).val(value)
            });
            admin.brand_type_data.forEach((value, index) => {
                if (value.id == $("#val0").val()) {
                    $("#val0Name").val(value.value)
                }
            });
            admin.export_benefits_data.forEach((value, index) => {
                if (value.id == $("#val1").val()) {
                    $("#val1Name").val(value.value)
                }
            });
            layer.close(his_dec_index);
        });
        /**征免方式回车保存商品**/
        $("body").on("keyup", "#duty_mode_name", function (event) {
            edit_row = $('#g_no').val();
            const eCode = event.keyCode ? event.keyCode : event.which ? event.which : event.charCode;
            if (event.shiftKey != 1 && eCode == 13) {
                if ($("#decListShow").attr("isdeclistshow") == 0) {
                    $("#order_pros_submit").click();
                } else {
                    $("#goods_spec_open").click();
                }
            }
            $('.layui-table-main').scrollTop(38 * (edit_row - 1));
        });
        /**用途回车保存商品**/
        $("body").on("keyup", "#purpose_name", function (event) {
            edit_row = $('#g_no').val();
            const eCode = event.keyCode ? event.keyCode : event.which ? event.which : event.charCode;
            if (event.shiftKey != 1 && eCode == 13) {
                $("#order_pros_submit").click();
            }
            $('.layui-table-main').scrollTop(38 * (edit_row - 1));
        });
        /**保存商品**/
        let insert = false;
        let insert_index = null;
        form.on('submit(order_pros_submit)', function (data) {
            if ($('#g_no').val() > 50) {
                layer.msg("商品数量不能超过50条！");
                return
            }
            if (order_pros_data.length >= 51) {
                layer.msg("商品数量不能超过50条！");
                return
            }
            if (order_pros_data.length > 1) {
                if (data.field.trade_curr != order_pros_data[0].trade_curr) {
                    return layer.msg("币制不一致！");
                }
            }
            if (data.field.g_unit == data.field.first_unit) {
                if (data.field.g_qty != data.field.first_qty) {
                    return layer.msg("成交数量与法定第一数量不一致！");
                }
            }
            if (data.field.g_unit == data.field.second_unit) {
                if (data.field.g_qty != data.field.second_qty) {
                    return layer.msg("成交数量与法定第二数量不一致！");
                }
            }
            if (data.field.code_t_s.toString().length != 10) {
                return layer.msg("商品编号必须为10位！");
            }
            data.field.decl_price = admin.formatnumber(data.field.decl_price, 4);
            data.field.decl_total = admin.formatnumber(data.field.decl_total, 4);

            /**
             * foreach data array
             *
             * add item+'_string'
             *
             * **/
            for (let item in data.field) {
                if ($.inArray(item, admin.order_pros_arr) >= 0) {
                    data.field[item + '_string'] = data.field[item]
                }
            }

            /** insert */
            if (insert) {
                data.field.goods_spec_data = goods_spec_data;
                data.field.quas = quas;
                data.field.goods_attr_data = goods_attr_data;
                data.field.dang_data = dang_data;
                order_pros_data.splice(insert_index - 1, 0, data.field);
                insert = true;
            } else {

                /** edit */
                if (order_pros_index) {
                    for (let item in data.field) {
                        order_pros_data[order_pros_index - 1][item] = data.field[item]
                    }
                    order_pros_data[order_pros_index - 1].goods_spec_data = goods_spec_data;
                    order_pros_data[order_pros_index - 1].quas = quas;
                    order_pros_data[order_pros_index - 1].goods_attr_data = goods_attr_data;
                    order_pros_data[order_pros_index - 1].dang_data = dang_data;
                } else {
                    /** create */
                    data.field.goods_spec_data = goods_spec_data;
                    data.field.quas = quas;
                    data.field.goods_attr_data = goods_attr_data;
                    data.field.dang_data = dang_data;
                    order_pros_data.push(data.field);
                }
            }

            order_pros_data.forEach((value, index) => {
                value.g_no = index + 1
            });
            table.reload('order_pros', {
                data: order_pros_data,
                limit: order_pros_data.length
            });
            order_pros_index = null;
            $("#order_pros_form input").each(function () {
                if ($(this).attr("id") == "g_no") {
                    $(this).val(order_pros_data.length + 1)
                } else {
                    if (!($(this).is("[no_empty]"))) {
                        $(this).val("")
                    }
                }
            });
            goods_spec_data = {};
            quas = [];
            goods_attr_data = [];
            dang_data = {};
            if ($("#contr_item").is(":disabled")) {
                $("#code_t_s").focus()
            } else {
                $("#contr_item").focus()
            }

            /**tip:总价/成交数量合计/法定第一数量合计/法定第二数量合计**/
            admin.is_total_number(order_pros_data);
        });
        /**商品--点击行反填数据**/
        table.on('row(order_pros)', function (obj) {

            for (let item in obj.data) {
                if ($.inArray(item, admin.order_pros_arr) >= 0) {
                    $(`#${item}`).val(obj.data[item + '_string'])
                } else {
                    $(`#${item}`).val(obj.data[item])
                }
            }
            goods_spec_data = obj.data.goods_spec_data;
            quas = obj.data.quas;
            goods_attr_data = obj.data.goods_attr_data;
            dang_data = obj.data.dang_data;
            let arr = [];
            if (goods_attr_data) {
                const data = goods_attr_data.map(item => item.name);
                $("#goods_attr").val(data.join(","));
            }
            if (goods_spec_data) {
                for (let item in goods_spec_data) {
                    if (goods_spec_data[item]) {
                        arr.push(goods_spec_data[item])
                    }
                }
                $("#goods_spec_data").val(arr.join(";"));
            }
            order_pros_index = obj.data.g_no;
            insert_index = obj.data.g_no;
            declaration_index = order_pros_index - 1;
        });
        /**商品列表--按钮操作**/
        table.on('toolbar(order_pros)', async function (obj) {
            let checkStatus = table.checkStatus(obj.config.id);
            let checkData = checkStatus.data;
            switch (obj.event) {
                case 'add':
                    insert = false;
                    order_pros_index = null;
                    declaration_index = null;
                    $("#order_pros_form input").each(function () {
                        if ($(this).attr("id") == "g_no") {
                            $(this).val(order_pros_data.length + 1)
                        } else {
                            if (!($(this).is("[no_empty]"))) {
                                $(this).val("")
                            }
                        }
                    });
                    goods_spec_data = {};
                    quas = [];
                    goods_attr_data = [];
                    dang_data = {};
                    if ($("#contr_item").is(":disabled")) {
                        $("#code_t_s").focus()
                    } else {
                        $("#contr_item").focus()
                    }
                    break;
                case 'delete':
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    layer.confirm('真的删除么', { title: '提示' }, async (index) => {
                        let sup_data = order_pros_data.filter(item => {
                            return checkData.every(item2 => {
                                return item.g_no != item2.g_no;
                            })
                        });
                        for (let item_apply of checkData) {
                            if (item_apply.id) {
                                admin.order_pros_delete_ids.push(item_apply.id);
                            }
                        }
                        order_pros_data = sup_data;
                        order_pros_data.forEach((value, index) => {
                            value.g_no = index + 1
                        });
                        table.reload('order_pros', {
                            data: order_pros_data,
                            limit: order_pros_data.length
                        });
                        insert = false;
                        order_pros_index = null;
                        declaration_index = null;
                        $("#order_pros_form input").each(function () {
                            if ($(this).attr("id") == "g_no") {
                                $(this).val(order_pros_data.length + 1)
                            } else {
                                if (!($(this).is("[no_empty]"))) {
                                    $(this).val("")
                                }
                            }
                        });
                        goods_spec_data = {};
                        quas = [];
                        goods_attr_data = [];
                        dang_data = {};
                        if ($("#contr_item").is(":disabled")) {
                            $("#code_t_s").focus()
                        } else {
                            $("#contr_item").focus()
                        }
                        /**tip:总价/成交数量合计/法定第一数量合计/法定第二数量合计**/
                        admin.is_total_number(order_pros_data);
                        layer.close(index);
                    });
                    break;
                case 'copy':
                    insert = false;
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    if (checkData.length != 1) {
                        return layer.msg('只能选择一条数据复制');
                    }
                    if (order_pros_data.length >= 50) {
                        return layer.msg("商品数量不能超过50条！");
                    }
                    checkData[0].g_no = order_pros_data.length + 1;
                    if (checkData[0].id) {
                        delete checkData[0].id
                    }
                    order_pros_data.push(checkData[0]);
                    table.reload('order_pros', {
                        data: order_pros_data,
                        limit: order_pros_data.length
                    });
                    $(`.layui-table-view[lay-id='order_pros'] .layui-table-body tr input`).each(function (index, el) {
                        el.checked = false;
                    });
                    form.render('checkbox');
                    $(`.layui-table-view[lay-id='order_pros'] .layui-table-body tr[data-index=${order_pros_data.length - 1}]`).click();
                    /**tip:总价/成交数量合计/法定第一数量合计/法定第二数量合计**/
                    admin.is_total_number(order_pros_data);
                    break;
                case 'up':
                    insert = false;
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    if (checkData.length != 1) {
                        return layer.msg('只能选择一条数据移动');
                    }
                    if (checkData[0].g_no == 1) {
                        return layer.msg("已是第一条商品表体，无法执行上移操作。");
                    }
                    const upIndex = parseInt(checkData[0].g_no) - 1;
                    order_pros_data = admin.swapItems(order_pros_data, upIndex, upIndex - 1);
                    order_pros_data.forEach((value, index) => {
                        value.g_no = index + 1
                    });
                    table.reload('order_pros', {
                        data: order_pros_data,
                        limit: order_pros_data.length
                    });
                    break;
                case 'down':
                    insert = false;
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    if (checkData.length != 1) {
                        return layer.msg('只能选择一条数据移动');
                    }
                    if (checkData[0].g_no == order_pros_data.length) {
                        return layer.msg("已是最后一条数据，无法执行下移移操作。");
                    }
                    const downIndex = parseInt(checkData[0].g_no) - 1;
                    order_pros_data = admin.swapItems(order_pros_data, downIndex, downIndex + 1);
                    order_pros_data.forEach((value, index) => {
                        value.g_no = index + 1
                    });
                    table.reload('order_pros', {
                        data: order_pros_data,
                        limit: order_pros_data.length
                    });
                    break;
                case 'insert':
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    if (checkData.length != 1) {
                        return layer.msg('只能选择一条数据插入');
                    }
                    insert = true;
                    order_pros_index = null;
                    declaration_index = null;
                    $("#order_pros_form input").each(function () {
                        if ($(this).attr("id") == "g_no") {
                            $(this).val(order_pros_data.length + 1)
                        } else {
                            if (!($(this).is("[no_empty]"))) {
                                $(this).val("")
                            }
                        }
                    });
                    goods_spec_data = {};
                    quas = [];
                    goods_attr_data = [];
                    dang_data = {};
                    if ($("#contr_item").is(":disabled")) {
                        $("#code_t_s").focus()
                    } else {
                        $("#contr_item").focus()
                    }
                    break;
                case 'again':
                    if (!($("#code_t_s").val().trim())) {
                        return layer.msg('请先填写商品编码');
                    }
                    const hs_code_data = await admin.get(`/hs_code/lists?limit=0&search=${order_pros_data[declaration_index].code_t_s}`);
                    if (hs_code_data.data.length === 0) {
                        return layer.msg("此商品没有申报要素")
                    }
                    declaration_edit = true;
                    const declaration_data_array = hs_code_data.data[0].declaration.split(";");
                    const data = declaration_data_array.map((item, index) => {
                        if (index < 9) {
                            return item.substr(1)
                        } else {
                            return item.substr(2)
                        }
                    });
                    data.type = 'again';
                    layui.laytpl($("#declaration_template").html()).render(data, function (html) {
                        $('#declaration_list').html(html)
                    });
                    layer.open({
                        type: 1,
                        title: '商品规范申报-商品申报要素',
                        shadeClose: true,
                        area: admin.screen() < 2 ? ['80%', '300px'] : ['1180px', '600px'],
                        content: `<div id="declaration_list_reload">${$('#declaration_list').html()}</div>`,
                        end: function () {
                            $(`.layui-table-view[lay-id='order_pros'] .layui-table-body tr input`).each(function (index, el) {
                                el.checked = false;
                            });
                            form.render('checkbox');
                            $(`.layui-table-view[lay-id='order_pros'] .layui-table-body tr[data-index=${declaration_index}]`).click();
                            $("#g_qty").focus();
                            declaration_edit = false;
                            $(`.layui-table-view[lay-id='order_pros'] .layui-table-body`).animate({
                                scrollTop: $(`.layui-table-view[lay-id='order_pros'] .layui-table-body tr[data-index=${declaration_index}]`)[0].offsetTop
                            }, 0);
                        }
                    });
                    $(".manual_no_btn_next").show();
                    $("body #val0Name").focus();
                    $("#selectCodeTs").val(order_pros_data[declaration_index].g_name);
                    $("#good_information").val(`${order_pros_data[declaration_index].g_no}-${order_pros_data[declaration_index].contr_item ? order_pros_data[declaration_index].contr_item : '一般贸易'}-${order_pros_data[declaration_index].g_name}`);
                    let brand_type = await admin.post(`/clearance/no_paginate`,{Type:"品牌类型"});
                    let data_filter = [];
                    for (let item of brand_type) {
                        data_filter.push({
                            id: item.customs_code,
                            label: `${item.customs_code}-${item.name}`,
                            value: `${item.name}`
                        })
                    }
                    admin.brand_type_data = data_filter;
                    $("#val0Name").AutoComplete({
                        'data': data_filter,
                        'width': 280,
                        'itemHeight': 20,
                        'listStyle': 'custom',
                        'listDirection': 'down',
                        'createItemHandler': function (index, data) {
                            return `<p class="auto_list_p">${data.label}</p>`
                        },
                        'afterSelectedHandler': function (data) {
                            $("#val0").val(data.id);
                        }
                    });

                    let export_benefits = await admin.post(`/clearance/no_paginate`,{Type:"出口享惠情况"});
                    let data_filter_benefits = [];
                    for (let item of export_benefits) {
                        data_filter_benefits.push({
                            id: item.customs_code,
                            label: `${item.customs_code}-${item.name}`,
                            value: `${item.name}`
                        })
                    }
                    admin.export_benefits_data = data_filter_benefits;
                    $("#val1Name").AutoComplete({
                        'data': data_filter_benefits,
                        'width': 280,
                        'itemHeight': 20,
                        'listStyle': 'custom',
                        'listDirection': 'down',
                        'createItemHandler': function (index, data) {
                            return `<p class="auto_list_p">${data.label}</p>`
                        },
                        'afterSelectedHandler': function (data) {
                            $("#val1").val(data.id);
                        }
                    });
                    $("#elements").val(order_pros_data[declaration_index].g_model);
                    const arrAnge = order_pros_data[declaration_index].g_model.split("|");
                    arrAnge.forEach((value, index) => {
                        $(`#val${index}`).val(value)
                    });
                    admin.brand_type_data.forEach((value, index) => {
                        if (value.id == $("#val0").val()) {
                            $("#val0Name").val(value.value)
                        }
                    });
                    admin.export_benefits_data.forEach((value, index) => {
                        if (value.id == $("#val1").val()) {
                            $("#val1Name").val(value.value)
                        }
                    });
                    break;
                case 'look':
                    if (!($("#code_t_s").val().trim())) {
                        return layer.msg('请先填写商品编码');
                    }
                    const hs_code_data_show = await admin.get(`/hs_code/lists?limit=0&search=${order_pros_data[declaration_index].code_t_s}`);
                    if (hs_code_data_show.data.length === 0) {
                        return layer.msg("此商品没有申报要素")
                    }
                    declaration_edit = false;
                    const declaration_data_array_show = hs_code_data_show.data[0].declaration.split(";");
                    const data_show = declaration_data_array_show.map((item, index) => {
                        if (index < 9) {
                            return item.substr(1)
                        } else {
                            return item.substr(2)
                        }
                    });
                    layui.laytpl($("#declaration_template").html()).render(data_show, function (html) {
                        $('#declaration_list').html(html)
                    });
                    layer.open({
                        type: 1,
                        title: '商品规范申报-商品申报要素',
                        shadeClose: true,
                        area: admin.screen() < 2 ? ['80%', '300px'] : ['1180px', '600px'],
                        content: `<div id="declaration_list_reload">${$('#declaration_list').html()}</div>`,
                        end: function () {
                            $("#g_qty").focus()
                        }
                    });
                    $(".manual_no_btn_next").show();
                    $("#declaration_save").hide();
                    $(".order_table_form_show input").each(function () {
                        $(this).attr("disabled", "disabled")
                    });
                    $("#selectCodeTs").val(order_pros_data[declaration_index].g_name);
                    $("#good_information").val(`${order_pros_data[declaration_index].g_no}-${order_pros_data[declaration_index].contr_item ? order_pros_data[declaration_index].contr_item : '一般贸易'}-${order_pros_data[declaration_index].g_name}`);
                    let brand_type_show = await admin.post(`/clearance/no_paginate`,{Type:"品牌类型"});
                    let data_filter_show = [];
                    for (let item of brand_type_show) {
                        data_filter_show.push({
                            id: item.customs_code,
                            label: `${item.customs_code}-${item.name}`,
                            value: `${item.name}`
                        })
                    }
                    admin.brand_type_data = data_filter_show;
                    $("#val0Name").AutoComplete({
                        'data': data_filter_show,
                        'width': 280,
                        'itemHeight': 20,
                        'listStyle': 'custom',
                        'listDirection': 'down',
                        'createItemHandler': function (index, data) {
                            return `<p class="auto_list_p">${data.label}</p>`
                        },
                        'afterSelectedHandler': function (data) {
                            $("#val0").val(data.id);
                        }
                    });

                    let export_benefits_show = await admin.post(`/clearance/no_paginate`,{Type:"出口享惠情况"});
                    let data_filter_benefits_show = [];
                    for (let item of export_benefits_show) {
                        data_filter_benefits_show.push({
                            id: item.customs_code,
                            label: `${item.customs_code}-${item.name}`,
                            value: `${item.name}`
                        })
                    }
                    admin.export_benefits_data = data_filter_benefits_show;
                    $("#val1Name").AutoComplete({
                        'data': data_filter_benefits_show,
                        'width': 280,
                        'itemHeight': 20,
                        'listStyle': 'custom',
                        'listDirection': 'down',
                        'createItemHandler': function (index, data) {
                            return `<p class="auto_list_p">${data.label}</p>`
                        },
                        'afterSelectedHandler': function (data) {
                            $("#val1").val(data.id);
                        }
                    });
                    $("#elements").val(order_pros_data[declaration_index].g_model);
                    const arrAnge_show = order_pros_data[declaration_index].g_model.split("|");
                    arrAnge_show.forEach((value, index) => {
                        $(`#val${index}`).val(value)
                    });
                    admin.brand_type_data.forEach((value, index) => {
                        if (value.id == $("#val0").val()) {
                            $("#val0Name").val(value.value)
                        }
                    });
                    admin.export_benefits_data.forEach((value, index) => {
                        if (value.id == $("#val1").val()) {
                            $("#val1Name").val(value.value)
                        }
                    });
                    break;
                case 'batch_edit':
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    layer.open({
                        type: 1,
                        title: '商品批量修改',
                        shadeClose: true,
                        area: admin.screen() < 2 ? ['80%', '300px'] : ['680px', '500px'],
                        content: $('#batch_edit_list').html()
                    });
                    $("#trade_curr_name_batch").focus();
                    laydate.render({ elem: '#prod_valid_dt_batch', theme: '#1E9FFF', fixed: true });
                    laydate.render({ elem: '#produce_date_batch', theme: '#1E9FFF', fixed: true });
                    /**自动完成--批量修改--币制**/
                    auto_fn({
                        data: admin.all_complete_data.currency,
                        listDirection: false,
                        id: ['#trade_curr_name_batch'],
                        after: ['#trade_curr_batch']
                    });
                    /**自动完成--批量修改--原产国（地区）**/
                    auto_fn({
                        data: admin.all_complete_data.country_area,
                        listDirection: false,
                        id: ['#origin_country_name_batch'],
                        after: ['#origin_country_batch']
                    });
                    /**自动完成--批量修改--境内目的地代码**/
                    auto_fn({
                        data: admin.all_complete_data.domestic_area,
                        listDirection: false,
                        id: ['#district_code_name_batch'],
                        after: ['#district_code_batch']
                    });
                    /**自动完成--征免方式**/
                    auto_fn({
                        data: admin.all_complete_data.exempting_method,
                        listDirection: false,
                        id: ['#duty_mode_name_batch'],
                        after: ['#duty_mode_batch']
                    });
                    break;
            }
        });
        /**查看商品申报要素-上一条**/
        /*$("body").on("click", "#declaration_last_look", async function () {
            if (declaration_index == 0) {
                return layer.msg("已经是第一条了");
            }
            declaration_index -= 1;
            const hs_code_data_show = await admin.get(`/hs_code/lists?limit=0&search=${order_pros_data[declaration_index].code_t_s}`);
            const declaration_data_array_show = hs_code_data_show.data[0].declaration.split(";");
            const data_show = declaration_data_array_show.map((item, index) => {
                if (index < 9) {
                    return item.substr(1)
                } else {
                    return item.substr(2)
                }
            });
            layui.laytpl($("#declaration_template").html()).render(data_show, function (html) {
                $('#declaration_list_reload').html(html)
            });
            $(".manual_no_btn_next").show();
            $(".order_table_form_show input").each(function () {
                if(!declaration_edit) {
                    $(this).attr("disabled", "disabled");
                }
                $(this).val("");
            });
            $("#selectCodeTs").val(order_pros_data[declaration_index].g_name);
            $("#good_information").val(`${order_pros_data[declaration_index].g_no}-${order_pros_data[declaration_index].contr_item ? order_pros_data[declaration_index].contr_item : '一般贸易'}-${order_pros_data[declaration_index].g_name}`);
            let brand_type_show = await admin.post(`/clearance/no_paginate`,{Type:"品牌类型"});
            let data_filter_show = [];
            for (let item of brand_type_show) {
                data_filter_show.push({
                    id: item.customs_code,
                    label: `${item.customs_code}-${item.name}`,
                    value: `${item.name}`
                })
            }
            admin.brand_type_data = data_filter_show;
            $("#val0Name").AutoComplete({
                'data': data_filter_show,
                'width': 280,
                'itemHeight': 20,
                'listStyle': 'custom',
                'listDirection': 'down',
                'createItemHandler': function (index, data) {
                    return `<p class="auto_list_p">${data.label}</p>`
                },
                'afterSelectedHandler': function (data) {
                    $("#val0").val(data.id);
                }
            });

            let export_benefits_show = await admin.post(`/clearance/no_paginate`,{Type:"出口享惠情况"});
            let data_filter_benefits_show = [];
            for (let item of export_benefits_show) {
                data_filter_benefits_show.push({
                    id: item.customs_code,
                    label: `${item.customs_code}-${item.name}`,
                    value: `${item.name}`
                })
            }
            admin.export_benefits_data = data_filter_benefits_show;
            $("#val1Name").AutoComplete({
                'data': data_filter_benefits_show,
                'width': 280,
                'itemHeight': 20,
                'listStyle': 'custom',
                'listDirection': 'down',
                'createItemHandler': function (index, data) {
                    return `<p class="auto_list_p">${data.label}</p>`
                },
                'afterSelectedHandler': function (data) {
                    $("#val1").val(data.id);
                }
            });
            $("#elements").val(order_pros_data[declaration_index].g_model);
            const arrAnge_show = order_pros_data[declaration_index].g_model.split("|");
            arrAnge_show.forEach((value, index) => {
                $(`#val${index}`).val(value)
            });
            admin.brand_type_data.forEach((value, index) => {
                if (value.id == $("#val0").val()) {
                    $("#val0Name").val(value.value)
                }
            });
            admin.export_benefits_data.forEach((value, index) => {
                if (value.id == $("#val1").val()) {
                    $("#val1Name").val(value.value)
                }
            });
        });*/


        async function last_dec(type) {
            if (declaration_index == 0) {
                return layer.msg("已经是第一条了");
            }
            declaration_index -= 1;
            const hs_code_data_show = await admin.get(`/hs_code/lists?limit=0&search=${order_pros_data[declaration_index].code_t_s}`);
            const declaration_data_array_show = hs_code_data_show.data[0].declaration.split(";");
            const data_show = declaration_data_array_show.map((item, index) => {
                if (index < 9) {
                    return item.substr(1)
                } else {
                    return item.substr(2)
                }
            });
            layui.laytpl($("#declaration_template").html()).render(data_show, function (html) {
                $('#declaration_list_reload').html(html)
            });

            if(type == 'look'){
                $('#declaration_next').attr('id','declaration_next_look');
                $('#declaration_last').attr('id','declaration_last_look');
            }
            submit_btn();
            form.render();

            $('#elements').focus();
            $(".manual_no_btn_next").show();
            $(".order_table_form_show input").each(function () {
                if (!declaration_edit) {
                    $(this).attr("disabled", "disabled");
                }
                $(this).val("");
            });
            $("#selectCodeTs").val(order_pros_data[declaration_index].g_name);
            $("#good_information").val(`${order_pros_data[declaration_index].g_no}-${order_pros_data[declaration_index].contr_item ? order_pros_data[declaration_index].contr_item : '一般贸易'}-${order_pros_data[declaration_index].g_name}`);
            let brand_type_show = await admin.post(`/clearance/no_paginate`,{Type:"品牌类型"});
            let data_filter_show = [];
            for (let item of brand_type_show) {
                data_filter_show.push({
                    id: item.customs_code,
                    label: `${item.customs_code}-${item.name}`,
                    value: `${item.name}`
                })
            }
            admin.brand_type_data = data_filter_show;
            $("#val0Name").AutoComplete({
                'data': data_filter_show,
                'width': 280,
                'itemHeight': 20,
                'listStyle': 'custom',
                'listDirection': 'down',
                'createItemHandler': function (index, data) {
                    return `<p class="auto_list_p">${data.label}</p>`
                },
                'afterSelectedHandler': function (data) {
                    $("#val0").val(data.id);
                }
            });

            let export_benefits_show = await admin.post(`/clearance/no_paginate`,{Type:"出口享惠情况"});
            let data_filter_benefits_show = [];
            for (let item of export_benefits_show) {
                data_filter_benefits_show.push({
                    id: item.customs_code,
                    label: `${item.customs_code}-${item.name}`,
                    value: `${item.name}`
                })
            }
            admin.export_benefits_data = data_filter_benefits_show;
            $("#val1Name").AutoComplete({
                'data': data_filter_benefits_show,
                'width': 280,
                'itemHeight': 20,
                'listStyle': 'custom',
                'listDirection': 'down',
                'createItemHandler': function (index, data) {
                    return `<p class="auto_list_p">${data.label}</p>`
                },
                'afterSelectedHandler': function (data) {
                    $("#val1").val(data.id);
                }
            });
            if(!order_pros_data[declaration_index].g_model){
                order_pros_data[declaration_index].g_model='新';
            }
            $("#elements").val(order_pros_data[declaration_index].g_model);
            const arrAnge_show = order_pros_data[declaration_index].g_model.split("|");
            arrAnge_show.forEach((value, index) => {
                $(`#val${index}`).val(value)
            });
            admin.brand_type_data.forEach((value, index) => {
                if (value.id == $("#val0").val()) {
                    $("#val0Name").val(value.value)
                }
            });
            admin.export_benefits_data.forEach((value, index) => {
                if (value.id == $("#val1").val()) {
                    $("#val1Name").val(value.value)
                }
            });
            $("#val1Name").val('不适用于进口报关单');
            $(`#val1`).val(3);
        }


        /**查看商品申报要素-下一条**/
        /*
        $("body").on("click", "#declaration_next_look", async function () {
            if ((declaration_index + 1) == order_pros_data.length) {
                return layer.msg("已经是最后一条了");
            }
            declaration_index += 1;
            const hs_code_data_show = await admin.get(`/hs_code/lists?limit=0&search=${order_pros_data[declaration_index].code_t_s}`);
            const declaration_data_array_show = hs_code_data_show.data[0].declaration.split(";");
            const data_show = declaration_data_array_show.map((item, index) => {
                if (index < 9) {
                    return item.substr(1)
                } else {
                    return item.substr(2)
                }
            });
            layui.laytpl($("#declaration_template").html()).render(data_show, function (html) {
                $('#declaration_list_reload').html(html)
            });
            $(".manual_no_btn_next").show();
            $(".order_table_form_show input").each(function () {
                if(!declaration_edit) {
                    $(this).attr("disabled", "disabled");
                }
                $(this).val("");
            });
            $("#selectCodeTs").val(order_pros_data[declaration_index].g_name);
            $("#good_information").val(`${order_pros_data[declaration_index].g_no}-${order_pros_data[declaration_index].contr_item ? order_pros_data[declaration_index].contr_item : '一般贸易'}-${order_pros_data[declaration_index].g_name}`);
            let brand_type_show = await admin.post(`/clearance/no_paginate`,{Type:"品牌类型"});
            let data_filter_show = [];
            for (let item of brand_type_show) {
                data_filter_show.push({
                    id: item.customs_code,
                    label: `${item.customs_code}-${item.name}`,
                    value: `${item.name}`
                })
            }
            admin.brand_type_data = data_filter_show;
            $("#val0Name").AutoComplete({
                'data': data_filter_show,
                'width': 280,
                'itemHeight': 20,
                'listStyle': 'custom',
                'listDirection': 'down',
                'createItemHandler': function (index, data) {
                    return `<p class="auto_list_p">${data.label}</p>`
                },
                'afterSelectedHandler': function (data) {
                    $("#val0").val(data.id);
                }
            });

            let export_benefits_show = await admin.post(`/clearance/no_paginate`,{Type:"出口享惠情况"});
            let data_filter_benefits_show = [];
            for (let item of export_benefits_show) {
                data_filter_benefits_show.push({
                    id: item.customs_code,
                    label: `${item.customs_code}-${item.name}`,
                    value: `${item.name}`
                })
            }
            admin.export_benefits_data = data_filter_benefits_show;
            $("#val1Name").AutoComplete({
                'data': data_filter_benefits_show,
                'width': 280,
                'itemHeight': 20,
                'listStyle': 'custom',
                'listDirection': 'down',
                'createItemHandler': function (index, data) {
                    return `<p class="auto_list_p">${data.label}</p>`
                },
                'afterSelectedHandler': function (data) {
                    $("#val1").val(data.id);
                }
            });
            $("#elements").val(order_pros_data[declaration_index].g_model);
            const arrAnge_show = order_pros_data[declaration_index].g_model.split("|");
            arrAnge_show.forEach((value, index) => {
                $(`#val${index}`).val(value)
            });
            admin.brand_type_data.forEach((value, index) => {
                if (value.id == $("#val0").val()) {
                    $("#val0Name").val(value.value)
                }
            });
            admin.export_benefits_data.forEach((value, index) => {
                if (value.id == $("#val1").val()) {
                    $("#val1Name").val(value.value)
                }
            });
        });
        */

        async function next_dec(type) {
            if ((declaration_index + 1) == order_pros_data.length) {
                return layer.msg("已经是最后一条了");
            }
            declaration_index += 1;
            const hs_code_data_show = await admin.get(`/hs_code/lists?limit=0&search=${order_pros_data[declaration_index].code_t_s}`);
            const declaration_data_array_show = hs_code_data_show.data[0].declaration.split(";");
            const data_show = declaration_data_array_show.map((item, index) => {
                if (index < 9) {
                    return item.substr(1)
                } else {
                    return item.substr(2)
                }
            });
            layui.laytpl($("#declaration_template").html()).render(data_show, function (html) {
                $('#declaration_list_reload').html(html)
            });

            if(type == 'look'){
                $('#declaration_next').attr('id','declaration_next_look');
                $('#declaration_last').attr('id','declaration_last_look');
            }
            submit_btn();
            form.render();


            $('#elements').focus();
            $(".manual_no_btn_next").show();
            $(".order_table_form_show input").each(function () {
                if (!declaration_edit) {
                    $(this).attr("disabled", "disabled");
                }
                $(this).val("");
            });
            $("#selectCodeTs").val(order_pros_data[declaration_index].g_name);
            $("#good_information").val(`${order_pros_data[declaration_index].g_no}-${order_pros_data[declaration_index].contr_item ? order_pros_data[declaration_index].contr_item : '一般贸易'}-${order_pros_data[declaration_index].g_name}`);
            let brand_type_show = await admin.post(`/clearance/no_paginate`,{Type:"品牌类型"});
            let data_filter_show = [];
            for (let item of brand_type_show) {
                data_filter_show.push({
                    id: item.customs_code,
                    label: `${item.customs_code}-${item.name}`,
                    value: `${item.name}`
                })
            }
            admin.brand_type_data = data_filter_show;
            $("#val0Name").AutoComplete({
                'data': data_filter_show,
                'width': 280,
                'itemHeight': 20,
                'listStyle': 'custom',
                'listDirection': 'down',
                'createItemHandler': function (index, data) {
                    return `<p class="auto_list_p">${data.label}</p>`
                },
                'afterSelectedHandler': function (data) {
                    $("#val0").val(data.id);
                }
            });

            let export_benefits_show = await admin.post(`/clearance/no_paginate`,{Type:"出口享惠情况"});
            let data_filter_benefits_show = [];
            for (let item of export_benefits_show) {
                data_filter_benefits_show.push({
                    id: item.customs_code,
                    label: `${item.customs_code}-${item.name}`,
                    value: `${item.name}`
                })
            }
            admin.export_benefits_data = data_filter_benefits_show;
            $("#val1Name").AutoComplete({
                'data': data_filter_benefits_show,
                'width': 280,
                'itemHeight': 20,
                'listStyle': 'custom',
                'listDirection': 'down',
                'createItemHandler': function (index, data) {
                    return `<p class="auto_list_p">${data.label}</p>`
                },
                'afterSelectedHandler': function (data) {
                    $("#val1").val(data.id);
                }
            });
            if(!order_pros_data[declaration_index].g_model){
                order_pros_data[declaration_index].g_model='新';
            }
            $("#elements").val(order_pros_data[declaration_index].g_model);
            const arrAnge_show = order_pros_data[declaration_index].g_model.split("|");
            arrAnge_show.forEach((value, index) => {
                $(`#val${index}`).val(value)
            });
            admin.brand_type_data.forEach((value, index) => {
                if (value.id == $("#val0").val()) {
                    $("#val0Name").val(value.value)
                }
            });
            admin.export_benefits_data.forEach((value, index) => {
                if (value.id == $("#val1").val()) {
                    $("#val1Name").val(value.value)
                }
            });
            $("#val1Name").val('不适用于进口报关单');
            $(`#val1`).val(3);
        }
        /**批量修改--货物属性**/
        $("body").on("click", "#goods_attr_open_batch", async function () {
            const data = await admin.post(`/clearance/no_paginate`,{Type:"货物属性代码"});
            const len = data.length;
            let result = [];
            const sliceNum = 4;
            for (let i = 0; i < len / sliceNum; i++) {
                result.push(data.slice(i * sliceNum, (i + 1) * sliceNum))
            }
            laytpl($("#goods_attr_batch_data_template").html()).render(result, function (html) {
                $('#goods_attr_batch_data_list').html(html)
            });
            const goods_attr_batch_index = layer.open({
                type: 1,
                title: '货物属性',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['680px', '300px'],
                content: $('#goods_attr_batch_data_list').html(),
                success: function (layero, index) {
                    this.enterEsc = function (event) {
                        if (event.keyCode === 13) {
                            layer.close(goods_attr_batch_index);
                            return false;
                        }
                    };
                    $(document).on('keydown', this.enterEsc);
                },
                end: function () {
                    $(document).off('keydown', this.enterEsc);
                    $("#stuff_batch").focus();
                }
            });
            $(".tableDiv_attr_batch_td").each(function () {
                const result = goods_attr_batch_data.some(item => {
                    if (item.code == $(this).data("id")) {
                        return true
                    }
                });
                if (result) {
                    $(this).addClass("bgcolor");
                }
            })
        });
        /**批量修改--货物属性选择**/
        $("body").on("click", ".tableDiv_attr_batch_td", function () {
            if ($(this).hasClass("bgcolor")) {
                $(this).removeClass("bgcolor")
            } else {
                $(this).addClass("bgcolor");
            }
            const result = goods_attr_batch_data.some(item => {
                if (item.code == $(this).data("id")) {
                    return true
                }
            });
            if (result) {
                const index = goods_attr_batch_data.findIndex(item => item.code == $(this).data("id"));
                goods_attr_batch_data.splice(index, 1);
            } else {
                goods_attr_batch_data.push({
                    code: $(this).data("id"),
                    name: $(this).attr("name")
                })
            }
            const data = goods_attr_batch_data.map(item => item.name);
            $("#goods_attr_batch").val(data.join(","));
        });
        /**批量修改--征免方式回车**/
        $("body").on("keyup", "#duty_mode_name_batch", function (event) {
            const eCode = event.keyCode ? event.keyCode : event.which ? event.which : event.charCode;
            if (event.shiftKey != 1 && eCode == 13) {
                $("#goods_attr_open_batch").click();
            }
        });
        /**批量修改--生产批次回车**/
        $("body").on("keyup", "#prod_batch_no_batch", function (event) {
            const eCode = event.keyCode ? event.keyCode : event.which ? event.which : event.charCode;
            if (event.shiftKey != 1 && eCode == 13) {
                $("#batch_edit_save").click();
                return false
            }
        });
        /**批量修改--保存**/
        form.on('submit(batch_edit_save)', async (data) => {
            const checkData = table.checkStatus("order_pros").data;
            const data_map = checkData.map(item => item.g_no);
            for (let item of data_map) {
                for (let value in data.field) {
                    if (value.indexOf("type") > -1) {
                        if (data.field[value]) {
                            const value_rep = value.replace("type_", "");
                            order_pros_data[parseInt(item) - 1].goods_spec_data[value_rep] = data.field[value]
                        }
                    } else {
                        if (data.field[value]) {
                            order_pros_data[parseInt(item) - 1][value] = data.field[value]
                        }
                    }
                }
                if (goods_attr_batch_data) {
                    order_pros_data[parseInt(item) - 1].goods_attr_data = goods_attr_batch_data
                }
            }
            table.reload('order_pros', {
                data: order_pros_data,
                limit: order_pros_data.length
            });
            layer.closeAll();
            for (let item of data_map) {
                $(`.layui-table-view[lay-id='order_pros'] .layui-table-body tr[data-index=${parseInt(item) - 1}] input`).prop('checked', true);
            }
            form.render('checkbox');
            $("#dec_users_submit").click();
        });
        /**集装箱**/
        let order_containers_index = null;
        table.render({
            elem: '#order_containers_table'
            , toolbar: '#order_containers_tool'
            , defaultToolbar: ['filter']
            , colFilterRecord: 'local'
            , cols: [[
                { type: 'checkbox' }
                , { field: 'container_id', title: '集装箱号', width: 120 }
                , { field: 'container_md_name', title: '集装箱规格', width: 150 }
                , { field: 'lcl_flag_name', title: '拼箱标识', width: 120 }
            ]]
            , data: order_containers
            , limit: order_containers.length
            , height: 200
        });
        /**自动完成--集装箱规格**/
        auto_fn({
            data: admin.chemicals_data,
            listDirection: false,
            id: ['#lcl_flag_name'],
            after: ['#lcl_flag']
        });
        /**集装箱--商品序号默认全选**/
        $("body").on("keyup", "#container_id", function (event) {
            const eCode = event.keyCode ? event.keyCode : event.which ? event.which : event.charCode;
            if (event.shiftKey != 1 && eCode == 13) {
                if (order_pros_data.length > 0) {
                    const data = order_pros_data.map(item => item.g_no);
                    $("#goods_no").val(data.join(","));
                }
            }
        });
        /**集装箱--点击行反填数据**/
        table.on('row(order_containers_table)', function (obj) {
            for (let item in obj.data) {
                $(`#${item}`).val(obj.data[item])
            }
            order_containers_index = obj.data.index;
            $("#container_id").focus();
        });
        /**集装箱--保存**/
        form.on('submit(order_containers_submit)', async (data) => {
            delete data.field.layTableCheckbox;
            if (order_containers_index) {
                for (let item in data.field) {
                    order_containers[order_containers_index - 1][item] = data.field[item]
                }
            } else {
                order_containers.push(data.field);
            }
            order_containers.forEach((value, index) => {
                value.index = index + 1
            });
            table.reload('order_containers_table', {
                data: order_containers,
                limit: order_containers.length
            });
            $("#order_containers_form input").each(function () {
                $(this).val("")
            });
            $("#container_id").focus();
            order_containers_index = null;
            let contaCount = 0;
            for (let item of order_containers) {
                const containerMdCode = item.container_md;
                if (containerMdCode == "11" || containerMdCode == "12" || containerMdCode == "13" || containerMdCode == "32") {
                    contaCount += 2;
                } else if (containerMdCode == "21" || containerMdCode == "22" || containerMdCode == "23" || containerMdCode == "31") {
                    contaCount += 1;
                }
            }
            if (!contaCount || contaCount == "0") {
                contaCount = "";
            }
            $("#container_counts").val(contaCount);
        });
        /**集装箱--删除**/
        table.on('toolbar(order_containers_table)', function (obj) {
            let checkStatus = table.checkStatus(obj.config.id);
            let checkData = checkStatus.data;
            switch (obj.event) {
                case 'add':
                    $("#order_containers_form input").each(function () {
                        $(this).val("")
                    });
                    $("#container_id").focus();
                    order_containers_index = null;
                    break;
                case 'delete':
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    layer.confirm('真的删除么', { title: '提示' }, async (index) => {
                        let sup_data = order_containers.filter(item => {
                            return checkData.every(item2 => {
                                return item.index != item2.index;
                            })
                        });
                        for (let item_apply of checkData) {
                            if (item_apply.id) {
                                admin.order_containers_delete_ids.push(item_apply.id);
                            }
                        }
                        order_containers = sup_data;
                        order_containers.forEach((value, index) => {
                            value.index = index + 1
                        });
                        table.reload('order_containers_table', {
                            data: order_containers,
                            limit: order_containers.length
                        });
                        $("#order_containers_form input").each(function () {
                            $(this).val("")
                        });
                        $("#container_id").focus();
                        order_containers_index = null;
                        layer.close(index);
                    });
                    break;
            }
        });
        /**集装箱--按enter保存**/
        $("body").on("keyup", "#goods_no", function (event) {
            const eCode = event.keyCode ? event.keyCode : event.which ? event.which : event.charCode;
            if (event.shiftKey != 1 && eCode == 13) {
                $("#order_containers_submit").click();
                $("#container_id").focus();
            }
        });
        /**集装箱--商品序号关系列表**/
        $("body").on("click", "#goods_no_open", function () {
            layer.open({
                type: 1,
                title: '商品序号关系',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['910px', '480px'],
                content: $('#goods_no_list').html(),
                success: function (layero, index) {
                    this.enterEsc = function (event) {
                        if (event.keyCode === 13) {
                            $("#goods_no_save").click();
                            return false;
                        }
                    };
                    $(document).on('keydown', this.enterEsc);
                },
                end: function () {
                    $(document).off('keydown', this.enterEsc);
                }
            });
            const check_data = $("#goods_no").val().split(",").map(Number);
            table.render({
                elem: '#goods_no_table'
                , toolbar: true
                , defaultToolbar: ['filter']
                , colFilterRecord: 'local'
                , primaryKey: 'g_no'
                , checkStatus: {
                    default: check_data
                }
                , cols: [[
                    { type: 'checkbox' }
                    , { field: 'g_no', title: '序号', width: 100 }
                    , { field: 'code_t_s', title: '商品编号' }
                    , { field: 'g_name', title: '商品名称' }
                ]]
                , data: order_pros_data
                , limit: order_pros_data.length
                , height: 350
            });
        });
        /**集装箱--商品序号关系列表--保存**/
        $("body").on("click", "#goods_no_save", function () {
            let checkData = table.checkStatus("goods_no_table").data;
            const data = checkData.map(item => item.g_no);
            $("#goods_no").val(data.join(","));
            $("#goods_no").focus();
            layer.closeAll();
        });
        /**随附单证**/
        let order_documents_index = null;
        table.render({
            elem: '#order_documents_table'
            , toolbar: '#order_documents_tool'
            , defaultToolbar: ['filter']
            , colFilterRecord: 'local'
            , cols: [[
                { type: 'checkbox' }
                , { field: 'docu_code', title: '单证代码', width: 120 }
                , { field: 'cert_code', title: '单证编号' }
            ]]
            , data: order_documents
            , limit: order_documents.length
            , height: 200
        });
        /**随附单证--点击行反填数据**/
        table.on('row(order_documents_table)', function (obj) {
            for (let item in obj.data) {
                $(`#${item}`).val(obj.data[item])
            }
            order_documents_index = obj.data.index;
            $("#docu_code_name").focus();
        });
        /**随附单证--保存**/
        form.on('submit(order_documents_submit)', async (data) => {
            delete data.field.layTableCheckbox;
            if (order_documents_index) {
                for (let item in data.field) {
                    order_documents[order_documents_index - 1][item] = data.field[item]
                }
            } else {
                order_documents.push(data.field);
            }
            order_documents.forEach((value, index) => {
                value.index = index + 1
            });
            table.reload('order_documents_table', {
                data: order_documents,
                limit: order_documents.length
            });
            $("#order_documents_form input").each(function () {
                $(this).val("")
            });
            $("#docu_code_name").focus();
            order_documents_index = null;
            const document_code_string = order_documents.map(item => item.docu_code);
            $("#document_code_string").val(document_code_string.join(","));
        });
        /**随附单证--操作按钮**/
        /**原产地对应关系录入--key**/
        let order_document_eco_relations_index = null;
        table.on('toolbar(order_documents_table)', function (obj) {
            let checkStatus = table.checkStatus(obj.config.id);
            let checkData = checkStatus.data;
            switch (obj.event) {
                case 'add':
                    $("#order_documents_form input").each(function () {
                        $(this).val("")
                    });
                    $("#docu_code_name").focus();
                    order_documents_index = null;
                    break;
                case 'delete':
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    layer.confirm('真的删除么', { title: '提示' }, async (index) => {
                        let sup_data = order_documents.filter(item => {
                            return checkData.every(item2 => {
                                return item.index != item2.index;
                            })
                        });
                        for (let item_apply of checkData) {
                            if (item_apply.id) {
                                admin.order_documents_delete_ids.push(item_apply.id);
                            }
                        }
                        order_documents = sup_data;
                        order_documents.forEach((value, index) => {
                            value.index = index + 1
                        });
                        table.reload('order_documents_table', {
                            data: order_documents,
                            limit: order_documents.length
                        });
                        $("#order_documents_form input").each(function () {
                            $(this).val("")
                        });
                        $("#docu_code_name").focus();
                        order_documents_index = null;
                        const document_code_string = order_documents.map(item => item.docu_code);
                        $("#document_code_string").val(document_code_string.join(","));
                        layer.close(index);
                    });
                    break;
                case 'ciq_rel':
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    const acmpFormCode = $("#docu_code").val();
                    if ("Y" != acmpFormCode && "E" != acmpFormCode && "R" != acmpFormCode && "F" != acmpFormCode && "J" != acmpFormCode) {
                        layer.msg('单证代码不是Y/E/R/F/J!');
                        return;
                    }
                    layer.open({
                        type: 1,
                        title: '原产地对应关系录入',
                        shadeClose: true,
                        area: admin.screen() < 2 ? ['80%', '300px'] : ['910px', '530px'],
                        content: $('#order_document_eco_relations_list').html(),
                        end: function () {
                            $(`.layui-table-view[lay-id='order_documents_table'] .layui-table-body tr input`).each(function (index, el) {
                                el.checked = false;
                            });
                            form.render('checkbox');
                            $("#order_documents_add").click();
                        }
                    });
                    if (!(order_documents[order_documents_index - 1].eco_relations)) {
                        order_documents[order_documents_index - 1].eco_relations = [];
                    }
                    table.render({
                        elem: '#order_document_eco_relations_table'
                        , toolbar: '#order_document_eco_relations_tool'
                        , defaultToolbar: ['filter']
                        , colFilterRecord: 'local'
                        , cols: [[
                            { type: 'checkbox' }
                            , { field: 'dec_g_no', title: '报关单商品序号' }
                            , { field: 'eco_g_no', title: '对应随附单证商品序号' }
                        ]]
                        , data: order_documents[order_documents_index - 1].eco_relations
                        , limit: order_documents[order_documents_index - 1].eco_relations.length
                        , height: 350
                    });
                    $("#dec_g_no").focus();
                    break;
            }
        });
        /**随附单证--按enter保存**/
        $("body").on("keyup", "#cert_code", function (event) {
            const eCode = event.keyCode ? event.keyCode : event.which ? event.which : event.charCode;
            if (event.shiftKey != 1 && eCode == 13) {
                const acmpFormCode = $("#docu_code").val();
                if ("Y" != acmpFormCode && "E" != acmpFormCode && "R" != acmpFormCode && "F" != acmpFormCode && "J" != acmpFormCode) {
                    $("#order_documents_submit").click();
                } else {
                    $("#order_documents_submit").click();
                    /**let index = order_documents.length - 1;
                     $(`.layui-table-view[lay-id='order_documents_table'] .layui-table-body tr[data-index=${index}] .layui-form-checkbox`).click();
                     $("#order_documents_ciq_rel").click();**/
                }
            }
        });
        /**随附单证--按enter判断不能有重复的随附单证代码**/
        $("body").on("keyup", "#docu_code_name", function (event) {
            const eCode = event.keyCode ? event.keyCode : event.which ? event.which : event.charCode;
            let docu_code_name = 0;
            if (event.shiftKey != 1 && eCode == 13) {
                if (order_documents.length > 0) {
                    for (let item of order_documents) {
                        if (item.docu_code_name == $("#docu_code_name").val()) {
                            docu_code_name = 1
                        }
                    }
                    if (docu_code_name) {
                        $("#docu_code_name").val("");
                        $("#docu_code_name").focus();
                        return layer.msg("随附单证代码不可重复!")
                    }
                }
            }
        });
        /**原产地对应关系录入--点击行反填数据**/
        table.on('row(order_document_eco_relations_table)', function (obj) {
            $("#dec_g_no").val(obj.data.dec_g_no);
            $("#eco_g_no").val(obj.data.eco_g_no);
            order_document_eco_relations_index = obj.data.index;
            $("#dec_g_no").focus();
        });
        /**原产地对应关系录入--保存**/
        form.on('submit(order_document_eco_relations_submit)', async (data) => {
            delete data.field.layTableCheckbox;
            if (order_document_eco_relations_index) {
                for (let item in data.field) {
                    order_documents[order_documents_index - 1].eco_relations[order_document_eco_relations_index - 1][item] = data.field[item]
                }
            } else {
                order_documents[order_documents_index - 1].eco_relations.push(data.field);
            }
            order_documents[order_documents_index - 1].eco_relations.forEach((value, index) => {
                value.index = index + 1
            });
            table.reload('order_document_eco_relations_table', {
                data: order_documents[order_documents_index - 1].eco_relations,
                limit: order_documents[order_documents_index - 1].eco_relations.length
            });
            $("#dec_g_no").val("");
            $("#eco_g_no").val("");
            $("#dec_g_no").focus();
            order_document_eco_relations_index = null;
        });
        /**原产地对应关系录入--删除**/
        table.on('toolbar(order_document_eco_relations_table)', function (obj) {
            let checkStatus = table.checkStatus(obj.config.id);
            let checkData = checkStatus.data;
            switch (obj.event) {
                case 'add':
                    $("#dec_g_no").val("");
                    $("#eco_g_no").val("");
                    $("#dec_g_no").focus();
                    order_document_eco_relations_index = null;
                    break;
                case 'delete':
                    if (checkData.length === 0) {
                        return layer.msg('请选择数据');
                    }
                    layer.confirm('真的删除么', { title: '提示' }, async (index) => {
                        let sup_data = order_documents[order_documents_index - 1].eco_relations.filter(item => {
                            return checkData.every(item2 => {
                                return item.index != item2.index;
                            })
                        });
                        for (let item_apply of checkData) {
                            if (item_apply.id) {
                                admin.order_eco_relations_delete_ids.push({
                                    index: order_documents_index - 1,
                                    id: item_apply.id
                                });
                            }
                        }
                        order_documents[order_documents_index - 1].eco_relations = sup_data;
                        order_documents[order_documents_index - 1].eco_relations.forEach((value, index) => {
                            value.index = index + 1
                        });
                        table.reload('order_document_eco_relations_table', {
                            data: order_documents[order_documents_index - 1].eco_relations,
                            limit: order_documents[order_documents_index - 1].eco_relations.length
                        });
                        $("#dec_g_no").val("");
                        $("#eco_g_no").val("");
                        $("#dec_g_no").focus();
                        order_document_eco_relations_index = null;
                        layer.close(index);
                    });
                    break;
            }
        });
        /**原产地对应关系录入--按enter保存**/
        $("body").on("keyup", "#eco_g_no", function (event) {
            const eCode = event.keyCode ? event.keyCode : event.which ? event.which : event.charCode;
            if (event.shiftKey != 1 && eCode == 13) {
                $("#order_document_eco_relations_submit").click();
                $("#dec_g_no").focus();
            }
        });

        /**新建进口报关保存刷新后反填数据**/
        async function order_i_edit_back_filling(order_i_edit_data) {
            $("#order_i_form input, #order_i_form textarea").each(function () {
                for (let item in order_i_edit_data) {
                    if ($(this).attr("name") == item) {
                        if (item == 'dec_other_packs') {
                            $(this).val(JSON.stringify(order_i_edit_data[item]))
                        } else if (item == 'i_e_date') {
                            if (order_i_edit_data[item]) {
                                $(this).val(order_i_edit_data[item]);
                            }
                        } else {
                            if (item == 'apl_date') {
                                if (order_i_edit_data[item]) {
                                    const time = admin.getyyyymmdd(order_i_edit_data[item]);
                                    $(this).val(time);
                                }
                            } else if (item == 'gross_wet') {
                                $(this).val(order_i_edit_data['gross_wet_string']);
                            } else if (item == 'net_wt') {
                                $(this).val(order_i_edit_data['net_wt_string']);
                            } else {
                                $(this).val(order_i_edit_data[item])
                            }
                        }
                    }
                }
            });
            $("#customMasterName").val(order_i_edit_data.order_status_string);
            $("#client_seq_no").val(order_i_edit_data.client_seq_no);
            $("#entry_id").val(order_i_edit_data.entry_id);
            if (order_i_edit_data.promise_itmes) {
                $(`input[name='promise_itmes[]']`).each(function (index, element) {
                    for (let item in order_i_edit_data.promise_itmes) {
                        if (index == item) {
                            if (order_i_edit_data.promise_itmes[item] == '1') {
                                $(this).next('div').click();
                            } else if (order_i_edit_data.promise_itmes[item] == '0') {
                                $(this).data('first', 1);
                                $(this).prev("span").text("否")
                            }
                        }
                    }
                });
            }
            if (!(order_i_edit_data.type)) {
                $(`input[name='type[]']`).each(function (index, element) {
                    if (index == 1) {
                        $(this).next('div').click();
                    }
                });
            }
            await order_i_add_back_filling(order_i_edit_data);
        }

        /**新建进口报关保存反填数据**/
        async function order_i_add_back_filling(order_i_edit_data) {
            if (parent.layui.admin.get_iframe_index()) {
                parent.layui.admin.get_iframe_index().find("span").append(order_i_edit_data.client_seq_no);
            }
            if (order_i_edit_data.quas) {
                if (order_i_edit_data.quas.length > 0) {
                    ent_qualif_data = order_i_edit_data.quas;
                    ent_qualif_data.forEach((value, index) => {
                        value.index = index + 1
                    });
                    $("#ent_qualif_type_code").val(order_i_edit_data.quas[0].ent_qualif_type_code);
                    $("#ent_qualif_type_name").val(order_i_edit_data.quas[0].ent_qualif_type_name);
                    if (order_i_edit_data.declaratio_material_code == "101040") {
                        declaratio_material_code = 1
                    } else {
                        declaratio_material_code = 0
                    }
                }
            }
            if (order_i_edit_data.orig_box_flag == '0') {
                $("input[name='orig_box_flag']").next('div').click()
            }
            if (order_i_edit_data.is_other) {
                if (!($("input[name='is_other']").prop("checked"))) {
                    $("input[name='is_other']").next('div').click()
                }
            }
            if (order_i_edit_data.spec_decl_flag_array) {
                if (order_i_edit_data.spec_decl_flag_array.length > 0) {
                    spec_decl_flag = order_i_edit_data.spec_decl_flag_array;
                    let spec_decl_flag_name_data = [];
                    admin.spec_decl_flag_data.forEach((value, index) => {
                        if (spec_decl_flag[index] == '1') {
                            spec_decl_flag_name_data.push(value)
                        }
                    });
                    $("#spec_decl_flag").val(spec_decl_flag_name_data.join(","));
                }
            }
            if (order_i_edit_data.dec_request_certs) {
                if (order_i_edit_data.dec_request_certs.length > 0) {
                    const data = order_i_edit_data.dec_request_certs.map(item => item.app_cert_name);
                    $("#dec_request_certs").val(data.join(","));
                    dec_request_certs_data.checkData = order_i_edit_data.dec_request_certs;
                }
            }
            dec_request_certs_data.domestic_consignee_ename = order_i_edit_data.domestic_consignee_ename;
            dec_request_certs_data.overseas_consignor_cname = order_i_edit_data.overseas_consignor_cname;
            dec_request_certs_data.overseas_consignor_addr = order_i_edit_data.overseas_consignor_addr;
            dec_request_certs_data.cmpl_dschrg_dt = order_i_edit_data.cmpl_dschrg_dt;
            if (order_i_edit_data.dec_users) {
                dec_users_data = order_i_edit_data.dec_users;
                dec_users_data.forEach((value, index) => {
                    value.index = index + 1
                });
            }
            if (order_i_edit_data.pros) {
                order_i_edit_data.pros.sort(admin.compare('g_no'));
                const goods_attr_filter = await admin.post(`/clearance/no_paginate`,{Type:"货物属性代码"});
                for (let item of order_i_edit_data.pros) {
                    item.goods_spec_data = {
                        stuff: item.stuff,
                        prod_valid_dt: item.prod_valid_dt,
                        prod_qgp: item.prod_qgp,
                        eng_man_ent_cnm: item.eng_man_ent_cnm,
                        goods_spec: item.goods_spec,
                        goods_model: item.goods_model,
                        goods_brand: item.goods_brand,
                        produce_date: item.produce_date,
                        prod_batch_no: item.prod_batch_no
                    };
                    let goods_attr_data_item = [];
                    if (item.goods_attr_array) {
                        item.goods_attr_array.forEach(function (item_attr) {
                            goods_attr_filter.forEach(function (item_filter) {
                                if (item_filter.customs_code == item_attr) {
                                    goods_attr_data_item.push({
                                        code: item_attr,
                                        name: item_filter.name
                                    })
                                }
                            })
                        });
                    }
                    item.goods_attr_data = goods_attr_data_item;
                    item.quas.forEach(function (item_quas) {
                        if (item_quas.order_pro_qua_vins) {
                            item_quas.order_pro_qua_vins.forEach((value, index) => {
                                value.index = index + 1
                            });
                        }
                    });
                    let no_dang_flag_name;
                    for (let item_chemicals of admin.chemicals_data) {
                        if (item_chemicals.id == item.no_dang_flag) {
                            no_dang_flag_name = item_chemicals.value
                        }
                    }
                    let dang_pack_type_name;
                    for (let item_category of admin.category_data) {
                        if (item_category.id == item.dang_pack_type) {
                            dang_pack_type_name = item_category.value
                        }
                    }
                    item.dang_data = {
                        no_dang_flag: item.no_dang_flag,
                        no_dang_flag_name: no_dang_flag_name,
                        un_code: item.un_code,
                        dang_name: item.dang_name,
                        dang_pack_type: item.dang_pack_type,
                        dang_pack_type_name: dang_pack_type_name,
                        dang_pack_spec: item.dang_pack_spec
                    }
                }
                order_pros_data = order_i_edit_data.pros;
                $("#g_no").val(order_i_edit_data.pros.length + 1);
                table.reload('order_pros', {
                    data: order_pros_data,
                    limit: order_pros_data.length
                });
            }
            if (order_i_edit_data.containers) {
                order_containers = order_i_edit_data.containers;
                order_containers.forEach((value, index) => {
                    value.index = index + 1;
                    let lcl_flag_name;
                    for (let item_chemicals of admin.chemicals_data) {
                        if (item_chemicals.id == value.lcl_flag) {
                            lcl_flag_name = item_chemicals.value
                        }
                    }
                    value.lcl_flag_name = lcl_flag_name;
                });
                table.reload('order_containers_table', {
                    data: order_containers,
                    limit: order_containers.length
                });
            }
            if (order_i_edit_data.documents) {
                order_documents = order_i_edit_data.documents;
                order_documents.forEach((value, index) => {
                    value.index = index + 1;
                    if (value.eco_relations) {
                        value.eco_relations.forEach((value2, index2) => {
                            value2.index = index2 + 1;
                        })
                    }
                });
                table.reload('order_documents_table', {
                    data: order_documents,
                    limit: order_documents.length
                });
            }
            order_note_data = order_i_edit_data.remark;
            if (order_note_data) {
                $("#order_note_dot").show()
            } else {
                $("#order_note_dot").hide()
            }
            /**tip:总价/成交数量合计/法定第一数量合计/法定第二数量合计**/
            admin.is_total_number(order_i_edit_data.pros);

            if (order_i_edit_data.company) {
                $("#order_dispatch").data("user", order_i_edit_data.company.user_id);
            }
            admin.order_pros_delete_ids = [];
            admin.order_quas_delete_ids = [];
            admin.order_quas_vin_delete_ids = [];
            admin.order_containers_delete_ids = [];
            admin.order_documents_delete_ids = [];
            admin.order_eco_relations_delete_ids = [];
            admin.order_ent_qualif_delete_ids = [];
        };
        /**进口报关整合申报保存**/
        form.on('submit(order_save)', async (data) => {
            if (data.field.net_wt) {
                if (order_pros_data.length === 0) {
                    return layer.msg("存在净重，必须填写表体商品！");
                }
            }
            let is_pros = 1;
            if (order_pros_data.length > 0) {
                for (let item of order_pros_data) {
                    if ((item.g_unit ? item.g_unit != '035' : 1) && (item.first_unit ? item.first_unit != '035' : 1) && (item.second_unit ? item.second_unit != '035' : 1)) {
                        is_pros = 0;
                        break;
                    }
                }
            } else {
                is_pros = 0;
            }
            if (is_pros) {
                if (!(data.field.net_wt)) {

                    return layer.msg("存在表体商品且是千克，必填净重！");
                }
                let total = 0;
                for (let item of order_pros_data) {
                    if (item.g_unit == '035') {
                        total += parseFloat(item.g_qty) * 100000;
                    } else if (item.first_unit == '035') {
                        total += parseFloat(item.first_qty) * 100000;
                    } else if (item.second_unit == '035') {
                        total += parseFloat(item.second_qty) * 100000;
                    }
                }
                if (admin.cutZero((total / 100000).toString()) != admin.cutZero(data.field.net_wt.toString())) {
                    if (!no_weight_save) {
                        layer.open({
                            title: '商品总重量与净重不符'
                            , content: '商品总重量与净重不符！是否要保存'
                            , btn: ['是', '否']
                            , yes: function (index) {
                                no_weight_save = true;
                                layer.close(index);
                                $('#order_save').click();
                            }
                            , btn2: function () {
                                layer.msg('请重新修改净重。');
                                $('#net_wt').focus();
                            }
                        });
                        if (!no_weight_save) {
                            return false;
                        }
                    }
                }
            }
            layer.load(2);
            /**企业承诺**/
            data.field.declaratio_material_code = declaratio_material_code;
            /**其他包装**/
            if (data.field.dec_other_packs) {
                data.field.dec_other_packs = JSON.parse(data.field.dec_other_packs);
            }
            /**企业资质数据**/
            for (let item of admin.order_ent_qualif_delete_ids) {
                ent_qualif_data.push({
                    id: item
                })
            }
            data.field.quas = ent_qualif_data;
            /**使用人数据**/
            data.field.dec_users = dec_users_data;
            /**特殊业务标识数据**/
            data.field.spec_decl_flag = spec_decl_flag;
            /**表体商品数据**/
            for (let item of admin.order_quas_vin_delete_ids) {
                order_pros_data[item.pros_index].quas[item.index].order_pro_qua_vins.push({
                    id: item.id
                })
            }
            for (let item of admin.order_quas_delete_ids) {
                order_pros_data[item.index].quas.push({
                    id: item.id
                })
            }
            for (let item of admin.order_pros_delete_ids) {
                order_pros_data.push({
                    id: item
                })
            }
            data.field.pros = order_pros_data;
            /**检验检疫申报要素数据**/
            data.field.dec_request_certs = dec_request_certs_data.checkData || [];
            data.field.domestic_consignee_ename = dec_request_certs_data.domestic_consignee_ename;
            data.field.overseas_consignor_cname = dec_request_certs_data.overseas_consignor_cname;
            data.field.overseas_consignor_addr = dec_request_certs_data.overseas_consignor_addr;
            data.field.cmpl_dschrg_dt = dec_request_certs_data.cmpl_dschrg_dt;
            /**集装箱数据**/
            for (let item of admin.order_containers_delete_ids) {
                order_containers.push({
                    id: item
                })
            }
            data.field.containers = order_containers;
            /**随附单证数据**/
            for (let item of admin.order_documents_delete_ids) {
                order_documents.push({
                    id: item
                })
            }
            data.field.documents = order_documents;
            /**业务选项**/
            data.field.promise_itmes = [];
            /**特殊关系确认**/
            if (!(data.field['promise_itmes[0]'])) {
                if ($("#promise_itmes0").data('first') == '1') {
                    data.field.promise_itmes.push('0');
                } else {
                    data.field.promise_itmes.push('9');
                }
            } else {
                delete data.field['promise_itmes[0]'];
                data.field.promise_itmes.push('1');
            }
            /**价格影响确认**/
            if (!(data.field['promise_itmes[1]'])) {
                if ($("#promise_itmes1").data('first') == '1') {
                    data.field.promise_itmes.push('0');
                } else {
                    data.field.promise_itmes.push('9');
                }
            } else {
                delete data.field['promise_itmes[1]'];
                data.field.promise_itmes.push('1');
            }
            /**支付特权使用费确认**/
            if (!(data.field['promise_itmes[2]'])) {
                if ($("#promise_itmes2").data('first') == '1') {
                    data.field.promise_itmes.push('0');
                } else {
                    data.field.promise_itmes.push('9');
                }
            } else {
                delete data.field['promise_itmes[2]'];
                data.field.promise_itmes.push('1');
            }
            /**异地报关**/
            if ($("#is_other").prop("checked")) {
                data.field.is_other = 1;
            } else {
                data.field.is_other = 0;
            }

            data.field.id = order_id;
            data.field.order_status = order_i_edit_data.order_status_string;
            const order_save_data = await admin.post(`/order/i/${order_id}/store_change`, data.field);
            if (order_save_data.status) {
                await order_i_add_back_filling(order_save_data.data);
                order_i_edit_data = order_save_data.data;
            }
            layer.closeAll('loading');
        });
        /**派单**/
        let distribute_index;
        $("body").on("click", "#order_dispatch", async function () {
            if (!order_id) {
                return layer.msg("请先保存订单！")
            }
            distribute_index = layer.open({
                type: 1,
                title: '派单',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['650px', '340px'],
                content: $('#distribute_template').html()
            });
            const user_id = $(this).data("user");
            $("select[name='user_id']").find(`option[value=${user_id}]`).prop("selected", true);
            form.render();
        });
        $("body").on("input", "#remark_distribute", function () {
            $("#remark_distribute_number span").text($(this).val().length);
        });
        /**派单保存**/
        form.on('submit(distribute_submit)', async (data) => {
            await admin.post(`/order/i/${order_id}/distribute`, data.field);
            layer.close(distribute_index);
        });
        /**打印**/
        $("body").on("click", "#order_print", async function () {
            if (!order_id) {
                return layer.msg("请先保存订单！")
            }
            layer.open({
                type: 1,
                title: '打印',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['1020px', '580px'],
                content: $('#print_lists_template').html()
            });
            layer.load(2);
            await admin.getPdf(order_id);
            table.render({
                elem: '#print_lists'
                , data: admin.order_i_print_list
                , cols: [[
                    {
                        field: 'name', title: '类型', templet: function (data) {
                            if (data.is_enclosure) {
                                return `<p>${data.name}<span class="enclosure_span">附件</span></p>`
                            } else {
                                return data.name
                            }
                        }
                    }
                    , { title: '操作', toolbar: '#print_toolbar', width: 280 }
                ]]
                , limit: admin.order_i_print_list.length
            });
        });
        /**打印按钮操作**/
        let print_edit_index
            , file_type_id
            , orientation_save
            , save_pdf_id;
        table.on('tool(print_lists)', async function (obj) {
            switch (obj.event) {
                case 'preview':
                    if (obj.data.code == 'local_5') {
                        orientation_save = ''
                    } else {
                        if (obj.data.code == '00000004' || obj.data.code == '00000001' || obj.data.code == '00000002' || obj.data.code == 'local_9') {
                            orientation_save = `&orientation=portrait&margin_top=0cm&margin_left=0cm&margin_right=0cm&margin_bottom=0cm`
                        }
                        if (obj.data.code == 'local_6') {
                            orientation_save = `&orientation=landscape&margin_top=0cm&margin_left=0cm&margin_right=0cm&margin_bottom=0cm`
                        }
                    }
                    if (obj.data.code == 'local_5') {
                        if (obj.data.is_enclosure) {
                            const preview_data = await admin.get(`/order/i/pdf/${obj.data.pdf_id}/show`);
                            window.open(preview_data);
                        } else {
                            window.open(`/order/i/${order_id}/downloads/prints?file_type=${obj.data.id}`);
                        }
                    } else {
                        if (obj.data.is_enclosure) {
                            const preview_data = await admin.get(`/order/i/pdf/${obj.data.pdf_id}/show`);
                            window.open(preview_data);
                        } else {
                            if (obj.data.code == '00000004') {
                                file_type_id = obj.data.id;
                                layer.confirm('是否需要修改内容？', {
                                    btn: ['确定', '取消']
                                }, async function (index) {
                                    layer.close(index);
                                    print_edit_index = layer.open({
                                        type: 1,
                                        title: '打印前修改(进口订购合同)',
                                        shadeClose: true,
                                        area: admin.screen() < 2 ? ['80%', '300px'] : ['650px', '250px'],
                                        content: $('#print_edit_contract_template').html()
                                    });
                                    let contact_safe = $("#contact_safe").val();
                                    $("#contact_safe").val("").focus().val(contact_safe);
                                    $("#dest_code_print").val($("#dest_code").val());
                                    $("#dest_code_name_print").val(order_i_edit_data.pros[0].district_code_name);
                                    $("#trans_mode_print").val($("#trans_mode").val());
                                    $("#trans_mode_name_print").val($("#trans_mode_name").val());
                                    $("#contact_sign_date").val(order_i_edit_data.contact_sign_date);
                                    /*auto_fn({
                                        data: admin.all_complete_data.destination,
                                        listDirection: false,
                                        id: ['#dest_code_name_print'],
                                        after: ['#dest_code_print']
                                    });*/
                                    auto_fn({
                                        data: admin.all_complete_data.terms_delivery,
                                        listDirection: false,
                                        id: ['#trans_mode_name_print'],
                                        after: ['#trans_mode_print']
                                    });
                                    laydate.render({ elem: '#contact_sign_date', theme: '#1E9FFF', fixed: true });
                                }, async function (index) {
                                    layer.close(index);
                                    window.open(`/order/i/${order_id}/downloads/prints?file_type=${obj.data.id}${orientation_save}`);
                                });
                            }
                            if (obj.data.code == '00000001') {
                                file_type_id = obj.data.id;
                                layer.confirm('是否需要修改内容？', {
                                    btn: ['确定', '取消']
                                }, async function (index) {
                                    layer.close(index);
                                    print_edit_index = layer.open({
                                        type: 1,
                                        title: '打印前修改(进口发票类型)',
                                        shadeClose: true,
                                        area: admin.screen() < 2 ? ['80%', '300px'] : ['650px', '150px'],
                                        content: $('#print_edit_invoice_template').html()
                                    });
                                    $("#dest_code_name_invoice").val("").focus().val(order_i_edit_data.pros[0].district_code_name);
                                    $("#dest_code_invoice").val($("#dest_code").val());
                                    // auto_fn({
                                    //     data: admin.all_complete_data.domestic_area,
                                    //     listDirection: false,
                                    //     id: ['#dest_code_name_invoice'],
                                    //     after: ['#dest_code_invoice']
                                    // });
                                }, async function (index) {
                                    layer.close(index);
                                    window.open(`/order/i/${order_id}/downloads/prints?file_type=${obj.data.id}${orientation_save}`);
                                });
                            }
                            if (obj.data.code == 'local_6') {
                                file_type_id = obj.data.id;
                                layer.confirm('是否需要修改内容？', {
                                    btn: ['确定', '取消']
                                }, async function (index) {
                                    layer.close(index);
                                    print_edit_index = layer.open({
                                        type: 1,
                                        title: '打印前修改(进口司机纸)',
                                        shadeClose: true,
                                        area: admin.screen() < 2 ? ['80%', '300px'] : ['1050px', '630px'],
                                        content: $('#print_edit_driver_template').html()
                                    });
                                    layer.load(2);
                                    const driver_data = await admin.get(`/order/i/${order_id}/order_driver_paper?voy_no=${order_i_edit_data.voy_no}`);
                                        /** 货物申报司机纸打印弹窗表格渲染 */
                                        admin.getOrderDriverTable(driver_data, order_i_edit_data, '#driver_table')
                                    auto_fn({
                                        data: admin.all_complete_data.domestic_area,
                                        listDirection: false,
                                        id: ['#district_code_name_driver']
                                    });
                                    auto_fn({
                                        data: admin.all_complete_data.entry_clearance,
                                        listDirection: false,
                                        id: ['#i_e_port_name_driver'],
                                        after: ['#i_e_port_driver']
                                    });
                                    auto_fn({
                                        data: admin.all_complete_data.entry_clearance,
                                        listDirection: false,
                                        id: ['#i_e_port_name_two_driver'],
                                        after: ['#i_e_port_two_driver']
                                    });
                                    auto_fn({
                                        data: admin.all_complete_data.objectives_based,
                                        listDirection: false,
                                        id: ['#trade_mode_name_driver']
                                    });
                                    auto_fn({
                                        data: admin.all_complete_data.country_area,
                                        listDirection: false,
                                        id: ['#origin_country_name_driver']
                                    });
                                    auto_fn({
                                        data: admin.all_complete_data.entry_clearance,
                                        listDirection: false,
                                        id: ['#custom_master_name_driver'],
                                        after: ['#custom_master_driver']
                                    });
                                    $("#contr_no_driver").val("").focus().val(driver_data.contr_no);
                                    $("#district_code_name_driver").val(driver_data.first_pros[0].district_code_name);
                                    $("#i_e_port_name_driver").val(driver_data.i_e_port_name);
                                    $("#i_e_port_driver").val(driver_data.i_e_port);
                                    $("#i_e_port_name_two_driver").val(driver_data.i_e_port_name);
                                    $("#trade_mode_name_driver").val(driver_data.trade_mode_name);
                                    $("#gross_wet_driver").val(driver_data.gross_wet);
                                    $("#net_wt_driver").val(driver_data.net_wt);
                                    $("#origin_country_name_driver").val(driver_data.first_pros[0].origin_country_name);
                                    $("#pack_no_driver").val(driver_data.pack_no);
                                    $("#bill_no_driver").val(driver_data.bill_no);
                                    $("#car_code_driver").val(driver_data.car_code);
                                    $("#voy_no_driver").val(driver_data.voy_no);
                                    $("#custom_master_name_driver").val(driver_data.custom_master_name);
                                    $("#custom_master_driver").val(driver_data.custom_master);
                                    $("#container_ids_driver").val(driver_data.container_ids);
                                }, async function (index) {
                                    layer.close(index);
                                    window.open(`/order/i/${order_id}/downloads/prints?file_type=${obj.data.id}${orientation_save}`);
                                });
                            }
                            if (obj.data.code == '00000002' || obj.data.code == 'local_9') {
                                window.open(`/order/i/${order_id}/downloads/prints?file_type=${obj.data.id}${orientation_save}`);
                            }
                        }
                    }
                    break;
                case 'download':
                    if (obj.data.code == 'local_5') {
                        orientation_save = ''
                    } else {
                        if (obj.data.code == '00000004' || obj.data.code == '00000001' || obj.data.code == '00000002' || obj.data.code == 'local_9') {
                            orientation_save = `&orientation=portrait&margin_top=0cm&margin_left=0cm&margin_right=0cm&margin_bottom=0cm`
                        }
                        if (obj.data.code == 'local_6') {
                            orientation_save = `&orientation=landscape&margin_top=0cm&margin_left=0cm&margin_right=0cm&margin_bottom=0cm`
                        }
                    }
                    if (obj.data.code == 'local_5') {
                        if (obj.data.is_enclosure) {
                            window.open(`/order/i/pdf/${obj.data.pdf_id}/downloads`);
                        } else {
                            window.open(`/order/i/${order_id}/downloads/prints?file_type=${obj.data.id}`);
                        }
                    } else {
                        if (obj.data.is_enclosure) {
                            window.open(`/order/i/pdf/${obj.data.pdf_id}/downloads`);
                        } else {
                            if (obj.data.code == '00000004') {
                                file_type_id = obj.data.id;
                                layer.confirm('是否需要修改内容？', {
                                    btn: ['确定', '取消']
                                }, async function (index) {
                                    layer.close(index);
                                    print_edit_index = layer.open({
                                        type: 1,
                                        title: '打印前修改(进口订购合同)',
                                        shadeClose: true,
                                        area: admin.screen() < 2 ? ['80%', '300px'] : ['650px', '250px'],
                                        content: $('#print_edit_contract_template').html()
                                    });
                                    let contact_safe = $("#contact_safe").val();
                                    $("#contact_safe").val("").focus().val(contact_safe);
                                    $("#dest_code_print").val($("#dest_code").val());
                                    $("#dest_code_name_print").val(order_i_edit_data.pros[0].district_code_name);
                                    $("#trans_mode_print").val($("#trans_mode").val());
                                    $("#trans_mode_name_print").val($("#trans_mode_name").val());
                                    $("#contact_sign_date").val(order_i_edit_data.contact_sign_date);
                                    // auto_fn({
                                    //     data: admin.all_complete_data.destination,
                                    //     listDirection: false,
                                    //     id: ['#dest_code_name_print'],
                                    //     after: ['#dest_code_print']
                                    // });
                                    auto_fn({
                                        data: admin.all_complete_data.terms_delivery,
                                        listDirection: false,
                                        id: ['#trans_mode_name_print'],
                                        after: ['#trans_mode_print']
                                    });
                                    laydate.render({ elem: '#contact_sign_date', theme: '#1E9FFF', fixed: true });
                                }, async function (index) {
                                    layer.close(index);
                                    window.open(`/order/i/${order_id}/downloads/prints?file_type=${obj.data.id}${orientation_save}`);
                                });
                            }
                            if (obj.data.code == '00000001') {
                                file_type_id = obj.data.id;
                                layer.confirm('是否需要修改内容？', {
                                    btn: ['确定', '取消']
                                }, async function (index) {
                                    layer.close(index);
                                    print_edit_index = layer.open({
                                        type: 1,
                                        title: '打印前修改(进口发票类型)',
                                        shadeClose: true,
                                        area: admin.screen() < 2 ? ['80%', '300px'] : ['650px', '150px'],
                                        content: $('#print_edit_invoice_template').html()
                                    });
                                    $("#dest_code_name_invoice").val("").focus().val(order_i_edit_data.pros[0].district_code_name);
                                    $("#dest_code_invoice").val($("#dest_code").val());
                                    // auto_fn({
                                    //     data: admin.all_complete_data.domestic_area,
                                    //     listDirection: false,
                                    //     id: ['#dest_code_name_invoice'],
                                    //     after: ['#dest_code_invoice']
                                    // });
                                }, async function (index) {
                                    layer.close(index);
                                    window.open(`/order/i/${order_id}/downloads/prints?file_type=${obj.data.id}${orientation_save}`);
                                });
                            }
                            if (obj.data.code == 'local_6') {
                                file_type_id = obj.data.id;
                                layer.confirm('是否需要修改内容？', {
                                    btn: ['确定', '取消']
                                }, async function (index) {
                                    layer.close(index);
                                    print_edit_index = layer.open({
                                        type: 1,
                                        title: '打印前修改(进口司机纸)',
                                        shadeClose: true,
                                        area: admin.screen() < 2 ? ['80%', '300px'] : ['1050px', '630px'],
                                        content: $('#print_edit_driver_template').html()
                                    });
                                    layer.load(2);
                                    const driver_data = await admin.get(`/order/i/${order_id}/order_driver_paper?voy_no=${order_i_edit_data.voy_no}`);
                                        /** 货物申报司机纸打印弹窗表格渲染 */
                                        admin.getOrderDriverTable(driver_data, order_i_edit_data, '#driver_table')
                                    auto_fn({
                                        data: admin.all_complete_data.domestic_area,
                                        listDirection: false,
                                        id: ['#district_code_name_driver']
                                    });
                                    auto_fn({
                                        data: admin.all_complete_data.entry_clearance,
                                        listDirection: false,
                                        id: ['#i_e_port_name_driver'],
                                        after: ['#i_e_port_driver']
                                    });
                                    auto_fn({
                                        data: admin.all_complete_data.entry_clearance,
                                        listDirection: false,
                                        id: ['#i_e_port_name_two_driver'],
                                        after: ['#i_e_port_two_driver']
                                    });
                                    auto_fn({
                                        data: admin.all_complete_data.objectives_based,
                                        listDirection: false,
                                        id: ['#trade_mode_name_driver']
                                    });
                                    auto_fn({
                                        data: admin.all_complete_data.country_area,
                                        listDirection: false,
                                        id: ['#origin_country_name_driver']
                                    });
                                    auto_fn({
                                        data: admin.all_complete_data.entry_clearance,
                                        listDirection: false,
                                        id: ['#custom_master_name_driver'],
                                        after: ['#custom_master_driver']
                                    });
                                    $("#contr_no_driver").val("").focus().val(driver_data.contr_no);
                                    $("#district_code_name_driver").val(driver_data.first_pros[0].district_code_name);
                                    $("#i_e_port_name_driver").val(driver_data.i_e_port_name);
                                    $("#i_e_port_driver").val(driver_data.i_e_port);
                                    $("#i_e_port_name_two_driver").val(driver_data.i_e_port_name);
                                    $("#trade_mode_name_driver").val(driver_data.trade_mode_name);
                                    $("#gross_wet_driver").val(driver_data.gross_wet);
                                    $("#net_wt_driver").val(driver_data.net_wt);
                                    $("#origin_country_name_driver").val(driver_data.first_pros[0].origin_country_name);
                                    $("#pack_no_driver").val(driver_data.pack_no);
                                    $("#bill_no_driver").val(driver_data.bill_no);
                                    $("#car_code_driver").val(driver_data.car_code);
                                    $("#voy_no_driver").val(driver_data.voy_no);
                                    $("#custom_master_name_driver").val(driver_data.custom_master_name);
                                    $("#custom_master_driver").val(driver_data.custom_master);
                                    $("#container_ids_driver").val(driver_data.container_ids);
                                }, async function (index) {
                                    layer.close(index);
                                    window.open(`/order/i/${order_id}/downloads/prints?file_type=${obj.data.id}${orientation_save}`);
                                });
                            }
                            if (obj.data.code == '00000002' || obj.data.code == 'local_9') {
                                window.open(`/order/i/${order_id}/downloads/prints?file_type=${obj.data.id}${orientation_save}`);
                            }
                        }
                    }
                    break;
                case 'enclosure_save':
                    if (obj.data.code == 'local_5') {
                        orientation_save = ''
                    } else {
                        if (obj.data.code == '00000004' || obj.data.code == '00000001' || obj.data.code == '00000002' || obj.data.code == 'local_9') {
                            orientation_save = `&orientation=portrait&margin_top=0cm&margin_left=0cm&margin_right=0cm&margin_bottom=0cm`
                        }
                        if (obj.data.code == 'local_6') {
                            orientation_save = `&orientation=landscape&margin_top=0cm&margin_left=0cm&margin_right=0cm&margin_bottom=0cm`
                        }
                    }
                    if (obj.data.code == 'local_5') {
                        layer.load(2);
                        const save_pdf_data = await admin.get(`/order/i/${order_id}/downloads/save_pdf?file_type=${obj.data.id}${orientation_save}`, 'show');
                        layer.closeAll('loading');
                        if (save_pdf_data.status) {
                            await admin.getPdf(order_id);
                            table.reload('print_lists', {
                                data: admin.order_i_print_list,
                                limit: admin.order_i_print_list.length
                            });
                        }
                    } else {
                        if (obj.data.code == '00000004') {
                            file_type_id = obj.data.id;
                            save_pdf_id = null;
                            layer.confirm('是否需要修改内容？', {
                                btn: ['确定', '取消']
                            }, async function (index) {
                                layer.close(index);
                                print_edit_index = layer.open({
                                    type: 1,
                                    title: '保存前修改(进口订购合同)',
                                    shadeClose: true,
                                    area: admin.screen() < 2 ? ['80%', '300px'] : ['650px', '250px'],
                                    content: $('#print_save_contract_template').html()
                                });
                                let contact_safe = $("#contact_safe_save").val();
                                $("#contact_safe_save").val("").focus().val(contact_safe);
                                $("#dest_code_save").val($("#dest_code").val());
                                $("#dest_code_name_save").val(order_i_edit_data.pros[0].district_code_name);
                                $("#trans_mode_save").val($("#trans_mode").val());
                                $("#trans_mode_name_save").val($("#trans_mode_name").val());
                                $("#contact_sign_date_save").val(order_i_edit_data.contact_sign_date);
                                // auto_fn({
                                //     data: admin.all_complete_data.domestic_area,
                                //     listDirection: false,
                                //     id: ['#dest_code_name_save'],
                                //     after: ['#dest_code_save']
                                // });
                                auto_fn({
                                    data: admin.all_complete_data.terms_delivery,
                                    listDirection: false,
                                    id: ['#trans_mode_name_save'],
                                    after: ['#trans_mode_save']
                                });
                                laydate.render({ elem: '#contact_sign_date_save', theme: '#1E9FFF', fixed: true });
                            }, async function (index) {
                                layer.close(index);
                                layer.load(2);
                                const save_pdf_data = await admin.get(`/order/i/${order_id}/downloads/save_pdf?file_type=${obj.data.id}${orientation_save}`, 'show');
                                layer.closeAll('loading');
                                if (save_pdf_data.status) {
                                    await admin.getPdf(order_id);
                                    table.reload('print_lists', {
                                        data: admin.order_i_print_list,
                                        limit: admin.order_i_print_list.length
                                    });
                                }
                            });
                        }
                        if (obj.data.code == '00000001') {
                            file_type_id = obj.data.id;
                            save_pdf_id = null;
                            layer.confirm('是否需要修改内容？', {
                                btn: ['确定', '取消']
                            }, async function (index) {
                                layer.close(index);
                                print_edit_index = layer.open({
                                    type: 1,
                                    title: '保存前修改(进口发票类型)',
                                    shadeClose: true,
                                    area: admin.screen() < 2 ? ['80%', '300px'] : ['650px', '150px'],
                                    content: $('#print_save_invoice_template').html()
                                });
                                $("#dest_code_name_save").val("").focus().val(order_i_edit_data.pros[0].district_code_name);
                                $("#dest_code_save").val($("#dest_code").val());
                                // auto_fn({
                                //    data: admin.all_complete_data.domestic_area,
                                //    listDirection: false,
                                //    id: ['#dest_code_name_save'],
                                //    after: ['#dest_code_save']
                                // });
                            }, async function (index) {
                                layer.close(index);
                                layer.load(2);
                                const save_pdf_data = await admin.get(`/order/i/${order_id}/downloads/save_pdf?file_type=${obj.data.id}${orientation_save}`, 'show');
                                layer.closeAll('loading');
                                if (save_pdf_data.status) {
                                    await admin.getPdf(order_id);
                                    table.reload('print_lists', {
                                        data: admin.order_i_print_list,
                                        limit: admin.order_i_print_list.length
                                    });
                                }
                            });
                        }
                        if (obj.data.code == 'local_6') {
                            file_type_id = obj.data.id;
                            save_pdf_id = null;
                            layer.confirm('是否需要修改内容？', {
                                btn: ['确定', '取消']
                            }, async function (index) {
                                layer.close(index);
                                print_edit_index = layer.open({
                                    type: 1,
                                    title: '保存前修改(进口司机纸)',
                                    shadeClose: true,
                                    area: admin.screen() < 2 ? ['80%', '300px'] : ['1050px', '630px'],
                                    content: $('#print_save_driver_template').html()
                                });
                                layer.load(2);
                                const driver_data = await admin.get(`/order/i/${order_id}/order_driver_paper?voy_no=${order_i_edit_data.voy_no}`);
                                    /** 货物申报司机纸打印弹窗表格渲染 */
                                    admin.getOrderDriverTable(driver_data, order_i_edit_data, '#driver_table_save')
                                auto_fn({
                                    data: admin.all_complete_data.domestic_area,
                                    listDirection: false,
                                    id: ['#district_code_name_save']
                                });
                                auto_fn({
                                    data: admin.all_complete_data.entry_clearance,
                                    listDirection: false,
                                    id: ['#i_e_port_name_save'],
                                    after: ['#i_e_port_save']
                                });
                                auto_fn({
                                    data: admin.all_complete_data.entry_clearance,
                                    listDirection: false,
                                    id: ['#i_e_port_name_two_save'],
                                    after: ['#i_e_port_two_save']
                                });
                                auto_fn({
                                    data: admin.all_complete_data.objectives_based,
                                    listDirection: false,
                                    id: ['#trade_mode_name_save']
                                });
                                auto_fn({
                                    data: admin.all_complete_data.country_area,
                                    listDirection: false,
                                    id: ['#origin_country_name_save']
                                });
                                auto_fn({
                                    data: admin.all_complete_data.entry_clearance,
                                    listDirection: false,
                                    id: ['#custom_master_name_save'],
                                    after: ['#custom_master_save']
                                });
                                $("#contr_no_save").val("").focus().val(driver_data.contr_no);
                                $("#district_code_name_save").val(driver_data.first_pros[0].district_code_name);
                                $("#i_e_port_name_save").val(driver_data.i_e_port_name);
                                $("#i_e_port_save").val(driver_data.i_e_port);
                                $("#i_e_port_name_two_save").val(driver_data.i_e_port_name);
                                $("#trade_mode_name_save").val(driver_data.trade_mode_name);
                                $("#gross_wet_save").val(driver_data.gross_wet);
                                $("#net_wt_save").val(driver_data.net_wt);
                                $("#origin_country_name_save").val(driver_data.first_pros[0].origin_country_name);
                                $("#pack_no_save").val(driver_data.pack_no);
                                $("#bill_no_save").val(driver_data.bill_no);
                                $("#car_code_save").val(driver_data.car_code);
                                $("#voy_no_save").val(driver_data.voy_no);
                                $("#custom_master_name_save").val(driver_data.custom_master_name);
                                $("#custom_master_save").val(driver_data.custom_master);
                                $("#container_ids_save").val(driver_data.container_ids);
                            }, async function (index) {
                                layer.close(index);
                                layer.load(2);
                                const save_pdf_data = await admin.get(`/order/i/${order_id}/downloads/save_pdf?file_type=${obj.data.id}${orientation_save}`, 'show');
                                layer.closeAll('loading');
                                if (save_pdf_data.status) {
                                    await admin.getPdf(order_id);
                                    table.reload('print_lists', {
                                        data: admin.order_i_print_list,
                                        limit: admin.order_i_print_list.length
                                    });
                                }
                            });
                        }
                        if (obj.data.code == '00000002' || obj.data.code == 'local_9') {
                            layer.load(2);
                            const save_pdf_data = await admin.get(`/order/i/${order_id}/downloads/save_pdf?file_type=${obj.data.id}${orientation_save}`, 'show');
                            layer.closeAll('loading');
                            if (save_pdf_data.status) {
                                await admin.getPdf(order_id);
                                table.reload('print_lists', {
                                    data: admin.order_i_print_list,
                                    limit: admin.order_i_print_list.length
                                });
                            }
                        }
                    }
                    break;
                case 'reload_save':
                    if (obj.data.code == 'local_5') {
                        orientation_save = ''
                    } else {
                        if (obj.data.code == '00000004' || obj.data.code == '00000001' || obj.data.code == '00000002' || obj.data.code == 'local_9') {
                            orientation_save = `&orientation=portrait&margin_top=0cm&margin_left=0cm&margin_right=0cm&margin_bottom=0cm`
                        }
                        if (obj.data.code == 'local_6') {
                            orientation_save = `&orientation=landscape&margin_top=0cm&margin_left=0cm&margin_right=0cm&margin_bottom=0cm`
                        }
                    }
                    if (obj.data.code == 'local_5') {
                        layer.load(2);
                        const reload_save_pdf_data = await admin.get(`/order/i/${order_id}/downloads/save_pdf?file_type=${obj.data.id}&pdf_id=${obj.data.pdf_id}${orientation_save}`, 'show');
                        layer.closeAll('loading');
                        if (reload_save_pdf_data.status) {
                            await admin.getPdf(order_id);
                            table.reload('print_lists', {
                                data: admin.order_i_print_list,
                                limit: admin.order_i_print_list.length
                            });
                        }
                    } else {
                        if (obj.data.code == '00000004') {
                            file_type_id = obj.data.id;
                            save_pdf_id = obj.data.pdf_id;
                            layer.confirm('是否需要修改内容？', {
                                btn: ['确定', '取消']
                            }, async function (index) {
                                layer.close(index);
                                print_edit_index = layer.open({
                                    type: 1,
                                    title: '保存前修改(进口订购合同)',
                                    shadeClose: true,
                                    area: admin.screen() < 2 ? ['80%', '300px'] : ['650px', '250px'],
                                    content: $('#print_save_contract_template').html()
                                });
                                let contact_safe = $("#contact_safe_save").val();
                                $("#contact_safe_save").val("").focus().val(contact_safe);
                                $("#dest_code_save").val($("#dest_code").val());
                                $("#dest_code_name_save").val(order_i_edit_data.pros[0].district_code_name);
                                $("#trans_mode_save").val($("#trans_mode").val());
                                $("#trans_mode_name_save").val($("#trans_mode_name").val());
                                $("#contact_sign_date_save").val(order_i_edit_data.contact_sign_date);
                                // auto_fn({
                                //     data: admin.all_complete_data.domestic_area,
                                //     listDirection: false,
                                //     id: ['#dest_code_name_save'],
                                //     after: ['#dest_code_save']
                                // });
                                auto_fn({
                                    data: admin.all_complete_data.terms_delivery,
                                    listDirection: false,
                                    id: ['#trans_mode_name_save'],
                                    after: ['#trans_mode_save']
                                });
                                laydate.render({ elem: '#contact_sign_date_save', theme: '#1E9FFF', fixed: true });
                            }, async function (index) {
                                layer.close(index);
                                layer.load(2);
                                const reload_save_pdf_data = await admin.get(`/order/i/${order_id}/downloads/save_pdf?file_type=${obj.data.id}&pdf_id=${obj.data.pdf_id}${orientation_save}`, 'show');
                                layer.closeAll('loading');
                                if (reload_save_pdf_data.status) {
                                    await admin.getPdf(order_id);
                                    table.reload('print_lists', {
                                        data: admin.order_i_print_list,
                                        limit: admin.order_i_print_list.length
                                    });
                                }
                            });
                        }
                        if (obj.data.code == '00000001') {
                            file_type_id = obj.data.id;
                            save_pdf_id = obj.data.pdf_id;
                            layer.confirm('是否需要修改内容？', {
                                btn: ['确定', '取消']
                            }, async function (index) {
                                layer.close(index);
                                print_edit_index = layer.open({
                                    type: 1,
                                    title: '保存前修改(进口发票类型)',
                                    shadeClose: true,
                                    area: admin.screen() < 2 ? ['80%', '300px'] : ['650px', '150px'],
                                    content: $('#print_save_invoice_template').html()
                                });
                                $("#dest_code_name_save").val("").focus().val(order_i_edit_data.pros[0].district_code_name);
                                $("#dest_code_save").val($("#dest_code").val());
                                // auto_fn({
                                //     data: admin.all_complete_data.domestic_area,
                                //     listDirection: false,
                                //    id: ['#dest_code_name_save'],
                                //     after: ['#dest_code_save']
                                // });
                            }, async function (index) {
                                layer.close(index);
                                layer.load(2);
                                const reload_save_pdf_data = await admin.get(`/order/i/${order_id}/downloads/save_pdf?file_type=${obj.data.id}&pdf_id=${obj.data.pdf_id}${orientation_save}`, 'show');
                                layer.closeAll('loading');
                                if (reload_save_pdf_data.status) {
                                    await admin.getPdf(order_id);
                                    table.reload('print_lists', {
                                        data: admin.order_i_print_list,
                                        limit: admin.order_i_print_list.length
                                    });
                                }
                            });
                        }
                        if (obj.data.code == 'local_6') {
                            file_type_id = obj.data.id;
                            save_pdf_id = obj.data.pdf_id;
                            layer.confirm('是否需要修改内容？', {
                                btn: ['确定', '取消']
                            }, async function (index) {
                                layer.close(index);
                                print_edit_index = layer.open({
                                    type: 1,
                                    title: '保存前修改(进口司机纸)',
                                    shadeClose: true,
                                    area: admin.screen() < 2 ? ['80%', '300px'] : ['1050px', '630px'],
                                    content: $('#print_save_driver_template').html()
                                });
                                layer.load(2);
                                const driver_data = await admin.get(`/order/i/${order_id}/order_driver_paper?voy_no=${order_i_edit_data.voy_no}`);
                                    /** 货物申报司机纸打印弹窗表格渲染 */
                                    admin.getOrderDriverTable(driver_data, order_i_edit_data, '#driver_table_save')
                                auto_fn({
                                    data: admin.all_complete_data.domestic_area,
                                    listDirection: false,
                                    id: ['#district_code_name_save']
                                });
                                auto_fn({
                                    data: admin.all_complete_data.entry_clearance,
                                    listDirection: false,
                                    id: ['#i_e_port_name_save'],
                                    after: ['#i_e_port_save']
                                });
                                auto_fn({
                                    data: admin.all_complete_data.entry_clearance,
                                    listDirection: false,
                                    id: ['#i_e_port_name_two_save'],
                                    after: ['#i_e_port_two_save']
                                });
                                auto_fn({
                                    data: admin.all_complete_data.objectives_based,
                                    listDirection: false,
                                    id: ['#trade_mode_name_save']
                                });
                                auto_fn({
                                    data: admin.all_complete_data.country_area,
                                    listDirection: false,
                                    id: ['#origin_country_name_save']
                                });
                                auto_fn({
                                    data: admin.all_complete_data.entry_clearance,
                                    listDirection: false,
                                    id: ['#custom_master_name_save'],
                                    after: ['#custom_master_save']
                                });
                                $("#contr_no_save").val("").focus().val(driver_data.contr_no);
                                $("#district_code_name_save").val(driver_data.first_pros[0].district_code_name);
                                $("#i_e_port_name_save").val(driver_data.i_e_port_name);
                                $("#i_e_port_save").val(driver_data.i_e_port);
                                $("#i_e_port_name_two_save").val(driver_data.i_e_port_name);
                                $("#trade_mode_name_save").val(driver_data.trade_mode_name);
                                $("#gross_wet_save").val(driver_data.gross_wet);
                                $("#net_wt_save").val(driver_data.net_wt);
                                $("#origin_country_name_save").val(driver_data.first_pros[0].origin_country_name);
                                $("#pack_no_save").val(driver_data.pack_no);
                                $("#bill_no_save").val(driver_data.bill_no);
                                $("#car_code_save").val(driver_data.car_code);
                                $("#voy_no_save").val(driver_data.voy_no);
                                $("#custom_master_name_save").val(driver_data.custom_master_name);
                                $("#custom_master_save").val(driver_data.custom_master);
                                $("#container_ids_save").val(driver_data.container_ids);
                            }, async function (index) {
                                layer.close(index);
                                layer.load(2);
                                const reload_save_pdf_data = await admin.get(`/order/i/${order_id}/downloads/save_pdf?file_type=${obj.data.id}&pdf_id=${obj.data.pdf_id}${orientation_save}`, 'show');
                                layer.closeAll('loading');
                                if (reload_save_pdf_data.status) {
                                    await admin.getPdf(order_id);
                                    table.reload('print_lists', {
                                        data: admin.order_i_print_list,
                                        limit: admin.order_i_print_list.length
                                    });
                                }
                            });
                        }
                        if (obj.data.code == '00000002' || obj.data.code == 'local_9') {
                            layer.load(2);
                            const reload_save_pdf_data = await admin.get(`/order/i/${order_id}/downloads/save_pdf?file_type=${obj.data.id}&pdf_id=${obj.data.pdf_id}${orientation_save}`, 'show');
                            layer.closeAll('loading');
                            if (reload_save_pdf_data.status) {
                                await admin.getPdf(order_id);
                                table.reload('print_lists', {
                                    data: admin.order_i_print_list,
                                    limit: admin.order_i_print_list.length
                                });
                            }
                        }
                    }
                    break;
            }
        });
        /**打印前修改-合同-预览**/
        form.on("submit(print_edit_contract_save)", async (data) => {
            let str = '';
            for (let item in data.field) {
                str += `&${item}=${data.field[item]}`
            }
            layer.close(print_edit_index);
            window.open(encodeURI(`/order/i/${order_id}/downloads/prints?file_type=${file_type_id}${orientation_save}${str}`));
        });
        /**打印前修改-发票-预览**/
        form.on("submit(print_edit_invoice_save)", async (data) => {
            let str = '';
            for (let item in data.field) {
                str += `&${item}=${data.field[item]}`
            }
            layer.close(print_edit_index);
            window.open(`/order/i/${order_id}/downloads/prints?file_type=${file_type_id}${orientation_save}${str}`);
        });
        /**打印前修改-司机纸-预览**/
        form.on("submit(print_edit_driver_save)", async (data) => {
            let str = '';
            for (let item in data.field) {
                str += `&${item}=${data.field[item]}`
            }
            layer.close(print_edit_index);
            window.open(`/order/i/${order_id}/downloads/prints?file_type=${file_type_id}${orientation_save}${str}`);
        });
        /**保存前修改-合同-保存**/
        form.on("submit(print_save_contract_save)", async (data) => {
            let str = ''
                , save_pdf_data;
            for (let item in data.field) {
                str += `&${item}=${data.field[item]}`
            }
            layer.close(print_edit_index);
            layer.load(2);
            if (save_pdf_id) {
                save_pdf_data = await admin.get(`/order/i/${order_id}/downloads/save_pdf?file_type=${file_type_id}&pdf_id=${save_pdf_id}${orientation_save}`, 'show');
            } else {
                save_pdf_data = await admin.get(encodeURI(`/order/i/${order_id}/downloads/save_pdf?file_type=${file_type_id}${orientation_save}${str}`), 'show');
            }
            layer.closeAll('loading');
            if (save_pdf_data.status) {
                await admin.getPdf(order_id);
                table.reload('print_lists', {
                    data: admin.order_i_print_list,
                    limit: admin.order_i_print_list.length
                });
            }
        });
        /**保存前修改-发票-保存**/
        form.on("submit(print_save_invoice_save)", async (data) => {
            let str = ''
                , save_pdf_data;
            for (let item in data.field) {
                str += `&${item}=${data.field[item]}`
            }
            layer.close(print_edit_index);
            layer.load(2);
            if (save_pdf_id) {
                save_pdf_data = await admin.get(`/order/i/${order_id}/downloads/save_pdf?file_type=${file_type_id}&pdf_id=${save_pdf_id}${orientation_save}`, 'show');
            } else {
                save_pdf_data = await admin.get(encodeURI(`/order/i/${order_id}/downloads/save_pdf?file_type=${file_type_id}${orientation_save}${str}`), 'show');
            }
            layer.closeAll('loading');
            if (save_pdf_data.status) {
                await admin.getPdf(order_id);
                table.reload('print_lists', {
                    data: admin.order_i_print_list,
                    limit: admin.order_i_print_list.length
                });
            }
        });
        /**保存前修改-司机纸-保存**/
        form.on("submit(print_save_driver_save)", async (data) => {
            let str = ''
                , save_pdf_data;
            for (let item in data.field) {
                str += `&${item}=${data.field[item]}`
            }
            layer.close(print_edit_index);
            layer.load(2);
            if (save_pdf_id) {
                save_pdf_data = await admin.get(`/order/i/${order_id}/downloads/save_pdf?file_type=${file_type_id}&pdf_id=${save_pdf_id}${orientation_save}`, 'show');
            } else {
                save_pdf_data = await admin.get(encodeURI(`/order/i/${order_id}/downloads/save_pdf?file_type=${file_type_id}${orientation_save}${str}`), 'show');
            }
            layer.closeAll('loading');
            if (save_pdf_data.status) {
                await admin.getPdf(order_id);
                table.reload('print_lists', {
                    data: admin.order_i_print_list,
                    limit: admin.order_i_print_list.length
                });
            }
        });
        /**附注**/
        let order_note_index;
        $("body").on("click", "#order_note", async function () {
            if (!order_id) {
                return layer.msg("请先保存订单！")
            }
            order_note_index = layer.open({
                type: 1,
                title: '附注',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['650px', '340px'],
                content: $('#remark_note_template').html()
            });
            form.render();
            if (order_note_data) {
                $("#remark_note").val(order_note_data)
            }
            $("#order_note_dot").hide();
        });
        $("body").on("input", "#remark_note", function () {
            $("#remark_note_number span").text($(this).val().length);
        });
        /**附注保存**/
        form.on('submit(remark_note_submit)', async (data) => {
            order_note_data = data.field.remark;
            await admin.post(`/order/i/${order_id}/remark`, data.field);
            layer.close(order_note_index);
        });
        /**新建保存后刷新提取数据**/
        if ($("#order_i_edit_data").val()) {
            layer.load(2);
            order_i_edit_data = JSON.parse($("#order_i_edit_data").val());
            order_id = order_i_edit_data.id;
            laytpl($("#order_button_template").html()).render(order_i_edit_data.order_status_string, function (html) {
                $("#order_button").html(html)
            });
            laytpl($("#order_edit_bill_template").html()).render(order_i_edit_data, function (html) {
                $("#order_edit_bill").html(html)
            });
            await order_i_edit_back_filling(order_i_edit_data);
            if ($("#manual_no").val().trim()) {
                const data = await admin.get(`/order/account_manual?limit=0&search=${$("#manual_no").val()}`);
                if (data.length === 0) {
                    layer.msg(`备案号：${$("#manual_no").val()} 不存在！`, {
                        offset: '15px'
                        , icon: 2
                        , time: 1000
                        , id: 'Message'
                    });
                } else {
                    let data_detail;
                    if (data[0].is_account) {
                        data_detail = await admin.get(`/account/${data[0].id}`);
                    } else {
                        data_detail = await admin.get(`/manual/manual/${data[0].id}`);
                    }

                    admin.materials_data = data_detail.data.materials.data;
                    admin.goods_data = data_detail.data.goods.data;

                    $("#contr_item").removeAttr("disabled", "disabled");
                    $("#trade_code").attr("disabled", "disabled");
                    $("#owner_code").attr("disabled", "disabled");
                    $("#duty_mode").val("3");
                    $("#duty_mode_name").val("全免");
                    if (data[0].company_number.substring(0, 5) == '44199') {
                        $("#district_code").val('44199');
                        $("#district_code_name").val('东莞');
                    }
                }
            }
            admin.transModeControl(cusIEFlag);
            layer.closeAll('loading');
        } else {
            laytpl($("#order_button_template").html()).render('新建', function (html) {
                $("#order_button").html(html)
            });
        }
        /**导入**/
        const upload_order_i = upload.render({
            elem: '#order_i_import'
            , url: '/order/i/import'
            , accept: 'file'
            , acceptMime: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,application/vnd.ms-excel,application/vnd.ms-excel.sheet.macroEnabled.12'
            , exts: 'xlsx|xls|xlsm'
            , field: 'order_file'
            , before: function (obj) {
                layer.load(2);
            }
            , done: async (res) => {
                if (res.status) {
                    layer.msg(res.msg, {
                        offset: '15px'
                        , icon: 1
                        , time: 2000
                        , id: 'Message'
                    });
                } else {
                    layer.msg(res.msg, {
                        offset: '15px'
                        , icon: 2
                        , time: 2000
                        , id: 'Message'
                    });
                }
                layer.closeAll('loading');
                if (!order_id) {
                    order_id = res.data.id;
                    let url = document.URL;
                    let index = url.lastIndexOf("\/");
                    let url_first = url.substring(0, index + 1);
                    let url_str = url.substring(index + 1, url.length);
                    if (admin.isNumber(url_str)) {
                        history.pushState(null, null, `${url_first}${order_id}`);
                    } else {
                        history.pushState(null, null, `${url}/${order_id}`);
                    }
                }
                upload_order_i.reload({
                    data: {
                        order_id: order_id
                    }
                });

                await order_i_edit_back_filling(res.data);
                order_i_edit_data = res.data;
                if ($("#manual_no").val().trim()) {
                    const data = await admin.get(`/order/account_manual?limit=0&search=${$("#manual_no").val()}`);
                    if (data.length === 0) {
                        layer.msg(`备案号：${$("#manual_no").val()} 不存在！`, {
                            offset: '15px'
                            , icon: 2
                            , time: 1000
                            , id: 'Message'
                        });
                    } else {
                        let data_detail;
                        if (data[0].is_account) {
                            data_detail = await admin.get(`/account/${data[0].id}`);
                        } else {
                            data_detail = await admin.get(`/manual/manual/${data[0].id}`);
                        }

                        admin.materials_data = data_detail.data.materials.data;
                        admin.goods_data = data_detail.data.goods.data;

                        $("#contr_item").removeAttr("disabled", "disabled");
                        $("#trade_code").attr("disabled", "disabled");
                        $("#owner_code").attr("disabled", "disabled");
                        $("#duty_mode").val("3");
                        $("#duty_mode_name").val("全免");
                        if (data[0].company_number.substring(0, 5) == '44199') {
                            $("#district_code").val('44199');
                            $("#district_code_name").val('东莞');
                        }
                    }
                }
                admin.transModeControl(cusIEFlag);
                if (order_i_edit_data.company) {
                    $("#order_dispatch").data("user", order_i_edit_data.company.user_id);
                }
            }
            , error: function () {
                layer.closeAll('loading');
            }
        });
        /**审核通过**/
        form.on('submit(order_save_pass)', async (data) => {
            if (data.field.net_wt) {
                if (order_pros_data.length === 0) {
                    return layer.msg("存在净重，必须填写表体商品！");
                }
            }
            let is_pros = 1;
            if (order_pros_data.length > 0) {
                for (let item of order_pros_data) {
                    if ((item.g_unit ? item.g_unit != '035' : 1) && (item.first_unit ? item.first_unit != '035' : 1) && (item.second_unit ? item.second_unit != '035' : 1)) {
                        is_pros = 0;
                        break;
                    }
                }
            } else {
                is_pros = 0;
            }
            if (is_pros) {
                if (!(data.field.net_wt)) {
                    return layer.msg("存在表体商品且是千克，必填净重！");
                }
                let total = 0;
                for (let item of order_pros_data) {
                    if (item.g_unit == '035') {
                        total += parseFloat(item.g_qty) * 100000;
                    } else if (item.first_unit == '035') {
                        total += parseFloat(item.first_qty) * 100000;
                    } else if (item.second_unit == '035') {
                        total += parseFloat(item.second_qty) * 100000;
                    }
                }
                if (admin.cutZero((total / 100000).toString()) != admin.cutZero(data.field.net_wt.toString())) {
                    if (!no_weight_save) {
                        layer.open({
                            title: '商品总重量与净重不符'
                            , content: '商品总重量与净重不符！是否要保存'
                            , btn: ['是', '否']
                            , yes: function (index) {
                                no_weight_save = true;
                                layer.close(index);
                                $('#order_save').click();
                            }
                            , btn2: function () {
                                layer.msg('请重新修改净重。');
                                $('#net_wt').focus();
                            }
                        });
                        if (!no_weight_save) {
                            return false;
                        }
                    }
                }
            }
            layer.load(2);
            /**企业承诺**/
            data.field.declaratio_material_code = declaratio_material_code;
            /**其他包装**/
            if (data.field.dec_other_packs) {
                data.field.dec_other_packs = JSON.parse(data.field.dec_other_packs);
            }
            /**企业资质数据**/
            for (let item of admin.order_ent_qualif_delete_ids) {
                ent_qualif_data.push({
                    id: item
                })
            }
            data.field.quas = ent_qualif_data;
            /**使用人数据**/
            data.field.dec_users = dec_users_data;
            /**特殊业务标识数据**/
            data.field.spec_decl_flag = spec_decl_flag;
            /**表体商品数据**/
            for (let item of admin.order_quas_vin_delete_ids) {
                order_pros_data[item.pros_index].quas[item.index].order_pro_qua_vins.push({
                    id: item.id
                })
            }
            for (let item of admin.order_quas_delete_ids) {
                order_pros_data[item.index].quas.push({
                    id: item.id
                })
            }
            for (let item of admin.order_pros_delete_ids) {
                order_pros_data.push({
                    id: item
                })
            }
            data.field.pros = order_pros_data;
            /**检验检疫申报要素数据**/
            data.field.dec_request_certs = dec_request_certs_data.checkData || [];
            data.field.domestic_consignee_ename = dec_request_certs_data.domestic_consignee_ename;
            data.field.overseas_consignor_cname = dec_request_certs_data.overseas_consignor_cname;
            data.field.overseas_consignor_addr = dec_request_certs_data.overseas_consignor_addr;
            data.field.cmpl_dschrg_dt = dec_request_certs_data.cmpl_dschrg_dt;
            /**集装箱数据**/
            for (let item of admin.order_containers_delete_ids) {
                order_containers.push({
                    id: item
                })
            }
            data.field.containers = order_containers;
            /**随附单证数据**/
            for (let item of admin.order_documents_delete_ids) {
                order_documents.push({
                    id: item
                })
            }
            data.field.documents = order_documents;
            /**业务选项**/
            data.field.promise_itmes = [];
            /**特殊关系确认**/
            if (!(data.field['promise_itmes[0]'])) {
                if ($("#promise_itmes0").data('first') == '1') {
                    data.field.promise_itmes.push('0');
                } else {
                    data.field.promise_itmes.push('9');
                }
            } else {
                delete data.field['promise_itmes[0]'];
                data.field.promise_itmes.push('1');
            }
            /**价格影响确认**/
            if (!(data.field['promise_itmes[1]'])) {
                if ($("#promise_itmes1").data('first') == '1') {
                    data.field.promise_itmes.push('0');
                } else {
                    data.field.promise_itmes.push('9');
                }
            } else {
                delete data.field['promise_itmes[1]'];
                data.field.promise_itmes.push('1');
            }
            /**支付特权使用费确认**/
            if (!(data.field['promise_itmes[2]'])) {
                if ($("#promise_itmes2").data('first') == '1') {
                    data.field.promise_itmes.push('0');
                } else {
                    data.field.promise_itmes.push('9');
                }
            } else {
                delete data.field['promise_itmes[2]'];
                data.field.promise_itmes.push('1');
            }
            /**异地报关**/
            if ($("#is_other").prop("checked")) {
                data.field.is_other = 1;
            } else {
                data.field.is_other = 0;
            }

            data.field.order_status = '审核通过';
            if (!order_id) {
                const order_save_data = await admin.post(`/order/i/`, data.field);
                if (order_save_data.status) {
                    order_id = order_save_data.data.id;
                    const url = document.URL;
                    history.pushState(null, null, `${url}/${order_id}`);
                    await order_i_add_back_filling(order_save_data.data);
                }
            } else {
                const data_pass = await admin.get(`/order/i/${order_id}/audit_pass`, 'show');
            }
            layer.closeAll('loading');
        });
        /**附件**/
        $("body").on("click", "#order_enclosure", function () {
            if (!order_id) {
                return layer.msg("请先保存订单！")
            }
            layer.open({
                type: 1,
                title: '附件',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['1020px', '600px'],
                content: $('#order_enclosure_template').html()
            });
            table.render({
                elem: '#order_enclosure_lists'
                , skin: 'line'
                , url: `/order/i/${order_id}/pdf/lists?sortedBy=desc&orderBy=created_at`
                , response: {
                    countName: 'total'
                }
                , cols: [[
                    { field: 'edoc_code_name', title: '类型', width: 150 }
                    , { field: 'edoc_cop_id', title: '文件名称' }
                    , { field: 'creator', title: '操作人', width: 120 }
                    , { field: 'version', title: '版本号', width: 120 }
                    , { field: 'created_at', title: '上传时间' }
                    , { title: '操作', toolbar: '#order_enclosure_toolbar', width: 280 }
                ]]
                , page: true
                , limit: 10
            });
        });
        /**新增附件**/
        let order_enclosure_data = []
            , order_enclosure_upload
            , add_enclosure_index
            , demoListView
            , demoListView_tr;
        $("body").on("click", "#add_enclosure", function () {
            add_enclosure_index = layer.open({
                type: 1,
                title: '新增附件',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['1020px', '500px'],
                content: $('#order_enclosure_template_add').html(),
                end: function () {
                    order_enclosure_data = []
                }
            });
            /**附件上传**/
            order_enclosure_upload = upload.render({
                elem: '#order_enclosure_upload'
                , url: '/image/pdf/upload'
                , accept: 'file'
                , data: {
                    order_id: order_id
                }
                , size: 1024 * admin.UPLOAD_PDF_SIZE
                , multiple: true
                , auto: true
                , bindAction: '#order_enclosure_upload_action'
                , field: 'pdf'
                , before: function (obj) {
                    layer.load(2);
                }
                , choose: function (obj) {
                    demoListView = $('#order_enclosure_upload_list');
                    if (!edoc_code) {
                        layer.msg('请先选择类型');
                        order_enclosure_upload.config.elem.next()[0].value = '';
                        return
                    }
                    if (edoc_code != 'local_2') {
                        if (demoListView[0].childNodes.length === 1) {
                            layer.msg('只能选择一个文件');
                            order_enclosure_upload.config.elem.next()[0].value = '';
                            return
                        }
                    }
                    let files = this.files = obj.pushFile();
                    obj.preview(function (index, file, result) {
                        demoListView_tr = $(['<tr id="upload-' + index + '">'
                            , '<td>' + file.name + '</td>'
                            , '<td>' + (file.size / 1014).toFixed(1) + 'kb</td>'
                            , '<td>等待上传</td>'
                            , '<td>'
                            , '<button class="layui-btn custom-create_btn test-upload-demo-reload layui-hide" type="button">重传</button>'
                            , `<button class="layui-btn layui-btn-mini layui-btn-danger test-upload-demo-delete" data-index="${index}" type="button">删除</button>`
                            , '</td>'
                            , '</tr>'].join(''));
                        demoListView_tr.find('.test-upload-demo-reload').on('click', function () {
                            obj.upload(index, file);
                        });
                        demoListView_tr.find('.test-upload-demo-delete').on('click', function () {
                            delete files[index];
                            demoListView_tr.remove();
                            order_enclosure_upload.config.elem.next()[0].value = '';
                            order_enclosure_data.forEach((value, key) => {
                                if (value.index == $(this).data("index")) {
                                    order_enclosure_data.splice(key, 1)
                                }
                            });
                        });
                        demoListView.append(demoListView_tr);
                    });
                }
                , done: function (res, index, upload) {
                    if (res.status) {
                        layer.msg('上传成功', {
                            offset: '15px'
                            , icon: 1
                            , time: 1000
                            , id: 'Message'
                        });
                        order_enclosure_data.push({
                            index: index,
                            ...res
                        });
                        let tr = demoListView.find('tr#upload-' + index)
                            , tds = tr.children();
                        tds.eq(2).html('<span style="color: #5FB878;">上传成功</span>');
                        tds.eq(3).find('.test-upload-demo-reload').addClass('layui-hide');
                        layer.closeAll('loading');
                        return delete this.files[index];
                    } else {
                        layer.msg(res.msg, {
                            offset: '15px'
                            , icon: 2
                            , time: 1000
                            , id: 'Message'
                        });
                    }
                    layer.closeAll('loading');
                    this.error(index, upload);
                }
                , error: function (index, upload) {
                    let tr = demoListView.find('tr#upload-' + index)
                        , tds = tr.children();
                    tds.eq(2).html('<span style="color: #FF5722;">上传失败</span>');
                    tds.eq(3).find('.test-upload-demo-reload').removeClass('layui-hide');
                    layer.closeAll('loading');
                }
            });
            form.render();
        });
        /**选择附件类型**/
        form.on('select(edoc_code)', function (data) {
            const data_value = data.value.split(",");
            edoc_code = data_value[0];
            edoc_code_name = data_value[1];
        });
        /**编辑附件-选择附件类型**/
        form.on('select(edoc_code_edit)', function (data) {
            const data_value = data.value.split(",");
            edoc_code = data_value[0];
            edoc_code_name = data_value[1];
        });
        /**新建附件**/
        form.on('submit(order_enclosure_submit)', async function (data) {
            if (order_enclosure_data.length === 0) {
                layer.msg('请先上传文件');
            }
            if (edoc_code != 'local_2') {
                if (order_enclosure_data.length > 1) {
                    layer.msg('只能选择一个文件');
                    return
                }
                const form_data = [{
                    edoc_code: edoc_code,
                    edoc_code_name: edoc_code_name,
                    edoc_cop_id: order_enclosure_data[0].edoc_cop_id,
                    edoc_cop_url: order_enclosure_data[0].edoc_cop_url
                }];
                const res = await admin.post(`/order/i/${order_id}/pdf/create`, {
                    data: form_data
                });
                if (res.status) {
                    layer.close(add_enclosure_index);
                    table.reload('order_enclosure_lists');
                }
            } else {
                const form_data = {
                    edoc_code: edoc_code,
                    edoc_code_name: edoc_code_name
                };
                order_enclosure_data.forEach((value, index) => {
                    delete value.status;
                    delete value.edoc_size;
                    value.edoc_code = edoc_code;
                    value.edoc_code_name = edoc_code_name;
                });
                const res = await admin.post(`/order/i/${order_id}/pdf/create`, {
                    data: order_enclosure_data
                });
                if (res.status) {
                    layer.close(add_enclosure_index);
                    table.reload('order_enclosure_lists');
                }
            }
        });
        /**附件按钮操作**/
        let order_enclosure_data_edit = []
            , order_enclosure_upload_edit
            , edit_enclosure_index
            , demoListView_edit
            , demoListView_tr_edit
            , order_enclosure_id;
        table.on('tool(order_enclosure_lists)', async function (obj) {
            switch (obj.event) {
                case 'preview':
                    const preview_data = await admin.get(`/order/i/pdf/${obj.data.id}/show`);
                    window.open(preview_data);
                    break;
                case 'download':
                    window.open(`/order/i/pdf/${obj.data.id}/downloads`);
                    break;
                case 'edit':
                    order_enclosure_id = obj.data.id;
                    edit_enclosure_index = layer.open({
                        type: 1,
                        title: '编辑附件',
                        shadeClose: true,
                        area: admin.screen() < 2 ? ['80%', '300px'] : ['1020px', '500px'],
                        content: $('#order_enclosure_template_edit').html(),
                        end: function () {
                            order_enclosure_data_edit = []
                        }
                    });
                    /**附件上传**/
                    order_enclosure_upload_edit = upload.render({
                        elem: '#order_enclosure_upload_edit'
                        , url: '/image/pdf/upload'
                        , accept: 'file'
                        , data: {
                            order_id: order_id
                        }
                        , size: 1024 * admin.UPLOAD_PDF_SIZE
                        , multiple: true
                        , auto: true
                        , bindAction: '#order_enclosure_upload_action_edit'
                        , field: 'pdf'
                        , before: function (obj) {
                            layer.load(2);
                        }
                        , choose: function (obj) {
                            demoListView_edit = $('#order_enclosure_upload_list_edit');
                            if (!edoc_code) {
                                layer.msg('请先选择类型');
                                order_enclosure_upload_edit.config.elem.next()[0].value = '';
                                return
                            }
                            if (demoListView_edit[0].childNodes.length === 1) {
                                layer.msg('只能选择一个文件');
                                order_enclosure_upload_edit.config.elem.next()[0].value = '';
                                return
                            }
                            let files = this.files = obj.pushFile();
                            obj.preview(function (index, file, result) {
                                demoListView_tr_edit = $(['<tr id="upload-' + index + '">'
                                    , '<td>' + file.name + '</td>'
                                    , '<td>' + (file.size / 1014).toFixed(1) + 'kb</td>'
                                    , '<td>等待上传</td>'
                                    , '<td>'
                                    , '<button class="layui-btn custom-create_btn test-upload-demo-reload layui-hide" type="button">重传</button>'
                                    , `<button class="layui-btn layui-btn-mini layui-btn-danger test-upload-demo-delete" data-index="${index}" type="button">删除</button>`
                                    , '</td>'
                                    , '</tr>'].join(''));
                                demoListView_tr_edit.find('.test-upload-demo-reload').on('click', function () {
                                    obj.upload(index, file);
                                });
                                demoListView_tr_edit.find('.test-upload-demo-delete').on('click', function () {
                                    delete files[index];
                                    demoListView_tr_edit.remove();
                                    order_enclosure_upload_edit.config.elem.next()[0].value = '';
                                    order_enclosure_data_edit.forEach((value, key) => {
                                        if ($(this).data("id")) {
                                            if (value.id == $(this).data("id")) {
                                                order_enclosure_data_edit.splice(key, 1)
                                            }
                                        } else {
                                            if (value.index == $(this).data("index")) {
                                                order_enclosure_data_edit.splice(key, 1)
                                            }
                                        }
                                    });
                                });
                                demoListView_edit.append(demoListView_tr_edit);
                            });
                        }
                        , done: function (res, index, upload) {
                            if (res.status) {
                                layer.msg('上传成功', {
                                    offset: '15px'
                                    , icon: 1
                                    , time: 1000
                                    , id: 'Message'
                                });
                                order_enclosure_data_edit.push({
                                    index: index,
                                    ...res
                                });
                                let tr = demoListView_edit.find('tr#upload-' + index)
                                    , tds = tr.children();
                                tds.eq(2).html('<span style="color: #5FB878;">上传成功</span>');
                                tds.eq(3).find('.test-upload-demo-reload').addClass('layui-hide');
                                layer.closeAll('loading');
                                return delete this.files[index];
                            } else {
                                layer.msg(res.msg, {
                                    offset: '15px'
                                    , icon: 2
                                    , time: 1000
                                    , id: 'Message'
                                });
                            }
                            layer.closeAll('loading');
                            this.error(index, upload);
                        }
                        , error: function (index, upload) {
                            let tr = demoListView_edit.find('tr#upload-' + index)
                                , tds = tr.children();
                            tds.eq(2).html('<span style="color: #FF5722;">上传失败</span>');
                            tds.eq(3).find('.test-upload-demo-reload').removeClass('layui-hide');
                            layer.closeAll('loading');
                        }
                    });
                    demoListView_edit = $('#order_enclosure_upload_list_edit');
                    order_enclosure_data_edit.push(obj.data);
                    $("#edoc_code_edit").val(`${obj.data.edoc_code},${obj.data.edoc_code_name}`);
                    demoListView_tr_edit = $(['<tr id="upload-' + obj.data.id + '">'
                        , '<td>' + obj.data.edoc_cop_id + '</td>'
                        , '<td></td>'
                        , '<td>已上传</td>'
                        , '<td>'
                        , '<button class="layui-btn custom-create_btn test-upload-demo-reload layui-hide" type="button">重传</button>'
                        , `<button class="layui-btn layui-btn-mini layui-btn-danger test-upload-demo-delete_edit" data-id="${obj.data.id}" type="button">删除</button>`
                        , '</td>'
                        , '</tr>'].join(''));
                    demoListView_tr_edit.find('.test-upload-demo-delete_edit').on('click', function () {
                        demoListView_tr_edit.remove();
                        order_enclosure_upload_edit.config.elem.next()[0].value = '';
                        order_enclosure_data_edit.forEach((value, key) => {
                            if ($(this).data("id")) {
                                if (value.id == $(this).data("id")) {
                                    order_enclosure_data_edit.splice(key, 1)
                                }
                            }
                        });
                    });
                    edoc_code = obj.data.edoc_code;
                    edoc_code_name = obj.data.edoc_code_name;
                    demoListView_edit.append(demoListView_tr_edit);
                    form.render('select');
                    break;
                case 'delete':
                    layer.confirm('真的删除么', { title: '提示' }, async (index) => {
                        const data = await admin.delete(`/order/i/pdf/${obj.data.id}`);
                        layer.close(index);
                        if (data.status) {
                            table.reload('order_enclosure_lists');
                        }
                    });
                    break;
            }
        });
        /**编辑附件**/
        form.on('submit(order_enclosure_submit_edit)', async function (data) {
            if (order_enclosure_data_edit.length === 0) {
                layer.msg('请先上传文件');
            }
            if (order_enclosure_data_edit.length > 1) {
                layer.msg('只能选择一个文件');
                return
            }
            const form_data = {
                edoc_code: edoc_code,
                edoc_code_name: edoc_code_name
            };
            order_enclosure_data_edit.forEach((value, index) => {
                delete value.status;
                delete value.edoc_size;
                value.edoc_code = edoc_code;
                value.edoc_code_name = edoc_code_name;
            });
            const res = await admin.patch(`/order/i/pdf/${order_enclosure_id}/update`, order_enclosure_data_edit[0]);
            if (res.status) {
                layer.close(edit_enclosure_index);
                table.reload('order_enclosure_lists');
            }
        });
        /**办理记录**/
        $("body").on("click", "#order_take", async function () {
            if (!order_id) {
                return layer.msg("请先保存订单！")
            }
            layer.load(2);
            const data_order = await admin.get(`/order/i/${order_id}/order_logs`);
            layer.closeAll('loading');
            const data = {
                data_order: data_order
            };
            laytpl($("#order_take_template").html()).render(data, function (html) {
                $("#order_take_list").html(html)
            });
            layer.open({
                type: 1,
                title: '办理记录',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['600px', '500px'],
                content: `<div id="order_take_list_content">${$('#order_take_list').html()}</div>`
            });
        });
        /**初始值模板**/
        $("body").on("click", "#order_init", async function () {
            layer.open({
                type: 1,
                title: '初始值模板',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['850px', '580px'],
                content: $('#order_init_template').html()
            });
            table.render({
                elem: '#order_init_lists'
                , skin: 'line'
                , url: `/template/lists?sortedBy=desc&orderBy=created_at&i_e_flag=${cusIEFlag}`
                , response: {
                    countName: 'total'
                }
                , cols: [[
                    { type: 'radio' }
                    , { field: 'name', title: '模版名称', width: 150 }
                    , { field: 'trade_name', title: '境内收发货人' }
                    , {
                        field: 'creator', title: '企业信用代码', width: 120, templet(d) {
                            return d.company ? d.company.credit_code : ''
                        }
                    }
                    , {
                        field: 'version', title: '企业名称', width: 120, templet(d) {
                            return d.company ? d.company.name : ''
                        }
                    }
                    , { field: 'creator', title: '创建人' }
                ]]
                , page: true
                , limit: 5
            });
        });
        /**初始值模板搜索**/
        $("body").on("input", "#order_init_search", function () {
            table.reload('order_init_lists', {
                where: {
                    search: $(this).val()
                }
            });
        });
        /**保存初始值模板**/
        $('body').on('click', '#order_init_save', async function () {
            const data = table.checkStatus('order_init_lists').data;
            if (data.length == 0) {
                layer.msg("请选择");
                return
            }
            const id = data[0].id;
            const data_init = await admin.get(`/template/${id}/get_template`);
            await order_i_edit_back_filling(data_init);
            if ($("#manual_no").val().trim()) {
                const data = await admin.get(`/order/account_manual?limit=0&search=${$("#manual_no").val()}`);
                if (data.length === 0) {
                    layer.msg(`备案号：${$("#manual_no").val()} 不存在！`, {
                        offset: '15px'
                        , icon: 2
                        , time: 1000
                        , id: 'Message'
                    });
                } else {
                    let data_detail;
                    if (data[0].is_account) {
                        data_detail = await admin.get(`/account/${data[0].id}`);
                    } else {
                        data_detail = await admin.get(`/manual/manual/${data[0].id}`);
                    }

                    admin.materials_data = data_detail.data.materials.data;
                    admin.goods_data = data_detail.data.goods.data;

                    $("#contr_item").removeAttr("disabled", "disabled");
                    $("#trade_code").attr("disabled", "disabled");
                    $("#owner_code").attr("disabled", "disabled");
                    $("#duty_mode").val("3");
                    $("#duty_mode_name").val("全免");
                    if (data[0].company_number.substring(0, 5) == '44199') {
                        $("#district_code").val('44199');
                        $("#district_code_name").val('东莞');
                    }
                }
            }
            admin.transModeControl(cusIEFlag);
            layer.closeAll();
        });
        /**生成公路舱单**/
        $("body").on("click", "#order_create_bill", async function () {
            if (!order_id) {
                return layer.msg("请先保存订单！")
            }
            if ($("#traf_mode").val() != 4) {
                return layer.msg("请选择公路运输！")
            }
            layer.load(2);
            const data_order = await admin.get(`/order/i/${order_id}/order_logs`);
        });
        /**表单验证规则**/
        form.verify({
            gross_wet(value, item) {
                if (value.trim()) {
                    if (isNaN(value.trim())) {
                        return "毛重不足1，按1填报";
                    } else if (value.trim() < '1') {
                        return "毛重不足1，按1填报";
                    }
                }
            },
            net_wt(value, item) {
                if (value.trim()) {
                    if (parseFloat(value.trim()) > parseFloat($("#gross_wet").val().trim())) {
                        return "净重大于毛重，请确认后重新填写!";
                    }
                }
            },
            traf_mode_name(value, item) {
                if (value.trim()) {
                    const i_e_port = $("#i_e_port").val();
                    const traf_mode = $("#traf_mode").val();
                    if (!(value.trim())) {
                        return "请输入运输方式!"
                    }
                    if (i_e_port == '5301' || i_e_port == '5320' || i_e_port == '5303' || i_e_port == '5345') {
                        if (traf_mode != '4') {
                            return "运输方式必须为公路运输"
                        }
                    }
                }
            },
            enty_port_name(value, item) {
                if (value.trim()) {
                    const i_e_port = $("#i_e_port").val();
                    const enty_port_code = $("#enty_port_code").val();
                    if (!(value.trim())) {
                        return "请输入（入/出境口岸）!"
                    }
                    if (i_e_port == '5301') {
                        if (enty_port_code != '470201') {
                            return "入/出境口岸必须为470201"
                        }
                    }
                    if (i_e_port == '5320') {
                        if (enty_port_code != '470401') {
                            return "入/出境口岸必须为470401"
                        }
                    }
                    if (i_e_port == '5303') {
                        if (enty_port_code != '470501') {
                            return "入/出境口岸必须为470501"
                        }
                    }
                    if (i_e_port == '5304') {
                        if (!(enty_port_code == '470101' || enty_port_code == '470102' ||  enty_port_code == '471801')) {
                            return "入/出境口岸必须为470101或者471801或者470102"
                        }
                    }
                    if (i_e_port == '5316') {
                        if (enty_port_code != '470601') {
                            return "入/出境口岸必须为470601"
                        }
                    }
                }
            },
            dec_g_no: function (value, item) {
                let gNo = 1;
                for (let item of order_pros_data) {
                    if (item.g_no == value) {
                        gNo = 0
                    }
                }
                if (gNo) {
                    $("#dec_g_no").val("");
                    return "输入商品序号不在范围内!"
                }
                let dec_g_no = 0;
                if (order_documents[order_documents_index - 1].eco_relations.length > 0) {
                    for (let item of order_documents[order_documents_index - 1].eco_relations) {
                        if (item.dec_g_no == value) {
                            dec_g_no = 1
                        }
                    }
                    if (dec_g_no) {
                        $("#dec_g_no").val("");
                        return "报关单商品序号不能重复!"
                    }
                }
            }
        });
        /**根据屏幕等比例缩小**/
        admin.sideFlexible_window();
        /**进口报关整合申报自动完成汇总**/
        try {
            await admin.order_i_auto(auto_fn);
        } catch (err) {
            console.log(err)
        }

         /** 一般贸易禁止输入备案序号 **/
         $(check_manual_no());
         $('#manual_no').bind('input propertychange', function() {
             check_manual_no();
         });
 
         function check_manual_no() {
             if (!$("#manual_no").val()) {
                 $('#contr_item').attr('disabled', 'disabled');
                 $('#contr_item').val('');
                 form.render();
                 return
             }
         };
 
    });
    exports('order_i_change', {});
});

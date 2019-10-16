layui.define(function (exports) {
    layui.use(['form', 'admin', 'table', 'AutoComplete', 'laydate', 'laytpl', 'upload'],async function () {
        const {form, admin, table, $, laydate, laytpl, upload} = layui;
        const cusIEFlag = admin.cusIEFlag = 'I';
        let order_i_edit_data;
        /**订单ID**/
        let order_id
            /**附注数据**/
            , order_note_data
            /**获取附件类型**/
            , edoc_code
            , edoc_code_name;
        /**回车键光标跳转**/
        admin.keydown_input_textarea();
        /**自动完成方法再次封装**/
        const auto_fn = admin.auto_fn;
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
                    , {title: '操作', toolbar: '#print_toolbar', width: 280}
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
                                    laydate.render({elem: '#contact_sign_date', theme: '#1E9FFF', fixed: true});
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
                                    laydate.render({elem: '#contact_sign_date', theme: '#1E9FFF', fixed: true});
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
                                laydate.render({elem: '#contact_sign_date_save', theme: '#1E9FFF', fixed: true});
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
                                //     data: admin.all_complete_data.domestic_area,
                                //     listDirection: false,
                                //     id: ['#dest_code_name_save'],
                                //     after: ['#dest_code_save']
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
                                laydate.render({elem: '#contact_sign_date_save', theme: '#1E9FFF', fixed: true});
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
                                //     id: ['#dest_code_name_save'],
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
            console.log(order_i_edit_data);
            order_id = order_i_edit_data.id;
            if (parent.layui.admin.get_iframe_index()) {
                parent.layui.admin.get_iframe_index().find("span").append(order_i_edit_data.client_seq_no);
            }
            laytpl($("#order_i_declaration_template").html()).render(order_i_edit_data, function (html) {
                $("#order_i_declaration").html(html)
            });
            order_note_data = order_i_edit_data.remark;
            if (order_note_data) {
                $("#order_note_dot").show()
            } else {
                $("#order_note_dot").hide()
            }
            layer.closeAll('loading');
        }
        /**通过**/
        $("body").on("click", "#order_review_pass", async function () {
            if(order_i_edit_data.order_status_string == '复核不通过') {
                layer.closeAll('loading');
                return layer.msg('该订单已经驳回');
            }
            if(order_i_edit_data.order_status_string == '复核通过' || order_i_edit_data.order_status_string == '待暂存') {
                layer.closeAll('loading');
                return layer.msg('该订单已通过复核');
            }
            if(recheck_data.length > 0) {
                return layer.msg('有错误内容，无法通过！')
            }
            layer.load(2);
            const data_pass = await admin.get(`/order/i/${order_id}/pass`, 'show');
            layer.closeAll('loading');
            setTimeout(() => {
                window.location.reload();
            }, 1000);
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
                    {field: 'edoc_code_name', title: '类型', width: 150}
                    , {field: 'edoc_cop_id', title: '文件名称'}
                    , {field: 'creator', title: '操作人', width: 120}
                    , {field: 'version', title: '版本号', width: 120}
                    , {field: 'created_at', title: '上传时间'}
                    , {title: '操作', toolbar: '#order_enclosure_toolbar', width: 280}
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
                    layer.confirm('真的删除么', {title: '提示'}, async (index) => {
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
        /**驳回**/
        $("body").on("click", "#order_reject", function () {
            if(order_i_edit_data.order_status_string == '复核通过' || order_i_edit_data.order_status_string == '待暂存') {
                layer.closeAll('loading');
                return layer.msg('该订单已通过复核');
            }
            if(order_i_edit_data.order_status_string == '复核不通过') {
                layer.closeAll('loading');
                return layer.msg('该订单已经驳回');
            }
            if(recheck_data.length === 0) {
                return layer.msg('请先选择错误内容！')
            }
            layer.open({
                type: 1,
                title: '填写驳回原因',
                shadeClose: true,
                area: admin.screen() < 2 ? ['80%', '300px'] : ['650px', '340px'],
                content: $('#order_reject_template').html()
            });
            form.render();
        });
        $("body").on("input", "#order_reject_remark", function () {
            $("#remark_reject_number span").text($(this).val().length);
        });
        /**驳回保存**/
        form.on('submit(order_reject_submit)', async (data) => {
            layer.load(2);
            try {
                data.field.recheck_error_input_ids = recheck_data;
                const res = await admin.post(`/order/i/${order_id}/reject`, data.field);
                layer.closeAll('loading');
                if (res.status) {
                    setTimeout(() => {
                        layer.closeAll();
                        window.location.reload();
                    }, 1000);
                }
            } catch (e) {
                layer.closeAll('loading');
                return layer.msg('接口错误！', {
                    offset: '15px'
                    , icon: 2
                    , time: 2000
                    , id: 'Message'
                });
            }
        });
        /**点击报关单显示错误**/
        let p_timer
            , recheck_data = [];
        $("body").on("click", "p, span", function () {
            const id = $(this).data('id');
            const index = $(this).data('index');
            if(id) {
                clearTimeout(p_timer);
                p_timer = setTimeout(() => {
                    if ($(this).hasClass('is_warn')) {
                        $(this).removeClass('is_warn');
                    } else {
                        if ($(this).hasClass('is_error')) {
                            $(this).removeClass('is_error');
                            if(index) {
                                const i = recheck_data.findIndex((item) => item.index == index);
                                const x = recheck_data[i].id.findIndex((item) => item == id);
                                recheck_data[i].id.splice(x, 1);
                                if(recheck_data[i].id.length === 0) {
                                    recheck_data.splice(i, 1);
                                }
                            } else {
                                const i = recheck_data.findIndex((item) => item.id == id);
                                recheck_data.splice(i, 1);
                            }
                        } else {
                            $(this).addClass('is_error');
                            if(index) {
                                if(recheck_data.some((item) => item.index == index)) {
                                    const i = recheck_data.findIndex((item) => item.index == index);
                                    recheck_data[i].id.push(id);
                                } else {
                                    recheck_data.push({
                                        index: index,
                                        id: [id]
                                    });
                                }
                            } else {
                                recheck_data.push({
                                    id: id
                                });
                            }
                        }
                        console.log(recheck_data);
                    }
                }, 150);
            }
        });
        /**双击报关单显示警告**/
        $("body").on("dblclick", "p, span", function (event) {
            const id = $(this).data('id');
            const index = $(this).data('index');
            if(id) {
                event.stopPropagation();
                clearTimeout(p_timer);
                if ($(this).hasClass('is_error')) {
                    $(this).removeClass('is_error');
                    $(this).addClass('is_warn');
                    if(index) {
                        const i = recheck_data.findIndex((item) => item.index == index);
                        const x = recheck_data[i].id.findIndex((item) => item == id);
                        recheck_data[i].id.splice(x, 1);
                        if(recheck_data[i].id.length === 0) {
                            recheck_data.splice(i, 1);
                        }
                    } else {
                        const i = recheck_data.findIndex((item) => item.id == id);
                        recheck_data.splice(i, 1);
                    }
                } else {
                    $(this).addClass('is_warn');
                }
            }
        });
        try {
            await admin.order_i_auto(auto_fn);
        } catch (err) {
            console.log(err)
        }
    });
    exports('order_i_recheck', {});
});

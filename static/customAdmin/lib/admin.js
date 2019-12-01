layui.define('view', function (exports) {
  let $ = layui.jquery,
    laytpl = layui.laytpl,
    element = layui.element,
    setter = layui.setter,
    view = layui.view,
    device = layui.device(),
    $win = $(window),
    $body = $('body'),
    container = $('#' + setter.container),
    SHOW = 'layui-show',
    HIDE = 'layui-hide',
    THIS = 'layui-this',
    DISABLED = 'layui-disabled',
    APP_BODY = '#LAY_app_body',
    APP_FLEXIBLE = 'LAY_app_flexible',
    FILTER_TAB_TBAS = 'layadmin-layout-tabs',
    APP_SPREAD_SM = 'layadmin-side-spread-sm',
    TABS_BODY = 'layadmin-tabsbody-item',
    ICON_SHRINK = 'layui-icon-shrink-right',
    ICON_SPREAD = 'layui-icon-spread-left',
    SIDE_SHRINK = 'layadmin-side-shrink',
    SIDE_MENU = 'LAY-system-side-menu',

    //通用方法
    admin = {

      //刷新指定iframe
      reloadFrame: function (frameId) {
        parent.document.getElementById(frameId).
          contentWindow.
          location.
          reload()
      },
      UPLOAD_PDF_SIZE: 10, //pdf上传大小 M

      /**监管方式判断备案序号数据调取**/
      e_account_manual_items_arr: [
        '0700',
        '0664',
        '0300',
        '0265',
        '0657',
        '0258',
      ],

      i_account_manual_items_arr: [
        '4600',
        '4400',
      ],
      get_goods_materials_data: function (cusIEFlag) {
        /**
         * 在加工贸易手册和账册情况下 进口 需调取成品数据的贸易方式：
         * 进料成品退换 4600
         * 来料成品退换 4400
         *
         * 在加工贸易手册和账册情况下 出口 需调取料件数据的贸易方式：
         * 进料料件退换 0700
         * 进料料件复出 0664
         * 来料料件退换 0300
         * 来料料件复出 0265
         * 进料余料结转 0657
         * 来料余料结转 0258
         */
        if (
          (cusIEFlag == 'I' && $.inArray($('#trade_mode').val(),
            admin.i_account_manual_items_arr) >= 0) ||
          (cusIEFlag == 'E' && $.inArray($('#trade_mode').val(),
            admin.e_account_manual_items_arr) < 0)
        ) {
          admin.goods_materials_data = admin.goods_data
        } else {
          admin.goods_materials_data = admin.materials_data
        }

        return admin.goods_materials_data
      },

      get_ann_goods_materials_data: function (cusIEFlag) {
        /**
         * 在加工贸易手册和账册情况下 进口 需调取成品数据的贸易方式：
         * 进料成品退换 4600
         * 来料成品退换 4400
         *
         * 在加工贸易手册和账册情况下 出口 需调取料件数据的贸易方式：
         * 进料料件退换 0700
         * 进料料件复出 0664
         * 来料料件退换 0300
         * 来料料件复出 0265
         * 进料余料结转 0657
         * 来料余料结转 0258
         */
        let data
        if (
          (cusIEFlag === 'I' && $.inArray($('#supv_modecd').val(),
            admin.i_account_manual_items_arr) >= 0) ||
          (cusIEFlag === 'E' && $.inArray($('#supv_modecd').val(),
            admin.e_account_manual_items_arr) < 0)
        ) {
          data = admin.ann_goods_data
        } else {
          data = admin.ann_materials_data
        }

        return data
      },

      v: '1.2.1 std',

      //数据的异步请求
      req: view.req,

      //清除本地 token，并跳转到登入页
      exit: view.exit,

      //xss 转义
      escape: function (html) {
        return String(html || '').
          replace(/&(?!#?[a-zA-Z0-9]+;)/g, '&amp;').
          replace(/</g, '&lt;').
          replace(/>/g, '&gt;').
          replace(/'/g, '&#39;').
          replace(/"/g, '&quot;')
      },

      //事件监听
      on: function (events, callback) {
        return layui.onevent.call(this, setter.MOD_NAME, events, callback)
      },

      //table点击行选中单选框
      table_radio_click: function () {
        $(document).on('click',
          '.layui-table-body table.layui-table tbody tr, .layui-table-body table.layui-table tbody tr td',
          function (e) {
            if ($(e.target).hasClass('layui-table-col-special') ||
              $(e.target).closest('.layui-table-col-special').length) {
              return false
            }
            let index = $(this).attr('data-index'),
              tableBox = $(this).closest('.layui-table-box'),
              tableFixed = tableBox.find(
                '.layui-table-fixed.layui-table-fixed-l'),
              tableBody = tableBox.find(
                '.layui-table-body.layui-table-main'),
              tableDiv = tableFixed.length ? tableFixed : tableBody,
              checkCell = tableDiv.find('tr[data-index=' + index + ']').
                find(
                  'td div.laytable-cell-checkbox div.layui-form-checkbox i'),
              radioCell = tableDiv.find('tr[data-index=' + index + ']').
                find(
                  'td div.laytable-cell-radio div.layui-form-radio i')
            if (checkCell.length) {
              checkCell.click()
            }
            if (radioCell.length) {
              radioCell.click()
            }
          })
        $(document).on('click',
          'td div.laytable-cell-checkbox div.layui-form-checkbox, td div.laytable-cell-radio div.layui-form-radio',
          function (e) {
            e.stopPropagation()
          })
      },

      //table左右拖动
      table_mousedown: function () {
        $('body').on('mousedown', '.layui-table-main', function (event) {
          if (event.button == 0) {
            gapX = event.clientX
            startx = $(this).scrollLeft()
            $(this).on('mousemove', function (ev) {
              let left = ev.clientX - gapX
              $(this).scrollLeft(startx - left)
              return false
            })
            $(this).on('mouseup', function (et) {
              $(this).off('mousemove')
              $(this).off('mouseup')
            })
          }
        })
      },

      //进口报关整合申报提示语
      tipsJson: [
        {
          'id': 'custom_master_name',
          'name': '申报地海关：输入4位代码或名称（如‘北京海关’应输入‘0100’或‘北京海关’）',
        },
        {
          'id': 'i_e_port_name',
          'name': '进/出境关别：输入4位代码或名称（如‘北京海关’应输入‘0100’或‘北京海关’）',
        },
        {
          'id': 'manual_no',
          'name': '备案号：请输入12位备案号',
        },
        {
          'id': 'contr_no',
          'name': '合同协议号：请输入合同的全部字头和号码',
        },
        {
          'id': 'i_e_date',
          'name': '进（出）口日期：输入进（出）口日期，格式为‘年月日’，如：‘20180712’',
        },
        {
          'id': 'apl_date',
          'name': '申报日期：输入申报日期，格式为‘年月日’，如：‘20180712’',
        },
        {
          'id': 'trade_co_scc',
          'name': '境内收发货人统一社会信用代码：请输入统一社会信用代码',
        },
        {
          'id': 'trade_code',
          'name': '境内收发货人海关编码：境内收发货人在海关备案的10位代码',
        },
        {
          'id': 'trade_name',
          'name': '境内收发货人名称：输入30个字以内海关注册单位名称',
        },
        {
          'id': 'overseas_consignor_code',
          'name': '境外收发货人：对于AEO互认国家（地区）企业的，编码填报AEO编码，特殊情况下无境外收发货人的，填报‘NO’',
        },
        {
          'id': 'overseas_consignor_ename',
          'name': '境外收发货人名称（外文）：名称一般填报英文名称，检验检疫要求填报其他外文名称的，在英文名称后填报，以半角括号分隔，特殊情况下无境外收发货人的，填报‘NO’',
        },
        {
          'id': 'overseas_consignee_ename',
          'name': '境外收发货人名称（外文）：名称一般填报英文名称，检验检疫要求填报其他外文名称的，在英文名称后填报，以半角括号分隔，特殊情况下无境外收发货人的，填报‘NO’',
        },
        {
          'id': 'owner_code_scc',
          'name': '消费使用单位统一社会信用代码：请输入统一社会信用代码',
        },
        {
          'id': 'owner_code',
          'name': '消费使用单位海关编码：消费使用单位在海关备案的10位代码',
        },
        {
          'id': 'owner_name',
          'name': '消费使用单位名称：输入30个字以内海关注册单位名称',
        },
        {
          'id': 'agent_code_scc',
          'name': '申报单位统一社会信用代码：请输入统一社会信用代码',
        },
        {
          'id': 'agent_code',
          'name': '申报单位海关编码：申报单位在海关备案的10位代码',
        },
        {
          'id': 'agent_name',
          'name': '申报单位名称：输入30个字以内海关注册单位名称',
        },
        {
          'id': 'traf_mode_name',
          'name': '运输方式：输入运输代码（1位）或名称',
        },
        {
          'id': 'traf_name',
          'name': '运输工具名称：请输入运输工具名称，转关运输的格式为：@+载货清单号',
        },
        {
          'id': 'voy_no',
          'name': '航次号：根据业务类型填写运输工具的航次编号，无实际进出境的货物不填',
        },
        {
          'id': 'bill_no',
          'name': '提运单号：填报进出口货物提单或运单的编号',
        },
        {
          'id': 'trade_mode_name',
          'name': '监管方式：输入贸易代码（4位，不够请在前面补0）或名称（如，‘一般贸易’应输入‘0110’或‘一般贸易’）',
        },
        {
          'id': 'cut_mode_name',
          'name': '征免性质：输入征免性质代码（3位）或名称，可以为空',
        },
        {
          'id': 'license_no',
          'name': '许可证号：输入许可证号（许可证号格式：年-XX-顺序号，例经贸部发：00-AA-000001）',
        },
        {
          'id': 'trade_country_name',
          'name': '启运国/运抵国(地区）：输入启运国/运抵国代码（3位）或名称',
        },
        {
          'id': 'distinate_port_name',
          'name': '经停港/指运港：输入经停港/指运港代码（6位）或名称',
        },
        {
          'id': 'trans_mode_name',
          'name': '成交方式：输入成交方式代码（成交方式代码：1-CIF,2-C&F,3-FOB,4-C&I,5-市场价,6-垫仓,7-EXW）',
        },
        {
          'id': 'fee_mark_name',
          'name': '运费标志：输入运费标志：1-运费率；2-运费单价；3-运费总价',
        },
        {
          'id': 'fee_rate',
          'name': '运费：输入运费/率',
        },
        {
          'id': 'fee_curr_name',
          'name': '运费币制：输入运费币制',
        },
        {
          'id': 'insur_mark_name',
          'name': '保费标志：输入保费标志：1-保费率；3-保费总价',
        },
        {
          'id': 'insur_rate',
          'name': '保费：输入保费/率',
        },
        {
          'id': 'insur_curr_name',
          'name': '保费币制：输入保费币制',
        },
        {
          'id': 'other_mark_name',
          'name': '杂费标志：输入杂费标志：1-杂费率；3-杂费总价',
        },
        {
          'id': 'other_rate',
          'name': '杂费：输入杂费/率',
        },
        {
          'id': 'other_curr_name',
          'name': '杂费币制：输入杂费币制',
        },
        {
          'id': 'pack_no',
          'name': '件数：输入件数，不得填报0，散装货物填报1',
        },
        {
          'id': 'wrap_type_name',
          'name': '包装种类：输入包装种类（2位）或名称',
        },
        {
          'id': 'gross_wet',
          'name': '毛重：进出口货物实际毛重，计算单位为千克，不足一千克的填报为‘1’',
        },
        {
          'id': 'net_wt',
          'name': '净重：进出口货物实际净重，计算单位为千克，不足一千克的填报为‘1’',
        },
        {
          'id': 'trade_area_name',
          'name': '贸易国别(地区)：输入贸易国别（地区）代码（3位）或名称',
        },
        {
          'id': 'enty_port_name',
          'name': '入境口岸/离境口岸：输入入境口岸/离境口岸代码（6位）或名称',
        },
        {
          'id': 'goods_place',
          'name': '货物存放地点：填报货物进境后存放的场所或地点，包括海关监管作业场所、分拨仓库、定点加工厂、隔离检疫场、企业自有仓库等',
        },
        {
          'id': 'desp_port_name',
          'name': '启运港：输入启运港代码（6位）或名称',
        },
        {
          'id': 'entry_type_name',
          'name': '报关单类型：请输入报关单类型',
        },
        {
          'id': 'billTypeName',
          'name': '清单类型：请输入清单类型',
        },
        {
          'id': 'note_s',
          'name': '备注：请输入报关单的备注信息',
        },
        {
          'id': 'mark_no',
          'name': '标记唛码：填报标记唛码中除图形以外的文字、数字，无标记唛码的填报“N/M”',
        },
        {
          'id': 'org_code_name',
          'name': '检验检疫受理机关：填报提交报关单和随附单据的检验检疫机关，输入代码（6位）或名称',
        },
        {
          'id': 'vsa_org_code_name',
          'name': '领证机关：填报领取证单的检验检疫机关，输入代码（6位）或名称',
        },
        {
          'id': 'insp_org_name',
          'name': '口岸检验检疫机关：填报对入境货物实施检验检疫的检验检疫机关，输入代码（6位）或名称',
        },
        {
          'id': 'desp_date',
          'name': '启运日期：填报装载入境货物的运输工具离开启运口岸的日期',
        },
        {
          'id': 'b_l_no',
          'name': 'B/L号：填报入境货物的提货单或出库单号码。当运输方式为“航空运输”时无需填写',
        },
        {
          'id': 'purp_org_name',
          'name': '目的地检验检疫机关：需要在目的地检验检疫机关实施检验检疫的，在本栏填写对应的检验检疫机关，输入代码（6位）或名称',
        },
        {
          'id': 'correlation_no',
          'name': '关联号码：录入关联号码',
        },
        {
          'id': 'correlation_reason_flag_name',
          'name': '关联理由：在下拉菜单中选择关联报关单的关联理由',
        },
        {
          'id': 'orig_box_flag',
          'name': '原箱运输：申报使用集装箱运输的货物, 根据是否原集装箱原箱运输，勾选‘是’或‘否’',
        },
        {
          'id': 'rel_id',
          'name': '关联报关单号：输入关联报关单编号',
        },
        {
          'id': 'rel_man_no',
          'name': '关联的备案号：输入关联的备案号',
        },
        {
          'id': 'bon_no',
          'name': '保税/监管场地：输入保税或监管仓号',
        },
        {
          'id': 'cus_fie',
          'name': '场地代码：输入场地代码或名称',
        },

        {
          'id': 'g_no',
          'name': '序号：商品表体的序号，为连续的流水号，系统自动生成',
        },
        {
          'id': 'contr_item',
          'name': '备案序号：输入备案表中的商品序号，不允许修改',
        },
        {
          'id': 'code_t_s',
          'name': '商品编号：该项货物对应的商品编码',
        },
        {
          'id': 'g_name',
          'name': '商品名称：输入商品名称',
        },
        {
          'id': 'g_model',
          'name': '规格型号：输入商品的规格型号',
        },
        {
          'id': 'g_qty',
          'name': '成交数量：该项商品的成交数量，与成交单位相对应，即申报数量',
        },
        {
          'id': 'g_unit_name',
          'name': '成交计量单位：该项商品的成交时的实际计量单位',
        },
        {
          'id': 'decl_price',
          'name': '单价：该项商品的成交时的商品单位价格，即申报单价',
        },
        {
          'id': 'decl_total',
          'name': '总价：总价=单价*成交数量',
        },
        {
          'id': 'trade_curr_name',
          'name': '币制：请输入币制的代码（3位）或名称',
        },
        {
          'id': 'first_qty',
          'name': '法定第一数量：该项商品的法定成交数量，与法定单位对应',
        },
        {
          'id': 'first_unit_name',
          'name': '法定第一计量单位：该项商品的商品编码对应的海关统计第一单位，由海关决定',
        },
        {
          'id': 'exg_version',
          'name': '加工成品单耗版本号：所加工成品对应的版本号',
        },
        {
          'id': 'exg_no',
          'name': '货号：加工料件/成品货号，即企业内部的货物编号',
        },
        {
          'id': 'destination_country_name',
          'name': '最终目的国(地区)：输入最终目的国(地区)代码（3位）或名称',
        },
        {
          'id': 'second_qty',
          'name': '法定第二数量：与第二单位对应的第二成交数量',
        },
        {
          'id': 'second_unit_name',
          'name': '法定第二计量单位：该项商品的商品编码对应的海关统计第二单位，由海关决定',
        },
        {
          'id': 'origin_country_name',
          'name': '原产国(地区)：输入原产国(地区)代码（3位）或名称',
        },
        {
          'id': 'orig_place_code_name',
          'name': '原产地区：输入原产地区（3位/6位）或名称',
        },
        {
          'id': 'district_code_name',
          'name': '境内目的地/货源地：输入境内目的地/货源地五位代码',
        },
        {
          'id': 'dest_code_name',
          'name': '目的地/产地代码：输入目的地/产地代码（6位）或名称',
        },
        {
          'id': 'duty_mode_name',
          'name': '征免方式：输入征免规定，如下：1-照章；2-折半；3-全免；4-特案；5-征免性质；6-保金；7-保函；9-出口全额退税',
        },
        {
          'id': 'purpose_name',
          'name': '用途：输入用途代码（2位）或名称',
        },

        {
          'id': 'container_id',
          'name': '集装箱号：输入集装箱号',
        },
        {
          'id': 'container_md_name',
          'name': '集装箱规格：输入集装箱规格代码（2位）或名称',
        },
        {
          'id': 'lcl_flag_name',
          'name': '拼箱标识：进出口货物装运集装箱为拼箱时，在本栏下拉菜单中选择‘是’或‘否’',
        },
        {
          'id': 'goods_no',
          'name': '商品序号关系：集装箱商品序号关系信息填报单个集装箱对应的商品序号，半角逗号分隔',
        },
        {
          'id': 'goodsContaWt',
          'name': '集装箱货重：集装箱货重录入集装箱箱体自重（千克）+装载货物重量（千克）',
        },

        {
          'id': 'docu_code_name',
          'name': '随附单证代码：输入随附单证代码',
        },
        {
          'id': 'cert_code',
          'name': '随附单证编号：输入随附单证编号',
        },
      ],

      //只允许输入正整数
      decCheckInt (obj) {
        let t = obj.value.replace(/[^(\(\)\d\&\|)]/g, '')
        if (obj.value != t)
          obj.value = t
      },

      //运保杂费标志为率时币制灰掉
      markSelect (mark, code, id) {
        if ($('#' + mark).val() == '1') {
          $('#' + code).val('')
          $('#' + id).val('')
          $('#' + id).attr('disabled', 'disabled')
        } else {
          $('#' + id).removeAttr('disabled', 'disabled')
        }
      },

      /**核注清单-报关单草稿数据**/
      decListData: [],

      //进口时，成交方式是CIF/出口成交方式是FOB，则不允许录入运费和保费；
      //进口成交方式是C&I/出口成交方式是C&F，则允许录入运费，而不允许录入保费；
      //进口成交方式是C&F/出口成交方式是C&I，则不允许录入运费，而允许录入保费；否则，运费和保费都允许录入
      transModeControl: function (cusIEFlag) {
        const transMode = $('#trans_mode').val()
        if (('I' == cusIEFlag && transMode == '1') ||
          ('E' == cusIEFlag && transMode == '3')) {
          $('#fee_mark').val('')
          $('#fee_mark_name').val('')
          $('#fee_rate').val('')
          $('#fee_curr_name').val('')
          $('#fee_curr').val('')
          $('#fee_mark_name').attr('disabled', 'disabled')
          $('#fee_mark_name').removeAttr('lay-verify')

          $('#fee_curr_name').attr('disabled', 'disabled')
          $('#fee_rate').attr('disabled', 'disabled')

          $('#insur_mark').val('')
          $('#insur_mark_name').val('')
          $('#insur_rate').val('')
          $('#insur_curr_name').val('')
          $('#insur_curr').val('')
          $('#insur_mark_name').attr('disabled', 'disabled')
          $('#insur_curr_name').attr('disabled', 'disabled')
          $('#insur_rate').attr('disabled', 'disabled')
        } else if (('I' == cusIEFlag && transMode == '2') ||
          ('E' == cusIEFlag && transMode == '4')) {
          $('#fee_mark_name').val('')
          $('#fee_mark').val('')
          $('#fee_rate').val('')
          $('#fee_curr_name').val('')
          $('#fee_curr').val('')
          $('#fee_mark_name').attr('disabled', 'disabled')
          $('#fee_mark_name').removeAttr('lay-verify')

          $('#fee_curr_name').attr('disabled', 'disabled')
          $('#fee_rate').attr('disabled', 'disabled')

          $('#insur_mark_name').removeAttr('disabled', 'disabled')
          $('#insur_rate').removeAttr('disabled', 'disabled')
          if ($('#insur_mark').val() != '1') {
            $('#insur_curr_name').removeAttr('disabled', 'disabled')
          }
        } else if (('I' == cusIEFlag && transMode == '4') ||
          ('E' == cusIEFlag && transMode == '2')) {
          $('#insur_mark_name').val('')
          $('#insur_mark').val('')
          $('#insur_rate').val('')
          $('#insur_curr_name').val('')
          $('#insur_curr').val('')
          $('#insur_mark_name').attr('disabled', 'disabled')
          $('#insur_curr_name').attr('disabled', 'disabled')
          $('#insur_rate').attr('disabled', 'disabled')

          $('#fee_mark_name').removeAttr('disabled', 'disabled')
          $('#fee_mark_name').attr('lay-verify', 'required')

          if ($('#fee_mark').val() != '1') {
            $('#fee_curr_name').removeAttr('disabled', 'disabled')
          }
          $('#fee_rate').removeAttr('disabled', 'disabled')
        } else {
          $('#fee_mark_name').removeAttr('disabled', 'disabled')
          $('#fee_mark_name').attr('lay-verify', 'required')

          $('#fee_rate').removeAttr('disabled', 'disabled')
          if ($('#fee_mark').val() != '1') {
            $('#fee_curr_name').removeAttr('disabled', 'disabled')
          }

          $('#insur_mark_name').removeAttr('disabled', 'disabled')
          $('#insur_rate').removeAttr('disabled', 'disabled')
          if ($('#insur_mark').val() != '1') {
            $('#insur_curr_name').removeAttr('disabled', 'disabled')
          }
        }
      },

      //回车键focus定位    禁止textarea回车换行
      keydown_input_textarea: function () {
        $('body').on('keydown', 'textarea', function (e) {
          let self = $(this)
          let eCode = e.keyCode ? e.keyCode : e.which ? e.which : e.charCode
          if (eCode == 13) {
            e.preventDefault()
          }
        })
        $('body').on('keyup', 'input, select, textarea', function (e) {
          let self = $(this),
            form = self.parents('form:eq(0)'),
            focusable, next, prev
          let eCode = e.keyCode ? e.keyCode : e.which ? e.which :
            e.charCode
          // shift+enter 光标向上个元素移动
          if (e.shiftKey) {
            if (e.keyCode == 13) {
              // 排除只读,disabled元素
              focusable = form.find('input,a,select,textarea').
                filter(':visible').
                not(':input[readonly]').
                not(':input[disabled]')
              // focusable =
              // form.find('input,a,select,textarea').filter(':visible');
              prev = focusable.eq(focusable.index(this) - 1)

              if (prev.length) {
                if ($(this).attr('shiftEnter') == 'no') {
                  return false
                } else {
                  prev.focus()
                }
              }
              // else {
              // form.submit();
              // }
            }
          } else

          // Ctrl+enter 在textaera中换行
          if (e.ctrlKey && eCode == 13 &&
            this.localName == 'textarea') {
            let myValue = '\n'
            let $t = $(this)[0]
            if (document.selection) { // ie<9
              this.focus()
              let sel = document.selection.createRange()
              sel.text = myValue
              this.focus()
              sel.moveStart('character', -l)
              let wee = sel.text.length
            }
            // 现代浏览器
            else if ($t.selectionStart || $t.selectionStart == '0') {
              let startPos = $t.selectionStart
              let endPos = $t.selectionEnd
              let scrollTop = $t.scrollTop
              $t.value = $t.value.substring(0, startPos) +
                myValue +
                $t.value.substring(endPos,
                  $t.value.length)
              this.focus()
              // 因为myValue回车显示为\n
              $t.selectionStart = startPos + myValue.length
              $t.selectionEnd = startPos + myValue.length
              $t.scrollTop = scrollTop

            } else {
              this.value += myValue
              this.focus()
            }
          } else
          // enter 光标向下个元素移动
          if (eCode == 13) {
            if (this.localName == 'textarea') {
              e.preventDefault()
              e.stopPropagation()
            }
            // this.context.css("background-color","#b3d7f4");
            // 排除只读,disabled元素
            focusable = form.find('input,select,textarea').filter(
              ':visible').not(':input[readonly]').not(
              ':input[disabled]').not(':input[enter=-1]')
            // focusable =
            // form.find('input,select,textarea').filter(':visible');
            console.log()
            next = focusable.eq(focusable.index(this) + 1)
            // 下个元素存在
            if (next.length) {
              // console.log(this.id +" "+ next[0].id);
              // var nid = next[0].id;
              // $("#" + nid).css("background-color", "#b3d7f4");
              // $("#" + this.id).css("background-color", "");
              if ($(this).attr('enter') == 'no') {
                return false
              } else {
                next.focus()
              }
            }
            return false
          }
        })
      },

      //导入模板下载
      downloads_order_template (dom) {
        window.open(`/order/${$(dom).data('flag')}/downloads/order_template`)
      },

      //ajax-get
      get: function (url, show) {
        return new Promise(async (resolve, reject) => {
          let ajax_abort = $.ajax({
            url: url,
            type: 'get',
            dataType: 'JSON',
            success: function (res) {
              if (show) {
                if (res.status === 1) {
                  layer.msg(res.msg, {
                    offset: '15px',
                    icon: 1,
                    time: 2000,
                    id: 'Message',
                  })
                } else {
                  layer.msg(res.msg, {
                    offset: '15px',
                    icon: 2,
                    time: 2000,
                    id: 'Message',
                  })
                }
              }
              resolve(res)
            },
            error: function (error) {
              if (error.responseJSON) {
                for (let i in error.responseJSON.errors) {
                  layer.msg(error.responseJSON.errors[i].join('、'), {
                    offset: '15px',
                    icon: 2,
                    time: 2000,
                    id: 'Message',
                  })
                }
              }
              layer.closeAll('loading')
              reject(error.responseJSON)
            },
            complete: function (XMLHttpRequest, status) {
              if (status === 'timeout') {
                ajax_abort.abort()
                layer.msg('会话请求超时，请重新登录！', {
                  offset: '15px',
                  icon: 2,
                  time: 2000,
                  id: 'Message',
                })
              }
              layer.closeAll('loading')
              reject(status)
            },
          })
        })
      },

      //ajax-post
      post: function (url, data, isNotShow) {
        return new Promise(async (resolve, reject) => {
          let ajax_abort = $.ajax({
            url: url,
            type: 'POST',
            data: data,
            dataType: 'JSON',
            timeout: 8000,
            success: function (res) {
              if (!isNotShow) {
                if (res.status === 1) {
                  layer.msg(res.msg, {
                    offset: '15px',
                    icon: 1,
                    time: 1000,
                    id: 'Message',
                  })
                } else {
                  layer.msg(res.msg, {
                    offset: '15px',
                    icon: 2,
                    time: 1000,
                    id: 'Message',
                  })
                }
              }

              resolve(res)
            },
            error: function (error) {
              if (error.responseJSON) {
                for (let i in error.responseJSON.errors) {
                  layer.msg(error.responseJSON.errors[i].join('、'), {
                    offset: '15px',
                    icon: 2,
                    time: 2000,
                    id: 'Message',
                  })
                }
              }
              layer.closeAll('loading')
              reject(error.responseJSON)
            },
            complete: function (XMLHttpRequest, status) {
              if (status === 'timeout') {
                ajax_abort.abort()
                layer.msg('会话请求超时，请重新登录！', {
                  offset: '15px',
                  icon: 2,
                  time: 2000,
                  id: 'Message',
                })
              }
              layer.closeAll('loading')
              reject(status)
            },
          })
        })
      },

      //ajax-patch
      patch: function (url, data) {
        return new Promise(async (resolve, reject) => {
          let ajax_abort = $.ajax({
            url: url,
            type: 'PATCH',
            data: data,
            dataType: 'JSON',
            timeout: 8000,
            success: function (res) {
              if (res.status === 1) {
                layer.msg(res.msg, {
                  offset: '15px',
                  icon: 1,
                  time: 1000,
                  id: 'Message',
                })
              } else {
                layer.msg(res.msg, {
                  offset: '15px',
                  icon: 2,
                  time: 1000,
                  id: 'Message',
                })
              }
              resolve(res)
            },
            error: function (error) {
              if (error.responseJSON) {
                for (let i in error.responseJSON.errors) {
                  layer.msg(error.responseJSON.errors[i].join('、'), {
                    offset: '15px',
                    icon: 2,
                    time: 2000,
                    id: 'Message',
                  })
                }
              }
              layer.closeAll('loading')
              reject(error.responseJSON)
            },
            complete: function (XMLHttpRequest, status) {
              if (status === 'timeout') {
                ajax_abort.abort()
                layer.msg('会话请求超时，请重新登录！', {
                  offset: '15px',
                  icon: 2,
                  time: 2000,
                  id: 'Message',
                })
              }
              layer.closeAll('loading')
              reject(status)
            },
          })
        })
      },

      //ajax-delete
      delete: function (url) {
        return new Promise(async (resolve, reject) => {
          let ajax_abort = $.ajax({
            url: url,
            type: 'DELETE',
            headers: {
              'Content-Type': 'application/json',
              'X-HTTP-Method-Override': 'DELETE',
            },
            dataType: 'JSON',
            timeout: 8000,
            success: function (res) {
              if (res.status === 1) {
                layer.msg(res.msg, {
                  offset: '15px',
                  icon: 1,
                  time: 2000,
                  id: 'Message',
                })
              } else {
                layer.msg(res.msg, {
                  offset: '15px',
                  icon: 2,
                  time: 2000,
                  id: 'Message',
                })
              }
              resolve(res)
            },
            error: function (error) {
              if (error.responseJSON) {
                for (let i in error.responseJSON.errors) {
                  layer.msg(error.responseJSON.errors[i].join('、'), {
                    offset: '15px',
                    icon: 2,
                    time: 2000,
                    id: 'Message',
                  })
                }
              }
              layer.closeAll('loading')
              reject(error.responseJSON)
            },
            complete: function (XMLHttpRequest, status) {
              if (status === 'timeout') {
                ajax_abort.abort()
                layer.msg('会话请求超时，请重新登录！', {
                  offset: '15px',
                  icon: 2,
                  time: 2000,
                  id: 'Message',
                })
              }
              layer.closeAll('loading')
              reject(status)
            },
          })
        })
      },

      //区分进出口
      cusIEFlag: '',

      //自动完成
      async auto_fn (type) {
        let data_filter = []
        let requestData = JSON.stringify({ Limit: 5000, TypeString: type.clearanceType })
        let data = await admin.post(type.url, requestData,true)
        type.filter(data.rows, data_filter)
        //参数默认规则
        type.id.forEach((value, index) => {
          $(value).AutoComplete({
            'data': data_filter,
            'itemHeight': 20,
            'listStyle': 'custom',
            'listDirection': type.listDirection ? 'up' : 'down',
            'createItemHandler': function (index, data) {
              return `<p class="auto_list_p">${data.label}</p>`
            },
            'afterSelectedHandler': function (data) {
              if (type.after) {
                $(type.after[index]).val(data.id)
                if (type.after[index] === '#trans_mode') {
                  admin.transModeControl(admin.cusIEFlag)
                }
                if (type.after[index] === '#fee_mark') {
                  admin.markSelect('fee_mark', 'fee_curr', 'fee_curr_name')
                }
                if (type.after[index] === '#insur_mark') {
                  admin.markSelect('insur_mark', 'insur_curr',
                    'insur_curr_name')
                }
                if (type.after[index] === '#other_mark') {
                  admin.markSelect('other_mark', 'other_curr',
                    'other_curr_name')
                }
                if (type.after[index] === '#traf_mode') {
                  if ($('#traf_mode').val() === 4) {
                    //$("#bill_no").removeAttr("disabled", "disabled");
                    //启运国(地区)
                    $('#trade_country').val('HKG')
                    $('#trade_country_name').val('中国香港')
                    //经停港
                    $('#distinate_port').val('HKG003')
                    $('#distinate_port_name').val('香港（中国香港）')
                    //贸易国别（地区）
                    $('#trade_area_code').val('HKG')
                    $('#trade_area_name').val('中国香港')
                    //启运港
                    $('#desp_port_code').val('HKG003')
                    $('#desp_port_name').val('香港（中国香港）')
                  } else {
                    //$("#bill_no").attr("disabled", "disabled");
                  }
                }
                if (type.after[index] == '#trsp_modecd') {
                  if ($('#trsp_modecd').val() === 4) {
                    $('#stship_trsarv_natcd').val('110')
                    $('#stship_trsarv_natcd_name').val('中国香港')
                  }
                }

                if (type.after[index] === '#cus_fie') {
                  const value = $(type.after[index]).val()
                  if (value === '5284') {
                    $('#note_s').val('[装卸口岸：长安车检场]')
                  }
                  if (value === '5299') {
                    $('#note_s').val('[装卸口岸：其它业务]')
                  }
                  if (value === '5238') {
                    $('#note_s').val('[装卸口岸：凤岗车检场]')
                  }
                  if (value === '5298') {
                    $('#note_s').val('[装卸口岸：外关区]')
                  }
                  if (value === '5297') {
                    $('#note_s').val('[装卸口岸：加贸结转]')
                  }
                }
              }
            },
          })
        })
      },

      //核注清单商品删除存储
      annotation_items_delete_ids: [],
      //货物申报商品删除存储
      order_pros_delete_ids: [],
      //货物申报产品资质删除存储
      order_quas_delete_ids: [],
      //货物申报产品资质VIN删除存储
      order_quas_vin_delete_ids: [],
      //货物申报集装箱删除存储
      order_containers_delete_ids: [],
      //货物申报随附单证删除存储
      order_documents_delete_ids: [],
      //货物申报原产地删除存储
      order_eco_relations_delete_ids: [],
      //货物申报企业资质删除存储
      order_ent_qualif_delete_ids: [],

      //商品库
      dec_users_commodity_index: '',
      dec_users_commodity_click () {
        if (!($('#trade_code').val().trim())) {
          return layui.layer.msg('请输入境内收发货人海关编码')
        }
        admin.dec_users_commodity_index = layui.layer.open({
          type: 1,
          title: '一般贸易商品库',
          shadeClose: true,
          area: admin.screen() < 2 ? ['80%', '300px'] : ['910px', '730px'],
          content: $('#commodity_storehouse_list').html(),
        })
        const i_e_flag = admin.cusIEFlag === 'I' ? 0 : 1
        layui.table.render({
          elem: '#commodity_storehouse_table',
          autoSort: false,
          url: `/history_item/list?sortedBy=desc&orderBy=created_at&trade_code=${$(
            '#trade_code').val()}&i_e_flag=${i_e_flag}`,
          response: {
            countName: 'total',
          },
          toolbar: true,
          defaultToolbar: ['filter'],
          colFilterRecord: 'local',
          cols: [
            [
              {
                type: 'radio',
              }, {
              field: 'code_t_s',
              title: '商品编码',
              width: 160,
            }, {
              field: 'ciq_code',
              title: '商检编码',
              width: 160,
            }, {
              field: 'g_name',
              title: '商品名称',
              width: 320,
            }, {
              field: 'g_model',
              title: '规格型号',
              width: 340,
            }],
          ],
          limit: 10,
          page: true,
          height: 550,
        })
      },

      //商品库搜索
      commodity_storehouse_search (dom) {
        layui.table.reload('commodity_storehouse_table', {
          where: {
            search: $(dom).val(),
          },
          page: {
            curr: 1,
          },
        })
      },

      //商品库保存
      commodity_storehouse_save () {
        const data = layui.table.checkStatus(
          'commodity_storehouse_table').data
        if (data.length == 0) {
          layui.layer.msg('请选择商品')
          return
        }
        for (let item in data[0]) {
          console.log($(`#${item}`))
          $(`#${item}`).val(data[0][item])
        }
        layui.layer.close(admin.dec_users_commodity_index)
      },

      //选择合同备案号
      async manual_no_open () {
        layui.layer.open({
          type: 1,
          title: '选择合同备案号',
          shadeClose: true,
          area: admin.screen() < 2 ? ['80%', '300px'] : ['910px', '730px'],
          content: $('#manual_no_list').html(),
        })
        const data = await admin.get(`/order/account_manual`)
        layui.table.render({
          elem: '#manual_no_list_table',
          toolbar: true,
          defaultToolbar: ['filter'],
          colFilterRecord: 'local',
          cols: [
            [
              {
                type: 'radio',
              }, {
              field: 'name',
              title: '合同备案号',
              width: 200,
            }, {
              field: 'company_manage_name',
              title: '境内收发货人',
              width: 220,
            }, {
              field: 'company_credit_code',
              title: '企业信用代码',
              width: 220,
            }, {
              field: 'foreign_trade_company_name',
              title: '外商公司名称',
              width: 220,
            }, {
              field: 'taxationxz_code',
              title: '征免性质代码',
              width: 220,
            }, {
              field: 'taxationxz_name',
              title: '征免性质',
              width: 220,
            }, {
              field: 'company_credit_code',
              title: '收发货单位信用代码',
              width: 240,
            }, {
              field: 'company_credit_code',
              title: '消费使用单位信用代码',
              width: 240,
            }, {
              field: 'company_registration',
              title: '收发货单位CIQ编码',
              width: 240,
            }, {
              field: 'company_registration',
              title: '消费使用单位CIQ编码',
              width: 240,
            }],
          ],
          data: data,
          limit: 10,
          page: true,
          height: 550,
        })
      },

      //派单
      distribute (clickEnum, id, impexpMarkcdName) {

        /**派单**/
        $(document).on('click', clickEnum, function () {
          if (!id) {
            id = $(this).data('id')
          }
          layer.open({
            type: 1,
            title: '派单',
            shadeClose: true,
            area: admin.screen() < 2 ? ['80%', '300px'] : ['450px', '340px'],
            content: $('#distribute_template').html(),
          })

          layui.form.render()
        })

        layui.form.on('submit(distribute_submit)', async (data) => {
          let url = `/annotation/distribute/${id}`
          let result = await admin.post(url, data.field)
          if (result.status === 1) {
            setTimeout(() => {
              if (clickEnum === '#order_dispatch') {
                admin.reloadFrame(impexpMarkcdName + '核注清单iframe')
                parent.layui.admin.closeThisTabs()
              } else {
                window.location.reload()
              }

            }, 500)
            layer.closeAll()
          }
        })

      },

      //保存合同备案号
      async manual_no_save () {
        const check_data = layui.table.checkStatus(
          'manual_no_list_table').data
        if (check_data.length == 0) {
          layui.layer.msg('请选择备案号')
          return
        }
        let data
        if (check_data[0].is_account) {
          layer.load(2)
          data = await admin.get(`/account/${check_data[0].id}`)
          layer.closeAll('loading')
        } else {
          layer.load(2)
          data = await admin.get(`/manual/manual/${check_data[0].id}`)
          layer.closeAll('loading')
        }
        $('#manual_no').val('').focus().val(check_data[0].name)
        $('#foreign_company_name').
          val(check_data[0].foreign_trade_company_name)
        $('#contr_no').val(check_data[0].contr_no)

        admin.materials_data = data.data.materials.data
        admin.goods_data = data.data.goods.data

        $('#trade_co_scc').val(check_data[0].company_credit_code)
        $('#trade_code').val(check_data[0].company_number)
        $('#trade_ciq_code').val(check_data[0].company_registration)
        $('#trade_name').val(check_data[0].company_manage_name)

        $('#owner_code_scc').val(check_data[0].company_credit_code)
        $('#owner_code').val(check_data[0].company_number)
        $('#owner_ciq_code').val(check_data[0].company_registration)
        $('#owner_name').val(check_data[0].company_manage_name)

        for (let item in data.data) {
          if (item == 'supervise_mode_code') {
            if (data.data[item]) {
              $('#trade_mode').val(data.data.supervise_mode_code)
              $('#trade_mode_name').val(data.data.supervise_mode)
            } else {
              $('#trade_mode').val('')
              $('#trade_mode_name').val('')
            }
          }
          if (item == 'trade_mode_code') {
            if (data.data[item]) {
              $('#trade_mode').val(data.data.trade_mode_code)
              $('#trade_mode_name').val(data.data.trade_mode)
            } else {
              $('#trade_mode').val('')
              $('#trade_mode_name').val('')
            }
          }
          if (item == 'taxationxz_code') {
            if (data.data[item]) {
              $('#cut_mode').val(data.data.taxationxz_code)
              $('#cut_mode_name').val(data.data.taxationxz_name)
            } else {
              $('#cut_mode').val('')
              $('#cut_mode_name').val('')
            }
          }
        }

        $('#contr_item').removeAttr('disabled', 'disabled')
        $('#trade_code').attr('disabled', 'disabled')
        $('#owner_code').attr('disabled', 'disabled')
        $('#duty_mode').val('3')
        $('#duty_mode_name').val('全免')
        if (check_data[0].company_number.substring(0, 5) == '44199') {
          $('#district_code').val('44199')
          $('#district_code_name').val('东莞')
        }
        layui.layer.closeAll()
      },

      //货物申报进出境关别反填
      i_e_port_name_blur (dom) {
        setTimeout(() => {
          const i_e_port = $('#i_e_port').val()
          if (i_e_port == '5301' || i_e_port == '5320' || i_e_port ==
            '5303' || i_e_port == '5345') {
            $('#traf_mode').val('4')
            $('#traf_mode_name').val('公路运输')
            $('#trade_country').val('HKG')
            $('#trade_country_name').val('中国香港')
            $('#distinate_port').val('HKG003')
            $('#distinate_port_name').val('香港（中国香港）')
          }
          if (i_e_port == '5301') {
            $('#enty_port_code').val('470201')
            $('#enty_port_name').val('皇岗')
          }
          if (i_e_port == '5320') {
            $('#enty_port_code').val('470401')
            $('#enty_port_name').val('文锦渡')
          }
          if (i_e_port == '5303') {
            $('#enty_port_code').val('470501')
            $('#enty_port_name').val('沙头角')
          }
        }, 100)
      },

      //新建--进口报关整合申报--料件数据--成品数据
      materials_data: [],
      goods_data: [],

      //监听进出口报关整合申报备案序号
      contr_item_change: async function (dom) {
        if ($(dom).val().trim()) {
          const flag = $(dom).data('flag')
          let data = admin.get_goods_materials_data(flag)

          for (let item of data) {
            if (item.serial == $(dom).val()) {
              $('#code_t_s').val(item.hs_code)
              $('#g_name').val(item.name)
              $('#g_model').val(item.special)

              $('#g_unit').val(item.unit_one_code)
              $('#g_unit_name').val(item.unit_one)

              $('#decl_price').val(item.price)

              $('#trade_curr').val(item.moneyunit_code)
              $('#trade_curr_name').val(item.moneyunit)

              $('#first_unit').val(item.unit_two_code)
              $('#first_unit_name').val(item.unit_two)

              if (item.manuplace_code) {
                $('#origin_country').val(item.manuplace_code)
              }
              if (item.manuplace) {
                $('#origin_country_name').val(item.manuplace)
              }

              $('#g_qty').focus()
              const code_data = await admin.get(
                `/hs_code/lists?limit=0&search=${item.hs_code}`)
              if (code_data.data.length > 0) {
                if (code_data.data[0].unit2) {
                  $('#second_unit').val(code_data.data[0].unit2)
                  $('#second_unit_name').val(code_data.data[0].unit2_name)
                }
              }
            }
          }
        }
      },

      //毛重判断
      gross_wet_blur (dom) {
        if ($(dom).val().trim()) {
          if (isNaN($(dom).val().trim())) {
            $(dom).focus()
            return layer.msg('毛重不足1，按1填报')
          } else if ($(dom).val().trim() < '1') {
            $(dom).focus()
            return layer.msg('毛重不足1，按1填报')
          }
        }
      },

      //净重判断
      net_wt_blur (dom) {
        if ($(dom).val().trim()) {
          if (parseFloat($(dom).val().trim()) >
            parseFloat($('#GrossWet').val().trim())) {
            $(dom).focus()
            return layer.msg('净重大于毛重，请确认后重新填写!')
          }
        }
      },

      //监听进出口货物申报境内收发货人海关编码
      async trade_code_change (dom) {
        setTimeout(() => {
          if (!($(dom).val().trim())) {
            $('#foreign_company_name').val('')
            $('#trade_co_scc').val('')
            $('#trade_ciq_code').val('')
            $('#trade_name').val('')
          } else {
            const company_list = layui.data('company_list').data
            const index = company_list.findIndex(
              (item) => item.value == $(dom).val())
            if (company_list[index].credit_code) {
              $('#trade_co_scc').val(company_list[index].credit_code)
            }
            if (company_list[index].registration) {
              $('#trade_ciq_code').val(company_list[index].registration)
            }
            if (company_list[index].name) {
              $('#trade_name').val(company_list[index].name)
            }
            if ($('#owner_code').val().trim() == '') {
              if (company_list[index].credit_code) {
                $('#owner_code_scc').val(company_list[index].credit_code)
              }
              if (company_list[index].registration) {
                $('#owner_ciq_code').val(company_list[index].registration)
              }
              if (company_list[index].name) {
                $('#owner_name').val(company_list[index].name)
              }
              $('#owner_code').val($(dom).val())
            }
            if (company_list[index].foreign_company_name) {
              $('#foreign_company_name').
                val(company_list[index].foreign_company_name)
            }
          }
        }, 150)
      },

      //监听进出口货物申报消费使用单位海关编码
      async owner_code_change (dom) {
        setTimeout(() => {
          if (!($(dom).val().trim())) {
            $('#owner_code_scc').val('')
            $('#owner_ciq_code').val('')
            $('#owner_name').val('')
          } else {
            const company_list = layui.data('company_list').data
            const index = company_list.findIndex(
              (item) => item.value === $(dom).val())
            if (company_list[index].credit_code) {
              $('#owner_code_scc').val(company_list[index].credit_code)
            }
            if (company_list[index].registration) {
              $('#owner_ciq_code').val(company_list[index].registration)
            }
            if (company_list[index].name) {
              $('#owner_name').val(company_list[index].name)
            }
            if ($('#trade_code').val().trim() === '') {
              if (company_list[index].credit_code) {
                $('#trade_co_scc').val(company_list[index].credit_code)
              }
              if (company_list[index].registration) {
                $('#trade_ciq_code').val(company_list[index].registration)
              }
              if (company_list[index].name) {
                $('#trade_name').val(company_list[index].name)
              }
              $('#trade_code').val($(dom).val())
            }
          }
        }, 150)
      },

      //监听进出口清单经营单位海关编码
      async bizop_etpsno_change (dom) {
        if ($('#rvsngd_etps_sccd').val() == '' || $('#rcvgd_etps_nm').val() ==
          '') {
          let CompanyList = await admin.post('/company/datagrid',
            JSON.stringify({
              NameLike: $(dom).val(),
            }))
          if (CompanyList.rows.length > 0) {
            $('#bizop_etps_sccd').val(CompanyList.rows[0].CreditCode)
            $('#bizop_etps_nm').val(CompanyList.rows[0].Name)
          }
        }
      },

      //监听进出口清单加工单位海关编码
      async rcvgd_etpsno_change (dom) {
        if ($('#rvsngd_etps_sccd').val() == '' || $('#rcvgd_etps_nm').val() ==
          '') {
          let CompanyList = await admin.post('/company/datagrid',
            JSON.stringify({
              NameLike: $(dom).val(),
            }))
          if (CompanyList.rows.length > 0) {

            $('#rvsngd_etps_sccd').val(CompanyList.rows[0].CreditCode)
            $('#rcvgd_etps_nm').val(CompanyList.rows[0].Name)
          }
        }
      },

      //监听进口报关整合申报备案号
      manual_no_change
        (dom) {
        setTimeout(async () => {
          if ($(dom).val().trim()) {
            const data = await admin.get(
              `/order/account_manual?limit=0&search=${$(dom).val()}`)
            if (data.length === 0) {
              layui.layer.msg(`备案号：${$(dom).val()} 不存在！`, {
                offset: '15px',
                icon: 2,
                time: 1000,
                id: 'Message',
              })
              $('#contr_item').removeAttr('disabled', 'disabled')
              $('#trade_code').removeAttr('disabled', 'disabled')
              $('#owner_code').removeAttr('disabled', 'disabled')
            } else if (data.length > 1) {
              $('#manual_no_open').click()
              $('#manual_no_search').val($(dom).val())
              const data_manual = await admin.get(
                `/order/account_manual?search=${$(dom).val()}`)
              layui.table.reload('manual_no_list_table', {
                data: data_manual,
              })
            } else {
              let data_detail
              if (data[0].is_account) {
                data_detail = await admin.get(`/account/${data[0].id}`)
              } else {
                data_detail = await admin.get(`/manual/manual/${data[0].id}`)
              }
              const flag = $(dom).data('flag')
              $('#foreign_company_name').
                val(data[0].foreign_trade_company_name)
              $('#contr_no').val(data[0].contr_no)
              $('#trade_co_scc').val(data[0].company_credit_code)
              $('#trade_code').val(data[0].company_number)
              $('#trade_ciq_code').val(data[0].company_registration)
              $('#trade_name').val(data[0].company_manage_name)

              $('#owner_code_scc').val(data[0].company_credit_code)
              $('#owner_code').val(data[0].company_number)
              $('#owner_ciq_code').val(data[0].company_registration)
              $('#owner_name').val(data[0].company_manage_name)

              for (let item in data_detail.data) {
                if (item == 'supervise_mode_code') {
                  if (data_detail.data[item]) {
                    $('#trade_mode').
                      val(data_detail.data.supervise_mode_code)
                    $('#trade_mode_name').
                      val(data_detail.data.supervise_mode)
                  } else {
                    $('#trade_mode').val('')
                    $('#trade_mode_name').val('')
                  }
                }
                if (item == 'trade_mode_code') {
                  if (data_detail.data[item]) {
                    $('#trade_mode').val(data_detail.data.trade_mode_code)
                    $('#trade_mode_name').val(data_detail.data.trade_mode)
                  } else {
                    $('#trade_mode').val('')
                    $('#trade_mode_name').val('')
                  }
                }
                if (item == 'taxationxz_code') {
                  if (data_detail.data[item]) {
                    $('#cut_mode').val(data_detail.data.taxationxz_code)
                    $('#cut_mode_name').val(data_detail.data.taxationxz_name)
                  } else {
                    $('#cut_mode').val('')
                    $('#cut_mode_name').val('')
                  }
                }
              }

              $('#contr_item').removeAttr('disabled', 'disabled')
              $('#trade_code').attr('disabled', 'disabled')
              $('#owner_code').attr('disabled', 'disabled')

              admin.materials_data = data_detail.data.materials.data
              admin.goods_data = data_detail.data.goods.data

              $('#duty_mode').val('3')
              $('#duty_mode_name').val('全免')
              if (data[0].company_number.substring(0, 5) == '44199') {
                $('#district_code').val('44199')
                $('#district_code_name').val('东莞')
              }
            }
          } else {
            $('#trade_code').removeAttr('disabled', 'disabled')
            $('#owner_code').removeAttr('disabled', 'disabled')
          }
        }, 100)
      }
      ,

      //新建--清单--料件数据--成品数据
      /* 监听进出口清单备案序号*/
      putrec_seqno_change: async function (dom, flag, handBookId) {
        let type = 0
        if (flag === 'I') {
          type = 2
        } else if (flag === 'E') {
          type = 1
        }
        if ($(dom).val().trim()) {
          let data = {
            HandBookId: parseInt(handBookId),
            Type: type,
            Serial: $(dom).val(),
          }
          let handbookgoods = await admin.post(
            '/handbook/get_hand_book_good_by_hand_book_id',
            JSON.stringify(data), true)

          if (handbookgoods.rows.length === 1) {
            let item = handbookgoods.rows[0]
            $('#GdsMtno').val(item.RecordNo)
            $('#Gdecd').val(item.HsCode)
            $('#GdsNm').val(item.Name)
            $('#GdsSpcfModelDesc').val(item.Special)
            $('#DclCurrcd').val(item.MoneyunitCode)
            $('#DclCurrcdName').val(item.Moneyunit)
            $('#DclUnitcd').val(item.UnitOneCode)
            $('#DclUnitcdName').val(item.UnitOne)
            $('#DclUprcAmt').val(item.Price)
            $('#LawfUnitcd').val(item.UnitTwoCode)
            $('#LawfUnitcdName').val(item.UnitTwo)
            $('#EntryGdsSeqno').focus()
          } else {
            setTimeout(async () => {
              layer.msg('查询失败！', {
                offset: '15px',
                icon: 2,
                time: 1000,
                id: 'Message',
              })
            }, 100)
          }
        }
      },

      /*货物申报-监听修改成交单位*/
      g_unit_name_change (dom) {
        setTimeout(() => {
          if ($('#first_unit').val() && $('#g_unit').val() ==
            $('#first_unit').val() && $('#g_qty').val() &&
            !($('#first_qty').val())) {
            $('#first_qty').val($('#g_qty').val())
          }
          if ($('#second_unit').val() && $('#g_unit').val() ==
            $('#second_unit').val() && $('#g_qty').val() &&
            !($('#second_qty').val())) {
            $('#second_qty').val($('#g_qty').val())
          }
        }, 100)
      }
      ,

      /*货物申报-监听修改成交数量*/
      g_qty_change_index: '',
      g_qty_change (dom) {
        if (admin.cutZero($('#decl_price').val().trim()) > 0 && admin.cutZero(
          $('#decl_total').val().trim()) > 0) {
          admin.g_qty_change_index = layer.open({
            type: 1,
            title: '报关修改单价/总价？',
            shadeClose: false,
            closeBtn: 0,
            area: admin.screen() < 2 ? ['80%', '300px'] : ['380px', '180px'],
            content: `<div id="selectDeclPriceOrTotal" style="text-align: center; margin-top: 36px;"><button class="layui-btn layui-btn-primary custom-create_btn_primary order_decl_price" id="selectDeclPrice" type="button">修改单价</button><button class="layui-btn layui-btn-primary custom-create_btn_primary order_decl_price" id="selectDeclTotal" type="button">修改总价</button></div>`,
            success: function (layero, index) { //层创建完毕时触发
              $('#selectDeclPrice').focus() //默认光标定位到“修改单价”按钮上
              $('body').on('click', '#selectDeclPrice', function () {
                admin.calculationDeclPrice()
              })
              $('body').on('keydown', '#selectDeclPrice', function (e) {
                const eCode = e.keyCode ? e.keyCode : e.which
                  ? e.which
                  : e.charCode
                if (eCode == 39) { //->键
                  $('#selectDeclTotal').focus()
                } else if (eCode == 13) { //回车
                  admin.calculationDeclPrice()
                }
              })
              $('body').on('click', '#selectDeclTotal', function () {
                admin.calculationDeclTotal()
              })
              $('body').on('keydown', '#selectDeclTotal', function (e) {
                const eCode = e.keyCode ? e.keyCode : e.which
                  ? e.which
                  : e.charCode
                if (eCode == 37) { //<-键
                  $('#selectDeclPrice').focus()
                } else if (eCode == 13) { //回车
                  admin.calculationDeclTotal()
                }
              })
            },
            end: function () {
              $('#g_unit_name').focus()
            },
          })
        }
        if ($('#decl_price').val().trim() &&
          !($('#decl_total').val().trim()) && admin.cutZero(
            $('#decl_price').val().trim()) > 0) {
          admin.calculationDeclTotal()
        }
        if ($('#g_unit').val().trim()) {
          if ($('#g_unit').val().trim() == $('#first_unit').val().trim()) {
            $('#first_qty').val($(dom).val().trim())
          }
          if ($('#g_unit').val().trim() == $('#second_unit').val().trim()) {
            $('#second_qty').val($(dom).val().trim())
          }
        }
      }
      ,

      /**货物申报-计算单价*/
      calculationDeclPrice () {
        layer.close(admin.g_qty_change_index)
        let gQty = $.trim($('#g_qty').val())
        let declTotal = $.trim($('#decl_total').val())
        if (gQty && !isNaN(gQty)) {
          declPrice = admin.decToDecimal(declTotal, gQty, '4', '3', '6')
          $('#decl_price').val(declPrice)
        }
      }
      ,

      /**货物申报-计算总价*/
      calculationDeclTotal () {
        layer.close(admin.g_qty_change_index)
        let gQty = $.trim($('#g_qty').val())
        let declPrice = $.trim($('#decl_price').val())
        if (isNaN(gQty) || isNaN(declPrice)) {
          return false
        }
        declTotal = admin.decToDecimal(gQty, declPrice, '2', '2', '6')
        $('#decl_total').val(declTotal)
      }
      ,

      /*货物申报-单价修改*/
      priceChange () {
        let gQty = $.trim($('#g_qty').val())
        let declPrice = $.trim($('#decl_price').val())
        if (!gQty || isNaN(gQty) || isNaN(declPrice)) {
          return false
        }
        if (gQty && declPrice) {
          let declTotal = admin.decToDecimal(gQty, declPrice, '2', '2', '6')
          $('#decl_total').val(declTotal)
        }
      }
      ,

      /*货物申报-总价修改*/
      totalChange () {
        let gQty = $.trim($('#g_qty').val())
        let declTotal = $.trim($('#decl_total').val())

        if (!gQty) {
          return false
        }
        if (gQty && gQty != '0' && declTotal && !isNaN(gQty) &&
          !isNaN(declTotal)) {
          let declPrice = admin.decToDecimal(declTotal, gQty, '4', '3', '6')
          $('#decl_price').val(declPrice)
        }
      }
      ,

      /*清单-监听修改成交数量*/
      dcl_qty_change_index: '',
      dcl_qty_change (dom) {
        if (admin.cutZero($('#DclUprcAmt').val().trim()) > 0 && admin.cutZero(
          $('#DclTotalAmt').val().trim()) > 0) {
          admin.dcl_qty_change_index = layer.open({
            type: 1,
            title: '修改单价/总价？',
            shadeClose: false,
            closeBtn: 0,
            area: admin.screen() < 2 ? ['80%', '300px'] : ['380px', '180px'],
            content: `<div id="selectDeclPriceOrTotal_ann" style="text-align: center; margin-top: 36px;"><button class="layui-btn layui-btn-primary custom-create_btn_primary order_decl_price" id="selectDeclPrice_ann" type="button">修改单价</button><button class="layui-btn layui-btn-primary custom-create_btn_primary order_decl_price" id="selectDeclTotal_ann" type="button">修改总价</button></div>`,
            success: function (layero, index) { //层创建完毕时触发
              $('#selectDeclPrice_ann').focus() //默认光标定位到“修改单价”按钮上
              $('body').on('click', '#selectDeclPrice_ann',
                function () {
                  console.log(1)
                  admin.calculationDeclPrice_ann()
                })
              $('body').on('keydown', '#selectDeclPrice_ann',
                function (e) {
                  const eCode = e.keyCode ? e.keyCode : e.which
                    ? e.which
                    : e.charCode
                  if (eCode == 39) { //->键
                    $('#selectDeclTotal_ann').focus()
                  } else if (eCode == 13) { //回车
                    admin.calculationDeclPrice_ann()
                  }
                })
              $('body').on('click', '#selectDeclTotal_ann',
                function () {
                  admin.calculationDeclTotal_ann()
                })
              $('body').on('keydown', '#selectDeclTotal_ann',
                function (e) {
                  const eCode = e.keyCode ? e.keyCode : e.which
                    ? e.which
                    : e.charCode
                  if (eCode == 37) { //<-键
                    $('#selectDeclPrice_ann').focus()
                  } else if (eCode == 13) { //回车
                    admin.calculationDeclTotal_ann()
                  }
                })
            },
            end: function () {
              $('#LawfQty').focus()
            },
          })
        }
        if ($('#DclUprcAmt').val().trim() &&
          !($('#DclTotalAmt').val().trim()) && admin.cutZero(
            $('#DclUprcAmt').val().trim()) > 0) {
          admin.calculationDeclTotal_ann()
        }
        if ($('#DclUnitcd').val().trim()) {
          if ($('#DclUnitcd').val().trim() == $('#LawfUnitcd').val().trim()) {
            $('#LawfQty').val($(dom).val().trim())
          }
          if ($('#DclUnitcd').val().trim() ==
            $('#SecdLawfUnitcd').val().trim()) {
            $('#SecdLawfQty').val($(dom).val().trim())
          }
        }
      }
      ,

      /**清单-计算单价*/
      calculationDeclPrice_ann () {
        layer.close(admin.dcl_qty_change_index)
        let gQty = $.trim($('#DclQty').val())
        let declTotal = $.trim($('#DclTotalAmt').val())
        if (gQty && !isNaN(gQty)) {
          declPrice = admin.decToDecimal(declTotal, gQty, '4', '3', '6')
          $('#DclUprcAmt').val(declPrice)
        }
      }
      ,

      /**清单-计算总价*/
      calculationDeclTotal_ann () {
        layer.close(admin.dcl_qty_change_index)
        let gQty = $.trim($('#DclQty').val())
        let declPrice = $.trim($('#DclUprcAmt').val())
        if (isNaN(gQty) || isNaN(declPrice)) {
          return false
        }
        declTotal = admin.decToDecimal(gQty, declPrice, '2', '2', '6')
        $('#DclTotalAmt').val(declTotal)
      },

      /**清单-单价修改*/
      priceChange_ann () {
        let gQty = $.trim($('#DclQty').val())
        let declPrice = $.trim($('#DclUprcAmt').val())
        if (!gQty || isNaN(gQty) || isNaN(declPrice)) {
          return false
        }
        if (gQty && declPrice) {
          let declTotal = admin.decToDecimal(gQty, declPrice, '2', '2', '6')
          $('#DclTotalAmt').val(declTotal)
        }
      }
      ,

      /*清单-总价修改*/
      totalChange_ann () {
        let gQty = $.trim($('#DclQty').val())
        let declTotal = $.trim($('#DclTotalAmt').val())
        if (!gQty) {
          return false
        }
        if (gQty && gQty != '0' && declTotal && !isNaN(gQty) &&
          !isNaN(declTotal)) {
          let declPrice = admin.decToDecimal(declTotal, gQty, '4', '3', '6')
          $('#DclUprcAmt').val(declPrice)
        }
      }
      ,

      /**
       * 四舍六入五成双
       * @param num
       * @param digit   小数点多少位
       * calculationType：计算类型（加减乘除对应0,1,2,3）
       * @returns {Number}
       */
      decToDecimal (num1, num2, digit, calculationType, roundingMode) {
        let calculationDatas = [num1 + ',' + num2]
        let digits = [digit]
        let calculationTypes = [calculationType]
        let roundingModes = [roundingMode]
        let resultList = admin.decCalculation(calculationDatas, digits,
          calculationTypes, roundingModes, '1')
        return resultList[0]
      }
      ,

      /**
       *
       * @param calculations    计算值的集合
       * @param calculationType    计算类型,0,1,2,3分别对应加减乘除
       * @param roundingMode    计算结果小数保留的方式，四舍五入等
       * @param digit    保留小数的位数
       * @param isZeroNoShow    小数点后全为0是否显示，1显示，0不显示
       */
      decCalculation (
        calculationDatas, digits, calculationTypes, roundingModes,
        isZeroNoShow) {
        let resultList = []
        if (!calculationDatas || calculationDatas.length < 1) {
          return resultList
        }

        if (!isZeroNoShow) {
          isZeroNoShow = 1
        }
        for (var i = 0; i < calculationDatas.length; i++) {
          if (!calculationDatas[i]) {
            return resultList
          }
          let calculationData = calculationDatas[i].split(',')
          let calculationDataMap = admin.calculation(calculationData,
            calculationTypes[i],
            roundingModes[i], digits[i], isZeroNoShow)
          if (calculationDataMap)
            resultList[resultList.length] = calculationDataMap.result
        }
        return resultList
      }
      ,

      /**
       *
       * @param calculations    计算值的集合
       * @param calculationType    计算类型,0,1,2,3分别对应加减乘除
       * @param roundingMode    计算结果小数保留的方式，四舍五入等
       * @param digit    保留小数的位数
       * @param isZeroNoShow    小数点后全为0是否显示，1显示，0不显示
       */
      calculation (
        calculations, calculationType, roundingMode, digit, isZeroNoShow) {
        let resultData = {}
        let result = ''
        if (!calculations || calculations.length < 1 || !calculationType ||
          calculationType.length < 0) {
          resultData.result = result
          return resultData
        }
        if (!digit) {
          digit = '0'
        }
        if (!isZeroNoShow) {
          digit = '1'
        }
        let bigDecimalResult = new BigDecimal(calculations[0])
        let bigDecimal = null
        let calculation = ''
        for (var i = 1; i < calculations.length; i++) {
          if (!calculations[i]) {
            calculation = '0'
          } else {
            calculation = calculations[i]
          }
          bigDecimal = new BigDecimal(calculation)
          if (calculationType == '0') {
            bigDecimalResult = bigDecimalResult.add(bigDecimal)
          } else if (calculationType == '1') {
            bigDecimalResult = bigDecimalResult.subtract(bigDecimal)
          } else if (calculationType == '2') {
            bigDecimalResult = bigDecimalResult.multiply(bigDecimal).
              setScale(parseInt(digit), parseInt(roundingMode))
          } else if (calculationType == '3') {
            bigDecimal = new BigDecimal(calculations[i])
            if (calculation == '0') {
              resultData.result = ''
              return resultData
            } else {
              bigDecimalResult = bigDecimalResult.divide(bigDecimal,
                parseInt(digit),
                parseInt(roundingMode))
            }
          } else {
            resultData.result = ''
            return resultData
          }
        }
        if (isZeroNoShow == '1' && parseInt(bigDecimalResult) ==
          bigDecimalResult.toString()) {
          result = parseInt(bigDecimalResult)
        } else {
          result = bigDecimalResult.toString()
        }
        resultData.result = result
        return resultData
      }
      ,

      /*选中商品编码申报要素*/
      declaration_data: '',

      /*标志是否已经打开了窗口*/
      isOpenGoodsWindow:
        false,
      isOpenGoodsElementWindow:
        false,

      /*商品规范申报--商品申报要素*/
      async openMerchElement () {
        const declaration_data_array = admin.declaration_data.split(';')
        const data = declaration_data_array.map((item, index) => {
          if (index < 9) {
            return item.substr(1)
          } else {
            return item.substr(2)
          }
        })
        layui.laytpl($('#declaration_template').html()).
          render(data, function (html) {
            $('#declaration_list').html(html)
          })
        layer.open({
          type: 1,
          title: '商品规范申报-商品申报要素',
          shadeClose: true,
          area: admin.screen() < 2 ? ['80%', '300px'] : ['1180px', '600px'],
          content: $('#declaration_list').html(),
          end: function () {
            $('#g_qty').focus()
          },
        })
        $('body #val0Name').focus()
        $('#selectCodeTs').val($('#g_name').val())
        let brand_type = await admin.post(`/clearance/no_paginate`,
          JSON.stringify({ Type: '品牌类型' }))
        let data_filter = []
        for (let item of brand_type) {
          data_filter.push({
            id: item.customs_code,
            label: `${item.customs_code}-${item.name}`,
            value: `${item.name}`,
          })
        }
        admin.brand_type_data = data_filter
        $('#val0Name').AutoComplete({
          'data': data_filter,
          'width': 280,
          'itemHeight': 20,
          'listStyle': 'custom',
          'listDirection': 'down',
          'createItemHandler': function (index, data) {
            return `<p class="auto_list_p">${data.label}</p>`
          },
          'afterSelectedHandler': function (data) {
            $('#val0').val(data.id)
          },
        })

        let export_benefits = await admin.post(`/clearance/no_paginate`,
          { Type: '出口享惠情况' })
        let data_filter_benefits = []
        for (let item of export_benefits) {
          data_filter_benefits.push({
            id: item.customs_code,
            label: `${item.customs_code}-${item.name}`,
            value: `${item.name}`,
          })
        }
        admin.export_benefits_data = data_filter_benefits
        $('#val1Name').AutoComplete({
          'data': data_filter_benefits,
          'width': 280,
          'itemHeight': 20,
          'listStyle': 'custom',
          'listDirection': 'down',
          'createItemHandler': function (index, data) {
            return `<p class="auto_list_p">${data.label}</p>`
          },
          'afterSelectedHandler': function (data) {
            $('#val1').val(data.id)
          },
        })
      }
      ,

      /*编辑商品申报要素-规格型号*/
      elements_change (dom) {
        const data = $(dom).val().split('|')
        for (let item in data) {
          if (item == 0) {
            const data_filter = admin.brand_type_data.filter(
              t => t.id == data[item])
            $('#val0').val(data_filter[0].id)
            $('#val0Name').val(data_filter[0].value)
          } else if (item == 1) {
            const data_filter = admin.export_benefits_data.filter(
              t => t.id == data[item])
            $('#val1').val(data_filter[0].id)
            $('#val1Name').val(data_filter[0].value)
          } else {
            $(`#val${item}`).val(data[item])
          }
        }
      }
      ,

//非危险化学品/集装箱规格
      chemicals_data: [
        {
          id: '0',
          label: '0-否',
          value: '否',
        }, {
          id: '1',
          label: '1-是',
          value: '是',
        }],
      //危包类别
      category_data:
        [
          {
            id: '1',
            label: '1-一类',
            value: '一类',
          }, {
          id: '2',
          label: '2-二类',
          value: '二类',
        }, {
          id: '3',
          label: '3-三类',
          value: '三类',
        }],

      //历史申报要素数据
      elements_his_data:
        [],
      //品牌类型数据
      brand_type_data:
        [],
      //出口享惠情况数据
      export_benefits_data:
        [],

      //规范要素回车事件
      decFocus (e, id) {
        let fanalVal = ''
        let surtax = ''
        let surtaxIndex = -1
        let objList = document.getElementsByName('dyanInput')
        for (var i = 0; i < objList.length; i++) {
          if (objList[i].getAttribute('decElemType') == '9') {
            fanalVal += '<' + objList[i].value + '>'
            surtaxIndex = i
          } else {
            if (surtaxIndex > -1) {
              fanalVal += '|'
            }
            if (i != objList.length - 1) {
              fanalVal += objList[i].value + '|'
            } else {
              fanalVal += objList[i].value
            }
          }
        }
        let flag = 0
        if (fanalVal && fanalVal.length > 0) {
          for (var j = fanalVal.length; j > 0; j--) {
            if (fanalVal[j - 1] != '|') {
              flag = j
              break
            }
          }
        }
        let str = fanalVal.substring(0, flag)
        $('#elements').val(str)
        if (!admin.checkGModel(str)) {
          return false
        }
        const eCode = e.keyCode ? e.keyCode : e.which ? e.which : e.charCode
        if (e.shiftKey != 1 && eCode == 13) {
          if (id == '0') {
            $('#declaration_save').click()
          }
        }
      }
      ,

//规格型号校验
      checkGModel (str) {
        var charCode
        var realLength = 0
        if (str) {
          for (var i = 0; i < str.length; i++) {
            charCode = str.charCodeAt(i)
            if (charCode >= 0 && charCode <= 128) {
              realLength += 1
            } else {
              // 如果是中文则长度加2
              realLength += 2
            }
            document.getElementById(
              'elementsByteTotal').innerHTML = realLength
          }
        } else {
          document.getElementById('elementsByteTotal').innerHTML = 0
        }
        if (realLength > 255) {
          $('#elementsByteTotal').focus()
          layer.open({
            title: '提示信息',
            content: '规范申报要素归类后，规格型号数据<br>' + str + '<br>超长，大于255字节',
          })
          return false
        }
        return true
      }
      ,

//表体商品编号弹窗
      async openGoodsWindow (e) {
        if (admin.isOpenGoodsWindow) {
          return false
        }
        const eCode = e.keyCode ? e.keyCode : e.which ? e.which : e.charCode
        if (eCode == 13 || !eCode) {
          // 输入商品编号大于等于4位时才弹商品列表框
          const codeTs = $('#code_t_s').val()
          const data = await admin.get(
            `/hs_code/lists?limit=0&search=${codeTs}`)
          if (data.data.length === 0) {
            layer.msg('无此商品编码')
            return
          }
          if (codeTs.length >= 4) {
            admin.isOpenGoodsWindow = true //标志窗口已经打开
            layer.open({
              type: 1,
              title: '商品列表',
              shadeClose: true,
              area: admin.screen() < 2 ? ['80%', '300px'] : [
                '910px',
                '480px'],
              content: $('#code_t_s_list').html(),
              success: function (layero, index) {
                this.enterEsc = function (event) {
                  if (event.keyCode === 13) {
                    $('#code_t_s_save').click()
                    return false
                  }
                }
                $(document).on('keydown', this.enterEsc)
              },
              end: function () {
                admin.isOpenGoodsWindow = false
                $(document).off('keydown', this.enterEsc)
                admin.openMerchElement()
              },
            })
            layui.table.render({
              elem: '#code_t_s_list_table',
              toolbar: true,
              defaultToolbar: ['filter'],
              colFilterRecord: 'local',
              cols: [
                [
                  {
                    type: 'radio',
                  }, {
                  field: 'code',
                  title: '商品编号',
                  width: 180,
                }, {
                  field: 'name',
                  title: '商品名称',
                  width: 320,
                }, {
                  field: 'remark',
                  title: '备注',
                }],
              ],
              data: data.data,
              limit: data.data.length,
              height: 300,
            })
            $(`.layui-table-view[lay-id='code_t_s_list_table'] .layui-table-body tr[data-index='0'] .layui-form-radio`).
              click()
          } else if (!codeTs) {
            return false
          } else {
            layer.msg('请至少输入四位商品编码！')
            $(`input[id="code_t_s"]`).focus()
            return false
          }
        }
      }
      ,

//编辑检验检疫货物规格回车
      changeFoucsTogoodsAttr (event) {
        const eCode = event.keyCode ? event.keyCode : event.which
          ? event.which
          : event.charCode
        if (event.shiftKey != 1 && eCode == 13) {
          $('#goods_spec_save').click()
          layer.closeAll()
        }
      }
      ,

//编辑货物危险信息回车
      changeFoucsDanger (event) {
        const eCode = event.keyCode ? event.keyCode : event.which
          ? event.which
          : event.charCode
        if (event.shiftKey != 1 && eCode == 13) {
          $('#dang_save').click()
          layer.closeAll()
        }
      }
      ,

//tips计算总价/成交数量合计/法定第一数量/法定第二数量
      is_total_number (order_pros_data) {
        let totalPrice = 0,
          totalGQty = 0,
          totalQty1 = 0,
          totalQty2 = 0
        for (let item of order_pros_data) {

          if (item.decl_total_string != '' && item.decl_total_string !=
            'null' && item.decl_total_string) {
            totalPrice += parseFloat(item.decl_total_string) * 100000
          }
          if (item.g_qty_string != '' && item.g_qty_string != 'null' &&
            item.g_qty_string) {
            totalGQty += parseFloat(item.g_qty_string) * 100000
          }
          if (item.first_qty_string != '' && item.first_qty_string !=
            'null' && item.first_qty_string) {
            totalQty1 += parseFloat(item.first_qty_string) * 100000
          }
          if (item.second_qty_string != '' && item.second_qty_string !=
            'null' && item.second_qty_string) {
            totalQty2 += parseFloat(item.second_qty_string) * 100000
          }

        }

        $('#totalPrice').
          text(admin.cutZero((totalPrice / 100000).toString()))
        $('#totalGQty').text(admin.cutZero((totalGQty / 100000).toString()))
        $('#totalQty1').text(admin.cutZero((totalQty1 / 100000).toString()))
        $('#totalQty2').text(admin.cutZero((totalQty2 / 100000).toString()))
      },

      async base_clearance_data_auto (obj) {
        const url = `/clearance/datagrid`
        //货物
        await admin.auto_fn({
          url: url,
          clearanceType: obj.name,
          listDirection: false,
          id: obj.id,
          after: obj.after,
          filter: function (data, data_filter) {
            for (let item of data) {
              let value = `${item.Name}`
              // if (`${item.ShortName}` !== 'null') {
              //   value = `${item.ShortName}`
              // } else if (`${item.Name}` !== 'null') {
              //   value = `${item.Name}`
              // }

              let label = `${item.CustomsCode}-${value}`
              let data_filter_id = item.CustomsCode
              if (obj.filter_type === 'l') {
                label = `<span class="auto_list_p_left">${item.CustomsCode}-${value}</span><span class="auto_list_p_right"><i>${item.OldCustomCode}</i><i>${item.OldCiqCode}</i></span>`
              }else if(obj.filter_type === 'anns'){
                label = `${item.OldCustomCode}-${value}`
                data_filter_id = item.OldCustomCode
              }

              data_filter.push({
                id:data_filter_id,
                label: label,
                value: value,
              })

            }
          },
        })
      },

      /* 通关参数加载 */
      async clearance_data_auto (loadArray) {
        if ($.inArray('list_types', loadArray) !== -1) {
          //`${item.customs_code}-${item.name}`
          let list_types = {
            name: '清单类型',
            filter_type: 's',
            id: ['#InvtTypeName'],
            after: ['#invt_type'],
          }
          await admin.base_clearance_data_auto(list_types)
        }

        if ($.inArray('finished_product', loadArray) !== -1) {
          //`${item.customs_code}-${item.name}`
          let finished_product = {
            name: '料件成品标记',
            filter_type: 's',
            id: ['#mtpck_endprd_markcd_name'],
            after: ['#mtpck_endprd_markcd'],
          }
          await admin.base_clearance_data_auto(finished_product)
        }

        if ($.inArray('types_transfer', loadArray) !== -1) {
          //`${item.customs_code}-${item.name}`
          let types_transfer = {
            name: '流转类型',
            filter_type: 's',
            id: ['#list_type_name'],
            after: ['#list_type'],
          }
          await admin.base_clearance_data_auto(types_transfer)
        }
        if ($.inArray('nuclear_declaration_lis', loadArray) !== -1) {
          //`${item.customs_code}-${item.name}`
          let nuclear_declaration_lis = {
            name: '清单报关标志',
            filter_type: 's',
            id: ['#dclcus_flag_name'],
            after: ['#dclcus_flag'],
          }
          await admin.base_clearance_data_auto(nuclear_declaration_lis)
        }
        if ($.inArray('customs_declaration_type', loadArray) !== -1) {
          //`${item.customs_code}-${item.name}`
          let customs_declaration_type = {
            name: '报关类型',
            filter_type: 's',
            id: ['#dclcus_typecd_name'],
            after: ['#dclcus_typecd'],
          }
          await admin.base_clearance_data_auto(customs_declaration_type)
        }

        if ($.inArray('type_declaration_list', loadArray) !== -1) {
          //`${item.customs_code}-${item.name}`
          let type_declaration_list = {
            name: '清单报关单类型',
            filter_type: 's',
            id: ['#dec_type_name'],
            after: ['#dec_type'],
          }
          await admin.base_clearance_data_auto(type_declaration_list)
        }

        if ($.inArray('gen_dec_flag_list', loadArray) !== -1) {
          //`${item.customs_code}-${item.name}`
          let gen_dec_flag_list = {
            name: '生成报关单标志',
            filter_type: 's',
            id: ['#gen_dec_flag_name'],
            after: ['#gen_dec_flag'],
          }
          await admin.base_clearance_data_auto(gen_dec_flag_list)
        }

        if ($.inArray('modf_markcd_list', loadArray) !== -1) {
          //`${item.customs_code}-${item.name}`
          let modf_markcd_list = {
            name: '清单表体修改标志',
            filter_type: 's',
            id: ['#ModfMarkcdName'],
            after: ['#ModfMarkcd'],
          }
          await admin.base_clearance_data_auto(modf_markcd_list)
        }

        if ($.inArray('entry_clearance', loadArray) !== -1) {
          let entry_clearance = {
            name: '关区代码',
            filter_type: null,
            id: [
              '#custom_master_name',
              '#i_e_port_name',
            ],
            after: [
              '#custom_master',
              '#i_e_port',
            ],
          }
          await admin.base_clearance_data_auto(entry_clearance)
        }

        //清单
       if ($.inArray('ann_entry_clearance', loadArray) !== -1) {
                let ann_entry_clearance = {
                  name: '关区代码',
                  filter_type: "s",
                  id: [
                    '#impexp_portcd_name',
                    '#dcl_plc_cuscd_name',
                  ],
                  after: [
                    '#impexp_portcd',
                    '#dcl_plc_cuscd',
                  ],
                }
                await admin.base_clearance_data_auto(ann_entry_clearance)
              }

        if ($.inArray('mode_shipping', loadArray) !== -1) {
          let mode_shipping = {
            name: '运输方式代码',
            filter_type: 'l',
            id: ['#traf_mode_name'],
            after: ['#traf_mode'],
          }
          await admin.base_clearance_data_auto(mode_shipping)
        }

        if ($.inArray('ann_mode_shipping', loadArray) !== -1) {
          let ann_mode_shipping = {
            name: '运输方式代码',
            filter_type: 's',
            id: [ '#trsp_modecd_name'],
            after: [ '#trsp_modecd'],
          }
          await admin.base_clearance_data_auto(ann_mode_shipping)
        }

        if ($.inArray('objectives_based', loadArray) !== -1) {
          let objectives_based = {
            name: '监管方式代码',
            filter_type: 'l',
            id: [
              '#trade_mode_name',
            ],
            after: [
              '#trade_mode',
            ],
          }
          await admin.base_clearance_data_auto(objectives_based)
        }

        if ($.inArray('ann_objectives_based', loadArray) !== -1) {
          let ann_objectives_based = {
            name: '监管方式代码',
            filter_type: 's',
            id: [
              '#supv_modecd_name',
            ],
            after: [
              '#supv_modecd',
            ],
          }
          await admin.base_clearance_data_auto(ann_objectives_based)
        }

        if ($.inArray('nature_exemption', loadArray) !== -1) {
          let nature_exemption = {
            name: '征免性质代码',
            filter_type: "s",
            id: ['#cut_mode_name'],
            after: ['#cut_mode'],
          }
          await admin.base_clearance_data_auto(nature_exemption)
        }

        if ($.inArray('country_area', loadArray) !== -1) {
          let country_area = {
            name: '国别地区代码',
            filter_type: 'l',
            id: [
              '#trade_country_name',
              '#trade_area_name',
              '#destination_country_name',
              '#origin_country_name',
            ],
            after: [
              '#trade_country',
              '#trade_area_code',
              '#destination_country',
              '#origin_country',
            ],
          }
          await admin.base_clearance_data_auto(country_area)
        }

        if ($.inArray('ann_country_area', loadArray) !== -1) {
          let ann_country_area = {
            name: '国别地区代码',
            filter_type: 'anns',
            id: [
              '#stship_trsarv_natcd_name',
              '#DestinationNatcdName',
              '#NatcdName',
            ],
            after: [
              '#stship_trsarv_natcd',
              '#DestinationNatcd',
              '#Natcd',
            ],
          }
          await admin.base_clearance_data_auto(ann_country_area)
        }

        if ($.inArray('harbour', loadArray) !== -1) {
          let harbour = {
            name: '港口代码',
            filter_type: 'l',
            id: ['#distinate_port_name', '#desp_port_name'],
            after: ['#distinate_port', '#desp_port_code'],
          }
          await admin.base_clearance_data_auto(harbour)
        }

        if ($.inArray('terms_delivery', loadArray) !== -1) {
          let terms_delivery = {
            name: '成交方式代码',
            filter_type: 's',
            id: ['#trans_mode_name'],
            after: ['#trans_mode'],
          }
          await admin.base_clearance_data_auto(terms_delivery)
        }

        if ($.inArray('cost_tag', loadArray) !== -1) {
          let cost_tag = {
            name: '费用标记',
            filter_type: 's',
            id: ['#fee_mark_name', '#insur_mark_name', '#other_mark_name'],
            after: ['#fee_mark', '#insur_mark', '#other_mark'],
          }
          await admin.base_clearance_data_auto(cost_tag)
        }

        if ($.inArray('currency', loadArray) !== -1) {
          let currency = {
            name: '货币代码',
            filter_type: 'l',
            id: [
              '#fee_curr_name',
              '#insur_curr_name',
              '#other_curr_name',
              '#trade_curr_name',
            ],
            after: [
              '#fee_curr',
              '#insur_curr',
              '#other_curr',
              '#trade_curr',
            ],
          }
          await admin.base_clearance_data_auto(currency)
        }

        if ($.inArray('ann_currency', loadArray) !== -1) {
          let ann_currency = {
            name: '货币代码',
            filter_type: 'anns',
            id: [
              '#DclCurrcdName',
            ],
            after: [
              '#DclCurrcd',
            ],
          }
          await admin.base_clearance_data_auto(ann_currency)
        }

        if ($.inArray('kind_packages', loadArray) !== -1) {
          let kind_packages = {
            name: '包装种类代码',
            filter_type: 'l',
            id: ['#wrap_type_name'],
            after: ['#wrap_type'],
          }
          await admin.base_clearance_data_auto(kind_packages)
        }

        if ($.inArray('domestic_ports', loadArray) !== -1) {
          let domestic_ports = {
            name: '国内口岸代码',
            filter_type: 's',
            id: ['#enty_port_name'],
            after: ['#enty_port_code'],
          }
          await admin.base_clearance_data_auto(domestic_ports)
        }

        if ($.inArray('types_customs', loadArray) !== -1) {
          let types_customs = {
            name: '报关单类型',
            filter_type: 's',
            id: ['#entry_type_name'],
            after: ['#entry_type'],
          }
          await admin.base_clearance_data_auto(types_customs)
        }

        if ($.inArray('inspection_quarantine', loadArray) !== -1) {
          let inspection_quarantine = {
            name: '检验检疫机关代码',
            filter_type: 's',
            id: [
              '#org_code_name',
              '#vsa_org_code_name',
              '#insp_org_name',
              '#purp_org_name'],
            after: [
              '#org_code',
              '#vsa_org_code',
              '#insp_org_code',
              '#purp_org_code'
            ],
          }
          await admin.base_clearance_data_auto(inspection_quarantine)
        }

        if ($.inArray('related_reasons', loadArray) !== -1) {
          let related_reasons = {
            name: '关联理由代码',
            filter_type: 's',
            id: ['#correlation_reason_flag_name'],
            after: ['#correlation_reason_flag'],
          }
          await admin.base_clearance_data_auto(related_reasons)
        }

        if ($.inArray('unit_measurement', loadArray) !== -1) {
          let unit_measurement = {
            name: '计量单位代码',
            filter_type: 's',
            id: [
              '#g_unit_name',
              '#first_unit_name',
              '#second_unit_name',

            ],
            after: [
              '#g_unit',
              '#first_unit',
              '#second_unit',
            ],
          }
          await admin.base_clearance_data_auto(unit_measurement)
        }

        if ($.inArray('ann_unit_measurement', loadArray) !== -1) {
          let ann_unit_measurement = {
            name: '计量单位代码',
            filter_type: 's',
            id: [
              '#DclUnitcdName',
              '#LawfUnitcdName',
              '#SecdLawfUnitcdName',
            ],
            after: [
              '#DclUnitcd',
              '#LawfUnitcd',
              '#SecdLawfUnitcd',
            ],
          }
          await admin.base_clearance_data_auto(ann_unit_measurement)
        }

        if ($.inArray('origin_area', loadArray) !== -1) {
          let origin_area = {
            name: '原产地区代码',
            filter_type: 's',
            id: ['#orig_place_code_name'],
            after: ['#orig_place_code'],
          }
          await admin.base_clearance_data_auto(origin_area)
        }

        if ($.inArray('domestic_area', loadArray) !== -1) {
          let domestic_area = {
            name: '国内地区代码',
            filter_type: 's',
            id: ['#district_code_name'],
            after: ['#district_code'],
          }
          await admin.base_clearance_data_auto(domestic_area)
        }

        if ($.inArray('destination', loadArray) !== -1) {
          let destination = {
            name: '中华人民共和国行政区划代码',
            filter_type: 's',
            id: ['#dest_code_name'],
            after: ['#dest_code'],
          }
          await admin.base_clearance_data_auto(destination)
        }

        if ($.inArray('exempting_method', loadArray) !== -1) {
          let exempting_method = {
            name: '征减免税方式代码',
            filter_type: 's',
            id: ['#duty_mode_name'],
            after: ['#duty_mode'],
          }
          await admin.base_clearance_data_auto(exempting_method)
        }

        if ($.inArray('ann_exempting_method', loadArray) !== -1) {
          let ann_exempting_method = {
            name: '征减免税方式代码',
            filter_type: 's',
            id: ['#LvyrlfModecdName'],
            after: ['#LvyrlfModecd'],
          }
          await admin.base_clearance_data_auto(ann_exempting_method)
        }

        if ($.inArray('use', loadArray) !== -1) {
          let use = {
            name: '用途代码',
            filter_type: 's',
            id: ['#purpose_name'],
            after: ['#purpose'],
          }
          await admin.base_clearance_data_auto(use)
        }

        if ($.inArray('type_container', loadArray) !== -1) {
          let type_container = {
            name: '集装箱规格代码',
            filter_type: 's',
            id: ['#container_md_name'],
            after: ['#container_md'],
          }
          await admin.base_clearance_data_auto(type_container)
        }

        if ($.inArray('documents_attached', loadArray) !== -1) {
          let documents_attached = {
            name: '随附单证代码',
            filter_type: 's',
            id: ['#docu_code_name'],
            after: ['#docu_code'],
          }
          await admin.base_clearance_data_auto(documents_attached)
        }

        if ($.inArray('enterprise_product', loadArray) !== -1) {
          let enterprise_product = {
            name: '企业产品许可类别代码',
            filter_type: 's',
            id: ['#ent_qualif_type_name_tem'],
            after: ['#ent_qualif_type_code_tem'],
          }
          await admin.base_clearance_data_auto(enterprise_product)
        }

        if ($.inArray('site_code', loadArray) !== -1) {
          let site_code = {
            name: '场地代码',
            filter_type: 's',
            id: ['#cus_fie_name'],
            after: ['#cus_fie'],
          }
          await admin.base_clearance_data_auto(site_code)
        }

        if ($.inArray('two_account_manual', loadArray) !== -1) {
          /**自动完成--手帐册**/
          await admin.auto_fn({
            url: `/handbook/datagrid`,
            listDirection: false,
            filter: function (data, data_filter) {
              for (let item of data) {
                data_filter.push({
                  id: item.Id,
                  label: `${item.ContractNumber}-${item.CompanyClientName}`,
                  value: item.ContractNumber,
                })
              }
            },
            id: ['#putrec_no'],
            after: ['#hand_book_id'],
          })
        }
      },

//数组上移、下移
      swapItems (arr, index1, index2) {
        arr[index1] = arr.splice(index2, 1, arr[index1])[0]
        return arr
      }
      ,

//进口报关打印列表
      order_i_print_list: [
        {
          id: 1,
          code: '00000004',
          name: '进口订购合同',
          is_enclosure: false,
        }, {
          id: 3,
          code: '00000001',
          name: '进口发票形式',
          is_enclosure: false,
        }, {
          id: 6,
          code: 'local_6',
          name: '进口六联司机',
          is_enclosure: false,
        }, {
          id: 7,
          code: '00000002',
          name: '进口货物装箱单',
          is_enclosure: false,
        }, {
          id: 12,
          code: 'local_9',
          name: '进口货物装箱单形式',
          is_enclosure: false,
        }, {
          id: 9,
          code: 'local_5',
          name: '关检合一新版进口报关单',
          is_enclosure: false,
        }, {
          id: 14,
          code: 'local_11',
          name: '审结通知书',
          is_enclosure: false,
        }, {
          id: 13,
          code: 'local_10',
          name: '放行通知书',
          is_enclosure: false,
        }],

      //出口报关打印列表
      order_e_print_list:
        [
          {
            id: 1,
            code: '00000004',
            name: '出口销售合同',
            is_enclosure: false,
          }, {
          id: 3,
          code: '00000001',
          name: '出口发票形式',
          is_enclosure: false,
        }, {
          id: 6,
          code: 'local_6',
          name: '出口六联司机',
          is_enclosure: false,
        }, {
          id: 7,
          code: '00000002',
          name: '出口货物装箱单',
          is_enclosure: false,
        }, {
          id: 12,
          code: 'local_9',
          name: '出口货物装箱单形式',
          is_enclosure: false,
        }, {
          id: 9,
          code: 'local_5',
          name: '关检合一新版出口报关单',
          is_enclosure: false,
        }, {
          id: 14,
          code: 'local_11',
          name: '审结通知书',
          is_enclosure: false,
        }, {
          id: 13,
          code: 'local_10',
          name: '放行通知书',
          is_enclosure: false,
        }],

      //特殊业务标识反填比对数据
      spec_decl_flag_data:
        ['国际赛事', '特殊进出军工物资', '国际援助物资', '国际会议', '直通放行', '外交礼遇', '转关'],

      //获取附件数据赋值打印
      async getPdf (order_id, i_e_flag, print_list) {
        const data = await admin.get(
          `/order/${i_e_flag}/${order_id}/pdf/lists?skip_paginate=true`)
        for (let item_print of print_list) {
          let data_judge = [],
            item_edoc_id = null
          for (let item_edoc of data) {
            if (item_edoc.edoc_code === item_print.code) {
              data_judge.push(1)
              item_edoc_id = item_edoc.id
            } else {
              data_judge.push(0)
            }
          }
          if (data_judge.includes(1)) {
            item_print.is_enclosure = true
            item_print.pdf_id = item_edoc_id
          } else {
            item_print.is_enclosure = false
          }
        }
        layer.closeAll('loading')
      }
      ,
// 重载PDF列表表格
      async render_pdf_list (order_id, i_e_flag) {
        if (i_e_flag == 'i') {
          print_list = admin.order_i_print_list
        } else {
          print_list = admin.order_e_print_list
        }
        await admin.getPdf(order_id, i_e_flag, print_list)
        layui.table.reload('print_lists', {
          data: print_list,
        })
      }
      ,

//判断一个字符串是否为数字
      isNumber (val) {
        var regPos = /^\d+(\.\d+)?$/ //非负浮点数
        var regNeg = /^(-(([0-9]+\.[0-9]*[1-9][0-9]*)|([0-9]*[1-9][0-9]*\.[0-9]+)|([0-9]*[1-9][0-9]*)))$/ //负浮点数
        if (regPos.test(val) || regNeg.test(val)) {
          return true
        } else {
          return false
        }
      }
      ,

//当前日期
      getCurrDate () {
        let now = new Date()
        let year = now.getFullYear() //年
        let month = now.getMonth() + 1 //月
        let day = now.getDate() //日
        if (month < 10) {
          month = '0' + month
        }
        if (day < 10) {
          day = '0' + day
        }
        return year + '' + month + '' + day
      }
      ,

//日期格式化
      getyyyymmdd (item) {
        let now = new Date(item)
        let year = now.getFullYear() //年
        let month = now.getMonth() + 1 //月
        let day = now.getDate() //日
        if (month < 10) {
          month = '0' + month
        }
        if (day < 10) {
          day = '0' + day
        }
        return year + '' + month + '' + day
      }
      ,

//图片base64 转 blob
      dataURItoBlob (dataURI) {
        // convert base64/URLEncoded data component to raw binary data held in a string
        let byteString
        if (dataURI.split(',')[0].indexOf('base64') >= 0) {
          byteString = atob(dataURI.split(',')[1])
        } else byteString = unescape(dataURI.split(',')[1])

        // separate out the mime component
        const mimeString = dataURI.split(',')[0].split(':')[1].split(';')[0]

        // write the bytes of the string to a typed array
        const ia = new Uint8Array(byteString.length)
        for (let i = 0; i < byteString.length; i++) {
          ia[i] = byteString.charCodeAt(i)
        }
        return new Blob([ia], {
          type: mimeString,
        })
      }
      ,
      data_item (index, item) {
        const jsonData = JSON.stringify(item)
        console.log(jsonData)
        return `<a class="seel_flex_edit_btn" data-index="${index}" data-item="${jsonData}">编辑</a>`
      }
      ,

      DecListGNoCheck (data, decData, InvtList) {
        let gNo = $.trim(data.entry_gds_seqno)
        if (gNo != '') {
          if (gNo <= decData.length) {
            let _decRow = decData[gNo - 1]
            let putrecSeqnos = _decRow.putrec_seqno.toString().
              replace(/\@/g, '').
              split(',')
            let bool
            //20181214 modefied by zongjie 修复老数据最终目的国字段没有值无法归并的问题
            if ('1' != $('#invt_type').val() &&
              ('' == _decRow.destination_natcd || undefined ==
                _decRow.destination_natcd)) {
              bool = data.gdecd == _decRow.gdecd && data.gds_nm ==
                _decRow.gds_nm &&
                data.dcl_unitcd == _decRow.dcl_unitcd &&
                data.dcl_currcd == _decRow.dcl_currcd &&
                data.natcd == _decRow.natcd &&
                data.entry_gds_seqno == _decRow.entry_gds_seqno
            } else if ('1' == $('#invt_type').val() &&
              ('' == _decRow.destination_natcd || undefined ==
                _decRow.destination_natcd)) {
              //20181213 modefied by zongjie
              bool = data.gdecd == _decRow.gdecd && data.gds_nm ==
                _decRow.gds_nm &&
                data.dcl_unitcd == _decRow.dcl_unitcd &&
                data.dcl_currcd == _decRow.dcl_currcd &&
                data.entry_gds_seqno == _decRow.entry_gds_seqno
            } else {
              bool = data.gdecd == _decRow.gdecd && data.gds_nm ==
                _decRow.gds_nm &&
                data.dcl_unitcd == _decRow.dcl_unitcd &&
                data.natcd == _decRow.natcd &&
                //20181018 优化 新增目的国
                data.destination_natcd == _decRow.destination_natcd &&
                data.dcl_currcd == _decRow.dcl_currcd &&
                data.entry_gds_seqno == _decRow.entry_gds_seqno
            }
            if (bool) {
              return true
            } else {
              layer.open({
                title: '提示',
                content: '表体序号[' + data.gds_seqno +
                  ']：报关单序号已经存在，但是与对应的报关商品信息归并条件不同，不能进行归并',
                icon: 2,
              })
              return false
            }
          } else if (gNo == (decData.length + 1)) {
            return true
          } else {
            layer.open({
              title: '提示',
              content: '报关单序号必须按照自然数顺序排号，不允许断号',
              icon: 2,
            })
            return false
          }
        }
        return true
      }
      ,

      DecListCountCheck (invtListData, invtDecListData) {
        //小于50项可以操作数据
        if (invtDecListData.length < 50) {
          return true
        } else { //等于50项时
          //判断新增或修改的数据是否可进行归并
          for (let i = 0; i < invtDecListData.length; i++) {
            let decList = invtDecListData[i]
            if (invtListData.gdecd === decList.gdecd &&
              invtListData.gds_nm === decList.gds_nm &&
              invtListData.dcl_unitcd === decList.dcl_unitcd &&
              invtListData.natcd === decList.natcd &&
              //20181018 优化 新增目的国
              invtListData.destination_natcd ===
              decList.destination_natcd &&
              invtListData.dcl_currcd === decList.dcl_currcd) {
              if (invtListData.entry_gds_seqno &&
                invtListData.entry_gds_seqno !== decList.entry_gds_seqno) {
                if (i == 49) {
                  return false
                }
                continue
              } else {
                return true
              }
            }
          }
          return false
        }
      }
      ,

      InvtList2DecList (invtList) {
        let decList = {}
        for (let item in invtList) {
          decList[item] = invtList[item]
        }
        decList.putrec_seqno = '@' + invtList.gds_seqno //报关单备案序号里面存储核注商品序号（为了生成报关单序号时返填核注商品信息）
        return decList
      }
      ,

      combineDec (decList, tmpDecList) {
        if (!!decList.dcl_qty) {
          let __t = new BigDecimal(decList.dcl_qty.toString()).add(
            new BigDecimal(tmpDecList.dcl_qty.toString() === ''
              ? '0'
              : tmpDecList.dcl_qty.toString())).
            setScale(5, MathContext.ROUND_HALF_UP).
            toString()
          tmpDecList.dcl_qty = admin.cutZero(__t)
        }
        if (!!decList.dcl_total_amt) {
          let __t = new BigDecimal(decList.dcl_total_amt.toString()).add(
            new BigDecimal(
              tmpDecList.dcl_total_amt.toString() === ''
                ? '0'
                : tmpDecList.dcl_total_amt.toString())).
            setScale(2, MathContext.ROUND_HALF_UP).
            toString()
          tmpDecList.dcl_total_amt = admin.cutZero(__t)
        }
        if (!!decList.lawf_qty) {
          let __t = new BigDecimal(decList.lawf_qty.toString()).add(
            new BigDecimal(
              tmpDecList.lawf_qty.toString() === ''
                ? '0'
                : tmpDecList.lawf_qty.toString())).
            setScale(5, MathContext.ROUND_HALF_UP).
            toString()
          tmpDecList.lawf_qty = admin.cutZero(__t)
        }
        if (!!decList.secd_lawf_qty) {
          let __t = new BigDecimal(decList.secd_lawf_qty.toString()).add(
            new BigDecimal(
              tmpDecList.secd_lawf_qty.toString() === ''
                ? '0'
                : tmpDecList.secd_lawf_qty.toString())).
            setScale(5, MathContext.ROUND_HALF_UP).
            toString()
          tmpDecList.secd_lawf_qty = admin.cutZero(__t)
        }
        if (!!tmpDecList.dcl_qty && tmpDecList.dcl_qty.toString() !== '0' &&
          !!tmpDecList.dcl_total_amt) {
          let __t = new BigDecimal(tmpDecList.dcl_total_amt).divide(
            new BigDecimal(tmpDecList.dcl_qty.toString()), 5,
            MathContext.ROUND_HALF_UP).
            setScale(4, MathContext.ROUND_HALF_UP).
            toString()
          tmpDecList.dcl_uprc_amt = admin.cutZero(__t)
        }

        let _tmpPutrecSeqNos = tmpDecList.putrec_seqno.toString().split(',')
        _tmpPutrecSeqNos.push(decList.putrec_seqno)
        tmpDecList.putrec_seqno = _tmpPutrecSeqNos.join(',')
      }
      ,

//去除末尾多余的零
      cutZero (old) {
        let newstr = old
        let leng = old.length - old.indexOf('.') - 1
        // 无小数点不处理
        if (old.indexOf('.') > -1) {
          // 循环小数部分
          for (let i = leng; i > 0; i--) {
            // 如果newstr末尾有0
            if (newstr.lastIndexOf('0') > -1 &&
              newstr.substr(newstr.length - 1, 1) === 0) {
              let k = newstr.lastIndexOf('0')
              // 如果小数点后只有一个0 去掉小数点
              if (newstr.charAt(k - 1) === '.') {
                return newstr.substring(0, k - 1)
              } else {
                // 否则 去掉一个0
                newstr = newstr.substring(0, k)
              }
            } else {
              // 如果末尾没有0
              return newstr
            }
          }
        }
        return old
      }
      ,

      /**如果报关则生成报关商品**/
      InitInvtList2DecList (decListData, order_pros_data) {
        let entryGdsSeqNos = []
        decListData = []
        if ($('#dclcus_flag').val() === '1') {
          if (order_pros_data.length > 0) {
            for (let i in order_pros_data) {
              let l = order_pros_data[i]
              let _tmpDec = admin.InvtList2DecList(l)
              if (decListData.length === 0) {
                decListData.push(_tmpDec)
                entryGdsSeqNos.push(decListData.length)
              } else {
                let flag = false
                for (let j = 0; j < decListData.length; j++) {
                  let _tmp = decListData[j]
                  /**商品编码+商品名称+计量单位+币制+产销国**/
                  if (l.gdecd == _tmp.gdecd &&
                    l.gds_nm == _tmp.gds_nm &&
                    l.dcl_unitcd == _tmp.dcl_unitcd &&
                    l.natcd == _tmp.natcd &&
                    l.destination_natcd == _tmp.destination_natcd &&
                    l.dcl_currcd == _tmp.dcl_currcd) {
                    if (l.entry_gds_seqno && l.entry_gds_seqno !=
                      _tmp.entry_gds_seqno) {
                      if (entryGdsSeqNos.filter(function (v) {
                        return l.entry_gds_seqno == v
                      }).length > 0) {
                        flag = false
                        continue
                      }
                      decListData.push(_tmpDec)
                      entryGdsSeqNos.push(decListData.length)
                      flag = true
                      break
                    } else {
                      admin.combineDec(_tmpDec, _tmp)
                      flag = true
                      break
                    }
                  } else {
                    flag = false
                  }
                }
                if (!flag) {
                  decListData.push(_tmpDec)
                  entryGdsSeqNos.push(decListData.length)
                }
              }
            }
          }
        }
        return decListData
      }
      ,

      /**清单产品/报关草稿列表**/
      getAnnotationTable (annotationItems, decListData) {
        /**产品**/
        layui.table.render({
          elem: '#order_pros',
          toolbar: '#order_pros_tool',
          defaultToolbar: ['filter'],
          colFilterRecord: 'local',
          primaryKey: 'GdsSeqno',
          cols: [
            [
              {
                type: 'radio',
              }, {
              field: 'GdsSeqno',
              title: '序号',
              width: 100,
            }, {
              field: 'PutrecSeqno',
              title: '备案序号',
              width: 120,
            }, {
              field: 'GdsMtno',
              title: '商品料号',
              width: 140,
            }, {
              field: 'EntryGdsSeqno',
              title: '报关单商品序号',
              width: 180,
            }, {
              field: 'ApplyTbSeqno',
              title: '流转申报表序号',
              width: 180,
            }, {
              field: 'Gdecd',
              title: '商品编码',
              width: 140,
            }, {
              field: 'GdsNm',
              title: '商品名称',
              width: 240,
            }, {
              field: 'GdsSpcfModelDesc',
              title: '规格型号',
              width: 240,
            }, {
              field: 'DclUnitcdName',
              title: '申报计量单位',
              width: 160,
            }, {
              field: 'LawfUnitcdName',
              title: '法定计量单位',
              width: 160,
            }, {
              field: 'SecdLawfUnitcdName',
              title: '法定第二计量单位',
              width: 180,
            }, {
              field: 'NatcdName',
              title: '原产国(地区)',
              width: 150,
            }, {
              field: 'DestinationNatcdName',
              title: '最终目的国(地区)',
              width: 200,
            }, {
              field: 'DclUprcAmt',
              title: '企业申报单价',
              width: 160,
            }, {
              field: 'DclCurrcdName',
              title: '币制',
              width: 100,
            }, {
              field: 'ModfMarkcdName',
              title: '修改标志',
              width: 100,
            }, {
              field: 'LawfQty',
              title: '法定数量',
              width: 120,
            }, {
              field: 'SecdLawfQty',
              title: '第二法定数量',
              width: 160,
            }, {
              field: 'WtSfVal',
              title: '重量比例因子',
              width: 160,
            }, {
              field: 'FstSfVal',
              title: '第一比例因子',
              width: 160,
            }, {
              field: 'SecdSfVal',
              title: '第二比例因子',
              width: 160,
            }, {
              field: 'DclQty',
              title: '申报数量',
              width: 120,
            }, {
              field: 'GrossWt',
              title: '毛重',
              width: 100,
            }, {
              field: 'NetWt',
              title: '净重',
              width: 100,
            }, {
              field: 'LvyrlfModecdName',
              title: '征免方式',
              width: 120,
            }, {
              field: 'UcnsVerno',
              title: '单耗版本号',
              width: 150,
            }, {
              field: 'ClyMarkcd',
              title: '归类标志',
              width: 130,
            }, {
              field: 'Param3',
              title: '自动备案序号',
              width: 160,
            }],
          ],
          data: annotationItems.rows,
          limit: annotationItems.total,
          height: 250,
        })

        /**报关单草稿**/
        layui.table.render({
          elem: '#customs_declaration_list',
          toolbar: true,
          defaultToolbar: ['filter'],
          colFilterRecord: 'local',
          cols: [
            [
              {
                field: 'DecSeqNo',
                title: '报关单统一编号',
                width: 180,
              }, {
              field: 'EntryGdsSeqno',
              title: '报关单商品序号',
              width: 180,
            }, {
              field: 'Gdecd',
              title: '商品编码',
              width: 140,
            }, {
              field: 'GdsNm',
              title: '商品名称',
              width: 240,
            }, {
              field: 'GdsSpcfModelDesc',
              title: '规格型号',
              width: 240,
            }, {
              field: 'DclQtyString',
              title: '申报数量',
              width: 120,
            }, {
              field: 'DclUprcAmt',
              title: '申报单价',
              width: 140,
            }, {
              field: 'DclTotalAmt',
              title: '申报总价',
              width: 140,
            }, {
              field: 'DclCurrcdName',
              title: '申报币制',
              width: 120,
            }, {
              field: 'DclUnitcdName',
              title: '申报计量单位',
              width: 160,
            }, {
              field: 'LawfUnitcdName',
              title: '法定计量单位',
              width: 160,
            }, {
              field: 'SecdLawfUnitcdName',
              title: '法定第二计量单位',
              width: 180,
            }, {
              field: 'NatcdName',
              title: '原产国(地区)',
              width: 150,
            }, {
              field: 'DestinationNatcdName',
              title: '最终目的国(地区)',
              width: 200,
            }, {
              field: 'LawfQty',
              title: '法定数量',
              width: 140,
            }, {
              field: 'SecdLawfQty',
              title: '第二法定数量',
              width: 160,
            }],
          ],
          data: decListData,
          limit: decListData.length,
          height: 200,
        })
      }
      ,

      /** 货物申报表体表格渲染 */
      getOrderTable (order_pros_data) {
        let destination_name = '原产国(地区)'
        let origin_name = '最终目的国'

        if (admin.cusIEFlag == 'I') {
          destination_name = '最终目的国'
          origin_name = '原产国(地区)'
        }

        layui.table.render({
          elem: '#order_pros',
          toolbar: '#order_pros_tool',
          defaultToolbar: ['filter'],
          colFilterRecord: 'local',
          primaryKey: 'g_no',
          cols: [
            [
              {
                type: 'checkbox',
              }, {
              field: 'g_no',
              title: '序号',
              width: 80,
            }, {
              field: 'contr_item',
              title: '备案序号',
              width: 100,
            }, {
              field: 'code_t_s',
              title: '商品编号',
              width: 120,
            }, {
              field: 'ciq_name',
              title: '检验检疫名称',
              width: 160,
            }, {
              field: 'g_name',
              title: '商品名称',
              width: 180,
            }, {
              field: 'g_model',
              title: '规格',
              width: 200,
            }, {
              field: 'g_qty_string',
              title: '成交数量',
              width: 100,
            }, {
              field: 'g_unit_name',
              title: '成交单位',
              width: 100,
            }, {
              field: 'decl_price_string',
              title: '单价',
              width: 100,
            }, {
              field: 'decl_total_string',
              title: '总价',
              width: 120,
            }, {
              field: 'trade_curr_name',
              title: '币制',
              width: 100,
            }, {
              field: 'destination_country_name',
              title: destination_name,
              width: 120,
            }, {
              field: 'origin_country_name',
              title: origin_name,
              width: 120,
            }, {
              field: 'duty_mode_name',
              title: '征免类型',
              width: 150,
            }, {
              field: 'cus_supv_dmd',
              title: '监管要求',
              width: 150,
            }],
          ],
          data: order_pros_data,
          limit: order_pros_data.length,
          height: 250,
        })
      }
      ,

      /** 货物申报司机纸打印弹窗表格渲染 */
      getOrderDriverTable (driver_data, order_i_edit_data, elem) {
        const pros_data = [...driver_data.first_pros, ...driver_data.pros]
        let order_driver_data_index = 1
        for (let pros_item of pros_data) {
          if (pros_item.g_no == '1') {
            pros_item.pack_no = order_i_edit_data.pack_no
            pros_item.net_wt = order_i_edit_data.net_wt
            pros_item.gross_wet = order_i_edit_data.gross_wet
            if (order_driver_data_index === 1) {
              pros_item.foreign_company_name = order_i_edit_data.foreign_company_name
              pros_item.trade_name = order_i_edit_data.trade_name
            }
          }

          order_driver_data_index += 1
        }
        console.log(order_driver_data_index)
        layui.table.render({
          elem: elem,
          toolbar: true,
          defaultToolbar: ['filter'],
          colFilterRecord: 'local',
          cols: [
            [
              {
                field: 'g_no',
                title: '货物序号',
                width: 90,
              }, {
              field: 'g_name',
              title: '商品名称',
              width: 120,
            }, {
              field: 'g_model',
              title: '规格型号',
              width: 150,
            }, {
              field: 'pack_no',
              title: '货物包装数量',
              width: 120,
            }, {
              field: 'net_wt',
              title: '货物净重',
              width: 100,
            }, {
              field: 'gross_wet',
              title: '货物毛重',
              width: 100,
            }, {
              field: 'decl_total_string',
              title: '货物总价',
              width: 100,
            }, {
              field: 'foreign_company_name',
              title: '货物外商名称',
              width: 140,
            }, {
              field: 'trade_name',
              title: '货物收发名称',
              width: 140,
            }, {
              field: 'default',
              title: '特殊包装数量种类',
              width: 160,
            }],
          ],
          data: pros_data,
          limit: pros_data.length,
          height: 350,
        })
      }
      ,

//货物申报--备注双击弹出
      note_s_index: '',
      note_s_dbclick () {
        admin.note_s_index = layer.open({
          type: 1,
          title: false,
          skin: 'layui-layer-admin',
          closeBtn: false,
          area: '350px',
          anim: 5,
          shadeClose: true,
          content: $('#note_s_selset').html(),
        })
      }
      ,
//货物申报--备注双击弹出选择
      note_s_selset_p (dom) {
        $('#note_s_fu').val($(dom).text())
        $('#note_s').val($(dom).text())
        $('#cus_fie').val($(dom).data('code'))
        $('#cus_fie_name').val($(dom).data('name'))
        layer.close(admin.note_s_index)
      }
      ,

//排序数字从小到大规则
      compare (prop) {
        return (obj1, obj2) => {
          let val1 = obj1[prop]
          let val2 = obj2[prop]
          if (!isNaN(Number(val1)) && !isNaN(Number(val2))) {
            val1 = Number(val1)
            val2 = Number(val2)
          }
          if (val1 < val2) {
            return -1
          } else if (val1 > val2) {
            return 1
          } else {
            return 0
          }
        }
      }
      ,

//小数点自动补零
      formatnumber (value, num) {
        let a, b, c, i
        a = value.toString()
        b = a.indexOf('.')
        c = a.length
        if (num == 0) {
          if (b != -1) {
            a = a.substring(0, b)
          }
        } else { //如果没有小数点
          if (b == -1) {
            a = a + '.'
            for (i = 1; i <= num; i++) {
              a = a + '0'
            }
          } else { //有小数点，超出位数自动截取，否则补0
            a = a.substring(0, b + num + 1)
            for (i = c; i <= b + num; i++) {
              a = a + '0'
            }
          }
        }
        return a
      }
      ,

      //货物申报-办理记录
      async order_jilu_click (dom, type) {
        const id = $(dom).data('id')
        layer.load(2)
        const data_order = await admin.get(`/order/${type}/${id}/order_logs`)
        layer.closeAll('loading')
        const data = {
          data_order: data_order,
        }
        layui.laytpl($('#order_take_template').html()).
          render(data, function (html) {
            $('#order_take_list').html(html)
          })
        layer.open({
          type: 1,
          title: '办理记录',
          shadeClose: true,
          area: admin.screen() < 2 ? ['80%', '300px'] : ['600px', '500px'],
          content: `<div id="order_take_list_content">${$('#order_take_list').
            html()}</div>`,
        })
      }
      ,

      //核注清单-办理记录
      async annotation_jilu_click (dom) {
        const id = $(dom).data('id')
        if (!id) {
          return layer.msg('请先保存订单！')
        }
        layer.load(2)
        const annotationRecords = await admin.post(
          `/annotation_record/datagrid`,
          JSON.stringify({ AnnotationId: parseInt(id) }), true)
        layer.closeAll('loading')
        let rows = annotationRecords.rows ? annotationRecords.rows : []
        laytpl($('#order_take_template').html()).render(rows, function (html) {
          $('#order_take_list').html(html)
        })
        layer.open({
          type: 1,
          title: '办理记录',
          shadeClose: true,
          area: admin.screen() < 2 ? ['80%', '300px'] : ['600px', '500px'],
          content: `<div id="order_take_list_content">${$('#order_take_list').
            html()}</div>`,
        })
      },

//补录报关单
      add_entry: {
        id: '',
        flag:
          '',
      }
      ,
      add_entry_click (dom) {
        admin.add_entry.id = $(dom).data('id')
        admin.add_entry.flag = $(dom).data('flag')
        layer.open({
          type: 1,
          title: '补录报关单',
          shadeClose: true,
          area: admin.screen() < 2 ? ['80%', '300px'] : ['600px', 'auto'],
          content: $('#order_add_entry_template').html(),
        })
      }
      ,

//提交补录报关单
      order_add_entry_submit () {
        layui.form.on('submit(order_add_entry_submit)', async (form) => {
          try {
            const data = await admin.post(
              `/order/${admin.add_entry.flag}/${admin.add_entry.id}/add_entry_id`,
              form.field)
            if (data.status === 1) {
              setTimeout(() => {
                layer.closeAll()
              }, 2000)
            }
          } catch (e) {
            return layer.msg('接口错误！', {
              offset: '15px',
              icon: 2,
              time: 2000,
              id: 'Message',
            })
          }
        })
      }
      ,

//发送验证码
      sendAuthCode: function (options) {
        options = $.extend({
          seconds: 60,
          elemPhone: '#LAY_phone',
          elemVercode: '#LAY_vercode',
        }, options)

        var seconds = options.seconds,
          btn = $(options.elem),
          token = null,
          timer, countDown = function (loop) {
            seconds--
            if (seconds < 0) {
              btn.removeClass(DISABLED).html('获取验证码')
              seconds = options.seconds
              clearInterval(timer)
            } else {
              btn.addClass(DISABLED).html(seconds + '秒后重获')
            }

            if (!loop) {
              timer = setInterval(function () {
                countDown(true)
              }, 1000)
            }
          }

        options.elemPhone = $(options.elemPhone)
        options.elemVercode = $(options.elemVercode)

        btn.on('click', function () {
          var elemPhone = options.elemPhone,
            value = elemPhone.val()

          if (seconds !== options.seconds ||
            $(this).hasClass(DISABLED)) return

          if (!/^1\d{10}$/.test(value)) {
            elemPhone.focus()
            return layer.msg('请输入正确的手机号')
          }

          if (typeof options.ajax === 'object') {
            var success = options.ajax.success
            delete options.ajax.success
          }

          admin.req($.extend(true, {
            url: '/auth/code',
            type: 'post',
            data: {
              phone: value,
            },
            success: function (res) {
              layer.msg('验证码已发送至你的手机，请注意查收')
              options.elemVercode.focus()
              countDown()
              success && success(res)
            },
          }, options.ajax))
        })
      }
      ,

//屏幕根据分辨率等比例缩小--收缩侧边栏
      sideFlexible_window () {
        const s = (window.screen.width - 270) / 1920
        document.body.style.zoom = s
        parent.document.body.style.zoom = s
        if (window.screen.width != 1920) {
          parent.layui.admin.sideFlexible()
        }

      }
      ,

//订单列表
      list_page: 1,
      list_limit: 10,
      async get_data_list (OrderIndexRequestData, isClickStatusTab) {

        let url = OrderIndexRequestData.List.Url
        let impexpMarkcd = OrderIndexRequestData.ImpexpMarkcd
        let StatusString = OrderIndexRequestData.StatusString

        /**订单列表**/
        let OrderIndexRequestListData = JSON.stringify($.extend(
          OrderIndexRequestData.List.Request, {
            ImpexpMarkcd: impexpMarkcd,
            StatusString: StatusString,
            offset: admin.list_page,
            limit: admin.list_limit,
          },
        ))
        let ListDatas = await admin.post(url, OrderIndexRequestListData,
          true)

        //点击状态 tab 不触发
        if (!isClickStatusTab) {
          //列表数量
          let StatusCount = await admin.post(
            OrderIndexRequestData.StatusCount.Url,
            JSON.stringify($.extend(
              OrderIndexRequestData.StatusCount.Request, {
                StatusString: StatusString,
                ImpexpMarkcd: impexpMarkcd,
              },
            )), true,
          )

          /**订单状态数量**/
          layui.laytpl($('#status_flex_list_template').html()).
            render(StatusCount.rows, function (html) {
              $('#status_flex_list').html(html)
            })
        }

        $('#order-i-table tbody').remove()
        if (ListDatas.total === 0) {
          $('#order-i-table').append(
            `<tbody><tr class="sep-row"><td colspan="5"><div class="no_data">无数据</div></td></tr></tbody>`)
        } else {
          layui.laytpl($('#order_i_list').html()).
            render(ListDatas.rows, function (html) {
              $('#order-i-table').append(html)
            })
        }

        layui.form.render('select')

        /**订单列表分页**/
        layui.laypage.render({
          elem: 'order_page',
          count: ListDatas.total,
          limit: admin.list_limit,
          limits: [10, 20, 30, 40, 50, 100, 200],
          theme: '#1E9FFF',
          layout: ['count', 'prev', 'page', 'next', 'limit', 'skip'],
          jump: async function (obj, first) {
            if (!first) {
              admin.list_page = obj.curr
              admin.list_limit = obj.limit
              OrderIndexRequestListData = JSON.stringify($.extend(
                OrderIndexRequestData.List.Request, {
                  ImpexpMarkcd: impexpMarkcd,
                  StatusString: StatusString,
                  offset: admin.list_page,
                  limit: admin.list_limit,
                },
              ))
              ListDatas = await admin.post(`${url}`,
                OrderIndexRequestListData,
                true)
              $('#order-i-table tbody').remove()
              if (ListDatas.total === 0) {
                $('#order-i-table').append(
                  `<tbody><tr class="sep-row"><td colspan="5"><div class="no_data">无数据</div></td></tr></tbody>`)
              } else {
                layui.laytpl($('#order_i_list').html()).
                  render(ListDatas.rows, function (html) {
                    $('#order-i-table').append(html)
                  })
              }
            }

          },
        })

        return ListDatas
      }
      ,

//只允许数字
      is_onlynumber (dom) {
        $(dom).val($(dom).val().replace(/\D/g, ''))
      }
      ,

//只能输入数字，小数点，不能有空格
      is_nolyNorD (dom) {
        $(dom).val($(dom).val().replace(/[^0-9\.\/]/g, ''))
      }
      ,

//不允许中文和空格
      is_noCork (dom) {
        $(dom).val($(dom).val().replace(/[\u4E00-\u9FA5\s]/g, ''))
      }
      ,

//只允许数字和-
      is_onlynumberLine (dom) {
        $(dom).val($(dom).val().replace(/[^\d-]/g, ''))
      }
      ,

//只能输入小数点后两位的数字
      is_onlyNumFloat (dom, number) {
        let value = $(dom).val()
        value = value.replace(/[^\d.]/g, '') //清理"数字"和"."以外的字符
        value = value.replace(/^\./g, '') //验证第一个字符是数字
        value = value.replace(/\.{2,}/g, '') //只保留第一个, 清理多余的
        value = value.replace('.', '$#$').
          replace(/\./g, '').
          replace('$#$', '.')
        if (number == 'two') {
          value = value.replace(/^(\-)*(\d+)\.(\d\d).*$/, '$1$2.$3') //只能输入两个小数
        }
        if (number == 'four') {
          value = value.replace(/^(\-)*(\d+)\.(\d\d\d\d).*$/, '$1$2.$3') //只能输入四个小数
        }
        if (number == 'sixteen') {
          value = value.replace(
            /^(\-)*(\d+)\.(\d\d\d\d\d\d\d\d\d\d\d\d\d\d\d\d).*$/,
            '$1$2.$3') //只能输入十六个小数
        }
        $(dom).val(value)
      }
      ,

//监听进出口货物申报业务选项
      promise_items_change (dom) {
        const first = $(dom).data('first')
        if (first == 0) {
          $(dom).data('first', 1)
        }
        if ($(dom).prop('checked')) {
          $(dom).prev('span').text('是')
        } else {
          $(dom).prev('span').text('否')
        }
      }
      ,

//获取cookie
      getCookie (name) {
        let arr, reg = new RegExp('(^| )' + name + '=([^;]*)(;|$)')
        if (arr = document.cookie.match(reg)) {
          return unescape(arr[2])
        } else {
          return false
        }
      }
      ,

//异地报关
      is_other_change (dom) {
        const is_other = $(dom).prop('checked') ? 1 : 0
        if (is_other) {
          $('#agent_code_scc').
            val('').
            removeAttr('lay-verify').
            removeAttr('lay-vertype').
            removeClass('required').
            attr('disabled', 'disabled')
          $('#agent_code').
            val('').
            removeAttr('lay-verify').
            removeAttr('lay-vertype').
            removeClass('required').
            attr('disabled', 'disabled')
          $('#decl_ciq_code').val('').attr('disabled', 'disabled')
          $('#agent_name').
            val('').
            removeAttr('lay-verify').
            removeAttr('lay-vertype').
            removeClass('required').
            attr('disabled', 'disabled')
        } else {
          $('#agent_code_scc').
            val('914419007962277606').
            attr('lay-verify', 'required').
            attr('lay-vertype', 'tips').
            addClass('required').
            removeAttr('disabled')
          $('#agent_code').
            val('4419986507').
            attr('lay-verify', 'required').
            attr('lay-vertype', 'tips').
            addClass('required').
            removeAttr('disabled')
          $('#decl_ciq_code').val('4400510052').removeAttr('disabled')
          $('#agent_name').
            val('广东东华报关服务有限公司').
            attr('lay-verify', 'required').
            attr('lay-vertype', 'tips').
            addClass('required').
            removeAttr('disabled')
        }
      }
      ,

//货物申报添加提运单号
      bill_no_create () {
        layer.open({
          type: 2,
          title: '提运单号',
          shadeClose: true,
          area: admin.screen() < 2 ? ['80%', '300px'] : ['880px', '700px'],
          content: `/car/create`,
        })
      }
      ,

//同步关务通--填写平台单证号
      sync_annotation_to_order (dom) {
        layer.open({
          type: 1,
          title: '同步关务通',
          shadeClose: true,
          area: admin.screen() < 2 ? ['80%', '300px'] : ['680px', 'auto'],
          content: `<div class="manual_no_list">
                    <form class="layui-form" onkeydown="if(event.keyCode===13){return false}">
                        <p>请向客户获取关务通的核注清单订单号，格式为：D012345678912345</p>
                        <table class="layui-table order_table_form">
                            <tbody>
                            <tr>
                                <td class="td-scale-01 text-align-right">平台单证号</td>
                                <td class="td-scale-02 td-btn">
                                    <div class="input-group valid">
                                        <input type="text" class="non-empty required" name="docNo" id="docNo" autocomplete="off"
                                            lay-verify="required" lay-vertype="tips">
                                    </div>
                                </td>
                            </tr>
                            </tbody>
                        </table>
                        <div class="manual_no_btn">
                            <button class="layui-btn custom-create_btn" lay-submit lay-filter="sync_annotation_to_order_save"
                                    id="sync_annotation_to_order_save" type="button" data-flag="${$(
            dom).data('flag')}"
                                    onclick="layui.admin.sync_annotation_to_order_save(this)">
                                提交
                            </button>
                        </div>
                    </form>
                </div>`,
        })
        $('#docNo').focus()
      }
      ,
//同步关务通--填写平台单证号-保存
      sync_annotation_to_order_save (dom) {
        layui.form.on('submit(sync_annotation_to_order_save)',
          async (data) => {
            try {
              const res = admin.post(
                `/order/${$(dom).data('flag')}/sync_annotation_to_order`,
                data.field)
              if (res.status) {
                setTimeout(() => {
                  layer.closeAll()
                }, 2000)
              }
            } catch (error) {
              return layer.msg('接口错误！', {
                offset: '15px',
                icon: 2,
                time: 2000,
                id: 'Message',
              })
            }
          })
      }
      ,

//屏幕类型
      screen: function () {
        var width = $win.width()
        if (width > 1200) {
          return 3 //大屏幕
        } else if (width > 992) {
          return 2 //中屏幕
        } else if (width > 768) {
          return 1 //小屏幕
        } else {
          return 0 //超小屏幕
        }
      }
      ,

//侧边伸缩
      sideFlexible: function (status) {
        var app = container,
          iconElem = $('#' + APP_FLEXIBLE),
          screen = admin.screen()

        //设置状态，PC：默认展开、移动：默认收缩
        if (status === 'spread') {
          //切换到展开状态的 icon，箭头：←
          iconElem.removeClass(ICON_SPREAD).addClass(ICON_SHRINK)

          //移动：从左到右位移；PC：清除多余选择器恢复默认
          if (screen < 2) {
            app.addClass(APP_SPREAD_SM)
          } else {
            app.removeClass(APP_SPREAD_SM)
          }

          app.removeClass(SIDE_SHRINK)
        } else {
          //切换到搜索状态的 icon，箭头：→
          iconElem.removeClass(ICON_SHRINK).addClass(ICON_SPREAD)

          //移动：清除多余选择器恢复默认；PC：从右往左收缩
          if (screen < 2) {
            app.removeClass(SIDE_SHRINK)
          } else {
            app.addClass(SIDE_SHRINK)
          }

          app.removeClass(APP_SPREAD_SM)
        }

        layui.event.call(this, setter.MOD_NAME, 'side({*})', {
          status: status,
        })
      }
      ,

//弹出面板
      popup: view.popup,

      //右侧面板
      popupRight:

        function (options) {
          //layer.close(admin.popup.index);
          return admin.popup.index = layer.open($.extend({
            type: 1,
            id: 'LAY_adminPopupR',
            anim: -1,
            title: false,
            closeBtn: false,
            offset: 'r',
            shade: 0.1,
            shadeClose: true,
            skin: 'layui-anim layui-anim-rl layui-layer-adminRight',
            area: '300px',
          }, options))
        }

      ,

//主题设置
      theme: function (options) {
        var theme = setter.theme,
          local = layui.data(setter.tableName),
          id = 'LAY_layadmin_theme',
          style = document.createElement('style'),
          styleText = laytpl([
            //主题色
            '.layui-side-menu,',
            '.layadmin-pagetabs .layui-tab-title li:after,',
            '.layadmin-pagetabs .layui-tab-title li.layui-this:after,',
            '.layui-layer-admin .layui-layer-title,',
            '.layadmin-side-shrink .layui-side-menu .layui-nav>.layui-nav-item>.layui-nav-child',
            '{background-color:{{d.color.main}} !important;}'

            //选中色
            ,
            '.layui-nav-tree .layui-this,',
            '.layui-nav-tree .layui-this>a,',
            '.layui-nav-tree .layui-nav-child dd.layui-this,',
            '.layui-form-select dl dd.layui-this,',
            '.layui-nav-tree .layui-nav-child dd.layui-this a',
            '{background-color:{{d.color.selected}} !important;}'

            ,
            '.layui-form-radio>i:hover, .layui-form-radioed>i',
            '{color:{{d.color.selected}} !important;}'

            //logo
            ,
            '.layui-layout-admin .layui-logo{background-color:{{d.color.logo || d.color.main}} !important;}'

            //头部色
            ,
            '{{# if(d.color.header){ }}',
            '.layui-layout-admin .layui-header{background-color:{{ d.color.header }};}',
            '.layui-layout-admin .layui-header a,',
            '.layui-layout-admin .layui-header a cite{color: #f8f8f8;}',
            '.layui-layout-admin .layui-header a:hover{color: #fff;}',
            '.layui-layout-admin .layui-header .layui-nav .layui-nav-more{border-top-color: #fbfbfb;}',
            '.layui-layout-admin .layui-header .layui-nav .layui-nav-mored{border-color: transparent; border-bottom-color: #fbfbfb;}',
            '.layui-layout-admin .layui-header .layui-nav .layui-this:after, .layui-layout-admin .layui-header .layui-nav-bar{background-color: #fff; background-color: rgba(255,255,255,.5);}',
            '.layadmin-pagetabs .layui-tab-title li:after{display: none;}',
            '{{# } }}',
          ].join('')).render(options = $.extend({}, local.theme, options)),
          styleElem = document.getElementById(id)

        //添加主题样式
        if ('styleSheet' in style) {
          style.setAttribute('type', 'text/css')
          style.styleSheet.cssText = styleText
        } else {
          style.innerHTML = styleText
        }
        style.id = id

        styleElem && $body[0].removeChild(styleElem)
        $body[0].appendChild(style)
        $body.attr('layadmin-themealias', options.color.alias)

        //本地存储记录
        local.theme = local.theme || {}
        layui.each(options, function (key, value) {
          local.theme[key] = value
        })
        layui.data(setter.tableName, {
          key: 'theme',
          value: local.theme,
        })
      }
      ,

//初始化主题
      initTheme: function (index) {
        var theme = setter.theme
        index = index || 0
        if (theme.color[index]) {
          theme.color[index].index = index
          admin.theme({
            color: theme.color[index],
          })
        }
      }
      ,

//记录最近一次点击的页面标签数据
      tabsPage: {}
      ,

//获取页面标签主体元素
      tabsBody: function (index) {
        return $(APP_BODY).find('.' + TABS_BODY).eq(index || 0)
      }
      ,

      //切换页面标签主体
      tabsBodyChange: function (index, options) {
        options = options || {}

        admin.tabsBody(index).addClass(SHOW).siblings().removeClass(SHOW)
        events.rollPage('auto', index)

        //执行 {setter.MOD_NAME}.tabsPage 下的事件
        layui.event.call(this, setter.MOD_NAME, 'tabsPage({*})', {
          url: options.url,
          text: options.text,
        })
      }
      ,

      //resize事件管理
      resize: function (fn) {
        var router = layui.router(),
          key = router.path.join('-')

        if (admin.resizeFn[key]) {
          $win.off('resize', admin.resizeFn[key])
          delete admin.resizeFn[key]
        }

        if (fn === 'off') return //如果是清除 resize 事件，则终止往下执行

        fn(), admin.resizeFn[key] = fn
        $win.on('resize', admin.resizeFn[key])
      }
      ,
      resizeFn: {}
      ,
      runResize: function () {
        var router = layui.router(),
          key = router.path.join('-')
        admin.resizeFn[key] && admin.resizeFn[key]()
      }
      ,
      delResize: function () {
        this.resize('off')
      }
      ,

      //关闭当前 pageTabs
      closeThisTabs: function () {
        if (!admin.tabsPage.index) return
        $(TABS_HEADER).
          eq(admin.tabsPage.index).
          find('.layui-tab-close').
          trigger('click')
      }
      ,

      //获取当前iframe的标签
      get_iframe_index () {
        if (!admin.tabsPage.index) return
        return $(TABS_HEADER).eq(admin.tabsPage.index)
      }
      ,

      //全屏
      fullScreen: function () {
        var ele = document.documentElement,
          reqFullScreen = ele.requestFullScreen ||
            ele.webkitRequestFullScreen ||
            ele.mozRequestFullScreen || ele.msRequestFullscreen
        if (typeof reqFullScreen !== 'undefined' && reqFullScreen) {
          reqFullScreen.call(ele)
        }
      }
      ,

//退出全屏
      exitScreen: function () {
        var ele = document.documentElement
        if (document.exitFullscreen) {
          document.exitFullscreen()
        } else if (document.mozCancelFullScreen) {
          document.mozCancelFullScreen()
        } else if (document.webkitCancelFullScreen) {
          document.webkitCancelFullScreen()
        } else if (document.msExitFullscreen) {
          document.msExitFullscreen()
        }
      },

//……
    }

//事件
  var events = admin.events = {
    //伸缩
    flexible: function (othis) {
      var iconElem = othis.find('#' + APP_FLEXIBLE),
        isSpread = iconElem.hasClass(ICON_SPREAD)
      admin.sideFlexible(isSpread ? 'spread' : null)
    }

    //刷新
    ,
    refresh: function () {
      var ELEM_IFRAME = '.layadmin-iframe',
        length = $('.' + TABS_BODY).length

      if (admin.tabsPage.index >= length) {
        admin.tabsPage.index = length - 1
      }

      var iframe = admin.tabsBody(admin.tabsPage.index).find(ELEM_IFRAME)
      iframe[0].contentWindow.location.reload(true)
    }

    //输入框搜索
    ,
    serach: function (othis) {
      othis.off('keypress').on('keypress', function (e) {
        if (!this.value.replace(/\s/g, '')) return
        //回车跳转
        if (e.keyCode === 13) {
          var href = othis.attr('lay-action'),
            text = othis.attr('lay-text') || '搜索'

          href = href + this.value
          text = text + ' <span style="color: #FF5722;">' +
            admin.escape(this.value) + '</span>'

          //打开标签页
          layui.index.openTabsPage(href, text)

          //如果搜索关键词已经打开，则刷新页面即可
          events.serach.keys || (events.serach.keys = {})
          events.serach.keys[admin.tabsPage.index] = this.value
          if (this.value === events.serach.keys[admin.tabsPage.index]) {
            events.refresh(othis)
          }

          //清空输入框
          this.value = ''
        }
      })
    }

    //点击消息
    ,
    message: function (othis) {
      othis.find('.layui-badge-dot').remove()
    }

    //弹出主题面板
    ,
    theme: function () {
      admin.popupRight({
        id: 'LAY_adminPopupTheme',
        success: function () {
          view(this.id).render('system/theme')
        },
      })
    }

    //便签
    ,
    note: function (othis) {
      var mobile = admin.screen() < 2,
        note = layui.data(setter.tableName).note

      events.note.index = admin.popup({
        title: '便签',
        shade: 0,
        offset: [
          '41px', (mobile ? null : (othis.offset().left - 250) + 'px'),
        ],
        anim: -1,
        id: 'LAY_adminNote',
        skin: 'layadmin-note layui-anim layui-anim-upbit',
        content: '<textarea placeholder="内容"></textarea>',
        resize: false,
        success: function (layero, index) {
          var textarea = layero.find('textarea'),
            value = note === undefined
              ? '便签中的内容会存储在本地，这样即便你关掉了浏览器，在下次打开时，依然会读取到上一次的记录。是个非常小巧实用的本地备忘录'
              : note

          textarea.val(value).focus().on('keyup', function () {
            layui.data(setter.tableName, {
              key: 'note',
              value: this.value,
            })
          })
        },
      })
    }

    //全屏
    ,
    fullscreen: function (othis) {
      var SCREEN_FULL = 'layui-icon-screen-full',
        SCREEN_REST = 'layui-icon-screen-restore',
        iconElem = othis.children('i')

      if (iconElem.hasClass(SCREEN_FULL)) {
        admin.fullScreen()
        iconElem.addClass(SCREEN_REST).removeClass(SCREEN_FULL)
      } else {
        admin.exitScreen()
        iconElem.addClass(SCREEN_FULL).removeClass(SCREEN_REST)
      }
    }

    //弹出关于面板
    ,
    about: function () {
      admin.popupRight({
        id: 'LAY_adminPopupAbout',
        success: function () {
          view(this.id).render('system/about')
        },
      })
    }

    //弹出更多面板
    ,
    more: function () {
      admin.popupRight({
        id: 'LAY_adminPopupMore',
        success: function () {
          view(this.id).render('system/more')
        },
      })
    }

    //返回上一页
    ,
    back: function () {
      history.back()
    }

    //主题设置
    ,
    setTheme: function (othis) {
      var index = othis.data('index'),
        nextIndex = othis.siblings('.layui-this').data('index')

      if (othis.hasClass(THIS)) return

      othis.addClass(THIS).siblings('.layui-this').removeClass(THIS)
      admin.initTheme(index)
    }

    //左右滚动页面标签
    ,
    rollPage: function (type, index) {
      var tabsHeader = $('#LAY_app_tabsheader'),
        liItem = tabsHeader.children('li'),
        scrollWidth = tabsHeader.prop('scrollWidth'),
        outerWidth = tabsHeader.outerWidth(),
        tabsLeft = parseFloat(tabsHeader.css('left'))

      //右左往右
      if (type === 'left') {
        if (!tabsLeft && tabsLeft <= 0) return

        //当前的left减去可视宽度，用于与上一轮的页标比较
        var prefLeft = -tabsLeft - outerWidth

        liItem.each(function (index, item) {
          var li = $(item),
            left = li.position().left

          if (left >= prefLeft) {
            tabsHeader.css('left', -left)
            return false
          }
        })
      } else if (type === 'auto') { //自动滚动
        (function () {
          var thisLi = liItem.eq(index),
            thisLeft

          if (!thisLi[0]) return
          thisLeft = thisLi.position().left

          //当目标标签在可视区域左侧时
          if (thisLeft < -tabsLeft) {
            return tabsHeader.css('left', -thisLeft)
          }

          //当目标标签在可视区域右侧时
          if (thisLeft + thisLi.outerWidth() >= outerWidth - tabsLeft) {
            var subLeft = thisLeft + thisLi.outerWidth() -
              (outerWidth - tabsLeft)
            liItem.each(function (i, item) {
              var li = $(item),
                left = li.position().left

              //从当前可视区域的最左第二个节点遍历，如果减去最左节点的差 > 目标在右侧不可见的宽度，则将该节点放置可视区域最左
              if (left + tabsLeft > 0) {
                if (left - tabsLeft > subLeft) {
                  tabsHeader.css('left', -left)
                  return false
                }
              }
            })
          }
        }())
      } else {
        //默认向左滚动
        liItem.each(function (i, item) {
          var li = $(item),
            left = li.position().left

          if (left + li.outerWidth() >= outerWidth - tabsLeft) {
            tabsHeader.css('left', -left)
            return false
          }
        })
      }
    }

    //向右滚动页面标签
    ,
    leftPage: function () {
      events.rollPage('left')
    }

    //向左滚动页面标签
    ,
    rightPage: function () {
      events.rollPage()
    }

    //关闭当前标签页
    ,
    closeThisTabs: function () {
      var topAdmin = parent === self ? admin : parent.layui.admin
      topAdmin.closeThisTabs()
    }

    //关闭其它标签页
    ,
    closeOtherTabs: function (type) {
      var TABS_REMOVE = 'LAY-system-pagetabs-remove'
      if (type === 'all') {
        $(TABS_HEADER + ':gt(0)').remove()
        $(APP_BODY).find('.' + TABS_BODY + ':gt(0)').remove()

        $(TABS_HEADER).eq(0).trigger('click')
      } else {
        $(TABS_HEADER).each(function (index, item) {
          if (index && index != admin.tabsPage.index) {
            $(item).addClass(TABS_REMOVE)
            admin.tabsBody(index).addClass(TABS_REMOVE)
          }
        })
        $('.' + TABS_REMOVE).remove()
      }
    }

    //关闭全部标签页
    ,
    closeAllTabs: function () {
      events.closeOtherTabs('all')
      //location.hash = '';
    }

    //遮罩
    ,
    shade: function () {
      admin.sideFlexible()
    }

    //呼出IM 示例
    ,
    im: function () {
      admin.popup({
        id: 'LAY-popup-layim-demo' //定义唯一ID，防止重复弹出
        ,
        shade: 0,
        area: ['800px', '300px'],
        title: '面板外的操作示例',
        offset: 'lb',
        success: function () {
          //将 views 目录下的某视图文件内容渲染给该面板
          layui.view(this.id).render('layim/demo').then(function () {
            layui.use('im')
          })
        },
      })
    },
  }

//初始
  !function () {
    //主题初始化，本地主题记录优先，其次为 initColorIndex
    var local = layui.data(setter.tableName)
    if (local.theme) {
      admin.theme(local.theme)
    } else if (setter.theme) {
      admin.initTheme(setter.theme.initColorIndex)
    }

    //常规版默认开启多标签页
    if (!('pageTabs' in layui.setter)) layui.setter.pageTabs = true

    //不开启页面标签时
    if (!setter.pageTabs) {
      $('#LAY_app_tabs').addClass(HIDE)
      container.addClass('layadmin-tabspage-none')
    }

    //低版本IE提示
    if (device.ie && device.ie < 10) {
      view.error(
        'IE' + device.ie + '下访问可能不佳，推荐使用：Chrome / Firefox / Edge 等高级浏览器', {
          offset: 'auto',
          id: 'LAY_errorIE',
        })
    }

  }()

//admin.prevRouter = {}; //上一个路由

//监听 tab 组件切换，同步 index
  element.on('tab(' + FILTER_TAB_TBAS + ')', function (data) {
    admin.tabsPage.index = data.index
  })

//监听选项卡切换，改变菜单状态
  admin.on('tabsPage(setMenustatus)', function (router) {
    var pathURL = router.url,
      getData = function (item) {
        return {
          list: item.children('.layui-nav-child'),
          a: item.children('*[lay-href]'),
        }
      },
      sideMenu = $('#' + SIDE_MENU),
      SIDE_NAV_ITEMD = 'layui-nav-itemed'

      //捕获对应菜单
      ,
      matchMenu = function (list) {
        list.each(function (index1, item1) {
          var othis1 = $(item1),
            data1 = getData(othis1),
            listChildren1 = data1.list.children('dd'),
            matched1 = pathURL === data1.a.attr('lay-href')

          listChildren1.each(function (index2, item2) {
            var othis2 = $(item2),
              data2 = getData(othis2),
              listChildren2 = data2.list.children('dd'),
              matched2 = pathURL === data2.a.attr('lay-href')

            listChildren2.each(function (index3, item3) {
              var othis3 = $(item3),
                data3 = getData(othis3),
                matched3 = pathURL === data3.a.attr('lay-href')

              if (matched3) {
                var selected = data3.list[0] ? SIDE_NAV_ITEMD : THIS
                othis3.addClass(selected).siblings().removeClass(selected) //标记选择器
                return false
              }

            })

            if (matched2) {
              var selected = data2.list[0] ? SIDE_NAV_ITEMD : THIS
              othis2.addClass(selected).siblings().removeClass(selected) //标记选择器
              return false
            }

          })

          if (matched1) {
            var selected = data1.list[0] ? SIDE_NAV_ITEMD : THIS
            othis1.addClass(selected).siblings().removeClass(selected) //标记选择器
            return false
          }

        })
      }

    //重置状态
    sideMenu.find('.' + THIS).removeClass(THIS)

    //移动端点击菜单时自动收缩
    if (admin.screen() < 2) admin.sideFlexible()

    //开始捕获
    matchMenu(sideMenu.children('li'))
  })

//监听侧边导航点击事件
  element.on('nav(layadmin-system-side-menu)', function (elem) {
    if (elem.siblings('.layui-nav-child')[0] &&
      container.hasClass(SIDE_SHRINK)) {
      admin.sideFlexible('spread')
      layer.close(elem.data('index'))
    }

    admin.tabsPage.type = 'nav'
  })

//监听选项卡的更多操作
  element.on('nav(layadmin-pagetabs-nav)', function (elem) {
    var dd = elem.parent()
    dd.removeClass(THIS)
    dd.parent().removeClass(SHOW)
  })

//同步路由
  var setThisRouter = function (othis) {
      var layid = othis.attr('lay-id'),
        attr = othis.attr('lay-attr'),
        index = othis.index()

      admin.tabsBodyChange(index, {
        url: attr,
      })
      //location.hash = layid === setter.entry ? '/' : attr;
    },
    TABS_HEADER = '#LAY_app_tabsheader>li'

//标签页标题点击
  $body.on('click', TABS_HEADER, function () {
    var othis = $(this),
      index = othis.index()

    admin.tabsPage.type = 'tab'
    admin.tabsPage.index = index

    setThisRouter(othis)
  })

//监听 tabspage 删除
  element.on('tabDelete(' + FILTER_TAB_TBAS + ')', function (obj) {
    var othis = $(TABS_HEADER + '.layui-this')

    obj.index && admin.tabsBody(obj.index).remove()
    setThisRouter(othis)

    //移除resize事件
    admin.delResize()
  })

//页面跳转
  $body.on('click', '*[lay-href]', function () {
    var othis = $(this),
      href = othis.attr('lay-href'),
      text = othis.attr('lay-text'),
      router = layui.router()

    admin.tabsPage.elem = othis
    //admin.prevRouter[router.path[0]] = router.href; //记录上一次各菜单的路由信息

    //执行跳转
    var topLayui = parent === self ? layui : top.layui
    topLayui.index.openTabsPage(href, text || othis.text())
  })

//点击事件
  $body.on('click', '*[layadmin-event]', function () {
    var othis = $(this),
      attrEvent = othis.attr('layadmin-event')
    events[attrEvent] && events[attrEvent].call(this, othis)
  })

//tips
  $body.on('mouseenter', '*[lay-tips]', function () {
    var othis = $(this)

    if (othis.parent().hasClass('layui-nav-item') &&
      !container.hasClass(SIDE_SHRINK)) return

    var tips = othis.attr('lay-tips'),
      offset = othis.attr('lay-offset'),
      direction = othis.attr('lay-direction'),
      index = layer.tips(tips, this, {
        tips: direction || 1,
        time: -1,
        success: function (layero, index) {
          if (offset) {
            layero.css('margin-left', offset + 'px')
          }
        },
      })
    othis.data('index', index)
  }).on('mouseleave', '*[lay-tips]', function () {
    layer.close($(this).data('index'))
  })

//窗口resize事件
  var resizeSystem = layui.data.resizeSystem = function () {
    //layer.close(events.note.index);
    layer.closeAll('tips')

    if (!resizeSystem.lock) {
      setTimeout(function () {
        admin.sideFlexible(admin.screen() < 2 ? '' : 'spread')
        delete resizeSystem.lock
      }, 100)
    }

    resizeSystem.lock = true
  }

  $win.on('resize', layui.data.resizeSystem)

//接口输出
  exports('admin', admin)
})


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
                parent.document.getElementById(frameId).contentWindow.location.reload()
            },
            UPLOAD_PDF_SIZE: 10, //pdf上传大小 M
            v: '1.2.1 std',
            //数据的异步请求
            req: view.req,
            //清除本地 token，并跳转到登入页
            exit: view.exit,
            //xss 转义
            escape: function (html) {
                return String(html || '').replace(/&(?!#?[a-zA-Z0-9]+;)/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/'/g, '&#39;').replace(/"/g, '&quot;')
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
                            checkCell = tableDiv.find('tr[data-index=' + index + ']').find(
                                'td div.laytable-cell-checkbox div.layui-form-checkbox i'),
                            radioCell = tableDiv.find('tr[data-index=' + index + ']').find(
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
            //只允许输入正整数
            decCheckInt(obj) {
                let t = obj.value.replace(/[^(\(\)\d\&\|)]/g, '')
                if (obj.value != t)
                    obj.value = t
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
                            focusable = form.find('input,a,select,textarea').filter(':visible').not(':input[readonly]').not(':input[disabled]')
                            prev = focusable.eq(focusable.index(this) - 1)
                            if (prev.length) {
                                if ($(this).attr('shiftEnter') == 'no') {
                                    return false
                                } else {
                                    prev.focus()
                                }
                            }
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

                        focusable = form.find('input,select,textarea').filter(
                            ':visible').not(':input[readonly]').not(
                            ':input[disabled]').not(':input[enter=-1]')
                        console.log()
                        next = focusable.eq(focusable.index(this) + 1)
                        // 下个元素存在
                        if (next.length) {
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



            //自动完成
            async auto_fn(type) {
                let data_filter = [];
                let requestData = JSON.stringify({Limit: 5000, TypeString: type.clearanceType});
                let data = await admin.post(type.url, requestData, true);
                type.filter(data.rows, data_filter);
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
                                type.after.forEach((avalue, aindex) => {
                                    $(type.after[aindex]).val(data.id[aindex]);
                                });
                                if (type.after[index] === '#TransMode') {
                                    admin.transModeControl(admin.cusIEFlag)
                                }
                                if (type.after[index] === '#FeeMark') {
                                    admin.markSelect('FeeMark', 'FeeCurr', 'FeeCurrName')
                                }
                                if (type.after[index] === '#InsurMark') {
                                    admin.markSelect('InsurMark', 'InsurCurr', 'insur curr name')
                                }
                                if (type.after[index] === '#OtherMark') {
                                    admin.markSelect('OtherMark', 'OtherCurr', 'OtherCurrName')
                                }
                                if (type.after[index] === '#TrafMode') {
                                    if ($('#TrafMode').val() === 4) {
                                        //$("#bill_no").removeAttr("disabled", "disabled");
                                        //启运国(地区)
                                        $('#TradeCountry').val('HKG');
                                        $('#TradeCountryName').val('中国香港');
                                        //经停港
                                        $('#DistinatePort').val('HKG003');
                                        $('#DistinatePortName').val('香港（中国香港）');
                                        //贸易国别（地区）
                                        $('#TradeAreaCode').val('HKG');
                                        $('#TradeAreaName').val('中国香港');
                                        //启运港
                                        $('#DespPortCode').val('HKG003');
                                        $('#DespPortName').val('香港（中国香港）')
                                    } else {
                                        //$("#bill_no").attr("disabled", "disabled");
                                    }
                                }
                                if (type.after[index] == '#TrspModecd') {
                                    if ($('#TrspModecd').val() === 4) {
                                        $('#StshipTrsarvNatcd').val('110');
                                        $('#StshipTrsarvNatcdName').val('中国香港')
                                    }
                                }

                                if (type.after[index] === '#CusFie') {
                                    const value = $(type.after[index]).val();
                                    if (value === '5284') {
                                        $('#NoteS').val('[装卸口岸：长安车检场]')
                                    }
                                    if (value === '5299') {
                                        $('#NoteS').val('[装卸口岸：其它业务]')
                                    }
                                    if (value === '5238') {
                                        $('#NoteS').val('[装卸口岸：凤岗车检场]')
                                    }
                                    if (value === '5298') {
                                        $('#NoteS').val('[装卸口岸：外关区]')
                                    }
                                    if (value === '5297') {
                                        $('#NoteS').val('[装卸口岸：加贸结转]')
                                    }
                                }
                            }
                        },
                    })
                })
            },
            //毛重判断
            gross_wet_blur(dom) {
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
            net_wt_blur(dom) {
                if ($(dom).val().trim()) {
                    if (parseFloat($(dom).val().trim()) >
                        parseFloat($('#GrossWet').val().trim())) {
                        $(dom).focus()
                        return layer.msg('净重大于毛重，请确认后重新填写!')
                    }
                }
            },

            /**
             * 四舍六入五成双
             * @param num
             * @param digit   小数点多少位
             * calculationType：计算类型（加减乘除对应0,1,2,3）
             * @returns {Number}
             */
            decToDecimal(num1, num2, digit, calculationType, roundingMode) {
                let calculationDatas = [num1 + ',' + num2]
                let digits = [digit]
                let calculationTypes = [calculationType]
                let roundingModes = [roundingMode]
                let resultList = admin.decCalculation(calculationDatas, digits,
                    calculationTypes, roundingModes, '1')
                return resultList[0]
            },

            /**
             *
             * @param calculations    计算值的集合
             * @param calculationType    计算类型,0,1,2,3分别对应加减乘除
             * @param roundingMode    计算结果小数保留的方式，四舍五入等
             * @param digit    保留小数的位数
             * @param isZeroNoShow    小数点后全为0是否显示，1显示，0不显示
             */
            decCalculation(
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
            calculation(
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
                        bigDecimalResult = bigDecimalResult.multiply(bigDecimal).setScale(parseInt(digit), parseInt(roundingMode))
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
            },
            //判断一个字符串是否为数字
            isNumber(val) {
                var regPos = /^\d+(\.\d+)?$/ //非负浮点数
                var regNeg = /^(-(([0-9]+\.[0-9]*[1-9][0-9]*)|([0-9]*[1-9][0-9]*\.[0-9]+)|([0-9]*[1-9][0-9]*)))$/ //负浮点数
                if (regPos.test(val) || regNeg.test(val)) {
                    return true
                } else {
                    return false
                }
            },
            //当前日期
            getCurrDate() {
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
            },

            //日期格式化
            getyyyymmdd(item) {
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
            },

            //图片base64 转 blob
            dataURItoBlob(dataURI) {
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
            },
            data_item(index, item) {
                const jsonData = JSON.stringify(item)
                console.log(jsonData)
                return `<a class="seel_flex_edit_btn" data-index="${index}" data-item="${jsonData}">编辑</a>`
            },

            //去除末尾多余的零
            cutZero(old) {
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
            },
            //排序数字从小到大规则
            compare(prop) {
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
            },

            //小数点自动补零
            formatnumber(value, num) {
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
            },

            //屏幕根据分辨率等比例缩小--收缩侧边栏
            sideFlexible_window() {
                const s = (window.screen.width - 270) / 1920
                document.body.style.zoom = s
                parent.document.body.style.zoom = s
                if (window.screen.width != 1920) {
                    parent.layui.admin.sideFlexible()
                }
            },

            //只允许数字
            is_onlynumber(dom) {
                $(dom).val($(dom).val().replace(/\D/g, ''))
            },

            //只能输入数字，小数点，不能有空格
            is_nolyNorD(dom) {
                $(dom).val($(dom).val().replace(/[^0-9\.\/]/g, ''))
            },

            //不允许中文和空格
            is_noCork(dom) {
                $(dom).val($(dom).val().replace(/[\u4E00-\u9FA5\s]/g, ''))
            },

            //只允许数字和-
            is_onlynumberLine(dom) {
                $(dom).val($(dom).val().replace(/[^\d-]/g, ''))
            },

            //只能输入小数点后两位的数字
            is_onlyNumFloat(dom, number) {
                let value = $(dom).val()
                value = value.replace(/[^\d.]/g, '') //清理"数字"和"."以外的字符
                value = value.replace(/^\./g, '') //验证第一个字符是数字
                value = value.replace(/\.{2,}/g, '') //只保留第一个, 清理多余的
                value = value.replace('.', '$#$').replace(/\./g, '').replace('$#$', '.')
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
            },

            //获取cookie
            getCookie(name) {
                let arr, reg = new RegExp('(^| )' + name + '=([^;]*)(;|$)')
                if (arr = document.cookie.match(reg)) {
                    return unescape(arr[2])
                } else {
                    return false
                }
            },

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
            },

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
            },
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
                },
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
            },

            //记录最近一次点击的页面标签数据
            tabsPage: {},

            //获取页面标签主体元素
            tabsBody: function (index) {
                return $(APP_BODY).find('.' + TABS_BODY).eq(index || 0)
            },

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
            },

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
            },
            resizeFn: {},
            runResize: function () {
                var router = layui.router(),
                    key = router.path.join('-')
                admin.resizeFn[key] && admin.resizeFn[key]()
            },
            delResize: function () {
                this.resize('off')
            },
            //关闭当前 pageTabs
            closeThisTabs: function () {
                if (!admin.tabsPage.index) return
                $(TABS_HEADER).eq(admin.tabsPage.index).find('.layui-tab-close').trigger('click')
            },

            //获取当前iframe的标签
            get_iframe_index() {
                if (!admin.tabsPage.index) return
                return $(TABS_HEADER).eq(admin.tabsPage.index)
            },

            //全屏
            fullScreen: function () {
                var ele = document.documentElement,
                    reqFullScreen = ele.requestFullScreen ||
                        ele.webkitRequestFullScreen ||
                        ele.mozRequestFullScreen || ele.msRequestFullscreen
                if (typeof reqFullScreen !== 'undefined' && reqFullScreen) {
                    reqFullScreen.call(ele)
                }
            },

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
        }
        //事件
    var events = admin.events = {
        //伸缩
        flexible: function (othis) {
            var iconElem = othis.find('#' + APP_FLEXIBLE),
                isSpread = iconElem.hasClass(ICON_SPREAD)
            admin.sideFlexible(isSpread ? 'spread' : null)
        } ,
        //刷新
        refresh: function () {
            var ELEM_IFRAME = '.layadmin-iframe',
                length = $('.' + TABS_BODY).length
            if (admin.tabsPage.index >= length) {
                admin.tabsPage.index = length - 1
            }
            var iframe = admin.tabsBody(admin.tabsPage.index).find(ELEM_IFRAME)
            iframe[0].contentWindow.location.reload(true)
        },
        //输入框搜索
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
        } ,
        //点击消息
        message: function (othis) {
            othis.find('.layui-badge-dot').remove()
        },
        //弹出主题面板
        theme: function () {
            admin.popupRight({
                id: 'LAY_adminPopupTheme',
                success: function () {
                    view(this.id).render('system/theme')
                },
            })
        } ,
        //便签
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
        }  ,
        //全屏
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
        } ,
        //弹出关于面板
        about: function () {
            admin.popupRight({
                id: 'LAY_adminPopupAbout',
                success: function () {
                    view(this.id).render('system/about')
                },
            })
        } ,
        //弹出更多面板
        more: function () {
            admin.popupRight({
                id: 'LAY_adminPopupMore',
                success: function () {
                    view(this.id).render('system/more')
                },
            })
        },
        //返回上一页
         back: function () {
            history.back()
        }
        //主题设置
        , setTheme: function (othis) {
            var index = othis.data('index'),
                nextIndex = othis.siblings('.layui-this').data('index')
            if (othis.hasClass(THIS)) return

            othis.addClass(THIS).siblings('.layui-this').removeClass(THIS)
            admin.initTheme(index)
        },
        //左右滚动页面标签
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
        } ,
        //向右滚动页面标签
        leftPage: function () {
            events.rollPage('left')
        } ,
        //向左滚动页面标签
        rightPage: function () {
            events.rollPage()
        } ,
        //关闭当前标签页
        closeThisTabs: function () {
            var topAdmin = parent === self ? admin : parent.layui.admin
            topAdmin.closeThisTabs()
        },
        //关闭其它标签页
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
        } ,
        //关闭全部标签页
        closeAllTabs: function () {
            events.closeOtherTabs('all')
            //location.hash = '';
        } ,
        //遮罩
        shade: function () {
            admin.sideFlexible()
        },
        //呼出IM 示例,
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
        //执行跳转
        var topLayui = parent === self ? layui : top.layui;
        topLayui.index.openTabsPage(href, text || othis.text())
    })
//点击事件
    $body.on('click', '*[layadmin-event]', function () {
        var othis = $(this),
            attrEvent = othis.attr('layadmin-event');
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
            });
        othis.data('index', index)
    }).on('mouseleave', '*[lay-tips]', function () {
        layer.close($(this).data('index'))
    });

//窗口resize事件
    var resizeSystem = layui.data.resizeSystem = function () {
        //layer.close(events.note.index);
        layer.closeAll('tips');

        if (!resizeSystem.lock) {
            setTimeout(function () {
                admin.sideFlexible(admin.screen() < 2 ? '' : 'spread');
                delete resizeSystem.lock
            }, 100)
        }

        resizeSystem.lock = true
    };

    $win.on('resize', layui.data.resizeSystem);

//接口输出
    exports('admin', admin)
});

